#!/usr/bin/env python3
"""Detects schema changes in Microsoft Graph SDK updates and creates GitHub issues.

This script automates the detection of schema changes when the msgraph-beta-sdk-go
is updated by Dependabot. It:
- Parses go.mod to get current SDK version
- Validates version increment (ensures single version jump)
- Fetches and analyzes the SDK changelog
- Analyzes commit diffs in the models/ folder
- Detects struct field additions and removals
- Creates detailed GitHub issues for required schema updates

Usage:
    # Analyze a specific PR
    ./sdkSchemaChangeDetector.py --pr-number 1686

    # Analyze version change directly
    ./sdkSchemaChangeDetector.py --current v0.156.0 --new v0.157.0

    # Dry run (don't create issues)
    ./sdkSchemaChangeDetector.py --pr-number 1686 --dry-run

Args:
    --pr-number: Dependabot PR number to analyze
    --current: Current SDK version (optional, auto-detected from go.mod)
    --new: New SDK version (optional, auto-detected from PR)
    --dry-run: Analyze changes without creating GitHub issues
    --repo: Repository in owner/repo format (default: current repo)
"""

import sys
import json
import subprocess
import re
import argparse
from pathlib import Path
from typing import Optional, Dict, List, Tuple
from dataclasses import dataclass, field
from datetime import datetime


@dataclass
class FieldChange:
    """Represents a field change in a Go struct."""
    field_name: str
    field_type: str
    change_type: str  # 'added' or 'removed'
    line_number: Optional[int] = None


@dataclass
class ModelChange:
    """Represents changes to a Go model file."""
    file_path: str
    model_name: str
    added_fields: List[FieldChange] = field(default_factory=list)
    removed_fields: List[FieldChange] = field(default_factory=list)

    @property
    def has_changes(self) -> bool:
        """Check if this model has any changes."""
        return bool(self.added_fields or self.removed_fields)

    @property
    def change_summary(self) -> str:
        """Get a summary of changes."""
        parts = []
        if self.added_fields:
            parts.append(f"+{len(self.added_fields)} fields")
        if self.removed_fields:
            parts.append(f"-{len(self.removed_fields)} fields")
        return ", ".join(parts)


