#!/usr/bin/env python3
"""
Manages GitHub issues for job-level failures.
Usage: ./manage-job-failure-issues.py <owner> <repo> <run-id> <job-failures-json>

Creates issues for:
- Infrastructure failures
- Job timeouts
- Runner failures
- OOM errors
- Step failures that prevented tests from running
"""

import sys
import json
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Optional


def run_gh_command(args: list[str]) -> str:
    """Run a GitHub CLI command and return output."""
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


def ensure_label_exists(owner: str, repo: str, label_name: str, color: str, description: str) -> None:
    """Create a label if it doesn't exist."""
    try:
        run_gh_command([
            "label", "create", label_name,
            "--repo", f"{owner}/{repo}",
            "--color", color,
            "--description", description,
            "--force"
        ])
        print(f"  ✅ Ensured label exists: {label_name}")
    except subprocess.CalledProcessError as e:
        print(f"  ⚠️  Warning: Could not create/update label '{label_name}': {e}", file=sys.stderr)


def find_existing_issue(owner: str, repo: str, job_name: str) -> Optional[str]:
    """Check if an issue already exists for this job failure."""
    try:
        # Search for issues with the job name
        result = run_gh_command([
            "issue", "list",
            "--repo", f"{owner}/{repo}",
            "--state", "open",
            "--label", "job-failure",
            "--search", f'in:title "{job_name}"',
            "--json", "number,title"
        ])
        
        if not result:
            return None
        
        issues = json.loads(result)
        if issues and len(issues) > 0:
            return str(issues[0]["number"])
        
        return None
    except Exception:
        return None


def create_job_failure_issue(owner: str, repo: str, failure: dict, 
                            date: str, run_id: str, workflow_url: str) -> str:
    """Create a new issue for job failure."""
    job_name = failure["job_name"]
    failure_type = failure["failure_type"].replace("_", " ").title()
    failed_step = failure.get("failed_step", "Unknown")
    
    # Map failure types to severity emoji
    severity_map = {
        "timeout": "⏱️",
        "out_of_memory": "💥",
        "runner_failure": "🔧",
        "infrastructure": "🚨",
        "cancelled": "🛑",
        "step_failure": "❌",
        "unknown": "❓"
    }
    emoji = severity_map.get(failure["failure_type"], "❌")
    
    issue_title = f"{emoji} Job Failure: {job_name}"
    
    issue_body = f"""## Job-Level Failure

**Job:** `{job_name}`  
**Failure Type:** {failure_type}  
**Failed Step:** `{failed_step}`  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})

### Details

- **Job ID:** {failure['job_id']}
- **Conclusion:** {failure['conclusion']}
- **Runner:** {failure.get('runner_name', 'Unknown')}
- **Started:** {failure.get('started_at', 'Unknown')}
- **Completed:** {failure.get('completed_at', 'Unknown')}

### Job Logs

[View Job Logs]({failure['html_url']})

### Possible Causes

"""
    
    # Add specific troubleshooting based on failure type
    if failure["failure_type"] == "timeout":
        issue_body += """
- Job exceeded maximum execution time
- Tests may be running too slowly
- Consider increasing timeout or optimizing tests
"""
    elif failure["failure_type"] == "out_of_memory":
        issue_body += """
- Job ran out of memory (OOM)
- Consider reducing parallel test execution
- May need larger runner or memory optimization
- Check if `-race` flag is causing excessive memory usage
"""
    elif failure["failure_type"] == "runner_failure":
        issue_body += """
- GitHub Actions runner infrastructure issue
- May be transient - retry workflow
- Check GitHub Status: https://www.githubstatus.com/
"""
    elif failure["failure_type"] == "infrastructure":
        issue_body += """
- Setup or dependency installation failed
- Check if external services are available
- May be network connectivity issue
"""
    else:
        issue_body += """
- Check job logs for specific error details
- May require manual investigation
"""
    
    issue_body += """
---
*Automated report from nightly test pipeline*"""
    
    issue_url = run_gh_command([
        "issue", "create",
        "--repo", f"{owner}/{repo}",
        "--title", issue_title,
        "--body", issue_body,
        "--label", "job-failure,automated,infrastructure"
    ])
    
    return issue_url


