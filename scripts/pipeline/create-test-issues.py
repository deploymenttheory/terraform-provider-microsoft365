#!/usr/bin/env python3
"""
Creates individual GitHub issues for each failing test.
Usage: ./create-test-issues.py <owner> <repo> <run-id> <failures-json>
"""

import sys
import json
import subprocess
from datetime import datetime
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


def find_existing_issue(owner: str, repo: str, test_name: str) -> Optional[str]:
    """Check if an issue already exists for this test failure."""
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


def update_existing_issue(owner: str, repo: str, issue_number: str, 
                         test_name: str, service_path: str, context: str, 
                         date: str, run_id: str, workflow_url: str) -> None:
    """Add a comment to existing issue with latest failure details."""
    comment_body = f"""## Failure Recurrence: {date}

**Workflow Run:** [{run_id}]({workflow_url})  
**Service:** `{service_path}`

### Error Output

```
{context}
```

---
*Automated update from nightly tests*"""
    
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


def create_new_issue(owner: str, repo: str, test_name: str, 
                    service_path: str, context: str, date: str, 
                    run_id: str, workflow_url: str) -> str:
    """Create a new issue for test failure."""
    issue_title = test_name
    
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
                         failures_json_path: str) -> None:
    """Process test failures and create or update issues."""
    if not all([owner, repo, run_id]):
        print("Usage: create-test-issues.py <owner> <repo> <run-id> <failures-json>", 
              file=sys.stderr)
        sys.exit(1)
    
    failures_path = Path(failures_json_path)
    if not failures_path.exists():
        print(f"Error: Failures JSON file not found: {failures_json_path}", 
              file=sys.stderr)
        sys.exit(1)
    
    with open(failures_path) as f:
        failures = json.load(f)
    
    failure_count = len(failures)
    
    if failure_count == 0:
        print("✅ No test failures to process")
        return
    
    print(f"\n{'='*60}")
    print(f"Creating GitHub issues for {failure_count} test failure(s)")
    print(f"{'='*60}\n")
    
    date = datetime.utcnow().strftime("%Y-%m-%d")
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
    print(f"Summary: {created_count} created, {updated_count} updated")
    print(f"{'='*60}")


def main():
    if len(sys.argv) < 4:
        print("Usage: create-test-issues.py <owner> <repo> <run-id> [failures-json]", 
              file=sys.stderr)
        sys.exit(1)
    
    owner = sys.argv[1]
    repo = sys.argv[2]
    run_id = sys.argv[3]
    failures_json = sys.argv[4] if len(sys.argv) > 4 else "test-failures.json"
    
    process_test_failures(owner, repo, run_id, failures_json)


if __name__ == "__main__":
    main()
