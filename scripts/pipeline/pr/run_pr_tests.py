#!/usr/bin/env python3
"""Unified PR test orchestrator for both unit tests and coverage workflows.

This script consolidates the logic from pr-unit-tests.yml and pr-coverage.yml
into a single executable with different modes.

Usage:
    ./run_pr_tests.py --mode unit-tests --base-ref origin/main
    ./run_pr_tests.py --mode coverage --base-ref origin/main

Modes:
    unit-tests: Run tests with race detection, block on missing tests
    coverage: Run tests with coverage, upload to Codecov, comment on PR
"""

import sys
import os
import subprocess
import argparse
import json
from pathlib import Path
from typing import List, Dict, Optional, Any

# Import local utilities
from common import get_packages_from_input, write_github_output


def get_changed_packages(base_ref: str) -> List[str]:
    """Get list of changed Go packages using git diff.
    
    Args:
        base_ref: Base branch reference.
    
    Returns:
        List of unique package directories.
    """
    try:
        result = subprocess.run(
            ["git", "diff", "--name-only", f"{base_ref}...HEAD"],
            capture_output=True,
            text=True,
            check=True
        )
        
        go_files = [
            line.strip()
            for line in result.stdout.split('\n')
            if line.strip().endswith('.go')
        ]
        
        if not go_files:
            return []
        
        print(f"📝 Found {len(go_files)} changed Go files")
        
        packages = set()
        for file_path in go_files:
            package_dir = str(Path(file_path).parent)
            packages.add(package_dir)
        
        return sorted(list(packages))
    
    except subprocess.CalledProcessError as e:
        print(f"❌ Error running git diff: {e.stderr}", file=sys.stderr)
        sys.exit(1)


def should_skip_package(package: str, skip_patterns: List[str]) -> tuple:
    """Check if package should be skipped from test requirements.
    
    Args:
        package: Package directory path.
        skip_patterns: List of directory patterns to skip.
    
    Returns:
        Tuple of (should_skip: bool, reason: str).
    """
    for pattern in skip_patterns:
        if pattern in package:
            return True, pattern
    return False, ""


def check_packages_have_tests(packages: List[str]) -> Dict[str, Any]:
    """Check if packages have test files.
    
    Args:
        packages: List of package paths.
    
    Returns:
        Dict with test coverage status.
    """
    skip_patterns = ["/mocks", "/schema", "/shared_models/", "internal/acceptance", "internal/constants"]
    
    packages_with_tests = []
    missing_tests = []
    
    for package in packages:
        should_skip, reason = should_skip_package(package, skip_patterns)
        
        if should_skip:
            print(f"⏭️  Skipping: {package} ({reason})")
            continue
        
        package_path = Path(package)
        test_files = list(package_path.glob("*_test.go")) if package_path.exists() else []
        
        if test_files:
            print(f"✅ Found {len(test_files)} test file(s) in {package}")
            packages_with_tests.append(package)
        else:
            print(f"❌ No test files found in {package}")
            missing_tests.append(package)
    
    return {
        "packages_with_tests": packages_with_tests,
        "missing_tests": missing_tests,
        "has_missing": len(missing_tests) > 0
    }


def determine_service_areas(packages: List[str]) -> List[str]:
    """Determine M365 service areas from package paths.
    
    Args:
        packages: List of package paths.
    
    Returns:
        List of unique service area names.
    """
    import re
    
    service_areas = set()
    
    for package in packages:
        # Check for service directories
        for category in ['resources', 'datasources', 'actions']:
            match = re.search(rf'internal/services/{category}/([^/]+)', package)
            if match:
                service_areas.add(match.group(1))
                break
        
        # Check for provider core
        core_patterns = ['internal/client', 'internal/helpers', 'internal/provider', 'internal/utilities']
        if any(pattern in package for pattern in core_patterns):
            service_areas.add('provider-core')
    
    return sorted(list(service_areas))


def run_unit_tests_with_race(packages: List[str]) -> int:
    """Run unit tests with race detection.
    
    Args:
        packages: List of package paths to test.
    
    Returns:
        Exit code: 0 if all pass, 1 if any fail.
    """
    print("\n" + "="*60)
    print("🧪 Running Unit Tests with Race Detection")
    print("="*60)
    
    has_failures = False
    
    for idx, package in enumerate(packages, 1):
        print(f"\n[{idx}/{len(packages)}] Testing: {package}")
        
        cmd = ["go", "test", "-v", "-race", f"./{package}"]
        
        result = subprocess.run(
            cmd,
            env={"TF_ACC": "0", **os.environ},
            check=False
        )
        
        if result.returncode != 0:
            has_failures = True
            print(f"❌ Tests failed in {package}")
        else:
            print(f"✅ Tests passed in {package}")
    
    return 1 if has_failures else 0


