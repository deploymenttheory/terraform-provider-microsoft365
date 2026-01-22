#!/usr/bin/env python3
"""Run Go race detector on packages with goroutines.

Exit codes:
- 0: All race detection tests passed
- 1: One or more race detection tests failed
"""

import argparse
import sys
from pathlib import Path

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from lib.go_tests import run_race_detection  # noqa: E402


def main():
    """Run race detection tests."""
    parser = argparse.ArgumentParser(description='Run Go race detection')
    parser.add_argument('--packages', required=True,
                        help='Space-separated list of packages to test')
    
    args = parser.parse_args()
    
    # Parse packages from space-separated string
    packages = args.packages.split()
    
    print(f"🔍 Running race detection on {len(packages)} package(s)...")
    return run_race_detection(packages)


if __name__ == "__main__":
    sys.exit(main())
