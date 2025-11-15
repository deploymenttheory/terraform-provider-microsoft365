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
                    merged.extend(failures)
                    print(f"  Merged {failure_file.relative_to(failure_file.parents[2])}: {len(failures)} failures")
        except (json.JSONDecodeError, IOError) as e:
            print(f"  Warning: Could not read {failure_file}: {e}", file=sys.stderr)
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
        print("No test-failures.json files found")
        with open(output_file, 'w') as f:
            json.dump([], f)
        print(f"Created empty {output_file}")
        return
    
    print(f"Found {len(failure_files)} test failure file(s)")
    merged = merge_failures(failure_files)
    
    with open(output_file, 'w') as f:
        json.dump(merged, f, indent=2)
    
    print(f"\nTotal failures: {len(merged)}")
    print(f"Merged results written to {output_file}")


if __name__ == "__main__":
    main()

