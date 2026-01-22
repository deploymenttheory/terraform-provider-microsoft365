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
from diff_analyzer import ChangeAnalyzer, generate_summary_stats, ImpactLevel
from github_api import compare_versions, get_sdk_repo_name, parse_breaking_changes, get_latest_release


def main():
    """Analyze SDK changes between versions and categorize by impact.
    
    Loads usage data, compares SDK versions via GitHub API, analyzes file
    changes and field additions, then outputs categorized results for reporting.
    """
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
    
    args = parser.parse_args()
    
    print("=" * 60)
    print(f"🔬 Analyzing {args.sdk} Changes")
    print("=" * 60)
    print(f"Version: {args.current_version} → {args.latest_version}")
    print()
    
    # Load usage data
    print(f"📂 Loading usage data from {args.usage_file}...")
    with open(args.usage_file, 'r', encoding='utf-8') as f:
        usage_data = json.load(f)
    
    # Get SDK repo name
    repo = get_sdk_repo_name(args.sdk)
    
    # Compare versions
    print(f"📊 Fetching changes from GitHub ({repo})...")
    try:
        comparison = compare_versions(repo, args.current_version, args.latest_version)
        print(f"✅ Found {comparison['commits']} commits, {comparison['files_changed']} files changed")
    except (RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError) as e:
        print(f"❌ Failed to compare versions: {e}")
        sys.exit(1)
    
    # Get breaking changes from release notes
    print("\n📝 Checking release notes for breaking changes...")
    try:
        release = get_latest_release(repo)
        breaking_changes = parse_breaking_changes(release.get("body", ""))
        if breaking_changes:
            print(f"⚠️  Found {len(breaking_changes)} breaking changes in release notes")
        else:
            print("✅ No breaking changes mentioned in release notes")
    except (RuntimeError, urllib.error.URLError, json.JSONDecodeError, KeyError) as e:
        print(f"⚠️  Could not fetch release notes: {e}")
        breaking_changes = []
    
    # Analyze changes
    print("\n🔍 Analyzing impact of file changes...")
    analyzer = ChangeAnalyzer(usage_data)
    categorized = analyzer.analyze_file_changes(comparison["files"])
    
    # Analyze field additions in types we use
    print("\n🔍 Analyzing field additions in used types...")
    try:
        field_opportunities = analyzer.analyze_field_additions(
            repo, 
            args.current_version, 
            args.latest_version
        )
        if field_opportunities:
            print(f"✨ Found {len(field_opportunities)} new field(s) in types you use")
            # Add to categorized as opportunities
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
    
    # Generate statistics
    stats = generate_summary_stats(categorized)
    
    print("\n📊 Analysis Results:")
    print(f"  Total changes: {stats['total_changes']:,}")
    print(f"  Relevant:      {stats['relevant_changes']} ({stats['relevant_changes'] / max(stats['total_changes'], 1) * 100:.1f}%)")
    print(f"    🚨 Critical: {stats['critical']}")
    print(f"    ⚠️  Warning:  {stats['warning']}")
    print(f"    ✅ Safe:     {stats['safe']}")
    print(f"    ✨ Opportunities: {stats['opportunity']}")
    print(f"  Noise:         {stats['noise']:,}")
    
    # Save results
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
            f.write(f"has-breaking-changes={'true' if stats['critical'] > 0 else 'false'}\n")
    
    print("\n✅ Change analysis complete")


if __name__ == "__main__":
    main()
