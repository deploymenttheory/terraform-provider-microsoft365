#!/usr/bin/env python3
"""Collect a completed discovery run using read-only GitHub API calls."""

import argparse
import io
import json
from pathlib import Path
import subprocess
import zipfile

from summarize import summarize_log

REPO = "deploymenttheory/terraform-provider-microsoft365"


def api(path):
    return subprocess.check_output(["gh", "api", f"repos/{REPO}/{path}"])


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("run_id", type=int)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    run = json.loads(api(f"actions/runs/{args.run_id}"))
    if run["status"] != "completed":
        raise SystemExit("Run is not complete; collect it after all experiments terminate.")
    jobs = json.loads(api(f"actions/runs/{args.run_id}/jobs?per_page=100"))["jobs"]
    # Read entries directly: never extract arbitrary archive paths.
    archive = zipfile.ZipFile(io.BytesIO(api(f"actions/runs/{args.run_id}/logs")))
    logs = {
        name.split("_", 1)[1].removesuffix(".txt"): archive.read(name).decode(errors="replace")
        for name in archive.namelist() if "/" not in name and "_" in name
    }
    results = []
    for job in jobs:
        if job["name"] == "matrix":
            continue
        log = logs.get(job["name"])
        result = summarize_log(log) if log else summarize_log("")
        replays = []
        for line in (log or "").splitlines():
            if "[VET REPLAY RESULT] " in line:
                replays.append(json.loads(line.split("[VET REPLAY RESULT] ", 1)[1]))
        result.update(
            name=job["name"], job_id=job["id"], conclusion=job["conclusion"],
            url=job["html_url"], log_available=log is not None,
            started_at_utc=job["started_at"], completed_at_utc=job["completed_at"],
            annotations=json.loads(api(f"check-runs/{job['id']}/annotations")),
            vet_replays=replays,
        )
        # Job failure remains a failure even if an earlier child reported zero.
        result["candidate_pass"] = result["candidate_pass"] and job["conclusion"] == "success"
        results.append(result)
    output = {
        "run_id": args.run_id, "harness_sha": run["head_sha"],
        "url": run["html_url"], "conclusion": run["conclusion"], "experiments": results,
    }
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(output, indent=2) + "\n")
    for result in results:
        print(json.dumps({
            key: result.get(key)
            for key in ("name", "conclusion", "exit_code", "oom_kill_delta",
                        "min_available_mib", "vet_invocations", "candidate_pass",
                        "elapsed_seconds", "log_available")
        }))


if __name__ == "__main__":
    main()
