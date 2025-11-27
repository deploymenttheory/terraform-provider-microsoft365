#!/usr/bin/env python3
"""Manages GitHub issues for test failures (create, update, close).

This script automates GitHub issue management for test failures by:
- Creating new issues for first-time test failures
- Updating existing issues when tests continue to fail
- Marking recurring failures with 'recurring' label
- Auto-closing issues when tests pass

Usage:
    ./manage-test-issues.py <owner> <repo> <run-id> <failures-json> [successes-json]

Args:
    owner: GitHub repository owner.
    repo: GitHub repository name.
    run-id: GitHub Actions workflow run ID.
    failures-json: Path to merged test failures JSON file.
    successes-json: Optional path to merged test successes JSON file.
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


def find_existing_issue(owner: str, repo: str, test_name: str) -> Optional[str]:
    """Check if an issue already exists for this test failure.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        test_name: Name of the failed test.

    Returns:
        Issue number as string if found, None otherwise.
    """
    try:
        result = run_gh_command([
            "issue", "list",
            "--repo", f"{owner}/{repo}",
            "--state", "open",
            "--label", "test-failure",
            "--search", f'in:title "{test_name}"',
            "--json", "number,title"
        ])
        
        if not result:
            return None
        
        # Parse JSON response in Python instead of using jq
        issues = json.loads(result)
        if issues and len(issues) > 0:
            return str(issues[0]["number"])
        
        return None
    except Exception:
        return None


def get_all_open_test_issues(owner: str, repo: str) -> list[dict]:
    """Get all open test-failure issues.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.

    Returns:
        List of issue dictionaries with number and title fields.
    """
    try:
        result = run_gh_command([
            "issue", "list",
            "--repo", f"{owner}/{repo}",
            "--state", "open",
            "--label", "test-failure",
            "--json", "number,title",
            "--limit", "1000"
        ])
        
        if not result:
            return []
        
        return json.loads(result)
    except Exception:
        return []


def update_existing_issue(owner: str, repo: str, issue_number: str, 
                         test_name: str, service_path: str, context: str, 
                         date: str, run_id: str, workflow_url: str) -> None:
    """Add a comment to existing issue with latest failure details.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        issue_number: Issue number to update.
        test_name: Name of the failed test.
        service_path: Service path (e.g., 'resources/identity_and_access').
        context: Error context from test output.
        date: Date string for the failure.
        run_id: GitHub Actions workflow run ID.
        workflow_url: URL to the workflow run.
    """
    # Get current timestamp
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")
    
    comment_body = f"""## Still Failing

**Timestamp:** {timestamp}  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})  
**Service:** `{service_path}`

### Latest Error Output

