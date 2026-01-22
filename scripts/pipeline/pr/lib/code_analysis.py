#!/usr/bin/env python3
"""Code analysis utilities for PR validation.

Provides functions for analyzing Go code to detect patterns like:
- Service domains (which Microsoft 365 services are affected)
- Goroutines (for race detection)
"""

import re
from pathlib import Path
from typing import Dict, Any, List, Set


def detect_service_domains(packages: List[str], config: Dict[str, Any]) -> List[str]:
    """Detect which service domains are affected by changed packages.
    
    Args:
        packages: List of changed package paths.
        config: Configuration dict with service_domain_patterns.
    
    Returns:
        List of unique service domain names (sorted).
    """
    domains: Set[str] = set()
    
    # Get service domain patterns from config
    service_patterns = config.get('service_domain_patterns', {})
    provider_core_paths = config.get('provider_core_paths', [])
    
    for package in packages:
        # Check if package is provider core
        if any(core_path in package for core_path in provider_core_paths):
            domains.add('provider-core')
            continue
        
        # Check service domain patterns
        for pattern in service_patterns.values():
            match = re.search(pattern, package)
            if match:
                service_name = match.group(1)
                domains.add(service_name)
                break
    
    return sorted(domains)


def detect_goroutines(packages: List[str]) -> List[str]:
    """Detect packages that contain goroutines.
    
    Scans Go files for 'go func()' pattern to identify packages
    that spawn goroutines and need race detection testing.
    
    Args:
        packages: List of package paths to scan.
    
    Returns:
        List of packages that contain goroutines.
    """
    packages_with_goroutines = []
    goroutine_pattern = re.compile(r'\bgo\s+func\s*\(')
    
    for package in packages:
        pkg_path = Path(package)
        if not pkg_path.exists():
            continue
        
        # Scan all .go files in package
        for go_file in pkg_path.glob('*.go'):
            try:
                with open(go_file, 'r', encoding='utf-8') as f:
                    content = f.read()
                    if goroutine_pattern.search(content):
                        packages_with_goroutines.append(package)
                        break  # Found goroutine in this package, move to next
            except (OSError, UnicodeDecodeError) as e:
                print(f"⚠️  Could not read {go_file}: {e}")
                continue
    
    return packages_with_goroutines


def has_goroutines(packages: List[str]) -> bool:
    """Check if any packages contain goroutines.
    
    Args:
        packages: List of package paths to check.
    
    Returns:
        True if at least one package has goroutines, False otherwise.
    """
    return len(detect_goroutines(packages)) > 0
