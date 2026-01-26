#!/usr/bin/env python3
"""Test runner for Terraform provider acceptance tests.

This script runs Terraform provider tests sequentially, one Go package at a time,
to conserve memory on resource-constrained runners. Supports configurable parallelism
and memory management options.

Usage:
    ./run_tests.py <type> [service] [coverage-file] [test-output-file] [options]

Positional Arguments:
    type              Test type: provider-core, resources, datasources, actions, list-resources, or ephemerals
    service           Service name (required for resources/datasources/actions/list-resources/ephemerals)
    coverage-file     Coverage output file (default: coverage.txt)
    test-output-file  Test log output file (default: test-output.log)

Memory & Parallelism Options:
    --max-procs N          GOMAXPROCS - max CPU cores for Go (default: 2)
    --test-parallel N      Tests to run in parallel within a package (default: 1)
    --pkg-parallel N       Packages to build/test in parallel (default: 1)
    --race / --no-race     Enable/disable race detection (auto: on for provider-core)
    --skip-enumeration     Skip test enumeration phase for faster startup
    --force-gc             Force garbage collection between packages (default: enabled)
    --no-force-gc          Disable forced garbage collection

Examples:
    # Low memory mode (8GB runners) - default settings
    ./run_tests.py resources identity_and_access coverage.txt output.log

    # Ultra-conservative mode (minimize memory)
    ./run_tests.py resources identity_and_access coverage.txt output.log \
        --max-procs 1 --test-parallel 1 --pkg-parallel 1 --skip-enumeration --no-race

    # Higher performance mode (16GB+ runners)
    ./run_tests.py resources identity_and_access coverage.txt output.log \
        --max-procs 4 --test-parallel 2 --pkg-parallel 1

    # Show help
    ./run_tests.py --help
"""

import os
import sys
import json
import re
import subprocess
import time
import gc
import argparse
from datetime import datetime
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
    with open(output_file, mode, encoding='utf-8', buffering=8192) as f:  # Use small buffer
        process = subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1  # Line buffered
        )
        
        for line in process.stdout:
            print(line, end='')
            f.write(line)
            f.flush()  # Ensure data is written immediately
        
        process.wait()
        return process.returncode


def discover_test_packages(base_path: Path) -> List[str]:
    """Discover all Go packages that contain test files.

    Args:
        base_path: Base directory path to search for test files.

    Returns:
        Sorted list of package paths (relative to workspace root) containing tests.
    """
    print(f"🔍 [DISCOVERY] Starting test package discovery in: {base_path}")
    
    if not base_path.exists():
        print(f"⚠️  [DISCOVERY] Path does not exist: {base_path}")
        return []
    
    base_path = base_path.resolve()
    cwd = Path.cwd().resolve()
    print(f"🔍 [DISCOVERY] Resolved base_path: {base_path}")
    print(f"🔍 [DISCOVERY] Current working directory: {cwd}")
    
    print("🔍 [DISCOVERY] Searching for *_test.go files...")
    packages = set()
    test_files = list(base_path.rglob("*_test.go"))
    print(f"🔍 [DISCOVERY] Found {len(test_files)} test files")
    
    for idx, test_file in enumerate(test_files, 1):
        if idx % 10 == 0:
            print(f"🔍 [DISCOVERY] Processing test file {idx}/{len(test_files)}...")
        package_dir = test_file.parent.resolve()
        rel_path = f"./{package_dir.relative_to(cwd)}"
        packages.add(rel_path)
    
    sorted_packages = sorted(list(packages))
    print(f"✅ [DISCOVERY] Discovered {len(sorted_packages)} unique packages")
    return sorted_packages


