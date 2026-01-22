#!/usr/bin/env python3
"""Unified PR test orchestrator for the pr-tests workflow.

This script handles test execution, coverage analysis, and goroutine detection
for pull request validation.

Usage:
    ./run_pr_tests.py --mode unit-tests --base-ref origin/main [--is-draft]
    ./run_pr_tests.py --mode race-detection --base-ref origin/main

Modes:
    unit-tests: Run tests with coverage profiling, enforce threshold
    race-detection: Run race detector on packages with goroutines
"""

import argparse
import os
import re
import subprocess
import sys
import time
import json
from pathlib import Path
from typing import Any, Dict, List, Optional

# Import local utilities
from common import load_pr_checks_config, write_github_output


def get_changed_go_packages(base_ref: str) -> List[str]:
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


def check_packages_have_tests(packages: List[str], skip_patterns: List[str]) -> Dict[str, Any]:
    """Check if packages have test files.
    
    Args:
        packages: List of package paths.
        skip_patterns: List of directory patterns to skip from test requirements.
    
    Returns:
        Dict with test coverage status.
    """
    
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


def determine_service_domains(packages: List[str], config: Dict[str, Any]) -> List[str]:
    """Determine M365 service domains from package paths.
    
    Args:
        packages: List of package paths.
        config: Test configuration dictionary.
    
    Returns:
        List of unique service domain names.
    """
    service_domains = set()
    service_patterns = config.get('service_domain_patterns', {})
    core_paths = config.get('provider_core_paths', [])
    
    for package in packages:
        # Check for service directories using patterns from config
        for pattern in service_patterns.values():
            match = re.search(pattern, package)
            if match:
                service_domains.add(match.group(1))
                break
        
        # Check for provider core
        if any(pattern in package for pattern in core_paths):
            service_domains.add('provider-core')
    
    return sorted(list(service_domains))


def detect_goroutines_in_packages(packages: List[str]) -> List[str]:
    """Detect packages that use goroutines (go func() patterns).
    
    Args:
        packages: List of package paths to scan.
    
    Returns:
        List of packages that contain goroutine usage.
    """
    goroutine_pattern = re.compile(r'\bgo\s+func\s*\(')
    packages_with_goroutines = []
    
    for package in packages:
        package_path = Path(package)
        if not package_path.exists():
            continue
        
        go_files = list(package_path.glob("*.go"))
        # Exclude test files from goroutine detection
        go_files = [f for f in go_files if not f.name.endswith('_test.go')]
        
        has_goroutines = False
        for go_file in go_files:
            try:
                with open(go_file, 'r', encoding='utf-8') as f:
                    content = f.read()
                    if goroutine_pattern.search(content):
                        has_goroutines = True
                        print(f"🔍 Found goroutine in {go_file}")
                        break
            except (IOError, OSError) as e:
                print(f"⚠️  Error reading {go_file}: {e}")
                continue
        
        if has_goroutines:
            packages_with_goroutines.append(package)
    
    return packages_with_goroutines


