#!/usr/bin/env python3
"""Generate impact reports from analysis data.

Outputs:
- Multiple structured Markdown reports
- Actionable JSON for PR automation
- Usage baseline JSON
"""

import argparse
import json
import sys
from pathlib import Path

# Add lib directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "lib"))
# noqa: E402
# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from report_generator import ReportInputs, generate_all_reports  # pylint: disable=import-error

# pylint: enable=wrong-import-position


def main():
    """Generate all SDK diff reports from analysis data."""
    parser = argparse.ArgumentParser(description="Generate SDK diff reports")
    parser.add_argument(
        "--changes-file",
        type=Path,
        required=True,
        help="Path to changes data JSON from analyze_changes.py",
    )
    parser.add_argument(
        "--usage-file",
        type=Path,
        required=True,
        help="Path to SDK usage baseline JSON from map_usage.py",
    )
    parser.add_argument(
        "--output-dir",
        type=Path,
        default=Path.cwd(),
        help="Directory to output all reports (default: current directory)",
    )

    args = parser.parse_args()

    print("=" * 60)
    print("📝 Generating SDK Diff Reports")
    print("=" * 60)

    # Load changes data
    print(f"\n📂 Loading changes data from {args.changes_file}...")
    with open(args.changes_file, "r", encoding="utf-8") as f:
        changes_data = json.load(f)

    # Load usage data
    print(f"📂 Loading usage baseline from {args.usage_file}...")
    with open(args.usage_file, "r", encoding="utf-8") as f:
        usage_data = json.load(f)

    # Extract metadata
    sdk_name = changes_data.get("sdk", "unknown")
    versions = changes_data.get("versions", {})
    current_version = versions.get("current", "unknown")
    latest_version = versions.get("latest", "unknown")
    comparison_url = changes_data.get("comparison_url", "")

    # Generate all reports
    print(f"\n📄 Generating reports in {args.output_dir}...")
    generated_files = generate_all_reports(
        ReportInputs(
            metadata={
                "sdk": sdk_name,
                "current_version": current_version,
                "latest_version": latest_version,
                "comparison_url": comparison_url,
            },
            categorized_changes=changes_data["categorized_changes"],
            stats=changes_data["statistics"],
            breaking_changes=changes_data.get("breaking_changes", []),
            usage_data=usage_data,
            output_dir=args.output_dir,
        )
    )

    print(f"\n✅ Generated {len(generated_files)} report files:")
    for file_path in generated_files:
        file_name = Path(file_path).name
        print(f"  - {file_name}")

    # Display summary
    stats = changes_data["statistics"]
    print("\n📊 Report Summary:")
    print(f"  🚨 Critical: {stats.get('critical', 0)}")
    print(f"  🔥 Enum Removed: {stats.get('enum_removed', 0)}")
    print(f"  ⚠️  Warnings: {stats.get('warning', 0)}")
    print(f"  ✅ Safe: {stats.get('safe', 0)}")
    print(f"  ✨ Opportunities: {stats.get('opportunity', 0)}")
    print(f"  🎯 Enum Added: {stats.get('enum_added', 0)}")
    print(f"  ⚪ Metadata: {stats.get('metadata', 0)}")
    print(f"  🗑️  Filtered: {stats.get('noise', 0)}")

    # Determine exit status
    critical_count = stats.get("critical", 0) + stats.get("enum_removed", 0)

    if critical_count > 0:
        print("\n⚠️  CRITICAL: Breaking changes detected! Review 02_breaking_changes.md")
        sys.exit(1)
    elif stats.get("warning", 0) > 0:
        print("\n⚠️  Warnings detected. Review 03_provider_updates_required.md")
    else:
        print("\n✅ No breaking changes detected!")

    print("\n✅ Report generation complete")
    print(f"\n📁 All reports saved to: {args.output_dir}")


if __name__ == "__main__":
    main()