def count_tests_in_package(package_path: str) -> int:
    """Count the number of test functions in a package.

    Args:
        package_path: Path to the Go package.

    Returns:
        Number of test functions found in the package, or 0 if error occurs.
    """
    print(f"  🔢 [COUNT] Enumerating tests in: {package_path}")
    try:
        result = subprocess.run(
            ["go", "test", "-list=.", package_path],
            capture_output=True,
            text=True,
            timeout=30,
            check=False
        )
        
        # Count lines that start with "Test" (test function names)
        test_count = 0
        for line in result.stdout.split('\n'):
            if line.startswith('Test'):
                test_count += 1
        
        if result.returncode != 0 and test_count == 0:
            print(f"  ⚠️  [COUNT] go test -list returned error: {result.stderr[:200]}")
        
        print(f"  ✅ [COUNT] Found {test_count} tests in {package_path}")
        return test_count
    except subprocess.TimeoutExpired:
        print(f"  ⚠️  [COUNT] Timeout counting tests in {package_path}")
        return 0
    except (OSError, ValueError) as e:
        print(f"  ⚠️  [COUNT] Error counting tests in {package_path}: {e}")
        return 0


def print_separator(char: str = "=", length: int = 70) -> None:
    """Print a separator line to stdout.

    Args:
        char: Character to use for the separator line.
        length: Length of the separator line in characters.
    """
    print(char * length)


def parse_test_results(output_file: str, configuration_block_type: str, service: str) -> None:
    """Parse test output and create JSON reports of failures and successes.

    Reads the Go test output file and extracts failed and passed tests,
    creating test-failures.json and test-successes.json files.

    Args:
        output_file: Path to the test output log file.
        configuration_block_type: Test configuration_block_type (e.g., 'provider-core', 'resources', 'datasources', 'actions', 'list-resources', 'ephemerals').
        service: Service name (e.g., 'identity_and_access'), empty string for provider-core.
    """
    failures_file = "test-failures.json"
    successes_file = "test-successes.json"
    failures = []
    successes = []
    
    with open(output_file, 'r', encoding='utf-8') as f:
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
            "configuration_block_type": configuration_block_type,
            "service": service,
            "context": context
        })
    
    pass_pattern = re.compile(r'^--- PASS: (\S+)', re.MULTILINE)
    for match in pass_pattern.finditer(content):
        test_name = match.group(1)
        
        successes.append({
            "test_name": test_name,
            "configuration_block_type": configuration_block_type,
            "service": service
        })
    
    with open(failures_file, 'w', encoding='utf-8') as f:
        json.dump(failures, f, indent=2)
    
    with open(successes_file, 'w', encoding='utf-8') as f:
        json.dump(successes, f, indent=2)
    
    print(f"✅ Test results: {len(failures)} failures, {len(successes)} successes")


