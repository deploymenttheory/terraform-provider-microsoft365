"""Verify package partition completeness, dependent tests, and coverage aggregation."""

import os
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import yaml
from support import RepositoryTest

# isort: split

from lib.partitions import (
    affected_test_packages,
    fixture_packages,
    partition_lint_packages,
    partition_packages,
    select_lint_packages,
)
from merge_coverage import merge_profiles, merge_shard_profiles
from run_lint import run_lint


class PartitionTests(RepositoryTest):
    """Exercise partition behavior."""

    def test_every_package_is_assigned_once_and_shards_are_bounded(self):
        """Verify every package is assigned once and shards are bounded."""
        packages = [f"internal/example/package_{index}" for index in range(243)]
        matrix = partition_packages(packages)
        assigned = [package for shard in matrix["include"] for package in shard["packages"].split()]
        self.assertEqual(sorted(assigned), sorted(packages))
        self.assertEqual(len(assigned), len(set(assigned)))
        self.assertTrue(all(len(shard["packages"].split()) <= 5 for shard in matrix["include"]))
        self.assertEqual(matrix, partition_packages(list(reversed(packages))))

    def test_empty_selection_launches_no_jobs(self):
        """Verify empty selection launches no jobs."""
        self.assertEqual(partition_packages([]), {"include": []})

    def test_production_registry_is_alone_and_every_lint_package_is_assigned_once(self):
        """Verify the aggregate registry cannot crowd out the ordinary package workers."""
        packages = [f"internal/example/package_{index}" for index in range(31)]
        packages += ["internal/provider", "internal/provider/core"]
        matrix = partition_lint_packages(packages)
        assigned = [package for entry in matrix["include"] for package in entry["packages"].split()]
        self.assertEqual(sorted(assigned), sorted(packages))
        self.assertEqual(len(assigned), len(set(assigned)))
        registry = next(
            entry for entry in matrix["include"] if entry["packages"] == "internal/provider"
        )
        self.assertEqual(registry["label"], "Provider Registry")
        self.assertEqual(registry["gomemlimit"], "20GiB")
        self.assertEqual(registry["gogc"], "100")
        for entry in matrix["include"]:
            self.assertLessEqual(len(entry["packages"].split()), 5)
            if entry is not registry:
                self.assertEqual((entry["gomemlimit"], entry["gogc"]), ("", ""))
        self.assertEqual(matrix, partition_lint_packages(reversed(packages)))

    def test_manual_lint_selection_can_only_narrow_the_changed_packages(self):
        """Verify diagnostic dispatches cannot lint an unselected production package."""
        changed = ["internal/provider", "internal/provider/core"]
        self.assertEqual(select_lint_packages(changed), sorted(changed))
        self.assertEqual(select_lint_packages(changed, "internal/provider"), ["internal/provider"])
        for requested in ("", "unrelated", "internal/provider unrelated"):
            with self.assertRaisesRegex(ValueError, "only packages changed"):
                select_lint_packages(changed, requested)
        self.assertEqual(partition_lint_packages([]), {"include": []})

    def test_shared_dependency_change_includes_transitive_tests(self):
        """Verify shared dependency change includes transitive tests."""
        self.write("internal/core/config.go", "package core\n")
        self.write(
            "internal/mocks/mocks.go",
            """package mocks
import core "github.com/deploymenttheory/terraform-provider-microsoft365/internal/core"
""",
        )
        self.write(
            "internal/services/first/resource_test.go",
            """package first_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
""",
        )
        self.write("internal/services/second/resource_test.go", "package second_test\n")
        self.assertEqual(affected_test_packages(["internal/core"]), ["internal/services/first"])

    def test_production_registry_without_tests_is_not_compiled_by_unit_workers(self):
        """Verify production registry without tests is not compiled by unit workers."""
        self.write("internal/provider/registry.go", "package provider\n")
        self.write("internal/provider/core/config_test.go", "package core_test\n")
        self.assertEqual(
            affected_test_packages(["internal/provider", "internal/provider/core"]),
            ["internal/provider/core"],
        )

    def test_lint_shard_cannot_expand_its_changed_package_scope(self):
        """Verify lint shard cannot expand its changed package scope."""
        with (
            patch("run_lint.get_merge_base", return_value=self.base),
            patch("run_lint.get_changed_packages", return_value=["selected"]),
            patch("run_lint.subprocess.run") as run,
        ):
            with self.assertRaisesRegex(ValueError, "only packages changed"):
                run_lint(self.base, selected_packages=["unrelated"])
        run.assert_not_called()

    def test_test_only_changes_do_not_propagate_as_production_dependencies(self):
        """Verify test only changes do not propagate as production dependencies."""
        self.write("internal/shared/helper.go", "package shared\n")
        self.write("internal/shared/helper_test.go", "package shared_test\n")
        self.write(
            "internal/consumer/consumer_test.go",
            """package consumer_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/shared"
""",
        )
        self.assertEqual(
            affected_test_packages(["internal/shared"], production_changed=set()),
            ["internal/shared"],
        )

    def test_dependency_of_shared_helpers_tests_does_not_fan_out_to_all_consumers(self):
        """Verify dependency of shared helpers tests does not fan out to all consumers."""
        self.write("internal/target/resource.go", "package target\n")
        self.write("internal/target/resource_test.go", "package target_test\n")
        self.write("internal/mocks/mocks.go", "package mocks\n")
        self.write(
            "internal/mocks/mocks_test.go",
            """package mocks_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/target"
""",
        )
        self.write(
            "internal/unrelated/resource_test.go",
            """package unrelated_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
""",
        )
        self.assertEqual(
            affected_test_packages(["internal/target"]), ["internal/mocks", "internal/target"]
        )

    def test_changed_fixtures_resolve_to_the_package_that_serves_them(self):
        """Verify changed fixtures resolve to the package that serves them."""
        self.write("internal/resource/mocks/responders.go", "package mocks\n")
        self.write(
            "internal/resource/resource_test.go",
            """package resource_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resource/mocks"
""",
        )
        files = ["internal/resource/mocks/responses/new.json", "docs/example.json"]
        owners = fixture_packages(files)
        self.assertEqual(owners, {"internal/resource/mocks"})
        self.assertEqual(affected_test_packages(sorted(owners)), ["internal/resource"])

    def test_deleted_helper_still_selects_tests_that_import_it(self):
        """Verify deleted helper still selects tests that import it."""
        self.write(
            "internal/consumer/consumer_test.go",
            """package consumer_test
import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/deleted"
""",
        )
        self.assertEqual(
            affected_test_packages([], production_changed={"internal/deleted"}),
            ["internal/consumer"],
        )


