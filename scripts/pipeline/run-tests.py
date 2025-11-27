#!/usr/bin/env python3
"""Test runner for nightly acceptance tests.

This script runs Terraform provider tests sequentially, one package at a time,
to conserve memory and provide better progress visibility.

Usage:
    ./run-tests.py <type> [service] [coverage-file] [test-output-file]

Types:
    provider-core: Core provider tests (client, helpers, provider, utilities)
    resources: Resource tests for a specific service
    datasources: Datasource tests for a specific service
"""

import os
import sys
import json
import re
import subprocess
from pathlib import Path
from typing import List, Dict


def run_command(cmd: List[str], output_file: str, append: bool = False) -> int:
    """Run a command and capture output to file and stdout.

    Args:
        cmd: Command and arguments to execute as a list.
        output_file: Path to file where output will be written.
        append: If True, append to output_file; if False, overwrite it.

    Returns:
        The exit code of the command.
    """
    mode = 'a' if append else 'w'
    with open(output_file, mode) as f:
        process = subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True
        )
        
        for line in process.stdout:
            print(line, end='')
            f.write(line)
        
        process.wait()
        return process.returncode


def discover_test_packages(base_path: Path) -> List[str]:
    """Discover all Go packages that contain test files.

    Args:
        base_path: Base directory path to search for test files.

    Returns:
        Sorted list of package paths (relative to workspace root) containing tests.
    """
    if not base_path.exists():
        return []
    
    base_path = base_path.resolve()
    cwd = Path.cwd().resolve()
    
    packages = set()
    test_files = list(base_path.rglob("*_test.go"))
    
    for test_file in test_files:
        # Get the package directory (parent of the test file)
        package_dir = test_file.parent.resolve()
        # Convert to relative path from workspace root
        rel_path = f"./{package_dir.relative_to(cwd)}"
        packages.add(rel_path)
    
    return sorted(list(packages))


def count_tests_in_package(package_path: str) -> int:
    """Count the number of test functions in a package.

    Args:
        package_path: Path to the Go package.

    Returns:
        Number of test functions found in the package, or 0 if error occurs.
    """
    try:
        result = subprocess.run(
            ["go", "test", "-list", ".", package_path],
            capture_output=True,
            text=True,
            timeout=30
        )
        
        # Count lines that start with "Test" (test function names)
        test_count = 0
        for line in result.stdout.split('\n'):
            if line.startswith('Test'):
                test_count += 1
        
        return test_count
    except Exception:
        return 0


def print_separator(char: str = "=", length: int = 70) -> None:
    """Print a separator line to stdout.

    Args:
        char: Character to use for the separator line.
        length: Length of the separator line in characters.
    """
    print(char * length)


def parse_test_results(output_file: str, category: str, service: str) -> None:
    """Parse test output and create JSON reports of failures and successes.

    Reads the Go test output file and extracts failed and passed tests,
    creating test-failures.json and test-successes.json files.

    Args:
        output_file: Path to the test output log file.
        category: Test category (e.g., 'provider-core', 'resources', 'datasources').
        service: Service name (e.g., 'identity_and_access'), empty string for provider-core.
    """
    failures_file = "test-failures.json"
    successes_file = "test-successes.json"
    failures = []
    successes = []
    
    with open(output_file, 'r') as f:
        content = f.read()
    
    fail_pattern = re.compile(r'^--- FAIL: (\S+)', re.MULTILINE)
    for match in fail_pattern.finditer(content):
        test_name = match.group(1)
        
        run_pattern = re.compile(rf'^=== RUN\s+{re.escape(test_name)}', re.MULTILINE)
        run_match = None
        
        search_start = max(0, match.start() - 50000)  # Look back up to 50KB
        search_content = content[search_start:match.start()]
        
        for run_match_candidate in run_pattern.finditer(search_content):
            run_match = run_match_candidate
        
        if run_match:
            context_start = search_start + run_match.end()
            context_end = match.start()
            full_context = content[context_start:context_end].strip()
            
            lines = full_context.split('\n')
            error_start_idx = None
            
            for idx, line in enumerate(lines):
                if '[DEBUG]' in line or '[INFO]' in line:
                    continue
                if any(indicator in line for indicator in ['.go:', 'Error:', 'panic:', '    ', '\t']):
                    error_start_idx = idx
                    break
            
            if error_start_idx is not None:
                context = '\n'.join(lines[error_start_idx:])
            else:
                context = full_context
            
            if len(context) > 1000:
                context = context[:1000] + "\n... (truncated)"
        else:
            context_start = max(0, match.start() - 500)
            context = content[context_start:match.start()].strip()
            if len(context) > 500:
                context = "... " + context[-500:]
        
        failures.append({
            "test_name": test_name,
            "category": category,
            "service": service,
            "context": context
        })
    
    pass_pattern = re.compile(r'^--- PASS: (\S+)', re.MULTILINE)
    for match in pass_pattern.finditer(content):
        test_name = match.group(1)
        
        successes.append({
            "test_name": test_name,
            "category": category,
            "service": service
        })
    
    with open(failures_file, 'w') as f:
        json.dump(failures, f, indent=2)
    
    with open(successes_file, 'w') as f:
        json.dump(successes, f, indent=2)
    
    print(f"✅ Test results: {len(failures)} failures, {len(successes)} successes")


