"""Change data models for OpenAPI schema detection."""

from dataclasses import dataclass, field
from typing import Optional, List, Any


@dataclass
class PropertyChange:
    """Represents a property change in an OpenAPI schema."""
    property_name: str
    change_type: str  # 'added', 'removed', 'type_changed', 'required_changed', 'nullable_changed'
    old_type: Optional[str] = None
    new_type: Optional[str] = None
    old_required: bool = False
    new_required: bool = False
    old_nullable: Optional[bool] = None
    new_nullable: Optional[bool] = None
    
    # Rich metadata from OpenAPI
    description: Optional[str] = None
    old_description: Optional[str] = None
    new_description: Optional[str] = None
    enum_values: Optional[List[Any]] = None
    old_enum_values: Optional[List[Any]] = None
    new_enum_values: Optional[List[Any]] = None
    format: Optional[str] = None  # e.g., 'date-time', 'email', 'uuid'
    pattern: Optional[str] = None  # regex pattern
    min_length: Optional[int] = None
    max_length: Optional[int] = None
    minimum: Optional[float] = None
    maximum: Optional[float] = None
    default: Optional[Any] = None
    example: Optional[Any] = None
    deprecated: bool = False
    read_only: bool = False
    write_only: bool = False
    
    @property
    def is_breaking(self) -> bool:
        """Check if this is a breaking change."""
        return self.change_type in ['removed', 'type_changed'] or (
            self.change_type == 'required_changed' and self.new_required and not self.old_required
        )


@dataclass
class SchemaChange:
    """Represents changes to an OpenAPI schema (model)."""
    schema_name: str  # e.g., 'microsoft.graph.user'
    file_path: str  # For compatibility with existing filter
    added_properties: List[PropertyChange] = field(default_factory=list)
    removed_properties: List[PropertyChange] = field(default_factory=list)
    type_changed_properties: List[PropertyChange] = field(default_factory=list)
    required_changed_properties: List[PropertyChange] = field(default_factory=list)
    nullable_changed_properties: List[PropertyChange] = field(default_factory=list)

    @property
    def has_changes(self) -> bool:
        """Check if this schema has any changes."""
        return bool(
            self.added_properties or
            self.removed_properties or
            self.type_changed_properties or
            self.required_changed_properties or
            self.nullable_changed_properties
        )

    @property
    def has_breaking_changes(self) -> bool:
        """Check if this schema has breaking changes."""
        return bool(
            self.removed_properties or
            self.type_changed_properties or
            [p for p in self.required_changed_properties if p.is_breaking]
        )

    @property
    def change_summary(self) -> str:
        """Get a summary of changes."""
        parts = []
        if self.added_properties:
            parts.append(f"+{len(self.added_properties)} properties")
        if self.removed_properties:
            parts.append(f"-{len(self.removed_properties)} properties")
        if self.type_changed_properties:
            parts.append(f"~{len(self.type_changed_properties)} type changes")
        if self.required_changed_properties:
            parts.append(f"!{len(self.required_changed_properties)} required changes")
        if self.nullable_changed_properties:
            parts.append(f"?{len(self.nullable_changed_properties)} nullable changes")
        return ", ".join(parts)
    
    @property
    def model_name(self) -> str:
        """Alias for schema_name for compatibility with provider filter.
        
        OpenAPI schema names are like: microsoft.graph.cloudPcDeviceImage
        Provider usage has Go model names like: CloudPcDeviceImage
        
        We need to match them by just taking the last part and capitalizing first letter.
        """
        parts = self.schema_name.split('.')
        if len(parts) >= 3:
            # Take last part: cloudPcDeviceImage
            name = parts[-1]
            
            # Handle underscore-separated names: cloud_pc_device_image → CloudPcDeviceImage
            if '_' in name:
                return ''.join(word.capitalize() for word in name.split('_'))
            
            # Handle camelCase names: cloudPcDeviceImage → CloudPcDeviceImage
            # Just capitalize the first letter, preserve rest
            return name[0].upper() + name[1:] if name else name
        
        return self.schema_name
