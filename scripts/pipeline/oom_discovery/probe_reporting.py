#!/usr/bin/env python3
"""Reproduce reporting edge cases locally, with all external calls blocked."""

import argparse
from contextlib import redirect_stdout, redirect_stderr
import importlib.util
import io
import json
from pathlib import Path
import subprocess
import sys
import tempfile
from unittest.mock import patch


def load(root, name):
    path = root / "scripts/pipeline/bi_weekly_tests" / f"{name}.py"
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def probe(root):
    issues = load(root, "manage_test_issues")
    detector = load(root, "detect_job_failures")
    merger = load(root, "merge_test_results")
    source = subprocess.check_output(["git", "-C", str(root), "rev-parse", "HEAD"], text=True).strip()
    with tempfile.TemporaryDirectory() as directory, redirect_stdout(io.StringIO()), redirect_stderr(io.StringIO()):
        temporary = Path(directory)
        failures = temporary / "failures.json"
        successes = temporary / "successes.json"
        failures.write_text("[]")
        successes.write_text(json.dumps([{"test_name": "TestActuallyPassed", "category": "resources", "service": "example"}]))
        with (
            patch.object(issues.subprocess, "run", side_effect=AssertionError("External calls are prohibited")) as external,
            patch.object(issues, "ensure_label_exists"),
            patch.object(issues, "get_all_open_test_issues", return_value=[
                {"number": 123, "title": "Bug: TestNeverRan Failing"},
            ]),
            patch.object(issues, "close_resolved_issue") as close,
            patch.object(issues.time, "sleep"),
        ):
            issues.process_test_failures("example", "fixture", "123", str(failures), str(successes))
            closure = {
                "scenario": "One test passes; an unrelated open test incident has no execution result",
                "closure_attempts": close.call_count,
                "explicit_pass_argument": close.call_args.args[-1] if close.called else None,
                "external_calls_performed": external.call_count,
            }
        cancelled_step = {
            "id": 123, "name": "Resources: example", "conclusion": "failure",
            "steps": [
                {"name": "Set up job", "conclusion": "success", "number": 1},
                {"name": "Run resources tests", "conclusion": "cancelled", "number": 2},
            ],
        }
        classification = detector.analyze_job_failure(cancelled_step)
        with patch.object(sys, "argv", ["merge_test_results.py", str(temporary / "missing-artifacts"), str(temporary / "merged.json")]):
            try:
                merger.main()
                merge_exit = 0
            except SystemExit as error:
                merge_exit = error.code
        return {
            "source_commit": source,
            "synthetic_fixtures": True,
            "issue_closure": closure,
            "classification": {
                "scenario": "Failed job with a cancelled test step, matching the incident job shape",
                "failure_type": classification["failure_type"],
                "failed_step": classification["failed_step"],
                "skipped_job_is_incident": detector.analyze_job_failure(
                    {
                        "name": "Resources: example", "conclusion": "success",
                        "steps": [{"name": "Run resources tests", "conclusion": "skipped"}],
                    }
                ) is not None,
            },
            "missing_artifacts": {
                "exit_code": merge_exit,
                "merged_file_created": (temporary / "merged.json").exists(),
            },
        }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source", type=Path, default=Path(__file__).resolve().parents[3])
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    result = probe(args.source.resolve())
    args.output.write_text(json.dumps(result, indent=2) + "\n")
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