def run_provider_core_tests(coverage_file: str, test_output_file: str, 
                            max_procs: int = 2, test_parallel: int = 1, 
                            pkg_parallel: int = 1, use_race: bool = True,
                            skip_enumeration: bool = False, force_gc: bool = True) -> int:
    """Run provider core tests sequentially, one package at a time to conserve memory.

    Discovers all test packages in core directories (client, helpers, provider, utilities),
    runs tests package-by-package with race detection, and collects coverage data.

    Args:
        coverage_file: Path where merged coverage data will be written.
        test_output_file: Path where test output logs will be written.
        max_procs: Maximum number of CPU cores for Go to use (GOMAXPROCS).
        test_parallel: Number of tests to run in parallel within a package.
        pkg_parallel: Number of packages to build/test in parallel.
        use_race: Whether to enable race detection (-race flag).

    Returns:
        0 if all tests passed, 1 if any test failed.
    """
    print("\n" + "="*70)
    print("🔍 [START] Running provider core tests")
    print("="*70 + "\n")
    
    # Define core directories to test
    core_dirs = [
        "./internal/client",
        "./internal/helpers",
        "./internal/provider",
        "./internal/utilities"
    ]
    
    print(f"📂 [CONFIG] Core directories to test: {len(core_dirs)}")
    for core_dir in core_dirs:
        print(f"   - {core_dir}")
    print(f"⚙️  [CONFIG] GOMAXPROCS: {max_procs}")
    print(f"⚙️  [CONFIG] Test parallel: {test_parallel}")
    print(f"⚙️  [CONFIG] Package parallel: {pkg_parallel}")
    print(f"⚙️  [CONFIG] Race detection: {use_race}")
    print()
    
    # Discover all test packages across core directories
    all_packages = []
    for idx, core_dir in enumerate(core_dirs, 1):
        print(f"🔍 [DISCOVERY] Processing directory {idx}/{len(core_dirs)}: {core_dir}")
        packages = discover_test_packages(Path(core_dir))
        all_packages.extend(packages)
        print(f"✅ [DISCOVERY] Found {len(packages)} packages in {core_dir}\n")
    
    if not all_packages:
        print("⚠️  No test packages found in provider core, creating empty coverage file")
        with open(coverage_file, 'w', encoding='utf-8') as f:
            f.write("mode: atomic\n")
        return 0
    
    # Count tests per package (optional)
    package_test_counts: Dict[str, int] = {}
    total_tests = 0
    
    if not skip_enumeration:
        print("="*70)
        print(f"📊 [ENUMERATE] Starting test enumeration for {len(all_packages)} package(s)")
        print("="*70 + "\n")
        
        enumerate_start = time.time()
        for idx, pkg in enumerate(all_packages, 1):
            pkg_start = time.time()
            print(f"📊 [ENUMERATE] Package {idx}/{len(all_packages)}: {pkg}")
            count = count_tests_in_package(pkg)
            package_test_counts[pkg] = count
            total_tests += count
            pkg_elapsed = time.time() - pkg_start
            print(f"⏱️  [TIMING] Package enumeration took {pkg_elapsed:.2f}s\n")
        enumerate_elapsed = time.time() - enumerate_start
        print(f"⏱️  [TIMING] Total enumeration time: {enumerate_elapsed:.2f}s\n")
    else:
        print("⏭️  [ENUMERATE] Skipping test enumeration (--skip-enumeration flag)\n")
        for pkg in all_packages:
            package_test_counts[pkg] = 0
    
    print_separator("=")
    print("📋 Test Discovery Summary for provider-core")
    print_separator("=")
    print(f"Total Packages: {len(all_packages)}")
    if not skip_enumeration:
        print(f"Total Tests: {total_tests}")
    print_separator("-")
    
    for pkg in all_packages:
        count = package_test_counts[pkg]
        if skip_enumeration:
            print(f"  📦 {pkg}")
        else:
            print(f"  📦 {pkg:<50} {count:>4} test(s)")
    
    print_separator("=")
    print()
    
    # Run tests package by package
    print("="*70)
    print(f"🚀 [EXECUTE] Starting sequential execution ({len(all_packages)} package(s))")
    print("="*70 + "\n")
    
    has_failures = False
    
    # Initialize coverage file with mode line
    print(f"📝 [SETUP] Initializing coverage file: {coverage_file}")
    with open(coverage_file, 'w', encoding='utf-8') as f:
        f.write("mode: atomic\n")
    print()
    
    for idx, pkg in enumerate(all_packages, 1):
        pkg_exec_start = time.time()
        test_count = package_test_counts[pkg]
        
        print_separator("-", 70)
        print(f"📦 Package {idx}/{len(all_packages)}: {pkg}")
        print(f"   Tests: {test_count}")
        print(f"⏱️  [TIMING] Started at: {datetime.now().strftime('%H:%M:%S')}")
        print_separator("-", 70)
        
        # Create temporary coverage file for this package
        temp_coverage = f"{coverage_file}.tmp"
        
        # Build command with configurable parallelism
        cmd = ["go", "test", "-v"]
        
        if use_race:
            cmd.append("-race")
        
        cmd.extend([
            "-timeout=90m",
            f"-p={pkg_parallel}",
            f"-parallel={test_parallel}",
            f"-coverprofile={temp_coverage}",
            "-covermode=atomic",
            pkg
        ])
        
        race_flag = "-race" if use_race else ""
        print(f"▶️  [RUN] Executing: go test {race_flag} -p={pkg_parallel} -parallel={test_parallel} {pkg}")
        print(f"📄 [RUN] Output mode: {'append' if idx > 1 else 'write'} to {test_output_file}\n")
        
        # Append to the output file for each package
        exit_code = run_command(cmd, test_output_file, append=(idx > 1))
        
        pkg_exec_elapsed = time.time() - pkg_exec_start
        print(f"⏱️  [TIMING] Package execution took {pkg_exec_elapsed:.2f}s")
        
        print("\n🔍 [COVERAGE] Processing coverage for package...")
        # Append coverage data (skip mode line)
        if Path(temp_coverage).exists():
            print(f"✅ [COVERAGE] Found coverage file: {temp_coverage}")
            with open(temp_coverage, 'r', encoding='utf-8') as tmp_f:
                lines = tmp_f.readlines()
                print(f"📊 [COVERAGE] Coverage file has {len(lines)} lines")
                with open(coverage_file, 'a', encoding='utf-8') as cov_f:
                    for line in lines:
                        if not line.startswith('mode:'):
                            cov_f.write(line)
            Path(temp_coverage).unlink()
            print("🗑️  [COVERAGE] Deleted temporary coverage file")
        else:
            print("⚠️  [COVERAGE] No coverage file generated")
        
        if exit_code != 0:
            has_failures = True
            print(f"\n❌ [RESULT] Package {pkg} completed with failures (exit code: {exit_code})")
        else:
            print(f"\n✅ [RESULT] Package {pkg} completed successfully")
        
        # Force garbage collection to free memory
        if force_gc:
            print("🗑️  [MEMORY] Running garbage collection...")
            gc.collect()
        print()
    
    # Parse all test results from the combined output
    print("📝 Parsing test results...")
    parse_test_results(test_output_file, "provider-core", "")
    
    print_separator("=")
    print("🏁 Sequential execution complete for provider-core")
    print_separator("=")
    
    return 1 if has_failures else 0


