#!/usr/bin/env python3
"""Manages GitHub issues for test failures (create, update, close).

This script automates GitHub issue management for test failures by:
- Creating new issues for first-time test failures
- Updating existing issues when tests continue to fail
- Marking recurring failures with 'recurring' label
- Auto-closing matching incidents only after scoped verification and freshness checks

Usage:
    ./manage_test_issues.py <owner> <repo> <run-id> <failures-json> [successes-json] [status-json]

Args:
    owner: GitHub repository owner.
    repo: GitHub repository name.
    run-id: GitHub Actions workflow run ID.
    failures-json: Path to merged test failures JSON file.
    successes-json: Optional path to merged test successes JSON file.
    status-json: Scoped package/test verification report, required for automatic closure.
"""

import json
import os
import re
import subprocess
import sys
import time
from dataclasses import dataclass, field
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
            result = subprocess.run(["gh"] + args, capture_output=True, text=True, check=True)
            return result.stdout.strip()
        except subprocess.CalledProcessError as e:
            error_msg = e.stderr.lower()

            # Check if it's a rate limit or transient API error
            is_retryable = any(
                phrase in error_msg
                for phrase in [
                    "graphql",
                    "rate limit",
                    "abuse",
                    "timeout",
                    "temporarily unavailable",
                    "server error",
                    "502",
                    "503",
                ]
            )

            if is_retryable and attempt < max_retries - 1:
                wait_time = (2**attempt) * 2  # Exponential backoff: 2, 4, 8 seconds
                print(
                    f"⚠️  API error (attempt {attempt + 1}/{max_retries}), retrying in"
                    f" {wait_time}s...",
                    file=sys.stderr,
                )
                print(f"   Error: {e.stderr[:200]}", file=sys.stderr)  # Show first 200 chars
                time.sleep(wait_time)
                continue

            # Not retryable or out of retries
            print(f"Error running gh command: {e.stderr}", file=sys.stderr)
            raise

    raise RuntimeError("Unexpected: exhausted retries without returning")


def ensure_label_exists(
    owner: str, repo: str, label_name: str, color: str, description: str
) -> None:
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
        run_gh_command(
            [
                "label",
                "create",
                label_name,
                "--repo",
                f"{owner}/{repo}",
                "--color",
                color,
                "--description",
                description,
                "--force",  # Update if exists, create if not
            ]
        )
        print(f"  ✅ Ensured label exists: {label_name}")
    except subprocess.CalledProcessError as e:
        print(f"  ⚠️  Warning: Could not create/update label '{label_name}': {e}", file=sys.stderr)


def test_identity(result: dict) -> tuple[str, str] | None:
    """Return the package and test name that identify a failure."""
    if not isinstance(result, dict):
        return None
    package, name = result.get("package"), result.get("test_name")
    return (
        (package, name)
        if isinstance(package, str) and package and isinstance(name, str) and name
        else None
    )


def issue_identity(issue: dict) -> tuple[str, str] | None:
    """Match structured identities or exact legacy body fields, never a title alone."""
    body = issue.get("body") or ""
    if not isinstance(body, str):
        return None
    if "<!-- acceptance-test:" in body:
        matches = re.findall(r"<!-- acceptance-test: (.+?) -->", body)
        if len(matches) != 1:
            return None
        try:
            return test_identity(json.loads(matches[0]))
        except (ValueError, AttributeError):
            pass
        return None
    packages = re.findall(r"^\*\*Package:\*\*[ \t]*`([^`\r\n]+)`", body, re.MULTILINE)
    names = re.findall(r"^\*\*Test:\*\*[ \t]*`([^`\r\n]+)`", body, re.MULTILINE)
    if len(packages) == len(names) == 1:
        return test_identity({"package": packages[0], "test_name": names[0]})
    return None


def get_all_open_test_issues(owner: str, repo: str) -> list[dict]:
    """Get all open test-failure issues.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.

    Returns:
        List of issue dictionaries with number, title and body fields.
    """
    result = run_gh_command(
        [
            "issue",
            "list",
            "--repo",
            f"{owner}/{repo}",
            "--state",
            "open",
            "--label",
            "test-failure",
            "--json",
            "number,title,body,updatedAt",
            "--limit",
            "1000",
        ]
    )

    if not result:
        return []

    return json.loads(result)


def _utc_timestamp(value: str) -> datetime:
    """Parse a timestamp with its explicit timezone before comparing report freshness."""
    timestamp = datetime.fromisoformat(value.replace("Z", "+00:00"))
    if timestamp.tzinfo is None:
        raise ValueError("timestamp has no timezone")
    return timestamp


