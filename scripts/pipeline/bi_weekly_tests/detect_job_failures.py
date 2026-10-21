#!/usr/bin/env python3
"""Classify failed jobs using step status, OOM counters and runner annotations.

Usage: detect_job_failures.py OWNER REPO RUN_ID [FAILURES_JSON] [JOBS_JSON]
Missing logs leave the cause unknown; a signal or cancelled step alone is not OOM.
"""

import json
import os
import re
import subprocess
import sys
from pathlib import Path


def run_gh_command(args: list[str]) -> str:
    return subprocess.run(["gh"] + args, capture_output=True, text=True, check=True).stdout.strip()


def get_workflow_jobs(owner: str, repo: str, run_id: str) -> list[dict]:
    # --slurp produces one JSON document even when there are more than 100 jobs.
    endpoint = f"/repos/{owner}/{repo}/actions/runs/{run_id}"
    attempt = os.environ.get("GITHUB_RUN_ATTEMPT")
    if attempt:
        endpoint += f"/attempts/{attempt}"
    pages = json.loads(run_gh_command(["api", f"{endpoint}/jobs?per_page=100", "--paginate", "--slurp"]))
    jobs = [job for page in pages for job in page["jobs"]]
    if not jobs:
        raise ValueError("GitHub returned no workflow jobs")
    return jobs


def analyze_job_failure(job: dict, log: str = "", annotations: list[dict] | None = None) -> dict | None:
    if job.get("conclusion") not in ("failure", "cancelled", "timed_out"):
        return None
    failed_step = next((step for step in job.get("steps", [])
                        if step.get("conclusion") in ("failure", "cancelled", "timed_out")), None)
    evidence = []
    counters = {}
    for path, value in re.findall(r"(\S+/memory.events):[^\n]*?\boom_kill (\d+)", log):
        counters.setdefault(path, []).append(int(value))
    # Each ancestor reports the same kill: don't sum hierarchical counters.
    delta = max((max(values) - values[0] for values in counters.values()), default=0)
    messages = "\n".join(annotation.get("message", "") for annotation in annotations or []).lower()
    if delta > 0:
        failure_type = "out_of_memory"
        evidence.append(f"cgroup memory.events oom_kill increased by {delta} during this job")
    elif job["conclusion"] == "timed_out" or "exceeded the maximum execution time" in messages:
        failure_type = "timeout"
        evidence.append("GitHub reported a timeout")
    elif "lost communication with the server" in messages:
        failure_type = "runner_failure"
        evidence.append("GitHub reported lost runner communication; underlying cause is unconfirmed")
    elif job["conclusion"] == "cancelled":
        failure_type = "cancelled"
        evidence.append("GitHub marked the job cancelled; no OOM counter increase was observed")
    elif failed_step:
        failure_type = "step_failure" if failed_step["name"].lower().startswith("run ") else "infrastructure"
        evidence.append(f"Step conclusion: {failed_step.get('conclusion')}; underlying cause is unconfirmed")
    else:
        failure_type = "unknown"
        evidence.append("GitHub reported job failure without a failed step or conclusive resource evidence")
    return {
        "job_name": job.get("name", "Unknown"), "job_id": job.get("id"),
        "conclusion": job.get("conclusion"), "failure_type": failure_type,
        "failed_step": failed_step.get("name") if failed_step else "Unknown",
        "step_number": failed_step.get("number") if failed_step else None,
        "started_at": job.get("started_at"), "completed_at": job.get("completed_at"),
        "html_url": job.get("html_url"), "runner_name": job.get("runner_name"),
        "evidence": evidence,
    }


def detect_job_failures(owner: str, repo: str, run_id: str, output_file: str,
                        jobs_file: str = "workflow-jobs.json") -> None:
    # Fail visibly on an API error. Never emit [] as though the query succeeded.
    Path(output_file).unlink(missing_ok=True)
    Path(jobs_file).unlink(missing_ok=True)
    jobs = get_workflow_jobs(owner, repo, run_id)
    Path(jobs_file).write_text(json.dumps(jobs, indent=2), encoding="utf-8")
    failures = []
    for job in jobs:
        if job.get("conclusion") not in ("failure", "cancelled", "timed_out"):
            continue
        log, annotations, unavailable = "", [], []
        try:
            log = run_gh_command(["api", f"/repos/{owner}/{repo}/actions/jobs/{job['id']}/logs",
                                  "--allow-escape-sequences"])
        except subprocess.CalledProcessError:
            unavailable.append("Job logs unavailable")
        check_id = str(job.get("check_run_url", "")).rsplit("/", 1)[-1]
        if check_id.isdigit():
            try:
                pages = json.loads(run_gh_command([
                    "api", f"/repos/{owner}/{repo}/check-runs/{check_id}/annotations?per_page=100",
                    "--paginate", "--slurp",
                ]))
                annotations = [annotation for page in pages for annotation in page]
            except (subprocess.CalledProcessError, ValueError):
                unavailable.append("Runner annotations unavailable")
        failure = analyze_job_failure(job, log, annotations)
        failure["evidence"].extend(unavailable)
        failures.append(failure)
        print(f"{failure['job_name']}: {failure['failure_type']} at {failure['failed_step']}")
    Path(output_file).write_text(json.dumps(failures, indent=2), encoding="utf-8")
    print(f"Detected {len(failures)} job-level failures")


def main():
    if len(sys.argv) < 4:
        sys.exit("Usage: detect_job_failures.py OWNER REPO RUN_ID [FAILURES_JSON] [JOBS_JSON]")
    detect_job_failures(*sys.argv[1:4], sys.argv[4] if len(sys.argv) > 4 else "job-failures.json",
                        sys.argv[5] if len(sys.argv) > 5 else "workflow-jobs.json")


if __name__ == "__main__":
    main()
