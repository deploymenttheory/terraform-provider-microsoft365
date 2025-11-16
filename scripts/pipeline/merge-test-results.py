#!/usr/bin/env python3
"""
Merges multiple test result JSON files from artifacts directory.
Usage: ./merge-test-results.py <artifacts-dir> <output-file> [filename-to-merge]
"""

import sys
import json
from pathlib import Path


def find_result_files(artifacts_dir: Path, filename: str) -> list[Path]:
    """Find all test result files with given filename in artifacts directory."""
    return list(artifacts_dir.rglob(filename))


def merge_results(result_files: list[Path], show_details: bool = False) -> list[dict]:
    """Merge all result JSON files into a single list."""
    merged = []
    
    for result_file in result_files:
        try:
            with open(result_file) as f:
                results = json.load(f)
                if isinstance(results, list):
                    if show_details:
                        for result in results:
                            test_name = result.get("test_name", "Unknown")
                            category = result.get("category", "")
                            service = result.get("service", "")
                            service_path = f"{category}/{service}" if service else category
                            print(f"  ❌ {test_name} ({service_path})")
                    
                    merged.extend(results)
        except (json.JSONDecodeError, IOError) as e:
            print(f"  ⚠️  Could not read {result_file}: {e}", file=sys.stderr)
            continue
    
    return merged


def main():
    if len(sys.argv) < 3:
        print("Usage: merge-test-results.py <artifacts-dir> <output-file> [filename-to-merge]", file=sys.stderr)
        sys.exit(1)
    
    artifacts_dir = Path(sys.argv[1])
    output_file = Path(sys.argv[2])
    filename = sys.argv[3] if len(sys.argv) > 3 else "test-failures.json"
    
    if not artifacts_dir.exists():
        print(f"Error: Artifacts directory not found: {artifacts_dir}", file=sys.stderr)
        sys.exit(1)
    
    is_failures = "failure" in filename
    result_type = "failures" if is_failures else "successes"
    
    print(f"Searching for {filename} in {artifacts_dir}...")
    result_files = find_result_files(artifacts_dir, filename)
    
    if not result_files:
        print(f"✅ No {filename} files found")
        with open(output_file, 'w') as f:
            json.dump([], f)
        return
    
    print(f"\nProcessing {len(result_files)} {filename} file(s):\n")
    merged = merge_results(result_files, show_details=is_failures)
    
    with open(output_file, 'w') as f:
        json.dump(merged, f, indent=2)
    
    if merged:
        print(f"\n{'='*60}")
        print(f"Total test {result_type}: {len(merged)}")
        print(f"{'='*60}")
    else:
        print(f"✅ No test {result_type} found")


if __name__ == "__main__":
    main()

