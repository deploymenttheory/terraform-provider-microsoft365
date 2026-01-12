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
    ./kiota_graph_sdk_schema_change_detector.py --pr-number 1686

    # Analyze version change directly
    ./kiota_graph_sdk_schema_change_detector.py --current v0.156.0 --new v0.157.0

    # Dry run (don't create issues)
    ./kiota_graph_sdk_schema_change_detector.py --pr-number 1686 --dry-run

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
import traceback
from pathlib import Path
from typing import Optional, Dict, List, Tuple
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum


# ============================================================================
# REGEX PATTERNS - Named constants for all regular expressions
# ============================================================================

class RegexPatterns:
    """Collection of compiled regex patterns used throughout the script."""
    
    # Version patterns
    VERSION_FULL = re.compile(r'v(\d+)\.(\d+)\.(\d+)')
    VERSION_IN_TEXT = re.compile(r'v\d+\.\d+\.\d+')
    
    # GitHub URL patterns
    GITHUB_REPO_URL = re.compile(r'github\.com[:/](.+/.+?)(\.git)?$')
    
    # File patterns
    MODEL_FILE_PATH = re.compile(r'models/[\w_]+\.go')
    
    # Go code patterns - Structs
    GO_STRUCT_FIELD = re.compile(r'(\w+)\s+([\*\[\]]?[\w\.]+(?:\[[\w\.]+\])?)\s*(?:`.*`)?')
    
    # Go code patterns - Interfaces
    GO_INTERFACE_METHOD = re.compile(r'(\w+)\s*\(([^)]*)\)\s*(\([^)]*\))?')
    GO_EMBEDDED_TYPE = re.compile(r'^\s*(\w+[\w\.]*)\s*$')
    
    # Go declarations
    GO_TYPE_STRUCT = re.compile(r'type\s+(\w+)\s+struct')
    GO_TYPE_INTERFACE = re.compile(r'type\s+(\w+)\s+interface')
    
    # Changelog patterns
    CHANGELOG_VERSION_HEADER = r'##'


# ============================================================================
# RESULT TYPES - Structured return values instead of Optional
# ============================================================================

class ResultStatus(Enum):
    """Status of an operation result."""
    SUCCESS = "success"
    ERROR = "error"
    NOT_FOUND = "not_found"
    INVALID = "invalid"


@dataclass
class VersionResult:
    """Result of version parsing or extraction."""
    status: ResultStatus
    version: Optional[str] = None
    version_tuple: Optional[Tuple[int, int, int]] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, version: str, version_tuple: Tuple[int, int, int]) -> 'VersionResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, version, version_tuple)
    
    @classmethod
    def error(cls, message: str) -> 'VersionResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def not_found(cls, message: str = "Version not found") -> 'VersionResult':
        """Create a not found result."""
        return cls(ResultStatus.NOT_FOUND, error_message=message)


@dataclass
class VersionChangeResult:
    """Result of PR version change detection."""
    status: ResultStatus
    old_version: Optional[str] = None
    new_version: Optional[str] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, old_version: str, new_version: str) -> 'VersionChangeResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, old_version, new_version)
    
    @classmethod
    def error(cls, message: str) -> 'VersionChangeResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def not_found(cls) -> 'VersionChangeResult':
        """Create a not found result."""
        return cls(ResultStatus.NOT_FOUND, error_message="Version change not detected in PR")


@dataclass
class DetectionResult:
    """Complete schema change detection result."""
    current_version: str
    new_version: str
    timestamp: str
    total_models_changed: int
    models_with_changes: int
    filtered_models: int
    model_changes: List['ModelChange']
    changelog_section: str
    statistics: 'ParseStatistics'
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            "current_version": self.current_version,
            "new_version": self.new_version,
            "timestamp": self.timestamp,
            "total_models_changed": self.total_models_changed,
            "models_with_changes": self.models_with_changes,
            "filtered_models": self.filtered_models,
            "model_changes": [
                {
                    "file_path": mc.file_path,
                    "model_name": mc.model_name,
                    "added_fields": [{"name": f.name, "field_type": f.field_type} for f in mc.added_fields],
                    "removed_fields": [{"name": f.name, "field_type": f.field_type} for f in mc.removed_fields],
                    "added_methods": [{"signature": m.signature} for m in mc.added_methods],
                    "removed_methods": [{"signature": m.signature} for m in mc.removed_methods],
                    "added_embedded_types": [{"type_name": e.type_name} for e in mc.added_embedded_types],
                    "removed_embedded_types": [{"type_name": e.type_name} for e in mc.removed_embedded_types],
                }
                for mc in self.model_changes
            ],
            "changelog_section": self.changelog_section,
            "statistics": {
                "total_files_in_diff": self.statistics.total_files_in_diff,
                "files_with_changes": self.statistics.files_with_changes,
                "files_without_changes": self.statistics.files_without_changes,
                "total_lines_processed": self.statistics.total_lines_processed,
                "added_lines_processed": self.statistics.added_lines_processed,
                "removed_lines_processed": self.statistics.removed_lines_processed,
                "struct_fields_added": self.statistics.struct_fields_added,
                "struct_fields_removed": self.statistics.struct_fields_removed,
                "interface_methods_added": self.statistics.interface_methods_added,
                "interface_methods_removed": self.statistics.interface_methods_removed,
                "embedded_types_added": self.statistics.embedded_types_added,
                "embedded_types_removed": self.statistics.embedded_types_removed,
            }
        }


@dataclass
class IssueCreationResult:
    """Result of GitHub issue creation."""
    status: ResultStatus
    issue_number: Optional[str] = None
    issue_url: Optional[str] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, issue_number: str, issue_url: str) -> 'IssueCreationResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, issue_number, issue_url)
    
    @classmethod
    def error(cls, message: str) -> 'IssueCreationResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def dry_run(cls) -> 'IssueCreationResult':
        """Create a dry run result."""
        return cls(ResultStatus.SUCCESS, issue_number="DRY_RUN")


