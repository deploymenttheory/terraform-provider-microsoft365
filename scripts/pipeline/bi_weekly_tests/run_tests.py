#!/usr/bin/env python3
"""Test runner for Terraform provider acceptance tests.

This script runs Terraform provider tests sequentially, one Go package at a time,
to conserve memory on resource-constrained runners. Supports configurable parallelism
and memory management options.

Usage:
    ./run_tests.py <type> [service] [coverage-file] [test-output-file] [options]

Positional Arguments:
    type              Test type: provider-core, resources, datasources, actions, list-resources, or
    ephemerals
    service           Service name (required for
    resources/datasources/actions/list-resources/ephemerals)
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

import argparse
import gc
import json
import os
import re
import signal
import subprocess
import sys
import time
from dataclasses import dataclass, replace
from datetime import datetime
from pathlib import Path
from typing import Dict, List


def run_command(cmd: List[str], output_file: str, append: bool = False) -> int:
    """Run a command and capture output to file and stdout.

    Args:
        cmd: Command and arguments to execute as a list.
        output_file: Path to file where output will be written.
        append: If True, append to output_file; if False, overwrite it.

    Returns:
        The exit code of the command.
    """
    mode = "a" if append else "w"
    started = time.monotonic()
    with open(output_file, mode, encoding="utf-8", buffering=8192) as f:  # Use small buffer
        with subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,  # Line buffered
        ) as process:
            message = f"[COMMAND] Started package={cmd[-1]} pid={process.pid}\n"
            print(message, end="", flush=True)
            f.write(message)
            f.flush()

            for line in process.stdout:
                print(line, end="", flush=True)
                f.write(line)
                f.flush()  # Ensure data is written immediately

            process.stdout.close()
            process.wait()
            status = process.returncode
            if status < 0:
                signal_name = signal.Signals(-status).name
                status = 128 - status
                message = f"[COMMAND] Terminated by {signal_name}; exit={status}"
            else:
                message = f"[COMMAND] Exit={status}"
            message += f" package={cmd[-1]} elapsed={time.monotonic() - started:.2f}s\n"
            print(message, end="", flush=True)
            f.write(message)
            return status


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
        rel_path = f"./{package_dir.relative_to(cwd).as_posix()}"
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
            check=False,
        )

        # Count lines that start with "Test" (test function names)
        test_count = 0
        for line in result.stdout.split("\n"):
            if line.startswith("Test"):
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


# Only top-level Go test functions are incident identities. Subtest lines and the
# CheckDestroy helper's per-resource lines (which embed random resource names) are
# indented or lack a Go test prefix; they remain in the parent test's context.
TOP_LEVEL_TEST = r"(?:Test|Fuzz|Example)[^\s/]*"
PACKAGE_FAILURE = re.compile(
    r"^FAIL[ \t]+\S+(?:[ \t]+(?:\d+(?:\.\d+)?s|\[(?:build|setup) failed\]))?[ \t]*$",
    re.MULTILINE,
)


def _failure_context(content: str, match, test_name: str) -> str:
    """Keep the relevant error output for an explicit top-level test failure."""
    search_start = max(0, match.start() - 50000)
    run_pattern = re.compile(rf"^=== RUN\s+{re.escape(test_name)}[ \t]*$", re.MULTILINE)
    starts = list(run_pattern.finditer(content[search_start : match.start()]))
    if not starts:
        return content[max(0, match.start() - 500) : match.start()].strip()
    context = content[search_start + starts[-1].end() : match.start()].strip()
    lines = context.split("\n")
    for index, line in enumerate(lines):
        if "[DEBUG]" in line or "[INFO]" in line:
            continue
        if any(indicator in line for indicator in [".go:", "Error:", "panic:", "    ", "\t"]):
            context = "\n".join(lines[index:])
            break
    return context[:1000] + "\n... (truncated)" if len(context) > 1000 else context


def _panic_context(output: str) -> str:
    """Keep bounded panic diagnostics without the package failure trailer."""
    lines = []
    capture = False
    for line in output.split("\n"):
        if not capture and any(
            indicator in line
            for indicator in ["panic:", "SIGSEGV", "runtime error:", "fatal error:"]
        ):
            capture = True
        if capture:
            if PACKAGE_FAILURE.fullmatch(line):
                break
            lines.append(line)
            if len(lines) > 50:
                break
    context = "\n".join(lines).strip()
    if len(context) > 1000:
        context = context[:1000] + "\n... (truncated)"
    return context or output[-500:].strip()


def _incomplete_test_failures(content: str, completed: set, metadata: dict) -> list:
    """Identify started tests terminated by a panic or package failure."""
    failures = []
    pattern = re.compile(rf"^=== RUN\s+({TOP_LEVEL_TEST})[ \t]*$", re.MULTILINE)
    for match in pattern.finditer(content):
        name = match.group(1)
        if name in completed:
            continue
        next_run = pattern.search(content, match.end())
        output = content[match.end() : next_run.start() if next_run else len(content)]
        if any(
            marker in output for marker in ["panic:", "SIGSEGV", "segmentation violation"]
        ) or PACKAGE_FAILURE.search(output):
            failures.append(
                {
                    **metadata,
                    "test_name": name,
                    "context": f"[PANIC DETECTED]\n{_panic_context(output)}",
                }
            )
            completed.add(name)
    return failures


def parse_package_results(
    content: str, configuration_block_type: str, service: str
) -> tuple[list, list]:
    """Extract outcomes within one package so repeated test names cannot collide."""
    metadata = {"configuration_block_type": configuration_block_type, "service": service}
    failures = [
        {
            **metadata,
            "test_name": match.group(1),
            "context": _failure_context(content, match, match.group(1)),
        }
        for match in re.finditer(rf"^--- FAIL: ({TOP_LEVEL_TEST})(?=\s)", content, re.MULTILINE)
    ]
    successes = [
        {**metadata, "test_name": match.group(1)}
        for match in re.finditer(rf"^--- PASS: ({TOP_LEVEL_TEST})(?=\s)", content, re.MULTILINE)
    ]
    completed = {result["test_name"] for result in failures + successes}
    completed.update(re.findall(rf"^--- SKIP: ({TOP_LEVEL_TEST})(?=\s)", content, re.MULTILINE))
    failures.extend(_incomplete_test_failures(content, completed, metadata))
    return failures, successes


def _command_sections(content: str):
    """Keep separate command frames, including legacy output without a package identity."""
    starts = list(
        re.finditer(r"^\[COMMAND\] Started package=(\S+) pid=\d+\n", content, re.MULTILINE)
    )
    sections = []
    if not starts:
        sections.append(("", content))
    for index, start in enumerate(starts):
        end = starts[index + 1].start() if index + 1 < len(starts) else len(content)
        sections.append((start.group(1), content[start.end() : end]))

    return sections, len(starts)


def _section_results(sections: list, identity: tuple):
    """Attach package completeness independently from the individual test verdicts."""
    failures, successes, completed = [], [], []
    for package, output in sections:
        exit_match = (
            re.search(
                rf"^\[COMMAND\] Exit=(\d+) package={re.escape(package)} elapsed=[\d.]+s$",
                output,
                re.MULTILINE,
            )
            if package
            else None
        )
        status = int(exit_match.group(1)) if exit_match else None
        package_failures, package_successes = parse_package_results(output, *identity)
        # Exit 1 without a failed test can be a build/vet failure before tests ran.
        if status == 0 or (status == 1 and package_failures):
            completed.append(package)
        for result in package_failures + package_successes:
            result["package"] = package
            result["package_complete"] = status == 0
        failures.extend(package_failures)
        successes.extend(package_successes)

    return failures, successes, completed


def parse_test_results(
    output_file: str,
    configuration_block_type: str,
    service: str,
    expected_packages: list[str] | None = None,
) -> None:
    """Preserve partial results, package identity and command completion separately.

    Only an explicit PASS in a package that exited successfully can resolve an
    incident. Unframed legacy output is retained for diagnosis without claiming
    package identity or completeness.
    """
    content = Path(output_file).read_text(encoding="utf-8")
    sections, start_count = _command_sections(content)
    failures, successes, completed = _section_results(sections, (configuration_block_type, service))

    Path("test-failures.json").write_text(json.dumps(failures, indent=2), encoding="utf-8")
    Path("test-successes.json").write_text(json.dumps(successes, indent=2), encoding="utf-8")
    manifest = {
        "schema_version": 1,
        "run_id": os.environ.get("GITHUB_RUN_ID", ""),
        "run_attempt": os.environ.get("GITHUB_RUN_ATTEMPT", ""),
        "expected_packages": expected_packages or [],
        "completed_packages": completed,
        "complete": (
            expected_packages is not None
            and sorted(completed) == sorted(expected_packages)
            and start_count == len(expected_packages)
        ),
    }
    Path("test-run.json").write_text(json.dumps(manifest, indent=2), encoding="utf-8")
    print(
        f"Test results: {len(failures)} failures, {len(successes)} successes;"
        f" complete={manifest['complete']}"
    )


@dataclass
# Explicit data fields mirror the configuration or external metadata contract.
class TestRunOptions:  # pylint: disable=too-many-instance-attributes
    """Control bounded Go execution without changing the inherited credentials or memory limit."""

    max_procs: int = 2
    test_parallel: int = 1
    pkg_parallel: int = 1
    use_race: bool = False
    skip_enumeration: bool = False
    force_gc: bool = True
    shard_resources: str = ""
    test_filters: Dict[str, str] | None = None


@dataclass
class PackageTestRun:
    """Execute selected packages and retain their coverage and incident identities."""

    coverage_file: str
    test_output_file: str
    options: TestRunOptions
    identity: tuple[str, str]

    def _label(self, package: str) -> str:
        if self.identity[0] == "provider-core":
            return package
        prefix = f"./internal/services/{self.identity[0]}/{self.identity[1]}"
        return package.replace(prefix, "").lstrip("/") or "."

    def _enumerate(self, packages: list[str]) -> dict[str, int]:
        if self.options.skip_enumeration:
            print("⏭️  [ENUMERATE] Skipping test enumeration (--skip-enumeration flag)")
            return dict.fromkeys(packages, 0)
        counts = {}
        started = time.time()
        for index, package in enumerate(packages, 1):
            print(f"📊 [ENUMERATE] Package {index}/{len(packages)}: {package}")
            counts[package] = count_tests_in_package(package)
        print(f"⏱️  [TIMING] Total enumeration time: {time.time() - started:.2f}s")
        return counts

    def _summary(self, packages: list[str], counts: dict[str, int]) -> None:
        print_separator("=")
        print(f"📋 Test Discovery Summary for {'/'.join(filter(None, self.identity))}")
        print(f"Total Packages: {len(packages)}")
        if not self.options.skip_enumeration:
            print(f"Total Tests: {sum(counts.values())}")
        for package in packages:
            label = self._label(package)
            print(
                f"  📦 {label}"
                if self.options.skip_enumeration
                else f"  📦 {label:<50} {counts[package]:>4} test(s)"
            )
        print_separator("=")

    def _command(self, package: str) -> list[str]:
        command = ["go", "test", "-v"]
        if self.options.use_race:
            command.append("-race")
        if package in (self.options.test_filters or {}):
            command.extend(["-run", self.options.test_filters[package]])
        command.extend(
            [
                "-timeout=90m",
                f"-p={self.options.pkg_parallel}",
                f"-parallel={self.options.test_parallel}",
                f"-coverprofile={self.coverage_file}.tmp",
                "-covermode=atomic",
                package,
            ]
        )
        return command

    def _append_coverage(self) -> None:
        temporary = Path(f"{self.coverage_file}.tmp")
        if not temporary.exists():
            print("⚠️  [COVERAGE] No coverage file generated")
            return
        with (
            temporary.open(encoding="utf-8") as source,
            open(self.coverage_file, "a", encoding="utf-8") as destination,
        ):
            for line in source:
                if not line.startswith("mode:"):
                    destination.write(line)
        temporary.unlink()
        print("🗑️  [COVERAGE] Deleted temporary coverage file")

    def _run_package(self, package: str, index: int, total: int, count: int) -> int:
        started = time.time()
        print_separator("-", 70)
        print(f"📦 Package {index}/{total}: {self._label(package)}")
        print(f"   Tests: {count}")
        print(f"⏱️  [TIMING] Started at: {datetime.now().strftime('%H:%M:%S')}")
        print_separator("-", 70)
        command = self._command(package)
        print(f"▶️  [RUN] Executing: {' '.join(command)}")
        status = run_command(command, self.test_output_file, append=index > 1)
        print(f"⏱️  [TIMING] Package execution took {time.time() - started:.2f}s")
        self._append_coverage()
        if status:
            print(f"❌ [RESULT] Package {self._label(package)} failed (exit code: {status})")
        else:
            print(f"✅ [RESULT] Package {self._label(package)} completed successfully")
        return status

    def run(self, packages: list[str]) -> int:
        """Run packages serially and stop after termination while preserving partial results."""
        print(f"⚙️  [CONFIG] GOMAXPROCS: {self.options.max_procs}")
        print(f"⚙️  [CONFIG] GOMEMLIMIT: {os.environ.get('GOMEMLIMIT', 'unset')}")
        print(f"⚙️  [CONFIG] Test parallel: {self.options.test_parallel}")
        print(f"⚙️  [CONFIG] Package parallel: {self.options.pkg_parallel}")
        print(f"⚙️  [CONFIG] Race detection: {self.options.use_race}")
        Path(self.coverage_file).write_text("mode: atomic\n", encoding="utf-8")
        if not packages:
            print("⚠️  No test packages found, creating empty coverage file")
            return 0
        counts = self._enumerate(packages)
        self._summary(packages, counts)
        failure_status = 0
        for index, package in enumerate(packages, 1):
            status = self._run_package(package, index, len(packages), counts[package])
            if status:
                failure_status = status if status >= 128 else 1
            if status >= 128:
                print(
                    "[TERMINATION] Stopping after command termination; preserving partial test"
                    " results",
                    flush=True,
                )
                break
            if self.options.force_gc:
                print("🗑️  [MEMORY] Running garbage collection...")
                gc.collect()
        print("📝 Parsing test results...")
        parse_test_results(self.test_output_file, *self.identity, packages)
        print_separator("=")
        print(f"🏁 Sequential execution complete for {'/'.join(filter(None, self.identity))}")
        print_separator("=")
        return failure_status


def run_provider_core_tests(
    coverage_file: str, test_output_file: str, *positional, **options
) -> int:
    """Discover core test packages and execute them with bounded run options."""
    if len(positional) < 4:
        options.setdefault("use_race", True)
    packages = []
    for directory in [
        "./internal/client",
        "./internal/helpers",
        "./internal/provider",
        "./internal/utilities",
    ]:
        packages.extend(discover_test_packages(Path(directory)))
    return PackageTestRun(
        coverage_file,
        test_output_file,
        TestRunOptions(*positional, **options),
        ("provider-core", ""),
    ).run(packages)


def _select_service_packages(directory: Path, options: TestRunOptions) -> tuple[list[str], dict]:
    """Validate the shard/filter contract and select exact packages in that service."""
    if options.test_filters is not None and (
        not isinstance(options.test_filters, dict)
        or any(
            not isinstance(key, str) or not isinstance(value, str) or not value
            for key, value in options.test_filters.items()
        )
    ):
        raise ValueError("Test filters must map package paths to non-empty test expressions")
    prefix = f"./{directory.as_posix()}"
    filters = {f"{prefix}/{key}": value for key, value in (options.test_filters or {}).items()}
    if not directory.exists():
        if options.shard_resources or filters:
            raise ValueError(f"Shard service directory does not exist: {prefix}")
        return [], filters
    packages = discover_test_packages(directory)
    if options.shard_resources:
        selected = []
        for resource_path in [
            item.strip() for item in options.shard_resources.split(",") if item.strip()
        ]:
            package = f"{prefix}/{resource_path}"
            if package not in packages:
                raise ValueError(f"Shard resource not found: {package}")
            selected.append(package)
        packages = selected
    unmatched = set(filters) - set(packages)
    if unmatched:
        raise ValueError(f"Test filters reference packages outside this shard: {sorted(unmatched)}")
    return packages, filters


def run_service_tests(
    configuration_block_type: str,
    service: str,
    coverage_file: str,
    test_output_file: str,
    *positional,
    **options,
) -> int:
    """Select one service shard and execute its packages with bounded run options."""
    run_options = TestRunOptions(*positional, **options)
    packages, filters = _select_service_packages(
        Path(f"internal/services/{configuration_block_type}/{service}"), run_options
    )
    return PackageTestRun(
        coverage_file,
        test_output_file,
        replace(run_options, test_filters=filters),
        (configuration_block_type, service),
    ).run(packages)


def main():
    """Run the requested provider acceptance test packages."""
    parser = argparse.ArgumentParser(
        description=(
            "Run Terraform Provider tests with configurable parallelism and memory management"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Run with default settings (low memory mode)
  ./run_tests.py resources identity_and_access coverage.txt output.log
  
  # Run with more parallelism (requires more memory)
  ./run_tests.py resources identity_and_access coverage.txt output.log --max-procs 4 --test-parallel 4
  
  # Disable race detection for faster execution
  ./run_tests.py resources identity_and_access coverage.txt output.log --no-race
        """,
    )

    # Positional arguments
    parser.add_argument(
        "type",
        choices=[
            "provider-core",
            "resources",
            "datasources",
            "actions",
            "list-resources",
            "ephemerals",
        ],
        help="Type of tests to run",
    )
    parser.add_argument(
        "service",
        nargs="?",
        default="",
        help="Service name (required for resources/datasources/actions/list-resources/ephemerals)",
    )
    parser.add_argument(
        "coverage_file",
        nargs="?",
        default="coverage.txt",
        help="Output file for coverage data (default: coverage.txt)",
    )
    parser.add_argument(
        "output_file",
        nargs="?",
        default="test-output.log",
        help="Output file for test logs (default: test-output.log)",
    )

    # Optional arguments for memory/parallelism control
    parser.add_argument(
        "--max-procs",
        type=int,
        default=2,
        help="GOMAXPROCS value - max CPU cores for Go (default: 2)",
    )
    parser.add_argument(
        "--test-parallel",
        type=int,
        default=1,
        help="Number of tests to run in parallel within a package (default: 1)",
    )
    parser.add_argument(
        "--pkg-parallel",
        type=int,
        default=1,
        help="Number of packages to build/test in parallel (default: 1)",
    )
    parser.add_argument(
        "--race",
        dest="use_race",
        action="store_true",
        default=None,
        help="Enable race detection (default for provider-core)",
    )
    parser.add_argument(
        "--no-race",
        dest="use_race",
        action="store_false",
        help="Disable race detection (default for resources/datasources)",
    )
    parser.add_argument(
        "--skip-enumeration",
        action="store_true",
        default=False,
        help="Skip test enumeration phase (faster startup, no test counts)",
    )
    parser.add_argument(
        "--force-gc",
        action="store_true",
        default=True,
        help="Force garbage collection between packages (default: enabled)",
    )
    parser.add_argument(
        "--no-force-gc",
        dest="force_gc",
        action="store_false",
        help="Disable forced garbage collection",
    )
    parser.add_argument(
        "--shard-resources",
        type=str,
        default="",
        help=(
            "Comma-separated list of resource paths in this shard (e.g.,"
            ' "graph_beta/resource1,graph_v1.0/resource2")'
        ),
    )
    parser.add_argument(
        "--test-filters",
        type=json.loads,
        default={},
        help="JSON object mapping relative package paths to Go test expressions",
    )

    args = parser.parse_args()

    # Set default for use_race if not specified
    if args.use_race is None:
        args.use_race = args.type == "provider-core"

    print("\n" + "=" * 70)
    print("🚀 [MAIN] Terraform Provider Test Runner Started")
    print("=" * 70)
    print(f"📋 [MAIN] Command: {' '.join(sys.argv)}")
    print(f"📂 [MAIN] Working directory: {Path.cwd()}")
    print("=" * 70 + "\n")

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

    if os.environ.get("SKIP_TESTS", "false") == "true":
        print("⏭️  Skipping tests - no credentials configured")
        with open(args.coverage_file, "w", encoding="utf-8") as f:
            f.write("mode: atomic\n")
        sys.exit(0)

    if args.type == "provider-core":
        exit_code = run_provider_core_tests(
            args.coverage_file,
            args.output_file,
            max_procs=args.max_procs,
            test_parallel=args.test_parallel,
            pkg_parallel=args.pkg_parallel,
            use_race=args.use_race,
            skip_enumeration=args.skip_enumeration,
            force_gc=args.force_gc,
        )
    elif args.type in ["resources", "datasources", "actions", "list-resources", "ephemerals"]:
        if not args.service:
            print(f"Error: service name required for {args.type} tests", file=sys.stderr)
            sys.exit(1)
        exit_code = run_service_tests(
            args.type,
            args.service,
            args.coverage_file,
            args.output_file,
            max_procs=args.max_procs,
            test_parallel=args.test_parallel,
            pkg_parallel=args.pkg_parallel,
            use_race=args.use_race,
            skip_enumeration=args.skip_enumeration,
            force_gc=args.force_gc,
            shard_resources=args.shard_resources,
            test_filters=args.test_filters,
        )
    else:
        print(f"Error: unknown test type: {args.type}", file=sys.stderr)
        print(
            "Valid types: provider-core, resources, datasources, actions, list-resources,"
            " ephemerals",
            file=sys.stderr,
        )
        sys.exit(1)

    print("Tests completed")
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