class CoverageMergeTests(unittest.TestCase):
    """Exercise coverage merge behavior."""

    def test_overlapping_blocks_are_covered_once_after_shards_are_combined(self):
        """Verify overlapping blocks are covered once after shards are combined."""
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            first, second, output = root / "first", root / "second", root / "merged"
            first.write_text("mode: atomic\nmodule/a.go:1.1,2.1 2 0\n", encoding="utf-8")
            second.write_text(
                "mode: atomic\nmodule/a.go:1.1,2.1 2 3\nmodule/b.go:1.1,2.1 1 1\n", encoding="utf-8"
            )
            merge_profiles([first, second], output)
            self.assertEqual(
                output.read_text(encoding="utf-8"),
                "mode: atomic\nmodule/a.go:1.1,2.1 2 3\nmodule/b.go:1.1,2.1 1 1\n",
            )

    def test_missing_artifacts_cannot_produce_successful_empty_coverage(self):
        """Verify missing artifacts cannot produce successful empty coverage."""
        with self.assertRaisesRegex(ValueError, "No shard coverage"):
            merge_profiles([], Path("must-not-exist"))

    def test_one_successful_shard_cannot_hide_a_missing_shard_artifact(self):
        """Verify one successful shard cannot hide a missing shard artifact."""
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            profile = root / "unit-test-results-0" / "unit-coverage.txt"
            profile.parent.mkdir()
            profile.write_text("mode: atomic\nmodule/a.go:1.1,2.1 1 1\n", encoding="utf-8")
            output = root / "merged"
            with self.assertRaisesRegex(ValueError, "Missing shard coverage profiles"):
                merge_shard_profiles(root, output, 2)
            self.assertFalse(output.exists())

    def test_single_artifact_flat_download_merges_successfully(self):
        """Verify single artifact flat download merges successfully."""
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            profile = root / "unit-coverage.txt"
            profile.write_text("mode: atomic\nmodule/a.go:1.1,2.1 1 1\n", encoding="utf-8")
            output = root / "merged"
            merge_shard_profiles(root, output, 1)
            self.assertEqual(
                output.read_text(encoding="utf-8"), profile.read_text(encoding="utf-8")
            )
            with self.assertRaisesRegex(ValueError, "Missing shard coverage profiles"):
                merge_shard_profiles(root, output, 2)


