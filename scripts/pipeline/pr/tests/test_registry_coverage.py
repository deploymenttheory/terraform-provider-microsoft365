"""Verify safe coverage recovery for uninstrumented provider wiring and models."""

import json
import subprocess
import unittest
from pathlib import Path
from unittest.mock import patch

from support import RepositoryTest

from enforce_coverage import _recover_patch_coverage, _verify_unreachable_tools
from lib.registry_changes import registry_only_additions


BASE_REGISTRY = '''package provider
import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/resource"
)
func (p *M365Provider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{}
}
'''

HEAD_REGISTRY = BASE_REGISTRY.replace(
    '    "context"',
    '    "context"\n    widget "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/example/graph_beta/widget"',
).replace('resource.Resource{}', 'resource.Resource{widget.NewWidgetResource,}')

ROOT = Path(__file__).resolve().parents[4]


class RegistryCoverageTests(RepositoryTest):
    """Reproduce the real gate with Git evidence and complete resource coverage."""

    def fixture(self, model=None):
        """Create a registration-only diff alongside a new resource and optional model."""
        self.write("internal/provider/resources.go", BASE_REGISTRY)
        self.commit()
        base = self.git("rev-parse", "HEAD")
        self.write("internal/provider/resources.go", HEAD_REGISTRY)
        name = "internal/services/resources/example/graph_beta/widget/resource.go"
        self.write(name, '''package widget
import "github.com/hashicorp/terraform-plugin-framework/resource"
func NewWidgetResource() resource.Resource { return nil }
''')
        if model is not None:
            self.write(str(Path(name).with_name("model.go")), model)
        self.commit()
        head = self.git("rev-parse", "HEAD")
        report = {
            "commitid": head,
            "state": "complete",
            "report": {"files": [{
                "name": name,
                "line_coverage": [[3, 0]],
                "totals": {"lines": 1, "hits": 1, "misses": 0, "partials": 0},
            }]},
        }
        args = type("Args", (), {
            "base_ref": base, "commit_sha": head,
            "repo_slug": "example/repo", "codecov_token": None,
        })()
        return args, report

    def test_registration_without_aggregate_registry_coverage(self):
        """The real PR wiring must not require an aggregate provider build."""
        args, report = self.fixture()
        with patch("enforce_coverage.fetch_commit_report", return_value=report):
            stats = _recover_patch_coverage(args)
        self.assertIsNotNone(stats)
        self.assertEqual((stats["coverage_pct"], stats["total_lines"]), (100, 1))

    def test_model_with_only_declarations_needs_no_coverage(self):
        """Go does not emit counters for a file of type declarations."""
        args, report = self.fixture("package widget\ntype Model struct { Name string }\n")
        with patch("enforce_coverage.fetch_commit_report", return_value=report):
            stats = _recover_patch_coverage(args)
        self.assertIsNotNone(stats)
        self.assertEqual((stats["coverage_pct"], stats["total_lines"]), (100, 1))

    def test_function_in_a_model_file_still_requires_coverage(self):
        """Filename conventions cannot hide executable code."""
        args, report = self.fixture("package widget\nfunc Hidden() int { return 1 }\n")
        with patch("enforce_coverage.fetch_commit_report", return_value=report):
            self.assertIsNone(_recover_patch_coverage(args))

    def test_missing_resource_coverage_still_fails_closed(self):
        """A registry proof cannot waive coverage for the registered resource."""
        args, report = self.fixture()
        report["report"]["files"] = []
        with patch("enforce_coverage.fetch_commit_report", return_value=report):
            self.assertIsNone(_recover_patch_coverage(args))