def run_service_tests(configuration_block_type: str, service: str,
                    coverage_file: str, test_output_file: str,
                    max_procs: int = 2, test_parallel: int = 1,
                    pkg_parallel: int = 1, use_race: bool = False,
                    skip_enumeration: bool = False, force_gc: bool = True,
                    shard_resources: str = "") -> int:
    """Run tests for a specific service sequentially, one package at a time to conserve memory.

    Discovers all test packages in the service directory, runs tests package-by-package,
    and collects coverage data.

    Args:
        configuration_block_type: Test configuration_block_type ('resources', 'datasources', 'actions', 'list-resources', or 'ephemerals').
        service: Service name (e.g., 'identity_and_access', 'device_management').
        coverage_file: Path where merged coverage data will be written.
        test_output_file: Path where test output logs will be written.
        max_procs: Maximum number of CPU cores for Go to use (GOMAXPROCS).
        test_parallel: Number of tests to run in parallel within a package.
        pkg_parallel: Number of packages to build/test in parallel.
        use_race: Whether to enable race detection (-race flag).
        skip_enumeration: Whether to skip test enumeration phase.
        force_gc: Whether to force garbage collection between packages.
        shard_resources: Comma-separated list of resource paths in this shard.

    Returns:
        0 if all tests passed, 1 if any test failed.
    """
    print("\n" + "="*70)
    print(f"🔍 [START] Running {configuration_block_type}/{service} tests")
    print("="*70 + "\n")
    
    test_dir_str = f"./internal/services/{configuration_block_type}/{service}"
    test_dir = Path(test_dir_str)
    print(f"📂 [CONFIG] Test directory: {test_dir_str}")
    print(f"⚙️  [CONFIG] GOMAXPROCS: {max_procs}")
    print(f"⚙️  [CONFIG] Test parallel: {test_parallel}")
    print(f"⚙️  [CONFIG] Package parallel: {pkg_parallel}")
    print(f"⚙️  [CONFIG] Race detection: {use_race}")
    
    if not test_dir.exists():
        print(f"⚠️  Directory not found: {test_dir_str}, creating empty coverage file")
        with open(coverage_file, 'w', encoding='utf-8') as f:
            f.write("mode: atomic\n")
        return 0
    
    # Discover all test packages
    print("\n🔍 [DISCOVERY] Starting package discovery...")
    test_packages = discover_test_packages(test_dir)

    # Why shard? Parallel test execution by distributing resources across jobs
    if shard_resources:
        print(f"🔍 [SHARD] Testing resources in this shard: {shard_resources}")
        shard_resource_list = [
            r.strip()
            for r in shard_resources.split(',')
            if r.strip()
        ]
        print(f"🔍 [SHARD] This shard contains {len(shard_resource_list)} resource(s)")

        shard_packages = []
        for resource_path in shard_resource_list:
            # Why append test_dir_str? discover_test_packages returns
            # full paths like ./internal/services/resources/service/graph_beta/name
            expected_pkg = f"{test_dir_str}/{resource_path}"

            if expected_pkg in test_packages:
                shard_packages.append(expected_pkg)
                print(f"✅ [SHARD] Matched resource: {resource_path}")
            else:
                print(f"⚠️  [SHARD] Resource not found: {resource_path}")
                print(f"   Looking for: {expected_pkg}")

        test_packages = shard_packages
        print(f"✅ [SHARD] Testing {len(test_packages)} package(s) in this shard")

    if not test_packages:
        print(f"⚠️  [DISCOVERY] No test packages found in {test_dir_str}")
        print("✅ [COMPLETE] Creating empty coverage file and exiting")
        with open(coverage_file, 'w', encoding='utf-8') as f:
            f.write("mode: atomic\n")
        return 0
    
    # Count tests per package (optional)
    package_test_counts: Dict[str, int] = {}
    total_tests = 0
    
    if not skip_enumeration:
        print("\n" + "="*70)
        print(f"📊 [ENUMERATE] Starting test enumeration for {len(test_packages)} package(s)")
        print("="*70 + "\n")
        
        enumerate_start = time.time()
        for idx, pkg in enumerate(test_packages, 1):
            pkg_start = time.time()
            print(f"📊 [ENUMERATE] Package {idx}/{len(test_packages)}: {pkg}")
            count = count_tests_in_package(pkg)
            package_test_counts[pkg] = count
            total_tests += count
            pkg_elapsed = time.time() - pkg_start
            print(f"⏱️  [TIMING] Package enumeration took {pkg_elapsed:.2f}s\n")
        enumerate_elapsed = time.time() - enumerate_start
        print(f"⏱️  [TIMING] Total enumeration time: {enumerate_elapsed:.2f}s\n")
    else:
        print("\n⏭️  [ENUMERATE] Skipping test enumeration (--skip-enumeration flag)\n")
        # Initialize with zero counts
        for pkg in test_packages:
            package_test_counts[pkg] = 0
    
    print_separator("=")
    print(f"📋 Test Discovery Summary for {configuration_block_type}/{service}")
    print_separator("=")
    print(f"Total Packages: {len(test_packages)}")
    if not skip_enumeration:
        print(f"Total Tests: {total_tests}")
    print_separator("-")
    
    for pkg in test_packages:
        count = package_test_counts[pkg]
        # Show package relative to service directory for readability
        rel_pkg = pkg.replace(test_dir_str, "").lstrip("/")
        if skip_enumeration:
            print(f"  📦 {rel_pkg or '.'}")
        else:
            print(f"  📦 {rel_pkg or '.':<50} {count:>4} test(s)")
    
    print_separator("=")
    print()
    
    # Run tests package by package
    print("="*70)
    print(f"🚀 [EXECUTE] Starting sequential execution ({len(test_packages)} package(s))")
    print("="*70 + "\n")
    
    has_failures = False
    
    # Initialize coverage file with mode line
    print(f"📝 [SETUP] Initializing coverage file: {coverage_file}")
    with open(coverage_file, 'w', encoding='utf-8') as f:
        f.write("mode: atomic\n")
    print()
    
    for idx, pkg in enumerate(test_packages, 1):
        pkg_exec_start = time.time()
        rel_pkg = pkg.replace(test_dir_str, "").lstrip("/") or "."
        test_count = package_test_counts[pkg]
        
        print_separator("-", 70)
        print(f"📦 Package {idx}/{len(test_packages)}: {rel_pkg}")
        print(f"   Tests: {test_count}")
        print(f"⏱️  [TIMING] Started at: {datetime.now().strftime('%H:%M:%S')}")
        print_separator("-", 70)
        
        # Create temporary coverage file for this package
        temp_coverage = f"{coverage_file}.tmp"
        
        # Build command with configurable parallelism
        cmd = ["go", "test", "-v"]
        
        if use_race:
            cmd.append("-race")
        
        cmd.extend([
            "-timeout=90m",
            f"-p={pkg_parallel}",
            f"-parallel={test_parallel}",
            f"-coverprofile={temp_coverage}",
            "-covermode=atomic",
            pkg
        ])
        
        race_flag = "-race" if use_race else ""
        print(f"▶️  [RUN] Executing: go test {race_flag} -p={pkg_parallel} -parallel={test_parallel} {pkg}")
        print(f"📄 [RUN] Output mode: {'append' if idx > 1 else 'write'} to {test_output_file}\n")
        
        # Append to the output file for each package
        exit_code = run_command(cmd, test_output_file, append=(idx > 1))
        
        pkg_exec_elapsed = time.time() - pkg_exec_start
        print(f"⏱️  [TIMING] Package execution took {pkg_exec_elapsed:.2f}s")
        
        print("\n🔍 [COVERAGE] Processing coverage for package...")
        # Append coverage data (skip mode line)
        if Path(temp_coverage).exists():
            print(f"✅ [COVERAGE] Found coverage file: {temp_coverage}")
            with open(temp_coverage, 'r', encoding='utf-8') as tmp_f:
                lines = tmp_f.readlines()
                print(f"📊 [COVERAGE] Coverage file has {len(lines)} lines")
                with open(coverage_file, 'a', encoding='utf-8') as cov_f:
                    for line in lines:
                        if not line.startswith('mode:'):
                            cov_f.write(line)
            Path(temp_coverage).unlink()
            print("🗑️  [COVERAGE] Deleted temporary coverage file")
        else:
            print("⚠️  [COVERAGE] No coverage file generated")
        
        if exit_code != 0:
            has_failures = True
            print(f"\n❌ [RESULT] Package {rel_pkg} completed with failures (exit code: {exit_code})")
        else:
            print(f"\n✅ [RESULT] Package {rel_pkg} completed successfully")
        
        # Force garbage collection to free memory
        if force_gc:
            print("🗑️  [MEMORY] Running garbage collection...")
            gc.collect()
        print()
    
    # Parse all test results from the combined output
    print("📝 Parsing test results...")
    parse_test_results(test_output_file, configuration_block_type, service)
    
    print_separator("=")
    print(f"🏁 Sequential execution complete for {configuration_block_type}/{service}")
    print_separator("=")
    
    return 1 if has_failures else 0


