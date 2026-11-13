#!/usr/bin/env python3
"""Aggregate test artifacts without treating missing evidence as passing tests.

Usage: merge_test_results.py ARTIFACTS FAILURES SUCCESSES JOBS STATUS

Always write both result arrays and a completeness report. Valid partial results
can create incidents. Verified passes from completed packages can resolve their
own incidents even when unrelated jobs are incomplete.
"""

import json
import os
import sys
from dataclasses import dataclass, field
from pathlib import Path

TEST_STEPS = {
    "Run provider core tests",
    "Run resource tests",
    "Run datasource tests",
    "Run actions tests",
    "Run list-resources tests",
    "Run ephemerals tests",
}


def read_results(path: Path, failures: bool) -> list[dict]:
    """Load one test result artifact from disk."""
    results = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(results, list):
        raise ValueError("expected an array")
    fields = ("test_name", "package", "configuration_block_type", "service")
    for result in results:
        if not isinstance(result, dict) or any(
            not isinstance(result.get(key), str) for key in fields
        ):
            raise ValueError("invalid test identity")
        if not result["test_name"] or not result["package"]:
            raise ValueError("missing test or package identity")
        if not isinstance(result.get("package_complete"), bool):
            raise ValueError("missing package completion status")
        if failures and not isinstance(result.get("context"), str):
            raise ValueError("missing failure context")
    return results


@dataclass
class ArtifactReport:
    """Keep observed outcomes separate from evidence eligible to resolve incidents."""

    failures: list[dict] = field(default_factory=list)
    successes: list[dict] = field(default_factory=list)
    verified_passes: set[tuple[str, str]] = field(default_factory=set)
    problems: list[str] = field(default_factory=list)


class ForeignArtifactError(ValueError):
    """The manifest explicitly identifies another workflow run or attempt."""


def _job_inventory(jobs_file: Path) -> tuple[int | None, list[str]]:
    """Count executed test jobs only after validating the whole inventory."""
    try:
        jobs = json.loads(jobs_file.read_text(encoding="utf-8"))
        if not isinstance(jobs, list) or not jobs:
            raise ValueError("no workflow jobs available")
        # Expected credential skips never require an artifact or create an incident.
        expected = sum(
            any(
                step.get("name") in TEST_STEPS and step.get("conclusion") != "skipped"
                for step in job.get("steps", [])
            )
            for job in jobs
        )
        problems = (
            ["Workflow has failed or interrupted jobs"]
            if any(job.get("conclusion") in ("failure", "cancelled", "timed_out") for job in jobs)
            else []
        )
        return expected, problems
    except (OSError, ValueError, TypeError, AttributeError) as error:
        return None, [f"Workflow job inventory unavailable: {error}"]


def _read_manifest(path: Path, provenance: dict[str, str]) -> dict:
    """Validate provenance and the declared execution scope before trusting outcomes."""
    manifest = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(manifest, dict):
        raise ValueError("unrecognized test execution manifest")
    for key, expected in provenance.items():
        actual = manifest.get(key)
        if expected and actual is not None and actual != expected:
            raise ForeignArtifactError(f"artifact belongs to another {key}")
    for key, expected in provenance.items():
        if expected and manifest.get(key) != expected:
            raise ValueError(f"missing artifact {key}")
    if manifest.get("schema_version") != 1:
        raise ValueError("unrecognized test execution manifest")
    expected_packages = manifest.get("expected_packages")
    completed_packages = manifest.get("completed_packages")
    if (
        not isinstance(manifest.get("complete"), bool)
        or not isinstance(expected_packages, list)
        or not isinstance(completed_packages, list)
        or any(
            not isinstance(item, str) or not item for item in expected_packages + completed_packages
        )
        or not set(completed_packages) <= set(expected_packages)
    ):
        raise ValueError("missing or invalid package execution scope")
    return manifest


