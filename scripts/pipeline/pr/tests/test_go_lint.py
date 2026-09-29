"""Verify lint scope and changed-line enforcement without loading the provider."""

import json
import os
import subprocess
from pathlib import Path
from unittest.mock import patch

from support import RepositoryTest
from run_lint import run_lint


class LintRunnerTests(RepositoryTest):
    def test_no_changes_launches_no_linter(self):
        # The selector uses Git; only the linter subprocess is forbidden here.
        with patch('run_lint.get_merge_base', return_value=self.base):
            with patch('run_lint.get_changed_packages', return_value=[]):
                with patch('run_lint.subprocess.run') as run:
                    self.assertEqual(run_lint(self.base), 0)
        run.assert_not_called()
        self.assertFalse(Path('lint-results').exists())

    def test_explicit_packages_and_all_exit_statuses(self):
        with patch('run_lint.get_merge_base', return_value=self.base):
            with patch('run_lint.get_changed_packages', return_value=['.', 'selected', 'unrelated']):
                with patch('run_lint.subprocess.run', side_effect=[
                    subprocess.CompletedProcess([], 0),
                    subprocess.CompletedProcess([], 1),
                    subprocess.CompletedProcess([], 3),
                ]) as run:
                    self.assertEqual(run_lint(self.base), 1)
        commands = [call.args[0] for call in run.call_args_list]
        self.assertEqual([cmd[-1] for cmd in commands], ['./.', './selected', './unrelated'])
        for cmd in commands:
            self.assertIn('--concurrency=1', cmd)
            self.assertIn('--fix=false', cmd)
            self.assertIn('--whole-files=false', cmd)
            self.assertIn('--issues-exit-code=1', cmd)
            self.assertIn(f'--new-from-rev={self.base}', cmd)
            self.assertNotIn('./...', cmd)
        report = json.loads(Path('lint-results/results.json').read_text())
        self.assertEqual([r['status'] for r in report['results']],
                         ['passed', 'findings', 'execution failure'])

    def test_real_lint_filters_old_findings_fails_new_findings_and_does_not_fix(self):
        self.write('.golangci.yml', '''version: "2"
linters:
  default: none
  enable: [misspell]
formatters:
  enable: [gofumpt]
''')
        self.write('selected/value.go', '''package selected

// Value can recieve numbers.
func Value() int {
	return 1
}
''')
        # This package must never be passed to lint or compiled by it.
        self.write('unrelated/value.go', 'intentionally invalid Go\n')
        self.commit()
        base = self.git('rev-parse', 'HEAD')
        source = Path('selected/value.go').read_text().replace('return 1', 'return 2')
        self.write('selected/value.go', source)
        self.commit()
        with patch.dict(os.environ, {'GOTOOLCHAIN': 'local', 'GOPROXY': 'off'}):
            self.assertEqual(run_lint(base), 0)
            self.write('selected/value.go', source.replace('return 2', '// recieve another number\n\treturn 2'))
            self.commit()
            before = Path('selected/value.go').read_bytes()
            self.assertEqual(run_lint(base), 1)
            self.assertEqual(Path('selected/value.go').read_bytes(), before)
        self.assertEqual(self.git('diff', '--name-only'), '')
        self.assertTrue(Path('lint-results/package-1.sarif').is_file())
