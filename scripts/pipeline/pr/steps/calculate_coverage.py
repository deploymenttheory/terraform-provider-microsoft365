#!/usr/bin/env python3
"""Calculate coverage statistics from coverage file.

Outputs:
- coverage-pct: Coverage percentage
- total-lines: Total statements
- covered-lines: Covered statements
"""

import argparse
import os
import sys
from pathlib import Path

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from lib.coverage import calculate_coverage  # noqa: E402
from lib.github import write_output  # noqa: E402

# pylint: enable=wrong-import-position


def main():
    """Calculate and output coverage statistics."""
    parser = argparse.ArgumentParser(description="Calculate coverage from file")
    parser.add_argument("--coverage-file", required=True, help="Path to coverage file")
    parser.add_argument(
        "--github-output",
        default=os.environ.get("GITHUB_OUTPUT"),
        help="Path to GITHUB_OUTPUT file",
    )

    args = parser.parse_args()

    print("📊 Calculating coverage...")
    coverage_file = Path(args.coverage_file)
    stats = calculate_coverage(coverage_file)

    print(f"\n{'='*60}")
    print(f"Coverage: {stats['coverage_pct']}%")
    print(f"Total Statements: {stats['total_lines']}")
    print(f"Covered Statements: {stats['covered_lines']}")
    print(f"{'='*60}")

    write_output(
        {
            "coverage-pct": str(stats["coverage_pct"]),
            "total-lines": str(stats["total_lines"]),
            "covered-lines": str(stats["covered_lines"]),
        },
        args.github_output,
    )

    return 0


if __name__ == "__main__":
    sys.exit(main())