def run_provider_core_tests(coverage_file: str, test_output_file: str) -> int:
    """Run provider core tests sequentially, one package at a time to conserve memory.

    Discovers all test packages in core directories (client, helpers, provider, utilities),
    runs tests package-by-package with race detection, and collects coverage data.

    Args:
        coverage_file: Path where merged coverage data will be written.
        test_output_file: Path where test output logs will be written.

    Returns:
        0 if all tests passed, 1 if any test failed.
    """
    print("\n🔍 Discovering provider core test packages...\n")
    
    # Define core directories to test
    core_dirs = [
        "./internal/client",
        "./internal/helpers",
        "./internal/provider",
        "./internal/utilities"
    ]
    
    # Discover all test packages across core directories
    all_packages = []
    for core_dir in core_dirs:
        packages = discover_test_packages(Path(core_dir))
        all_packages.extend(packages)
    
    if not all_packages:
        print("⚠️  No test packages found in provider core, creating empty coverage file")
        with open(coverage_file, 'w') as f:
            f.write("mode: atomic\n")
        return 0
    
    # Count tests per package
    print(f"📊 Enumerating tests in {len(all_packages)} package(s)...\n")
    package_test_counts: Dict[str, int] = {}
    total_tests = 0
    
    for pkg in all_packages:
        count = count_tests_in_package(pkg)
        package_test_counts[pkg] = count
        total_tests += count
    
    # Display summary
    print_separator("=")
    print("📋 Test Discovery Summary for provider-core")
    print_separator("=")
    print(f"Total Packages: {len(all_packages)}")
    print(f"Total Tests: {total_tests}")
    print_separator("-")
    
    for pkg in all_packages:
        count = package_test_counts[pkg]
        print(f"  📦 {pkg:<50} {count:>4} test(s)")
    
    print_separator("=")
    print()
    
    # Run tests package by package
    print(f"🚀 Starting sequential execution ({len(all_packages)} package(s), one at a time)\n")
    
    has_failures = False
    
    # Initialize coverage file with mode line
    with open(coverage_file, 'w') as f:
        f.write("mode: atomic\n")
    
    for idx, pkg in enumerate(all_packages, 1):
        test_count = package_test_counts[pkg]
        
        print_separator("-", 70)
        print(f"📦 Package {idx}/{len(all_packages)}: {pkg}")
        print(f"   Tests: {test_count}")
        print_separator("-", 70)
        
        # Create temporary coverage file for this package
        temp_coverage = f"{coverage_file}.tmp"
        
        # Run tests for this package with -race flag (smaller scope than before)
        cmd = [
            "go", "test", "-v", "-race",
            "-timeout=90m",
            f"-coverprofile={temp_coverage}",
            "-covermode=atomic",
            pkg
        ]
        
        print(f"▶️  Running: go test -race {pkg}\n")
        
        # Append to the output file for each package
        exit_code = run_command(cmd, test_output_file, append=(idx > 1))
        
        # Append coverage data (skip mode line)
        if Path(temp_coverage).exists():
            with open(temp_coverage, 'r') as tmp_f:
                lines = tmp_f.readlines()
                with open(coverage_file, 'a') as cov_f:
                    for line in lines:
                        if not line.startswith('mode:'):
                            cov_f.write(line)
            Path(temp_coverage).unlink()
        
        if exit_code != 0:
            has_failures = True
            print(f"\n❌ Package {pkg} completed with failures (exit code: {exit_code})")
        else:
            print(f"\n✅ Package {pkg} completed successfully")
        
        print()
    
    # Parse all test results from the combined output
    print("📝 Parsing test results...")
    parse_test_results(test_output_file, "provider-core", "")
    
    print_separator("=")
    print("🏁 Sequential execution complete for provider-core")
    print_separator("=")
    
    return 1 if has_failures else 0


