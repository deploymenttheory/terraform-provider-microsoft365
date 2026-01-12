"""Data models for OpenAPI schema detection."""

from .results import (
    ResultStatus,
    VersionResult,
    DetectionResult,
)
from .changes import (
    PropertyChange,
    SchemaChange,
)
from .statistics import ParseStatistics

__all__ = [
    'ResultStatus',
    'VersionResult',
    'DetectionResult',
    'PropertyChange',
    'SchemaChange',
    'ParseStatistics',
]
