#!/usr/bin/env python3
"""Lint exact changed Go packages, reporting findings on changed lines only."""

import argparse
import json
import os
import subprocess
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent))

from lib.git_operations import get_changed_packages, get_merge_base  # noqa: E402


def run_lint(base_ref: str, output_dir: str = 'lint-results') -> int:
    """Run each selected package once, preserving exit codes and streamed logs."""
    base_sha = get_merge_base(base_ref)
    packages = get_changed_packages(base_sha)
    if not packages:
        print('No changed Go packages to lint.', flush=True)
        return 0

    reports = Path(output_dir)
    reports.mkdir(parents=True, exist_ok=True)
    results = []
    for index, package in enumerate(packages, 1):
        print(f'[{index}/{len(packages)}] Linting {package} against {base_sha}', flush=True)
        cmd = [
            'golangci-lint', 'run', '--config=.golangci.yml', '--verbose',
            '--concurrency=1', '--timeout=45m', '--issues-exit-code=1',
            '--fix=false', '--new=false', '--new-from-merge-base=',
            f'--new-from-rev={base_sha}', '--whole-files=false',
            f'--output.sarif.path={reports / f"package-{index}.sarif"}',
            f'./{package}',
        ]
        result = subprocess.run(
            cmd, check=False,
            env={**os.environ, 'GOMAXPROCS': '1', 'GOFLAGS': '-p=1'},
        )
        status = 'passed' if result.returncode == 0 else (
            'findings' if result.returncode == 1 else 'execution failure'
        )
        results.append({'package': package, 'exit_code': result.returncode, 'status': status})
        # Persist after each package so earlier results survive a later interruption.
        (reports / 'results.json').write_text(
            json.dumps({'base_sha': base_sha, 'results': results}, indent=2) + '\n',
            encoding='utf-8',
        )
        print(f'{package}: {status} (exit {result.returncode})', flush=True)
    return 1 if any(result['exit_code'] != 0 for result in results) else 0


def main() -> int:
    """Select packages from a local Git comparison, never from the patch API."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--base-ref', required=True)
    parser.add_argument('--output-dir', default='lint-results')
    args = parser.parse_args()
    return run_lint(args.base_ref, args.output_dir)


if __name__ == '__main__':
    sys.exit(main())
