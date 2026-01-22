#!/usr/bin/env python3
"""Fetch coverage results from Codecov API with intelligent retry logic."""

import json
import time
import urllib.request
import urllib.error
from typing import Optional, Dict, Any


def fetch_codecov_coverage(
    repo_slug: str,
    pr_number: str,
    codecov_token: str,
    max_wait_seconds: int = 180
) -> Optional[Dict[str, Any]]:
    """Fetch coverage results from Codecov API with exponential backoff.
    
    Args:
        repo_slug: Repository in format 'owner/repo'.
        pr_number: Pull request number.
        codecov_token: Codecov API token.
        max_wait_seconds: Maximum total wait time in seconds (default: 180 = 3 minutes).
    
    Returns:
        Dict with coverage statistics or None if fetch fails.
    """
    api_url = f"https://api.codecov.io/api/v2/github/{repo_slug}/pulls/{pr_number}"
    
    print(f"\n⏳ Fetching coverage from Codecov (max wait: {max_wait_seconds}s)...")
    print(f"   API: {api_url}")
    
    # Exponential backoff schedule: 2, 4, 8, 10, 10, 10... (capped at 10s)
    delays = []
    total_time = 0
    delay = 2
    
    while total_time < max_wait_seconds:
        delays.append(delay)
        total_time += delay
        delay = min(delay * 2, 10)  # Double but cap at 10 seconds
    
    print(f"   Retry schedule: {len(delays)} attempts with exponential backoff\n")
    
    start_time = time.time()
    
    for attempt, delay in enumerate(delays, 1):
        elapsed = int(time.time() - start_time)
        
        try:
            req = urllib.request.Request(api_url)
            req.add_header('Authorization', f'Bearer {codecov_token}')
            req.add_header('Accept', 'application/json')
            
            with urllib.request.urlopen(req, timeout=30) as response:
                data = json.loads(response.read().decode())
                
                # Check if we have patch coverage data
                if 'totals' in data and 'patch' in data['totals']:
                    patch_coverage = data['totals']['patch']
                    
                    if patch_coverage and 'coverage' in patch_coverage:
                        coverage_pct = patch_coverage['coverage']
                        
                        print(f"\n✅ Coverage retrieved after {elapsed}s ({attempt} attempts)")
                        print(f"   Patch Coverage: {coverage_pct}%")
                        
                        return {
                            "coverage_pct": round(float(coverage_pct), 2),
                            "total_lines": patch_coverage.get('lines', 0),
                            "covered_lines": patch_coverage.get('covered', 0)
                        }
                
                # Data exists but coverage not ready yet
                print(f"⏳ [{elapsed}s] Attempt {attempt}/{len(delays)}: "
                      f"PR found but coverage not processed yet, retrying in {delay}s...")
                time.sleep(delay)
                
        except urllib.error.HTTPError as e:
            if e.code == 404:
                # 404 could mean PR doesn't exist yet in Codecov (eventual consistency)
                print(f"⏳ [{elapsed}s] Attempt {attempt}/{len(delays)}: "
                      f"PR not found in Codecov yet (eventual consistency), retrying in {delay}s...")
                time.sleep(delay)
            elif e.code == 401:
                print("\n❌ Authentication Error: Invalid Codecov token")
                print("   Please verify CODECOV_TOKEN is correct")
                return None
            elif e.code == 403:
                print("\n❌ Authorization Error: Token lacks permissions")
                print("   Token may be valid but doesn't have access to this repository")
                return None
            else:
                print(f"\n❌ HTTP Error {e.code}: {e.reason}")
                if attempt < len(delays):
                    print(f"   Retrying in {delay}s...")
                    time.sleep(delay)
                else:
                    return None
        except urllib.error.URLError as e:
            print(f"⚠️  Network error: {e.reason}")
            if attempt < len(delays):
                print(f"   Retrying in {delay}s...")
                time.sleep(delay)
            else:
                return None
        except Exception as e:
            print(f"⚠️  Unexpected error: {type(e).__name__}: {e}")
            if attempt < len(delays):
                print(f"   Retrying in {delay}s...")
                time.sleep(delay)
            else:
                return None
    
    # Timeout reached
    elapsed = int(time.time() - start_time)
    print(f"\n❌ Timeout after {elapsed}s: Codecov did not process coverage within {max_wait_seconds}s")
    print("   This could indicate:")
    print("   1. Codecov is experiencing delays (check status.codecov.io)")
    print("   2. The PR upload failed (check earlier Codecov upload step)")
    print("   3. The PR number or repo slug is incorrect")
    print(f"      - Repo: {repo_slug}")
    print(f"      - PR: {pr_number}")
    
    return None


if __name__ == "__main__":
    import sys
    import argparse
    
    parser = argparse.ArgumentParser(description='Fetch coverage from Codecov')
    parser.add_argument('--repo-slug', required=True, help='Repository slug (owner/repo)')
    parser.add_argument('--pr-number', required=True, help='Pull request number')
    parser.add_argument('--codecov-token', required=True, help='Codecov API token')
    parser.add_argument('--max-wait', type=int, default=180, help='Max wait time in seconds')
    
    args = parser.parse_args()
    
    result = fetch_codecov_coverage(
        args.repo_slug,
        args.pr_number,
        args.codecov_token,
        args.max_wait
    )
    
    if result:
        print(f"\nCoverage: {result['coverage_pct']}%")
        sys.exit(0)
    else:
        sys.exit(1)
