"""Verify exact-head Codecov line coverage against the entire PR patch."""

import copy
import io
import sys
import urllib.error
import unittest
from unittest.mock import patch

import support  # pylint: disable=unused-import

# isort: split

from enforce_coverage import _recover_patch_coverage, main
from lib.patch_lines import added_go_lines, measure_patch


def commit_report(statuses):
    """Build a complete executable-line report with Codecov's actual status encoding."""
    return {
        "commitid": "head",
        "state": "complete",
        "report": {
            "files": [
                {
                    "name": "internal/runtime.go",
                    "line_coverage": [[i + 1, status] for i, status in enumerate(statuses)],
                    "totals": {
                        "lines": len(statuses),
                        "hits": statuses.count(0),
                        "misses": statuses.count(1),
                        "partials": statuses.count(2),
                    },
                }
            ]
        },
    }


class PatchLineTests(unittest.TestCase):
    """Reject missing, stale, inconsistent and partially covered data."""

    def test_partial_and_missed_lines_are_uncovered(self):
        """Verify partial and missed lines are uncovered."""
        result = measure_patch(
            commit_report([0, 0, 1, 2]), {"internal/runtime.go": {1, 2, 3, 4, 5}}, "head"
        )
        self.assertEqual((result["coverage_pct"], result["total_lines"]), (50, 4))

    def test_only_added_executable_lines_are_measured(self):
        """Verify only added executable lines are measured."""
        result = measure_patch(commit_report([0, 1, 2]), {"internal/runtime.go": {1, 4}}, "head")
        self.assertEqual((result["coverage_pct"], result["total_lines"]), (100, 1))

    def test_missing_changed_file_fails_closed(self):
        """Verify missing changed file fails closed."""
        with self.assertRaises(ValueError):
            measure_patch(commit_report([0]), {"new.go": {1}}, "head")

    def test_missing_or_malformed_report_fails_closed(self):
        """Verify missing or malformed report fails closed."""
        for report in (None, {}, {"commitid": "head", "state": "complete", "report": []}):
            with self.subTest(report=report), self.assertRaises(ValueError):
                measure_patch(report, {"internal/runtime.go": {1}}, "head")

    def test_unavailable_api_and_unverified_mocks_cannot_recover(self):
        """Verify unavailable api and unverified mocks cannot recover."""
        args = type(
            "Args",
            (),
            {
                "base_ref": "origin/main",
                "commit_sha": "head",
                "repo_slug": "example/repo",
                "codecov_token": None,
            },
        )()
        with (
            patch("enforce_coverage.added_go_lines", return_value={"internal/runtime.go": {1}}),
            patch(
                "enforce_coverage.fetch_commit_report",
                side_effect=urllib.error.URLError("unavailable"),
            ),
            patch("sys.stdout", new_callable=io.StringIO),
        ):
            self.assertIsNone(_recover_patch_coverage(args))
        with (
            patch("enforce_coverage.added_go_lines", return_value={"internal/mocks/mock.go": {1}}),
            patch("enforce_coverage._verify_unreachable_mocks", return_value=False),
            patch("enforce_coverage.fetch_commit_report") as fetch,
            patch("sys.stdout", new_callable=io.StringIO),
        ):
            self.assertIsNone(_recover_patch_coverage(args))
            fetch.assert_not_called()

    def test_stale_or_incomplete_commit_fails_closed(self):
        """Verify stale or incomplete commit fails closed."""
        for key, value in (("commitid", "old"), ("state", "pending")):
            report = commit_report([0])
            report[key] = value
            with self.subTest(key=key), self.assertRaises(ValueError):
                measure_patch(report, {"internal/runtime.go": {1}}, "head")

    def test_incomplete_duplicate_or_invalid_line_data_fails_closed(self):
        """Verify incomplete duplicate or invalid line data fails closed."""
        original = commit_report([0, 1, 2])
        for rows in (
            [[1, 0]],
            [[1, 0], [1, 1], [3, 2]],
            [[1, 0], [2, 1], [3, 3]],
            [[True, 0], [2, 1], [3, 2]],
        ):
            report = copy.deepcopy(original)
            report["report"]["files"][0]["line_coverage"] = rows
            with self.subTest(rows=rows), self.assertRaises(ValueError):
                measure_patch(report, {"internal/runtime.go": {1, 2, 3}}, "head")

    def test_inconsistent_totals_fail_closed(self):
        """Verify inconsistent totals fail closed."""
        report = commit_report([0])
        report["report"]["files"][0]["totals"]["lines"] = 2
        with self.assertRaises(ValueError):
            measure_patch(report, {"internal/runtime.go": {1}}, "head")

    def test_duplicate_file_reports_fail_closed(self):
        """Verify duplicate file reports fail closed."""
        report = commit_report([0])
        report["report"]["files"] *= 2
        with self.assertRaises(ValueError):
            measure_patch(report, {"internal/runtime.go": {1}}, "head")

    def test_git_patch_uses_merge_base_and_preserves_zero_length_hunks(self):
        """Verify git patch uses merge base and preserves zero length hunks."""
        with (
            patch("lib.patch_lines.subprocess.run") as run,
            patch("lib.patch_lines.get_merge_base", return_value="base"),
            patch("lib.patch_lines.get_changed_files", return_value=["runtime.go", "unit_test.go"]),
        ):
            run.return_value.stdout = "head"
            run.side_effect = [
                type("Result", (), {"stdout": "head"})(),
                type("Result", (), {"stdout": "@@ -1 +2,2 @@\n@@ -5 +8,0 @@\n@@ -9 +10 @@"})(),
            ]
            self.assertEqual(added_go_lines("origin/main", "head"), {"runtime.go": {2, 3, 10}})
            self.assertIn("base", run.call_args.args[0])

    def test_checked_out_commit_must_match(self):
        """Verify checked out commit must match."""
        with patch("lib.patch_lines.subprocess.run") as run:
            run.return_value.stdout = "wrong"
            with self.assertRaises(ValueError):
                added_go_lines("origin/main", "head")

    def test_gate_recovers_valid_patch_without_lowering_threshold(self):
        """Verify gate recovers valid patch without lowering threshold."""
        for statuses, expected in (([0, 0, 0, 1, 2], 0), ([0, 0, 1, 2], 1)):
            with (
                self.subTest(statuses=statuses),
                patch.object(
                    sys,
                    "argv",
                    [
                        "enforce_coverage.py",
                        "--repo-slug",
                        "example/repo",
                        "--pr-number",
                        "1",
                        "--commit-sha",
                        "head",
                        "--base-ref",
                        "origin/main",
                    ],
                ),
                patch(
                    "enforce_coverage.fetch_codecov_coverage",
                    return_value={"no_executable_changes": True},
                ),
                patch("enforce_coverage._verify_test_only_changes", return_value=False),
                patch(
                    "enforce_coverage.added_go_lines",
                    return_value={"internal/runtime.go": set(range(1, len(statuses) + 1))},
                ),
                patch("enforce_coverage.fetch_commit_report", return_value=commit_report(statuses)),
                patch(
                    "enforce_coverage.load_pr_checks_config",
                    return_value={"coverage_threshold": {"minimum_pct": 60}},
                ),
                patch("enforce_coverage.write_output"),
                patch("sys.stdout", new_callable=io.StringIO),
            ):
                self.assertEqual(main(), expected)
