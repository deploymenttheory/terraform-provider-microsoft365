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
import subprocess
import sys
from pathlib import Path

import yaml

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from lib.common import load_pr_checks_config  # noqa: E402
from lib.coverage import check_coverage_threshold  # noqa: E402
from lib.get_codecov_coverage import fetch_codecov_coverage  # noqa: E402
from lib.git_operations import get_changed_files  # noqa: E402
from lib.github import write_output  # noqa: E402
from lib.partitions import production_import_closure  # noqa: E402

# pylint: enable=wrong-import-position


def _verify_unreachable_mocks(changed):
    """Require the existing mocks exclusion and prove these packages cannot run in the provider."""
    if any(
        Path(name).parts[0] != "internal" or "mocks" not in Path(name).parts[:-1]
        for name in changed
    ):
        return False
    policy = yaml.safe_load(Path("codecov.yml").read_text(encoding="utf-8"))
    if not isinstance(policy, dict) or not isinstance(policy.get("ignore"), list):
        raise ValueError("Cannot verify the Codecov mocks exclusion")
    if "**/mocks/**" not in policy["ignore"]:
        raise ValueError("Codecov does not exclude the mocks namespace")
    reachable = production_import_closure()
    return all(
        Path(name).is_file() and Path(name).parent.as_posix() not in reachable for name in changed
    )


def _verify_test_only_changes(base_ref):
    """Require Git evidence before treating a completed empty Go patch as inapplicable."""
    if not base_ref:
        print("❌ ERROR: An empty patch requires a Git base reference for verification")
        return False
    try:
        changed = get_changed_files(base_ref)
    except (OSError, subprocess.CalledProcessError) as error:
        print(f"❌ ERROR: Cannot verify changed Go files: {error}")
        return False
    runtime_candidates = [
        name for name in changed if name.endswith(".go") and not name.endswith("_test.go")
    ]
    try:
        if runtime_candidates and not _verify_unreachable_mocks(runtime_candidates):
            print(
                "❌ ERROR: Production Go files changed but Codecov has no executable patch coverage"
            )
            return False
    except (OSError, ValueError, yaml.YAMLError) as error:
        print(f"❌ ERROR: Cannot verify test-only mock packages: {error}")
        return False
    return True


def main():
    """Fetch coverage from Codecov and enforce threshold."""
    parser = argparse.ArgumentParser(description="Enforce coverage threshold via Codecov")
    parser.add_argument("--repo-slug", required=True, help="Repository slug (owner/repo)")
    parser.add_argument("--pr-number", required=True, help="Pull request number")
    parser.add_argument("--commit-sha", help="Require Codecov coverage for this PR head commit")
    parser.add_argument(
        "--base-ref", help="Git base reference for verifying test-only empty patches"
    )
    parser.add_argument(
        "--codecov-token",
        default=os.environ.get("CODECOV_API_TOKEN"),
        help="Optional API access token for private repositories (not an upload token)",
    )
    parser.add_argument("--config", default=None, help="Path to PR checks config file")
    parser.add_argument(
        "--github-output",
        default=os.environ.get("GITHUB_OUTPUT"),
        help="Path to GITHUB_OUTPUT file",
    )

    args = parser.parse_args()

    if not args.codecov_token:
        print("No Codecov API token supplied; requesting public repository coverage")

    config = load_pr_checks_config(args.config)

    # Fetch coverage from Codecov
    print("\n📊 Fetching coverage from Codecov...")
    stats = fetch_codecov_coverage(
        args.repo_slug, args.pr_number, args.codecov_token, expected_commit=args.commit_sha
    )

    if not stats:
        print("❌ ERROR: Could not fetch coverage from Codecov")
        return 1

    no_executable_changes = stats.get("no_executable_changes", False)
    if no_executable_changes and not _verify_test_only_changes(args.base_ref):
        return 1
    coverage_pct = "N/A" if no_executable_changes else stats["coverage_pct"]

    print(f"\n{'='*60}")
    coverage_display = coverage_pct if no_executable_changes else f"{coverage_pct}%"
    print(f"Codecov Patch Coverage: {coverage_display}")
    print(f"Total Lines Changed: {stats['total_lines']}")
    print(f"Covered Lines: {stats['covered_lines']}")
    print(f"{'='*60}")

    write_output(
        {
            "coverage-pct": str(coverage_pct),
            "total-lines": str(stats["total_lines"]),
            "covered-lines": str(stats["covered_lines"]),
        },
        args.github_output,
    )

    if no_executable_changes:
        print("✅ No executable Go lines changed; patch coverage is not applicable")
        return 0

    # Enforce threshold
    min_coverage = config.get("coverage_threshold", {}).get("minimum_pct", 60)

    if not check_coverage_threshold(stats["coverage_pct"], min_coverage):
        print(
            f"\n❌ ERROR: Coverage {stats['coverage_pct']}% is below minimum threshold"
            f" {min_coverage}%"
        )
        print("   Please add tests to increase coverage for changed code.")
        return 1

    print(f"\n✅ Coverage {stats['coverage_pct']}% meets minimum threshold {min_coverage}%")
    return 0


if __name__ == "__main__":
    sys.exit(main())
