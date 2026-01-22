#!/usr/bin/env python3
"""Fetch coverage from Codecov and enforce threshold.

Fetches patch coverage from Codecov API and validates against minimum threshold.

Outputs:
- coverage-pct: Coverage percentage from Codecov
- total-lines: Total lines changed
- covered-lines: Covered lines

Exit codes:
- 0: Coverage meets threshold
- 1: Coverage below threshold or fetch failed
"""

import argparse
import os
import sys
from pathlib import Path

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from lib.get_codecov_coverage import fetch_codecov_coverage  # noqa: E402
from lib.coverage import check_coverage_threshold  # noqa: E402
from lib.common import load_pr_checks_config  # noqa: E402
from lib.github import write_output  # noqa: E402


def main():
    """Fetch coverage from Codecov and enforce threshold."""
    parser = argparse.ArgumentParser(description='Enforce coverage threshold via Codecov')
    parser.add_argument('--repo-slug', required=True,
                        help='Repository slug (owner/repo)')
    parser.add_argument('--pr-number', required=True,
                        help='Pull request number')
    parser.add_argument('--codecov-token', default=os.environ.get('CODECOV_TOKEN'),
                        help='Codecov API token')
    parser.add_argument('--config', default=None,
                        help='Path to PR checks config file')
    parser.add_argument('--github-output', default=os.environ.get('GITHUB_OUTPUT'),
                        help='Path to GITHUB_OUTPUT file')
    
    args = parser.parse_args()
    
    if not args.codecov_token:
        print("❌ ERROR: CODECOV_TOKEN is required")
        return 1
    
    config = load_pr_checks_config(args.config)
    
    # Fetch coverage from Codecov
    print("\n📊 Fetching coverage from Codecov...")
    stats = fetch_codecov_coverage(args.repo_slug, args.pr_number, args.codecov_token)
    
    if not stats:
        print("❌ ERROR: Could not fetch coverage from Codecov")
        return 1
    
    print(f"\n{'='*60}")
    print(f"Codecov Patch Coverage: {stats['coverage_pct']}%")
    print(f"Total Lines Changed: {stats['total_lines']}")
    print(f"Covered Lines: {stats['covered_lines']}")
    print(f"{'='*60}")
    
    write_output({
        "coverage-pct": str(stats['coverage_pct']),
        "total-lines": str(stats['total_lines']),
        "covered-lines": str(stats['covered_lines'])
    }, args.github_output)
    
    # Enforce threshold
    min_coverage = config.get('coverage_threshold', {}).get('minimum_pct', 60)
    
    if not check_coverage_threshold(stats['coverage_pct'], min_coverage):
        print(f"\n❌ ERROR: Coverage {stats['coverage_pct']}% is below minimum threshold {min_coverage}%")
        print("   Please add tests to increase coverage for changed code.")
        return 1
    
    print(f"\n✅ Coverage {stats['coverage_pct']}% meets minimum threshold {min_coverage}%")
    return 0


if __name__ == "__main__":
    sys.exit(main())
