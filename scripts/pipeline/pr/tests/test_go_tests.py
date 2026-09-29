"""Verify package scope, serial execution, acceptance isolation, and failure propagation."""

import os
import subprocess
from pathlib import Path
from unittest.mock import patch

from support import RepositoryTest
from lib.go_tests import run_unit_tests, run_race_detection


class TestRunnerTests(RepositoryTest):
    def test_empty_selections_do_not_execute_go(self):
        with patch('lib.go_tests.subprocess.run') as run:
            run_unit_tests([], 'coverage')
            self.assertEqual(run_race_detection([]), 0)
        run.assert_not_called()

    def test_test_and_race_commands_are_serial_and_disable_acceptance(self):
        with patch.dict(os.environ, {'TF_ACC': '1', 'GOFLAGS': '-p=8', 'GOMAXPROCS': '8'}):
            with patch('lib.go_tests.subprocess.run', return_value=subprocess.CompletedProcess([], 0)) as run:
                run_unit_tests(['selected'], 'coverage')
                run_race_detection(['selected'])
        for call in run.call_args_list:
            command = call.args[0]
            self.assertEqual(command[-1], './selected')
            self.assertIn('-p=1', command)
            self.assertIn('-parallel=1', command)
            self.assertIn('-skip=^TestAcc', command)
            self.assertNotIn('./...', command)
            self.assertEqual(call.kwargs['env']['TF_ACC'], '0')
            self.assertEqual(call.kwargs['env']['GOMAXPROCS'], '1')
            self.assertEqual(call.kwargs['env']['GOFLAGS'], '-p=1')

    def test_failure_is_not_hidden_by_coverage(self):
        def failed(command, **kwargs):
            profile = next(arg.split('=', 1)[1] for arg in command if arg.startswith('-coverprofile='))
            Path(profile).write_text('mode: atomic\nexample.com/selected/value.go:1.1,1.5 1 1\n')
            raise subprocess.CalledProcessError(1, command)
        with patch('lib.go_tests.subprocess.run', side_effect=failed):
            with self.assertRaises(subprocess.CalledProcessError):
                run_unit_tests(['selected'])
        self.assertIn('example.com/selected', Path('coverage/selected.out').read_text())
        self.assertFalse(Path('coverage/unit-coverage.txt').exists())

    def test_stale_coverage_is_removed_before_failed_build(self):
        self.write('coverage/selected.out', 'mode: atomic\nstale.go:1.1,1.5 1 1\n')
        self.write('coverage/unit-coverage.txt', 'mode: atomic\nstale.go:1.1,1.5 1 1\n')
        with patch('lib.go_tests.subprocess.run', side_effect=subprocess.CalledProcessError(1, ['go', 'test'])):
            with self.assertRaises(subprocess.CalledProcessError):
                run_unit_tests(['selected'])
        self.assertFalse(Path('coverage/selected.out').exists())
        self.assertFalse(Path('coverage/unit-coverage.txt').exists())

    def test_success_merges_coverage_from_multiple_packages(self):
        def successful(command, *, env, check):
            self.assertTrue(check)
            profile = next(arg.split('=', 1)[1] for arg in command if arg.startswith('-coverprofile='))
            Path(profile).write_text('mode: atomic\nexample.go:1.1,2.1 1 1\n')
            return subprocess.CompletedProcess(command, 0)
        with patch('lib.go_tests.subprocess.run', side_effect=successful) as run:
            coverage = run_unit_tests(['first', 'second'])
        self.assertEqual(run.call_count, 2)
        self.assertEqual(coverage.read_text().count('mode: atomic'), 1)
        self.assertEqual(coverage.read_text().count('example.go:'), 2)

    def test_failed_or_terminated_test_stops_execution_and_preserves_exit_code(self):
        for returncode in [1, -15]:
            with self.subTest(returncode=returncode):
                failure = subprocess.CalledProcessError(returncode, ['go', 'test'])
                with patch('lib.go_tests.subprocess.run', side_effect=failure) as run:
                    with self.assertRaises(subprocess.CalledProcessError) as error:
                        run_unit_tests(['first', 'must-not-run'])
                self.assertEqual(error.exception.returncode, returncode)
                self.assertEqual(run.call_count, 1)
                self.assertFalse(Path('coverage/unit-coverage.txt').exists())

    def test_race_failure_propagates(self):
        with patch('lib.go_tests.subprocess.run', return_value=subprocess.CompletedProcess([], 1)):
            self.assertEqual(run_race_detection(['selected']), 1)

    def test_real_go_runs_only_selected_package_and_forces_tf_acc_zero(self):
        self.write('selected/value_test.go', '''package selected
import ("os"; "testing")
func TestUnitSelected(t *testing.T) {
 if os.Getenv("TF_ACC") != "0" { t.Fatal("acceptance environment leaked") }
 if Value() != 1 { t.Fatal("wrong value") }
}
func TestAccSelected(t *testing.T) {
 if os.Getenv("TF_ACC") != "" { t.Fatal("acceptance test must not execute") }
}
''')
        self.write('unrelated/value_test.go', '''package unrelated
import "testing"
func TestUnrelated(t *testing.T) { t.Fatal("unrelated tests must not run") }
''')
        with patch.dict(os.environ, {'TF_ACC': '1', 'GOTOOLCHAIN': 'local', 'GOPROXY': 'off'}):
            coverage = run_unit_tests(['selected'])
        self.assertIn('/selected/', coverage.read_text())
        self.assertNotIn('/unrelated/', coverage.read_text())

    def test_real_go_failure_reaches_caller(self):
        self.write('selected/value_test.go', '''package selected
import "testing"
func TestFailure(t *testing.T) { t.Fatal("intentional failure") }
''')
        with patch.dict(os.environ, {'GOTOOLCHAIN': 'local', 'GOPROXY': 'off'}):
            with self.assertRaises(subprocess.CalledProcessError):
                run_unit_tests(['selected'])
