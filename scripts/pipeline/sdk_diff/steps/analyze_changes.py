#!/usr/bin/env python3
"""Analyze SDK changes between versions.

Requires:
- Usage data from map_usage.py
- Version information from detect_versions.py

Outputs:
- changes-file: Path to JSON file containing categorized changes
- critical-count: Number of critical changes
- warning-count: Number of warnings
- safe-count: Number of safe changes
"""

import argparse
import json
import sys
import urllib.error
from pathlib import Path

# Add lib directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "lib"))
# noqa: E402
from diff_analyzer import ChangeAnalyzer, generate_summary_stats, ImpactLevel # pylint: disable=import-error
from github_api import compare_versions, get_sdk_repo_name, parse_breaking_changes, get_latest_release # pylint: disable=import-error


def parse_arguments():
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(description="Analyze SDK changes")
    parser.add_argument(
        "--usage-file",
        type=Path,
        required=True,
        help="Path to usage data JSON"
    )
    parser.add_argument(
        "--sdk",
        required=True,
        choices=["msgraph-sdk-go", "msgraph-beta-sdk-go"],
        help="Which SDK to analyze"
    )
    parser.add_argument(
        "--current-version",
        required=True,
        help="Current version (e.g., v1.93.0)"
    )
    parser.add_argument(
        "--latest-version",
        required=True,
        help="Latest version to compare against"
    )
    parser.add_argument(
        "--output-file",
        help="GitHub Actions output file"
    )
    parser.add_argument(
        "--changes-output",
        type=Path,
        default=Path.cwd() / "sdk_changes.json",
        help="Path to save changes data JSON"
    )
    return parser.parse_args()


def get_version_comparison(repo, current_version, latest_version):
    """Get version comparison from GitHub API."""
    print(f"📊 Getting changes from GitHub ({repo})...")
    try:
        comparison = compare_versions(repo, current_version, latest_version)
        print(f"✅ Found {comparison['commits']} commits, {comparison['files_changed']} files changed")
        return comparison
    except (RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError) as e:
        print(f"❌ Failed to compare versions: {e}")
        sys.exit(1)


def get_breaking_changes(repo):
    """Get breaking changes from release notes."""
    print("\n📝 Checking release notes for breaking changes...")
    try:
        release = get_latest_release(repo)
        breaking_changes = parse_breaking_changes(release.get("body", ""))
        if breaking_changes:
            print(f"⚠️  Found {len(breaking_changes)} breaking changes in release notes")
        else:
            print("✅ No breaking changes mentioned in release notes")
        return breaking_changes
    except (RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError) as e:
        print(f"⚠️  Could not get release notes: {e}")
        return []


def analyze_field_additions(analyzer, repo, current_version, latest_version, categorized):
    """Analyze field additions in used types."""
    print("\n🔍 Analyzing field additions in used types...")
    try:
        field_opportunities = analyzer.analyze_field_additions(
            repo,
            current_version,
            latest_version
        )
        if field_opportunities:
            print(f"✨ Found {len(field_opportunities)} new field(s) in types you use")
            for opp in field_opportunities:
                categorized[ImpactLevel.OPPORTUNITY].append({
                    'type': opp['type'],
                    'field': opp['field'],
                    'field_type': opp['field_type'],
                    'file': opp['file'],
                    'description': opp['description'],
                    'reason': f"New field '{opp['field']}' added to {opp['type']}"
                })
        else:
            print("ℹ️  No new fields in types you currently use")
    except (RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError, TypeError) as e:
        print(f"⚠️  Could not analyze field additions: {e}")


