"""Temporary Git repositories for exercising CI selection without provider builds."""

import contextlib
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

from scripts.pipeline.test_support import TemporaryDirectoryTestCase

PIPELINE = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(PIPELINE))
sys.path.insert(0, str(PIPELINE / "steps"))


class RepositoryTest(TemporaryDirectoryTestCase):
    """Exercise repository behavior."""

    def setUp(self):
        super().setUp()
        self.change_directory()
        self.git("init", "-q")
        self.git("config", "user.name", "CI Tests")
        self.git("config", "user.email", "ci-tests@example.invalid")
        self.write(".gitignore", "lint-results/\ncoverage/\n")
        self.write("go.mod", "module example.com/ci-fixture\n\ngo 1.26.0\n")
        self.write("selected/value.go", "package selected\n\nfunc Value() int { return 1 }\n")
        self.write("unrelated/value.go", "package unrelated\n\nfunc Value() int { return 1 }\n")
        self.commit()
        self.base = self.git("rev-parse", "HEAD")

    def git(self, *args):
        """Run Git in the temporary fixture repository."""
        return subprocess.run(
            ["git", *args], check=True, capture_output=True, text=True
        ).stdout.strip()

    def write(self, name, contents):
        """Write a UTF-8 file into the temporary fixture repository."""
        path = Path(name)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(contents, encoding="utf-8")

    def commit(self):
        """Commit the current fixture changes."""
        self.git("add", ".")
        self.git("commit", "-qm", "Fixture changes")