class WorkflowGateTests(unittest.TestCase):
    """Exercise workflow gate behavior."""

    def test_unit_and_race_workers_use_the_preinstalled_terraform_binary(self):
        """Verify unit and race workers use the preinstalled terraform binary."""

        workflow = Path(__file__).resolve().parents[4] / ".github/workflows/pr-tests.yml"
        jobs = yaml.safe_load(workflow.read_text(encoding="utf-8"))["jobs"]
        for worker in ["unit-test-shards", "race-shards"]:
            with self.subTest(worker=worker):
                steps = jobs[worker]["steps"]
                setup = next(
                    step
                    for step in steps
                    if step.get("uses", "").startswith("hashicorp/setup-terraform@")
                )
                self.assertEqual(setup["with"]["terraform_version"], "latest")
                self.assertFalse(setup["with"]["terraform_wrapper"])
                run = next(
                    step["run"] for step in steps if "monitor_resources.sh" in step.get("run", "")
                )
                self.assertIn('export TF_ACC_TERRAFORM_PATH="$(command -v terraform)"', run)

    def test_unit_and_race_workers_restore_one_successful_prepared_snapshot(self):
        """Use one writer and exact-key readers, preserving preparation failures."""
        workflow = Path(__file__).resolve().parents[4] / ".github/workflows/pr-tests.yml"
        jobs = yaml.safe_load(workflow.read_text(encoding="utf-8"))["jobs"]
        producer = jobs["prepare-build-cache"]
        writers = [
            step for step in producer["steps"] if step.get("uses", "").startswith("actions/cache@")
        ]
        self.assertEqual(len(writers), 1)
        for worker, gate in (("unit-test-shards", "unit-tests"), ("race-shards", "race-detection")):
            self.assertIn("prepare-build-cache", jobs[worker]["needs"])
            self.assertIn("prepare-build-cache", jobs[gate]["needs"])
            self.assertIn("needs.prepare-build-cache.result", str(jobs[gate]["steps"]))
            steps = jobs[worker]["steps"]
            setup = next(
                step for step in steps if step.get("uses", "").startswith("actions/setup-go@")
            )
            self.assertFalse(setup["with"]["cache"])
            readers = [
                step for step in steps if step.get("uses", "").startswith("actions/cache/restore@")
            ]
            self.assertEqual(len(readers), 1)
            self.assertEqual(
                readers[0]["with"]["key"], "${{ needs.prepare-build-cache.outputs.cache-key }}"
            )
            self.assertTrue(readers[0]["with"]["fail-on-cache-miss"])
            self.assertFalse(
                any(step.get("uses", "").startswith("actions/cache@") for step in steps)
            )

    @unittest.skipIf(os.name == "nt", "Gate shell execution is verified on Linux")
    def test_preparation_failure_cannot_be_hidden_by_successful_workers(self):
        """Run the actual shell gates with deliberately contradictory worker outcomes."""
        workflows = Path(__file__).resolve().parents[4] / ".github/workflows"
        scenarios = (
            ("pr-tests.yml", "unit-tests", {"HAS_CHANGES": "true", "CACHE": "failure"}),
            (
                "pr-tests.yml",
                "race-detection",
                {
                    "HAS_GOROUTINES": "true",
                    "RUN_RACE": "true",
                    "CACHE": "failure",
                    "UNIT_RESULT": "success",
                },
            ),
            ("go-lint.yml", "golint", {"HAS_CHANGES": "true", "COMPILE": "failure"}),
        )
        for filename, job, values in scenarios:
            with self.subTest(gate=job):
                jobs = yaml.safe_load((workflows / filename).read_text(encoding="utf-8"))["jobs"]
                command = next(
                    step["run"]
                    for step in jobs[job]["steps"]
                    if 'test "$PLAN"' in step.get("run", "")
                )
                environment = {**os.environ, "PLAN": "success", "RESULT": "success", **values}
                result = subprocess.run(["bash", "-e", "-c", command], env=environment, check=False)
                self.assertNotEqual(result.returncode, 0)

    def test_race_and_unit_workers_run_independently_after_cache_preparation(self):
        """A unit or coverage result must not serialize or suppress race execution."""
        workflow = Path(__file__).resolve().parents[4] / ".github/workflows/pr-tests.yml"
        jobs = yaml.safe_load(workflow.read_text(encoding="utf-8"))["jobs"]
        self.assertEqual(set(jobs["race-shards"]["needs"]), {"plan-tests", "prepare-build-cache"})
        self.assertNotIn("unit-tests", jobs["race-detection"]["needs"])
        self.assertIn("prepare-build-cache", jobs["unit-test-shards"]["needs"])
        self.assertIn("race-shards", jobs["race-detection"]["needs"])

    def test_lint_workers_wait_for_their_successful_prepared_snapshot(self):
        """Prepare only selected graphs and fail the gate if the cache writer fails."""
        workflow = Path(__file__).resolve().parents[4] / ".github/workflows/go-lint.yml"
        jobs = yaml.safe_load(workflow.read_text(encoding="utf-8"))["jobs"]
        self.assertIn("compile-provider", jobs["lint-shards"]["needs"])
        self.assertIn("compile-provider", jobs["golint"]["needs"])
        self.assertIn("needs.compile-provider.result", str(jobs["golint"]["steps"]))
        producer = jobs["compile-provider"]
        self.assertIn("has-changes", producer["if"])
        writers = [
            step for step in producer["steps"] if step.get("uses", "").startswith("actions/cache@")
        ]
        self.assertEqual(len(writers), 1)
        full_build = next(
            step
            for step in producer["steps"]
            if "go build ./internal/provider" in step.get("run", "")
        )
        self.assertIn("has-provider-changes", full_build["if"])
        steps = jobs["lint-shards"]["steps"]
        reader = next(
            step for step in steps if step.get("uses", "").startswith("actions/cache/restore@")
        )
        self.assertEqual(reader["with"]["key"], "${{ needs.compile-provider.outputs.cache-key }}")
        for output in ("mod-cache", "build-cache"):
            self.assertIn(f"needs.compile-provider.outputs.{output}", reader["with"]["path"])
            self.assertIn(output, producer["outputs"])
        self.assertTrue(reader["with"]["fail-on-cache-miss"])
        self.assertFalse(any(step.get("uses", "").startswith("actions/cache@") for step in steps))
        self.assertIn("needs.compile-provider.result == 'success'", jobs["lint-shards"]["if"])

    def test_required_checks_keep_their_names_and_require_all_worker_shards(self):
        """Verify required checks keep their names and require all worker shards."""

        workflows = Path(__file__).resolve().parents[4] / ".github/workflows"
        unit = yaml.safe_load((workflows / "pr-tests.yml").read_text(encoding="utf-8"))["jobs"]
        lint = yaml.safe_load((workflows / "go-lint.yml").read_text(encoding="utf-8"))["jobs"]
        self.assertEqual(
            unit["unit-test-shards"]["name"], "🧪 Unit Test Partition ${{ matrix.shard_index }}"
        )
        for jobs, gate, worker, name in (
            (unit, "unit-tests", "unit-test-shards", "🧪 Go Unit Tests & Coverage"),
            (unit, "race-detection", "race-shards", "🔍 Race Detection"),
            (lint, "golint", "lint-shards", "✨ Run golangci-lint"),
        ):
            with self.subTest(gate=gate):
                self.assertEqual(jobs[gate]["name"], name)
                self.assertIn(worker, jobs[gate]["needs"])
                self.assertEqual(jobs[gate]["if"], "always()")
                self.assertFalse(jobs[worker]["strategy"]["fail-fast"])
                self.assertEqual(jobs[worker]["runs-on"], "ubuntu-24.04-arm")
                self.assertIn(f"needs.{worker}.result", str(jobs[gate]["steps"]))
