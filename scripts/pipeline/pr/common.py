#!/usr/bin/env python3
"""Common utilities for PR pipeline scripts.

Shared functions for GitHub Actions integration and configuration management.
"""

import os
import yaml
from pathlib import Path
from typing import Dict, Any, Optional


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


def load_pr_checks_config(config_path: Optional[str] = None) -> Dict[str, Any]:
    """Load PR checks configuration from YAML file.
    
    Args:
        config_path: Path to config file. If None, looks for pr-checks-config.yml
                     in the repository root.
    
    Returns:
        Dictionary with PR checks configuration settings.
    """
    if config_path is None:
        # Default to pr-checks-config.yml in repository root
        repo_root = Path(__file__).parent.parent.parent.parent
        config_path = repo_root / 'pr-checks-config.yml'
    else:
        config_path = Path(config_path)
    
    if not config_path.exists():
        # Return default configuration if file doesn't exist
        return {
            'coverage_threshold': {
                'minimum_pct': 60
            },
            'service_domain_patterns': {
                'resources': 'internal/services/resources/([^/]+)',
                'datasources': 'internal/services/datasources/([^/]+)',
                'actions': 'internal/services/actions/([^/]+)',
                'ephemerals': 'internal/services/ephemerals/([^/]+)',
                'list_resources': 'internal/services/list-resources/([^/]+)',
                'common': 'internal/services/common/([^/]+)'
            },
            'provider_core_paths': [
                'internal/client',
                'internal/helpers',
                'internal/provider',
                'internal/utilities'
            ]
        }
    
    with open(config_path, 'r', encoding='utf-8') as f:
        return yaml.safe_load(f)