@dataclass
class ValidationResult:
    """Result of version validation."""
    is_valid: bool
    message: Optional[str] = None
    
    @classmethod
    def valid(cls) -> 'ValidationResult':
        """Create a valid result."""
        return cls(True)
    
    @classmethod
    def invalid(cls, message: str) -> 'ValidationResult':
        """Create an invalid result."""
        return cls(False, message)


# ============================================================================
# DATA MODELS
# ============================================================================

@dataclass
class FieldChange:
    """Represents a field change in a Go struct."""
    field_name: str
    field_type: str
    change_type: str  # 'added' or 'removed'
    line_number: Optional[int] = None


@dataclass
class MethodChange:
    """Represents an interface method change."""
    method_name: str
    parameters: str
    return_type: str
    change_type: str  # 'added' or 'removed'
    line_number: Optional[int] = None
    
    @property
    def signature(self) -> str:
        """Get the full method signature."""
        params = f"({self.parameters})" if self.parameters else "()"
        returns = f" {self.return_type}" if self.return_type else ""
        return f"{self.method_name}{params}{returns}"


@dataclass
class EmbeddedTypeChange:
    """Represents a change in embedded types (interfaces or structs)."""
    type_name: str
    change_type: str  # 'added' or 'removed'
    context: str  # 'interface' or 'struct'
    line_number: Optional[int] = None


@dataclass
class ModelChange:
    """Represents changes to a Go model file."""
    file_path: str
    model_name: str
    added_fields: List[FieldChange] = field(default_factory=list)
    removed_fields: List[FieldChange] = field(default_factory=list)
    added_methods: List[MethodChange] = field(default_factory=list)
    removed_methods: List[MethodChange] = field(default_factory=list)
    added_embedded_types: List[EmbeddedTypeChange] = field(default_factory=list)
    removed_embedded_types: List[EmbeddedTypeChange] = field(default_factory=list)

    @property
    def has_changes(self) -> bool:
        """Check if this model has any changes."""
        return bool(
            self.added_fields or self.removed_fields or
            self.added_methods or self.removed_methods or
            self.added_embedded_types or self.removed_embedded_types
        )

    @property
    def change_summary(self) -> str:
        """Get a summary of changes."""
        parts = []
        if self.added_fields:
            parts.append(f"+{len(self.added_fields)} fields")
        if self.removed_fields:
            parts.append(f"-{len(self.removed_fields)} fields")
        if self.added_methods:
            parts.append(f"+{len(self.added_methods)} methods")
        if self.removed_methods:
            parts.append(f"-{len(self.removed_methods)} methods")
        if self.added_embedded_types:
            parts.append(f"+{len(self.added_embedded_types)} embedded")
        if self.removed_embedded_types:
            parts.append(f"-{len(self.removed_embedded_types)} embedded")
        return ", ".join(parts)


@dataclass
class ParseStatistics:
    """Statistics from parsing diff for diagnostic purposes."""
    total_files_in_diff: int = 0
    files_with_changes: int = 0
    files_without_changes: int = 0
    total_lines_processed: int = 0
    added_lines_processed: int = 0
    removed_lines_processed: int = 0
    
    # Struct tracking
    struct_fields_added: int = 0
    struct_fields_removed: int = 0
    
    # Interface tracking
    interface_methods_added: int = 0
    interface_methods_removed: int = 0
    
    # Embedded types tracking
    embedded_types_added: int = 0
    embedded_types_removed: int = 0
    
    # Filtering reasons
    lines_filtered_comments: int = 0
    lines_filtered_declarations: int = 0
    lines_filtered_func_impl: int = 0
    lines_filtered_no_match: int = 0
    lines_filtered_unexported: int = 0
    
    def get_summary(self) -> str:
        """Get a human-readable summary."""
        parts = [
            "  📊 Parsing Statistics:",
            f"     Files in diff: {self.total_files_in_diff}",
            f"     Files with changes: {self.files_with_changes}",
            f"     Files without changes: {self.files_without_changes}",
            "",
            "  📝 Lines Processed:",
            f"     Total change lines: {self.added_lines_processed + self.removed_lines_processed}",
            f"     Added lines (+): {self.added_lines_processed}",
            f"     Removed lines (-): {self.removed_lines_processed}",
            "",
            "  🔧 Changes Detected:",
            f"     Struct fields added: {self.struct_fields_added}",
            f"     Struct fields removed: {self.struct_fields_removed}",
            f"     Interface methods added: {self.interface_methods_added}",
            f"     Interface methods removed: {self.interface_methods_removed}",
            f"     Embedded types added: {self.embedded_types_added}",
            f"     Embedded types removed: {self.embedded_types_removed}",
            "",
            "  🔍 Filtering Breakdown:",
            f"     Comments (//): {self.lines_filtered_comments}",
            f"     Type/package/import declarations: {self.lines_filtered_declarations}",
            f"     Function implementations (func): {self.lines_filtered_func_impl}",
            f"     Unexported fields (lowercase): {self.lines_filtered_unexported}",
            f"     No regex match: {self.lines_filtered_no_match}",
        ]
        return "\n".join(parts)