def run_service_tests(category: str, service: str, 
                    coverage_file: str, test_output_file: str) -> int:
    """Run tests for a specific service sequentially, one package at a time to conserve memory.

    Discovers all test packages in the service directory, runs tests package-by-package,
    and collects coverage data.

    Args:
        category: Test category ('resources' or 'datasources').
        service: Service name (e.g., 'identity_and_access', 'device_management').
        coverage_file: Path where merged coverage data will be written.
        test_output_file: Path where test output logs will be written.

    Returns:
        0 if all tests passed, 1 if any test failed.
    """
    print(f"\n🔍 Discovering test packages for {category}/{service}...")
    
    test_dir_str = f"./internal/services/{category}/{service}"
    test_dir = Path(test_dir_str)
    
    if not test_dir.exists():
        print(f"⚠️  Directory not found: {test_dir_str}, creating empty coverage file")
        with open(coverage_file, 'w') as f:
            f.write("mode: atomic\n")
        return 0
    
    # Discover all test packages
    test_packages = discover_test_packages(test_dir)
    
    if not test_packages:
        print(f"⚠️  No test packages found in {test_dir_str}, creating empty coverage file")
        with open(coverage_file, 'w') as f:
            f.write("mode: atomic\n")
        return 0
    
    # Count tests per package
    print(f"\n📊 Enumerating tests in {len(test_packages)} package(s)...\n")
    package_test_counts: Dict[str, int] = {}
    total_tests = 0
    
    for pkg in test_packages:
        count = count_tests_in_package(pkg)
        package_test_counts[pkg] = count
        total_tests += count
    
    # Display summary
    print_separator("=")
    print(f"📋 Test Discovery Summary for {category}/{service}")
    print_separator("=")
    print(f"Total Packages: {len(test_packages)}")
    print(f"Total Tests: {total_tests}")
    print_separator("-")
    
    for pkg in test_packages:
        count = package_test_counts[pkg]
        # Show package relative to service directory for readability
        rel_pkg = pkg.replace(test_dir_str, "").lstrip("/")
        print(f"  📦 {rel_pkg or '.':<50} {count:>4} test(s)")
    
    print_separator("=")
    print()
    
    # Run tests package by package
    print(f"🚀 Starting sequential execution ({len(test_packages)} package(s), one at a time)\n")
    
    has_failures = False
    
    # Initialize coverage file with mode line
    with open(coverage_file, 'w') as f:
        f.write("mode: atomic\n")
    
    for idx, pkg in enumerate(test_packages, 1):
        rel_pkg = pkg.replace(test_dir_str, "").lstrip("/") or "."
        test_count = package_test_counts[pkg]
        
        print_separator("-", 70)
        print(f"📦 Package {idx}/{len(test_packages)}: {rel_pkg}")
        print(f"   Tests: {test_count}")
        print_separator("-", 70)
        
        # Create temporary coverage file for this package
        temp_coverage = f"{coverage_file}.tmp"
        
        # Run tests for this package only
        cmd = [
            "go", "test", "-v",
            "-timeout=90m",
            f"-coverprofile={temp_coverage}",
            "-covermode=atomic",
            pkg
        ]
        
        print(f"▶️  Running: go test {pkg}\n")
        
        # Append to the output file for each package
        exit_code = run_command(cmd, test_output_file, append=(idx > 1))
        
        # Append coverage data (skip mode line)
        if Path(temp_coverage).exists():
            with open(temp_coverage, 'r') as tmp_f:
                lines = tmp_f.readlines()
                with open(coverage_file, 'a') as cov_f:
                    for line in lines:
                        if not line.startswith('mode:'):
                            cov_f.write(line)
            Path(temp_coverage).unlink()
        
        if exit_code != 0:
            has_failures = True
            print(f"\n❌ Package {rel_pkg} completed with failures (exit code: {exit_code})")
        else:
            print(f"\n✅ Package {rel_pkg} completed successfully")
        
        print()
    
    # Parse all test results from the combined output
    print("📝 Parsing test results...")
    parse_test_results(test_output_file, category, service)
    
    print_separator("=")
    print(f"🏁 Sequential execution complete for {category}/{service}")
    print_separator("=")
    
    return 1 if has_failures else 0


def main():
    if len(sys.argv) < 2:
        print("Usage: run-tests.py <type> [service] [coverage-file] [test-output-file]", 
              file=sys.stderr)
        print("Types: provider-core, resources, datasources", file=sys.stderr)
        sys.exit(1)
    
    test_type = sys.argv[1]
    service = sys.argv[2] if len(sys.argv) > 2 else ""
    coverage_file = sys.argv[3] if len(sys.argv) > 3 else "coverage.txt"
    test_output_file = sys.argv[4] if len(sys.argv) > 4 else "test-output.log"
    
    if os.environ.get("SKIP_TESTS", "false") == "true":
        print("⏭️  Skipping tests - no credentials configured")
        with open(coverage_file, 'w') as f:
            f.write("mode: atomic\n")
        sys.exit(0)
    
    if test_type == "provider-core":
        exit_code = run_provider_core_tests(coverage_file, test_output_file)
    elif test_type in ["resources", "datasources"]:
        if not service:
            print(f"Error: service name required for {test_type} tests", 
                  file=sys.stderr)
            sys.exit(1)
        exit_code = run_service_tests(test_type, service, coverage_file, test_output_file)
    else:
        print(f"Error: unknown test type: {test_type}", file=sys.stderr)
        print("Valid types: provider-core, resources, datasources", file=sys.stderr)
        sys.exit(1)
    
    print("Tests completed")
    sys.exit(exit_code)


if __name__ == "__main__":
    main()