def run_tests_with_coverage(packages: List[str], output_dir: str = "coverage") -> Path:
    """Run unit tests with coverage profiling.
    
    Args:
        packages: List of package paths to test.
        output_dir: Directory for coverage files.
    
    Returns:
        Path to merged coverage file.
    """
    print("\n" + "="*60)
    print("📊 Running Unit Tests with Coverage")
    print("="*60)
    
    coverage_dir = Path(output_dir)
    coverage_dir.mkdir(parents=True, exist_ok=True)
    
    merged_file = coverage_dir / "unit-coverage.txt"
    coverage_files = []
    
    for idx, package in enumerate(packages, 1):
        safe_name = package.replace('/', '_').replace('.', '_').strip('_')
        coverage_file = coverage_dir / f"{safe_name}.out"
        
        print(f"\n[{idx}/{len(packages)}] Testing: {package}")
        
        cmd = [
            "go", "test", "-v",
            f"-coverprofile={coverage_file}",
            "-covermode=atomic",
            f"./{package}"
        ]
        
        result = subprocess.run(
            cmd,
            env={"TF_ACC": "0", **os.environ},
            check=False
        )
        
        if coverage_file.exists():
            coverage_files.append(coverage_file)
            print(f"✅ Coverage generated for {package}")
        else:
            print(f"⚠️  No coverage file for {package}")
    
    # Merge coverage files
    print(f"\n📊 Merging {len(coverage_files)} coverage file(s)...")
    with open(merged_file, 'w', encoding='utf-8') as out_f:
        out_f.write("mode: atomic\n")
        for cov_file in coverage_files:
            with open(cov_file, 'r', encoding='utf-8') as in_f:
                for line in in_f:
                    if not line.startswith('mode:'):
                        out_f.write(line)
    
    print(f"✅ Merged coverage: {merged_file}")
    return merged_file


def calculate_coverage(coverage_file: Path) -> Dict[str, Any]:
    """Calculate coverage statistics from Go coverage file.
    
    Args:
        coverage_file: Path to coverage file.
    
    Returns:
        Dict with coverage statistics.
    """
    if not coverage_file.exists():
        return {"total_lines": 0, "covered_lines": 0, "coverage_pct": 0.0}
    
    total_statements = 0
    covered_statements = 0
    
    with open(coverage_file, 'r', encoding='utf-8') as f:
        for line in f:
            if line.startswith('mode:'):
                continue
            
            parts = line.strip().split()
            if len(parts) >= 3:
                try:
                    statements = int(parts[1])
                    count = int(parts[2])
                    total_statements += statements
                    if count > 0:
                        covered_statements += statements
                except (ValueError, IndexError):
                    continue
    
    coverage_pct = (covered_statements / total_statements * 100) if total_statements > 0 else 0.0
    
    return {
        "total_lines": total_statements,
        "covered_lines": covered_statements,
        "coverage_pct": round(coverage_pct, 2)
    }


def main():
    """Main entry point for PR test orchestrator.
    
    Handles both unit-tests and coverage modes with consolidated logic.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    parser = argparse.ArgumentParser(
        description='Unified PR test orchestrator',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument('--mode', required=True, choices=['unit-tests', 'coverage'],
                        help='Test mode: unit-tests (with race) or coverage (with profiling)')
    parser.add_argument('--base-ref', required=True,
                        help='Base branch reference (e.g., origin/main)')
    parser.add_argument('--github-output', default=os.environ.get('GITHUB_OUTPUT'),
                        help='Path to GITHUB_OUTPUT file')
    parser.add_argument('--output-dir', default='coverage',
                        help='Directory for coverage files (coverage mode only)')
    
    args = parser.parse_args()
    
    print("="*60)
    print(f"🚀 PR Test Orchestrator - Mode: {args.mode}")
    print("="*60)
    
    # Step 1: Identify changed packages
    print("\n📦 Step 1: Identifying changed packages...")
    packages = get_changed_packages(args.base_ref)
    
    if not packages:
        print("✅ No Go files changed")
        write_github_output({
            "packages": "",
            "has-changes": "false"
        }, args.github_output)
        return 0
    
    print(f"📦 Found {len(packages)} changed package(s)")
    for pkg in packages:
        print(f"   - {pkg}")
    
    write_github_output({
        "packages": ' '.join(packages),
        "has-changes": "true"
    }, args.github_output)
    
    # Step 2: Mode-specific validation
    if args.mode == 'unit-tests':
        print("\n🧪 Step 2: Checking for test files...")
        test_status = check_packages_have_tests(packages)
        
        if test_status["has_missing"]:
            print("\n❌ ERROR: Some packages are missing tests")
            for pkg in test_status["missing_tests"]:
                print(f"   - {pkg}")
            return 1
        
        packages_to_test = test_status["packages_with_tests"]
        write_github_output({
            "packages-with-tests": ' '.join([f"./{pkg}" for pkg in packages_to_test])
        }, args.github_output)
        
        # Step 3: Run unit tests with race detection
        print("\n🧪 Step 3: Running unit tests...")
        return run_unit_tests_with_race(packages_to_test)
    
    else:  # coverage mode
        # Step 2: Determine service areas
        print("\n📊 Step 2: Determining service areas...")
        service_areas = determine_service_areas(packages)
        print(f"Service areas: {' '.join(service_areas)}")
        
        write_github_output({
            "service-areas": ' '.join(service_areas)
        }, args.github_output)
        
        # Step 3: Run tests with coverage
        print("\n📊 Step 3: Running tests with coverage...")
        merged_coverage = run_tests_with_coverage(packages, args.output_dir)
        
        # Step 4: Calculate coverage summary
        print("\n📊 Step 4: Calculating coverage...")
        stats = calculate_coverage(merged_coverage)
        
        print(f"Coverage: {stats['coverage_pct']}%")
        print(f"Total: {stats['total_lines']} statements")
        print(f"Covered: {stats['covered_lines']} statements")
        
        write_github_output({
            "coverage-pct": str(stats['coverage_pct']),
            "total-lines": str(stats['total_lines']),
            "covered-lines": str(stats['covered_lines'])
        }, args.github_output)
        
        return 0


if __name__ == "__main__":
    sys.exit(main())
