"""GitHub issue content builder."""

from datetime import datetime
from typing import List, TYPE_CHECKING

if TYPE_CHECKING:
    from models import ModelChange


class IssueBuilder:
    """Builds GitHub issue content from analysis results."""

    def __init__(self, sdk_repo: str):
        """Initialize issue builder.
        
        Args:
            sdk_repo: SDK repository (owner/repo)
        """
        self.sdk_repo = sdk_repo

    def build_title(self, old_version: str, new_version: str) -> str:
        """Build issue title.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            
        Returns:
            Issue title
        """
        return f"Schema Update Required: Microsoft Graph SDK {old_version} → {new_version}"

    def build_body(self, old_version: str, new_version: str, 
                  model_changes: List['ModelChange'], changelog_section: str) -> str:
        """Build complete issue body.
        
        Args:
            old_version: Old SDK version
            new_version: New SDK version
            model_changes: List of model changes
            changelog_section: Changelog text
            
        Returns:
            Issue body markdown
        """
        sections = [
            self._build_summary(old_version, new_version, model_changes),
            self._build_changed_models(model_changes),
            self._build_action_required(),
            self._build_references(old_version, new_version),
            self._build_changelog(changelog_section),
            self._build_footer()
        ]
        
        return '\n'.join(sections)

    def _build_summary(self, old_version: str, new_version: str, 
                      model_changes: List['ModelChange']) -> str:
        """Build summary section."""
        return f"""## Summary
The Microsoft Graph Beta SDK has been updated from `{old_version}` to `{new_version}`.
This update includes {len(model_changes)} model(s) with schema changes that require review and potential Terraform schema updates.
"""

    def _build_changed_models(self, model_changes: List['ModelChange']) -> str:
        """Build changed models section."""
        parts = ["## Changed Models", ""]
        
        for change in model_changes:
            parts.append(f"### `{change.model_name}` ({change.file_path})")
            parts.append(f"**Changes:** {change.change_summary}")
            parts.append("")

            # Struct fields
            if change.added_fields:
                parts.append("**Added Struct Fields:**")
                for fld in change.added_fields:
                    parts.append(f"- `{fld.field_name}` ({fld.field_type})")
                parts.append("")

            if change.removed_fields:
                parts.append("**Removed Struct Fields:**")
                for fld in change.removed_fields:
                    parts.append(f"- `{fld.field_name}` ({fld.field_type})")
                parts.append("")

            # Interface methods
            if change.added_methods:
                parts.append("**Added Interface Methods:**")
                for method in change.added_methods:
                    parts.append(f"- `{method.signature}`")
                parts.append("")

            if change.removed_methods:
                parts.append("**Removed Interface Methods:**")
                for method in change.removed_methods:
                    parts.append(f"- `{method.signature}`")
                parts.append("")

            # Embedded types
            if change.added_embedded_types:
                parts.append("**Added Embedded Types:**")
                for emb in change.added_embedded_types:
                    parts.append(f"- `{emb.type_name}` ({emb.context})")
                parts.append("")

            if change.removed_embedded_types:
                parts.append("**Removed Embedded Types:**")
                for emb in change.removed_embedded_types:
                    parts.append(f"- `{emb.type_name}` ({emb.context})")
                parts.append("")

        return '\n'.join(parts)

    def _build_action_required(self) -> str:
        """Build action required section."""
        return """## Action Required

1. Review each changed model listed above
2. For struct field changes:
   - Update corresponding Terraform resource schemas
   - Add/update field mappings in CRUD operations
3. For interface method changes:
   - Review API contract changes
   - Update method calls if signatures changed
   - Verify compatibility with existing code
4. For embedded type changes:
   - Review inheritance/composition changes
   - Check for breaking changes in type hierarchy
5. Add/update tests for all changes
6. Update documentation

⚠️ **Interface method changes may indicate breaking API changes!**
"""

    def _build_references(self, old_version: str, new_version: str) -> str:
        """Build references section."""
        return f"""## References

- [SDK Changelog](https://github.com/{self.sdk_repo}/blob/main/CHANGELOG.md)
- [Version Diff](https://github.com/{self.sdk_repo}/compare/{old_version}...{new_version})
- [Models Diff](https://github.com/{self.sdk_repo}/compare/{old_version}...{new_version}#files_bucket)
"""

    def _build_changelog(self, changelog_section: str) -> str:
        """Build changelog section."""
        if not changelog_section or "not found" in changelog_section.lower():
            return ""

        parts = ["## Changelog Excerpt", "", "```"]
        parts.append(changelog_section[:1000])
        if len(changelog_section) > 1000:
            parts.append("... (truncated)")
        parts.append("```")
        parts.append("")
        
        return '\n'.join(parts)

    def _build_footer(self) -> str:
        """Build footer."""
        timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        return f"---\n🤖 Auto-generated by kiota_graph_sdk_schema_change_detector.py on {timestamp}"
