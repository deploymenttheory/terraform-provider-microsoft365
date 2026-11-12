"""Exercise termination reporting without compiling the provider or calling APIs."""

import contextlib
import importlib.util
import io
import os
import signal
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from scripts.pipeline.test_support import TemporaryDirectoryTestCase

SCRIPT = Path(__file__).resolve().parents[1] / "run_tests.py"
SPEC = importlib.util.spec_from_file_location("bi_weekly_runner", SCRIPT)
runner = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(runner)


class CommandTests(unittest.TestCase):
    """Exercise command behavior."""

    def test_real_command_output_and_failure_are_preserved(self):
        """Verify real command output and failure are preserved."""
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / "output.log"
            status = runner.run_command(
                [sys.executable, "-c", 'print("test output"); raise SystemExit(7)'],
                str(output),
            )
            self.assertEqual(status, 7)
            self.assertIn("test output\n", output.read_text(encoding="utf-8"))
            self.assertIn("[COMMAND] Exit=7", output.read_text(encoding="utf-8"))

    @unittest.skipIf(os.name == "nt", "POSIX signals required")
    def test_real_sigterm_is_reported_as_143(self):
        """Verify real sigterm is reported as 143."""
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory) / "output.log"
            status = runner.run_command(
                [sys.executable, "-c", "import os, signal; os.kill(os.getpid(), signal.SIGTERM)"],
                str(output),
            )
            self.assertEqual(status, 128 + signal.SIGTERM)
            self.assertIn("Terminated by SIGTERM; exit=143", output.read_text(encoding="utf-8"))


class PackageTests(TemporaryDirectoryTestCase):
    """Exercise package behavior."""

    def setUp(self):
        super().setUp()
        self.change_directory()
        Path("internal/services/resources/example").mkdir(parents=True)

    def run_packages(self, core, statuses):
        """Run packages."""
        packages = ["./first", "./must-not-run"]
        Path("output.log").write_text(
            "=== RUN   TestPassed\n--- PASS: TestPassed (0.01s)\n", encoding="utf-8"
        )
        with patch.object(
            runner,
            "discover_test_packages",
            side_effect=[packages, [], [], []] if core else [packages],
        ):
            with (
                patch.object(runner, "run_command", side_effect=statuses) as command,
                contextlib.redirect_stdout(io.StringIO()),
            ):
                if core:
                    status = runner.run_provider_core_tests(
                        "coverage.txt", "output.log", skip_enumeration=True
                    )
                else:
                    status = runner.run_service_tests(
                        "resources", "example", "coverage.txt", "output.log", skip_enumeration=True
                    )
        return status, command

    def test_termination_stops_packages_and_preserves_partial_results(self):
        """Verify termination stops packages and preserves partial results."""
        for core in (False, True):
            with self.subTest(core=core):
                status, command = self.run_packages(core, [143])
                self.assertEqual(status, 143)
                self.assertEqual(command.call_count, 1)
                self.assertIn("TestPassed", Path("test-successes.json").read_text(encoding="utf-8"))
                self.assertEqual(Path("test-failures.json").read_text(encoding="utf-8"), "[]")

    def test_ordinary_failure_continues_packages_but_stays_failed(self):
        """Verify ordinary failure continues packages but stays failed."""
        for core in (False, True):
            with self.subTest(core=core):
                status, command = self.run_packages(core, [1, 0])
                self.assertEqual(status, 1)
                self.assertEqual(command.call_count, 2)

    def test_shard_applies_filters_to_their_own_packages_only(self):
        """Verify shard applies filters to their own packages only."""
        expression = "^(TestAccFirst|TestUnitSecond)$"
        Path("output.log").write_text(
            "=== RUN   TestAccFirst\n--- PASS: TestAccFirst (0.01s)\n", encoding="utf-8"
        )
        packages = [
            "./internal/services/resources/example/graph_beta/first",
            "./internal/services/resources/example/graph_beta/second",
            "./internal/services/resources/example/graph_beta/unfiltered",
        ]
        with patch.object(runner, "discover_test_packages", return_value=packages):
            with (
                patch.object(runner, "run_command", return_value=0) as command,
                contextlib.redirect_stdout(io.StringIO()),
            ):
                status = runner.run_service_tests(
                    "resources",
                    "example",
                    "coverage.txt",
                    "output.log",
                    skip_enumeration=True,
                    test_filters={
                        "graph_beta/first": expression,
                        "graph_beta/second": "^TestOther$",
                    },
                )
        self.assertEqual(status, 0)
        self.assertEqual(command.call_count, 3)
        for call, expected in zip(command.call_args_list, [expression, "^TestOther$", None]):
            arguments = call.args[0]
            if expected:
                self.assertEqual(arguments[arguments.index("-run") + 1], expected)
            else:
                self.assertNotIn("-run", arguments)
            self.assertIn("-timeout=90m", arguments)

    def test_positional_run_options_reach_each_go_command(self):
        """Verify existing positional callers retain their concurrency and race options."""
        for core in (False, True):
            with self.subTest(core=core):
                package = "./internal/services/resources/example/first"
                Path("output.log").write_text("", encoding="utf-8")
                discovery = [[package], [], [], []] if core else [[package]]
                with (
                    patch.object(runner, "discover_test_packages", side_effect=discovery),
                    patch.object(runner, "run_command", return_value=0) as command,
                    contextlib.redirect_stdout(io.StringIO()),
                ):
                    options = (2, 4, 3, True, True, False)
                    if core:
                        runner.run_provider_core_tests("coverage.txt", "output.log", *options)
                    else:
                        runner.run_service_tests(
                            "resources", "example", "coverage.txt", "output.log", *options
                        )
                arguments = command.call_args.args[0]
                self.assertIn("-p=3", arguments)
                self.assertIn("-parallel=4", arguments)
                self.assertIn("-race", arguments)

    def test_filter_for_missing_package_fails_before_running_tests(self):
        """Verify filter for missing package fails before running tests."""
        with patch.object(runner, "discover_test_packages", return_value=["./first"]):
            with (
                patch.object(runner, "run_command") as command,
                contextlib.redirect_stdout(io.StringIO()),
            ):
                with self.assertRaisesRegex(ValueError, "outside this shard"):
                    runner.run_service_tests(
                        "resources",
                        "example",
                        "coverage.txt",
                        "output.log",
                        test_filters={"graph_beta/missing": "^TestMissing$"},
                    )
        command.assert_not_called()

    def test_missing_unfiltered_shard_resource_also_fails(self):
        """Verify missing unfiltered shard resource also fails."""
        with patch.object(runner, "discover_test_packages", return_value=[]):
            with (
                patch.object(runner, "run_command") as command,
                contextlib.redirect_stdout(io.StringIO()),
            ):
                with self.assertRaisesRegex(ValueError, "Shard resource not found"):
                    runner.run_service_tests(
                        "resources",
                        "example",
                        "coverage.txt",
                        "output.log",
                        shard_resources="graph_beta/missing",
                    )
        command.assert_not_called()
