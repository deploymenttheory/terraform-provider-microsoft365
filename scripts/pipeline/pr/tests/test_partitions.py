"""Verify package partition completeness, dependent tests, and coverage aggregation."""

import tempfile
import unittest
from pathlib import Path

from support import RepositoryTest
from lib.partitions import affected_test_packages, partition_packages
from merge_coverage import merge_profiles, merge_shard_profiles
from run_lint import run_lint
from unittest.mock import patch


class PartitionTests(RepositoryTest):
    def test_every_package_is_assigned_once_and_shards_are_bounded(self):
        packages = [f'internal/example/package_{index}' for index in range(243)]
        matrix = partition_packages(packages)
        assigned = [package for shard in matrix['include'] for package in shard['packages'].split()]
        self.assertEqual(sorted(assigned), sorted(packages))
        self.assertEqual(len(assigned), len(set(assigned)))
        self.assertTrue(all(len(shard['packages'].split()) <= 5 for shard in matrix['include']))
        self.assertEqual(matrix, partition_packages(list(reversed(packages))))

    def test_empty_selection_launches_no_jobs(self):
        self.assertEqual(partition_packages([]), {'include': []})

    def test_shared_dependency_change_includes_transitive_tests(self):
        self.write('internal/core/config.go', 'package core\n')
        self.write('internal/mocks/mocks.go', '''package mocks
import core "github.com/deploymenttheory/terraform-provider-microsoft365/internal/core"
''')
        self.write('internal/services/first/resource_test.go', '''package first_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
''')
        self.write('internal/services/second/resource_test.go', 'package second_test\n')
        self.assertEqual(affected_test_packages(['internal/core']),
                         ['internal/core', 'internal/services/first'])

    def test_lint_shard_cannot_expand_its_changed_package_scope(self):
        with patch('run_lint.get_merge_base', return_value=self.base), \
             patch('run_lint.get_changed_packages', return_value=['selected']), \
             patch('run_lint.subprocess.run') as run:
            with self.assertRaisesRegex(ValueError, 'only packages changed'):
                run_lint(self.base, selected_packages=['unrelated'])
        run.assert_not_called()


class CoverageMergeTests(unittest.TestCase):
    def test_overlapping_blocks_are_covered_once_after_shards_are_combined(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            first, second, output = root / 'first', root / 'second', root / 'merged'
            first.write_text('mode: atomic\nmodule/a.go:1.1,2.1 2 0\n')
            second.write_text('mode: atomic\nmodule/a.go:1.1,2.1 2 3\nmodule/b.go:1.1,2.1 1 1\n')
            merge_profiles([first, second], output)
            self.assertEqual(output.read_text(), 'mode: atomic\nmodule/a.go:1.1,2.1 2 3\nmodule/b.go:1.1,2.1 1 1\n')

    def test_missing_artifacts_cannot_produce_successful_empty_coverage(self):
        with self.assertRaisesRegex(ValueError, 'No shard coverage'):
            merge_profiles([], Path('must-not-exist'))

    def test_one_successful_shard_cannot_hide_a_missing_shard_artifact(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            profile = root / 'unit-test-results-0' / 'unit-coverage.txt'
            profile.parent.mkdir()
            profile.write_text('mode: atomic\nmodule/a.go:1.1,2.1 1 1\n')
            output = root / 'merged'
            with self.assertRaisesRegex(ValueError, 'Missing shard coverage profiles'):
                merge_shard_profiles(root, output, 2)
            self.assertFalse(output.exists())


class WorkflowGateTests(unittest.TestCase):
    def test_required_checks_keep_their_names_and_require_all_worker_shards(self):
        import yaml
        workflows = Path(__file__).resolve().parents[4] / '.github/workflows'
        unit = yaml.safe_load((workflows / 'pr-tests.yml').read_text(encoding='utf-8'))['jobs']
        lint = yaml.safe_load((workflows / 'go-lint.yml').read_text(encoding='utf-8'))['jobs']
        for jobs, gate, worker, name in (
            (unit, 'unit-tests', 'unit-test-shards', '🧪 Go Unit Tests & Coverage'),
            (unit, 'race-detection', 'race-shards', '🔍 Race Detection'),
            (lint, 'golint', 'lint-shards', '✨ Run golangci-lint'),
        ):
            with self.subTest(gate=gate):
                self.assertEqual(jobs[gate]['name'], name)
                self.assertIn(worker, jobs[gate]['needs'])
                self.assertEqual(jobs[gate]['if'], 'always()')
                self.assertFalse(jobs[worker]['strategy']['fail-fast'])
                self.assertEqual(jobs[worker]['runs-on'], 'ubuntu-24.04-arm')
                self.assertIn(f'needs.{worker}.result', str(jobs[gate]['steps']))
