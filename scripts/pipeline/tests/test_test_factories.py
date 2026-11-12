"""Verify isolated factories retain fixture dependencies without a full-provider import."""

import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
from generate_test_factories import generate, inventory, selected_entries


class TestFactories(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        for filename in ('resources.go', 'datasources.go', 'actions.go', 'list_resources.go', 'ephemeral_resources.go'):
            self.write('internal/provider/' + filename, 'package provider\n')
        self.directory = Path('internal/services/resources/example/first')
        self.write(str(self.directory / 'resource.go'), 'package first\nconst ResourceName = "microsoft365_first"\n')
        self.write('internal/services/resources/example/second/resource.go', 'package second\nconst ResourceName = "microsoft365_second"\n')
        self.write('internal/services/datasources/example/second/datasource.go', 'package second\nconst DataSourceName = "microsoft365_second"\n')
        self.write('internal/provider/resources.go', '''package provider
import (
 first "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/example/first"
 second "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/example/second"
)
func resources() { return []func() resource.Resource{
 first.NewFirstResource,
 second.NewSecondResource,
} }
''')
        self.write('internal/provider/datasources.go', '''package provider
import second "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/example/second"
func datasources() { return []func() datasource.DataSource{
 second.NewSecondDataSource,
} }
''')
        self.write(str(self.directory / 'resource_test.go'), '''package first_test
var factories = mocks.TestUnitTestProtoV6ProviderFactories
''')

    def write(self, path, content):
        target = self.root / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding='utf-8')

    def test_selected_package_does_not_import_unrelated_registrations(self):
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual([entry[5] for entry in chosen], ['microsoft365_first'])

    def test_cross_resource_fixtures_and_data_sources_are_partitioned_by_kind(self):
        self.write(str(self.directory / 'tests/lifecycle.tf'), 'data "microsoft365_second" "dependency" {}\n')
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual([(entry[0], entry[5]) for entry in chosen],
                         [('Resources', 'microsoft365_first'), ('DataSources', 'microsoft365_second')])

    def test_escaped_inline_go_configuration_keeps_dependencies(self):
        self.write(str(self.directory / 'inline_test.go'), r'''package first_test
var config = "resource \"microsoft365_second\" \"dependency\" {}"
''')
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual([entry[5] for entry in chosen], ['microsoft365_first', 'microsoft365_second'])

    def test_generation_is_repeatable_and_detects_new_fixture_dependencies(self):
        generate(self.root)
        target = self.root / self.directory / 'provider_factories_test.go'
        original = target.read_bytes()
        generate(self.root, check=True)
        self.assertEqual(target.read_bytes(), original)
        self.assertNotIn('internal/provider"', target.read_text())
        self.write(str(self.directory / 'new.tf'), 'resource "microsoft365_second" "dependency" {}\n')
        with self.assertRaisesRegex(ValueError, 'Stale isolated factories'):
            generate(self.root, check=True)
        self.assertEqual(target.read_bytes(), original)
        generate(self.root)
        self.assertIn('NewSecondResource', target.read_text())

    def test_unreadable_metadata_fails_instead_of_silently_omitting_constructor(self):
        self.write(str(self.directory / 'resource.go'), 'package first\n')
        with self.assertRaisesRegex(ValueError, 'Expected one Terraform type name'):
            inventory(self.root)

    def test_deleting_consumers_marks_orphaned_factories_stale(self):
        generate(self.root)
        (self.root / self.directory / 'resource_test.go').unlink()
        target = self.root / self.directory / 'provider_factories_test.go'
        with self.assertRaisesRegex(ValueError, 'Stale isolated factories'):
            generate(self.root, check=True)
        generate(self.root)
        self.assertFalse(target.exists())

    def test_generator_does_not_overwrite_manual_factories(self):
        self.write(str(self.directory / 'provider_factories_test.go'), 'package first_test\n')
        with self.assertRaisesRegex(ValueError, 'manually maintained'):
            generate(self.root)

    def test_list_resource_includes_matching_managed_resource_schema_dependency(self):
        self.write('internal/services/list-resources/example/second/list.go', 'package second\nconst ListResourceName = "microsoft365_second"\n')
        self.write('internal/provider/list_resources.go', '''package provider
import second "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/list-resources/example/second"
func lists() { return []func() list.ListResource{
 second.NewSecondListResource,
} }
''')
        self.write(str(self.directory / 'list.tf'), 'list "microsoft365_second" "query" {}\n')
        chosen = selected_entries(self.root, self.directory, inventory(self.root))
        self.assertEqual([(entry[0], entry[5]) for entry in chosen],
                         [('Resources', 'microsoft365_first'), ('Resources', 'microsoft365_second'),
                          ('ListResources', 'microsoft365_second')])


if __name__ == '__main__':
    unittest.main()
