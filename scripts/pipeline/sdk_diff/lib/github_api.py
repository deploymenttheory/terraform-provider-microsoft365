#!/usr/bin/env python3
"""GitHub API client for SDK version and change detection.

Provides functions to:
- Get latest SDK releases from GitHub
- Compare versions between tags
- Extract changelog information
"""

import json
import os
import re
import urllib.request
import urllib.error
from typing import Dict, List, Optional



def _make_github_request(url: str, token: Optional[str] = None) -> Dict:
    """Make authenticated GitHub API request.
    
    Args:
        url: GitHub API URL
        token: Optional GitHub token (defaults to GITHUB_TOKEN env var)
        
    Returns:
        Parsed JSON response
    """
    if token is None:
        token = os.environ.get('GITHUB_TOKEN')
    
    headers = {
        'Accept': 'application/vnd.github.v3+json'
    }
    
    if token:
        headers['Authorization'] = f'token {token}'
    
    req = urllib.request.Request(url, headers=headers)
    
    try:
        with urllib.request.urlopen(req) as response:
            return json.loads(response.read().decode())
    except urllib.error.HTTPError as e:
        error_body = e.read().decode()
        raise RuntimeError(f"GitHub API request failed ({e.code}): {error_body}") from e


def get_latest_release(repo: str) -> Dict[str, str]:
    """Get latest release information for a GitHub repository.
    
    Args:
        repo: Repository in format "owner/repo"
        
    Returns:
        Dictionary with release information:
        {
            "tag": "v1.93.0",
            "name": "Release v1.93.0",
            "published_at": "2026-01-15T10:00:00Z",
            "body": "Release notes..."
        }
    """
    url = f"https://api.github.com/repos/{repo}/releases/latest"
    data = _make_github_request(url)
    
    return {
        "tag": data.get("tag_name", ""),
        "name": data.get("name", ""),
        "published_at": data.get("published_at", ""),
        "body": data.get("body", "")
    }


def get_all_releases(repo: str, limit: int = 10) -> List[Dict[str, str]]:
    """Get recent releases for a GitHub repository.
    
    Args:
        repo: Repository in format "owner/repo"
        limit: Maximum number of releases to return
        
    Returns:
        List of release dictionaries
    """
    url = f"https://api.github.com/repos/{repo}/releases?per_page={limit}"
    releases = _make_github_request(url)
    
    return [
        {
            "tag": r.get("tag_name", ""),
            "name": r.get("name", ""),
            "published_at": r.get("published_at", ""),
            "body": r.get("body", "")
        }
        for r in releases
    ]


def compare_versions(repo: str, base: str, head: str) -> Dict[str, any]:
    """Compare two versions of a repository.
    
    Args:
        repo: Repository in format "owner/repo"
        base: Base version tag (e.g., "v0.150.0")
        head: Head version tag (e.g., "v0.157.0")
        
    Returns:
        Dictionary with comparison data:
        {
            "commits": <count>,
            "files_changed": <count>,
            "files": [{"filename": "...", "status": "added|modified|removed"}],
            "url": "github.com/..."
        }
    """
    url = f"https://api.github.com/repos/{repo}/compare/{base}...{head}"
    data = _make_github_request(url)
    
    return {
        "commits": len(data.get("commits", [])),
        "files_changed": len(data.get("files", [])),
        "files": [
            {
                "filename": f.get("filename", ""),
                "status": f.get("status", ""),
                "additions": f.get("additions", 0),
                "deletions": f.get("deletions", 0),
                "changes": f.get("changes", 0)
            }
            for f in data.get("files", [])
        ],
        "url": data.get("html_url", "")
    }


def parse_breaking_changes(release_body: str) -> List[str]:
    """Extract breaking changes from release notes.
    
    Looks for common patterns in release notes:
    - "BREAKING CHANGE:"
    - "Breaking Changes"
    - Lines starting with "⚠️" or "🚨"
    
    Args:
        release_body: The release notes text
        
    Returns:
        List of breaking change descriptions
    """
    breaking_changes = []
    
    # Split into lines
    lines = release_body.split('\n')
    
    # Pattern 1: Dedicated "Breaking Changes" section
    in_breaking_section = False
    for line in lines:
        line = line.strip()
        
        # Start of breaking changes section
        if re.search(r'(?i)^#+\s*breaking\s+changes?', line):
            in_breaking_section = True
            continue
        
        # End of section (next heading)
        if in_breaking_section and re.match(r'^#+\s+', line):
            in_breaking_section = False
            continue
        
        # Collect items in breaking section
        if in_breaking_section and line and not line.startswith('#'):
            breaking_changes.append(line)
        
        # Pattern 2: Inline "BREAKING CHANGE:" markers
        if re.search(r'(?i)BREAKING\s+CHANGE:', line):
            breaking_changes.append(line)
        
        # Pattern 3: Warning emojis
        if line.startswith(('⚠️', '🚨', '❗')) and 'breaking' in line.lower():
            breaking_changes.append(line)
    
    return breaking_changes


def get_sdk_repo_name(sdk_package: str) -> str:
    """Map SDK package name to GitHub repository name.
    
    Args:
        sdk_package: Package like "msgraph-sdk-go" or "msgraph-beta-sdk-go"
        
    Returns:
        Repository name in format "owner/repo"
    """
    if "msgraph-beta-sdk-go" in sdk_package:
        return "microsoftgraph/msgraph-beta-sdk-go"
    if "msgraph-sdk-go-core" in sdk_package:
        return "microsoftgraph/msgraph-sdk-go-core"
    if "msgraph-sdk-go" in sdk_package:
        return "microsoftgraph/msgraph-sdk-go"

    raise ValueError(f"Unknown SDK package: {sdk_package}")
