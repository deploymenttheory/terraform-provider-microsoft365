#!/usr/bin/env python3
"""
Detect SDK schema changes and create GitHub issues.

This script analyzes Microsoft Graph SDK version changes and creates
GitHub issues for schema updates that affect the Terraform provider.

Usage:
    # Analyze a specific PR
    ./detect_schema_changes.py --pr-number 1686

    # Analyze version change directly
    ./detect_schema_changes.py --current v0.156.0 --new v0.157.0

    # Dry run (don't create issues)
    ./detect_schema_changes.py --pr-number 1686 --dry-run

    # Filter by provider usage
    ./detect_schema_changes.py --pr-number 1686 --filter-by-usage provider_model_usage.json
"""

import sys
import json
import argparse
import subprocess
import traceback
from pathlib import Path
from datetime import datetime
from typing import Optional, Tuple, List
from urllib.parse import urlparse

# Add current directory to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from core import (  # pylint: disable=import-error
    ProgressReporter,
    VersionParser,
    GitHubClient,
    DiffFetcher,
    StructParser,
    IssueBuilder,
)
from models import (  # pylint: disable=import-error
    VersionChangeResult,
    IssueCreationResult,
    DetectionResult,
    ModelChange,
)


class SchemaChangeDetector:
    """Main detector for SDK schema changes."""

    SDK_MODULE = "github.com/microsoftgraph/msgraph-beta-sdk-go"
    SDK_REPO = "microsoftgraph/msgraph-beta-sdk-go"
    REQUIRED_LABELS = [
        ("sdk-update", "0E8A16", "Microsoft Graph SDK update"),
        ("schema-change", "D93F0B", "Schema changes detected"),
        ("needs-review", "FBCA04", "Requires engineer review"),
    ]

    def __init__(self, repo: Optional[str] = None, dry_run: bool = False,
                 save_diff: bool = False, verbose: bool = True):
        """Initialize detector.
        
        Args:
            repo: Target repository (owner/repo)
            dry_run: If True, don't create actual GitHub issues
            save_diff: If True, save fetched diff to file
            verbose: If True, show detailed progress
        """
        self.dry_run = dry_run
        self.go_mod_path = Path.cwd() / "go.mod"
        
        self.reporter = ProgressReporter(verbose)
        self.version_parser = VersionParser()
        
        target_repo = repo or self._get_current_repo_from_git()
        
        self.github_client = GitHubClient(target_repo, self.reporter)
        self.diff_fetcher = DiffFetcher(
            self.SDK_REPO,
            self.github_client,
            self.reporter,
            save_diff
        )
        self.struct_parser = StructParser(self.reporter)
        self.issue_builder = IssueBuilder(self.SDK_REPO)

    def _get_current_repo_from_git(self) -> str:
        """Get current repository from git remote."""
        try:
            result = subprocess.run(
                ["git", "remote", "get-url", "origin"],
                capture_output=True,
                text=True,
                check=True
            )
            remote_url = result.stdout.strip()
            
            # Parse URL properly to avoid security issues
            # Handle both HTTPS and SSH formats
            if remote_url.startswith("git@"):
                # SSH format: git@github.com:owner/repo.git
                if remote_url.startswith("git@github.com:"):
                    path = remote_url.replace("git@github.com:", "").replace(".git", "")
                    return path
            else:
                # HTTPS format: https://github.com/owner/repo.git
                parsed = urlparse(remote_url)
                # Verify hostname is exactly github.com
                if parsed.hostname == "github.com":
                    # Remove leading slash and .git suffix
                    path = parsed.path.lstrip("/").replace(".git", "")
                    return path
        except (subprocess.CalledProcessError, ValueError):
            pass
        
        return "deploymenttheory/terraform-provider-microsoft365"

    def parse_go_mod_version(self) -> Optional[str]:
        """Parse SDK version from go.mod file."""
        try:
            with open(self.go_mod_path, 'r', encoding='utf-8') as f:
                for line in f:
                    if self.SDK_MODULE in line:
                        version = self.version_parser.extract_version_from_line(line)
                        if version:
                            return version
        except OSError as e:
            self.reporter.error(f"Error reading go.mod: {e}")
        
        return None

    def get_pr_version_change(self, pr_number: int) -> VersionChangeResult:
        """Get version change from a PR."""
        try:
            diff = self.github_client.get_pr_diff(pr_number)
            
            old_version = None
            new_version = None
            
            for line in diff.split('\n'):
                if self.SDK_MODULE in line:
                    if line.startswith('-') and not line.startswith('---'):
                        old_version = self.version_parser.extract_version_from_line(line)
                    elif line.startswith('+') and not line.startswith('+++'):
                        new_version = self.version_parser.extract_version_from_line(line)
            
            if old_version and new_version:
                return VersionChangeResult.success(old_version, new_version)
            
            return VersionChangeResult.not_found()
            
        except (subprocess.CalledProcessError, OSError, ValueError) as e:
            return VersionChangeResult.error(str(e))

    def analyze_pr(self, pr_number: int, save_results: Optional[Path] = None,
                   filter_by_usage: Optional[Path] = None) -> bool:
        """Analyze a PR for schema changes."""
        self.reporter.section(f"🔍 Analyzing PR #{pr_number}...")

        version_change = self.get_pr_version_change(pr_number)
        if not version_change.is_success:
            self.reporter.error(version_change.error_message or "Could not detect version change")
            return False

        return self.analyze_version_change(
            version_change.old_version,
            version_change.new_version,
            save_results,
            filter_by_usage
        )

    def analyze_version_change(self, old_version: str, new_version: str,
                              save_results: Optional[Path] = None,
                              filter_by_usage: Optional[Path] = None) -> bool:
        """Analyze version change for schema updates."""
        self.reporter.section(f"📊 Analyzing version change: {old_version} → {new_version}")

        # Validate version increment
        validation = self.version_parser.validate_increment(old_version, new_version, self.reporter)
        if not validation.is_valid:
            self.reporter.warning("Version increment validation failed, but continuing...")

        # Fetch changelog
        self.reporter.section("📖 Fetching changelog...")
        changelog = self.diff_fetcher.fetch_changelog_section(new_version)

        # Fetch diff
        self.reporter.section("🔄 Fetching version diff...")
        diff_text = self.diff_fetcher.fetch_version_diff(old_version, new_version)

        if not diff_text:
            self.reporter.warning("No model changes detected in diff")
            return True

        # Parse model changes
        self.reporter.section("🔬 Parsing model changes...")
        model_changes = self.struct_parser.parse_diff(diff_text)

        if not model_changes:
            self.reporter.success("No struct field changes detected in models")
            return True

        self.reporter.info(f"📝 Detected changes in {len(model_changes)} model(s)")

        # Filter by usage if requested
        original_count = len(model_changes)
        filtered_count = 0
        if filter_by_usage:
            model_changes, filtered_count = self._filter_by_usage(model_changes, filter_by_usage)
            self.reporter.info(
                f"🔍 Filtered to {len(model_changes)} relevant model(s) "
                f"({filtered_count} filtered out)"
            )

        # Save results if requested
        if save_results:
            self._save_results(old_version, new_version, model_changes, changelog,
                             original_count, filtered_count, save_results)

        if not model_changes:
            self.reporter.success("No relevant schema changes detected after filtering")
            return True

        # Create issue
        self.reporter.section("🎫 Creating GitHub issue...")
        result = self._create_issue(old_version, new_version, model_changes, changelog)

        if result.is_success and result.issue_number != "DRY_RUN":
            self.reporter.success(f"Analysis complete! Issue #{result.issue_number} created")
            return True
        elif result.is_success and result.issue_number == "DRY_RUN":
            self.reporter.success("Dry run complete!")
            return True
        else:
            self.reporter.error(result.error_message or "Failed to create issue")
            return False

    def _filter_by_usage(self, model_changes: List[ModelChange],
                        usage_file: Path) -> Tuple[List[ModelChange], int]:
        """Filter model changes to only those used in the provider."""
        self.reporter.section("🔍 Filtering changes by provider usage...")
        
        try:
            with open(usage_file, 'r', encoding='utf-8') as f:
                usage_data = json.load(f)
        except (OSError, json.JSONDecodeError) as e:
            self.reporter.error(f"Failed to load usage data: {e}")
            return model_changes, 0
        
        used_models = {model['model_file'] for model in usage_data.get('models', [])}
        self.reporter.info(f"  Loaded {len(used_models)} model(s) from provider usage data")
        
        filtered = [mc for mc in model_changes if mc.file_path in used_models]
        filtered_out = len(model_changes) - len(filtered)
        
        if filtered_out > 0:
            self.reporter.info(f"  ✓ Kept {len(filtered)} relevant model(s)")
            self.reporter.info(f"  ✗ Filtered {filtered_out} unused model(s)")
        
        return filtered, filtered_out

    def _save_results(self, old_version: str, new_version: str,
                     model_changes: List[ModelChange], changelog: str,
                     total_models: int, filtered_count: int,
                     output_file: Path):
        """Save detection results to JSON."""
        result = DetectionResult(
            current_version=old_version,
            new_version=new_version,
            timestamp=datetime.now().isoformat(),
            total_models_changed=total_models,
            models_with_changes=len(model_changes),
            filtered_models=filtered_count,
            model_changes=model_changes,
            changelog_section=changelog,
            statistics=self.struct_parser.statistics
        )
        
        try:
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(result.to_dict(), f, indent=2)
            self.reporter.success(f"💾 Results saved to: {output_file}")
        except OSError as e:
            self.reporter.error(f"Failed to save results: {e}")

    def _create_issue(self, old_version: str, new_version: str,
                     model_changes: List[ModelChange], changelog: str) -> IssueCreationResult:
        """Create GitHub issue for schema updates."""
        title = self.issue_builder.build_title(old_version, new_version)
        body = self.issue_builder.build_body(old_version, new_version, model_changes, changelog)

        if self.dry_run:
            self.reporter.print_dry_run_issue(title, body)
            return IssueCreationResult.dry_run()

        try:
            self.github_client.ensure_labels(self.REQUIRED_LABELS)
            issue_url = self.github_client.create_issue(
                title, body,
                [label[0] for label in self.REQUIRED_LABELS]
            )
            issue_number = issue_url.split('/')[-1]
            return IssueCreationResult.success(issue_number, issue_url)
        except (subprocess.CalledProcessError, OSError, ValueError) as e:
            return IssueCreationResult.error(str(e))


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Detect schema changes in Microsoft Graph SDK updates",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )

    parser.add_argument("--pr-number", type=int, help="PR number to analyze")
    parser.add_argument("--current", type=str, help="Current SDK version")
    parser.add_argument("--new", type=str, help="New SDK version")
    parser.add_argument("--dry-run", action="store_true", help="Analyze without creating issues")
    parser.add_argument("--save-diff", action="store_true", help="Save fetched diff to file")
    parser.add_argument("--repo", type=str, help="Repository (owner/repo)")
    parser.add_argument("--save-results", type=Path, help="Save results to JSON file")
    parser.add_argument("--filter-by-usage", type=Path, help="Filter by provider usage (path to JSON)")

    args = parser.parse_args()

    # Validate arguments
    if not args.pr_number and not (args.current and args.new):
        parser.error("Either --pr-number or both --current and --new must be provided")

    # Create detector
    detector = SchemaChangeDetector(
        repo=args.repo,
        dry_run=args.dry_run,
        save_diff=args.save_diff
    )

    # Run analysis
    try:
        if args.pr_number:
            success = detector.analyze_pr(
                args.pr_number,
                save_results=args.save_results,
                filter_by_usage=args.filter_by_usage
            )
        else:
            current = args.current or detector.parse_go_mod_version()
            if not current:
                print("❌ Could not determine current version", file=sys.stderr)
                sys.exit(1)
            success = detector.analyze_version_change(
                current, args.new,
                save_results=args.save_results,
                filter_by_usage=args.filter_by_usage
            )

        sys.exit(0 if success else 1)

    except KeyboardInterrupt:
        print("\n❌ Interrupted by user", file=sys.stderr)
        sys.exit(130)
    except Exception as e:  # pylint: disable=broad-except
        # Catch-all for unexpected errors at top level
        print(f"❌ Unexpected error: {e}", file=sys.stderr)
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
