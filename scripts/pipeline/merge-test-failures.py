#!/usr/bin/env python3
"""
Merges multiple test-failures.json files from artifacts directory.
Usage: ./merge-test-failures.py <artifacts-dir> <output-file>
"""

import sys
import json
from pathlib import Path


def find_failure_files(artifacts_dir: Path) -> list[Path]:
    """Find all test-failures.json files in artifacts directory."""
    return list(artifacts_dir.rglob("test-failures.json"))


def merge_failures(failure_files: list[Path]) -> list[dict]:
    """Merge all failure JSON files into a single list."""
    merged = []
    
    for failure_file in failure_files:
        try:
            with open(failure_file) as f:
                failures = json.load(f)
                if isinstance(failures, list):
                    # Show which tests failed from this file
                    for failure in failures:
                        test_name = failure.get("test_name", "Unknown")
                        category = failure.get("category", "")
                        service = failure.get("service", "")
                        service_path = f"{category}/{service}" if service else category
                        print(f"  ❌ {test_name} ({service_path})")
                    
                    merged.extend(failures)
        except (json.JSONDecodeError, IOError) as e:
            print(f"  ⚠️  Could not read {failure_file}: {e}", file=sys.stderr)
            continue
    
    return merged


def main():
    if len(sys.argv) < 3:
        print("Usage: merge-test-failures.py <artifacts-dir> <output-file>", file=sys.stderr)
        sys.exit(1)
    
    artifacts_dir = Path(sys.argv[1])
    output_file = Path(sys.argv[2])
    
    if not artifacts_dir.exists():
        print(f"Error: Artifacts directory not found: {artifacts_dir}", file=sys.stderr)
        sys.exit(1)
    
    print(f"Searching for test failure files in {artifacts_dir}...")
    failure_files = find_failure_files(artifacts_dir)
    
    if not failure_files:
        print("✅ No test-failures.json files found")
        with open(output_file, 'w') as f:
            json.dump([], f)
        return
    
    print(f"\nProcessing {len(failure_files)} test failure file(s):\n")
    merged = merge_failures(failure_files)
    
    with open(output_file, 'w') as f:
        json.dump(merged, f, indent=2)
    
    if merged:
        print(f"\n{'='*60}")
        print(f"Total test failures: {len(merged)}")
        print(f"{'='*60}")
    else:
        print("✅ No test failures found")


if __name__ == "__main__":
    main()

