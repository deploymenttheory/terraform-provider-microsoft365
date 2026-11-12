#!/usr/bin/env python3
"""Plan isolated lint, unit and race jobs without loading the Go provider."""

import argparse
import json
import os
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
from lib.git_operations import get_changed_packages, get_merge_base
from lib.code_analysis import detect_goroutines, detect_service_domains
from lib.github import write_output
from lib.common import load_pr_checks_config
from lib.partitions import affected_test_packages, partition_packages


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--base-ref', required=True)
    args = parser.parse_args()
    base = get_merge_base(args.base_ref)
    changed = get_changed_packages(base)
    tests = affected_test_packages(changed) if changed else []
    races = detect_goroutines(tests)
    unit_matrix = partition_packages(tests)
    outputs = {
        'base-sha': base,
        'has-changes': str(bool(changed)).lower(),
        'has-goroutines': str(bool(races)).lower(),
        'service-domains': ' '.join(detect_service_domains(tests, load_pr_checks_config())),
        'unit-matrix': json.dumps(unit_matrix, separators=(',', ':')),
        'unit-shards': str(len(unit_matrix['include'])),
        'lint-matrix': json.dumps(partition_packages(changed), separators=(',', ':')),
        'race-matrix': json.dumps(partition_packages(races), separators=(',', ':')),
    }
    print(json.dumps(outputs, indent=2))
    write_output(outputs, os.environ.get('GITHUB_OUTPUT'))
    return 0


if __name__ == '__main__':
    sys.exit(main())
