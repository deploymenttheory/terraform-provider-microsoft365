#!/usr/bin/env python3
"""Fetch coverage results from Codecov API v2 with intelligent retry logic.

This module provides functions to fetch PR coverage statistics from Codecov's
API v2, handling authentication, retries with exponential backoff, and error handling.

See: https://docs.codecov.com/reference/repos_pulls_retrieve
"""

import json
import time
import urllib.request
import urllib.error
from typing import Optional, Dict, Any, List, Tuple


# ============================================================================
# URL Building
# ============================================================================

def _build_api_url(service: str, owner: str, repo: str, pr_number: str) -> str:
    """Build Codecov API v2 URL for pull request details.
    
    Args:
        service: Git service provider (e.g., 'github', 'gitlab').
        owner: Repository owner username.
        repo: Repository name.
        pr_number: Pull request number.
    
    Returns:
        Complete API endpoint URL.
    
    See: https://docs.codecov.com/reference/repos_pulls_retrieve
    """
    return f"https://api.codecov.io/api/v2/{service}/{owner}/repos/{repo}/pulls/{pr_number}"


def _parse_repo_slug(repo_slug: str) -> Tuple[str, str]:
    """Parse repository slug into owner and repo name.
    
    Args:
        repo_slug: Repository slug in format 'owner/repo'.
    
    Returns:
        Tuple of (owner, repo).
    
    Raises:
        ValueError: If slug format is invalid.
    """
    parts = repo_slug.split('/')
    if len(parts) != 2:
        raise ValueError(f"Invalid repo_slug format: {repo_slug}. Expected 'owner/repo'")
    return parts[0], parts[1]


# ============================================================================
# HTTP Request Building & Execution
# ============================================================================

def _build_request(api_url: str, token: str) -> urllib.request.Request:
    """Build authenticated HTTP request for Codecov API.
    
    Args:
        api_url: Complete API endpoint URL.
        token: Codecov authentication token.
    
    Returns:
        Configured Request object with auth headers.
    """
    req = urllib.request.Request(api_url)
    req.add_header('Authorization', f'Bearer {token}')
    req.add_header('Accept', 'application/json')
    return req


def _execute_request(request: urllib.request.Request, timeout: int = 30) -> Dict[str, Any]:
    """Execute HTTP request and parse JSON response.
    
    Args:
        request: Configured Request object.
        timeout: Request timeout in seconds.
    
    Returns:
        Parsed JSON response as dictionary.
    
    Raises:
        urllib.error.HTTPError: For HTTP errors (4xx, 5xx).
        urllib.error.URLError: For network/connectivity errors.
        json.JSONDecodeError: For invalid JSON responses.
    """
    with urllib.request.urlopen(request, timeout=timeout) as response:
        return json.loads(response.read().decode())


# ============================================================================
# Response Parsing
# ============================================================================

def _parse_coverage_response(data: Dict[str, Any]) -> Optional[Dict[str, Any]]:
    """Extract patch coverage statistics from Codecov API response.
    
    Args:
        data: Codecov API response JSON.
    
    Returns:
        Dict with coverage stats if available, None otherwise:
        {
            "coverage_pct": float,
            "total_lines": int,
            "covered_lines": int
        }
    """
    if 'totals' not in data or 'patch' not in data['totals']:
        return None
    
    patch = data['totals']['patch']
    if not patch or 'coverage' not in patch:
        return None
    
    return {
        "coverage_pct": round(float(patch['coverage']), 2),
        "total_lines": patch.get('lines', 0),
        "covered_lines": patch.get('covered', 0)
    }


# ============================================================================
# Error Handling
# ============================================================================

def _handle_error(
    error: Exception,
    elapsed: int,
    attempt: int,
    total_attempts: int
) -> bool:
    """Handle errors and determine if retry should continue.
    
    Args:
        error: The exception that occurred.
        elapsed: Elapsed time in seconds since start.
        attempt: Current attempt number.
        total_attempts: Total number of attempts available.
    
    Returns:
        True if should retry, False if should abort.
    """
    if isinstance(error, urllib.error.HTTPError):
        return _handle_http_error(error, elapsed, attempt, total_attempts)
    elif isinstance(error, urllib.error.URLError):
        return _handle_network_error(error, attempt, total_attempts)
    else:
        print(f"⚠️  Unexpected error: {type(error).__name__}: {error}")
        return attempt < total_attempts


