"""Result types for schema detection operations."""

from dataclasses import dataclass
from enum import Enum
from typing import Optional, Tuple, List, Dict, TYPE_CHECKING

if TYPE_CHECKING:
    from models.changes import ModelChange
    from models.statistics import ParseStatistics


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
    version_tuple: Optional[Tuple[int, int, int]] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, version: str, version_tuple: Tuple[int, int, int]) -> 'VersionResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, version, version_tuple)
    
    @classmethod
    def error(cls, message: str) -> 'VersionResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def not_found(cls, message: str = "Version not found") -> 'VersionResult':
        """Create a not found result."""
        return cls(ResultStatus.NOT_FOUND, error_message=message)


@dataclass
class VersionChangeResult:
    """Result of PR version change detection."""
    status: ResultStatus
    old_version: Optional[str] = None
    new_version: Optional[str] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, old_version: str, new_version: str) -> 'VersionChangeResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, old_version, new_version)
    
    @classmethod
    def error(cls, message: str) -> 'VersionChangeResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def not_found(cls) -> 'VersionChangeResult':
        """Create a not found result."""
        return cls(ResultStatus.NOT_FOUND, error_message="Version change not detected in PR")


@dataclass
class IssueCreationResult:
    """Result of GitHub issue creation."""
    status: ResultStatus
    issue_number: Optional[str] = None
    issue_url: Optional[str] = None
    error_message: Optional[str] = None
    
    @property
    def is_success(self) -> bool:
        """Check if operation was successful."""
        return self.status == ResultStatus.SUCCESS
    
    @classmethod
    def success(cls, issue_number: str, issue_url: str) -> 'IssueCreationResult':
        """Create a successful result."""
        return cls(ResultStatus.SUCCESS, issue_number, issue_url)
    
    @classmethod
    def error(cls, message: str) -> 'IssueCreationResult':
        """Create an error result."""
        return cls(ResultStatus.ERROR, error_message=message)
    
    @classmethod
    def dry_run(cls) -> 'IssueCreationResult':
        """Create a dry run result."""
        return cls(ResultStatus.SUCCESS, issue_number="DRY_RUN")


@dataclass
class ValidationResult:
    """Result of version validation."""
    is_valid: bool
    message: Optional[str] = None
    
    @classmethod
    def valid(cls) -> 'ValidationResult':
        """Create a valid result."""
        return cls(True)
    
    @classmethod
    def invalid(cls, message: str) -> 'ValidationResult':
        """Create an invalid result."""
        return cls(False, message)


@dataclass
class DetectionResult:
    """Complete schema change detection result."""
    current_version: str
    new_version: str
    timestamp: str
    total_models_changed: int
    models_with_changes: int
    filtered_models: int
    model_changes: List['ModelChange']
    changelog_section: str
    statistics: 'ParseStatistics'
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            "current_version": self.current_version,
            "new_version": self.new_version,
            "timestamp": self.timestamp,
            "total_models_changed": self.total_models_changed,
            "models_with_changes": self.models_with_changes,
            "filtered_models": self.filtered_models,
            "model_changes": [
                {
                    "file_path": mc.file_path,
                    "model_name": mc.model_name,
                    "added_fields": [{"name": f.name, "field_type": f.field_type} for f in mc.added_fields],
                    "removed_fields": [{"name": f.name, "field_type": f.field_type} for f in mc.removed_fields],
                    "added_methods": [{"signature": m.signature} for m in mc.added_methods],
                    "removed_methods": [{"signature": m.signature} for m in mc.removed_methods],
                    "added_embedded_types": [{"type_name": e.type_name} for e in mc.added_embedded_types],
                    "removed_embedded_types": [{"type_name": e.type_name} for e in mc.removed_embedded_types],
                }
                for mc in self.model_changes
            ],
            "changelog_section": self.changelog_section,
            "statistics": {
                "total_files_in_diff": self.statistics.total_files_in_diff,
                "files_with_changes": self.statistics.files_with_changes,
                "files_without_changes": self.statistics.files_without_changes,
                "total_lines_processed": self.statistics.total_lines_processed,
                "added_lines_processed": self.statistics.added_lines_processed,
                "removed_lines_processed": self.statistics.removed_lines_processed,
                "struct_fields_added": self.statistics.struct_fields_added,
                "struct_fields_removed": self.statistics.struct_fields_removed,
                "interface_methods_added": self.statistics.interface_methods_added,
                "interface_methods_removed": self.statistics.interface_methods_removed,
                "embedded_types_added": self.statistics.embedded_types_added,
                "embedded_types_removed": self.statistics.embedded_types_removed,
            }
        }
