#!/usr/bin/env python3
"""SDK diff analysis and impact assessment.

Analyzes SDK changes and determines their impact on the Terraform provider.
"""

from typing import Dict, List, Any
import re

from .github_api import compare_versions


class ImpactLevel:
    """Impact severity levels."""
    CRITICAL = "critical"      # Breaking changes to used APIs
    WARNING = "warning"        # Deprecations, signature changes
    SAFE = "safe"             # New features, unused changes
    OPPORTUNITY = "opportunity"  # New fields in types we use - potential additions
    NOISE = "noise"           # Filtered out (tests, docs, unused packages)


class ChangeAnalyzer:
    """Analyzes SDK changes and assesses impact."""
    
    def __init__(self, usage_data: Dict[str, Any]):
        """Initialize analyzer with provider's SDK usage.
        
        Args:
            usage_data: Usage data from go_parser.extract_sdk_usage()
        """
        self.usage_data = usage_data
        self.used_packages = set(usage_data.get("packages", {}).keys())
        self.used_types = set(usage_data.get("types", {}).keys())
        self.used_methods = set(usage_data.get("methods", {}).keys())
    
    def analyze_file_changes(self, files: List[Dict[str, Any]]) -> Dict[str, List[Dict]]:
        """Analyze changed files and categorize by impact.
        
        Args:
            files: List of changed files from github_api.compare_versions()
            
        Returns:
            Dictionary categorizing changes:
            {
                "critical": [...],
                "warning": [...],
                "safe": [...],
                "noise": [...]
            }
        """
        categorized = {
            ImpactLevel.CRITICAL: [],
            ImpactLevel.WARNING: [],
            ImpactLevel.SAFE: [],
            ImpactLevel.OPPORTUNITY: [],
            ImpactLevel.NOISE: []
        }
        
        for file_change in files:
            filename = file_change.get("filename", "")
            status = file_change.get("status", "")
            
            # Determine impact level
            impact = self._assess_file_impact(filename, status)
            categorized[impact].append({
                "file": filename,
                "status": status,
                "additions": file_change.get("additions", 0),
                "deletions": file_change.get("deletions", 0),
                "reason": self._explain_impact(filename, status, impact)
            })
        
        return categorized
    
    def _assess_file_impact(self, filename: str, status: str) -> str:
        """Assess the impact level of a file change.
        
        Args:
            filename: Path to changed file
            status: Change type (added, modified, removed)
            
        Returns:
            Impact level (critical, warning, safe, noise)
        """
        # Filter noise
        if self._is_noise_file(filename):
            return ImpactLevel.NOISE
        
        # Check if file is in a package we use
        package_path = self._extract_package_path(filename)
        
        # If not used, it's noise
        if not self._is_used_package(package_path):
            return ImpactLevel.NOISE
        
        # File removed = critical
        if status == "removed":
            return ImpactLevel.CRITICAL
        
        # Check for breaking patterns in modified files
        if status == "modified" and self._likely_breaking(filename):
            return ImpactLevel.WARNING
        
        # New files in used packages = safe (new features)
        if status == "added":
            return ImpactLevel.SAFE
        
        # Default: modifications to used packages = warning
        return ImpactLevel.WARNING
    
    def _is_noise_file(self, filename: str) -> bool:
        """Check if file should be filtered as noise.
        
        Args:
            filename: Path to file
            
        Returns:
            True if file is noise (tests, docs, etc.)
        """
        noise_patterns = [
            "_test.go",           # Test files
            "/testdata/",         # Test data
            "/examples/",         # Example code
            "README.md",          # Documentation
            "CHANGELOG.md",       # Changelog
            ".github/",           # GitHub workflows
            "/internal/",         # Internal packages
        ]
        
        return any(pattern in filename for pattern in noise_patterns)
    
    def _extract_package_path(self, filename: str) -> str:
        """Extract Go package path from filename.
        
        Args:
            filename: Path like "models/user.go" or "users/users.go"
            
        Returns:
            Package path (directory containing the file)
        """
        parts = filename.split('/')
        if len(parts) > 1:
            return '/'.join(parts[:-1])
        return ""
    
    def _is_used_package(self, package_path: str) -> bool:
        """Check if a package path is used in the provider.
        
        Args:
            package_path: Package directory path
            
        Returns:
            True if package is imported anywhere
        """
        # Check if any used package contains this path
        for used_pkg in self.used_packages:
            if package_path in used_pkg:
                return True
        return False
    
    def _likely_breaking(self, filename: str) -> bool:
        """Heuristic: is this modification likely breaking?
        
        Args:
            filename: Path to modified file
            
        Returns:
            True if likely to be breaking
        """
        # Model/type files are more likely to have breaking changes
        if "models/" in filename or "model.go" in filename:
            return True
        
        # Client interface changes are likely breaking
        if "client.go" in filename or "client_" in filename:
            return True
        
        return False
    
    def _explain_impact(self, filename: str, status: str, impact: str) -> str:
        """Generate human-readable explanation of impact.
        
        Args:
            filename: Changed file
            status: Change type
            impact: Assessed impact level
            
        Returns:
            Explanation string
        """
        if impact == ImpactLevel.NOISE:
            return "Filtered out (test/doc/unused package)"
        
        if impact == ImpactLevel.CRITICAL:
            if status == "removed":
                return "File removed in used package"
            return "Breaking change detected"
        
        if impact == ImpactLevel.WARNING:
            return f"Modified in used package ({self._extract_package_path(filename)})"
        
        if impact == ImpactLevel.SAFE:
            return "New feature in used package"
        
        if impact == ImpactLevel.OPPORTUNITY:
            return "New field added to type you use"
        
        return "Unknown impact"
    
    def analyze_field_additions(self, repo: str, base_version: str, head_version: str) -> List[Dict[str, Any]]:
        """Analyze field additions in types we already use.
        
        This detects new fields added to SDK types that the provider is already using,
        which represent opportunities to adopt new functionality.
        
        Args:
            repo: Repository name (e.g., "microsoftgraph/msgraph-beta-sdk-go")
            base_version: Base version tag
            head_version: Head version tag
            
        Returns:
            List of field addition opportunities:
            [
                {
                    "type": "models.User",
                    "field": "PreferredLanguage",
                    "file": "models/user.go",
                    "line_added": 145
                }
            ]
        """
        opportunities = []
        
        # Get file changes
        comparison = compare_versions(repo, base_version, head_version)
        
        # Track which types we use
        used_types_simple = set()
        for full_type in self.used_types:
            # Extract simple type name from full path
            # e.g., "github.com/.../models.User" -> "User"
            if '.' in full_type:
                simple_name = full_type.split('.')[-1]
                used_types_simple.add(simple_name)
        
        # Also check fields dict for types
        for full_type in self.usage_data.get('fields', {}).keys():
            if '.' in full_type:
                simple_name = full_type.split('.')[-1]
                used_types_simple.add(simple_name)
        
        # Analyze each changed file in models directory
        for file_change in comparison['files']:
            filename = file_change.get('filename', '')
            
            # Only look at model files
            if 'models/' not in filename or not filename.endswith('.go'):
                continue
            
            # Skip test files
            if '_test.go' in filename:
                continue
            
            # Look for struct field additions in the patch
            patch = file_change.get('patch', '')
            if not patch:
                continue
            
            # Parse the patch for field additions
            added_fields = self._extract_field_additions_from_patch(patch, used_types_simple)
            
            for field_info in added_fields:
                opportunities.append({
                    'type': field_info['type'],
                    'field': field_info['field'],
                    'field_type': field_info['field_type'],
                    'file': filename,
                    'description': field_info.get('comment', ''),
                    'currently_used': field_info['type'] in used_types_simple
                })
        
        return opportunities
    
    def _extract_field_additions_from_patch(self, patch: str, used_types: set) -> List[Dict[str, str]]:
        """Extract field additions from a git patch.
        
        Args:
            patch: Git diff patch content
            used_types: Set of type names we're using
            
        Returns:
            List of field additions with metadata
        """
        additions = []
        current_struct = None
        last_comment = None
        
        for line in patch.split('\n'):
            # Check for struct declaration
            struct_name = self._parse_struct_declaration(line, used_types)
            if struct_name is not None:
                current_struct = struct_name
                continue
            
            # Check for comment
            comment = self._parse_field_comment(line)
            if comment:
                last_comment = comment
                continue
            
            # Check for field addition
            if current_struct and line.startswith('+'):
                field_info = self._parse_field_addition(line, current_struct)
                if field_info:
                    field_info['comment'] = last_comment or ''
                    additions.append(field_info)
                    last_comment = None
        
        return additions
    
    def _parse_struct_declaration(self, line: str, used_types: set) -> str:
        """Parse struct declaration and return struct name if used."""
        pattern = re.compile(r'^[\+\s]*type\s+(\w+)\s+struct\s*\{', re.MULTILINE)
        match = pattern.match(line)
        if match:
            struct_name = match.group(1)
            return struct_name if struct_name in used_types else ''
        return None
    
    def _parse_field_comment(self, line: str) -> str:
        """Extract comment from line if present."""
        if line.strip().startswith('+') and '//' in line:
            return line.split('//', 1)[1].strip()
        return ''
    
    def _parse_field_addition(self, line: str, struct_name: str) -> Dict[str, str]:
        """Parse field addition from patch line."""
        pattern = re.compile(r'^\+\s+(\w+)\s+\*?(\w+)\s+`json:"([^"]+)"', re.MULTILINE)
        match = pattern.match(line)
        
        if not match:
            return None
        
        field_name = match.group(1)
        field_type = match.group(2)
        json_name = match.group(3).split(',')[0]
        
        # Skip internal/private fields
        if field_name.startswith('_') or not field_name[0].isupper():
            return None
        
        return {
            'type': struct_name,
            'field': field_name,
            'field_type': field_type,
            'json_name': json_name
        }


def generate_summary_stats(categorized: Dict[str, List[Dict]]) -> Dict[str, int]:
    """Generate summary statistics from categorized changes.
    
    Args:
        categorized: Categorized changes from analyze_file_changes()
        
    Returns:
        Dictionary of counts:
        {
            "total_changes": 1234,
            "relevant_changes": 15,
            "critical": 3,
            "warning": 7,
            "safe": 5,
            "noise": 1219
        }
    """
    total = sum(len(changes) for changes in categorized.values())
    
    return {
        "total_changes": total,
        "relevant_changes": total - len(categorized[ImpactLevel.NOISE]),
        "critical": len(categorized[ImpactLevel.CRITICAL]),
        "warning": len(categorized[ImpactLevel.WARNING]),
        "safe": len(categorized[ImpactLevel.SAFE]),
        "opportunity": len(categorized.get(ImpactLevel.OPPORTUNITY, [])),
        "noise": len(categorized[ImpactLevel.NOISE]),
    }
