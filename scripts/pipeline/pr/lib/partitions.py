"""Partition Go jobs using the acceptance pipeline's bounded Rendezvous sharder."""

import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2] / 'bi_weekly_tests'))
from generate_test_matrix import calculate_shard_count, distribute_rendezvous, RENDEZVOUS_SEED

MODULE = 'github.com/deploymenttheory/terraform-provider-microsoft365'


def affected_test_packages(changed, root=Path('.')):
    """Include tests whose local import graph depends on a changed package."""
    imports = {}
    tests = set()
    for directory in sorted({path.parent for path in (root / 'internal').rglob('*.go')}):
        package = directory.relative_to(root).as_posix()
        sources = list(directory.glob('*.go'))
        if any(path.name.endswith('_test.go') for path in sources):
            tests.add(package)
        imports[package] = set()
        for path in sources:
            source = path.read_text(encoding='utf-8')
            for declaration in re.findall(r'\bimport\s*(\([\s\S]*?\)|(?:\w+\s+)?"[^"\n]+")', source):
                imports[package].update(re.findall('"' + re.escape(MODULE) + r'/([^"\n]+)"', declaration))
    affected = set(changed)
    while True:
        expanded = affected | {package for package, deps in imports.items() if deps & affected}
        if expanded == affected:
            break
        affected = expanded
    return sorted(set(changed) | (affected & tests))


def partition_packages(packages):
    """Assign every selected package exactly once to a bounded standard-runner job."""
    packages = sorted(set(packages))
    if not packages:
        return {'include': []}
    weights = dict.fromkeys(packages, 1)
    count = calculate_shard_count(len(packages))
    while count <= len(packages):
        try:
            shards = distribute_rendezvous(packages, count, RENDEZVOUS_SEED, weights)
            break
        except ValueError:
            count += 1
    occupied = [items for items in shards.values() if items]
    if len(occupied) > 256:
        raise ValueError('Go job matrix exceeds the GitHub Actions 256-job limit')
    return {'include': [{'shard_index': index, 'packages': ' '.join(items)}
                        for index, items in enumerate(occupied)]}
