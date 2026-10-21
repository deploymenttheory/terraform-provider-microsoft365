#!/usr/bin/env python3
"""Manages GitHub issues for test failures (create, update, close).

This script automates GitHub issue management for test failures by:
- Creating new issues for first-time test failures
- Updating existing issues when tests continue to fail
- Marking recurring failures with 'recurring' label
- Auto-closing issues when tests pass

Usage:
    ./manage_test_issues.py <owner> <repo> <run-id> <failures-json> [successes-json] [status-json]

Args:
    owner: GitHub repository owner.
    repo: GitHub repository name.
    run-id: GitHub Actions workflow run ID.
    failures-json: Path to merged test failures JSON file.
    successes-json: Optional path to merged test successes JSON file.
    status-json: Aggregation completeness report, required for automatic closure.
"""

import sys
import json
import re
import subprocess
import time
from datetime import datetime, timezone
from pathlib import Path
from typing import Optional


def run_gh_command(args: list[str], max_retries: int = 3) -> str:
    """Run a GitHub CLI command and return output with retry logic.

    Args:
        args: List of arguments to pass to 'gh' command.
        max_retries: Maximum number of retry attempts (default: 3).

    Returns:
        Command stdout as a string.

    Raises:
        subprocess.CalledProcessError: If the command fails after all retries.
    """
    for attempt in range(max_retries):
        try:
            result = subprocess.run(
                ["gh"] + args,
                capture_output=True,
                text=True,
                check=True
            )
            return result.stdout.strip()
        except subprocess.CalledProcessError as e:
            error_msg = e.stderr.lower()
            
            # Check if it's a rate limit or transient API error
            is_retryable = any(phrase in error_msg for phrase in [
                "graphql", "rate limit", "abuse", "timeout", 
                "temporarily unavailable", "server error", "502", "503"
            ])
            
            if is_retryable and attempt < max_retries - 1:
                wait_time = (2 ** attempt) * 2  # Exponential backoff: 2, 4, 8 seconds
                print(f"⚠️  API error (attempt {attempt + 1}/{max_retries}), retrying in {wait_time}s...", 
                      file=sys.stderr)
                print(f"   Error: {e.stderr[:200]}", file=sys.stderr)  # Show first 200 chars
                time.sleep(wait_time)
                continue
            
            # Not retryable or out of retries
            print(f"Error running gh command: {e.stderr}", file=sys.stderr)
            raise
    
    raise RuntimeError("Unexpected: exhausted retries without returning")


def ensure_label_exists(owner: str, repo: str, label_name: str, color: str, description: str) -> None:
    """Create a label if it doesn't exist, or update it if it does.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        label_name: Name of the label.
        color: Hex color code (without #).
        description: Label description.
    """
    # Simply try to create the label with --force flag to update if exists
    # This is simpler than checking first - gh handles the logic
    try:
        run_gh_command([
            "label", "create", label_name,
            "--repo", f"{owner}/{repo}",
            "--color", color,
            "--description", description,
            "--force"  # Update if exists, create if not
        ])
        print(f"  ✅ Ensured label exists: {label_name}")
    except subprocess.CalledProcessError as e:
        print(f"  ⚠️  Warning: Could not create/update label '{label_name}': {e}", file=sys.stderr)


def test_identity(result: dict) -> tuple[str, str] | None:
    if not isinstance(result, dict):
        return None
    package, name = result.get("package"), result.get("test_name")
    return (package, name) if isinstance(package, str) and package and isinstance(name, str) and name else None


def issue_identity(issue: dict) -> tuple[str, str] | None:
    """Legacy titles do not identify a package; leave those incidents open."""
    match = re.search(r"<!-- acceptance-test: (.+?) -->", issue.get("body", ""))
    if match:
        try:
            return test_identity(json.loads(match.group(1)))
        except (ValueError, AttributeError):
            pass
    return None


