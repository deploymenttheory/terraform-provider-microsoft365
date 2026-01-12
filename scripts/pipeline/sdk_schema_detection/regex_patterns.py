"""Compiled regex patterns for schema detection."""

import re


class RegexPatterns:
    """Collection of compiled regex patterns used throughout the schema detection."""
    
    # Version patterns
    VERSION_FULL = re.compile(r'v(\d+)\.(\d+)\.(\d+)')
    VERSION_IN_TEXT = re.compile(r'v\d+\.\d+\.\d+')
    
    # GitHub URL patterns
    GITHUB_REPO_URL = re.compile(r'github\.com[:/](.+/.+?)(\.git)?$')
    
    # File patterns
    MODEL_FILE_PATH = re.compile(r'models/[\w_]+\.go')
    
    # Go code patterns - Structs
    GO_STRUCT_FIELD = re.compile(r'(\w+)\s+([\*\[\]]?[\w\.]+(?:\[[\w\.]+\])?)\s*(?:`.*`)?')
    
    # Go code patterns - Interfaces
    GO_INTERFACE_METHOD = re.compile(r'(\w+)\s*\(([^)]*)\)\s*(\([^)]*\))?')
    GO_EMBEDDED_TYPE = re.compile(r'^\s*(\w+[\w\.]*)\s*$')
    
    # Go declarations
    GO_TYPE_STRUCT = re.compile(r'type\s+(\w+)\s+struct')
    GO_TYPE_INTERFACE = re.compile(r'type\s+(\w+)\s+interface')
    
    # Changelog patterns
    CHANGELOG_VERSION_HEADER = r'##'
