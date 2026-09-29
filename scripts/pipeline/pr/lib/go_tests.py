#!/usr/bin/env python3
"""Go test execution utilities.

Provides functions for running Go unit tests and race detection.
"""

import os
import subprocess
from pathlib import Path
from typing import List


def run_unit_tests(packages: List[str], output_dir: str = "coverage") -> Path:
    """Run Go unit tests with coverage profiling.
    
    Args:
        packages: List of package paths to test.
        output_dir: Directory for coverage output files.
    
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
    has_failures = False
    
    for idx, package in enumerate(packages, 1):
        safe_name = package.replace('/', '_').replace('.', '_').strip('_') or 'root'
        coverage_file = coverage_dir / f"{safe_name}.out"
        coverage_file.unlink(missing_ok=True)
        
        print(f"\n[{idx}/{len(packages)}] Testing: {package}")
        
        cmd = [
            "go", "test", "-v", "-p=1", "-parallel=1", "-timeout=30m", "-skip=^TestAcc",
            f"-coverprofile={coverage_file}",
            "-covermode=atomic",
            f"./{package}"
        ]
        
        result = subprocess.run(
            cmd,
            env={**os.environ, "TF_ACC": "0", "GOMAXPROCS": "1", "GOFLAGS": "-p=1"},
            check=False
        )

        if result.returncode != 0:
            has_failures = True
            print(f"❌ Unit tests failed in {package}")
        
        if coverage_file.exists():
            coverage_files.append(coverage_file)
            print(f"✅ Coverage generated for {package}")
        else:
            print(f"⚠️  No coverage file for {package}")
    
    # Merge coverage files
    print(f"\n📊 Merging {len(coverage_files)} coverage file(s)...")
    _merge_coverage_files(coverage_files, merged_file)

    if has_failures:
        raise subprocess.CalledProcessError(1, ["go", "test", *packages])
    
    print(f"✅ Merged coverage file: {merged_file}")
    return merged_file


def run_race_detection(packages: List[str]) -> int:
    """Run Go race detector on specified packages.
    
    Args:
        packages: List of package paths to test.
    
    Returns:
        Exit code: 0 if all tests pass, 1 if any fail.
    """
    print("\n" + "="*60)
    print("🔍 Running Race Detection Tests")
    print("="*60)
    
    has_failures = False
    
    for idx, package in enumerate(packages, 1):
        print(f"\n[{idx}/{len(packages)}] Testing: {package}")
        
        cmd = ["go", "test", "-v", "-race", "-p=1", "-parallel=1",
               "-timeout=30m", "-skip=^TestAcc", f"./{package}"]
        
        result = subprocess.run(
            cmd,
            env={**os.environ, "TF_ACC": "0", "GOMAXPROCS": "1", "GOFLAGS": "-p=1"},
            check=False
        )
        
        if result.returncode != 0:
            has_failures = True
            print(f"❌ Race detection failed in {package}")
        else:
            print(f"✅ Race detection passed in {package}")
    
    return 1 if has_failures else 0


def _merge_coverage_files(coverage_files: List[Path], output_file: Path) -> None:
    """Merge multiple Go coverage files into one.
    
    Args:
        coverage_files: List of coverage file paths.
        output_file: Path for merged output file.
    """
    with open(output_file, 'w', encoding='utf-8') as out_f:
        out_f.write("mode: atomic\n")
        for cov_file in coverage_files:
            with open(cov_file, 'r', encoding='utf-8') as in_f:
                for line in in_f:
                    if not line.startswith('mode:'):
                        out_f.write(line)
