"""Exercise coverage applicability with real Git diffs and an unavailable service."""

import io
import sys
from pathlib import Path
from unittest.mock import patch

from support import RepositoryTest

# isort: split

from enforce_coverage import main


class LocalCoverageTests(RepositoryTest):
    """Only verified non-runtime changes may pass without a Codecov response."""

    def run_gate(self, *, commit=None, base=None, response=None, report=None):
        """Invoke the actual gate against the temporary repository's history."""
        argv = [
            "enforce_coverage.py",
            "--repo-slug",
            "example/provider",
            "--pr-number",
            "1",
            "--commit-sha",
            commit or self.git("rev-parse", "HEAD"),
            "--base-ref",
            base or self.base,
        ]
        with (
            patch.object(sys, "argv", argv),
            patch(
                "enforce_coverage.load_pr_checks_config",
                return_value={"coverage_threshold": {"minimum_pct": 60}},
            ),
            patch("enforce_coverage.fetch_codecov_coverage", return_value=response) as fetch,
            patch("enforce_coverage.fetch_commit_report", return_value=report),
            patch("enforce_coverage.write_output") as output,
            patch("sys.stdout", new_callable=io.StringIO),
        ):
            result = main()
        return result, fetch, output

    def test_acceptance_only_diff_does_not_require_remote_coverage(self):
        """A Codecov outage must not fail a patch containing only acceptance tests."""
        self.write(
            "selected/resource_acceptance_test.go", "package selected\n// acceptance fixture\n"
        )
        self.commit()
        result, fetch, output = self.run_gate()
        self.assertEqual(result, 0)
        fetch.assert_not_called()
        self.assertEqual(
            output.call_args.args[0],
            {"coverage-pct": "N/A", "total-lines": "0", "covered-lines": "0"},
        )

    def test_pipeline_changes_without_go_lines_remain_inapplicable(self):
        """Non-Go tooling changes do not manufacture a patch percentage."""
        self.write("scripts/check.py", "print('check')\n")
        self.commit()
        result, fetch, _ = self.run_gate()
        self.assertEqual(result, 0)
        fetch.assert_not_called()

    def test_earlier_runtime_change_is_not_hidden_by_last_test_commit(self):
        """Verify the entire merge-base diff rather than only the final commit."""
        self.write("selected/value.go", "package selected\nfunc Value() int { return 2 }\n")
        self.commit()
        self.write("selected/value_test.go", "package selected\n// final test change\n")
        self.commit()
        result, fetch, _ = self.run_gate()
        self.assertEqual(result, 1)
        fetch.assert_called_once()
        self.assertEqual(fetch.call_args.kwargs["expected_commit"], self.git("rev-parse", "HEAD"))

    def test_runtime_deletion_still_requires_coverage_verification(self):
        """No added lines does not exempt a deleted production Go file."""
        Path("selected/value.go").unlink()
        self.commit()
        result, fetch, _ = self.run_gate()
        self.assertEqual(result, 1)
        fetch.assert_called_once()

    def test_wrong_head_or_missing_history_cannot_use_local_exemption(self):
        """The local proof must match the exact checked-out source and real base."""
        self.write("selected/value_test.go", "package selected\n")
        self.commit()
        for arguments in ({"commit": "0" * 40}, {"base": "missing-base"}):
            with self.subTest(arguments=arguments):
                result, fetch, _ = self.run_gate(**arguments)
                self.assertEqual(result, 1)
                fetch.assert_called_once()

    def test_runtime_patch_keeps_the_existing_threshold(self):
        """Both passing and failing runtime coverage still come from Codecov."""
        self.write("selected/value.go", "package selected\nfunc Value() int { return 2 }\n")
        self.commit()
        for percentage, expected in ((59, 1), (60, 0)):
            with self.subTest(percentage=percentage):
                result, fetch, _ = self.run_gate(
                    response={
                        "coverage_pct": percentage,
                        "total_lines": 100,
                        "covered_lines": percentage,
                    }
                )
                self.assertEqual(result, expected)
                fetch.assert_called_once()

    def test_complete_commit_report_recovers_comparison_outage(self):
        """The exact uploaded report may recover missing comparison data at the same threshold."""
        self.write("selected/value.go", "package selected\nfunc Value() int { return 2 }\n")
        self.commit()
        head = self.git("rev-parse", "HEAD")
        for covered, expected in ((True, 0), (False, 1)):
            report = {
                "commitid": head,
                "state": "complete",
                "report": {
                    "files": [
                        {
                            "name": "selected/value.go",
                            "line_coverage": [[2, 0 if covered else 1]],
                            "totals": {
                                "lines": 1,
                                "hits": int(covered),
                                "misses": int(not covered),
                                "partials": 0,
                            },
                        }
                    ]
                },
            }
            with self.subTest(covered=covered):
                result, fetch, _ = self.run_gate(
                    response={"coverage_pct": 100, "total_lines": 1, "covered_lines": 1},
                    report=report,
                )
                self.assertEqual(result, expected)
                fetch.assert_not_called()

    def test_stale_or_unfinished_commit_report_cannot_recover_outage(self):
        """Unavailable comparisons never authorize stale or incomplete executable-line data."""
        self.write("selected/value.go", "package selected\nfunc Value() int { return 2 }\n")
        self.commit()
        head = self.git("rev-parse", "HEAD")
        for commit, state in (("old", "complete"), (head, "pending")):
            with self.subTest(commit=commit, state=state):
                result, _, _ = self.run_gate(
                    report={"commitid": commit, "state": state, "report": {"files": []}}
                )
                self.assertEqual(result, 1)