class SDKSchemaChangeDetector:
    """Detects and reports schema changes in Microsoft Graph SDK updates."""

    SDK_MODULE = "github.com/microsoftgraph/msgraph-beta-sdk-go"
    SDK_REPO = "microsoftgraph/msgraph-beta-sdk-go"

    def __init__(self, repo: Optional[str] = None, dry_run: bool = False):
        """Initialize the detector.

        Args:
            repo: Target repository in owner/repo format
            dry_run: If True, don't create actual GitHub issues
        """
        self.dry_run = dry_run
        self.repo = repo or self._get_current_repo()
        self.go_mod_path = Path.cwd() / "go.mod"

    def _get_current_repo(self) -> str:
        """Get current repository from git remote."""
        try:
            result = subprocess.run(
                ["git", "remote", "get-url", "origin"],
                capture_output=True,
                text=True,
                check=True
            )
            remote_url = result.stdout.strip()
            # Parse owner/repo from URL
            match = re.search(r'github\.com[:/](.+/.+?)(\.git)?$', remote_url)
            if match:
                return match.group(1).rstrip('.git')
        except subprocess.CalledProcessError:
            pass
        return "deploymenttheory/terraform-provider-microsoft365"

    def run_command(self, args: List[str], check: bool = True) -> Tuple[str, str, int]:
        """Run a command and return output.

        Args:
            args: Command and arguments to run
            check: Whether to raise exception on non-zero exit

        Returns:
            Tuple of (stdout, stderr, returncode)
        """
        try:
            result = subprocess.run(
                args,
                capture_output=True,
                text=True,
                check=check
            )
            return result.stdout.strip(), result.stderr.strip(), result.returncode
        except subprocess.CalledProcessError as e:
            if check:
                print(f"Error running command: {' '.join(args)}", file=sys.stderr)
                print(f"Error: {e.stderr}", file=sys.stderr)
                raise
            return e.stdout.strip(), e.stderr.strip(), e.returncode

    def run_gh_command(self, args: List[str]) -> str:
        """Run a GitHub CLI command and return output.

        Args:
            args: Arguments to pass to 'gh' command

        Returns:
            Command stdout as string
        """
        stdout, _, _ = self.run_command(["gh"] + args)
        return stdout

    def parse_go_mod_version(self) -> Optional[str]:
        """Parse current SDK version from go.mod.

        Returns:
            Current version string (e.g., 'v0.156.0') or None if not found
        """
        if not self.go_mod_path.exists():
            print(f"Error: go.mod not found at {self.go_mod_path}", file=sys.stderr)
            return None

        with open(self.go_mod_path, 'r') as f:
            for line in f:
                if self.SDK_MODULE in line:
                    # Extract version using regex
                    match = re.search(r'v\d+\.\d+\.\d+', line)
                    if match:
                        return match.group(0)

        print(f"Error: Could not find {self.SDK_MODULE} in go.mod", file=sys.stderr)
        return None

    def get_pr_version_change(self, pr_number: int) -> Optional[Tuple[str, str]]:
        """Get version change from a PR.

        Args:
            pr_number: PR number

        Returns:
            Tuple of (old_version, new_version) or None if not found
        """
        try:
            # Get PR diff
            diff = self.run_gh_command([
                "pr", "diff", str(pr_number),
                "--repo", self.repo
            ])

            # Look for go.mod changes
            old_version = None
            new_version = None

            for line in diff.split('\n'):
                if self.SDK_MODULE in line:
                    if line.startswith('-') and not line.startswith('---'):
                        # Removed line (old version)
                        match = re.search(r'v\d+\.\d+\.\d+', line)
                        if match:
                            old_version = match.group(0)
                    elif line.startswith('+') and not line.startswith('+++'):
                        # Added line (new version)
                        match = re.search(r'v\d+\.\d+\.\d+', line)
                        if match:
                            new_version = match.group(0)

            if old_version and new_version:
                return old_version, new_version

        except Exception as e:
            print(f"Error getting PR version change: {e}", file=sys.stderr)

        return None

    def validate_version_increment(self, old_version: str, new_version: str) -> bool:
        """Validate that version increment is exactly one minor version.

        Args:
            old_version: Old version string (e.g., 'v0.156.0')
            new_version: New version string (e.g., 'v0.157.0')

        Returns:
            True if increment is valid (single minor version bump)
        """
        # Parse versions
        old_match = re.match(r'v(\d+)\.(\d+)\.(\d+)', old_version)
        new_match = re.match(r'v(\d+)\.(\d+)\.(\d+)', new_version)

        if not old_match or not new_match:
            print(f"Error: Invalid version format", file=sys.stderr)
            return False

        old_major, old_minor, old_patch = map(int, old_match.groups())
        new_major, new_minor, new_patch = map(int, new_match.groups())

        # Check if it's a single minor version increment (patch should reset to 0)
        # Or a patch increment on the same minor version
        if new_major == old_major:
            if new_minor == old_minor + 1 and new_patch == 0:
                # Valid minor version bump
                return True
            elif new_minor == old_minor and new_patch == old_patch + 1:
                # Valid patch version bump
                return True
            elif new_minor == old_minor and new_patch > old_patch:
                # Multiple patch version bump - still acceptable
                print(f"⚠️  Warning: Multiple patch version increment detected", file=sys.stderr)
                return True
            elif new_minor > old_minor + 1:
                print(f"⚠️  Warning: Multiple minor version jump detected: {old_version} -> {new_version}", file=sys.stderr)
                print(f"   This may indicate missing intermediate versions.", file=sys.stderr)
                return False

        print(f"⚠️  Warning: Unexpected version change: {old_version} -> {new_version}", file=sys.stderr)
        return False

    def fetch_changelog_section(self, version: str) -> str:
        """Fetch changelog section for a specific version.

        Args:
            version: Version to fetch changelog for

        Returns:
            Changelog section text
        """
        try:
            # Fetch changelog from GitHub
            changelog_url = f"https://raw.githubusercontent.com/{self.SDK_REPO}/main/CHANGELOG.md"
            stdout, _, _ = self.run_command([
                "curl", "-s", changelog_url
            ])

            # Extract section for this version
            lines = stdout.split('\n')
            section_lines = []
            in_section = False

            for line in lines:
                if line.startswith('##') and version.lstrip('v') in line:
                    in_section = True
                    section_lines.append(line)
                elif in_section:
                    if line.startswith('##'):
                        # Next version section, stop
                        break
                    section_lines.append(line)

            return '\n'.join(section_lines) if section_lines else "Changelog section not found"

        except Exception as e:
            print(f"Error fetching changelog: {e}", file=sys.stderr)
            return "Error fetching changelog"

    def get_sdk_commit_for_version(self, version: str) -> Optional[str]:
        """Get the commit SHA for a specific SDK version.

        Args:
            version: Version tag (e.g., 'v0.157.0')

        Returns:
            Commit SHA or None if not found
        """
        try:
            # Use GitHub API to get tag info
            result = self.run_gh_command([
                "api",
                f"/repos/{self.SDK_REPO}/git/ref/tags/{version}",
                "--jq", ".object.sha"
            ])
            return result.strip()
        except Exception as e:
            print(f"Error getting commit for version {version}: {e}", file=sys.stderr)
            return None

    def fetch_version_diff(self, old_version: str, new_version: str) -> str:
        """Fetch diff between two versions from SDK repository.

        Args:
            old_version: Old version tag
            new_version: New version tag

        Returns:
            Diff text
        """
        try:
            # Use GitHub API to get compare diff
            result = self.run_gh_command([
                "api",
                f"/repos/{self.SDK_REPO}/compare/{old_version}...{new_version}",
                "--jq", ".files[] | select(.filename | startswith(\"models/\")) | .filename + \"\\n\" + .patch"
            ])
            return result
        except Exception as e:
            print(f"Error fetching diff: {e}", file=sys.stderr)
            return ""

    def parse_go_struct_changes(self, diff_text: str) -> List[ModelChange]:
        """Parse Go struct changes from diff text.

        Args:
            diff_text: Diff text containing model changes

        Returns:
            List of ModelChange objects
        """
        model_changes: Dict[str, ModelChange] = {}
        current_file = None
        current_model = None

        lines = diff_text.split('\n')

        for i, line in enumerate(lines):
            # Detect file path
            if line.startswith('models/') and not line.startswith('+++') and not line.startswith('---'):
                # This is a filename line
                current_file = line.strip()
                # Extract model name from filename
                # e.g., models/agent_identity_blueprint.go -> AgentIdentityBlueprint
                filename = Path(current_file).stem
                # Convert snake_case to PascalCase
                current_model = ''.join(word.capitalize() for word in filename.split('_'))

                if current_file not in model_changes:
                    model_changes[current_file] = ModelChange(
                        file_path=current_file,
                        model_name=current_model
                    )
                continue

            if not current_file or not line.strip():
                continue

            # Look for struct field changes
            # Pattern: added/removed field in Go struct
            if line.startswith('+') and not line.startswith('+++'):
                # Added line
                field_info = self._parse_go_field(line[1:].strip())
                if field_info:
                    model_changes[current_file].added_fields.append(
                        FieldChange(
                            field_name=field_info[0],
                            field_type=field_info[1],
                            change_type='added',
                            line_number=i
                        )
                    )
            elif line.startswith('-') and not line.startswith('---'):
                # Removed line
                field_info = self._parse_go_field(line[1:].strip())
                if field_info:
                    model_changes[current_file].removed_fields.append(
                        FieldChange(
                            field_name=field_info[0],
                            field_type=field_info[1],
                            change_type='removed',
                            line_number=i
                        )
                    )

        # Filter out models with no changes
        return [change for change in model_changes.values() if change.has_changes]

    def _parse_go_field(self, line: str) -> Optional[Tuple[str, str]]:
        """Parse a Go struct field line.

        Args:
            line: Line of Go code

        Returns:
            Tuple of (field_name, field_type) or None if not a field
        """
        # Skip non-field lines
        if not line or line.startswith('//') or line.startswith('type ') or \
           line.startswith('package ') or line.startswith('import ') or \
           line.startswith('func ') or line.startswith('}') or line.startswith('{'):
            return None

        # Match Go struct field pattern: FieldName *type.Type `json:"fieldName"`
        # Or: FieldName string `json:"fieldName"`
        match = re.match(r'(\w+)\s+([\*\[\]]?[\w\.]+(?:\[[\w\.]+\])?)\s*(?:`.*`)?', line)
        if match:
            field_name = match.group(1)
            field_type = match.group(2)

            # Skip if it looks like a method or other non-field
            if field_name[0].isupper():  # Go exported field
                return (field_name, field_type)

        return None

    def create_schema_update_issue(
        self,
        old_version: str,
        new_version: str,
        model_changes: List[ModelChange],
        changelog_section: str
    ) -> Optional[str]:
        """Create a GitHub issue for schema updates.

        Args:
            old_version: Old SDK version
            new_version: New SDK version
            model_changes: List of detected model changes
            changelog_section: Relevant changelog section

        Returns:
            Issue number if created, None otherwise
        """
        if self.dry_run:
            print("\n🔍 DRY RUN: Would create issue with following content:")
            print("=" * 80)

        # Build issue title
        title = f"Schema Update Required: Microsoft Graph SDK {old_version} → {new_version}"

        # Build issue body
        body_parts = []
        body_parts.append("## Summary")
        body_parts.append(f"The Microsoft Graph Beta SDK has been updated from `{old_version}` to `{new_version}`.")
        body_parts.append(f"This update includes {len(model_changes)} model(s) with schema changes that require review and potential Terraform schema updates.")
        body_parts.append("")

        body_parts.append("## Changed Models")
        body_parts.append("")

        for change in model_changes:
            body_parts.append(f"### `{change.model_name}` ({change.file_path})")
            body_parts.append(f"**Changes:** {change.change_summary}")
            body_parts.append("")

            if change.added_fields:
                body_parts.append("**Added Fields:**")
                for field in change.added_fields:
                    body_parts.append(f"- `{field.field_name}` ({field.field_type})")
                body_parts.append("")

            if change.removed_fields:
                body_parts.append("**Removed Fields:**")
                for field in change.removed_fields:
                    body_parts.append(f"- `{field.field_name}` ({field.field_type})")
                body_parts.append("")

        body_parts.append("## Action Required")
        body_parts.append("")
        body_parts.append("1. Review each changed model listed above")
        body_parts.append("2. Update corresponding Terraform resource schemas")
        body_parts.append("3. Update resource CRUD operations if needed")
        body_parts.append("4. Add/update tests for new fields")
        body_parts.append("5. Update documentation")
        body_parts.append("")

        body_parts.append("## References")
        body_parts.append("")
        body_parts.append(f"- [SDK Changelog](https://github.com/{self.SDK_REPO}/blob/main/CHANGELOG.md)")
        body_parts.append(f"- [Version Diff](https://github.com/{self.SDK_REPO}/compare/{old_version}...{new_version})")
        body_parts.append(f"- [Models Diff](https://github.com/{self.SDK_REPO}/compare/{old_version}...{new_version}#files_bucket)")
        body_parts.append("")

        if changelog_section and "not found" not in changelog_section.lower():
            body_parts.append("## Changelog Excerpt")
            body_parts.append("")
            body_parts.append("```")
            body_parts.append(changelog_section[:1000])  # Limit length
            if len(changelog_section) > 1000:
                body_parts.append("... (truncated)")
            body_parts.append("```")
            body_parts.append("")

        body_parts.append("---")
        body_parts.append(f"🤖 Auto-generated by sdkSchemaChangeDetector.py on {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

        body = '\n'.join(body_parts)

        if self.dry_run:
            print(f"Title: {title}\n")
            print(body)
            print("=" * 80)
            return None

        # Create issue using gh CLI
        try:
            # Ensure labels exist
            self._ensure_labels()

            # Create issue
            issue_url = self.run_gh_command([
                "issue", "create",
                "--repo", self.repo,
                "--title", title,
                "--body", body,
                "--label", "sdk-update,schema-change,needs-review"
            ])

            # Extract issue number from URL
            issue_number = issue_url.split('/')[-1]
            print(f"✅ Created issue #{issue_number}: {issue_url}")
            return issue_number

        except Exception as e:
            print(f"Error creating issue: {e}", file=sys.stderr)
            return None

    def _ensure_labels(self):
        """Ensure required labels exist in the repository."""
        labels = [
            ("sdk-update", "0E8A16", "Microsoft Graph SDK update"),
            ("schema-change", "D93F0B", "Schema changes detected"),
            ("needs-review", "FBCA04", "Requires engineer review"),
        ]

        for label_name, color, description in labels:
            try:
                self.run_gh_command([
                    "label", "create", label_name,
                    "--repo", self.repo,
                    "--color", color,
                    "--description", description,
                    "--force"
                ])
            except Exception:
                # Label might already exist, that's fine
                pass

    def analyze_pr(self, pr_number: int) -> bool:
        """Analyze a PR for SDK schema changes.

        Args:
            pr_number: PR number to analyze

        Returns:
            True if analysis completed successfully
        """
        print(f"🔍 Analyzing PR #{pr_number}...")

        # Get version change from PR
        version_change = self.get_pr_version_change(pr_number)
        if not version_change:
            print("❌ Could not detect version change in PR", file=sys.stderr)
            return False

        old_version, new_version = version_change
        return self.analyze_version_change(old_version, new_version)

    def analyze_version_change(self, old_version: str, new_version: str) -> bool:
        """Analyze a version change for schema updates.

        Args:
            old_version: Old SDK version
            new_version: New SDK version

        Returns:
            True if analysis completed successfully
        """
        print(f"📊 Analyzing version change: {old_version} → {new_version}")

        # Validate version increment
        if not self.validate_version_increment(old_version, new_version):
            print("⚠️  Version increment validation failed, but continuing analysis...")

        # Fetch changelog
        print("📖 Fetching changelog...")
        changelog = self.fetch_changelog_section(new_version)

        # Fetch diff
        print("🔄 Fetching version diff...")
        diff_text = self.fetch_version_diff(old_version, new_version)

        if not diff_text:
            print("⚠️  No model changes detected in diff")
            return True

        # Parse model changes
        print("🔬 Parsing model changes...")
        model_changes = self.parse_go_struct_changes(diff_text)

        if not model_changes:
            print("✅ No struct field changes detected in models")
            return True

        print(f"📝 Detected changes in {len(model_changes)} model(s)")

        # Create issue
        print("🎫 Creating GitHub issue...")
        issue_number = self.create_schema_update_issue(
            old_version, new_version, model_changes, changelog
        )

        if issue_number:
            print(f"✅ Analysis complete! Issue #{issue_number} created")
            return True
        elif self.dry_run:
            print("✅ Dry run complete!")
            return True
        else:
            print("❌ Failed to create issue", file=sys.stderr)
            return False


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Detect schema changes in Microsoft Graph SDK updates",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )

    parser.add_argument(
        "--pr-number",
        type=int,
        help="Dependabot PR number to analyze"
    )
    parser.add_argument(
        "--current",
        type=str,
        help="Current SDK version (auto-detected if not provided)"
    )
    parser.add_argument(
        "--new",
        type=str,
        help="New SDK version (required if --pr-number not provided)"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Analyze without creating GitHub issues"
    )
    parser.add_argument(
        "--repo",
        type=str,
        help="Repository in owner/repo format (auto-detected if not provided)"
    )

    args = parser.parse_args()

    # Validate arguments
    if not args.pr_number and not (args.current and args.new):
        parser.error("Either --pr-number or both --current and --new must be provided")

    # Create detector
    detector = SDKSchemaChangeDetector(repo=args.repo, dry_run=args.dry_run)

    # Run analysis
    try:
        if args.pr_number:
            success = detector.analyze_pr(args.pr_number)
        else:
            current = args.current or detector.parse_go_mod_version()
            if not current:
                print("❌ Could not determine current version", file=sys.stderr)
                sys.exit(1)
            success = detector.analyze_version_change(current, args.new)

        sys.exit(0 if success else 1)

    except KeyboardInterrupt:
        print("\n❌ Interrupted by user", file=sys.stderr)
        sys.exit(130)
    except Exception as e:
        print(f"❌ Unexpected error: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
