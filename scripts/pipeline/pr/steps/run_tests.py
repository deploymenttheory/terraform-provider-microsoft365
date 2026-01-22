#!/usr/bin/env python3
"""Run Go unit tests with coverage profiling.

Runs tests for all specified packages and generates merged coverage file.
"""

import argparse
import sys
from pathlib import Path

# Add parent directory to path for lib imports
sys.path.insert(0, str(Path(__file__).parent.parent))

from lib.go_tests import run_unit_tests  # noqa: E402


def main():
    """Run unit tests with coverage."""
    parser = argparse.ArgumentParser(description='Run Go unit tests with coverage')
    parser.add_argument('--packages', required=True,
                        help='Space-separated list of packages to test')
    parser.add_argument('--output-dir', default='coverage',
                        help='Directory for coverage output files')
    
    args = parser.parse_args()
    
    # Parse packages from space-separated string
    packages = args.packages.split()
    
    print(f"🧪 Running unit tests on {len(packages)} package(s)...")
    coverage_file = run_unit_tests(packages, args.output_dir)
    
    print(f"\n✅ Coverage file generated: {coverage_file}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
