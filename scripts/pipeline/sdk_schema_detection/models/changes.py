"""Change data models for schema detection."""

from dataclasses import dataclass, field
from typing import Optional, List


@dataclass
class FieldChange:
    """Represents a field change in a Go struct."""
    field_name: str
    field_type: str
    change_type: str  # 'added' or 'removed'
    line_number: Optional[int] = None
    
    @property
    def name(self) -> str:
        """Alias for field_name for compatibility."""
        return self.field_name


@dataclass
class MethodChange:
    """Represents an interface method change."""
    method_name: str
    parameters: str
    return_type: str
    change_type: str  # 'added' or 'removed'
    line_number: Optional[int] = None
    
    @property
    def signature(self) -> str:
        """Get the full method signature."""
        params = f"({self.parameters})" if self.parameters else "()"
        returns = f" {self.return_type}" if self.return_type else ""
        return f"{self.method_name}{params}{returns}"


@dataclass
class EmbeddedTypeChange:
    """Represents a change in embedded types (interfaces or structs)."""
    type_name: str
    change_type: str  # 'added' or 'removed'
    context: str  # 'interface' or 'struct'
    line_number: Optional[int] = None


@dataclass
class ModelChange:
    """Represents changes to a Go model file."""
    file_path: str
    model_name: str
    added_fields: List[FieldChange] = field(default_factory=list)
    removed_fields: List[FieldChange] = field(default_factory=list)
    added_methods: List[MethodChange] = field(default_factory=list)
    removed_methods: List[MethodChange] = field(default_factory=list)
    added_embedded_types: List[EmbeddedTypeChange] = field(default_factory=list)
    removed_embedded_types: List[EmbeddedTypeChange] = field(default_factory=list)

    @property
    def has_changes(self) -> bool:
        """Check if this model has any changes."""
        return bool(
            self.added_fields or self.removed_fields or
            self.added_methods or self.removed_methods or
            self.added_embedded_types or self.removed_embedded_types
        )

    @property
    def change_summary(self) -> str:
        """Get a summary of changes."""
        parts = []
        if self.added_fields:
            parts.append(f"+{len(self.added_fields)} fields")
        if self.removed_fields:
            parts.append(f"-{len(self.removed_fields)} fields")
        if self.added_methods:
            parts.append(f"+{len(self.added_methods)} methods")
        if self.removed_methods:
            parts.append(f"-{len(self.removed_methods)} methods")
        if self.added_embedded_types:
            parts.append(f"+{len(self.added_embedded_types)} embedded")
        if self.removed_embedded_types:
            parts.append(f"-{len(self.removed_embedded_types)} embedded")
        return ", ".join(parts)