class ProgressReporter:
    """Handles all user-facing output and progress reporting."""

    def __init__(self, verbose: bool = True):
        """Initialize the reporter.
        
        Args:
            verbose: If True, show detailed progress messages
        """
        self.verbose = verbose

    def section(self, message: str):
        """Print a section header."""
        if self.verbose:
            print(f"\n{message}")

    def info(self, message: str, indent: int = 0):
        """Print an info message."""
        if self.verbose:
            prefix = "  " * indent
            print(f"{prefix}{message}")

    def success(self, message: str):
        """Print a success message."""
        print(f"✅ {message}")

    def warning(self, message: str):
        """Print a warning message."""
        print(f"⚠️  {message}", file=sys.stderr)

    def error(self, message: str):
        """Print an error message."""
        print(f"❌ {message}", file=sys.stderr)

    def print_parse_summary(self, model_changes: List[ModelChange], stats: 'ParseStatistics' = None,
                           files_without_changes: List['ModelChange'] = None):
        """Print summary of parsed model changes with statistics.
        
        Args:
            model_changes: List of model changes detected
            stats: Optional parsing statistics for diagnostics
            files_without_changes: Optional list of files that had no field changes
        """
        self.info("\n📋 Parse Summary:", indent=1)
        
        if model_changes:
            self.info(f"✓ Found {len(model_changes)} model(s) with changes:\n", indent=1)
            
            for change in model_changes:
                self.info(f"📄 {change.model_name} ({change.file_path})", indent=2)
                self.info(change.change_summary, indent=3)

                # Show struct fields
                if change.added_fields:
                    self.info("Added struct fields:", indent=3)
                    for fld in change.added_fields[:5]:
                        self.info(f"+ {fld.field_name}: {fld.field_type}", indent=4)
                    if len(change.added_fields) > 5:
                        self.info(f"... and {len(change.added_fields) - 5} more", indent=4)

                if change.removed_fields:
                    self.info("Removed struct fields:", indent=3)
                    for fld in change.removed_fields[:5]:
                        self.info(f"- {fld.field_name}: {fld.field_type}", indent=4)
                    if len(change.removed_fields) > 5:
                        self.info(f"... and {len(change.removed_fields) - 5} more", indent=4)

                # Show interface methods
                if change.added_methods:
                    self.info("Added interface methods:", indent=3)
                    for method in change.added_methods[:5]:
                        self.info(f"+ {method.signature}", indent=4)
                    if len(change.added_methods) > 5:
                        self.info(f"... and {len(change.added_methods) - 5} more", indent=4)

                if change.removed_methods:
                    self.info("Removed interface methods:", indent=3)
                    for method in change.removed_methods[:5]:
                        self.info(f"- {method.signature}", indent=4)
                    if len(change.removed_methods) > 5:
                        self.info(f"... and {len(change.removed_methods) - 5} more", indent=4)

                # Show embedded types
                if change.added_embedded_types:
                    self.info("Added embedded types:", indent=3)
                    for emb in change.added_embedded_types[:5]:
                        self.info(f"+ {emb.type_name} ({emb.context})", indent=4)
                    if len(change.added_embedded_types) > 5:
                        self.info(f"... and {len(change.added_embedded_types) - 5} more", indent=4)

                if change.removed_embedded_types:
                    self.info("Removed embedded types:", indent=3)
                    for emb in change.removed_embedded_types[:5]:
                        self.info(f"- {emb.type_name} ({emb.context})", indent=4)
                    if len(change.removed_embedded_types) > 5:
                        self.info(f"... and {len(change.removed_embedded_types) - 5} more", indent=4)

                print()
        else:
            self.info("ℹ️  No changes detected in diff", indent=1)
        
        # Print detailed statistics if provided
        if stats:
            print()
            print(stats.get_summary())
            
            # Explain why files were filtered
            if stats.files_without_changes > 0:
                print()
                self.info(f"ℹ️  {stats.files_without_changes} file(s) had changes but no detectable model modifications.", indent=1)
                self.info("   Possible reasons:", indent=1)
                self.info("   • Only comments, imports, or package declarations changed", indent=1)
                self.info("   • Method implementations (func body) changed", indent=1)
                self.info("   • Type aliases or constants changed", indent=1)
                self.info("   • Only unexported (lowercase) fields/methods changed", indent=1)
                self.info("   • Changes didn't match expected patterns", indent=1)
                
                # Show examples of files without field changes
                if files_without_changes:
                    print()
                    self.info("📁 Examples of files without field changes (showing up to 10):", indent=1)
                    for change in files_without_changes[:10]:
                        self.info(f"   • {change.model_name} ({change.file_path})", indent=1)
                    if len(files_without_changes) > 10:
                        self.info(f"   ... and {len(files_without_changes) - 10} more", indent=1)

    def print_dry_run_issue(self, title: str, body: str):
        """Print issue content for dry run."""
        print("\n🔍 DRY RUN: Would create issue with following content:")
        print("=" * 80)
        print(f"Title: {title}\n")
        print(body)
        print("=" * 80)


class VersionParser:
    """Handles version parsing and validation."""

    def parse_version(self, version_str: str) -> VersionResult:
        """Parse version string into structured result.
        
        Args:
            version_str: Version string (e.g., 'v0.156.0')
            
        Returns:
            VersionResult with parsed version data
        """
        match = RegexPatterns.VERSION_FULL.match(version_str)
        if match:
            version_tuple = tuple(map(int, match.groups()))
            return VersionResult.success(version_str, version_tuple)
        return VersionResult.error(f"Invalid version format: {version_str}")

    def extract_version_from_line(self, line: str) -> Optional[str]:
        """Extract version string from a line of text.
        
        Args:
            line: Line of text that may contain a version
            
        Returns:
            Version string or None if not found
        """
        match = RegexPatterns.VERSION_IN_TEXT.search(line)
        return match.group(0) if match else None

    def validate_increment(self, old_version: str, new_version: str, 
                          reporter: ProgressReporter) -> ValidationResult:
        """Validate that version increment is acceptable.
        
        Args:
            old_version: Old version string
            new_version: New version string
            reporter: Reporter for warnings
            
        Returns:
            ValidationResult with validation status
        """
        old_result = self.parse_version(old_version)
        new_result = self.parse_version(new_version)

        if not old_result.is_success or not new_result.is_success:
            reporter.error("Invalid version format")
            return ValidationResult.invalid("Invalid version format")

        old_major, old_minor, old_patch = old_result.version_tuple
        new_major, new_minor, new_patch = new_result.version_tuple

        if new_major == old_major:
            if new_minor == old_minor + 1 and new_patch == 0:
                return ValidationResult.valid()  # Valid minor version bump
            elif new_minor == old_minor and new_patch == old_patch + 1:
                return ValidationResult.valid()  # Valid patch version bump
            elif new_minor == old_minor and new_patch > old_patch:
                reporter.warning("Multiple patch version increment detected")
                return ValidationResult.valid()  # Multiple patch bump - acceptable
            elif new_minor > old_minor + 1:
                msg = f"Multiple minor version jump detected: {old_version} -> {new_version}"
                reporter.warning(msg)
                reporter.warning("This may indicate missing intermediate versions.")
                return ValidationResult.invalid(msg)

        msg = f"Unexpected version change: {old_version} -> {new_version}"
        reporter.warning(msg)
        return ValidationResult.invalid(msg)


