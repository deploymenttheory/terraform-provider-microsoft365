#!/usr/bin/env python3
"""
Generate provider coverage documentation by extracting metadata from templates and code.

This script:
1. Parses all .md.tmpl files to extract metadata (subcategory, version, status, permissions)
2. Checks for existence of unit and acceptance tests
3. Counts resources, data sources, list resources, ephemerals, and actions
4. Generates a comprehensive provider_coverage.md file organized by service domain
"""

import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, field
from collections import defaultdict
from datetime import datetime, timezone


@dataclass
class ComponentMetadata:
    """Metadata for a Terraform component (resource, data source, etc.)"""
    name: str
    component_type: str
    subcategory: str
    initial_version: Optional[str] = None
    last_updated_version: Optional[str] = None
    status: Optional[str] = None
    required_permissions: List[str] = field(default_factory=list)
    optional_permissions: List[str] = field(default_factory=list)
    example_count: int = 0
    has_unit_tests: bool = False
    has_acceptance_tests: bool = False
    ms_doc_links: List[str] = field(default_factory=list)
    template_path: Optional[str] = None
    code_path: Optional[str] = None


@dataclass
class ProjectPaths:
    """Project directory paths."""
    repo_root: Path
    templates: Path
    internal: Path
    provider: Path
    services: Path


