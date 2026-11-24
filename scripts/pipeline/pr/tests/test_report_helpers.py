"""Verify the report-only helper scan selects service packages and never fails the job."""

import shutil
import subprocess
import unittest
from pathlib import Path
from unittest.mock import Mock, patch

from support import RepositoryTest

# isort: split

from report_helpers import SKIPPED, report_helpers, service_packages

REPOSITORY = Path(__file__).resolve().parents[4]
SERVICE = "internal/services/resources/example/graph_beta/widget"


class ServicePackageTests(unittest.TestCase):
    """Exercise package selection."""

    def test_only_service_packages_are_selected(self):
        """Verify only non-common packages under internal/services are kept."""
        packages = [
            ".",
            "internal/provider",
            SERVICE,
            "tools/helperscan",
            "internal/servicesx",
            "internal/services/common/crud",
        ]
        self.assertEqual(service_packages(packages), [SERVICE])


class ReportHelpersTests(RepositoryTest):
    """Exercise the report step against a temporary repository."""

    def summary(self):
        """Return the step summary written so far."""
        path = Path("summary.md")
        return path.read_text(encoding="utf-8") if path.exists() else ""

    def test_no_service_changes_skips_the_scan(self):
        """Verify a change outside internal/services runs nothing and says so."""
        self.write("selected/value.go", "package selected\n\nfunc Value() int { return 2 }\n")
        self.commit()
        run = Mock()
        self.assertEqual(report_helpers(self.base, "summary.md", run), 0)
        run.assert_not_called()
        self.assertIn(SKIPPED, self.summary())

    def test_service_changes_print_and_summarize_the_report(self):
        """Verify the scanner is built once and run on the changed service package."""
        self.write(f"{SERVICE}/crud.go", "package widget\n")
        self.commit()
        outputs = [
            subprocess.CompletedProcess([], 0),
            subprocess.CompletedProcess([], 0, stdout="TEXT REPORT\n", stderr="skipping x\n"),
            subprocess.CompletedProcess([], 0, stdout="## Helper report\n"),
        ]
        run = Mock(side_effect=outputs)
        self.assertEqual(report_helpers(self.base, "summary.md", run), 0)
        commands = [call.args[0] for call in run.call_args_list]
        self.assertEqual(commands[0][:2], ["go", "build"])
        self.assertEqual(commands[1][1:], ["-format", "text", SERVICE])
        self.assertEqual(commands[2][1:], ["-format", "markdown", SERVICE])
        self.assertEqual(self.summary(), "## Helper report\n")

    def test_scanner_failure_is_only_a_warning(self):
        """Verify a failed build leaves the job passing."""
        self.write(f"{SERVICE}/crud.go", "package widget\n")
        self.commit()
        failure = subprocess.CalledProcessError(1, ["go", "build"])
        with patch("builtins.print") as printed:
            self.assertEqual(report_helpers(self.base, None, Mock(side_effect=failure)), 0)
        self.assertIn("::warning", printed.call_args_list[-1].args[0])

    def test_unknown_base_ref_is_only_a_warning(self):
        """Verify a Git failure leaves the job passing."""
        with patch("builtins.print") as printed:
            self.assertEqual(report_helpers("does-not-exist", None), 0)
        self.assertIn("::warning", printed.call_args_list[-1].args[0])

    @unittest.skipUnless(shutil.which("go"), "Go is not installed")
    def test_real_scanner_reports_the_changed_package(self):
        """Verify the real scanner runs end to end on a changed service package."""
        for source in (REPOSITORY / "tools" / "helperscan").glob("*.go"):
            if not source.name.endswith("_test.go"):
                self.write(f"tools/helperscan/{source.name}", source.read_text(encoding="utf-8"))
        self.commit()
        base = self.git("rev-parse", "HEAD")
        self.write(
            f"{SERVICE}/crud.go",
            "package widget\n\ntype R struct{}\n\nfunc (r *R) Create() { matches() }\n\n"
            "func (r *R) Read() { matches() }\n\nfunc matches() {}\n",
        )
        self.commit()
        with patch("builtins.print") as printed:
            self.assertEqual(report_helpers(base, "summary.md"), 0)
        log = "".join(str(call.args[0]) for call in printed.call_args_list)
        self.assertIn("Packages scanned: 1", log)
        self.assertIn("banned generic name", log)
        self.assertIn("| widget | crud.go | `matches` |", self.summary())


if __name__ == "__main__":
    unittest.main()