class GitHubClient:
    """Handles all GitHub CLI operations."""

    def __init__(self, repo: str, reporter: ProgressReporter):
        """Initialize GitHub client.
        
        Args:
            repo: Repository in owner/repo format
            reporter: Progress reporter for output
        """
        self.repo = repo
        self.reporter = reporter

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
                self.reporter.error(f"Command failed: {' '.join(args)}")
                self.reporter.error(f"Error: {e.stderr}")
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

    def get_pr_diff(self, pr_number: int) -> str:
        """Get diff from a pull request.
        
        Args:
            pr_number: PR number
            
        Returns:
            Diff text
        """
        return self.run_gh_command([
            "pr", "diff", str(pr_number),
            "--repo", self.repo
        ])

    def get_compare_data(self, sdk_repo: str, old_version: str, new_version: str) -> dict:
        """Get comparison data between two versions from GitHub API.
        
        Args:
            sdk_repo: SDK repository (owner/repo)
            old_version: Old version tag
            new_version: New version tag
            
        Returns:
            Parsed JSON response from GitHub API
        """
        json_result = self.run_gh_command([
            "api",
            f"/repos/{sdk_repo}/compare/{old_version}...{new_version}"
        ])
        
        if not json_result:
            raise ValueError("No response from GitHub API")
        
        return json.loads(json_result)

    def create_issue(self, title: str, body: str, labels: List[str]) -> str:
        """Create a GitHub issue.
        
        Args:
            title: Issue title
            body: Issue body
            labels: List of label names
            
        Returns:
            Issue URL
        """
        return self.run_gh_command([
            "issue", "create",
            "--repo", self.repo,
            "--title", title,
            "--body", body,
            "--label", ",".join(labels)
        ])

    def ensure_labels(self, labels: List[Tuple[str, str, str]]):
        """Ensure labels exist in the repository.
        
        Args:
            labels: List of (name, color, description) tuples
        """
        for label_name, color, description in labels:
            try:
                self.run_gh_command([
                    "label", "create", label_name,
                    "--repo", self.repo,
                    "--color", color,
                    "--description", description,
                    "--force"
                ])
            except subprocess.CalledProcessError:
                pass  # Label might already exist

    def fetch_changelog(self, url: str) -> str:
        """Fetch changelog from URL.
        
        Args:
            url: Changelog URL
            
        Returns:
            Changelog content
        """
        stdout, _, _ = self.run_command(["curl", "-s", url])
        return stdout

    def get_current_repo(self) -> str:
        """Get current repository from git remote.
        
        Returns:
            Repository in owner/repo format
        """
        try:
            result = subprocess.run(
                ["git", "remote", "get-url", "origin"],
                capture_output=True,
                text=True,
                check=True
            )
            remote_url = result.stdout.strip()
            match = RegexPatterns.GITHUB_REPO_URL.search(remote_url)
            if match:
                return match.group(1).rstrip('.git')
        except subprocess.CalledProcessError:
            pass
        return "deploymenttheory/terraform-provider-microsoft365"


class DiffFetcher:
    """Handles fetching and filtering diffs from SDK repository."""

    def __init__(self, sdk_repo: str, github_client: GitHubClient, reporter: ProgressReporter, save_diff: bool = False):
        """Initialize diff fetcher.
        
        Args:
            sdk_repo: SDK repository (owner/repo)
            github_client: GitHub client for API calls
            reporter: Progress reporter
            save_diff: If True, save diffs to files
        """
        self.sdk_repo = sdk_repo
        self.github_client = github_client
        self.reporter = reporter
        self.save_diff = save_diff

    def fetch_version_diff(self, old_version: str, new_version: str) -> str:
        """Fetch diff between two versions.
        
        Args:
            old_version: Old version tag
            new_version: New version tag
            
        Returns:
            Unified diff text for model files
        """
        try:
            self.reporter.info("Fetching diff from GitHub API...", indent=1)
            
            compare_data = self.github_client.get_compare_data(
                self.sdk_repo, old_version, new_version
            )
            
            files = compare_data.get('files', [])
            if not files:
                self.reporter.info("No files changed in this version", indent=1)
                return ""

            model_files = self._filter_model_files(files)
            if not model_files:
                self.reporter.info("No model files with changes found", indent=1)
                return ""

            diff_text = self._build_unified_diff(model_files)
            self.reporter.info(f"📦 Collected diffs from {len(model_files)} model file(s)", indent=1)

            if self.save_diff:
                self._save_diff_to_file(diff_text, old_version, new_version)

            return diff_text

        except (subprocess.CalledProcessError, json.JSONDecodeError, OSError, KeyError) as e:
            self.reporter.error(f"Error fetching diff: {e}")
            traceback.print_exc()
            return ""

    def _filter_model_files(self, files: List[dict]) -> List[dict]:
        """Filter for model files with patches.
        
        Args:
            files: List of file info from GitHub API
            
        Returns:
            Filtered list of model files
        """
        model_files = []
        
        for file_info in files:
            filename = file_info.get('filename', '')
            
            if not filename.startswith('models/'):
                continue

            patch = file_info.get('patch')
            status = file_info.get('status', 'unknown')

            if not patch:
                self.reporter.info(
                    f"⚠️  Skipping {filename} (no patch content, status: {status})",
                    indent=2
                )
                continue

            model_files.append(file_info)
            self.reporter.info(f"✓ Found changes in {filename} ({status})", indent=2)

        return model_files

    def _build_unified_diff(self, files: List[dict]) -> str:
        """Build unified diff format from file patches.
        
        Args:
            files: List of file info with patches
            
        Returns:
            Unified diff text
        """
        diff_parts = []
        
        for file_info in files:
            filename = file_info['filename']
            patch = file_info['patch']
            
            diff_parts.append(f"diff --git a/{filename} b/{filename}")
            diff_parts.append(f"--- a/{filename}")
            diff_parts.append(f"+++ b/{filename}")
            diff_parts.append(patch)
            diff_parts.append("")

        return '\n'.join(diff_parts)

    def _save_diff_to_file(self, diff_text: str, old_version: str, new_version: str):
        """Save diff to file for debugging.
        
        Args:
            diff_text: Diff content
            old_version: Old version
            new_version: New version
        """
        diff_file = f"sdk_diff_{old_version}_to_{new_version}.patch"
        with open(diff_file, 'w', encoding='utf-8') as f:
            f.write(diff_text)
        self.reporter.info(f"💾 Saved diff to {diff_file}", indent=1)

    def fetch_changelog_section(self, version: str) -> str:
        """Fetch changelog section for a specific version.
        
        Args:
            version: Version to fetch changelog for
            
        Returns:
            Changelog section text
        """
        try:
            changelog_url = f"https://raw.githubusercontent.com/{self.sdk_repo}/main/CHANGELOG.md"
            changelog_content = self.github_client.fetch_changelog(changelog_url)
            
            return self._extract_version_section(changelog_content, version)

        except subprocess.CalledProcessError as e:
            self.reporter.error(f"Error fetching changelog: {e}")
            return "Error fetching changelog"

    def _extract_version_section(self, changelog: str, version: str) -> str:
        """Extract version section from changelog.
        
        Args:
            changelog: Full changelog content
            version: Version to extract
            
        Returns:
            Version section text
        """
        lines = changelog.split('\n')
        section_lines = []
        in_section = False

        for line in lines:
            if line.startswith('##') and version.lstrip('v') in line:
                in_section = True
                section_lines.append(line)
            elif in_section:
                if line.startswith('##'):
                    break
                section_lines.append(line)

        return '\n'.join(section_lines) if section_lines else "Changelog section not found"


