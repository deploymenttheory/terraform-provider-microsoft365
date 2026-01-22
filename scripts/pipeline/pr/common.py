#!/usr/bin/env python3
"""Common utilities for PR pipeline scripts.

Shared functions for input handling, output formatting, and GitHub Actions integration.
"""

import os
import sys
from pathlib import Path
from typing import List, Dict, Any


def get_packages_from_input(packages_arg: List[str], use_stdin: bool = False, 
                            env_var: str = None) -> List[str]:
    """Get package list from various input sources.
    
    Args:
        packages_arg: Package list from command-line arguments.
        use_stdin: Whether to read from stdin.
        env_var: Environment variable name to read from.
    
    Returns:
        List of package paths.
    """
    if use_stdin:
        return [line.strip() for line in sys.stdin if line.strip()]
    elif env_var:
        env_value = os.environ.get(env_var, '')
        return [p.strip() for p in env_value.split() if p.strip()]
    else:
        return packages_arg


def write_github_output(outputs: Dict[str, Any], github_output_path: str = None) -> None:
    """Write outputs to GitHub Actions output file.
    
    Args:
        outputs: Dictionary of key-value pairs to write.
        github_output_path: Path to GITHUB_OUTPUT file (defaults to env var).
    """
    output_file = github_output_path or os.environ.get('GITHUB_OUTPUT')
    
    if not output_file:
        return
    
    with open(output_file, 'a', encoding='utf-8') as f:
        for key, value in outputs.items():
            f.write(f"{key}={value}\n")


def sanitize_package_path(package: str) -> str:
    """Sanitize package path for use as filename.
    
    Args:
        package: Package path (e.g., 'internal/services/common/state').
    
    Returns:
        Sanitized filename (e.g., 'internal_services_common_state').
    """
    return package.replace('/', '_').replace('.', '_').strip('_')