def main():
    parser = argparse.ArgumentParser(
        description='Run Terraform Provider tests with configurable parallelism and memory management',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Run with default settings (low memory mode)
  ./run_tests.py resources identity_and_access coverage.txt output.log
  
  # Run with more parallelism (requires more memory)
  ./run_tests.py resources identity_and_access coverage.txt output.log --max-procs 4 --test-parallel 4
  
  # Disable race detection for faster execution
  ./run_tests.py resources identity_and_access coverage.txt output.log --no-race
        """
    )
    
    # Positional arguments
    parser.add_argument('type', choices=['provider-core', 'resources', 'datasources', 'actions', 'list-resources', 'ephemerals'],
                       help='Type of tests to run')
    parser.add_argument('service', nargs='?', default='',
                       help='Service name (required for resources/datasources/actions/list-resources/ephemerals)')
    parser.add_argument('coverage_file', nargs='?', default='coverage.txt',
                       help='Output file for coverage data (default: coverage.txt)')
    parser.add_argument('output_file', nargs='?', default='test-output.log',
                       help='Output file for test logs (default: test-output.log)')
    
    # Optional arguments for memory/parallelism control
    parser.add_argument('--max-procs', type=int, default=2,
                       help='GOMAXPROCS value - max CPU cores for Go (default: 2)')
    parser.add_argument('--test-parallel', type=int, default=1,
                       help='Number of tests to run in parallel within a package (default: 1)')
    parser.add_argument('--pkg-parallel', type=int, default=1,
                       help='Number of packages to build/test in parallel (default: 1)')
    parser.add_argument('--race', dest='use_race', action='store_true', default=None,
                       help='Enable race detection (default for provider-core)')
    parser.add_argument('--no-race', dest='use_race', action='store_false',
                       help='Disable race detection (default for resources/datasources)')
    parser.add_argument('--skip-enumeration', action='store_true', default=False,
                       help='Skip test enumeration phase (faster startup, no test counts)')
    parser.add_argument('--force-gc', action='store_true', default=True,
                       help='Force garbage collection between packages (default: enabled)')
    parser.add_argument('--no-force-gc', dest='force_gc', action='store_false',
                       help='Disable forced garbage collection')
    parser.add_argument('--shard-resources', type=str, default='',
                       help='Comma-separated list of resource paths in this shard (e.g., "graph_beta/resource1,graph_v1.0/resource2")')

    args = parser.parse_args()
    
    # Set default for use_race if not specified
    if args.use_race is None:
        args.use_race = (args.type == 'provider-core')
    
    print("\n" + "="*70)
    print("🚀 [MAIN] Terraform Provider Test Runner Started")
    print("="*70)
    print(f"📋 [MAIN] Command: {' '.join(sys.argv)}")
    print(f"📂 [MAIN] Working directory: {Path.cwd()}")
    print("="*70 + "\n")
    
    # Set memory-friendly environment variables for Go
    if "GOMAXPROCS" not in os.environ:
        os.environ["GOMAXPROCS"] = str(args.max_procs)
        print(f"⚙️  [MEMORY] Set GOMAXPROCS={args.max_procs}")
    else:
        print(f"⚙️  [MEMORY] GOMAXPROCS already set to {os.environ['GOMAXPROCS']}")
    
    # Set GODEBUG to reduce memory usage
    os.environ["GODEBUG"] = "gctrace=0"
    print("⚙️  [MEMORY] Set GODEBUG to optimize garbage collection\n")
    
    print(f"⚙️  [CONFIG] Test type: {args.type}")
    print(f"⚙️  [CONFIG] Service: {args.service if args.service else 'N/A'}")
    print(f"⚙️  [CONFIG] Coverage file: {args.coverage_file}")
    print(f"⚙️  [CONFIG] Test output file: {args.output_file}")
    print(f"⚙️  [CONFIG] Max procs (GOMAXPROCS): {args.max_procs}")
    print(f"⚙️  [CONFIG] Test parallel (-parallel): {args.test_parallel}")
    print(f"⚙️  [CONFIG] Package parallel (-p): {args.pkg_parallel}")
    print(f"⚙️  [CONFIG] Race detection (-race): {'enabled' if args.use_race else 'disabled'}")
    
    # Calculate estimated memory usage
    base_memory = 500  # Base Go runtime memory in MB
    test_memory = args.test_parallel * 200  # ~200MB per parallel test
    race_memory = 1500 if args.use_race else 0  # Race detector overhead
    estimated_memory = base_memory + test_memory + race_memory
    
    print(f"⚙️  [CONFIG] Estimated memory usage: ~{estimated_memory}MB")
    if estimated_memory > 7000:
        print("⚠️  [WARNING] Estimated memory usage exceeds 7GB - may cause OOM on 8GB runners!")
    print()
    
    if os.environ.get("SKIP_TESTS", "false") == "true":
        print("⏭️  Skipping tests - no credentials configured")
        with open(args.coverage_file, 'w', encoding='utf-8') as f:
            f.write("mode: atomic\n")
        sys.exit(0)
    
    if args.type == "provider-core":
        exit_code = run_provider_core_tests(
            args.coverage_file, 
            args.output_file,
            args.max_procs,
            args.test_parallel,
            args.pkg_parallel,
            args.use_race,
            args.skip_enumeration,
            args.force_gc
        )
    elif args.type in ["resources", "datasources", "actions", "list-resources", "ephemerals"]:
        if not args.service:
            print(f"Error: service name required for {args.type} tests", 
                  file=sys.stderr)
            sys.exit(1)
        exit_code = run_service_tests(
            args.type,
            args.service,
            args.coverage_file,
            args.output_file,
            args.max_procs,
            args.test_parallel,
            args.pkg_parallel,
            args.use_race,
            args.skip_enumeration,
            args.force_gc,
            args.shard_resources
        )
    else:
        print(f"Error: unknown test type: {args.type}", file=sys.stderr)
        print("Valid types: provider-core, resources, datasources, actions, list-resources, ephemerals", file=sys.stderr)
        sys.exit(1)
    
    print("Tests completed")
    sys.exit(exit_code)


if __name__ == "__main__":
    main()