class StructParser:
    """Parses Go model changes from diff text (structs, interfaces, embedded types)."""

    def __init__(self, reporter: ProgressReporter):
        """Initialize parser.
        
        Args:
            reporter: Progress reporter
        """
        self.reporter = reporter
        self.stats = ParseStatistics()
        self.in_interface_context = False  # Track if we're parsing inside an interface
    
    @property
    def statistics(self) -> ParseStatistics:
        """Get parsing statistics."""
        return self.stats

    def parse_diff(self, diff_text: str) -> List[ModelChange]:
        """Parse Go model changes from diff text.
        
        Args:
            diff_text: Unified diff text
            
        Returns:
            List of ModelChange objects
        """
        model_changes: Dict[str, ModelChange] = {}
        current_file = None
        current_model = None
        self.stats = ParseStatistics()  # Reset stats for new parse
        self.in_interface_context = False

        lines = diff_text.split('\n')
        self.stats.total_lines_processed = len(lines)

        for i, line in enumerate(lines):
            # Check for file header
            if self._is_file_header(line):
                filename = self._extract_filename(line)
                if filename and (line.startswith('diff --git') or 
                               (line.startswith('+++') and current_file != filename)):
                    current_file = filename
                    current_model = self._filename_to_model_name(filename)
                    self.in_interface_context = False  # Reset for new file
                    
                    if current_file not in model_changes:
                        model_changes[current_file] = ModelChange(
                            file_path=current_file,
                            model_name=current_model
                        )
                        self.stats.total_files_in_diff += 1
                continue

            if not current_file or not line.strip():
                continue

            # Track added/removed lines
            if line.startswith('+') and not line.startswith('+++'):
                self.stats.added_lines_processed += 1
            elif line.startswith('-') and not line.startswith('---'):
                self.stats.removed_lines_processed += 1

            # Detect context switches (struct vs interface)
            self._update_context(line)

            # Parse changes based on context
            self._parse_line_change(line, i, current_file, model_changes)

        # Calculate final statistics
        result = [change for change in model_changes.values() if change.has_changes]
        files_without_changes = [change for change in model_changes.values() if not change.has_changes]
        
        # Update statistics
        self.stats.files_with_changes = len(result)
        self.stats.files_without_changes = len(files_without_changes)
        
        for change in result:
            self.stats.struct_fields_added += len(change.added_fields)
            self.stats.struct_fields_removed += len(change.removed_fields)
            self.stats.interface_methods_added += len(change.added_methods)
            self.stats.interface_methods_removed += len(change.removed_methods)
            self.stats.embedded_types_added += len(change.added_embedded_types)
            self.stats.embedded_types_removed += len(change.removed_embedded_types)
        
        self.reporter.print_parse_summary(result, self.stats, files_without_changes)
        
        return result

    def _update_context(self, line: str):
        """Update parsing context based on type declarations.
        
        Args:
            line: Current line being processed
        """
        cleaned = line.lstrip('+-').strip()
        
        # Check for interface declaration
        if RegexPatterns.GO_TYPE_INTERFACE.search(cleaned):
            self.in_interface_context = True
        # Check for struct declaration
        elif RegexPatterns.GO_TYPE_STRUCT.search(cleaned):
            self.in_interface_context = False
        # Closing brace resets context
        elif cleaned == '}':
            self.in_interface_context = False

    def _parse_line_change(self, line: str, line_number: int, current_file: str,
                          model_changes: Dict[str, ModelChange]):
        """Parse a line for any type of change (field, method, embedded type).
        
        Args:
            line: Line from diff
            line_number: Line number in diff
            current_file: Current file being processed
            model_changes: Dictionary of model changes being built
        """
        if line.startswith('+') and not line.startswith('+++'):
            self._parse_added_line(line[1:].strip(), line_number, current_file, model_changes)
        elif line.startswith('-') and not line.startswith('---'):
            self._parse_removed_line(line[1:].strip(), line_number, current_file, model_changes)

    def _parse_added_line(self, line: str, line_number: int, current_file: str,
                         model_changes: Dict[str, ModelChange]):
        """Parse an added line (+).
        
        Args:
            line: Cleaned line content
            line_number: Line number in diff
            current_file: Current file
            model_changes: Model changes dictionary
        """
        # Try embedded type first (works for both interface and struct)
        embedded_info = self._parse_embedded_type(line)
        if embedded_info:
            context = 'interface' if self.in_interface_context else 'struct'
            embedded_change = EmbeddedTypeChange(
                type_name=embedded_info,
                change_type='added',
                context=context,
                line_number=line_number
            )
            model_changes[current_file].added_embedded_types.append(embedded_change)
            self.stats.embedded_types_added += 1
            return

        # If in interface, try to parse as method
        if self.in_interface_context:
            method_info = self._parse_interface_method(line)
            if method_info:
                method_change = MethodChange(
                    method_name=method_info[0],
                    parameters=method_info[1],
                    return_type=method_info[2],
                    change_type='added',
                    line_number=line_number
                )
                model_changes[current_file].added_methods.append(method_change)
                self.stats.interface_methods_added += 1
                return

        # Otherwise, try to parse as struct field
        field_info = self._parse_field_line(line)
        if field_info:
            field_change = FieldChange(
                field_name=field_info[0],
                field_type=field_info[1],
                change_type='added',
                line_number=line_number
            )
            model_changes[current_file].added_fields.append(field_change)
            self.stats.struct_fields_added += 1

    def _parse_removed_line(self, line: str, line_number: int, current_file: str,
                           model_changes: Dict[str, ModelChange]):
        """Parse a removed line (-).
        
        Args:
            line: Cleaned line content
            line_number: Line number in diff
            current_file: Current file
            model_changes: Model changes dictionary
        """
        # Try embedded type first
        embedded_info = self._parse_embedded_type(line)
        if embedded_info:
            context = 'interface' if self.in_interface_context else 'struct'
            embedded_change = EmbeddedTypeChange(
                type_name=embedded_info,
                change_type='removed',
                context=context,
                line_number=line_number
            )
            model_changes[current_file].removed_embedded_types.append(embedded_change)
            self.stats.embedded_types_removed += 1
            return

        # If in interface, try to parse as method
        if self.in_interface_context:
            method_info = self._parse_interface_method(line)
            if method_info:
                method_change = MethodChange(
                    method_name=method_info[0],
                    parameters=method_info[1],
                    return_type=method_info[2],
                    change_type='removed',
                    line_number=line_number
                )
                model_changes[current_file].removed_methods.append(method_change)
                self.stats.interface_methods_removed += 1
                return

        # Otherwise, try to parse as struct field
        field_info = self._parse_field_line(line)
        if field_info:
            field_change = FieldChange(
                field_name=field_info[0],
                field_type=field_info[1],
                change_type='removed',
                line_number=line_number
            )
            model_changes[current_file].removed_fields.append(field_change)
            self.stats.struct_fields_removed += 1

    def _is_file_header(self, line: str) -> bool:
        """Check if line is a file header."""
        return line.startswith('diff --git') or line.startswith('+++') or line.startswith('---')

    def _extract_filename(self, line: str) -> Optional[str]:
        """Extract filename from diff header line."""
        match = RegexPatterns.MODEL_FILE_PATH.search(line)
        return match.group(0) if match else None

    def _filename_to_model_name(self, filename: str) -> str:
        """Convert filename to model name (snake_case to PascalCase)."""
        file_stem = Path(filename).stem
        return ''.join(word.capitalize() for word in file_stem.split('_'))

    def _parse_interface_method(self, line: str) -> Optional[Tuple[str, str, str]]:
        """Parse an interface method declaration.
        
        Args:
            line: Line of Go code
            
        Returns:
            Tuple of (method_name, parameters, return_type) or None
        """
        if not line or line.startswith('//') or line.startswith('}') or line.startswith('{'):
            return None
        
        # Skip function implementations (have body indicators)
        if line.startswith('func (') and '{' in line:
            self.stats.lines_filtered_func_impl += 1
            return None

        match = RegexPatterns.GO_INTERFACE_METHOD.match(line)
        if match:
            method_name = match.group(1)
            parameters = match.group(2) if match.group(2) else ""
            return_type = match.group(3) if match.group(3) else ""
            
            # Only track exported methods (uppercase first letter)
            if method_name and method_name[0].isupper():
                return (method_name, parameters.strip(), return_type.strip().strip('()'))
            else:
                self.stats.lines_filtered_unexported += 1
        
        return None

    def _parse_embedded_type(self, line: str) -> Optional[str]:
        """Parse an embedded type (interface or struct).
        
        Args:
            line: Line of Go code
            
        Returns:
            Type name or None
        """
        if not line or line.startswith('//') or line.startswith('}') or line.startswith('{'):
            return None
        
        # Skip type declarations and function implementations
        if line.startswith('type ') or line.startswith('func '):
            return None

        match = RegexPatterns.GO_EMBEDDED_TYPE.match(line)
        if match:
            type_name = match.group(1)
            # Check if it looks like a type name (not a field with type)
            # Embedded types are just the type name, no field name before it
            if type_name and (type_name[0].isupper() or '.' in type_name):
                return type_name
        
        return None

    def _parse_field_line(self, line: str) -> Optional[Tuple[str, str]]:
        """Parse a Go struct field line and track filtering reasons.
        
        Args:
            line: Line of Go code
            
        Returns:
            Tuple of (field_name, field_type) or None
        """
        if not line:
            return None
            
        # Track why lines are filtered
        if line.startswith('//'):
            self.stats.lines_filtered_comments += 1
            return None
        
        if line.startswith('type ') or line.startswith('package ') or line.startswith('import '):
            self.stats.lines_filtered_declarations += 1
            return None
        
        if line.startswith('func '):
            self.stats.lines_filtered_func_impl += 1
            return None
            
        if line.startswith('}') or line.startswith('{'):
            return None

        match = RegexPatterns.GO_STRUCT_FIELD.match(line)
        if match:
            field_name = match.group(1)
            field_type = match.group(2)

            if field_name[0].isupper():  # Go exported field
                return (field_name, field_type)
            else:
                # Unexported field (starts with lowercase)
                self.stats.lines_filtered_unexported += 1
                return None
        else:
            # Line didn't match the field pattern
            self.stats.lines_filtered_no_match += 1

        return None


