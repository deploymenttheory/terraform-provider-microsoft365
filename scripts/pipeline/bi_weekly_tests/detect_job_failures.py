#!/usr/bin/env python3
"""Detects job-level failures in GitHub Actions workflow using GitHub API.

This script analyzes GitHub Actions workflow runs to detect infrastructure-level
failures that are not test failures, including:
- Job timeouts
- Job cancellations
- Infrastructure failures (OOM, runner issues)
- Setup/dependency step failures that prevented test execution

Usage:
    ./detect_job_failures.py <owner> <repo> <run-id> <output-file>

Args:
    owner: GitHub repository owner.
    repo: GitHub repository name.
    run-id: GitHub Actions workflow run ID.
    output-file: Path to write job failures JSON output.
"""

import sys
import json
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Optional


def run_gh_command(args: list[str]) -> str:
    """Run a GitHub CLI command and return output.

    Args:
        args: List of arguments to pass to 'gh' command.

    Returns:
        Command stdout as a string.

    Raises:
        subprocess.CalledProcessError: If the command fails.
    """
    try:
        result = subprocess.run(
            ["gh"] + args,
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Error running gh command: {e.stderr}", file=sys.stderr)
        raise


def get_workflow_jobs(owner: str, repo: str, run_id: str) -> list[dict]:
    """Get all jobs for a workflow run.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        run_id: GitHub Actions workflow run ID.

    Returns:
        List of job dictionaries from GitHub API.
    """
    try:
        result = run_gh_command([
            "api",
            f"/repos/{owner}/{repo}/actions/runs/{run_id}/jobs",
            "--paginate"
        ])
        
        if not result:
            return []
        
        response = json.loads(result)
        return response.get("jobs", [])
    except (subprocess.CalledProcessError, json.JSONDecodeError) as e:
        print(f"Error getting workflow jobs: {e}", file=sys.stderr)
        return []


def analyze_job_failure(job: dict) -> Optional[dict]:
    """Analyze a failed job to determine failure type and details.

    Distinguishes between test failures (which are handled separately) and
    infrastructure/setup failures that require special attention.

    Args:
        job: Job dictionary from GitHub API.

    Returns:
        Dictionary with failure details if this is an infrastructure failure,
        None if it's a test failure or the job didn't fail.
    """
    if job.get("conclusion") not in ["failure", "cancelled", "timed_out"]:
        return None

    # Test failures use continue-on-error and are
    # tracked separately. We only report infrastructure/setup failures here.
    job_name = job.get("name", "Unknown")

    # Test jobs have continue-on-error, so we only
    # want to report their infrastructure step failures, not test failures
    test_keywords = ["test", "provider", "resource", "datasource"]
    if any(keyword in job_name.lower() for keyword in test_keywords):
        # Test step failures are handled by
        # test result tracking. We only care about setup/infra failures.
        steps = job.get("steps", [])
        for step in steps:
            step_name = step.get("name", "").lower()
            step_conclusion = step.get("conclusion")

            test_step_keywords = [
                "run test",
                "run provider",
                "run resource",
                "run datasource"
            ]
            is_test_step = any(kw in step_name for kw in test_step_keywords)

            if step_conclusion == "failure" and not is_test_step:
                return {
                    "job_name": job_name,
                    "job_id": job.get("id"),
                    "conclusion": job.get("conclusion"),
                    "failure_type": "infrastructure",
                    "failed_step": step.get("name"),
                    "step_number": step.get("number"),
                    "started_at": job.get("started_at"),
                    "completed_at": job.get("completed_at"),
                    "html_url": job.get("html_url"),
                    "runner_name": job.get("runner_name"),
                }
    
    # Non-test jobs don't use continue-on-error,
    # so any failure is a real infrastructure issue
    failed_step = None
    if job.get("conclusion") == "timed_out":
        failure_type = "timeout"
    elif job.get("conclusion") == "cancelled":
        failure_type = "cancelled"
    else:
        steps = job.get("steps", [])
        for step in steps:
            if step.get("conclusion") == "failure":
                failed_step = step
                break

        # Helps categorize infrastructure issues
        # for better debugging and metrics
        if failed_step:
            step_name = failed_step.get("name", "").lower()
            if "memory" in step_name or "oom" in step_name:
                failure_type = "out_of_memory"
            elif "runner" in step_name or "setup" in step_name:
                failure_type = "runner_failure"
            else:
                failure_type = "step_failure"
        else:
            failure_type = "unknown"
    
    return {
        "job_name": job_name,
        "job_id": job.get("id"),
        "conclusion": job.get("conclusion"),
        "failure_type": failure_type,
        "failed_step": failed_step.get("name") if failed_step else "Unknown",
        "step_number": failed_step.get("number") if failed_step else None,
        "started_at": job.get("started_at"),
        "completed_at": job.get("completed_at"),
        "html_url": job.get("html_url"),
        "runner_name": job.get("runner_name"),
    }


def detect_job_failures(owner: str, repo: str, run_id: str, output_file: str) -> None:
    """Detect job-level failures and write to JSON file.

    Main processing function that queries GitHub API for all jobs in a workflow run,
    analyzes each job for infrastructure failures, and writes results to JSON.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        run_id: GitHub Actions workflow run ID.
        output_file: Path to write job failures JSON output.
    """
    print(f"\n{'='*60}")
    print(f"Detecting job-level failures for run {run_id}")
    print(f"{'='*60}\n")
    
    print("Querying GitHub API for job statuses...")
    jobs = get_workflow_jobs(owner, repo, run_id)
    
    if not jobs:
        print("⚠️  No jobs found or API query failed")
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump([], f, indent=2)
        return
    
    print(f"Found {len(jobs)} total jobs in workflow run")
    
    job_failures = []
    for job in jobs:
        failure = analyze_job_failure(job)
        if failure:
            job_failures.append(failure)
            print("\n❌ Job Failure Detected:")
            print(f"   Job: {failure['job_name']}")
            print(f"   Type: {failure['failure_type']}")
            print(f"   Failed Step: {failure['failed_step']}")
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(job_failures, f, indent=2)
    
    if job_failures:
        print(f"\n{'='*60}")
        print(f"Total job-level failures: {len(job_failures)}")
        print(f"{'='*60}")
    else:
        print("\n✅ No job-level failures detected")
    
    print(f"\nResults written to {output_file}")


def main():
    """Main entry point for the script.

    Parses command-line arguments and invokes job failure detection.
    """
    if len(sys.argv) < 4:
        print(
            "Usage: detect_job_failures.py <owner> <repo> <run-id> "
            "[output-file]",
            file=sys.stderr
        )
        sys.exit(1)
    
    owner = sys.argv[1]
    repo = sys.argv[2]
    run_id = sys.argv[3]
    output_file = sys.argv[4] if len(sys.argv) > 4 else "job-failures.json"
    
    detect_job_failures(owner, repo, run_id, output_file)


if __name__ == "__main__":
    main()

