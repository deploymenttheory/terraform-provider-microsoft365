"""Exercise termination reporting without compiling the provider or calling APIs."""

import importlib.util
import contextlib
import io
import os
import signal
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

SCRIPT = Path(__file__).resolve().parents[1] / 'run_tests.py'
SPEC = importlib.util.spec_from_file_location('bi_weekly_runner', SCRIPT)
runner = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(runner)


class CommandTests(unittest.TestCase):
    def test_real_command_output_and_failure_are_preserved(self):
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / 'output.log'
            status = runner.run_command(
                [sys.executable, '-c', 'print("test output"); raise SystemExit(7)'],
                str(output),
            )
            self.assertEqual(status, 7)
            self.assertIn('test output\n', output.read_text())
            self.assertIn('[COMMAND] Exit=7', output.read_text())

    @unittest.skipIf(os.name == 'nt', 'POSIX signals required')
    def test_real_sigterm_is_reported_as_143(self):
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / 'output.log'
            status = runner.run_command(
                [sys.executable, '-c', 'import os, signal; os.kill(os.getpid(), signal.SIGTERM)'],
                str(output),
            )
            self.assertEqual(status, 128 + signal.SIGTERM)
            self.assertIn('Terminated by SIGTERM; exit=143', output.read_text())


class PackageTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        previous = Path.cwd()
        os.chdir(self.directory.name)
        self.addCleanup(os.chdir, previous)
        Path('internal/services/resources/example').mkdir(parents=True)

    def run_packages(self, core, statuses):
        packages = ['./first', './must-not-run']
        Path('output.log').write_text('=== RUN   TestPassed\n--- PASS: TestPassed (0.01s)\n')
        with patch.object(runner, 'discover_test_packages', side_effect=[packages, [], [], []] if core else [packages]):
            with patch.object(runner, 'run_command', side_effect=statuses) as command, contextlib.redirect_stdout(io.StringIO()):
                if core:
                    status = runner.run_provider_core_tests('coverage.txt', 'output.log', skip_enumeration=True)
                else:
                    status = runner.run_service_tests('resources', 'example', 'coverage.txt', 'output.log', skip_enumeration=True)
        return status, command

    def test_termination_stops_packages_and_preserves_partial_results(self):
        for core in (False, True):
            with self.subTest(core=core):
                status, command = self.run_packages(core, [143])
                self.assertEqual(status, 143)
                self.assertEqual(command.call_count, 1)
                self.assertIn('TestPassed', Path('test-successes.json').read_text())
                self.assertEqual(Path('test-failures.json').read_text(), '[]')

    def test_ordinary_failure_continues_packages_but_stays_failed(self):
        for core in (False, True):
            with self.subTest(core=core):
                status, command = self.run_packages(core, [1, 0])
                self.assertEqual(status, 1)
                self.assertEqual(command.call_count, 2)
