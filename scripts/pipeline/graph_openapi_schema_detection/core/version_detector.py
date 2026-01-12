"""Version detection from OpenAPI specifications."""

import re
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ..models import VersionResult
    from .progress_reporter import ProgressReporter


class VersionDetector:
    """Detects version information from OpenAPI spec."""
    
    # Pattern to extract version from info section
    VERSION_PATTERN = re.compile(r'^\s*version:\s*["\']?([^"\']+)["\']?', re.MULTILINE)
    
    def __init__(self, reporter: 'ProgressReporter'):
        """Initialize version detector.
        
        Args:
            reporter: Progress reporter
        """
        self.reporter = reporter
    
    def extract_version_from_spec(self, spec_content: str) -> 'VersionResult':
        """Extract version from 'info.version' in OpenAPI spec.
        
        Args:
            spec_content: Full OpenAPI spec content
            
        Returns:
            VersionResult with extracted version
        """
        from models import VersionResult  # type: ignore
        
        # Look for version in info section (should be near the top)
        # OpenAPI format:
        # info:
        #   version: beta
        #   title: ...
        
        lines = spec_content.splitlines()[:50]  # Version should be in first 50 lines
        
        for i, line in enumerate(lines):
            if 'info:' in line:
                # Found info section, look for version in next few lines
                for j in range(i + 1, min(i + 10, len(lines))):
                    version_line = lines[j]
                    match = self.VERSION_PATTERN.match(version_line)
                    if match:
                        version = match.group(1).strip()
                        self.reporter.info(f"   Detected version: {version}")
                        return VersionResult.success(version)
        
        return VersionResult.not_found("Version not found in OpenAPI spec")
    
    def compare_versions(self, old_version: str, new_version: str) -> bool:
        """Check if versions are different.
        
        Args:
            old_version: Previous version
            new_version: New version
            
        Returns:
            True if versions differ
        """
        return old_version != new_version