class IssueBuilder:
    """Builds GitHub issue content from analysis results."""

    def __init__(self, sdk_repo: str):
        """Initialize issue builder.
        
        Args:
            sdk_repo: SDK repository (owner/repo)
        """
        self.sdk_repo = sdk_repo

    def build_title(self, old_version: str, new_version: str) -> str:
        """Build issue title.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            
        Returns:
            Issue title
        """
        return f"Schema Update Required: Microsoft Graph SDK {old_version} → {new_version}"

    def build_body(self, old_version: str, new_version: str, 
                  model_changes: List[ModelChange], changelog_section: str) -> str:
        """Build complete issue body.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            model_changes: List of model changes
            changelog_section: Changelog text
            
        Returns:
            Issue body markdown
        """
        sections = [
            self._build_summary(old_version, new_version, model_changes),
            self._build_changed_models(model_changes),
            self._build_action_required(),
            self._build_references(old_version, new_version),
            self._build_changelog(changelog_section),
            self._build_footer()
        ]
        
        return '\n'.join(sections)

    def _build_summary(self, old_version: str, new_version: str, 
                      model_changes: List[ModelChange]) -> str:
        """Build summary section."""
        return f"""## Summary
The Microsoft Graph Beta SDK has been updated from `{old_version}` to `{new_version}`.
This update includes {len(model_changes)} model(s) with schema changes that require review and potential Terraform schema updates.
"""

    def _build_changed_models(self, model_changes: List[ModelChange]) -> str:
        """Build changed models section."""
        parts = ["## Changed Models", ""]
        
        for change in model_changes:
            parts.append(f"### `{change.model_name}` ({change.file_path})")
            parts.append(f"**Changes:** {change.change_summary}")
            parts.append("")

            # Struct fields
            if change.added_fields:
                parts.append("**Added Struct Fields:**")
                for fld in change.added_fields:
                    parts.append(f"- `{fld.field_name}` ({fld.field_type})")
                parts.append("")

            if change.removed_fields:
                parts.append("**Removed Struct Fields:**")
                for fld in change.removed_fields:
                    parts.append(f"- `{fld.field_name}` ({fld.field_type})")
                parts.append("")

            # Interface methods
            if change.added_methods:
                parts.append("**Added Interface Methods:**")
                for method in change.added_methods:
                    parts.append(f"- `{method.signature}`")
                parts.append("")

            if change.removed_methods:
                parts.append("**Removed Interface Methods:**")
                for method in change.removed_methods:
                    parts.append(f"- `{method.signature}`")
                parts.append("")

            # Embedded types
            if change.added_embedded_types:
                parts.append("**Added Embedded Types:**")
                for emb in change.added_embedded_types:
                    parts.append(f"- `{emb.type_name}` ({emb.context})")
                parts.append("")

            if change.removed_embedded_types:
                parts.append("**Removed Embedded Types:**")
                for emb in change.removed_embedded_types:
                    parts.append(f"- `{emb.type_name}` ({emb.context})")
                parts.append("")

        return '\n'.join(parts)

    def _build_action_required(self) -> str:
        """Build action required section."""
        return """## Action Required

1. Review each changed model listed above
2. For struct field changes:
   - Update corresponding Terraform resource schemas
   - Add/update field mappings in CRUD operations
3. For interface method changes:
   - Review API contract changes
   - Update method calls if signatures changed
   - Verify compatibility with existing code
4. For embedded type changes:
   - Review inheritance/composition changes
   - Check for breaking changes in type hierarchy
5. Add/update tests for all changes
6. Update documentation

⚠️ **Interface method changes may indicate breaking API changes!**
"""

    def _build_references(self, old_version: str, new_version: str) -> str:
        """Build references section."""
        return f"""## References

- [SDK Changelog](https://github.com/{self.sdk_repo}/blob/main/CHANGELOG.md)
- [Version Diff](https://github.com/{self.sdk_repo}/compare/{old_version}...{new_version})
- [Models Diff](https://github.com/{self.sdk_repo}/compare/{old_version}...{new_version}#files_bucket)
"""

    def _build_changelog(self, changelog_section: str) -> str:
        """Build changelog section."""
        if not changelog_section or "not found" in changelog_section.lower():
            return ""

        parts = ["## Changelog Excerpt", "", "```"]
        parts.append(changelog_section[:1000])
        if len(changelog_section) > 1000:
            parts.append("... (truncated)")
        parts.append("```")
        parts.append("")
        
        return '\n'.join(parts)

    def _build_footer(self) -> str:
        """Build footer."""
        timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        return f"---\n🤖 Auto-generated by kiota_graph_sdk_schema_change_detector.py on {timestamp}"


