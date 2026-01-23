#!/usr/bin/env python3
"""Map SDK usage across the Terraform provider codebase.

Outputs:
- usage-file: Path to JSON file containing SDK usage data
"""

import argparse
import json
import sys
from pathlib import Path

# Add lib directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "lib"))
# noqa: E402
from go_parser import extract_sdk_usage  # pylint: disable=import-error


def main():
    """Map SDK usage across provider codebase and output statistics."""
    parser = argparse.ArgumentParser(description="Map SDK usage in provider")
    parser.add_argument(
        "--repo-path",
        type=Path,
        default=Path.cwd(),
        help="Path to repository root"
    )
    parser.add_argument(
        "--output-file",
        help="GitHub Actions output file"
    )
    parser.add_argument(
        "--usage-output",
        type=Path,
        default=Path.cwd() / "sdk_usage.json",
        help="Path to save usage data JSON"
    )
    
    args = parser.parse_args()
    
    print("=" * 60)
    print("📊 Analyzing SDK Usage")
    print("=" * 60)
    
    # Extract usage
    usage_data = extract_sdk_usage(args.repo_path)
    
    # Save to file
    with open(args.usage_output, 'w', encoding='utf-8') as f:
        json.dump(usage_data, f, indent=2)
    
    print(f"\n💾 Usage data saved: {args.usage_output}")
    
    # Statistics
    stats = usage_data.get('statistics', {})
    print("\n📊 Terraform Entities:")
    print(f"  - Resources:     {stats.get('total_resources', 0)}")
    print(f"  - Actions:       {stats.get('total_actions', 0)}")
    print(f"  - List Actions:  {stats.get('total_list_actions', 0)}")
    print(f"  - Ephemerals:    {stats.get('total_ephemerals', 0)}")
    print(f"  - Data Sources:  {stats.get('total_data_sources', 0)}")
    
    print("\n📈 SDK Usage:")
    print(f"  - SDK types used:    {stats.get('total_sdk_types_used', 0)}")
    print(f"  - SDK methods used:  {stats.get('total_sdk_methods_used', 0)}")
    print(f"  - Enums tracked:     {stats.get('total_enums_tracked', 0)}")
    
    # Show sample resources
    resources = usage_data.get('terraform_resources', {})
    if resources:
        print("\n🔍 Sample Resources:")
        for i, (name, info) in enumerate(list(resources.items())[:5], 1):
            type_count = len(info.get('sdk_dependencies', {}).get('types', []))
            print(f"  {i}. {name} ({type_count} SDK types)")
    
    # Write outputs
    if args.output_file:
        with open(args.output_file, 'a', encoding='utf-8') as f:
            f.write(f"usage-file={args.usage_output.absolute()}\n")
            f.write(f"resources-count={stats.get('total_resources', 0)}\n")
            f.write(f"types-count={stats.get('total_sdk_types_used', 0)}\n")
    
    print("\n✅ Usage mapping complete")


if __name__ == "__main__":
    main()
