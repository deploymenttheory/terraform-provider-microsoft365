"""Verify real dependency compilation, race artifacts and failure propagation."""

import os
import subprocess
import unittest
from pathlib import Path
from unittest.mock import patch

from support import RepositoryTest

# isort: split

from prepare_go_cache import prepare_go_cache


class CachePreparationTests(RepositoryTest):
    """Compile tiny Go fixtures through the same interface used by cache preparation."""

    def setUp(self):
        super().setUp()
        self.write("dependency/value.go", "package dependency\nfunc Value() int { return 1 }\n")
        self.write("testonly/value.go", "package testonly\nfunc Value() int { return 1 }\n")
        self.write(
            "selected/value.go",
            """package selected
import "example.com/ci-fixture/dependency"
func Value() int { return dependency.Value() }
""",
        )
        self.write(
            "selected/value_test.go",
            """package selected_test
import (
    "os"
    "testing"
    "example.com/ci-fixture/testonly"
)
func TestMustNotRun(t *testing.T) {
    if err := os.WriteFile("test-executed", []byte("executed"), 0600); err != nil { t.Fatal(err) }
    t.Fatal(testonly.Value())
}
""",
        )
        # Unrelated invalid code must never be compiled while preparing this selection.
        self.write("unrelated/value.go", "intentionally invalid Go\n")

    def run_preparation(self, race=False):
        """Record actual commands without replacing the Go compiler."""
        actual_run = subprocess.run
        with patch("prepare_go_cache.subprocess.run", wraps=actual_run) as run:
            prepare_go_cache(["selected"], race=race)
        self.assertFalse(Path("selected/test-executed").exists())
        return [call.args[0] for call in run.call_args_list]

    def test_compiles_production_and_test_only_imports_without_executing_tests(self):
        """Prime the selected dependency graph and leave real test execution downstream."""
        commands = self.run_preparation()
        builds = [command for command in commands if command[:2] == ["go", "build"]]
        self.assertTrue(builds)
        for dependency in ("dependency", "testonly"):
            self.assertTrue(any(f"example.com/ci-fixture/{dependency}" in cmd for cmd in builds))
        self.assertFalse(any("test" in command or "./..." in command for command in commands))
        self.assertTrue(all("-mod=readonly" in command for command in commands))

    @unittest.skipIf(os.name == "nt", "Race compilation is verified on the Linux runner")
    def test_prepares_real_race_instrumented_dependencies(self):
        """Use the actual race compiler without executing test functions."""
        commands = self.run_preparation(race=True)
        self.assertTrue(all("-race" in command for command in commands))

    def test_invalid_test_only_dependency_fails_preparation(self):
        """Do not publish a snapshot when a required test dependency cannot compile."""
        self.write("testonly/value.go", "invalid Go dependency\n")
        with self.assertRaises(subprocess.CalledProcessError):
            self.run_preparation()

    def test_empty_selection_is_rejected(self):
        """A missing planner output cannot produce a successful empty cache."""
        with self.assertRaisesRegex(ValueError, "requires selected packages"):
            prepare_go_cache([])
