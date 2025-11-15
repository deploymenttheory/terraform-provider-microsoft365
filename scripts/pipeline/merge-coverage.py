#!/usr/bin/env python3
"""
Merges multiple unit and acceptance test coverage files into one for code cov.
Usage: ./merge-coverage.py <input-dir> <output-file>
"""

import sys
from pathlib import Path
from typing import List


def find_coverage_files(input_dir: Path) -> List[Path]:
    """Find all coverage .txt files in directory recursively."""
    return list(input_dir.rglob("*.txt"))


def merge_coverage_files(coverage_files: List[Path], output_file: Path) -> None:
    """Merge all coverage files into one."""
    total_lines = 0
    
    # Write mode line
    with open(output_file, 'w') as out_f:
        out_f.write("mode: atomic\n")
        
        # Merge each file
        for coverage_file in coverage_files:
            if not coverage_file.is_file():
                continue
            
            # Read file, skip first line (mode: atomic) and empty lines
            with open(coverage_file, 'r') as in_f:
                lines = [line for line in in_f.readlines()[1:] if line.strip()]
            
            if lines:
                out_f.writelines(lines)
                line_count = len(lines)
                total_lines += line_count
                print(f"  Merged {coverage_file.name} ({line_count} lines)")
    
    return total_lines


def main():
    if len(sys.argv) < 3:
        print("Usage: merge-coverage.py <input-dir> <output-file>", file=sys.stderr)
        sys.exit(1)
    
    input_dir = Path(sys.argv[1])
    output_file = Path(sys.argv[2])
    
    # Validate input directory
    if not input_dir.exists():
        print(f"Error: input directory not found: {input_dir}", file=sys.stderr)
        sys.exit(1)
    
    print(f"Searching for coverage files in {input_dir}...")
    
    # Find coverage files
    coverage_files = find_coverage_files(input_dir)
    
    if not coverage_files:
        print("⚠️  No coverage files found, creating empty output")
        with open(output_file, 'w') as f:
            f.write("mode: atomic\n")
        sys.exit(0)
    
    file_count = len(coverage_files)
    print(f"Found {file_count} coverage files to merge")
    
    # Merge files
    total_lines = merge_coverage_files(coverage_files, output_file)
    
    print(f"Total coverage lines: {total_lines}")
    print(f"✅ Successfully merged coverage files to {output_file}")


if __name__ == "__main__":
    main()

