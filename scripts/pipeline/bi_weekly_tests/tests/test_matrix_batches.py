"""Keep the large device configuration template suite within CI time limits."""

import importlib.util
import os
import re
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


SCRIPT = Path(__file__).resolve().parents[1] / 'generate_test_matrix.py'
SPEC = importlib.util.spec_from_file_location('test_matrix', SCRIPT)
matrix = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(matrix)


class TemplateBatchTests(unittest.TestCase):
    def test_template_tests_run_once_in_bounded_dedicated_jobs(self):
        package = 'graph_beta/device_configuration_templates_json'
        names = ([f'TestAccTemplate_{index:03d}' for index in range(31)] +
                 [f'TestUnitTemplate_{index:03d}' for index in range(31, 61)])
        with tempfile.TemporaryDirectory() as directory:
            previous = Path.cwd()
            os.chdir(directory)
            try:
                source = Path('internal/services/resources/device_management') / package
                source.mkdir(parents=True)
                (source / 'resource_acceptance_test.go').write_text('\n'.join(
                    f'func {name}(t *testing.T) {{}}' for name in names[:31]
                ))
                (source / 'resource_test.go').write_text('\n'.join(
                    f'func {name}(t *testing.T) {{}}' for name in names[31:]
                ))
                with patch.object(matrix, 'discover_resources', return_value=[package, 'graph_beta/other']):
                    entries = matrix.generate_matrix_for_service('device_management', 'resources')
            finally:
                os.chdir(previous)

        batches = [entry for entry in entries if entry['resources'] == package]
        self.assertEqual(len(batches), 3)
        selected = []
        for batch in batches:
            matches = [name for name in names if re.fullmatch(batch['test_filter'], name)]
            self.assertGreater(len(matches), 0)
            self.assertLessEqual(len(matches), 25)
            self.assertIsNone(re.fullmatch(batch['test_filter'], matches[0] + '_Extra'))
            selected.extend(matches)
        self.assertCountEqual(selected, names)
        self.assertEqual(len(selected), len(set(selected)))
        regular = [entry for entry in entries if entry['resources'] != package]
        self.assertEqual([entry['resources'] for entry in regular], ['graph_beta/other'])
        self.assertEqual(regular[0]['test_filter'], '')
        self.assertEqual([entry['shard_index'] for entry in entries], list(range(4)))
        self.assertTrue(all(entry['total_shards'] == 4 for entry in entries))

    def test_other_services_keep_their_unfiltered_package_jobs(self):
        with patch.object(matrix, 'discover_resources', return_value=['graph_beta/other']):
            entries = matrix.generate_matrix_for_service('users', 'resources')
        self.assertEqual(len(entries), 1)
        self.assertEqual(entries[0]['resources'], 'graph_beta/other')
        self.assertEqual(entries[0]['test_filter'], '')


if __name__ == '__main__':
    unittest.main()
