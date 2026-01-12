"""Version parsing and validation."""

from typing import Optional, TYPE_CHECKING

from regex_patterns import RegexPatterns  # type: ignore
from models import VersionResult, ValidationResult  # type: ignore

if TYPE_CHECKING:
    from core.progress_reporter import ProgressReporter


class VersionParser:
    """Handles version parsing and validation."""

    def parse_version(self, version_str: str) -> VersionResult:
        """Parse version string into structured result.
        
        Args:
            version_str: Version string (e.g., 'v0.156.0')
            
        Returns:
            VersionResult with parsed version data
        """
        match = RegexPatterns.VERSION_FULL.match(version_str)
        if match:
            version_tuple = tuple(map(int, match.groups()))
            return VersionResult.success(version_str, version_tuple)
        return VersionResult.error(f"Invalid version format: {version_str}")

    def extract_version_from_line(self, line: str) -> Optional[str]:
        """Extract version string from a line of text.
        
        Args:
            line: Line of text that may contain a version
            
        Returns:
            Version string or None if not found
        """
        match = RegexPatterns.VERSION_IN_TEXT.search(line)
        return match.group(0) if match else None

    def validate_increment(self, old_version: str, new_version: str, 
                          reporter: 'ProgressReporter') -> ValidationResult:
        """Validate that version increment is acceptable.
        
        Args:
            old_version: Old version string
            new_version: New version string
            reporter: Reporter for warnings
            
        Returns:
            ValidationResult with validation status
        """
        old_result = self.parse_version(old_version)
        new_result = self.parse_version(new_version)

        if not old_result.is_success or not new_result.is_success:
            reporter.error("Invalid version format")
            return ValidationResult.invalid("Invalid version format")

        old_major, old_minor, old_patch = old_result.version_tuple
        new_major, new_minor, new_patch = new_result.version_tuple

        if new_major == old_major:
            if new_minor == old_minor + 1 and new_patch == 0:
                return ValidationResult.valid()  # Valid minor version bump
            elif new_minor == old_minor and new_patch == old_patch + 1:
                return ValidationResult.valid()  # Valid patch version bump
            elif new_minor == old_minor and new_patch > old_patch:
                reporter.warning("Multiple patch version increment detected")
                return ValidationResult.valid()  # Multiple patch bump - acceptable
            elif new_minor > old_minor + 1:
                msg = f"Multiple minor version jump detected: {old_version} -> {new_version}"
                reporter.warning(msg)
                reporter.warning("This may indicate missing intermediate versions.")
                return ValidationResult.invalid(msg)

        msg = f"Unexpected version change: {old_version} -> {new_version}"
        reporter.warning(msg)
        return ValidationResult.invalid(msg)
