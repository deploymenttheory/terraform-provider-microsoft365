#!/usr/bin/env python3
"""
Scan the Terraform provider codebase to identify implemented Microsoft Graph API endpoints.

This script walks through the provider codebase and extracts information about
what Graph API resources and endpoints are currently implemented.
"""

import argparse
import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Set


class ProviderEndpoint:
    """Represents a Graph API endpoint implemented in the provider."""
    
    def __init__(self, resource_name: str, file_path: str):
        self.resource_name = resource_name
        self.file_path = file_path
        self.graph_resources = set()  # Graph API resource types
        self.graph_methods = set()    # HTTP methods and operations
        self.endpoints = set()         # Actual API endpoint paths
        self.operations = set()        # CRUD operations supported
        
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            'resource_name': self.resource_name,
            'file_path': self.file_path,
            'graph_resources': sorted(list(self.graph_resources)),
            'graph_methods': sorted(list(self.graph_methods)),
            'endpoints': sorted(list(self.endpoints)),
            'operations': sorted(list(self.operations))
        }


def find_resource_files(base_path: str) -> List[str]:
    """Find all resource implementation files in the provider."""
    resource_paths = []
    
    # Main resource directories
    resource_dirs = [
        'internal/resources',
        'internal/services',
    ]
    
    for resource_dir in resource_dirs:
        full_path = os.path.join(base_path, resource_dir)
        if not os.path.exists(full_path):
            continue
        
        # Walk through directories
        for root, dirs, files in os.walk(full_path):
            for file in files:
                # Look for Go files that likely contain resource implementations
                if file.endswith('.go') and not file.endswith('_test.go'):
                    if any(keyword in file for keyword in ['resource.go', 'crud.go', 'create.go', 'read.go', 'update.go', 'delete.go']):
                        resource_paths.append(os.path.join(root, file))
    
    return resource_paths


def extract_graph_api_calls(file_path: str) -> Set[str]:
    """Extract Graph API endpoint calls from a Go file."""
    endpoints = set()
    
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Pattern 1: Direct URL construction
        # Example: fmt.Sprintf("/deviceManagement/deviceHealthScripts/%s", id)
        url_patterns = [
            r'["\']/?([a-zA-Z0-9/{}%\-_]+)["\']',
            r'fmt\.Sprintf\(["\']([^"\']+)["\']',
            r'uri\s*:?=\s*["\']([^"\']+)["\']',
        ]
        
        for pattern in url_patterns:
            matches = re.findall(pattern, content)
            for match in matches:
                # Filter for Graph API paths
                if any(keyword in match.lower() for keyword in [
                    'devicemanagement', 'device', 'intune', 'policy',
                    'configuration', 'application', 'users', 'groups',
                    'conditional', 'identity', 'security'
                ]):
                    # Clean up the endpoint
                    endpoint = match.replace('%s', '{id}').replace('%d', '{id}')
                    endpoint = re.sub(r'\$\{[^}]+\}', '{id}', endpoint)
                    endpoints.add(endpoint)
        
        # Pattern 2: Model/SDK calls
        # Example: client.Get(ctx, "deviceManagement/deviceHealthScripts/...")
        sdk_pattern = r'client\.(Get|Post|Patch|Put|Delete|Create|Update|Read)\([^,]+,\s*["\']([^"\']+)["\']'
        sdk_matches = re.findall(sdk_pattern, content, re.IGNORECASE)
        for method, endpoint in sdk_matches:
            if any(keyword in endpoint.lower() for keyword in [
                'devicemanagement', 'device', 'intune', 'policy',
                'configuration', 'application', 'users', 'groups'
            ]):
                endpoints.add(endpoint)
        
        # Pattern 3: Resource type constants
        # Example: const resourceType = "microsoft.graph.deviceManagementScript"
        resource_pattern = r'microsoft\.graph\.([a-zA-Z0-9_]+)'
        resource_matches = re.findall(resource_pattern, content)
        for resource in resource_matches:
            endpoints.add(f"resources/{resource}")
    
    except Exception as e:
        print(f"Warning: Error reading {file_path}: {e}", file=sys.stderr)
    
    return endpoints


def extract_operations_from_file(file_path: str) -> Set[str]:
    """Determine which CRUD operations are implemented in a file."""
    operations = set()
    
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Look for function signatures
        if re.search(r'func.*Create\(', content):
            operations.add('create')
        if re.search(r'func.*Read\(', content):
            operations.add('read')
        if re.search(r'func.*Update\(', content):
            operations.add('update')
        if re.search(r'func.*Delete\(', content):
            operations.add('delete')
        
        # Also check for HTTP method calls
        if re.search(r'\.Post\(', content):
            operations.add('create')
        if re.search(r'\.Get\(', content):
            operations.add('read')
        if re.search(r'\.(Patch|Put)\(', content):
            operations.add('update')
        if re.search(r'\.Delete\(', content):
            operations.add('delete')
    
    except Exception as e:
        print(f"Warning: Error reading {file_path}: {e}", file=sys.stderr)
    
    return operations