def _handle_http_error(
    error: urllib.error.HTTPError,
    elapsed: int,
    attempt: int,
    total_attempts: int
) -> bool:
    """Handle HTTP-specific errors.
    
    Args:
        error: The HTTP error.
        elapsed: Elapsed time in seconds.
        attempt: Current attempt number.
        total_attempts: Total number of attempts.
    
    Returns:
        True if should retry, False if should abort.
    """
    if error.code == 404:
        print(f"⏳ [{elapsed}s] Attempt {attempt}/{total_attempts}: "
              f"PR not found in Codecov (eventual consistency delay)")
        return True
    elif error.code == 401:
        print("\n❌ Authentication Error: Invalid Codecov token")
        print("   Verify CODECOV_TOKEN is correct")
        return False
    elif error.code == 403:
        print("\n❌ Authorization Error: Token lacks repository access")
        print("   Token may be valid but lacks permissions for this repo")
        return False
    else:
        print(f"\n❌ HTTP {error.code}: {error.reason}")
        return attempt < total_attempts


def _handle_network_error(
    error: urllib.error.URLError,
    attempt: int,
    total_attempts: int
) -> bool:
    """Handle network connectivity errors.
    
    Args:
        error: The network error.
        attempt: Current attempt number.
        total_attempts: Total number of attempts.
    
    Returns:
        True if should retry, False if should abort.
    """
    print(f"⚠️  Network error: {error.reason}")
    return attempt < total_attempts


# ============================================================================
# Retry Logic
# ============================================================================

def _build_backoff_schedule(max_wait_seconds: int) -> List[int]:
    """Build exponential backoff schedule with 10-second cap.
    
    Schedule: 2s, 4s, 8s, 10s, 10s, 10s...
    
    Args:
        max_wait_seconds: Maximum total wait time across all retries.
    
    Returns:
        List of delay intervals in seconds.
    """
    delays = []
    total = 0
    delay = 2
    
    while total < max_wait_seconds:
        delays.append(delay)
        total += delay
        delay = min(delay * 2, 10)
    
    return delays


def _retry_with_backoff(
    api_url: str,
    token: str,
    delays: List[int],
    start_time: float
) -> Optional[Dict[str, Any]]:
    """Retry API requests with exponential backoff.
    
    Args:
        api_url: Codecov API endpoint URL.
        token: Authentication token.
        delays: List of delay intervals for retry schedule.
        start_time: Timestamp when retries started.
    
    Returns:
        Coverage statistics dict or None if all retries exhausted.
    """
    for attempt, delay in enumerate(delays, 1):
        elapsed = int(time.time() - start_time)
        
        try:
            # Build and execute request
            request = _build_request(api_url, token)
            data = _execute_request(request)
            
            # Parse response
            coverage = _parse_coverage_response(data)
            
            if coverage:
                print(f"\n✅ Coverage retrieved after {elapsed}s ({attempt} attempts)")
                print(f"   Patch Coverage: {coverage['coverage_pct']}%")
                return coverage
            
            # Data exists but coverage not ready
            print(f"⏳ [{elapsed}s] Attempt {attempt}/{len(delays)}: "
                  f"PR found but coverage not processed yet")
            
        except urllib.error.HTTPError as e:
            if not _handle_error(e, elapsed, attempt, len(delays)):
                return None
        
        except urllib.error.URLError as e:
            if not _handle_error(e, elapsed, attempt, len(delays)):
                return None
        
        except (json.JSONDecodeError, OSError, ValueError) as e:
            if not _handle_error(e, elapsed, attempt, len(delays)):
                return None
        
        # Wait before next retry (unless this was the last attempt)
        if attempt < len(delays):
            print(f"   Retrying in {delay}s...")
            time.sleep(delay)
    
    return None


# ============================================================================
# Public API
# ============================================================================

