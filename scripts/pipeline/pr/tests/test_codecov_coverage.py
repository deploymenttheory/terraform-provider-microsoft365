"""Exercise public/private Codecov access through the actual coverage gate."""

import io
import json
import os
import subprocess
import sys
import unittest
from contextlib import nullcontext
from pathlib import Path
from unittest.mock import patch
from urllib.error import HTTPError

import support  # pylint: disable=unused-import  # Configures fixture import paths.
from support import RepositoryTest

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
        environment = {
            name: os.environ[name]
            for name in (
                "PATH",
                "SYSTEMROOT",
                "WINDIR",
                "COMSPEC",
                "TEMP",
                "TMP",
                "TMPDIR",
                "HOME",
                "USERPROFILE",
            )
            if name in os.environ
        }
        if token is not None:
            environment["CODECOV_API_TOKEN"] = token
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
                        (
                            patch(
                                "enforce_coverage.get_changed_files",
                                return_value=changed_files or [],
                            )
                            if changed_files is not None
                            else nullcontext()
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


class MockOnlyCoverageTests(RepositoryTest):
    """Use real Git diffs and runtime import graphs before accepting an empty mock patch."""

    def setUp(self):
        super().setUp()
        self.module = "github.com/deploymenttheory/terraform-provider-microsoft365"
        self.mock = "internal/example/mocks/responders.go"
        self.write("go.mod", f"module {self.module}\n\ngo 1.26.0\n")
        self.write("main.go", f'package main\nimport "{self.module}/internal/provider"\n')
        self.write("internal/provider/provider.go", "package provider\n")
        self.write(self.mock, "package mocks\nfunc Response() int { return 1 }\n")
        self.write(
            "internal/example/resource_test.go",
            f'package example_test\nimport "{self.module}/internal/example/mocks"\n',
        )
        self.write("codecov.yml", 'ignore:\n  - "**/mocks/**"\n')
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 2 }\n")

    def empty_gate(self):
        """Exercise the real gate with a completed Codecov report containing no changed lines."""
        self.commit()
        result, _ = CodecovCoverageTests().run_gate(
            CodecovCoverageTests.comparison_response(None, 0, 0),
            commit="current-head",
            base_ref=self.base,
        )
        return result

    def test_unreachable_ignored_mock_package_accepts_completed_empty_patch(self):
        """A test-only mock package is not executable provider runtime code."""
        self.assertEqual(self.empty_gate(), 0)

    def test_any_runtime_source_change_still_requires_executable_coverage(self):
        """Changing a real resource alongside its mocks must not bypass coverage."""
        self.write(
            "internal/example/resource.go", "package example\nfunc Value() int { return 1 }\n"
        )
        self.assertEqual(self.empty_gate(), 1)

    def test_direct_production_import_of_mock_package_rejects_empty_patch(self):
        """A namespace convention cannot override actual runtime reachability."""
        self.write(
            "internal/provider/provider.go",
            f'package provider\nimport "{self.module}/internal/example/mocks"\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 3 }\n")
        self.assertEqual(self.empty_gate(), 1)

    def test_indirect_raw_production_import_rejects_empty_patch(self):
        """Trace runtime dependencies through raw-string and dot-import declarations."""
        self.write(
            "internal/provider/provider.go",
            f"package provider\nimport `{self.module}/internal/bridge`\n",
        )
        self.write(
            "internal/bridge/bridge.go",
            f'package bridge\nimport . "{self.module}/internal/example/mocks"\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 4 }\n")
        self.assertEqual(self.empty_gate(), 1)

    def test_commented_imports_do_not_make_test_mock_runtime_code(self):
        """Commented code must not create production graph edges."""
        self.write(
            "internal/provider/provider.go",
            f'package provider\n// import "{self.module}/internal/example/mocks"\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 5 }\n")
        self.assertEqual(self.empty_gate(), 0)

    def test_missing_runtime_entry_point_or_coverage_policy_rejects_empty_patch(self):
        """Missing proof must not authorize an empty-patch exception."""
        for missing in ("main.go", "internal/provider/provider.go", "codecov.yml"):
            with self.subTest(missing=missing):
                original = Path(missing).read_text(encoding="utf-8")
                self.write(self.mock, "package mocks\nfunc Response() int { return 1 }\n")
                Path(missing).unlink()
                self.commit()
                self.base = self.git("rev-parse", "HEAD")
                self.write(self.mock, "package mocks\nfunc Response() int { return 2 }\n")
                self.assertEqual(self.empty_gate(), 1)
                self.write(missing, original)

    def test_coverage_policy_must_explicitly_ignore_mock_namespace(self):
        """Reachability proof cannot invent a new coverage exclusion policy."""
        self.write("codecov.yml", "ignore: []\n")
        self.assertEqual(self.empty_gate(), 1)

    def test_escaped_runtime_import_rejects_empty_patch(self):
        """Go string escapes cannot hide a runtime dependency from the verifier."""
        escaped = self.module.replace("github", r"\x67ithub")
        self.write(
            "internal/provider/provider.go",
            f'package provider\nimport "{escaped}/internal/example/mocks"\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 6 }\n")
        self.assertEqual(self.empty_gate(), 1)

    def test_comments_between_import_tokens_cannot_hide_runtime_mock(self):
        """Go comments act as separators rather than merging the import keyword and alias."""
        self.write(
            "internal/provider/provider.go",
            "package"
            f' provider\nimport/*separator*/_/*separator*/"{self.module}/internal/example/mocks"\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 7 }\n")
        self.assertEqual(self.empty_gate(), 1)

    def test_local_module_alias_to_mock_sources_cannot_bypass_graph_proof(self):
        """Aliases to sources inside the provider tree require more evidence than import paths."""
        self.write(
            "go.mod",
            f"module {self.module}\n\ngo 1.26.0\nreplace example.com/alias =>"
            " ./internal/example/mocks\n",
        )
        self.assertEqual(self.empty_gate(), 1)

    def test_missing_reachable_package_prevents_runtime_graph_proof(self):
        """An incomplete runtime import graph must fail closed."""
        self.write(
            "internal/provider/provider.go",
            f'package provider\nimport "{self.module}/internal/missing"\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 8 }\n")
        self.assertEqual(self.empty_gate(), 1)

    def test_go_source_text_inside_strings_is_not_a_runtime_import(self):
        """Embedded sample code must not become an import edge."""
        self.write(
            "internal/provider/provider.go",
            f'package provider\nconst Example = `import "{self.module}/internal/example/mocks"`\n',
        )
        self.commit()
        self.base = self.git("rev-parse", "HEAD")
        self.write(self.mock, "package mocks\nfunc Response() int { return 9 }\n")
        self.assertEqual(self.empty_gate(), 0)

    def test_other_uncovered_namespace_cannot_use_mock_exception(self):
        """Do not turn every ignored fixture or unreachable package into a coverage exemption."""
        self.write(
            "internal/example/tests/helper.go", "package tests\nfunc Value() int { return 1 }\n"
        )
        self.assertEqual(self.empty_gate(), 1)


if __name__ == "__main__":
    unittest.main()
