"""Core components for schema detection."""

from .progress_reporter import ProgressReporter
from .version_parser import VersionParser
from .github_client import GitHubClient
from .diff_fetcher import DiffFetcher
from .struct_parser import StructParser
from .issue_builder import IssueBuilder

__all__ = [
    'ProgressReporter',
    'VersionParser',
    'GitHubClient',
    'DiffFetcher',
    'StructParser',
    'IssueBuilder',
]