def fetch_codecov_coverage(
    repo_slug: str,
    pr_number: str,
    codecov_token: str,
    max_wait_seconds: int = 180,
    service: str = "github"
) -> Optional[Dict[str, Any]]:
    """Fetch PR patch coverage from Codecov API v2 with exponential backoff.
    
    This function implements robust fetching with:
    - Exponential backoff retry (2s, 4s, 8s, 10s...)
    - Eventual consistency handling for 404s
    - Detailed error messages for auth/network issues
    - Timeout protection
    
    Args:
        repo_slug: Repository in format 'owner/repo'.
        pr_number: Pull request number.
        codecov_token: Codecov API authentication token.
        max_wait_seconds: Maximum total wait time (default: 180s = 3 minutes).
        service: Git service provider (default: 'github').
    
    Returns:
        Dict with coverage statistics or None if fetch fails:
        {
            "coverage_pct": 85.5,      # Patch coverage percentage
            "total_lines": 100,        # Total lines in patch
            "covered_lines": 85        # Covered lines in patch
        }
    
    Example:
        >>> coverage = fetch_codecov_coverage(
        ...     repo_slug="deploymenttheory/terraform-provider-microsoft365",
        ...     pr_number="123",
        ...     codecov_token="abc123..."
        ... )
        >>> if coverage:
        ...     print(f"Coverage: {coverage['coverage_pct']}%")
    
    See: https://docs.codecov.com/reference/repos_pulls_retrieve
    """

    owner, repo = _parse_repo_slug(repo_slug)
    api_url = _build_api_url(service, owner, repo, pr_number)
    delays = _build_backoff_schedule(max_wait_seconds)
    
    print("\n⏳ Fetching coverage from Codecov API v2")
    print(f"   Endpoint: {api_url}")
    print(f"   Max wait: {max_wait_seconds}s ({len(delays)} attempts with exponential backoff)\n")
    
    start_time = time.time()
    coverage_data = _retry_with_backoff(api_url, codecov_token, delays, start_time)
    
    if coverage_data:
        return coverage_data
    
    elapsed = int(time.time() - start_time)
    print(f"\n❌ Timeout after {elapsed}s: Coverage not available within {max_wait_seconds}s")
    print("   Possible causes:")
    print("   1. Codecov processing delays (check https://status.codecov.io)")
    print("   2. Upload failed (verify earlier 'Upload coverage to Codecov' step)")
    print("   3. Incorrect parameters:")
    print(f"      - Service: {service}")
    print(f"      - Owner: {owner}")
    print(f"      - Repo: {repo}")
    print(f"      - PR: {pr_number}")
    
    return None


# ============================================================================
# CLI Entry Point
# ============================================================================

if __name__ == "__main__":
    import sys
    import argparse

    parser = argparse.ArgumentParser(
        description='Fetch PR coverage from Codecov API v2',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python get_codecov_coverage.py \\
    --repo-slug deploymenttheory/terraform-provider-microsoft365 \\
    --pr-number 123 \\
    --codecov-token $CODECOV_TOKEN

  python get_codecov_coverage.py \\
    --repo-slug owner/repo \\
    --pr-number 456 \\
    --codecov-token abc123 \\
    --max-wait 300 \\
    --service github
        """
    )
    parser.add_argument('--repo-slug', required=True, 
                        help='Repository slug (owner/repo)')
    parser.add_argument('--pr-number', required=True, 
                        help='Pull request number')
    parser.add_argument('--codecov-token', required=True, 
                        help='Codecov API authentication token')
    parser.add_argument('--max-wait', type=int, default=180, 
                        help='Maximum wait time in seconds (default: 180)')
    parser.add_argument('--service', default='github', 
                        help='Git service provider (default: github)')

    args = parser.parse_args()

    result = fetch_codecov_coverage(
        args.repo_slug,
        args.pr_number,
        args.codecov_token,
        args.max_wait,
        args.service
    )

    if result:
        print(f"\n✓ Coverage: {result['coverage_pct']}%")
        print(f"  Lines: {result['covered_lines']}/{result['total_lines']}")
        sys.exit(0)
    else:
        print("\n✗ Failed to fetch coverage")
        sys.exit(1)