def _run_started_at(owner: str, repo: str, run_id: str) -> datetime | None:
    """Get the current attempt's start time; unavailable provenance disables closure."""
    attempt = os.environ.get("GITHUB_RUN_ATTEMPT")
    if not attempt:
        print("Cannot verify workflow freshness without an explicit run attempt")
        return None
    endpoint = f"repos/{owner}/{repo}/actions/runs/{run_id}/attempts/{attempt}"
    try:
        run = json.loads(run_gh_command(["api", endpoint]))
        if str(run.get("id")) != str(run_id) or str(run.get("run_attempt")) != attempt:
            raise ValueError("workflow attempt identity does not match")
        return _utc_timestamp(run["run_started_at"])
    except (
        subprocess.CalledProcessError,
        OSError,
        ValueError,
        KeyError,
        TypeError,
        AttributeError,
    ) as error:
        print(f"Cannot verify workflow freshness; automatic closure disabled: {error}")
        return None


def _issue_predates_run(owner: str, repo: str, issue: dict, started_at: datetime) -> bool:
    """Do not let an older passing run close an incident updated after that run began."""
    try:
        if _utc_timestamp(issue["updatedAt"]) > started_at:
            return False
        updated_at = run_gh_command(
            [
                "issue",
                "view",
                str(issue["number"]),
                "--repo",
                f"{owner}/{repo}",
                "--json",
                "updatedAt",
                "--jq",
                ".updatedAt",
            ]
        )
        return _utc_timestamp(updated_at) <= started_at
    except (
        subprocess.CalledProcessError,
        OSError,
        ValueError,
        KeyError,
        TypeError,
        AttributeError,
    ) as error:
        print(f"Cannot verify incident freshness; leaving issue open: {error}")
        return False


def update_existing_issue(
    owner: str,
    repo: str,
    issue_number: str,
    service_path: str,
    context: str,
    date: str,
    run_id: str,
    workflow_url: str,
) -> None:
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
        truncated_context = (
            context[:max_context_length] + "\n\n... (truncated, see workflow logs for full output)"
        )

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
    run_gh_command(
        ["issue", "comment", issue_number, "--repo", f"{owner}/{repo}", "--body", comment_body]
    )

    # Mark as recurring failure
    run_gh_command(
        ["issue", "edit", issue_number, "--repo", f"{owner}/{repo}", "--add-label", "recurring"]
    )


def close_resolved_issue(
    owner: str,
    repo: str,
    issue_number: str,
    date: str,
    run_id: str,
    workflow_url: str,
    test_passed: bool = False,
) -> None:
    """Close an issue only after an explicit matching pass.

    Args:
        owner: GitHub repository owner.
        repo: GitHub repository name.
        issue_number: Issue number to close.
        date: Date string for the resolution.
        run_id: GitHub Actions workflow run ID.
        workflow_url: URL to the workflow run.
        test_passed: Must be True after checking scoped evidence and incident freshness.
    """
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")

    if not test_passed:
        return
    message = (
        "This exact package and test passed in a normally completed package from the verified"
        " workflow attempt. Automatically closing this issue."
    )

    run_gh_command(
        [
            "issue",
            "close",
            issue_number,
            "--repo",
            f"{owner}/{repo}",
            "--comment",
            f"""## ✅ Resolved

**Timestamp:** {timestamp}  
**Date:** {date}  
**Workflow:** [{run_id}]({workflow_url})

{message}

---
*Automated closure from nightly tests*""",
            "--reason",
            "completed",
        ]
    )


def create_new_issue(
    owner: str,
    repo: str,
    test_name: str,
    service_path: str,
    context: str,
    date: str,
    run_id: str,
    workflow_url: str,
    package: str,
) -> str:
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
        truncated_context = (
            context[:max_context_length] + "\n\n... (truncated, see workflow logs for full output)"
        )

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

    issue_url = run_gh_command(
        [
            "issue",
            "create",
            "--repo",
            f"{owner}/{repo}",
            "--title",
            issue_title,
            "--body",
            issue_body,
            "--label",
            "test-failure,automated",
        ]
    )

    return issue_url


@dataclass(frozen=True)
class WorkflowRun:
    """Share one workflow attempt's repository and reporting metadata."""

    owner: str
    repo: str
    run_id: str
    date: str = field(default_factory=lambda: datetime.now(timezone.utc).strftime("%Y-%m-%d"))

    @property
    def url(self) -> str:
        """Link incident updates to their originating workflow run."""
        return f"https://github.com/{self.owner}/{self.repo}/actions/runs/{self.run_id}"


