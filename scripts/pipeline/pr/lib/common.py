#!/usr/bin/env python3
"""Configuration management for PR pipeline scripts.

Provides functions for loading and managing PR checks configuration.
"""

import yaml
from pathlib import Path
from typing import Dict, Any, Optional


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