class RegistryShapeTests(unittest.TestCase):
    """Reject executable changes, alias replacements, and arbitrary file exemptions."""

    def test_literal_registration_accepts_comments_and_layout(self):
        """Comments and formatting do not alter constructor wiring."""
        after = HEAD_REGISTRY.replace("return []func()", "// Register resources\n    return    []func()")
        self.assertTrue(registry_only_additions("internal/provider/resources.go", BASE_REGISTRY, after))

    def test_executable_and_non_additive_changes_are_rejected(self):
        """The proof must not become a general registry-file exclusion."""
        variants = {
            "constructor invocation": HEAD_REGISTRY.replace("widget.NewWidgetResource,", "widget.NewWidgetResource(),"),
            "extra statement": HEAD_REGISTRY.replace("return []func()", "ctx = context.Background(); return []func()"),
            "changed existing import": HEAD_REGISTRY.replace("terraform-plugin-framework/resource", "other/resource"),
            "new initializer": HEAD_REGISTRY + "\nfunc init() { panic(\"side effect\") }\n",
            "new variable": HEAD_REGISTRY + "\nvar hidden = execute()\n",
            "duplicate entry": HEAD_REGISTRY.replace("widget.NewWidgetResource,", "widget.NewWidgetResource, widget.NewWidgetResource,"),
            "wrong namespace": HEAD_REGISTRY.replace("internal/services/resources/", "internal/core/"),
            "mock registration": HEAD_REGISTRY.replace("graph_beta/widget\"", "graph_beta/widget/mocks\""),
        }
        for label, after in variants.items():
            with self.subTest(label=label):
                self.assertFalse(registry_only_additions("internal/provider/resources.go", BASE_REGISTRY, after))

    def test_existing_registration_cannot_be_removed_or_replaced(self):
        """Existing constructor references must remain unchanged and in order."""
        changed = HEAD_REGISTRY.replace("widget.NewWidgetResource", "widget.NewOtherResource")
        self.assertFalse(registry_only_additions("internal/provider/resources.go", HEAD_REGISTRY, changed))
        self.assertFalse(registry_only_additions("internal/provider/provider.go", BASE_REGISTRY, HEAD_REGISTRY))

    def test_all_known_registry_interfaces(self):
        """Apply the same restricted proof to the five provider registry interfaces."""
        for filename, method, package, interface, service in (
            ("resources.go", "Resources", "resource", "Resource", "resources"),
            ("datasources.go", "DataSources", "datasource", "DataSource", "datasources"),
            ("actions.go", "Actions", "action", "Action", "actions"),
            ("list_resources.go", "ListResources", "list", "ListResource", "list-resources"),
            ("ephemeral_resources.go", "EphemeralResources", "ephemeral", "EphemeralResource", "ephemerals"),
        ):
            with self.subTest(filename=filename):
                def variant(source):
                    return source.replace("Resources(", method + "(").replace(
                        "resource.Resource", package + "." + interface
                    ).replace("framework/resource\"", "framework/" + package + "\"").replace(
                        "services/resources/", "services/" + service + "/"
                    ).replace("widget.NewWidgetResource", "widget.NewWidget" + interface)
                self.assertTrue(registry_only_additions("internal/provider/" + filename, variant(BASE_REGISTRY), variant(HEAD_REGISTRY)))

    def test_go_parser_rejects_executable_or_invalid_declarations(self):
        """Use the actual Go parser, including variable initializers and function literals."""
        sources = {
            "types.go": "package p; type Model struct { Name string }; const Version = 1",
            "function.go": "package p; func Value() int { return 1 }",
            "variable.go": "package p; var Value = call()",
            "closure.go": "package p; var Value = func() int { return 1 }()",
            "side_effect.go": 'package p; import _ "example.com/side-effect"',
            "invalid.go": "package p; func",
        }
        result = subprocess.run(
            ["go", "run", str(ROOT / "scripts/pipeline/pr/tools/declarations/main.go")],
            input=json.dumps(sources), capture_output=True, text=True, check=True,
        )
        self.assertEqual(json.loads(result.stdout), ["types.go"])

    def test_ci_tools_exclusion_requires_policy_and_runtime_proof(self):
        """Honor the existing scripts policy only outside the production import closure."""
        name = "scripts/pipeline/pr/tools/declarations/main.go"
        with (
            patch("enforce_coverage.Path.read_text", return_value='ignore: ["scripts/"]'),
            patch("enforce_coverage.Path.is_file", return_value=True),
            patch("enforce_coverage.production_import_closure", return_value=set()) as closure,
        ):
            self.assertTrue(_verify_unreachable_tools([name]))
            closure.return_value = {str(Path(name).parent).replace("\\", "/")}
            self.assertFalse(_verify_unreachable_tools([name]))
            self.assertFalse(_verify_unreachable_tools(["internal/runtime.go"]))
        with patch("enforce_coverage.Path.read_text", return_value="ignore: []"):
            with self.assertRaises(ValueError):
                _verify_unreachable_tools([name])