class ProviderCoverageGenerator:
    """Generate provider coverage documentation from templates and code."""

    def __init__(self, repo_root: Path):
        self.paths = ProjectPaths(
            repo_root=repo_root,
            templates=repo_root / "templates",
            internal=repo_root / "internal",
            provider=repo_root / "internal" / "provider",
            services=repo_root / "internal" / "services"
        )
        self.components: List[ComponentMetadata] = []
        self.service_domains: Dict[str, List[ComponentMetadata]] = defaultdict(list)

    def parse_template_frontmatter(self, template_path: Path) -> Dict:
        """Extract YAML frontmatter from a template file."""
        with open(template_path, 'r', encoding='utf-8') as f:
            content = f.read()

        # Extract frontmatter between --- markers
        match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
        if not match:
            return {}

        frontmatter_text = match.group(1)

        # Handle template variables in frontmatter - extract subcategory directly
        subcategory_match = re.search(r'subcategory:\s*"([^"]+)"', frontmatter_text)
        if subcategory_match:
            return {'subcategory': subcategory_match.group(1)}

        return {}

    def extract_version_history(
        self,
        template_path: Path
    ) -> Tuple[Optional[str], Optional[str], Optional[str]]:
        """Extract initial version, last updated version, and status from version history."""
        with open(template_path, 'r', encoding='utf-8') as f:
            content = f.read()

        # Find version history section
        version_section = re.search(
            r'## Version History.*?\n\|.*?\n\|.*?\n((?:\|.*?\n)+)',
            content,
            re.DOTALL
        )

        if not version_section:
            return None, None, None

        # Parse all data rows
        rows = version_section.group(1).strip().split('\n')
        if not rows:
            return None, None, None

        # First row is the initial version (oldest)
        first_match = re.match(r'\|\s*([^\|]+?)\s*\|\s*([^\|]+?)\s*\|', rows[0])
        if not first_match:
            return None, None, None

        initial_version = first_match.group(1).strip()
        status = first_match.group(2).strip()

        # Last row is the most recent update
        last_updated = None
        if len(rows) > 1:
            last_match = re.match(r'\|\s*([^\|]+?)\s*\|', rows[-1])
            if last_match:
                last_updated = last_match.group(1).strip()

        return initial_version, last_updated, status

    def extract_permissions(self, template_path: Path) -> Tuple[List[str], List[str]]:
        """Extract required and optional Graph API permissions."""
        with open(template_path, 'r', encoding='utf-8') as f:
            content = f.read()

        required_perms = []
        optional_perms = []

        # Extract required permissions
        required_match = re.search(
            r'\*\*Required:\*\*\s*\n((?:- `[^`]+`\s*\n)+)',
            content
        )
        if required_match:
            required_perms = re.findall(r'- `([^`]+)`', required_match.group(1))

        # Extract optional permissions
        optional_match = re.search(
            r'\*\*Optional:\*\*\s*\n((?:- `[^`]+`\s*\n)+)',
            content
        )
        if optional_match:
            optional_perms = re.findall(r'- `([^`]+)`', optional_match.group(1))
            # Filter out N/A entries
            optional_perms = [p for p in optional_perms if p not in ['None', 'N/A']]

        return required_perms, optional_perms

    def count_examples(self, template_path: Path) -> int:
        """Count example references in template."""
        with open(template_path, 'r', encoding='utf-8') as f:
            content = f.read()

        # Count tffile references - these are validated during doc generation
        return len(re.findall(r'\{\{ tffile', content))

    def extract_ms_doc_links(self, template_path: Path) -> List[str]:
        """Extract Microsoft documentation links."""
        with open(template_path, 'r', encoding='utf-8') as f:
            content = f.read()

        # Find Microsoft Documentation section
        doc_section = re.search(
            r'## Microsoft Documentation\s*\n((?:- \[.*?\]\(.*?\)\s*\n)+)',
            content
        )

        if not doc_section:
            return []

        # Extract URLs
        return re.findall(r'\((https://[^\)]+)\)', doc_section.group(1))

    def get_resource_name_from_template(self, template_filename: str) -> str:
        """Convert template filename to Terraform resource name."""
        # Remove .md.tmpl extension
        name = template_filename.replace('.md.tmpl', '')

        # Add microsoft365_ prefix
        return f"microsoft365_{name}"

    def find_code_path(self, resource_name: str, component_type: str) -> Optional[Path]:
        """Find the code implementation path for a resource."""
        # Remove microsoft365_ prefix
        name_without_prefix = resource_name.replace('microsoft365_', '')

        # Determine search base
        search_bases = {
            'resource': self.paths.services / "resources",
            'data-source': self.paths.services / "datasources",
            'ephemeral': self.paths.services / "ephemerals",
            'action': self.paths.services / "actions",
            'list-resource': self.paths.services / "list-resources"
        }

        search_base = search_bases.get(component_type)
        if not search_base or not search_base.exists():
            return None

        # Strategy: Look for directories with resource.go, datasource.go, or list_resource.go
        # Different component types use different file names
        search_patterns = ['resource.go', 'datasource.go', 'list_resource.go']
        for pattern in search_patterns:
            for resource_file in search_base.rglob(pattern):
                resource_dir = resource_file.parent

                # Build a potential match by combining parent directory names
                rel_path = resource_dir.relative_to(search_base)
                path_parts = list(rel_path.parts)

                # Try different combinations to match
                # Remove api version directories (graph_beta, graph_v1.0, etc.)
                filtered_parts = [
                    p for p in path_parts
                    if p not in ['graph_beta', 'graph_v1.0', 'graph']
                ]

                # Create potential match strings
                potential_matches = [
                    '_'.join(filtered_parts),  # groups_group
                    '_'.join(path_parts),       # groups_graph_beta_group
                    path_parts[-1],             # group
                ]

                # Check if any match the end of our resource name
                for match_str in potential_matches:
                    if name_without_prefix.endswith(match_str):
                        return resource_dir

        return None

    def check_tests(self, code_path: Optional[Path], component_type: str) -> Tuple[bool, bool]:
        """Check if unit and acceptance tests exist."""
        if not code_path or not code_path.exists():
            return False, False

        # Different component types use different naming conventions
        if component_type == 'list-resource':
            has_unit = (
                (code_path / "list_test.go").exists() or
                (code_path / "list_resource_test.go").exists()
            )
            has_acceptance = (code_path / "list_acceptance_test.go").exists()
        elif component_type == 'data-source':
            has_unit = (code_path / "datasource_test.go").exists()
            has_acceptance = (code_path / "datasource_acceptance_test.go").exists()
        else:
            has_unit = (code_path / "resource_test.go").exists()
            has_acceptance = (code_path / "resource_acceptance_test.go").exists()

        return has_unit, has_acceptance

    def _create_component_metadata(
        self,
        template_file: Path,
        component_type: str
    ) -> ComponentMetadata:
        """Create metadata object from a template file."""
        frontmatter = self.parse_template_frontmatter(template_file)
        subcategory = frontmatter.get('subcategory', 'Unknown')

        initial_version, last_updated, status = self.extract_version_history(template_file)
        required_perms, optional_perms = self.extract_permissions(template_file)
        example_count = self.count_examples(template_file)
        ms_doc_links = self.extract_ms_doc_links(template_file)

        resource_name = self.get_resource_name_from_template(template_file.name)
        code_path = self.find_code_path(resource_name, component_type)
        has_unit, has_acceptance = self.check_tests(code_path, component_type)

        return ComponentMetadata(
            name=resource_name,
            component_type=component_type,
            subcategory=subcategory,
            initial_version=initial_version,
            last_updated_version=last_updated,
            status=status,
            required_permissions=required_perms,
            optional_permissions=optional_perms,
            example_count=example_count,
            has_unit_tests=has_unit,
            has_acceptance_tests=has_acceptance,
            ms_doc_links=ms_doc_links,
            template_path=str(template_file.relative_to(self.paths.repo_root)),
            code_path=str(code_path.relative_to(self.paths.repo_root)) if code_path else None
        )

    def process_templates(self, component_type: str, template_subdir: str):
        """Process all templates of a given type."""
        template_dir = self.paths.templates / template_subdir

        if not template_dir.exists():
            print(f"  ⚠️  Template directory not found: {template_subdir}")
            return

        template_count = 0
        for template_file in template_dir.glob("*.md.tmpl"):
            template_count += 1
            metadata = self._create_component_metadata(template_file, component_type)
            self.components.append(metadata)
            self.service_domains[metadata.subcategory].append(metadata)

        if template_count > 0:
            print(f"     Found {template_count} templates")

    def add_components_without_templates(self):
        """Add ephemerals and actions that don't have templates."""
        # Get set of existing component names to avoid duplicates
        existing_names = {comp.name for comp in self.components}

        # Process ephemerals from code (only if not already added from templates)
        ephemerals_file = self.paths.provider / "ephemeral_resources.go"
        if ephemerals_file.exists():
            with open(ephemerals_file, 'r', encoding='utf-8') as f:
                content = f.read()

            # Extract ephemeral registrations
            pattern = r'New(\w+)EphemeralResource'
            matches = re.findall(pattern, content)

            for func_name in matches:
                resource_name = re.sub(r'([A-Z])', r'_\1', func_name).lower().strip('_')
                full_name = f"microsoft365_graph_beta_{resource_name}"

                # Skip if already added from template
                if full_name in existing_names:
                    continue

                code_path = self.find_code_path(full_name, 'ephemeral')
                has_unit, has_acceptance = self.check_tests(code_path, 'ephemeral')

                code_path_str = (
                    str(code_path.relative_to(self.paths.repo_root))
                    if code_path else None
                )
                metadata = ComponentMetadata(
                    name=full_name,
                    component_type='ephemeral',
                    subcategory='Multitenant Management',
                    has_unit_tests=has_unit,
                    has_acceptance_tests=has_acceptance,
                    code_path=code_path_str
                )

                self.components.append(metadata)
                self.service_domains[metadata.subcategory].append(metadata)
                existing_names.add(full_name)

        # Process actions from code (only if not already added from templates)
        actions_file = self.paths.provider / "actions.go"
        if actions_file.exists():
            with open(actions_file, 'r', encoding='utf-8') as f:
                content = f.read()

            # Extract action registrations
            pattern = r'New(\w+)Action'
            matches = re.findall(pattern, content)

            for func_name in matches:
                resource_name = re.sub(r'([A-Z])', r'_\1', func_name).lower().strip('_')
                full_name = f"microsoft365_graph_beta_device_management_{resource_name}"

                # Skip if already added from template
                if full_name in existing_names:
                    continue

                code_path = self.find_code_path(full_name, 'action')

                code_path_str = (
                    str(code_path.relative_to(self.paths.repo_root))
                    if code_path else None
                )
                metadata = ComponentMetadata(
                    name=full_name,
                    component_type='action',
                    subcategory='Device Management',
                    code_path=code_path_str
                )

                self.components.append(metadata)
                self.service_domains[metadata.subcategory].append(metadata)
                existing_names.add(full_name)

    def count_provider_components(self) -> Dict[str, int]:
        """Count total components from provider registration files."""
        counts = {
            'resources': 0,
            'data_sources': 0,
            'list_resources': 0,
            'ephemerals': 0,
            'actions': 0
        }

        # Count resources
        resources_file = self.paths.provider / "resources.go"
        if resources_file.exists():
            with open(resources_file, 'r', encoding='utf-8') as f:
                content = f.read()
                matches = re.findall(r'New\w+Resource,', content)
                counts['resources'] = len(matches)

        # Count data sources
        datasources_file = self.paths.provider / "datasources.go"
        if datasources_file.exists():
            with open(datasources_file, 'r', encoding='utf-8') as f:
                content = f.read()
                matches = re.findall(r'New\w+DataSource,', content)
                counts['data_sources'] = len(matches)

        # Count list resources
        list_resources_file = self.paths.provider / "list_resources.go"
        if list_resources_file.exists():
            with open(list_resources_file, 'r', encoding='utf-8') as f:
                content = f.read()
                matches = re.findall(r'New\w+ListResource,', content)
                counts['list_resources'] = len(matches)

        # Count ephemerals
        ephemerals_file = self.paths.provider / "ephemeral_resources.go"
        if ephemerals_file.exists():
            with open(ephemerals_file, 'r', encoding='utf-8') as f:
                content = f.read()
                matches = re.findall(r'New\w+EphemeralResource,', content)
                counts['ephemerals'] = len(matches)

        # Count actions
        actions_file = self.paths.provider / "actions.go"
        if actions_file.exists():
            with open(actions_file, 'r', encoding='utf-8') as f:
                content = f.read()
                matches = re.findall(r'New\w+Action,', content)
                counts['actions'] = len(matches)

        return counts

    def _generate_header(self, lines: List[str]):
        """Generate the document header."""
        lines.append("# Provider Coverage")
        lines.append("")
        lines.append(
            "This provider offers extensive coverage across Microsoft 365 services "
            "including Intune, Microsoft 365, Teams, and Defender. Given the large size "
            "of this project, this page provides a comprehensive breakdown of available "
            "Terraform components organized by service domain."
        )
        lines.append("")
        timestamp = datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M UTC')
        lines.append(f"*Last updated: {timestamp}*")
        lines.append("")
        lines.append("---")
        lines.append("")

    def _generate_summary_table(self, lines: List[str]):
        """Generate the summary statistics table."""
        counts = self.count_provider_components()
        lines.append("## Summary Statistics")
        lines.append("")
        lines.append("| Terraform Block Type | Count |")
        lines.append("|---------------------|-------|")
        lines.append(f"| Resources | {counts['resources']} |")
        lines.append(f"| Data Sources | {counts['data_sources']} |")
        lines.append(f"| List Resources | {counts['list_resources']} |")
        lines.append(f"| Ephemerals | {counts['ephemerals']} |")
        lines.append(f"| Actions | {counts['actions']} |")
        total = sum(counts.values())
        lines.append(f"| **Total Components** | **{total}** |")
        lines.append("")
        lines.append("---")
        lines.append("")

    def generate_markdown(self) -> str:
        """Generate the provider_coverage.md content."""
        lines = []

        self._generate_header(lines)
        self._generate_summary_table(lines)

        # Sort service domains alphabetically
        sorted_domains = sorted(self.service_domains.items(), key=lambda x: x[0])

        # Generate service domain sections
        for domain_name, components in sorted_domains:
            self._generate_domain_section(lines, domain_name, components)

        return '\n'.join(lines)

    @staticmethod
    def _format_component_row(comp: ComponentMetadata) -> str:
        """Format a single component as a markdown table row."""
        version_intro = comp.initial_version or "—"
        version_last = comp.last_updated_version or "—"
        status = comp.status or "—"
        examples = str(comp.example_count) if comp.example_count > 0 else "—"
        unit = "✅" if comp.has_unit_tests else "❌"
        acceptance = "✅" if comp.has_acceptance_tests else "❌"

        return (
            f"| `{comp.name}` | {version_intro} | {version_last} | {status} | "
            f"{examples} | {unit} | {acceptance} |"
        )

    @staticmethod
    def _add_component_table(
        lines: List[str],
        title: str,
        components: List[ComponentMetadata],
        name_column: str
    ):
        """Add a collapsible component table section to the output."""
        count = len(components)
        lines.append("<details>")
        lines.append(f"<summary><b>{title} ({count})</b></summary>")
        lines.append("")
        lines.append(
            f"| {name_column} | Version Introduced | Last Updated | Status | "
            "Examples | Unit Tests | Acceptance Tests |"
        )
        lines.append(
            "|" + "-" * (len(name_column) + 2) + "|-------------------|--------------|"
            "--------|----------|------------|------------------|"
        )

        for comp in sorted(components, key=lambda x: x.name):
            lines.append(ProviderCoverageGenerator._format_component_row(comp))

        lines.append("")
        lines.append("</details>")
        lines.append("")

    @staticmethod
    def _pluralize(count: int, singular: str) -> str:
        """Return pluralized string based on count."""
        plural = 's' if count != 1 else ''
        return f"{count} {singular}{plural}"

    def _generate_domain_section(
        self,
        lines: List[str],
        domain_name: str,
        components: List[ComponentMetadata]
    ):
        """Generate a service domain section."""
        # Count by type
        resources = [c for c in components if c.component_type == 'resource']
        data_sources = [c for c in components if c.component_type == 'data-source']
        list_resources = [c for c in components if c.component_type == 'list-resource']
        ephemerals = [c for c in components if c.component_type == 'ephemeral']
        actions = [c for c in components if c.component_type == 'action']

        # Build header line
        header_parts = []
        if resources:
            header_parts.append(self._pluralize(len(resources), "Resource"))
        if data_sources:
            header_parts.append(self._pluralize(len(data_sources), "Data Source"))
        if list_resources:
            header_parts.append(self._pluralize(len(list_resources), "List Resource"))
        if ephemerals:
            header_parts.append(self._pluralize(len(ephemerals), "Ephemeral"))
        if actions:
            header_parts.append(self._pluralize(len(actions), "Action"))

        lines.append(f"## {domain_name}")
        lines.append("")
        lines.append(f"**{' • '.join(header_parts)}**")
        lines.append("")

        # Add component tables
        if resources:
            self._add_component_table(lines, "Resources", resources, "Resource Name")

        if data_sources:
            self._add_component_table(lines, "Data Sources", data_sources, "Data Source Name")

        if list_resources:
            self._add_component_table(
                lines, "List Resources", list_resources, "List Resource Name"
            )

        if ephemerals:
            self._add_component_table(lines, "Ephemerals", ephemerals, "Ephemeral Name")

        # Actions section (summary only, no table)
        if actions:
            lines.append("<details>")
            lines.append(f"<summary><b>Actions ({len(actions)})</b></summary>")
            lines.append("")
            lines.append(
                "Device management actions for managed devices including lifecycle "
                "operations, security actions, and maintenance tasks."
            )
            lines.append("")
            lines.append("</details>")
            lines.append("")

        lines.append("---")
        lines.append("")

    def run(self, output_dir: str = "docs/development"):
        """Main execution method."""
        print("🔍 Scanning templates and code...")

        # Process resources
        print("  📦 Processing resource templates...")
        self.process_templates('resource', 'resources')

        # Process data sources
        print("  📖 Processing data source templates...")
        self.process_templates('data-source', 'data-sources')

        # Process list resources
        print("  📋 Processing list resource templates...")
        self.process_templates('list-resource', 'list-resources')

        # Process ephemerals
        print("  ⚡ Processing ephemeral templates...")
        self.process_templates('ephemeral', 'ephemerals')

        # Process actions
        print("  🎬 Processing action templates...")
        self.process_templates('action', 'actions')

        # Add components without templates
        print("  🔧 Adding components without templates...")
        self.add_components_without_templates()

        print(
            f"✅ Found {len(self.components)} components across "
            f"{len(self.service_domains)} service domains"
        )

        # Generate markdown
        print("📝 Generating provider_coverage.md...")
        markdown = self.generate_markdown()

        # Write to file
        output_path = self.paths.repo_root / output_dir / "provider_coverage.md"
        output_path.parent.mkdir(parents=True, exist_ok=True)

        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(markdown)

        print(f"✅ Generated: {output_path}")
        print(f"📊 Total components: {len(self.components)}")


def main():
    """Main entry point."""
    # Determine repository root
    script_dir = Path(__file__).parent
    repo_root = script_dir.parent.parent

    # Verify we're in the right place
    if not (repo_root / "internal" / "provider").exists():
        print("❌ Error: Could not find internal/provider directory")
        print(f"   Looking in: {repo_root}")
        return 1

    # Run generator
    generator = ProviderCoverageGenerator(repo_root)
    generator.run(output_dir="docs/guides")

    return 0


if __name__ == "__main__":
    sys.exit(main())
