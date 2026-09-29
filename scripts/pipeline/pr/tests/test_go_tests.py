"""Regression tests for the PR unit-test runner."""

import importlib.util
import os
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch


spec = importlib.util.spec_from_file_location(
    "go_tests", Path(__file__).resolve().parents[1] / "lib" / "go_tests.py"
)
go_tests = importlib.util.module_from_spec(spec)
spec.loader.exec_module(go_tests)


class UnitTestRunnerTests(unittest.TestCase):
    def test_success_merges_coverage_with_serial_unit_test_environment(self):
        def successful_test(command, *, env, check):
            self.assertTrue(check)
            self.assertEqual(env["TF_ACC"], "0")
            self.assertIn("-p=1", command)
            self.assertIn("-parallel=1", command)
            coverage = next(arg.split("=", 1)[1] for arg in command
                            if arg.startswith("-coverprofile="))
            Path(coverage).write_text("mode: atomic\nexample.go:1.1,2.1 1 1\n")
            return subprocess.CompletedProcess(command, 0)

        with tempfile.TemporaryDirectory() as directory:
            with patch.dict(os.environ, {"TF_ACC": "1"}):
                with patch.object(go_tests.subprocess, "run", side_effect=successful_test) as run:
                    coverage = go_tests.run_unit_tests(["first", "second"], directory)
            self.assertEqual(run.call_count, 2)
            self.assertEqual(coverage.read_text().count("mode: atomic"), 1)
            self.assertEqual(coverage.read_text().count("example.go:"), 2)

    def test_failed_test_cannot_be_reported_as_successful_coverage(self):
        self.assert_failed_process_propagates(1)

    def test_terminated_test_cannot_be_reported_as_successful_coverage(self):
        self.assert_failed_process_propagates(-15)

    def assert_failed_process_propagates(self, returncode):
        def failed_test(command, *, env, check):
            self.assertTrue(check)
            coverage = next(arg.split("=", 1)[1] for arg in command
                            if arg.startswith("-coverprofile="))
            Path(coverage).write_text("mode: atomic\n")
            raise subprocess.CalledProcessError(returncode, command)

        with tempfile.TemporaryDirectory() as directory:
            with patch.object(go_tests.subprocess, "run", side_effect=failed_test) as run:
                with self.assertRaises(subprocess.CalledProcessError) as error:
                    go_tests.run_unit_tests(["first", "must-not-run"], directory)
            self.assertEqual(error.exception.returncode, returncode)
            self.assertEqual(run.call_count, 1)
            self.assertFalse((Path(directory) / "unit-coverage.txt").exists())


if __name__ == "__main__":
    unittest.main()
