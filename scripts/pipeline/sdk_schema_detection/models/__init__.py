"""Data models for schema detection."""

from .results import (
    ResultStatus,
    VersionResult,
    VersionChangeResult,
    IssueCreationResult,
    ValidationResult,
    DetectionResult,
)
from .changes import (
    FieldChange,
    MethodChange,
    EmbeddedTypeChange,
    ModelChange,
)
from .statistics import ParseStatistics

__all__ = [
    'ResultStatus',
    'VersionResult',
    'VersionChangeResult',
    'IssueCreationResult',
    'ValidationResult',
    'DetectionResult',
    'FieldChange',
    'MethodChange',
    'EmbeddedTypeChange',
    'ModelChange',
    'ParseStatistics',
]