def get_all_open_test_issues(owner: str, repo: str) -> list[dict]:
    """Get all open test-failure issues.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.

    Returns:
        List of issue dictionaries with number, title and body fields.
    """
    try:
        result = run_gh_command([
            "issue", "list",
            "--repo", f"{owner}/{repo}",
            "--state", "open",
            "--label", "test-failure",
            "--json", "number,title,body",
            "--limit", "1000"
        ])
        
        if not result:
            return []
        
        return json.loads(result)
    except (subprocess.CalledProcessError, json.JSONDecodeError):
        # Do not create duplicates or close issues after a failed inventory query.
        raise


def update_existing_issue(owner: str, repo: str, issue_number: str, 
                         service_path: str, context: str, 
                         date: str, run_id: str, workflow_url: str) -> None:
    """Add a comment to existing issue with latest failure details.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        issue_number: Issue number to update.
        service_path: Service path (e.g., 'resources/identity_and_access').
        context: Error context from test output.
        date: Date string for the failure.
        run_id: GitHub Actions workflow run ID.
        workflow_url: URL to the workflow run.
    """
    # Get current timestamp
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")
    
    # Truncate context if too long (GitHub has limits on comment body size)
    max_context_length = 5000
    truncated_context = context
    if len(context) > max_context_length:
        truncated_context = context[:max_context_length] + "\n\n... (truncated, see workflow logs for full output)"
    
    comment_body = f"""## Still Failing

**Timestamp:** {timestamp}  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})  
**Service:** `{service_path}`

### Latest Error Output

```
{truncated_context}
```

---
*Automated update from nightly tests*"""
    
    # Add comment with latest failure details
    run_gh_command([
        "issue", "comment", issue_number,
        "--repo", f"{owner}/{repo}",
        "--body", comment_body
    ])
    
    # Mark as recurring failure
    run_gh_command([
        "issue", "edit", issue_number,
        "--repo", f"{owner}/{repo}",
        "--add-label", "recurring"
    ])


def close_resolved_issue(owner: str, repo: str, issue_number: str,
                        date: str, run_id: str, workflow_url: str,
                        test_passed: bool = False) -> None:
    """Close an issue only after an explicit matching pass.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        issue_number: Issue number to close.
        date: Date string for the resolution.
        run_id: GitHub Actions workflow run ID.
        workflow_url: URL to the workflow run.
        test_passed: Must be True after checking identity and reporting completeness.
    """
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")
    
    if not test_passed:
        return
    message = "This package and test explicitly passed in a complete report. Automatically closing this issue."
    
    run_gh_command([
        "issue", "close", issue_number,
        "--repo", f"{owner}/{repo}",
        "--comment", f"""## ✅ Resolved

**Timestamp:** {timestamp}  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})

{message}

---
*Automated closure from nightly tests*""",
        "--reason", "completed"
    ])


def create_new_issue(owner: str, repo: str, test_name: str, 
                    service_path: str, context: str, date: str, 
                    run_id: str, workflow_url: str, package: str) -> str:
    """Create a new issue for test failure.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        test_name: Name of the failed test.
        service_path: Service path (e.g., 'resources/identity_and_access').
        context: Error context from test output.
        date: Date string for the failure.
        run_id: GitHub Actions workflow run ID.
        workflow_url: URL to the workflow run.

    Returns:
        URL of the created issue.
    """
    issue_title = f"Bug: {test_name} Failing"
    
    # Truncate context if too long (GitHub has limits on issue body size)
    max_context_length = 5000
    truncated_context = context
    if len(context) > max_context_length:
        truncated_context = context[:max_context_length] + "\n\n... (truncated, see workflow logs for full output)"
    
    identity = json.dumps({"package": package, "test_name": test_name}, sort_keys=True)
    issue_body = f"""<!-- acceptance-test: {identity} -->

## Test Failure

**Test:** `{test_name}`  
**Package:** `{package}`<br>
**Service:** `{service_path}`  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})

### Error Output

```
{truncated_context}
```

### Resources

- [Workflow Logs]({workflow_url})
- [Test Source](https://github.com/{owner}/{repo}/tree/main/{package.removeprefix("./")})

---
*Automated report from nightly tests*"""
    
    issue_url = run_gh_command([
        "issue", "create",
        "--repo", f"{owner}/{repo}",
        "--title", issue_title,
        "--body", issue_body,
        "--label", "test-failure,automated"
    ])
    
    return issue_url