def run_race_detection(packages: List[str]) -> int:
    """Run tests with race detection on packages with goroutines.
    
    Args:
        packages: List of package paths to test.
    
    Returns:
        Exit code: 0 if all pass, 1 if any fail.
    """
    print("\n" + "="*60)
    print("🔍 Running Race Detection Tests")
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
            print(f"❌ Race detection failed in {package}")
        else:
            print(f"✅ Race detection passed in {package}")
    
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
        
        subprocess.run(
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


def fetch_codecov_coverage(repo_slug: str, pr_number: str, codecov_token: str,
                           max_retries: int = 30, retry_delay: int = 10) -> Optional[Dict[str, Any]]:
    """Fetch coverage results from Codecov API.
    
    Args:
        repo_slug: Repository in format 'owner/repo'.
        pr_number: Pull request number.
        codecov_token: Codecov API token.
        max_retries: Maximum number of retry attempts.
        retry_delay: Delay between retries in seconds.
    
    Returns:
        Dict with coverage statistics or None if fetch fails.
    """
    import urllib.request
    import urllib.error
    
    api_url = f"https://api.codecov.io/api/v2/github/{repo_slug}/pulls/{pr_number}"
    
    print(f"\n⏳ Waiting for Codecov to process coverage (max {max_retries * retry_delay}s)...")
    
    for attempt in range(1, max_retries + 1):
        try:
            req = urllib.request.Request(api_url)
            req.add_header('Authorization', f'Bearer {codecov_token}')
            
            with urllib.request.urlopen(req, timeout=30) as response:
                data = json.loads(response.read().decode())
                
                # Check if we have patch coverage data
                if 'totals' in data and 'patch' in data['totals']:
                    patch_coverage = data['totals']['patch']
                    
                    if patch_coverage and 'coverage' in patch_coverage:
                        coverage_pct = patch_coverage['coverage']
                        
                        print(f"✅ Codecov coverage fetched: {coverage_pct}%")
                        
                        return {
                            "coverage_pct": round(float(coverage_pct), 2),
                            "total_lines": patch_coverage.get('lines', 0),
                            "covered_lines": patch_coverage.get('covered', 0)
                        }
                
                print(f"⏳ Attempt {attempt}/{max_retries}: Coverage not ready yet, retrying in {retry_delay}s...")
                time.sleep(retry_delay)
                
        except urllib.error.HTTPError as e:
            if e.code == 404 and attempt < max_retries:
                print(f"⏳ Attempt {attempt}/{max_retries}: PR not found in Codecov yet, retrying in {retry_delay}s...")
                time.sleep(retry_delay)
            else:
                print(f"❌ HTTP Error {e.code}: {e.reason}")
                return None
        except Exception as e:
            print(f"⚠️  Error fetching from Codecov: {e}")
            if attempt < max_retries:
                time.sleep(retry_delay)
            else:
                return None
    
    print("❌ Timeout: Codecov did not process coverage within timeout period")
    return None


def run_unit_tests_mode(packages: List[str], config: Dict[str, Any], 
                        args: argparse.Namespace) -> int:
    """Execute unit tests mode with coverage analysis.
    
    Args:
        packages: List of changed packages.
        config: Test configuration.
        args: Parsed command-line arguments.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    # If skip_tests is set, only fetch from Codecov
    if not args.skip_tests:
        # Step 2: Determine service domains
        print("\n📊 Step 2: Determining service domains...")
        service_domains = determine_service_domains(packages, config)
        print(f"Service domains: {' '.join(service_domains) if service_domains else 'N/A'}")
        
        # Step 3: Detect goroutines for race detection job
        print("\n🔍 Step 3: Scanning for goroutines...")
        packages_with_goroutines = detect_goroutines_in_packages(packages)
        has_goroutines = len(packages_with_goroutines) > 0
        
        print(f"\n{'✅' if has_goroutines else 'ℹ️ '} Found {len(packages_with_goroutines)} package(s) with goroutines")
        for pkg in packages_with_goroutines:
            print(f"   - {pkg}")
        
        write_github_output({
            "service-domains": ' '.join(service_domains),
            "has-goroutines": "true" if has_goroutines else "false",
            "goroutine-packages": ' '.join(packages_with_goroutines)
        }, args.github_output)
        
        # Step 4: Run tests with coverage
        print("\n📊 Step 4: Running tests with coverage...")
        run_tests_with_coverage(packages, args.output_dir)
        
        print("\n✅ Coverage file generated - will be uploaded to Codecov by workflow")
    
    # Step 5: Fetch coverage from Codecov (after upload by workflow)
    if args.codecov_token and args.repo_slug and args.pr_number:
        print("\n📊 Step 5: Fetching coverage from Codecov...")
        stats = fetch_codecov_coverage(args.repo_slug, args.pr_number, args.codecov_token)
        
        if not stats:
            print("⚠️  Could not fetch coverage from Codecov, skipping enforcement")
            return 0
        
        print(f"\n{'='*60}")
        print(f"Codecov Patch Coverage: {stats['coverage_pct']}%")
        print(f"Total Lines Changed: {stats['total_lines']}")
        print(f"Covered Lines: {stats['covered_lines']}")
        print(f"{'='*60}")
        
        write_github_output({
            "coverage-pct": str(stats['coverage_pct']),
            "total-lines": str(stats['total_lines']),
            "covered-lines": str(stats['covered_lines'])
        }, args.github_output)
        
        # Step 6: Enforce coverage threshold (only for ready PRs)
        if not args.is_draft:
            min_coverage = config.get('coverage_threshold', {}).get('minimum_pct', 60)
            
            if stats['coverage_pct'] < min_coverage:
                print(f"\n❌ ERROR: Coverage {stats['coverage_pct']}% is below minimum threshold {min_coverage}%")
                print("   Please add tests to increase coverage for changed code.")
                return 1
            
            print(f"\n✅ Coverage {stats['coverage_pct']}% meets minimum threshold {min_coverage}%")
            return 0
        
        print("\n📝 Draft PR: Coverage check skipped (informational only)")
        return 0
    else:
        print("\n⚠️  Codecov integration not configured (missing token/repo/pr), skipping coverage enforcement")
        return 0


def run_race_detection_mode(packages: List[str]) -> int:
    """Execute race detection mode on packages with goroutines.
    
    Args:
        packages: List of changed packages.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    # Step 2: Detect packages with goroutines
    print("\n🔍 Step 2: Detecting packages with goroutines...")
    packages_with_goroutines = detect_goroutines_in_packages(packages)
    
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
    
    Handles unit-tests mode (with coverage) and race-detection mode.
    
    Returns:
        Exit code: 0 on success, 1 on failure.
    """
    parser = argparse.ArgumentParser(
        description='Unified PR test orchestrator',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument('--mode', required=True, choices=['unit-tests', 'race-detection'],
                        help='Test mode: unit-tests (with coverage) or race-detection')
    parser.add_argument('--base-ref', required=True,
                        help='Base branch reference (e.g., origin/main)')
    parser.add_argument('--config', default=None,
                        help='Path to PR checks config file (defaults to pr-checks-config.yml in repo root)')
    parser.add_argument('--github-output', default=os.environ.get('GITHUB_OUTPUT'),
                        help='Path to GITHUB_OUTPUT file')
    parser.add_argument('--output-dir', default='coverage',
                        help='Directory for coverage files (unit-tests mode only)')
    parser.add_argument('--is-draft', action='store_true',
                        help='PR is in draft state (informational mode)')
    parser.add_argument('--codecov-token', default=os.environ.get('CODECOV_TOKEN'),
                        help='Codecov API token for fetching coverage results')
    parser.add_argument('--repo-slug', default=None,
                        help='Repository slug in format owner/repo')
    parser.add_argument('--pr-number', default=None,
                        help='Pull request number')
    parser.add_argument('--skip-tests', action='store_true',
                        help='Skip test execution, only fetch coverage from Codecov')
    
    args = parser.parse_args()
    
    config = load_pr_checks_config(args.config)
    
    print("="*60)
    print(f"🚀 PR Test Orchestrator - Mode: {args.mode}")
    if args.is_draft:
        print("📝 Draft PR Mode: Tests are informational only")
    print("="*60)
    
    # Step 1: Identify changed packages
    print("\n📦 Step 1: Identifying changed packages...")
    packages = get_changed_go_packages(args.base_ref)
    
    if not packages:
        print("✅ No Go files changed")
        write_github_output({
            "packages": "",
            "has-changes": "false",
            "has-goroutines": "false"
        }, args.github_output)
        return 0
    
    print(f"📦 Found {len(packages)} changed package(s)")
    for pkg in packages:
        print(f"   - {pkg}")
    
    write_github_output({
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
