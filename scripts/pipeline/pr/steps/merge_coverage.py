#!/usr/bin/env python3
"""Merge successful shard profiles without counting overlapping blocks more than once."""

import argparse
from pathlib import Path


def merge_profiles(paths, output):
    """Validate and combine atomic profiles, retaining execution counts across shards."""
    blocks = {}
    paths = list(paths)
    if not paths:
        raise ValueError('No shard coverage profiles were downloaded')
    for path in paths:
        lines = path.read_text(encoding='utf-8').splitlines()
        if not lines or lines[0] != 'mode: atomic':
            raise ValueError(f'Invalid coverage mode in {path}')
        for line in lines[1:]:
            location, statements, hits = line.rsplit(' ', 2)
            key = (location, int(statements))
            blocks[key] = blocks.get(key, 0) + int(hits)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text('mode: atomic\n' + ''.join(f'{location} {statements} {hits}\n'
                      for (location, statements), hits in sorted(blocks.items())), encoding='utf-8')


def merge_shard_profiles(artifacts, output, expected_shards):
    """Require a profile from every successful shard before producing merged coverage."""
    profiles = [artifacts / f'unit-test-results-{index}' / 'unit-coverage.txt'
                for index in range(expected_shards)]
    missing = [str(path) for path in profiles if not path.is_file()]
    if missing:
        raise ValueError('Missing shard coverage profiles: ' + ', '.join(missing))
    merge_profiles(profiles, output)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--artifacts', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    parser.add_argument('--expected-shards', type=int, required=True)
    args = parser.parse_args()
    merge_shard_profiles(args.artifacts, args.output, args.expected_shards)