def process_test_failures(owner: str, repo: str, run_id: str,
                          failures_json_path: str, successes_json_path: Optional[str] = None,
                          status_json_path: Optional[str] = None) -> None:
    """Report failures; resolve only proven passes from a complete aggregation."""
    failures = json.loads(Path(failures_json_path).read_text(encoding="utf-8"))
    successes = json.loads(Path(successes_json_path).read_text(encoding="utf-8")) if successes_json_path else []
    complete = False
    if status_json_path:
        try:
            complete = json.loads(Path(status_json_path).read_text(encoding="utf-8")).get("complete") is True
        except (OSError, ValueError, AttributeError):
            pass
    if not isinstance(failures, list) or not isinstance(successes, list):
        raise ValueError("Expected test result arrays")
    if any(test_identity(result) is None for result in failures + successes):
        raise ValueError("Test result lacks package and test identity")
    print(f"Test results: {len(failures)} failures, {len(successes)} explicit passes; complete={complete}")
    if not failures and not successes:
        print("No test outcomes to process; existing incidents retained")
        return
    open_issues = get_all_open_test_issues(owner, repo)
    date = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    workflow_url = f"https://github.com/{owner}/{repo}/actions/runs/{run_id}"
    if failures:
        ensure_label_exists(owner, repo, "test-failure", "d73a4a", "Automated test failure report")
        ensure_label_exists(owner, repo, "automated", "0366d6", "Automatically generated")
        ensure_label_exists(owner, repo, "recurring", "b60205", "Test has failed multiple times")
    reported = set()
    for failure in failures:
        identity = test_identity(failure)
        if identity in reported:
            continue
        reported.add(identity)
        package, test_name = identity
        kind, service = failure["configuration_block_type"], failure["service"]
        service_path = f"{kind}/{service}" if service and service != "null" else kind
        existing = next((issue for issue in open_issues if issue_identity(issue) == identity), None)
        if existing:
            update_existing_issue(owner, repo, str(existing["number"]), service_path,
                                  failure["context"], date, run_id, workflow_url)
        else:
            create_new_issue(owner, repo, test_name, service_path, failure["context"],
                             date, run_id, workflow_url, package)
        time.sleep(1)

    if not complete:
        print("Reporting incomplete or unverified; automatic issue closure disabled")
        return
    failed = {test_identity(result) for result in failures}
    passed = {test_identity(result) for result in successes if result.get("package_complete") is True}
    for issue in open_issues:
        identity = issue_identity(issue)
        if identity is not None and identity in passed and identity not in failed:
            close_resolved_issue(owner, repo, str(issue["number"]), date, run_id, workflow_url, True)
            time.sleep(1)


def main():
    if len(sys.argv) < 4:
        print("Usage: manage_test_issues.py <owner> <repo> <run-id> [failures-json] [successes-json] [status-json]",
            file=sys.stderr)
        sys.exit(1)
    
    owner = sys.argv[1]
    repo = sys.argv[2]
    run_id = sys.argv[3]
    failures_json = sys.argv[4] if len(sys.argv) > 4 else "test-failures.json"
    successes_json = sys.argv[5] if len(sys.argv) > 5 else None
    
    status_json = sys.argv[6] if len(sys.argv) > 6 else None
    process_test_failures(owner, repo, run_id, failures_json, successes_json, status_json)


if __name__ == "__main__":
    main()
