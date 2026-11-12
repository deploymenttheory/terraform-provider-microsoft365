"""Verify isolated factories retain fixture dependencies without a full-provider import."""

import contextlib
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from generate_test_factories import MODULE, generate, inventory, selected_entries

from scripts.pipeline.test_support import TemporaryDirectoryTestCase

# pylint: enable=wrong-import-position


class TestFactories(TemporaryDirectoryTestCase):
    """Exercise test factories behavior."""

    def setUp(self):
        super().setUp()
        self.root = Path(self.directory)
        for filename in (
            "resources.go",
            "datasources.go",
            "actions.go",
            "list_resources.go",
            "ephemeral_resources.go",
        ):
            self.write("internal/provider/" + filename, "package provider\n")
        self.directory = Path("internal/services/resources/example/first")
        self.write(
            str(self.directory / "resource.go"),
            'package first\nconst ResourceName = "microsoft365_first"\n',
        )
        self.write(
            "internal/services/resources/example/second/resource.go",
            'package second\nconst ResourceName = "microsoft365_second"\n',
        )
        self.write(
            "internal/services/datasources/example/second/datasource.go",
            'package second\nconst DataSourceName = "microsoft365_second"\n',
        )
        self.write(
            "internal/provider/resources.go",
            """package provider
import (
 first "{MODULE}/internal/services/resources/example/first"
 second "{MODULE}/internal/services/resources/example/second"
)
func resources() { return []func() resource.Resource{
 first.NewFirstResource,
 second.NewSecondResource,
} }
""",
        )
        self.write(
            "internal/provider/datasources.go",
            """package provider
import second "{MODULE}/internal/services/datasources/example/second"
func datasources() { return []func() datasource.DataSource{
 second.NewSecondDataSource,
} }
""",
        )
        self.write(
            str(self.directory / "resource_test.go"),
            """package first_test
var factories = mocks.TestUnitTestProtoV6ProviderFactories
""",
        )

    def write(self, path, content):
        """Write a UTF-8 file into the temporary fixture repository."""
        target = self.root / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content.replace("{MODULE}", MODULE), encoding="utf-8")

    def test_selected_package_does_not_import_unrelated_registrations(self):
        """Verify selected package does not import unrelated registrations."""
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual([entry[5] for entry in chosen], ["microsoft365_first"])

    def test_cross_resource_fixtures_and_data_sources_are_partitioned_by_kind(self):
        """Verify cross resource fixtures and data sources are partitioned by kind."""
        self.write(
            str(self.directory / "tests/lifecycle.tf"),
            'data "microsoft365_second" "dependency" {}\n',
        )
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual(
            [(entry[0], entry[5]) for entry in chosen],
            [("Resources", "microsoft365_first"), ("DataSources", "microsoft365_second")],
        )

    def test_escaped_inline_go_configuration_keeps_dependencies(self):
        """Verify escaped inline go configuration keeps dependencies."""
        self.write(
            str(self.directory / "inline_test.go"),
            r"""package first_test
var config = "resource \"microsoft365_second\" \"dependency\" {}"
""",
        )
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual(
            [entry[5] for entry in chosen], ["microsoft365_first", "microsoft365_second"]
        )

    def test_generation_is_repeatable_and_detects_new_fixture_dependencies(self):
        """Verify generation is repeatable and detects new fixture dependencies."""
        generate(self.root)
        target = self.root / self.directory / "provider_factories_test.go"
        original = target.read_bytes()
        generate(self.root, check=True)
        self.assertEqual(target.read_bytes(), original)
        self.assertNotIn('internal/provider"', target.read_text(encoding="utf-8"))
        self.write(
            str(self.directory / "new.tf"), 'resource "microsoft365_second" "dependency" {}\n'
        )
        with self.assertRaisesRegex(ValueError, "Stale isolated factories"):
            generate(self.root, check=True)
        self.assertEqual(target.read_bytes(), original)
        generate(self.root)
        self.assertIn("NewSecondResource", target.read_text(encoding="utf-8"))

    def test_unreadable_metadata_fails_instead_of_silently_omitting_constructor(self):
        """Verify unreadable metadata fails instead of silently omitting constructor."""
        self.write(str(self.directory / "resource.go"), "package first\n")
        with self.assertRaisesRegex(ValueError, "Expected one Terraform type name"):
            inventory(self.root)

    def test_deleting_consumers_marks_orphaned_factories_stale(self):
        """Verify deleting consumers marks orphaned factories stale."""
        generate(self.root)
        (self.root / self.directory / "resource_test.go").unlink()
        target = self.root / self.directory / "provider_factories_test.go"
        with self.assertRaisesRegex(ValueError, "Stale isolated factories"):
            generate(self.root, check=True)
        generate(self.root)
        self.assertFalse(target.exists())

    def test_generator_does_not_overwrite_manual_factories(self):
        """Verify generator does not overwrite manual factories."""
        self.write(str(self.directory / "provider_factories_test.go"), "package first_test\n")
        with self.assertRaisesRegex(ValueError, "manually maintained"):
            generate(self.root)

    def test_list_resource_includes_matching_managed_resource_schema_dependency(self):
        """Verify list resource includes matching managed resource schema dependency."""
        self.write(
            "internal/services/list-resources/example/second/list.go",
            'package second\nconst ListResourceName = "microsoft365_second"\n',
        )
        self.write(
            "internal/provider/list_resources.go",
            """package provider
import second "{MODULE}/internal/services/list-resources/example/second"
func lists() { return []func() list.ListResource{
 second.NewSecondListResource,
} }
""",
        )
        self.write(str(self.directory / "list.tf"), 'list "microsoft365_second" "query" {}\n')
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual(
            [(entry[0], entry[5]) for entry in chosen],
            [
                ("Resources", "microsoft365_first"),
                ("Resources", "microsoft365_second"),
                ("ListResources", "microsoft365_second"),
            ],
        )


if __name__ == "__main__":
    unittest.main()
