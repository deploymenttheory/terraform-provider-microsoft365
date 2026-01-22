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
from go_parser import extract_sdk_usage, get_most_used_packages


def main():
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
    
    # Display most used packages
    print("\n📦 Top 10 Most Used SDK Packages:")
    top_packages = get_most_used_packages(usage_data, top_n=10)
    
    for i, (pkg, count) in enumerate(top_packages, 1):
        # Shorten package names for display
        short_name = pkg.split('/')[-1] if '/' in pkg else pkg
        print(f"{i:2d}. {short_name:40s} ({count:3d} uses)")
    
    # Statistics
    print("\n📈 Usage Statistics:")
    print(f"  - Total packages: {len(usage_data['packages'])}")
    print(f"  - Total types: {len(usage_data['types'])}")
    print(f"  - Total methods: {len(usage_data['methods'])}")
    print(f"  - Total files analyzed: {len(set(f for files in usage_data['imports'].values() for f in files))}")
    
    # Write outputs
    if args.output_file:
        with open(args.output_file, 'a', encoding='utf-8') as f:
            f.write(f"usage-file={args.usage_output.absolute()}\n")
            f.write(f"packages-count={len(usage_data['packages'])}\n")
            f.write(f"types-count={len(usage_data['types'])}\n")
    
    print("\n✅ Usage mapping complete")


if __name__ == "__main__":
    main()