def _read_outcomes(
    failures_json_path: str, successes_json_path: Optional[str]
) -> tuple[list, list]:
    """Validate result identities before any external issue updates."""
    failures = json.loads(Path(failures_json_path).read_text(encoding="utf-8"))
    successes = (
        json.loads(Path(successes_json_path).read_text(encoding="utf-8"))
        if successes_json_path
        else []
    )
    if not isinstance(failures, list) or not isinstance(successes, list):
        raise ValueError("Expected test result arrays")
    if any(test_identity(result) is None for result in failures + successes):
        raise ValueError("Test result lacks package and test identity")
    if any(
        not isinstance(failure.get(key), str)
        for failure in failures
        for key in ("context", "configuration_block_type", "service")
    ):
        raise ValueError("Failure result lacks reporting metadata")
    return failures, successes


def _read_status(status_json_path: Optional[str]) -> dict:
    """Unavailable verification evidence never authorizes automatic closure."""
    if status_json_path:
        try:
            status = json.loads(Path(status_json_path).read_text(encoding="utf-8"))
            if isinstance(status, dict):
                return status
        except (OSError, ValueError):
            pass
    return {}


def _verified_identities(status: dict, run_id: str) -> set[tuple[str, str]]:
    """Accept scoped passes only for the explicitly selected run and attempt."""
    provenance = {
        "run_id": os.environ.get("GITHUB_RUN_ID", ""),
        "run_attempt": os.environ.get("GITHUB_RUN_ATTEMPT", ""),
    }
    if not all(provenance.values()) or provenance["run_id"] != str(run_id):
        return set()
    if status.get("schema_version") != 2 or any(
        status.get(key) != expected for key, expected in provenance.items()
    ):
        return set()
    verified = status.get("verified_passes")
    if not isinstance(verified, list) or any(test_identity(item) is None for item in verified):
        return set()
    return {test_identity(item) for item in verified}


def _report_failures(run: WorkflowRun, failures: list, open_issues: list) -> None:
    """Create or update each observed current failure once, keyed by package and test."""
    if not failures:
        return
    for label, color, description in (
        ("test-failure", "d73a4a", "Automated test failure report"),
        ("automated", "0366d6", "Automatically generated"),
        ("recurring", "b60205", "Test has failed multiple times"),
    ):
        ensure_label_exists(run.owner, run.repo, label, color, description)
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
            update_existing_issue(
                run.owner,
                run.repo,
                str(existing["number"]),
                service_path,
                failure["context"],
                run.date,
                run.run_id,
                run.url,
            )
        else:
            create_new_issue(
                run.owner,
                run.repo,
                test_name,
                service_path,
                failure["context"],
                run.date,
                run.run_id,
                run.url,
                package,
            )
        time.sleep(1)


def _close_recoveries(
    run: WorkflowRun, failures: list, successes: list, status: dict, open_issues: list
) -> None:
    """Resolve verified passes only when their workflow began after the last incident update."""
    if not status.get("complete"):
        print(
            "Reporting incomplete; only independently verified package/test passes can resolve"
            " incidents"
        )
    failed = {test_identity(result) for result in failures}
    verified_identities = _verified_identities(status, run.run_id)
    passed = {
        test_identity(result)
        for result in successes
        if result.get("package_complete") is True and test_identity(result) in verified_identities
    }
    candidates = [
        issue
        for issue in open_issues
        if issue_identity(issue) in passed and issue_identity(issue) not in failed
    ]
    started_at = _run_started_at(run.owner, run.repo, run.run_id) if candidates else None
    for issue in candidates:
        identity = issue_identity(issue)
        if (
            identity is not None
            and started_at
            and _issue_predates_run(run.owner, run.repo, issue, started_at)
        ):
            close_resolved_issue(
                run.owner, run.repo, str(issue["number"]), run.date, run.run_id, run.url, True
            )
            time.sleep(1)


def process_test_failures(
    owner: str,
    repo: str,
    run_id: str,
    failures_json_path: str,
    successes_json_path: Optional[str] = None,
    status_json_path: Optional[str] = None,
) -> None:
    """Report observed failures and resolve independently verified, sufficiently recent passes."""
    failures, successes = _read_outcomes(failures_json_path, successes_json_path)
    status = _read_status(status_json_path)
    print(
        f"Test results: {len(failures)} failures, {len(successes)} explicit passes;"
        f" report_complete={status.get('complete') is True}"
    )
    if not failures and not successes:
        print("No test outcomes to process; existing incidents retained")
        return
    run = WorkflowRun(owner, repo, str(run_id))
    open_issues = get_all_open_test_issues(owner, repo)
    _report_failures(run, failures, open_issues)
    _close_recoveries(run, failures, successes, status, open_issues)


def main():
    """Reconcile acceptance test failures with GitHub issues."""
    if len(sys.argv) < 4:
        print(
            "Usage: manage_test_issues.py <owner> <repo> <run-id> [failures-json] [successes-json]"
            " [status-json]",
            file=sys.stderr,
        )
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
