"""OpenAPI spec fetching from msgraph-metadata repository."""

import requests
from typing import Optional, TYPE_CHECKING

if TYPE_CHECKING:
    from .progress_reporter import ProgressReporter


class SpecFetcher:
    """Fetches OpenAPI specs from msgraph-metadata repo."""
    
    SPEC_URL = "https://raw.githubusercontent.com/microsoftgraph/msgraph-metadata/master/openapi/beta/openapi.yaml"
    REPO_API = "https://api.github.com/repos/microsoftgraph/msgraph-metadata"
    
    def __init__(self, reporter: 'ProgressReporter'):
        """Initialize spec fetcher.
        
        Args:
            reporter: Progress reporter
        """
        self.reporter = reporter
    
    def fetch_latest_spec(self) -> str:
        """Download current openapi.yaml from master branch.
        
        Returns:
            OpenAPI spec content as string
        """
        self.reporter.info("📥 Downloading latest OpenAPI spec from master...")
        
        try:
            response = requests.get(self.SPEC_URL, timeout=60)
            response.raise_for_status()
            
            content = response.text
            size_mb = len(content) / (1024 * 1024)
            self.reporter.info(f"   Downloaded {size_mb:.2f} MB")
            
            return content
            
        except requests.RequestException as e:
            self.reporter.error(f"Failed to download spec: {e}")
            raise
    
    def fetch_spec_at_commit(self, commit_sha: str) -> str:
        """Download openapi.yaml from specific commit.
        
        Args:
            commit_sha: Git commit SHA
            
        Returns:
            OpenAPI spec content as string
        """
        url = f"https://raw.githubusercontent.com/microsoftgraph/msgraph-metadata/{commit_sha}/openapi/beta/openapi.yaml"
        
        self.reporter.info(f"📥 Downloading OpenAPI spec from commit {commit_sha[:8]}...")
        
        try:
            response = requests.get(url, timeout=60)
            response.raise_for_status()
            
            content = response.text
            size_mb = len(content) / (1024 * 1024)
            self.reporter.info(f"   Downloaded {size_mb:.2f} MB")
            
            return content
            
        except requests.RequestException as e:
            self.reporter.error(f"Failed to download spec at commit {commit_sha}: {e}")
            raise
    
    def get_latest_commit_sha(self, file_path: str = "openapi/beta/openapi.yaml") -> Optional[str]:
        """Get latest commit SHA for the OpenAPI spec file.
        
        Args:
            file_path: Path to file in repo
            
        Returns:
            Commit SHA or None
        """
        url = f"{self.REPO_API}/commits"
        params = {"path": file_path, "per_page": 1}
        
        try:
            response = requests.get(url, params=params, timeout=30)
            response.raise_for_status()
            
            commits = response.json()
            if commits and len(commits) > 0:
                return commits[0]['sha']
            
            return None
            
        except requests.RequestException as e:
            self.reporter.error(f"Failed to get latest commit: {e}")
            return None
