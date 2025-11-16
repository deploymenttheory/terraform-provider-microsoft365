#!/usr/bin/env python3
"""
Test runner for nightly acceptance tests.
Usage: ./run-tests.py <type> [service] [coverage-file] [test-output-file]
"""

import os
import sys
import json
import re
import subprocess
from pathlib import Path
from typing import List, Tuple, Optional


def run_command(cmd: List[str], output_file: str) -> int:
    """Run a command and capture output to file and stdout."""
    with open(output_file, 'w') as f:
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


def parse_test_results(output_file: str, category: str, service: str) -> None:
    """Parse test output and create JSON reports of failures and successes."""
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
    """Run provider core tests."""
    print("Running provider core tests...")
    
    cmd = [
        "go", "test", "-v", "-race",
        "-timeout=90m",
        f"-coverprofile={coverage_file}",
        "-covermode=atomic",
        "./internal/client/...",
        "./internal/helpers/...",
        "./internal/provider/...",
        "./internal/utilities/..."
    ]
    
    exit_code = run_command(cmd, test_output_file)
    parse_test_results(test_output_file, "provider-core", "")
    
    return exit_code


def run_service_tests(category: str, service: str, 
                    coverage_file: str, test_output_file: str) -> int:
    """Run tests for a specific service."""
    print(f"Running tests for {category}/{service}...")
    
    test_dir_str = f"./internal/services/{category}/{service}"
    test_dir = Path(test_dir_str)
    
    if not test_dir.exists():
        print(f"Directory not found: {test_dir_str}, creating empty coverage file")
        with open(coverage_file, 'w') as f:
            f.write("mode: atomic\n")
        return 0
    
    test_files = list(test_dir.rglob("*_test.go"))
    test_count = len(test_files)
    
    if test_count == 0:
        print(f"No test files found in {test_dir_str}, creating empty coverage file")
        with open(coverage_file, 'w') as f:
            f.write("mode: atomic\n")
        return 0
    
    print(f"Found {test_count} test files")
    
    # Run tests without -race flag for acceptance tests
    # (prevents OOM on ARM runners)
    # Use string path with /... for recursive package matching
    cmd = [
        "go", "test", "-v",
        "-timeout=90m",
        f"-coverprofile={coverage_file}",
        "-covermode=atomic",
        f"{test_dir_str}/..."
    ]
    
    exit_code = run_command(cmd, test_output_file)
    parse_test_results(test_output_file, category, service)
    
    return exit_code


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

