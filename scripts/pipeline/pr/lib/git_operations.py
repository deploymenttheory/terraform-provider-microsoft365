#!/usr/bin/env python3
"""Git operations for PR analysis.

Provides functions for querying Git history and identifying changed files.
"""

import subprocess
from pathlib import Path
from typing import List, Set


def get_merge_base(base_ref: str) -> str:
    """Resolve the comparison baseline using local Git history."""
    return subprocess.run(
        ["git", "merge-base", base_ref, "HEAD"],
        capture_output=True, text=True, check=True
    ).stdout.strip()


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
        ["git", "diff", "--name-only", "--no-renames", "-z",
         get_merge_base(base_ref), "HEAD", "--"],
        capture_output=True,
        text=True,
        check=True
    )
    
    return [
        name
        for name in result.stdout.split('\0')
        if name.endswith(file_extension)
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
    
    # Include both sides of renames and directories affected by deletions.
    # Completely deleted packages cannot be passed to Go tooling.
    packages: Set[str] = set()
    for file_path in go_files:
        directory = Path(file_path).parent
        if directory.is_dir() and any(directory.glob('*.go')):
            packages.add(directory.as_posix())
    
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