def extract_resource_name_from_path(file_path: str, base_path: str) -> str:
    """Extract a human-readable resource name from the file path."""
    rel_path = os.path.relpath(file_path, base_path)
    
    # Remove common prefixes and suffixes
    name = rel_path.replace('internal/services/resources/', '')
    name = name.replace('internal/resources/', '')
    name = name.replace('/resource.go', '')
    name = name.replace('/crud.go', '')
    name = name.replace('.go', '')
    
    # Convert path separators to underscores
    name = name.replace('/', '_')
    
    return name


def scan_provider(base_path: str) -> List[ProviderEndpoint]:
    """Scan the provider codebase and extract endpoint information."""
    print(f"Scanning provider at: {base_path}")
    
    # Find all resource files
    resource_files = find_resource_files(base_path)
    print(f"Found {len(resource_files)} resource files")
    
    endpoints = []
    
    for file_path in resource_files:
        # Extract resource name
        resource_name = extract_resource_name_from_path(file_path, base_path)
        
        # Create endpoint object
        endpoint = ProviderEndpoint(resource_name, os.path.relpath(file_path, base_path))
        
        # Extract Graph API calls
        endpoint.endpoints = extract_graph_api_calls(file_path)
        
        # Extract operations
        endpoint.operations = extract_operations_from_file(file_path)
        
        # Extract Graph resource types from endpoints
        for ep in endpoint.endpoints:
            # Try to extract resource type from endpoint path
            parts = ep.split('/')
            if len(parts) > 0:
                # First part is usually the resource collection
                resource_type = parts[0]
                if resource_type and not resource_type.startswith('{'):
                    endpoint.graph_resources.add(resource_type)
                
                # Check for specific resource names in subsequent parts
                for part in parts:
                    if part and not part.startswith('{') and len(part) > 3:
                        # Looks like a resource name, not an ID
                        if not part.isdigit():
                            endpoint.graph_resources.add(part)
        
        # Only add if we found something
        if endpoint.endpoints or endpoint.operations:
            endpoints.append(endpoint)
    
    print(f"Extracted {len(endpoints)} implemented resources")
    
    return endpoints


def main():
    parser = argparse.ArgumentParser(
        description="Scan Terraform provider for implemented Graph API endpoints"
    )
    parser.add_argument(
        '--output',
        type=str,
        default='provider-endpoints.json',
        help='Output JSON file path'
    )
    parser.add_argument(
        '--base-path',
        type=str,
        default='.',
        help='Base path of the provider codebase (default: current directory)'
    )
    
    args = parser.parse_args()
    
    try:
        # Scan provider
        endpoints = scan_provider(args.base_path)
        
        # Build a lookup structure for efficient comparison
        endpoint_lookup = {
            'resources': {},
            'endpoints': {},
            'operations': {}
        }
        
        for ep in endpoints:
            # Index by resource name
            endpoint_lookup['resources'][ep.resource_name] = ep.to_dict()
            
            # Index by endpoint paths
            for endpoint_path in ep.endpoints:
                if endpoint_path not in endpoint_lookup['endpoints']:
                    endpoint_lookup['endpoints'][endpoint_path] = []
                endpoint_lookup['endpoints'][endpoint_path].append(ep.resource_name)
            
            # Index by Graph resource types
            for graph_resource in ep.graph_resources:
                if graph_resource not in endpoint_lookup['operations']:
                    endpoint_lookup['operations'][graph_resource] = {
                        'resources': [],
                        'operations': set()
                    }
                endpoint_lookup['operations'][graph_resource]['resources'].append(ep.resource_name)
                endpoint_lookup['operations'][graph_resource]['operations'].update(ep.operations)
        
        # Convert sets to lists for JSON serialization
        for resource in endpoint_lookup['operations'].values():
            resource['operations'] = sorted(list(resource['operations']))
        
        # Create output structure
        output_data = {
            'generated_at': json.dumps(sys.maxsize),  # Will be replaced with timestamp
            'total_resources': len(endpoints),
            'total_endpoints': len(endpoint_lookup['endpoints']),
            'total_graph_resources': len(endpoint_lookup['operations']),
            'resources': [ep.to_dict() for ep in endpoints],
            'lookup': endpoint_lookup
        }
        
        # Fix timestamp
        from datetime import datetime
        output_data['generated_at'] = datetime.now().isoformat()
        
        # Write output
        with open(args.output, 'w') as f:
            json.dump(output_data, f, indent=2)
        
        print(f"\n✓ Successfully wrote provider endpoint data to {args.output}")
        
        # Print summary
        print("\n" + "="*60)
        print("SUMMARY")
        print("="*60)
        print(f"Total resources: {len(endpoints)}")
        print(f"Total unique endpoints: {len(endpoint_lookup['endpoints'])}")
        print(f"Total Graph resource types: {len(endpoint_lookup['operations'])}")
        
        print(f"\nTop Graph resource types:")
        sorted_resources = sorted(
            endpoint_lookup['operations'].items(),
            key=lambda x: len(x[1]['resources']),
            reverse=True
        )[:10]
        for resource, data in sorted_resources:
            print(f"  - {resource}: {len(data['resources'])} resources, {data['operations']}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()