class KiotaGraphSdkSchemaChangeDetector:
    """Main orchestrator for SDK schema change detection."""

    SDK_MODULE = "github.com/microsoftgraph/msgraph-beta-sdk-go"
    SDK_REPO = "microsoftgraph/msgraph-beta-sdk-go"
    REQUIRED_LABELS = [
        ("sdk-update", "0E8A16", "Microsoft Graph SDK update"),
        ("schema-change", "D93F0B", "Schema changes detected"),
        ("needs-review", "FBCA04", "Requires engineer review"),
    ]

    def __init__(self, repo: Optional[str] = None, dry_run: bool = False, 
                 save_diff: bool = False, verbose: bool = True):
        """Initialize the detector.
        
        Args:
            repo: Target repository in owner/repo format
            dry_run: If True, don't create actual GitHub issues
            save_diff: If True, save fetched diff to file for debugging
            verbose: If True, show detailed progress
        """
        self.dry_run = dry_run
        self.go_mod_path = Path.cwd() / "go.mod"
        
        # Initialize components
        self.reporter = ProgressReporter(verbose)
        self.version_parser = VersionParser()
        self.github_client = GitHubClient(
            repo or self.github_client.get_current_repo() if hasattr(self, 'github_client') 
            else GitHubClient.get_current_repo(None),
            self.reporter
        )
        # Fix circular dependency
        if not repo:
            repo = self.github_client.get_current_repo()
        self.github_client = GitHubClient(repo, self.reporter)
        
        self.diff_fetcher = DiffFetcher(self.SDK_REPO, self.github_client, self.reporter, save_diff)
        self.struct_parser = StructParser(self.reporter)
        self.issue_builder = IssueBuilder(self.SDK_REPO)

    def parse_go_mod_version(self) -> Optional[str]:
        """Parse current SDK version from go.mod.
        
        Returns:
            Current version string or None
        """
        if not self.go_mod_path.exists():
            self.reporter.error(f"go.mod not found at {self.go_mod_path}")
            return None

        with open(self.go_mod_path, 'r', encoding='utf-8') as f:
            for line in f:
                if self.SDK_MODULE in line:
                    version = self.version_parser.extract_version_from_line(line)
                    if version:
                        return version

        self.reporter.error(f"Could not find {self.SDK_MODULE} in go.mod")
        return None

    def get_pr_version_change(self, pr_number: int) -> VersionChangeResult:
        """Get version change from a PR.
        
        Args:
            pr_number: PR number
            
        Returns:
            VersionChangeResult with old and new versions
        """
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

        except subprocess.CalledProcessError as e:
            error_msg = f"Error getting PR version change: {e}"
            self.reporter.error(error_msg)
            return VersionChangeResult.error(error_msg)

    def analyze_pr(self, pr_number: int, save_results: Optional[Path] = None,
                   filter_by_usage: Optional[Path] = None) -> bool:
        """Analyze a PR for SDK schema changes.
        
        Args:
            pr_number: PR number to analyze
            save_results: Optional path to save detection results JSON
            filter_by_usage: Optional path to model usage JSON for filtering
            
        Returns:
            True if analysis completed successfully
        """
        self.reporter.section(f"🔍 Analyzing PR #{pr_number}...")

        version_change = self.get_pr_version_change(pr_number)
        if not version_change.is_success:
            self.reporter.error(version_change.error_message or "Could not detect version change in PR")
            return False

        return self.analyze_version_change(version_change.old_version, version_change.new_version,
                                          save_results, filter_by_usage)

    def analyze_version_change(self, old_version: str, new_version: str, 
                              save_results: Optional[Path] = None,
                              filter_by_usage: Optional[Path] = None) -> bool:
        """Analyze a version change for schema updates.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            save_results: Optional path to save detection results JSON
            filter_by_usage: Optional path to model usage JSON for filtering
            
        Returns:
            True if analysis completed successfully
        """
        self.reporter.section(f"📊 Analyzing version change: {old_version} → {new_version}")

        # Validate version increment
        validation = self.version_parser.validate_increment(old_version, new_version, self.reporter)
        if not validation.is_valid:
            self.reporter.warning("Version increment validation failed, but continuing analysis...")

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
            self.reporter.info(f"🔍 Filtered to {len(model_changes)} relevant model(s) "
                             f"({filtered_count} filtered out)")

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
        """Filter model changes to only those used in the provider.
        
        Args:
            model_changes: List of detected model changes
            usage_file: Path to model usage JSON file
            
        Returns:
            Tuple of (filtered_changes, filtered_count)
        """
        self.reporter.section("🔍 Filtering changes by provider usage...")
        
        try:
            with open(usage_file, 'r', encoding='utf-8') as f:
                usage_data = json.load(f)
        except (OSError, json.JSONDecodeError) as e:
            self.reporter.error(f"Failed to load usage data: {e}")
            return model_changes, 0
        
        # Build set of used model files
        used_models = {model['model_file'] for model in usage_data.get('models', [])}
        
        self.reporter.info(f"  Loaded {len(used_models)} model(s) from provider usage data")
        
        # Filter changes
        filtered = [mc for mc in model_changes if mc.file_path in used_models]
        filtered_out = len(model_changes) - len(filtered)
        
        if filtered_out > 0:
            self.reporter.info(f"  ✓ Kept {len(filtered)} relevant model(s)")
            self.reporter.info(f"  ✗ Filtered {filtered_out} unused model(s)")
        
        return filtered, filtered_out
    
    def _save_results(self, old_version: str, new_version: str, 
                     model_changes: List[ModelChange], changelog: str,
                     total_models: int, filtered_count: int, 
                     output_file: Path) -> None:
        """Save detection results to JSON file.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            model_changes: List of detected model changes
            changelog: Changelog section
            total_models: Total models before filtering
            filtered_count: Number of models filtered out
            output_file: Path to save results
        """
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
        """Create GitHub issue for schema updates.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            model_changes: List of model changes
            changelog: Changelog section
            
        Returns:
            IssueCreationResult with creation status
        """
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

        except (subprocess.CalledProcessError, IndexError) as e:
            error_msg = f"Error creating issue: {e}"
            self.reporter.error(error_msg)
            return IssueCreationResult.error(error_msg)


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
        "--save-diff",
        action="store_true",
        help="Save fetched diff to file for debugging"
    )
    parser.add_argument(
        "--repo",
        type=str,
        help="Repository in owner/repo format (auto-detected if not provided)"
    )
    parser.add_argument(
        "--save-results",
        type=Path,
        help="Save detection results to JSON file"
    )
    parser.add_argument(
        "--filter-by-usage",
        type=Path,
        help="Filter changes by provider model usage (path to usage JSON)"
    )

    args = parser.parse_args()

    # Validate arguments
    if not args.pr_number and not (args.current and args.new):
        parser.error("Either --pr-number or both --current and --new must be provided")

    # Create detector
    detector = KiotaGraphSdkSchemaChangeDetector(
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
    except (subprocess.CalledProcessError, OSError, ValueError, RuntimeError) as e:
        print(f"❌ Unexpected error: {e}", file=sys.stderr)
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
