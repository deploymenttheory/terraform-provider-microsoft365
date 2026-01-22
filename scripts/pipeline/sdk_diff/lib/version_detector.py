#!/usr/bin/env python3
"""Version detection utilities for Go modules.

Parses go.mod to extract current SDK versions.
"""

import re
from pathlib import Path
from typing import Dict


def parse_go_mod(go_mod_path: Path) -> Dict[str, str]:
    """Parse go.mod file to extract dependency versions.
    
    Args:
        go_mod_path: Path to go.mod file
        
    Returns:
        Dictionary mapping module names to versions:
        {
            "github.com/microsoftgraph/msgraph-sdk-go": "v1.93.0",
            "github.com/microsoftgraph/msgraph-beta-sdk-go": "v0.157.0"
        }
    """
    if not go_mod_path.exists():
        raise FileNotFoundError(f"go.mod not found at {go_mod_path}")
    
    dependencies = {}
    
    with open(go_mod_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Parse require block
    # Pattern: github.com/module/path v1.2.3
    require_pattern = r'^\s*(github\.com/[^\s]+)\s+(v[\d.]+)'
    
    for line in content.split('\n'):
        match = re.match(require_pattern, line.strip())
        if match:
            module, version = match.groups()
            dependencies[module] = version
    
    return dependencies


def get_msgraph_versions(go_mod_path: Path) -> Dict[str, str]:
    """Extract Microsoft Graph SDK versions from go.mod.
    
    Args:
        go_mod_path: Path to go.mod file
        
    Returns:
        Dictionary of msgraph SDK versions:
        {
            "msgraph-sdk-go": "v1.93.0",
            "msgraph-beta-sdk-go": "v0.157.0",
            "msgraph-sdk-go-core": "v1.4.0"
        }
    """
    all_deps = parse_go_mod(go_mod_path)
    
    msgraph_deps = {}
    for module, version in all_deps.items():
        if "microsoftgraph/msgraph" in module:
            # Extract just the SDK name (last part of path)
            sdk_name = module.split('/')[-1]
            msgraph_deps[sdk_name] = version
    
    return msgraph_deps


def format_version_display(versions: Dict[str, str]) -> str:
    """Format SDK versions for display.
    
    Args:
        versions: Dictionary of SDK versions
        
    Returns:
        Formatted string like:
        - msgraph-sdk-go: v1.93.0
        - msgraph-beta-sdk-go: v0.157.0
    """
    lines = []
    for sdk, version in sorted(versions.items()):
        lines.append(f"  - {sdk}: {version}")
    
    return '\n'.join(lines)
