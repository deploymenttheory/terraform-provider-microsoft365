"""Result types for OpenAPI schema detection operations."""

from dataclasses import dataclass
from enum import Enum
from typing import Optional, List, Dict, TYPE_CHECKING

if TYPE_CHECKING:
    from .changes import SchemaChange
    from .statistics import ParseStatistics


class ResultStatus(Enum):
    """Status of an operation result."""
    SUCCESS = "success"
    ERROR = "error"
    NOT_FOUND = "not_found"
    INVALID = "invalid"


@dataclass
class VersionResult:
    """Result of version parsing or extraction."""
    status: ResultStatus
    version: Optional[str] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, version: str) -> 'VersionResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, version)
    
    @classmethod
    def error(cls, message: str) -> 'VersionResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def not_found(cls, message: str = "Version not found") -> 'VersionResult':
        """Create a not found result."""
        return cls(ResultStatus.NOT_FOUND, error_message=message)


@dataclass
class DetectionResult:
    """Complete OpenAPI schema change detection result."""
    spec_version: str
    previous_version: str
    detection_timestamp: str
    total_schemas_changed: int
    schemas_with_changes: int
    filtered_schemas: int
    breaking_changes_count: int
    schema_changes: List['SchemaChange']
    statistics: 'ParseStatistics'
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            "spec_version": self.spec_version,
            "previous_version": self.previous_version,
            "detection_timestamp": self.detection_timestamp,
            "total_schemas_changed": self.total_schemas_changed,
            "schemas_with_changes": self.schemas_with_changes,
            "filtered_schemas": self.filtered_schemas,
            "breaking_changes_count": self.breaking_changes_count,
            "schema_changes": [
                {
                    "schema_name": sc.schema_name,
                    "model_name": sc.model_name,
                    "file_path": sc.file_path,
                    "has_breaking_changes": sc.has_breaking_changes,
                    "change_summary": sc.change_summary,
                    "added_properties": [
                        {
                            "name": p.property_name,
                            "type": p.new_type,
                            "required": p.new_required,
                            "nullable": p.new_nullable,
                            "description": p.description,
                            "enum": p.enum_values,
                            "format": p.format,
                            "pattern": p.pattern,
                            "minLength": p.min_length,
                            "maxLength": p.max_length,
                            "minimum": p.minimum,
                            "maximum": p.maximum,
                            "default": p.default,
                            "example": p.example,
                            "deprecated": p.deprecated,
                            "readOnly": p.read_only,
                            "writeOnly": p.write_only,
                        }
                        for p in sc.added_properties
                    ],
                    "removed_properties": [
                        {"name": p.property_name, "type": p.old_type, "required": p.old_required}
                        for p in sc.removed_properties
                    ],
                    "type_changed_properties": [
                        {"name": p.property_name, "old_type": p.old_type, "new_type": p.new_type}
                        for p in sc.type_changed_properties
                    ],
                    "required_changed_properties": [
                        {"name": p.property_name, "old_required": p.old_required, "new_required": p.new_required}
                        for p in sc.required_changed_properties
                    ],
                }
                for sc in self.schema_changes
            ],
            "statistics": {
                "total_schemas_compared": self.statistics.total_schemas_compared,
                "schemas_with_changes": self.statistics.schemas_with_changes,
                "schemas_added": self.statistics.schemas_added,
                "schemas_removed": self.statistics.schemas_removed,
                "properties_added": self.statistics.properties_added,
                "properties_removed": self.statistics.properties_removed,
                "type_changes": self.statistics.type_changes,
                "required_changes": self.statistics.required_changes,
            }
        }
