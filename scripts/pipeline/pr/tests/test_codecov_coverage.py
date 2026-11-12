"""Exercise public/private Codecov access through the actual coverage gate."""

import io
import json
import os
import subprocess
import sys
import unittest
from unittest.mock import patch
from urllib.error import HTTPError

import support  # pylint: disable=unused-import  # Configures fixture import paths.

# isort: split

from enforce_coverage import main


class CodecovCoverageTests(unittest.TestCase):
    """Exercise codecov coverage behavior."""

    def run_gate(self, response, token=None, commit=None, base_ref=None, changed_files=None):
        """Run gate."""
        argv = ["enforce_coverage.py", "--repo-slug", "example/provider", "--pr-number", "123"]
        if commit:
            argv.extend(["--commit-sha", commit])
        if base_ref:
            argv.extend(["--base-ref", base_ref])
        responses = iter(response) if isinstance(response, list) else None
        environment = {} if token is None else {"CODECOV_API_TOKEN": token}
        requests = []

        def respond(request, **_kwargs):
            requests.append(request)
            item = next(responses) if responses is not None else response
            if isinstance(item, Exception):
                raise item
            return io.BytesIO(json.dumps(item).encode())

        with patch.object(sys, "argv", argv), patch.dict(os.environ, environment, clear=True):
            with patch(
                "enforce_coverage.load_pr_checks_config",
                return_value={"coverage_threshold": {"minimum_pct": 60}},
            ):
                with patch("lib.get_codecov_coverage.urllib.request.urlopen", side_effect=respond):
                    with (
                        patch("lib.get_codecov_coverage.time.sleep"),
                        patch(
                            "enforce_coverage.get_changed_files", return_value=changed_files or []
                        ),
                        patch("sys.stdout", new_callable=io.StringIO),
                    ):
                        result = main()
        return result, requests

    @staticmethod
    def comparison_response(coverage, hits, misses, partials=0, state="complete"):
        """Model the complete comparison shape returned by the live Codecov API."""
        lines = hits + misses + partials
        return {
            "head_commit": "current-head",
            "commit_uploads": [{"commitid": "current-head", "state": state}],
            "totals": {
                "patch": {"coverage": 0, "hits": 0, "misses": 0, "partials": 0},
                "head": {
                    "diff": [
                        int(lines > 0),
                        lines,
                        hits,
                        misses,
                        partials,
                        coverage,
                        0,
                        0,
                        0,
                        0,
                        0,
                        0,
                        0,
                    ]
                },
            },
        }

    def test_live_comparison_diff_totals_enforce_the_existing_threshold(self):
        """Use actual diff totals instead of the empty comparison patch placeholder."""
        for coverage, hits, misses, partials, expected in [
            ("100", 5, 0, 0, 0),
            ("66.66667", 6, 2, 1, 0),
            ("59.99", 5999, 4001, 0, 1),
            ("59.999", 59999, 40001, 0, 1),
            ("0", 0, 10, 0, 1),
        ]:
            with self.subTest(coverage=coverage):
                result, requests = self.run_gate(
                    self.comparison_response(coverage, hits, misses, partials),
                    commit="current-head",
                )
                self.assertEqual(result, expected)
                self.assertEqual(len(requests), 1)

    def test_unprocessed_or_malformed_diff_cannot_pass(self):
        """Incomplete uploads and invalid totals cannot be mistaken for coverage."""
        pending = self.comparison_response("100", 5, 0, state="pending")
        malformed = self.comparison_response("100", 5, 0)
        malformed["totals"]["head"]["diff"] = [1, 5]
        for response in [pending, malformed]:
            with self.subTest(response=response):
                result, requests = self.run_gate(response, commit="current-head")
                self.assertEqual(result, 1)
                self.assertGreater(len(requests), 1)

    def test_completed_empty_diff_requires_proof_of_test_only_changes(self):
        """Zero executable lines are applicable only to verified non-production changes."""
        empty = self.comparison_response(None, 0, 0)
        result, _ = self.run_gate(
            empty,
            commit="current-head",
            base_ref="base-head",
            changed_files=["internal/example/resource_acceptance_test.go"],
        )
        self.assertEqual(result, 0)
        result, _ = self.run_gate(
            empty,
            commit="current-head",
            base_ref="base-head",
            changed_files=["internal/example/resource.go"],
        )
        self.assertEqual(result, 1)
        result, _ = self.run_gate(empty, commit="current-head")
        self.assertEqual(result, 1)

    def test_empty_patch_cannot_pass_when_git_verification_fails(self):
        """Missing Git history must not turn unavailable coverage into a passing gate."""
        from enforce_coverage import (  # pylint: disable=import-outside-toplevel
            _verify_test_only_changes,
        )

        with (
            patch(
                "enforce_coverage.get_changed_files",
                side_effect=subprocess.CalledProcessError(1, "git"),
            ),
            patch("sys.stdout", new_callable=io.StringIO),
        ):
            self.assertFalse(_verify_test_only_changes("missing-base"))

    def test_public_coverage_passes_without_authorization_header(self):
        """Verify public coverage passes without authorization header."""
        for token in [None, ""]:
            with self.subTest(token=token):
                result, requests = self.run_gate(
                    {"patch": {"coverage": 71.9, "hits": 87, "misses": 30, "partials": 4}}, token
                )
                self.assertEqual(result, 0)
                self.assertEqual(len(requests), 1)
                self.assertNotIn("Authorization", dict(requests[0].header_items()))
                self.assertEqual(requests[0].get_header("Accept"), "application/json")

    def test_private_coverage_keeps_supplied_api_token(self):
        """Verify private coverage keeps supplied api token."""
        result, requests = self.run_gate(
            {"patch": {"coverage": 80, "hits": 8, "misses": 2}}, "test-api-token"
        )
        self.assertEqual(result, 0)
        self.assertEqual(requests[0].get_header("Authorization"), "bearer test-api-token")

    def test_below_threshold_fails_without_token(self):
        """Verify below threshold fails without token."""
        for coverage in [0, 59.99]:
            with self.subTest(coverage=coverage):
                result, _ = self.run_gate(
                    {"patch": {"coverage": coverage, "hits": 0, "misses": 10}}
                )
                self.assertEqual(result, 1)

    def test_exact_threshold_passes_without_token(self):
        """Verify exact threshold passes without token."""
        result, _ = self.run_gate({"patch": {"coverage": 60, "hits": 6, "misses": 4}})
        self.assertEqual(result, 0)

    def test_missing_coverage_fails_instead_of_skipping_the_gate(self):
        """Verify missing coverage fails instead of skipping the gate."""
        for response in [{}, {"patch": None}, {"patch": {}}]:
            with self.subTest(response=response):
                result, requests = self.run_gate(response)
                self.assertEqual(result, 1)
                self.assertGreater(len(requests), 1)

    def test_stale_passing_coverage_cannot_pass_the_current_commit(self):
        """Reject old green reports until the exact PR head has been processed."""
        response = {
            "head_commit": "old-head",
            "totals": {"patch": {"coverage": 100, "hits": 10, "misses": 0}},
        }
        result, requests = self.run_gate(response, commit="current-head")
        self.assertEqual(result, 1)
        self.assertGreater(len(requests), 1)
        self.assertTrue(all("compare/?pullid=123" in request.full_url for request in requests))

    def test_retry_uses_the_current_commit_coverage_and_preserves_the_threshold(self):
        """An old report cannot mask the current commit's failing or passing patch."""
        for coverage, expected in ((59.99, 1), (60, 0)):
            with self.subTest(coverage=coverage):
                responses = [
                    {"head_commit": "old-head", "totals": {"patch": {"coverage": 100}}},
                    {
                        "head_commit": "current-head",
                        "totals": {
                            "patch": {
                                "coverage": coverage,
                                "hits": 6,
                                "misses": 4,
                            }
                        },
                    },
                ]
                result, requests = self.run_gate(responses, commit="current-head")
                self.assertEqual(result, expected)
                self.assertEqual(len(requests), 2)

    def test_current_commit_without_processed_patch_cannot_pass(self):
        """Require a completed patch report as well as a matching commit."""
        result, requests = self.run_gate(
            {"head_commit": "current-head", "totals": {}}, commit="current-head"
        )
        self.assertEqual(result, 1)
        self.assertGreater(len(requests), 1)

    def test_access_denied_fails_without_retrying_anonymously(self):
        """Verify access denied fails without retrying anonymously."""
        for code in [401, 403]:
            for token in [None, "test-api-token"]:
                with self.subTest(code=code, token=token):
                    error = HTTPError(
                        "https://api.codecov.io/example", code, "denied", {}, io.BytesIO(b"")
                    )
                    self.addCleanup(error.close)
                    result, requests = self.run_gate(error, token)
                    self.assertEqual(result, 1)
                    self.assertEqual(len(requests), 1)


if __name__ == "__main__":
    unittest.main()