def update_job_failure_issue(owner: str, repo: str, issue_number: str,
                            failure: dict, date: str, run_id: str, workflow_url: str) -> None:
    """Add a comment to existing job failure issue."""
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")
    failure_type = failure["failure_type"].replace("_", " ").title()
    
    comment_body = f"""## Failure Recurrence

**Timestamp:** {timestamp}  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})  
**Type:** {failure_type}  
**Failed Step:** {failure.get('failed_step', 'Unknown')}

[View Job Logs]({failure['html_url']})

---
*Automated update from nightly test pipeline*"""
    
    run_gh_command([
        "issue", "comment", issue_number,
        "--repo", f"{owner}/{repo}",
        "--body", comment_body
    ])
    
    run_gh_command([
        "issue", "edit", issue_number,
        "--repo", f"{owner}/{repo}",
        "--add-label", "recurring"
    ])


def process_job_failures(owner: str, repo: str, run_id: str, failures_json_path: str) -> None:
    """Process job failures and create or update issues."""
    if not all([owner, repo, run_id]):
        print("Usage: manage-job-failure-issues.py <owner> <repo> <run-id> <job-failures-json>", 
            file=sys.stderr)
        sys.exit(1)
    
    failures_path = Path(failures_json_path)
    if not failures_path.exists():
        print(f"Error: Job failures JSON file not found: {failures_json_path}", 
            file=sys.stderr)
        sys.exit(1)
    
    with open(failures_path) as f:
        failures = json.load(f)
    
    failure_count = len(failures)
    
    if failure_count == 0:
        print("✅ No job-level failures to process")
        return
    
    # Ensure required labels exist
    print("Checking required labels...")
    ensure_label_exists(owner, repo, "job-failure", "8B0000", "Job-level infrastructure failure")
    ensure_label_exists(owner, repo, "infrastructure", "FFA500", "Infrastructure or runner issue")
    ensure_label_exists(owner, repo, "automated", "0366d6", "Automatically generated")
    ensure_label_exists(owner, repo, "recurring", "b60205", "Issue has occurred multiple times")
    
    print(f"\n{'='*60}")
    print(f"Processing {failure_count} job-level failure(s)")
    print(f"{'='*60}\n")
    
    date = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    workflow_url = f"https://github.com/{owner}/{repo}/actions/runs/{run_id}"
    
    created_count = 0
    updated_count = 0
    
    for failure in failures:
        job_name = failure["job_name"]
        
        print(f"• {job_name}")
        print(f"  Type: {failure['failure_type']}")
        
        existing_issue = find_existing_issue(owner, repo, job_name)
        
        if existing_issue:
            print(f"  Action: Updated existing issue #{existing_issue}")
            update_job_failure_issue(
                owner, repo, existing_issue, failure,
                date, run_id, workflow_url
            )
            updated_count += 1
        else:
            issue_url = create_job_failure_issue(
                owner, repo, failure, date, run_id, workflow_url
            )
            print(f"  Action: Created new issue → {issue_url}")
            created_count += 1
        print()
    
    print(f"{'='*60}")
    print(f"Summary: {created_count} created, {updated_count} updated")
    print(f"{'='*60}")


def main():
    if len(sys.argv) < 4:
        print("Usage: manage-job-failure-issues.py <owner> <repo> <run-id> [job-failures-json]", 
            file=sys.stderr)
        sys.exit(1)
    
    owner = sys.argv[1]
    repo = sys.argv[2]
    run_id = sys.argv[3]
    failures_json = sys.argv[4] if len(sys.argv) > 4 else "job-failures.json"
    
    process_job_failures(owner, repo, run_id, failures_json)


if __name__ == "__main__":
    main()

