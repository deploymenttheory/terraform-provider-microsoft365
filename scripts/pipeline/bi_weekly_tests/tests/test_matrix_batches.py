"""Verify general sharding capacity, coverage and deterministic assignment."""

import contextlib
import importlib.util
import json
import os
import re
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from scripts.pipeline.test_support import TemporaryDirectoryTestCase

SCRIPT = Path(__file__).resolve().parents[1] / "generate_test_matrix.py"
SPEC = importlib.util.spec_from_file_location("test_matrix", SCRIPT)
matrix = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(matrix)


class SharderTests(TemporaryDirectoryTestCase):
    """Exercise sharder behavior."""

    def setUp(self):
        super().setUp()
        self.change_directory()

    def generate(self, kind, service, sizes, reverse=False):
        """Generate."""
        names = {}
        for package, size in sizes.items():
            source = Path("internal/services") / kind / service / package
            source.mkdir(parents=True, exist_ok=True)
            names[package] = [f"TestAccCase_{index:03d}" for index in range(size)]
            (source / "resource_test.go").write_text(
                "\n".join(f"func {name}(t *testing.T) {{}}" for name in names[package]),
                encoding="utf-8",
            )
        packages = sorted(sizes, reverse=reverse)
        with patch.object(matrix, "discover_resources", return_value=packages):
            entries = matrix.generate_matrix_for_service(service, kind)
        return entries, names

    def assert_coverage(self, entries, names, limit=25):
        """Assert coverage."""
        selected = {package: [] for package in names}
        for index, entry in enumerate(entries):
            self.assertEqual(entry["shard_index"], index)
            self.assertEqual(entry["total_shards"], len(entries))
            self.assertLessEqual(entry["test_count"], limit)
            self.assertLessEqual(entry["item_count"], matrix.ITEMS_PER_SHARD)
            self.assertGreater(entry["item_count"], 0)
            packages = entry["resources"].split(",")
            self.assertEqual(len(packages), len(set(packages)))
            filters = json.loads(entry["test_filters"])
            for package in packages:
                if package in filters:
                    expression = filters[package]
                    matches = [name for name in names[package] if re.fullmatch(expression, name)]
                    self.assertTrue(matches)
                    self.assertIsNone(re.fullmatch(expression, matches[0] + "_Extra"))
                    selected[package].extend(matches)
                else:
                    selected[package].extend(names[package])
        for package, expected in names.items():
            self.assertCountEqual(selected[package], expected)
            self.assertEqual(len(selected[package]), len(set(selected[package])))

    def test_large_packages_split_in_every_service_and_kind(self):
        """Verify large packages split in every service and kind."""
        for kind in ("resources", "datasources", "actions", "list-resources", "ephemerals"):
            with self.subTest(kind=kind):
                entries, names = self.generate(
                    kind,
                    "arbitrary_service",
                    {
                        "graph_beta/large": 61,
                        "graph_beta/small": 7,
                    },
                )
                self.assert_coverage(entries, names)
                self.assertTrue(any(json.loads(entry["test_filters"]) for entry in entries))

    def test_large_and_small_packages_share_the_same_rendezvous_distribution(self):
        """Verify large and small packages share the same rendezvous distribution."""
        with patch.object(
            matrix, "distribute_rendezvous", wraps=matrix.distribute_rendezvous
        ) as distribute:
            entries, names = self.generate(
                "resources",
                "example",
                {
                    "graph_beta/large": 61,
                    "graph_beta/small": 7,
                },
            )
        identities = distribute.call_args.args[0]
        self.assertIn("graph_beta/small", identities)
        self.assertTrue(any(identity.startswith("graph_beta/large::") for identity in identities))
        self.assert_coverage(entries, names)

    def test_capacity_handles_hash_collisions_and_indivisible_packages(self):
        """Verify capacity handles hash collisions and indivisible packages."""
        with patch.object(matrix, "rendezvous_hash", return_value=1):
            entries, names = self.generate(
                "resources", "example", {f"graph_beta/package_{index}": 13 for index in range(5)}
            )
        self.assertEqual(len(entries), 5)
        self.assert_coverage(entries, names)

    def test_package_count_is_bounded_even_with_tiny_packages(self):
        """Verify package count is bounded even with tiny packages."""
        entries, names = self.generate(
            "resources", "example", {f"graph_beta/package_{index}": 1 for index in range(17)}
        )
        self.assert_coverage(entries, names)

    def test_assignment_is_deterministic_independent_of_discovery_order(self):
        """Verify assignment is deterministic independent of discovery order."""
        sizes = {"graph_beta/large": 61, "graph_beta/another": 19, "graph_beta/small": 7}
        first, names = self.generate("resources", "example", sizes)
        second, _ = self.generate("resources", "example", sizes, reverse=True)
        self.assertEqual(first, second)
        self.assert_coverage(first, names)

    def test_configurable_test_budget(self):
        """Verify configurable test budget."""
        with patch.object(matrix, "TESTS_PER_SHARD", 10):
            entries, names = self.generate("resources", "example", {"graph_beta/large": 31})
        self.assert_coverage(entries, names, limit=10)

    def test_empty_inventory_stays_unfiltered(self):
        """Verify empty inventory stays unfiltered."""
        entries, _ = self.generate("resources", "example", {"graph_beta/unknown": 0})
        self.assertEqual(entries[0]["resources"], "graph_beta/unknown")
        self.assertEqual(json.loads(entries[0]["test_filters"]), {})

    def test_empty_service_has_no_jobs(self):
        """Verify empty service has no jobs."""
        entries, _ = self.generate("resources", "example", {})
        self.assertEqual(entries, [])

    def test_matrix_limit_is_reported_instead_of_dropping_jobs(self):
        """Verify matrix limit is reported instead of dropping jobs."""
        with patch.object(matrix, "generate_matrix_for_service", return_value=[{}] * 257):
            with self.assertRaisesRegex(ValueError, "256"):
                matrix.generate_matrix("resources", ["example"])

    def test_inventory_handles_multiline_tests_fuzz_examples_and_comment_decoys(self):
        """Verify inventory handles multiline tests fuzz examples and comment decoys."""
        Path("resource_test.go").write_text(
            """
// func TestComment(t *testing.T) {}
/* func TestBlock(t *testing.T) {} */
var example = `
func TestString(t *testing.T) {}
`
func TestReal(
    t *testing.T,
) {}
func FuzzSeed(f *testing.F) {}
func ExampleWidget() {}
func TestMain(m *testing.M) {}
""",
            encoding="utf-8",
        )
        self.assertEqual(
            matrix.discover_test_names(Path(".")), ["ExampleWidget", "FuzzSeed", "TestReal"]
        )


if __name__ == "__main__":
    unittest.main()
