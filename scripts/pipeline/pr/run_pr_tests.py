#!/usr/bin/env python3
"""Unified PR test orchestrator for the pr-tests workflow.

This script orchestrates PR validation by coordinating:
- Git operations (detecting changed code)
- Code analysis (service domains, goroutines)
- Go test execution
- Coverage calculation and enforcement

Usage:
    # Run tests and calculate coverage locally (draft and ready PRs)
    ./run_pr_tests.py --mode unit-tests --base-ref origin/main [--is-draft]
    
    # Fetch coverage from Codecov and enforce threshold (ready PRs only)
    ./run_pr_tests.py --mode enforce-coverage --repo-slug owner/repo --pr-number 123 --codecov-token TOKEN
    
    # Run race detection
    ./run_pr_tests.py --mode race-detection --base-ref origin/main

Modes:
    unit-tests: Run tests with coverage profiling, calculate coverage locally
    enforce-coverage: Fetch coverage from Codecov API and enforce threshold
    race-detection: Run race detector on packages with goroutines
"""

import argparse
import os
import sys
from typing import Any, Dict, List

# Local imports
from code_analysis import detect_service_domains, detect_goroutines
from common import load_pr_checks_config
from coverage import calculate_coverage, check_coverage_threshold
from get_codecov_coverage import fetch_codecov_coverage
from git_operations import get_changed_packages
from github import write_output
from go_tests import run_unit_tests, run_race_detection


def run_unit_tests_mode(packages: List[str], config: Dict[str, Any], 
                        args: argparse.Namespace) -> int:
    """Execute unit tests mode with coverage analysis.
    
    Workflow:
    1. Detect service domains
    2. Scan for goroutines
    3. Run unit tests with coverage
    4. Calculate coverage locally
    5. Output results
    
    Args:
        packages: List of changed packages.
        config: Test configuration.
        args: Parsed command-line arguments.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    # Step 2: Determine service domains
    print("\n📊 Step 2: Determining service domains...")
    service_domains = detect_service_domains(packages, config)
    print(f"Service domains: {' '.join(service_domains) if service_domains else 'N/A'}")
    
    # Step 3: Detect goroutines for race detection job
    print("\n🔍 Step 3: Scanning for goroutines...")
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
    
    # Step 4: Run tests with coverage
    print("\n📊 Step 4: Running tests with coverage...")
    coverage_file = run_unit_tests(packages, args.output_dir)
    
    # Step 5: Calculate local coverage
    print("\n📊 Step 5: Calculating local coverage...")
    stats = calculate_coverage(coverage_file)
    
    print(f"\n{'='*60}")
    print(f"Local Coverage: {stats['coverage_pct']}%")
    print(f"Total Statements: {stats['total_lines']}")
    print(f"Covered Statements: {stats['covered_lines']}")
    print(f"{'='*60}")
    
    write_output({
        "coverage-pct": str(stats['coverage_pct']),
        "total-lines": str(stats['total_lines']),
        "covered-lines": str(stats['covered_lines'])
    }, args.github_output)
    
    # For draft PRs, inform but don't enforce
    if args.is_draft:
        print("\n📝 Draft PR: Coverage calculated locally (informational only)")
        print("   Mark PR as ready for review to upload to Codecov and enforce threshold")
        return 0
    
    # For ready PRs, coverage will be uploaded to Codecov by workflow
    print("\n✅ Coverage file generated - will be uploaded to Codecov by workflow")
    return 0


def run_enforce_coverage_mode(config: Dict[str, Any], args: argparse.Namespace) -> int:
    """Fetch coverage from Codecov and enforce threshold.
    
    Workflow:
    1. Fetch coverage from Codecov API
    2. Check against minimum threshold
    3. Return success/failure
    
    Args:
        config: Test configuration.
        args: Parsed command-line arguments.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    if not args.codecov_token or not args.repo_slug or not args.pr_number:
        print("❌ ERROR: Missing required Codecov parameters (token/repo/pr)")
        return 1
    
    print("\n📊 Fetching coverage from Codecov...")
    stats = fetch_codecov_coverage(args.repo_slug, args.pr_number, args.codecov_token)
    
    if not stats:
        print("❌ ERROR: Could not fetch coverage from Codecov")
        return 1
    
    print(f"\n{'='*60}")
    print(f"Codecov Patch Coverage: {stats['coverage_pct']}%")
    print(f"Total Lines Changed: {stats['total_lines']}")
    print(f"Covered Lines: {stats['covered_lines']}")
    print(f"{'='*60}")
    
    write_output({
        "coverage-pct": str(stats['coverage_pct']),
        "total-lines": str(stats['total_lines']),
        "covered-lines": str(stats['covered_lines'])
    }, args.github_output)
    
    # Enforce coverage threshold
    min_coverage = config.get('coverage_threshold', {}).get('minimum_pct', 60)
    
    if not check_coverage_threshold(stats['coverage_pct'], min_coverage):
        print(f"\n❌ ERROR: Coverage {stats['coverage_pct']}% is below minimum threshold {min_coverage}%")
        print("   Please add tests to increase coverage for changed code.")
        return 1
    
    print(f"\n✅ Coverage {stats['coverage_pct']}% meets minimum threshold {min_coverage}%")
    return 0


