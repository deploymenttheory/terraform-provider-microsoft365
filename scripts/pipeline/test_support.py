"""Managed temporary resources shared by pipeline regression tests."""

import contextlib
import os
import tempfile
import unittest
from pathlib import Path


class TemporaryDirectoryTestCase(unittest.TestCase):
    """Give each test a temporary directory and deterministic cleanup order."""

    def setUp(self):
        self.resources = contextlib.ExitStack()
        self.addCleanup(self.resources.close)
        self.directory = self.resources.enter_context(tempfile.TemporaryDirectory())

    def change_directory(self):
        """Restore the original directory before temporary files are cleaned up."""
        previous = Path.cwd()
        os.chdir(self.directory)
        self.addCleanup(os.chdir, previous)
