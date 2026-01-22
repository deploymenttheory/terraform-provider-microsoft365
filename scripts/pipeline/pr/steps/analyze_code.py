#!/usr/bin/env python3
"""Analyze code for service domains and goroutines.

Outputs:
- service-domains: Space-separated list of affected service domains
- has-goroutines: 'true' if goroutines detected, 'false' otherwise
- goroutine-packages: Space-separated list of packages with goroutines
"""

import argparse
import os
import sys
from pathlib import Path

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from lib.code_analysis import detect_service_domains, detect_goroutines  # noqa: E402
from lib.common import load_pr_checks_config  # noqa: E402
from lib.github import write_output  # noqa: E402


def main():
    """Analyze code for service domains and goroutines."""
    parser = argparse.ArgumentParser(description='Analyze code patterns')
    parser.add_argument('--packages', required=True,
                        help='Space-separated list of packages to analyze')
    parser.add_argument('--config', default=None,
                        help='Path to PR checks config file')
    parser.add_argument('--github-output', default=os.environ.get('GITHUB_OUTPUT'),
                        help='Path to GITHUB_OUTPUT file')
    
    args = parser.parse_args()
    
    # Parse packages from space-separated string
    packages = args.packages.split()
    config = load_pr_checks_config(args.config)
    
    # Detect service domains
    print("\n📊 Detecting service domains...")
    service_domains = detect_service_domains(packages, config)
    print(f"Service domains: {' '.join(service_domains) if service_domains else 'N/A'}")
    
    # Detect goroutines
    print("\n🔍 Scanning for goroutines...")
    packages_with_goroutines = detect_goroutines(packages)
    has_goroutines = len(packages_with_goroutines) > 0
    
    print(f"\n{'✅' if has_goroutines else 'ℹ️ '} Found {len(packages_with_goroutines)} package(s) with goroutines")
    for pkg in packages_with_goroutines:
        print(f"   - {pkg}")
    
    write_output({
        "service-domains": ' '.join(service_domains),
        "has-goroutines": "true" if has_goroutines else "false",
        "goroutine-packages": ' '.join(packages_with_goroutines)
    }, args.github_output)
    
    return 0


if __name__ == "__main__":
    sys.exit(main())