def run_race_detection_mode(packages: List[str]) -> int:
    """Execute race detection mode on packages with goroutines.
    
    Workflow:
    1. Detect packages with goroutines
    2. Run race detection tests
    
    Args:
        packages: List of changed packages.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    # Step 2: Detect packages with goroutines
    print("\n🔍 Step 2: Detecting packages with goroutines...")
    packages_with_goroutines = detect_goroutines(packages)
    
    if not packages_with_goroutines:
        print("✅ No packages with goroutines found, skipping race detection")
        return 0
    
    print(f"\n🔍 Found {len(packages_with_goroutines)} package(s) with goroutines")
    for pkg in packages_with_goroutines:
        print(f"   - {pkg}")
    
    # Step 3: Run race detection
    print("\n🔍 Step 3: Running race detection tests...")
    return run_race_detection(packages_with_goroutines)


def main():
    """Main entry point for PR test orchestrator.
    
    Handles three modes:
    - unit-tests: Run tests with coverage (calculate locally)
    - enforce-coverage: Fetch from Codecov and enforce threshold
    - race-detection: Run race detector
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    parser = argparse.ArgumentParser(
        description='Unified PR test orchestrator',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument('--mode', required=True, 
                        choices=['unit-tests', 'enforce-coverage', 'race-detection'],
                        help='Test mode: unit-tests, enforce-coverage, or race-detection')
    parser.add_argument('--base-ref', required=False,
                        help='Base branch reference (e.g., origin/main) - required for unit-tests and race-detection')
    parser.add_argument('--config', default=None,
                        help='Path to PR checks config file (defaults to pr-checks-config.yml in repo root)')
    parser.add_argument('--github-output', default=os.environ.get('GITHUB_OUTPUT'),
                        help='Path to GITHUB_OUTPUT file')
    parser.add_argument('--output-dir', default='coverage',
                        help='Directory for coverage files (unit-tests mode only)')
    parser.add_argument('--is-draft', action='store_true',
                        help='PR is in draft state (informational mode, unit-tests mode only)')
    parser.add_argument('--codecov-token', default=os.environ.get('CODECOV_TOKEN'),
                        help='Codecov API token (enforce-coverage mode only)')
    parser.add_argument('--repo-slug', default=None,
                        help='Repository slug in format owner/repo (enforce-coverage mode only)')
    parser.add_argument('--pr-number', default=None,
                        help='Pull request number (enforce-coverage mode only)')
    
    args = parser.parse_args()
    
    config = load_pr_checks_config(args.config)
    
    print("="*60)
    print(f"🚀 PR Test Orchestrator - Mode: {args.mode}")
    if args.is_draft and args.mode == 'unit-tests':
        print("📝 Draft PR Mode: Tests are informational only")
    print("="*60)
    
    # enforce-coverage mode doesn't need package detection
    if args.mode == 'enforce-coverage':
        return run_enforce_coverage_mode(config, args)
    
    # Validate base-ref for other modes
    if not args.base_ref:
        print("❌ ERROR: --base-ref is required for this mode")
        return 1
    
    # Step 1: Identify changed packages
    print("\n📦 Step 1: Identifying changed packages...")
    packages = get_changed_packages(args.base_ref)
    
    if not packages:
        print("✅ No Go files changed")
        write_output({
            "packages": "",
            "has-changes": "false",
            "has-goroutines": "false"
        }, args.github_output)
        return 0
    
    print(f"📦 Found {len(packages)} changed package(s)")
    for pkg in packages:
        print(f"   - {pkg}")
    
    write_output({
        "packages": ' '.join(packages),
        "has-changes": "true"
    }, args.github_output)
    
    # Mode-specific execution
    if args.mode == 'unit-tests':
        return run_unit_tests_mode(packages, config, args)
    
    # race-detection mode
    return run_race_detection_mode(packages)


if __name__ == "__main__":
    sys.exit(main())
