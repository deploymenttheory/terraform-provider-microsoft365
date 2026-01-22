#!/usr/bin/env python3
"""Detect changed Go packages in PR.

Outputs:
- packages: Space-separated list of changed packages
- has-changes: 'true' if changes found, 'false' otherwise
"""

import argparse
import os
import sys
from pathlib import Path

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from lib.git_operations import get_changed_packages  # noqa: E402
from lib.github import write_output  # noqa: E402


def main():
    """Detect and output changed Go packages."""
    parser = argparse.ArgumentParser(description='Detect changed Go packages')
    parser.add_argument('--base-ref', required=True,
                        help='Base branch reference (e.g., origin/main)')
    parser.add_argument('--github-output', default=os.environ.get('GITHUB_OUTPUT'),
                        help='Path to GITHUB_OUTPUT file')
    
    args = parser.parse_args()
    
    print("📦 Detecting changed Go packages...")
    packages = get_changed_packages(args.base_ref)
    
    if not packages:
        print("✅ No Go files changed")
        write_output({
            "packages": "",
            "has-changes": "false"
        }, args.github_output)
        return 0
    
    print(f"📦 Found {len(packages)} changed package(s):")
    for pkg in packages:
        print(f"   - {pkg}")
    
    write_output({
        "packages": ' '.join(packages),
        "has-changes": "true"
    }, args.github_output)
    
    return 0


if __name__ == "__main__":
    sys.exit(main())
