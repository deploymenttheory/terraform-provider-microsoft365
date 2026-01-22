#!/usr/bin/env python3
"""Git operations for PR analysis.

Provides functions for querying Git history and identifying changed files.
"""

import subprocess
from typing import List, Set


def get_changed_files(base_ref: str, file_extension: str = ".go") -> List[str]:
    """Get list of changed files with specific extension.
    
    Args:
        base_ref: Base branch reference (e.g., 'origin/main').
        file_extension: File extension to filter (default: '.go').
    
    Returns:
        List of changed file paths.
    
    Raises:
        subprocess.CalledProcessError: If git command fails.
    """
    result = subprocess.run(
        ["git", "diff", "--name-only", f"{base_ref}...HEAD"],
        capture_output=True,
        text=True,
        check=True
    )
    
    return [
        line.strip()
        for line in result.stdout.split('\n')
        if line.strip().endswith(file_extension)
    ]


def get_changed_packages(base_ref: str) -> List[str]:
    """Get list of changed Go packages.
    
    Args:
        base_ref: Base branch reference.
    
    Returns:
        List of unique package directories (sorted).
    """
    go_files = get_changed_files(base_ref, ".go")
    
    if not go_files:
        return []
    
    # Extract unique package directories
    packages: Set[str] = set()
    for file_path in go_files:
        parts = file_path.split('/')
        if len(parts) > 1:
            # Get package directory (exclude filename)
            pkg_path = '/'.join(parts[:-1])
            packages.add(pkg_path)
    
    return sorted(packages)


def read_file_content(file_path: str) -> str:
    """Read content of a file from the repository.
    
    Args:
        file_path: Path to file relative to repository root.
    
    Returns:
        File content as string.
    
    Raises:
        FileNotFoundError: If file doesn't exist.
    """
    with open(file_path, 'r', encoding='utf-8') as f:
        return f.read()
