#!/usr/bin/env python3
"""Aggregate test artifacts without treating missing evidence as passing tests.

Usage: merge_test_results.py ARTIFACTS FAILURES SUCCESSES JOBS STATUS

Always write both result arrays and a completeness report. Valid partial results
can create incidents, but incomplete reporting disables automatic issue closure.
"""

import json
import os
import sys
from pathlib import Path


TEST_STEPS = {
    "Run provider core tests", "Run resource tests", "Run datasource tests",
    "Run actions tests", "Run list-resources tests", "Run ephemerals tests",
}


def read_results(path: Path, failures: bool) -> list[dict]:
    results = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(results, list):
        raise ValueError("expected an array")
    fields = ("test_name", "package", "configuration_block_type", "service")
    for result in results:
        if not isinstance(result, dict) or any(not isinstance(result.get(key), str) for key in fields):
            raise ValueError("invalid test identity")
        if not result["test_name"] or not result["package"]:
            raise ValueError("missing test or package identity")
        if not isinstance(result.get("package_complete"), bool):
            raise ValueError("missing package completion status")
        if failures and not isinstance(result.get("context"), str):
            raise ValueError("missing failure context")
    return results


def aggregate(artifacts: Path, jobs_file: Path, download_ok: bool = True) -> tuple[list, list, dict]:
    failures, successes, problems = [], [], []
    expected = None
    if not download_ok:
        problems.append("Artifact download did not succeed")
    try:
        jobs = json.loads(jobs_file.read_text(encoding="utf-8"))
        if not isinstance(jobs, list) or not jobs:
            raise ValueError("no workflow jobs available")
        # Expected credential skips never require an artifact or create an incident.
        expected = sum(any(step.get("name") in TEST_STEPS and step.get("conclusion") != "skipped"
                           for step in job.get("steps", [])) for job in jobs)
        if any(job.get("conclusion") in ("failure", "cancelled", "timed_out") for job in jobs):
            problems.append("Workflow has failed or interrupted jobs")
    except (OSError, ValueError, TypeError, AttributeError) as error:
        problems.append(f"Workflow job inventory unavailable: {error}")

    directories = sorted({path.parent for name in ("test-failures.json", "test-successes.json", "test-run.json")
                          for path in artifacts.rglob(name)})
    # Diagnostic-only artifacts also indicate a job that produced no test results.
    directories = sorted(set(directories) | {path for path in artifacts.glob("test-results-*") if path.is_dir()})
    if expected is not None and len(directories) != expected:
        problems.append(f"Expected {expected} test artifacts; found {len(directories)}")
    if not artifacts.is_dir():
        problems.append("Test artifact directory is missing")

    for directory in directories:
        for filename, target in (("test-failures.json", failures), ("test-successes.json", successes)):
            try:
                target.extend(read_results(directory / filename, filename == "test-failures.json"))
            except (OSError, ValueError) as error:
                problems.append(f"{directory.name}/{filename}: {error}")
        try:
            manifest = json.loads((directory / "test-run.json").read_text(encoding="utf-8"))
            if manifest.get("schema_version") != 1 or manifest.get("complete") is not True:
                raise ValueError("test execution incomplete")
            for key, env in (("run_id", "GITHUB_RUN_ID"), ("run_attempt", "GITHUB_RUN_ATTEMPT")):
                if os.environ.get(env) and manifest.get(key) != os.environ[env]:
                    raise ValueError(f"artifact belongs to another {key}")
        except (OSError, ValueError, AttributeError) as error:
            problems.append(f"{directory.name}/test-run.json: {error}")

    status = {"complete": not problems, "expected_artifacts": expected,
              "received_artifacts": len(directories), "problems": problems}
    return failures, successes, status


def main():
    if len(sys.argv) != 6:
        sys.exit("Usage: merge_test_results.py ARTIFACTS FAILURES SUCCESSES JOBS STATUS")
    artifacts, failures_path, successes_path, jobs_path, status_path = map(Path, sys.argv[1:])
    failures, successes, status = aggregate(artifacts, jobs_path, os.environ.get("ARTIFACT_DOWNLOAD_OUTCOME", "success") == "success")
    for path, data in ((failures_path, failures), (successes_path, successes), (status_path, status)):
        path.write_text(json.dumps(data, indent=2), encoding="utf-8")
    message = f"Collected {len(failures)} test failures and {len(successes)} explicit passes."
    message += " Reporting complete." if status["complete"] else " Reporting incomplete; automatic issue closure disabled."
    print(message)
    for problem in status["problems"]:
        print(f"  Warning: {problem}")
    if os.environ.get("GITHUB_STEP_SUMMARY"):
        with open(os.environ["GITHUB_STEP_SUMMARY"], "a", encoding="utf-8") as summary:
            summary.write(message + "\n" + "\n".join(f"- {problem}" for problem in status["problems"]) + "\n")


if __name__ == "__main__":
    main()
