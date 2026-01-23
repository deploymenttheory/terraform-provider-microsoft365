#!/usr/bin/env python3
"""Generate impact report from analysis data.

Outputs:
- Markdown report
- JSON report (optional)
"""

import argparse
import json
import sys
from pathlib import Path

# Add lib directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "lib"))
# noqa: E402
from report_generator import generate_markdown_report, save_json_report  # pylint: disable=import-error


def main():
    """Generate Markdown and optional JSON reports from SDK change analysis data."""
    parser = argparse.ArgumentParser(description="Generate SDK diff report")
    parser.add_argument(
        "--changes-file",
        type=Path,
        required=True,
        help="Path to changes data JSON from analyze_changes.py"
    )
    parser.add_argument(
        "--current-versions",
        type=json.loads,
        required=True,
        help='Current versions as JSON, e.g. \'{"msgraph-sdk-go": "v1.93.0"}\''
    )
    parser.add_argument(
        "--latest-versions",
        type=json.loads,
        required=True,
        help='Latest versions as JSON, e.g. \'{"msgraph-sdk-go": "v1.95.0"}\''
    )
    parser.add_argument(
        "--output-markdown",
        type=Path,
        default=Path.cwd() / "SDK_DIFF_REPORT.md",
        help="Path to output Markdown report"
    )
    parser.add_argument(
        "--output-json",
        type=Path,
        help="Optional path to output JSON report"
    )
    
    args = parser.parse_args()
    
    print("=" * 60)
    print("📝 Generating SDK Diff Report")
    print("=" * 60)
    
    # Load changes data
    print(f"📂 Loading changes data from {args.changes_file}...")
    with open(args.changes_file, 'r', encoding='utf-8') as f:
        changes_data = json.load(f)
    
    # Generate Markdown report
    print("\n📄 Generating Markdown report...")
    markdown = generate_markdown_report(
        current_versions=args.current_versions,
        latest_versions=args.latest_versions,
        categorized_changes=changes_data["categorized_changes"],
        stats=changes_data["statistics"],
        breaking_changes=changes_data.get("breaking_changes", [])
    )
    
    # Save Markdown
    with open(args.output_markdown, 'w', encoding='utf-8') as f:
        f.write(markdown)
    
    print(f"✅ Markdown report saved: {args.output_markdown}")
    
    # Save JSON if requested
    if args.output_json:
        print("\n💾 Saving JSON report...")
        save_json_report(changes_data, args.output_json)
    
    # Display summary
    stats = changes_data["statistics"]
    print("\n📊 Report Summary:")
    print(f"  🚨 Critical issues: {stats['critical']}")
    print(f"  ⚠️  Warnings: {stats['warning']}")
    print(f"  ✅ Safe changes: {stats['safe']}")
    
    if stats['critical'] > 0:
        print("\n⚠️  WARNING: Critical changes detected! Review required before upgrading.")
        sys.exit(1)
    elif stats['warning'] > 0:
        print("\n⚠️  Warnings detected. Review recommended.")
    else:
        print("\n✅ No breaking changes detected!")
    
    print("\n✅ Report generation complete")


if __name__ == "__main__":
    main()
