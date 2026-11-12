#!/usr/bin/env python3
"""Plan isolated lint, unit and race jobs without loading the Go provider."""

import argparse
import json
import os
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from lib.code_analysis import detect_goroutines, detect_service_domains
from lib.common import load_pr_checks_config
from lib.git_operations import get_changed_files, get_changed_packages, get_merge_base
from lib.github import write_output
from lib.partitions import (
    affected_test_packages,
    fixture_packages,
    partition_lint_packages,
    partition_packages,
    select_lint_packages,
)

# pylint: enable=wrong-import-position


def main():
    """Select affected packages and emit bounded CI job matrices."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-ref", required=True)
    parser.add_argument(
        "--lint-packages", help="Optional changed-package subset for manual lint diagnostics"
    )
    args = parser.parse_args()
    base = get_merge_base(args.base_ref)
    changed = get_changed_packages(base)
    lint_packages = select_lint_packages(changed, args.lint_packages)
    files = get_changed_files(base, "")
    fixtures = fixture_packages(files)
    production = {
        Path(name).parent.as_posix()
        for name in files
        if name.endswith(".go") and not name.endswith("_test.go")
    }
    tests = affected_test_packages(
        sorted(set(changed) | fixtures), production_changed=production | fixtures
    )
    races = detect_goroutines(tests)
    unit_matrix = partition_packages(tests)
    outputs = {
        "base-sha": base,
        "has-changes": str(bool(tests)).lower(),
        "has-lint-changes": str(bool(lint_packages)).lower(),
        "lint-packages": " ".join(lint_packages),
        "has-provider-changes": str("internal/provider" in lint_packages).lower(),
        "has-goroutines": str(bool(races)).lower(),
        "service-domains": " ".join(detect_service_domains(tests, load_pr_checks_config())),
        "unit-matrix": json.dumps(unit_matrix, separators=(",", ":")),
        "unit-shards": str(len(unit_matrix["include"])),
        "unit-packages": " ".join(tests),
        "race-packages": " ".join(races),
        "lint-matrix": json.dumps(partition_lint_packages(lint_packages), separators=(",", ":")),
        "race-matrix": json.dumps(partition_packages(races), separators=(",", ":")),
    }
    print(json.dumps(outputs, indent=2))
    write_output(outputs, os.environ.get("GITHUB_OUTPUT"))
    return 0


if __name__ == "__main__":
    sys.exit(main())
