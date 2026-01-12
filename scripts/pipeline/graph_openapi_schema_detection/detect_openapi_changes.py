#!/usr/bin/env python3
"""
Detect OpenAPI schema changes and create GitHub issues.

This script analyzes Microsoft Graph OpenAPI specification changes and creates
GitHub issues for schema updates that affect the Terraform provider.

Usage:
    # Analyze latest spec vs cached version
    ./detect_openapi_changes.py

    # Compare specific commits
    ./detect_openapi_changes.py --old-commit abc123 --new-commit def456

    # Dry run (don't create issues)
    ./detect_openapi_changes.py --dry-run

    # Filter by provider usage
    ./detect_openapi_changes.py --filter-by-usage ../sdk_schema_detection/provider_model_usage.json
"""

import sys
import json
import argparse
import subprocess
import traceback
from pathlib import Path
from datetime import datetime
from typing import Optional, List
from urllib.parse import urlparse

# Add current directory to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from core import (  # pylint: disable=import-error
    ProgressReporter,
    SpecFetcher,
    VersionDetector,
    SchemaParser,
    SchemaComparer,
    IssueBuilder,
)
from models import (  # pylint: disable=import-error
    DetectionResult,
    SchemaChange,
)


class OpenAPIChangeDetector:
    """Main detector for OpenAPI schema changes."""

    OPENAPI_REPO = "microsoftgraph/msgraph-metadata"
    REQUIRED_LABELS = [
        ("sdk-update", "0E8A16", "Microsoft Graph SDK update"),
        ("schema-change", "D93F0B", "Schema changes detected"),
        ("needs-review", "FBCA04", "Requires engineer review"),
    ]

    def __init__(self, repo: Optional[str] = None, dry_run: bool = False, verbose: bool = True):
        """Initialize detector.
        
        Args:
            repo: Target repository for issues (owner/repo)
            dry_run: If True, don't create actual GitHub issues
            verbose: If True, show detailed progress
        """
        self.dry_run = dry_run
        self.target_repo = repo or self._get_current_repo_from_git()
        
        self.reporter = ProgressReporter(verbose)
        self.spec_fetcher = SpecFetcher(self.reporter)
        self.version_detector = VersionDetector(self.reporter)
        self.schema_parser = SchemaParser(self.reporter)
        self.schema_comparer = SchemaComparer(self.schema_parser, self.reporter)
        self.issue_builder = IssueBuilder(self.OPENAPI_REPO)

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

    def detect_changes(
        self,
        old_commit: Optional[str] = None,
        new_commit: Optional[str] = None,
        filter_by_usage: Optional[str] = None,
        save_results: Optional[str] = None
    ) -> bool:
        """Detect OpenAPI schema changes.
        
        Args:
            old_commit: Previous commit SHA (if None, fetches latest)
            new_commit: New commit SHA (if None, fetches latest)
            filter_by_usage: Path to provider usage JSON file
            save_results: Path to save detection results JSON
            
        Returns:
            True if successful
        """
        self.reporter.section("OpenAPI Schema Change Detection")
        
        # Fetch specs
        if new_commit:
            new_spec = self.spec_fetcher.fetch_spec_at_commit(new_commit)
        else:
            new_spec = self.spec_fetcher.fetch_latest_spec()
        
        if old_commit:
            old_spec = self.spec_fetcher.fetch_spec_at_commit(old_commit)
        else:
            # For now, require old_commit - in production, we'd cache previous version
            self.reporter.error("--old-commit is required (caching not yet implemented)")
            return False
        
        # Extract versions
        old_version_result = self.version_detector.extract_version_from_spec(old_spec)
        new_version_result = self.version_detector.extract_version_from_spec(new_spec)
        
        if not (old_version_result.is_success and new_version_result.is_success):
            self.reporter.error("Could not extract versions from specs")
            return False
        
        old_version = old_version_result.version
        new_version = new_version_result.version
        
        self.reporter.info(f"📊 Analyzing version change: {old_version} → {new_version}")
        
        # Parse schemas
        old_schemas_content = self.schema_parser.extract_schemas_section(old_spec)
        new_schemas_content = self.schema_parser.extract_schemas_section(new_spec)
        
        old_schemas = self.schema_parser.parse_schemas(old_schemas_content)
        new_schemas = self.schema_parser.parse_schemas(new_schemas_content)
        
        # Compare schemas
        schema_changes, statistics = self.schema_comparer.compare_schemas(old_schemas, new_schemas)
        
        self.reporter.info("")
        self.reporter.info(statistics.get_summary())
        
        # Filter by provider usage
        if filter_by_usage:
            schema_changes = self._filter_by_usage(schema_changes, filter_by_usage)
        
        # Create detection result
        detection_result = DetectionResult(
            spec_version=new_version,
            previous_version=old_version,
            detection_timestamp=datetime.now().isoformat(),
            total_schemas_changed=len(schema_changes),
            schemas_with_changes=statistics.schemas_with_changes,
            filtered_schemas=statistics.schemas_with_changes - len(schema_changes) if filter_by_usage else 0,
            breaking_changes_count=sum(1 for sc in schema_changes if sc.has_breaking_changes),
            schema_changes=schema_changes,
            statistics=statistics
        )
        
        # Save results
        if save_results:
            self._save_results(detection_result, save_results)
        
        # Create GitHub issue
        if schema_changes and not self.dry_run:
            self._create_issue(old_version, new_version, schema_changes)
        elif schema_changes:
            self._show_dry_run_output(old_version, new_version, schema_changes)
        else:
            self.reporter.info("\n✅ No relevant schema changes detected")
        
        return True

    def _filter_by_usage(self, schema_changes: List[SchemaChange], usage_file: str) -> List[SchemaChange]:
        """Filter schema changes to only those used by provider.
        
        Args:
            schema_changes: All detected schema changes
            usage_file: Path to provider_model_usage.json
            
        Returns:
            Filtered list of schema changes
        """
        self.reporter.info(f"\n🔍 Filtering changes by provider usage...")
        
        try:
            with open(usage_file, 'r', encoding='utf-8') as f:
                usage_data = json.load(f)
            
            used_models = {model['model_name'].lower() for model in usage_data.get('models', [])}
            self.reporter.info(f"   Loaded {len(used_models)} model(s) from provider usage data")
            
            # Debug: Show mapping for changed schemas
            self.reporter.info(f"\n   🔍 Mapping changed schemas to provider models:")
            filtered = []
            for sc in schema_changes:
                model_name_lower = sc.model_name.lower()
                is_used = model_name_lower in used_models
                status = "✓ MATCH" if is_used else "✗ not used"
                self.reporter.info(f"      {sc.schema_name} → {sc.model_name} [{status}]")
                if is_used:
                    filtered.append(sc)
            
            self.reporter.info(f"\n   ✓ Kept {len(filtered)} relevant model(s)")
            self.reporter.info(f"   ✗ Filtered {len(schema_changes) - len(filtered)} unused model(s)")
            
            return filtered
            
        except (FileNotFoundError, json.JSONDecodeError, KeyError) as e:
            self.reporter.error(f"Failed to load provider usage data: {e}")
            return schema_changes

    def _save_results(self, result: DetectionResult, output_path: str):
        """Save detection results to JSON file.
        
        Args:
            result: Detection result
            output_path: Output file path
        """
        self.reporter.info(f"\n💾 Saving results to: {output_path}")
        
        try:
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(result.to_dict(), f, indent=2)
            self.reporter.info("   ✓ Results saved successfully")
        except (IOError, OSError) as e:
            self.reporter.error(f"Failed to save results: {e}")

    def _create_issue(self, old_version: str, new_version: str, schema_changes: List[SchemaChange]):
        """Create GitHub issue for schema changes.
        
        Args:
            old_version: Previous version
            new_version: New version
            schema_changes: List of schema changes
        """
        self.reporter.info("\n🎫 Creating GitHub issue...")
        
        title = self.issue_builder.build_title(new_version, old_version)
        body = self.issue_builder.build_body(new_version, old_version, schema_changes)
        
        # TODO: Implement GitHub CLI integration (reuse from SDK detection)
        self.reporter.info(f"   Title: {title}")
        self.reporter.info(f"   Body length: {len(body)} characters")
        self.reporter.info("   ⚠️ GitHub issue creation not yet implemented")

    def _show_dry_run_output(self, old_version: str, new_version: str, schema_changes: List[SchemaChange]):
        """Show what would be created in dry run mode.
        
        Args:
            old_version: Previous version
            new_version: New version
            schema_changes: List of schema changes
        """
        title = self.issue_builder.build_title(new_version, old_version)
        body = self.issue_builder.build_body(new_version, old_version, schema_changes)
        
        self.reporter.info("\n🔍 DRY RUN: Would create issue with following content:")
        self.reporter.info("=" * 80)
        self.reporter.info(f"Title: {title}\n")
        self.reporter.info(body)
        self.reporter.info("=" * 80)


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Detect OpenAPI schema changes and create GitHub issues"
    )
    parser.add_argument(
        '--old-commit',
        help='Previous commit SHA'
    )
    parser.add_argument(
        '--new-commit',
        help='New commit SHA (default: latest)'
    )
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='Analyze without creating issues'
    )
    parser.add_argument(
        '--repo',
        help='Target repository (owner/repo)'
    )
    parser.add_argument(
        '--filter-by-usage',
        help='Filter by provider usage (path to JSON)'
    )
    parser.add_argument(
        '--save-results',
        help='Save results to JSON file'
    )
    parser.add_argument(
        '--verbose',
        action='store_true',
        default=True,
        help='Show detailed progress'
    )
    
    args = parser.parse_args()
    
    try:
        detector = OpenAPIChangeDetector(
            repo=args.repo,
            dry_run=args.dry_run,
            verbose=args.verbose
        )
        
        success = detector.detect_changes(
            old_commit=args.old_commit,
            new_commit=args.new_commit,
            filter_by_usage=args.filter_by_usage,
            save_results=args.save_results
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