def _artifact_report(
    directory: Path, provenance: dict[str, str], verify_passes: bool
) -> ArtifactReport:
    """Preserve valid partial evidence, but discard explicitly foreign outcomes."""
    report = ArtifactReport()
    for filename, target in (
        ("test-failures.json", report.failures),
        ("test-successes.json", report.successes),
    ):
        try:
            target.extend(read_results(directory / filename, filename == "test-failures.json"))
        except (OSError, ValueError) as error:
            report.problems.append(f"{directory.name}/{filename}: {error}")
    arrays_valid = not report.problems
    try:
        manifest = _read_manifest(directory / "test-run.json", provenance)
    except ForeignArtifactError as error:
        report.failures.clear()
        report.successes.clear()
        report.problems.append(f"{directory.name}/test-run.json: {error}")
        return report
    except (OSError, ValueError) as error:
        report.problems.append(f"{directory.name}/test-run.json: {error}")
        return report
    if any(
        result["package"] not in manifest["expected_packages"]
        for result in report.failures + report.successes
    ):
        report.problems.append(f"{directory.name}: test result outside declared package scope")
        return report
    if not manifest["complete"]:
        report.problems.append(f"{directory.name}/test-run.json: test execution incomplete")
    if arrays_valid and verify_passes:
        report.verified_passes.update(
            (result["package"], result["test_name"])
            for result in report.successes
            if result["package_complete"] and result["package"] in manifest["completed_packages"]
        )
    return report


def _artifact_directories(artifacts: Path) -> list[Path]:
    """Include result files and diagnostic-only artifact directories."""

    directories = sorted(
        {
            path.parent
            for name in ("test-failures.json", "test-successes.json", "test-run.json")
            for path in artifacts.rglob(name)
        }
    )
    # Diagnostic-only artifacts also indicate a job that produced no test results.
    return sorted(
        set(directories) | {path for path in artifacts.glob("test-results-*") if path.is_dir()}
    )


def aggregate(
    artifacts: Path, jobs_file: Path, download_ok: bool = True
) -> tuple[list, list, dict]:
    """Combine outcomes while independently verifying each completed package's passes."""
    expected, problems = _job_inventory(jobs_file)
    provenance = {
        "run_id": os.environ.get("GITHUB_RUN_ID", ""),
        "run_attempt": os.environ.get("GITHUB_RUN_ATTEMPT", ""),
    }
    if not download_ok:
        problems.append("Artifact download did not succeed")
    if not all(provenance.values()):
        problems.append("Current workflow run and attempt provenance is unavailable")
    directories = _artifact_directories(artifacts)
    if expected is not None and len(directories) != expected:
        problems.append(f"Expected {expected} test artifacts; found {len(directories)}")
    if not artifacts.is_dir():
        problems.append("Test artifact directory is missing")

    combined = ArtifactReport(problems=problems)
    verify_passes = bool(expected) and download_ok and all(provenance.values())
    for directory in directories:
        report = _artifact_report(directory, provenance, verify_passes)
        combined.failures.extend(report.failures)
        combined.successes.extend(report.successes)
        combined.verified_passes.update(report.verified_passes)
        combined.problems.extend(report.problems)
    combined.verified_passes.difference_update(
        (item["package"], item["test_name"]) for item in combined.failures
    )
    status = {
        "schema_version": 2,
        **provenance,
        "verified_passes": [
            {"package": package, "test_name": name}
            for package, name in sorted(combined.verified_passes)
        ],
        "complete": not combined.problems,
        "expected_artifacts": expected,
        "received_artifacts": len(directories),
        "problems": combined.problems,
    }
    return combined.failures, combined.successes, status


def main():
    """Merge test result artifacts into a workflow summary."""
    if len(sys.argv) != 6:
        sys.exit("Usage: merge_test_results.py ARTIFACTS FAILURES SUCCESSES JOBS STATUS")
    artifacts, failures_path, successes_path, jobs_path, status_path = map(Path, sys.argv[1:])
    failures, successes, status = aggregate(
        artifacts, jobs_path, os.environ.get("ARTIFACT_DOWNLOAD_OUTCOME", "success") == "success"
    )
    for path, data in (
        (failures_path, failures),
        (successes_path, successes),
        (status_path, status),
    ):
        path.write_text(json.dumps(data, indent=2), encoding="utf-8")
    message = f"Collected {len(failures)} test failures and {len(successes)} explicit passes."
    message += (
        " Reporting complete."
        if status["complete"]
        else " Reporting incomplete; only independently verified passes can resolve incidents."
    )
    print(message)
    for problem in status["problems"]:
        print(f"  Warning: {problem}")
    if os.environ.get("GITHUB_STEP_SUMMARY"):
        with open(os.environ["GITHUB_STEP_SUMMARY"], "a", encoding="utf-8") as summary:
            summary.write(
                message + "\n" + "\n".join(f"- {problem}" for problem in status["problems"]) + "\n"
            )


if __name__ == "__main__":
    main()