def analyze_enum_changes(analyzer, repo, current_version, latest_version, categorized):
    """Analyze enum value changes in used enums."""
    print("\n🔍 Analyzing enum value changes...")
    try:
        enum_changes = analyzer.analyze_enum_changes(
            repo,
            current_version,
            latest_version
        )
        if enum_changes:
            added_count = sum(1 for e in enum_changes if e['added_values'])
            removed_count = sum(1 for e in enum_changes if e['removed_values'])
            
            if added_count:
                print(f"✨ Found {added_count} enum(s) with new values")
            if removed_count:
                print(f"🚨 Found {removed_count} enum(s) with removed values (breaking!)")
            
            for enum_change in enum_changes:
                if enum_change['removed_values']:
                    # Removed values are CRITICAL
                    categorized[ImpactLevel.ENUM_REMOVED].append({
                        'enum_type': enum_change['enum_type'],
                        'removed_values': enum_change['removed_values'],
                        'added_values': enum_change['added_values'],
                        'file': enum_change['file'],
                        'reason': f"Enum values removed: {', '.join(enum_change['removed_values'])}"
                    })
                elif enum_change['added_values']:
                    # Added values are opportunities
                    categorized[ImpactLevel.ENUM_ADDED].append({
                        'enum_type': enum_change['enum_type'],
                        'added_values': enum_change['added_values'],
                        'file': enum_change['file'],
                        'reason': f"New enum values: {', '.join(enum_change['added_values'])}"
                    })
        else:
            print("ℹ️  No enum value changes in enums you use")
    except (RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError, TypeError) as e:
        print(f"⚠️  Could not analyze enum changes: {e}")


def print_analysis_results(stats):
    """Print analysis results to console."""
    print("\n📊 Analysis Results:")
    print(f"  Total changes: {stats['total_changes']:,}")
    relevant_pct = stats['relevant_changes'] / max(stats['total_changes'], 1) * 100
    print(f"  Relevant:      {stats['relevant_changes']} ({relevant_pct:.1f}%)")
    print(f"    🚨 Critical: {stats['critical']}")
    print(f"    ⚠️  Warning:  {stats['warning']}")
    print(f"    ✅ Safe:     {stats['safe']}")
    print(f"    ✨ Opportunities: {stats['opportunity']}")
    print(f"    🎯 Enum values added: {stats.get('enum_added', 0)}")
    print(f"    🔥 Enum values removed: {stats.get('enum_removed', 0)}")
    print(f"  Noise:         {stats['noise']:,}")


def save_results(args, comparison, breaking_changes, categorized, stats):
    """Save analysis results to files."""
    output_data = {
        "sdk": args.sdk,
        "versions": {
            "current": args.current_version,
            "latest": args.latest_version
        },
        "comparison_url": comparison["url"],
        "breaking_changes": breaking_changes,
        "categorized_changes": categorized,
        "statistics": stats
    }
    
    with open(args.changes_output, 'w', encoding='utf-8') as f:
        json.dump(output_data, f, indent=2)
    
    print(f"\n💾 Changes data saved: {args.changes_output}")
    
    if args.output_file:
        with open(args.output_file, 'a', encoding='utf-8') as f:
            f.write(f"changes-file={args.changes_output.absolute()}\n")
            f.write(f"critical-count={stats['critical']}\n")
            f.write(f"warning-count={stats['warning']}\n")
            f.write(f"safe-count={stats['safe']}\n")
            has_breaking = 'true' if stats['critical'] > 0 else 'false'
            f.write(f"has-breaking-changes={has_breaking}\n")


def main():
    """Analyze SDK changes between versions and categorize by impact.
    
    Loads usage data, compares SDK versions via GitHub API, analyzes file
    changes and field additions, then outputs categorized results for reporting.
    """
    args = parse_arguments()
    
    print("=" * 60)
    print(f"🔬 Analyzing {args.sdk} Changes")
    print("=" * 60)
    print(f"Version: {args.current_version} → {args.latest_version}")
    print()
    
    print(f"📂 Loading usage data from {args.usage_file}...")
    with open(args.usage_file, 'r', encoding='utf-8') as f:
        usage_data = json.load(f)
    
    repo = get_sdk_repo_name(args.sdk)
    comparison = get_version_comparison(repo, args.current_version, args.latest_version)
    breaking_changes = get_breaking_changes(repo)
    
    print("\n🔍 Analyzing impact of file changes...")
    analyzer = ChangeAnalyzer(usage_data)
    categorized = analyzer.analyze_file_changes(comparison["files"])
    
    analyze_field_additions(analyzer, repo, args.current_version, args.latest_version, categorized)
    analyze_enum_changes(analyzer, repo, args.current_version, args.latest_version, categorized)
    
    stats = generate_summary_stats(categorized)
    print_analysis_results(stats)
    save_results(args, comparison, breaking_changes, categorized, stats)
    
    print("\n✅ Change analysis complete")


if __name__ == "__main__":
    main()
