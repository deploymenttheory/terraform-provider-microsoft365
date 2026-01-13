"""GitHub CLI operations."""

import json
import subprocess
from typing import List, Tuple, TYPE_CHECKING

from regex_patterns import RegexPatterns  # type: ignore

if TYPE_CHECKING:
    from core.progress_reporter import ProgressReporter


class GitHubClient:
    """Handles all GitHub CLI operations."""

    def __init__(self, repo: str, reporter: 'ProgressReporter'):
        """Initialize GitHub client.
        
        Args:
            repo: Repository in owner/repo format
            reporter: Progress reporter for output
        """
        self.repo = repo
        self.reporter = reporter

    def run_command(self, args: List[str], check: bool = True) -> Tuple[str, str, int]:
        """Run a command and return output.
        
        Args:
            args: Command and arguments to run
            check: Whether to raise exception on non-zero exit
            
        Returns:
            Tuple of (stdout, stderr, returncode)
        """
        try:
            result = subprocess.run(
                args,
                capture_output=True,
                text=True,
                check=check
            )
            return result.stdout.strip(), result.stderr.strip(), result.returncode
        except subprocess.CalledProcessError as e:
            if check:
                self.reporter.error(f"Command failed: {' '.join(args)}")
                self.reporter.error(f"Error: {e.stderr}")
                raise
            return e.stdout.strip(), e.stderr.strip(), e.returncode

    def run_gh_command(self, args: List[str]) -> str:
        """Run a GitHub CLI command and return output.
        
        Args:
            args: Arguments to pass to 'gh' command
            
        Returns:
            Command stdout as string
        """
        stdout, _, _ = self.run_command(["gh"] + args)
        return stdout

    def get_pr_diff(self, pr_number: int) -> str:
        """Get diff from a pull request.
        
        Args:
            pr_number: PR number
            
        Returns:
            Diff text
        """
        return self.run_gh_command([
            "pr", "diff", str(pr_number),
            "--repo", self.repo
        ])

    def get_compare_data(self, sdk_repo: str, old_version: str, new_version: str) -> dict:
        """Get comparison data between two versions from GitHub API.
        
        Args:
            sdk_repo: SDK repository (owner/repo)
            old_version: Old version tag
            new_version: New version tag
            
        Returns:
            Parsed JSON response from GitHub API
        """
        json_result = self.run_gh_command([
            "api",
            f"/repos/{sdk_repo}/compare/{old_version}...{new_version}"
        ])
        
        if not json_result:
            raise ValueError("No response from GitHub API")
        
        return json.loads(json_result)

    def create_issue(self, title: str, body: str, labels: List[str]) -> str:
        """Create a GitHub issue.
        
        Args:
            title: Issue title
            body: Issue body
            labels: List of label names
            
        Returns:
            Issue URL
        """
        return self.run_gh_command([
            "issue", constants.TfOperationCreate,
            "--repo", self.repo,
            "--title", title,
            "--body", body,
            "--label", ",".join(labels)
        ])

    def ensure_labels(self, labels: List[Tuple[str, str, str]]):
        """Ensure labels exist in the repository.
        
        Args:
            labels: List of (name, color, description) tuples
        """
        for label_name, color, description in labels:
            try:
                self.run_gh_command([
                    "label", constants.TfOperationCreate, label_name,
                    "--repo", self.repo,
                    "--color", color,
                    "--description", description,
                    "--force"
                ])
            except subprocess.CalledProcessError:
                pass  # Label might already exist

    def fetch_changelog(self, url: str) -> str:
        """Fetch changelog from URL.
        
        Args:
            url: Changelog URL
            
        Returns:
            Changelog content
        """
        stdout, _, _ = self.run_command(["curl", "-s", url])
        return stdout

    def get_current_repo(self) -> str:
        """Get current repository from git remote.
        
        Returns:
            Repository in owner/repo format
        """
        try:
            result = subprocess.run(
                ["git", "remote", "get-url", "origin"],
                capture_output=True,
                text=True,
                check=True
            )
            remote_url = result.stdout.strip()
            match = RegexPatterns.GITHUB_REPO_URL.search(remote_url)
            if match:
                return match.group(1).rstrip('.git')
        except subprocess.CalledProcessError:
            pass
        return "deploymenttheory/terraform-provider-microsoft365"
