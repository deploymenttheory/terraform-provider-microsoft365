#!/usr/bin/env python3
"""Go AST parsing utilities for SDK usage extraction.

Uses a Go tool to analyze Go source code and extract SDK usage patterns.
"""

import json
import subprocess
from pathlib import Path
from typing import Dict, Any


def extract_sdk_usage(repo_path: Path) -> Dict[str, Any]:
    """Extract SDK usage from the Terraform provider codebase.
    
    Uses the Go AST parser tool to analyze all Go files in internal/services
    and extract:
    - Which SDK packages are imported
    - Which types are used
    - Which methods are called
    - Which fields are accessed
    
    Args:
        repo_path: Path to the repository root
        
    Returns:
        Dictionary containing SDK usage information:
        {
            "packages": {<package>: <usage_count>},
            "imports": {<package>: [<file_paths>]},
            "types": {<type>: {<field>: <count>}},
            "methods": {<method>: <count>},
            "fields": {<type>: {<field>: <count>}}
        }
    """
    extractor_path = repo_path / "scripts" / "pipeline" / "sdk_diff" / "tools" / "extract_usage.go"
    
    print(f"📊 Analyzing SDK usage in {repo_path}...")
    
    # Run the Go AST parser
    result = subprocess.run(
        ["go", "run", str(extractor_path), str(repo_path)],
        capture_output=True,
        text=True,
        check=False
    )
    
    if result.returncode != 0:
        print("❌ Error running Go AST parser:")
        print(result.stderr)
        raise RuntimeError(f"Go AST parser failed: {result.stderr}")
    
    try:
        usage_data = json.loads(result.stdout)
        
        # Print summary
        print(f"✅ Found {len(usage_data['packages'])} SDK packages in use")
        print(f"   - {len(usage_data['types'])} types")
        print(f"   - {len(usage_data['methods'])} methods")
        
        return usage_data
        
    except json.JSONDecodeError as e:
        print(f"❌ Failed to parse Go AST parser output:")
        print(result.stdout)
        raise RuntimeError(f"Invalid JSON from Go AST parser: {e}") from e


def filter_msgraph_usage(usage_data: Dict[str, Any]) -> Dict[str, Any]:
    """Filter usage data to only include Microsoft Graph SDK items.
    
    Args:
        usage_data: Raw usage data from extract_sdk_usage
        
    Returns:
        Filtered usage data containing only msgraph-related items
    """
    filtered = {
        "packages": {},
        "imports": {},
        "types": {},
        "methods": {},
        "fields": {}
    }
    
    # Only keep msgraph and kiota packages
    for pkg, count in usage_data.get("packages", {}).items():
        if "microsoftgraph" in pkg or "kiota" in pkg:
            filtered["packages"][pkg] = count
            if pkg in usage_data.get("imports", {}):
                filtered["imports"][pkg] = usage_data["imports"][pkg]
    
    # Filter types, methods, fields
    for key in ["types", "methods", "fields"]:
        for name, value in usage_data.get(key, {}).items():
            if any(pkg in name for pkg in filtered["packages"]):
                filtered[key][name] = value
    
    return filtered


def get_most_used_packages(usage_data: Dict[str, Any], top_n: int = 10) -> list:
    """Get the most frequently used SDK packages.
    
    Args:
        usage_data: Usage data from extract_sdk_usage
        top_n: Number of top packages to return
        
    Returns:
        List of (package_name, usage_count) tuples, sorted by count
    """
    packages = usage_data.get("packages", {})
    return sorted(packages.items(), key=lambda x: x[1], reverse=True)[:top_n]
