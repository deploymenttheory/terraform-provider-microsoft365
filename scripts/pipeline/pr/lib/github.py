#!/usr/bin/env python3
"""GitHub Actions integration utilities.

Provides functions for interacting with GitHub Actions environment,
including writing outputs and reading environment variables.
"""

from pathlib import Path
from typing import Dict, Optional


def write_output(outputs: Dict[str, str], output_file: Optional[str] = None) -> None:
    """Write key-value pairs to GitHub Actions output file.
    
    Args:
        outputs: Dictionary of output key-value pairs.
        output_file: Path to GITHUB_OUTPUT file (defaults to None for testing).
    """
    if not output_file:
        return
    
    output_path = Path(output_file)
    
    with open(output_path, 'a', encoding='utf-8') as f:
        for key, value in outputs.items():
            f.write(f"{key}={value}\n")


def is_draft_pr() -> bool:
    """Check if current PR is in draft state from GitHub environment.
    
    Returns:
        True if PR is draft, False otherwise.
    """
    import os
    return os.environ.get('GITHUB_EVENT_PULL_REQUEST_DRAFT', 'false').lower() == 'true'


def get_pr_number() -> Optional[str]:
    """Get current PR number from GitHub environment.
    
    Returns:
        PR number as string or None if not available.
    """
    import os
    return os.environ.get('GITHUB_EVENT_PULL_REQUEST_NUMBER')


def get_repo_slug() -> Optional[str]:
    """Get repository slug (owner/repo) from GitHub environment.
    
    Returns:
        Repository slug or None if not available.
    """
    import os
    return os.environ.get('GITHUB_REPOSITORY')
