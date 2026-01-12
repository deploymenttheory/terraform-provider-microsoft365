"""Statistics models for OpenAPI schema detection."""

from dataclasses import dataclass


@dataclass
class ParseStatistics:
    """Statistics from parsing and comparing OpenAPI schemas."""
    total_schemas_compared: int = 0
    schemas_with_changes: int = 0
    schemas_added: int = 0
    schemas_removed: int = 0
    
    # Property-level changes
    properties_added: int = 0
    properties_removed: int = 0
    type_changes: int = 0
    required_changes: int = 0
    nullable_changes: int = 0
    
    def get_summary(self) -> str:
        """Get a human-readable summary."""
        parts = [
            "  📊 Comparison Statistics:",
            f"     Schemas compared: {self.total_schemas_compared:,}",
            f"     Schemas with changes: {self.schemas_with_changes}",
            f"     Schemas added: {self.schemas_added}",
            f"     Schemas removed: {self.schemas_removed}",
            "",
            "  🔧 Property Changes:",
            f"     Properties added: {self.properties_added}",
            f"     Properties removed: {self.properties_removed}",
            f"     Type changes: {self.type_changes}",
            f"     Required changes: {self.required_changes}",
            f"     Nullable changes: {self.nullable_changes}",
        ]
        return "\n".join(parts)
