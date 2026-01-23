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
    and extract resource-centric SDK usage mapping:
    - Which SDK types each Terraform resource uses
    - Which fields are accessed by each resource
    - Which methods are called by each resource
    - Which enums are used by each resource
    
    Args:
        repo_path: Path to the repository root
        
    Returns:
        Dictionary containing resource-centric SDK usage:
        {
            "terraform_resources": {<resource_name>: <ResourceInfo>},
            "terraform_actions": {<action_name>: <ResourceInfo>},
            "terraform_list_actions": {<list_action_name>: <ResourceInfo>},
            "terraform_ephemerals": {<ephemeral_name>: <ResourceInfo>},
            "terraform_data_sources": {<data_source_name>: <ResourceInfo>},
            "sdk_to_resource_index": {<sdk_type>: [<resource_names>]},
            "statistics": {<key>: <value>}
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
        stats = usage_data.get('statistics', {})
        print(f"✅ Found {stats.get('total_resources', 0)} resources, "
              f"{stats.get('total_actions', 0)} actions, "
              f"{stats.get('total_data_sources', 0)} data sources")
        print(f"   - {stats.get('total_sdk_types_used', 0)} SDK types used")
        print(f"   - {stats.get('total_sdk_methods_used', 0)} SDK methods used")
        print(f"   - {stats.get('total_enums_tracked', 0)} enums tracked")
        
        return usage_data
        
    except json.JSONDecodeError as e:
        print(f"❌ Failed to parse Go AST parser output:")
        print(result.stdout)
        raise RuntimeError(f"Invalid JSON from Go AST parser: {e}") from e


def get_resources_using_sdk_type(usage_data: Dict[str, Any], sdk_type: str) -> list:
    """Get all Terraform entities that use a specific SDK type.
    
    Args:
        usage_data: Usage data from extract_sdk_usage
        sdk_type: SDK type name (e.g., "models.User")
        
    Returns:
        List of Terraform entity names using this SDK type
    """
    index = usage_data.get("sdk_to_resource_index", {})
    return index.get(sdk_type, [])


def get_all_sdk_types_used(usage_data: Dict[str, Any]) -> list:
    """Get all unique SDK types used across all Terraform entities.
    
    Args:
        usage_data: Usage data from extract_sdk_usage
        
    Returns:
        Sorted list of SDK type names
    """
    index = usage_data.get("sdk_to_resource_index", {})
    return sorted(index.keys())


def get_resource_dependencies(usage_data: Dict[str, Any], resource_name: str) -> Dict[str, Any]:
    """Get SDK dependencies for a specific Terraform resource.
    
    Args:
        usage_data: Usage data from extract_sdk_usage
        resource_name: Terraform resource name (e.g., "microsoft365_user")
        
    Returns:
        SDK dependencies for the resource, or None if not found
    """
    # Check all entity types
    for entity_type in ["terraform_resources", "terraform_actions", "terraform_list_actions", 
                        "terraform_ephemerals", "terraform_data_sources"]:
        entities = usage_data.get(entity_type, {})
        if resource_name in entities:
            return entities[resource_name].get("sdk_dependencies", {})
    
    return None