```
{context}
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
                        test_name: str, date: str, run_id: str, workflow_url: str) -> None:
    """Close an issue when test is now passing.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        issue_number: Issue number to close.
        test_name: Name of the test that is now passing.
        date: Date string for the resolution.
        run_id: GitHub Actions workflow run ID.
        workflow_url: URL to the workflow run.
    """
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")
    
    run_gh_command([
        "issue", "close", issue_number,
        "--repo", f"{owner}/{repo}",
        "--comment", f"""## ✅ Resolved

**Timestamp:** {timestamp}  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})

Test is now passing. Automatically closing this issue.

---
*Automated closure from nightly tests*""",
        "--reason", "completed"
    ])


def create_new_issue(owner: str, repo: str, test_name: str, 
                    service_path: str, context: str, date: str, 
                    run_id: str, workflow_url: str) -> str:
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
    
    issue_body = f"""## Test Failure

**Test:** `{test_name}`  
**Service:** `{service_path}`  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})

### Error Output

```
{context}
```

### Resources

- [Workflow Logs]({workflow_url})
- [Test Source](../../internal/services/{service_path.replace('-', '_')})

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
                        failures_json_path: str, successes_json_path: Optional[str] = None) -> None:
    """Process test failures and successes, create/update/close issues as needed.

    Main processing function that handles the complete lifecycle of test failure issues:
    - Creates new issues for first-time failures
    - Updates existing issues for recurring failures
    - Closes issues for tests that are now passing

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        run_id: GitHub Actions workflow run ID.
        failures_json_path: Path to JSON file containing test failures.
        successes_json_path: Optional path to JSON file containing test successes.
    """
    if not all([owner, repo, run_id]):
        print("Usage: manage-test-issues.py <owner> <repo> <run-id> <failures-json> [successes-json]", 
            file=sys.stderr)
        sys.exit(1)
    
    failures_path = Path(failures_json_path)
    if not failures_path.exists():
        print(f"Error: Failures JSON file not found: {failures_json_path}", 
            file=sys.stderr)
        sys.exit(1)
    
    with open(failures_path) as f:
        failures = json.load(f)
    
    successes = []
    if successes_json_path:
        successes_path = Path(successes_json_path)
        if successes_path.exists():
            with open(successes_path) as f:
                successes = json.load(f)
    
    failure_count = len(failures)
    success_count = len(successes)
    
    print(f"\nTest results: {failure_count} failures, {success_count} successes")
    
    if failure_count == 0 and success_count == 0:
        print("✅ No test results to process")
        return

    print("Checking required labels...")
    ensure_label_exists(owner, repo, "test-failure", "d73a4a", "Automated test failure report")
    ensure_label_exists(owner, repo, "automated", "0366d6", "Automatically generated")
    ensure_label_exists(owner, repo, "recurring", "b60205", "Test has failed multiple times")
    
    print(f"\n{'='*60}")
    print(f"Creating GitHub issues for {failure_count} test failure(s)")
    print(f"{'='*60}\n")
    
    date = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    workflow_url = f"https://github.com/{owner}/{repo}/actions/runs/{run_id}"
    
    created_count = 0
    updated_count = 0
    
    for failure in failures:
        test_name = failure["test_name"]
        category = failure["category"]
        service = failure["service"]
        context = failure["context"]
        
        service_path = f"{category}/{service}" if service and service != "null" else category
        
        print(f"• {test_name}")
        print(f"  Service: {service_path}")
        
        existing_issue = find_existing_issue(owner, repo, test_name)
        
        if existing_issue:
            print(f"  Action: Updated existing issue #{existing_issue}")
            update_existing_issue(
                owner, repo, existing_issue, test_name, 
                service_path, context, date, run_id, workflow_url
            )
            updated_count += 1
        else:
            issue_url = create_new_issue(
                owner, repo, test_name, service_path, 
                context, date, run_id, workflow_url
            )
            print(f"  Action: Created new issue → {issue_url}")
            created_count += 1
        print()
    
    print(f"{'='*60}")
    print(f"Issue updates: {created_count} created, {updated_count} updated")
    print(f"{'='*60}")
    
    if successes:
        print(f"\n{'='*60}")
        print("Checking for resolved issues to close")
        print(f"{'='*60}\n")
        
        open_issues = get_all_open_test_issues(owner, repo)
        
        failed_test_names = {f["test_name"] for f in failures}
        passed_test_names = {s["test_name"] for s in successes}
        
        closed_count = 0
        
        for issue in open_issues:
            issue_title = issue["title"]
            issue_number = str(issue["number"])
            
            test_name = issue_title
            if issue_title.startswith("Bug: ") and issue_title.endswith(" Failing"):
                test_name = issue_title[5:-8]  # Remove "Bug: " prefix and " Failing" suffix
            
            if test_name not in failed_test_names and test_name in passed_test_names:
                print(f"• {issue_title}")
                print(f"  Action: Closing resolved issue #{issue_number}")
                close_resolved_issue(
                    owner, repo, issue_number, 
                    test_name, date, run_id, workflow_url
                )
                closed_count += 1
                print()
        
        if closed_count > 0:
            print(f"{'='*60}")
            print(f"Closed {closed_count} resolved issue(s)")
            print(f"{'='*60}")
        else:
            print("No issues to close - all open issues still failing")


def main():
    if len(sys.argv) < 4:
        print("Usage: manage-test-issues.py <owner> <repo> <run-id> [failures-json] [successes-json]", 
            file=sys.stderr)
        sys.exit(1)
    
    owner = sys.argv[1]
    repo = sys.argv[2]
    run_id = sys.argv[3]
    failures_json = sys.argv[4] if len(sys.argv) > 4 else "test-failures.json"
    successes_json = sys.argv[5] if len(sys.argv) > 5 else None
    
    process_test_failures(owner, repo, run_id, failures_json, successes_json)


if __name__ == "__main__":
    main()
