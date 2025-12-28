#!/usr/bin/env python3
"""
Scan Terraform provider codebase to extract REAL Microsoft Graph API endpoints
from Go SDK method calls.

This version parses actual SDK method chains like:
  client.DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId(id).Get()
  
And converts them to Graph API paths:
  /deviceManagement/deviceHealthScripts/{id}
"""

import argparse
import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Set, Tuple


class GraphAPIEndpoint:
    """Represents a Graph API endpoint extracted from SDK calls."""
    
    def __init__(self, resource_name: str, file_path: str, service_type: str, service_domain: str, api_version: str):
        self.resource_name = resource_name
        self.file_path = file_path
        self.service_type = service_type  # resources, datasources, ephemeral, actions
        self.service_domain = service_domain  # device_management, identity_and_access, m365_admin, windows_365
        self.api_version = api_version  # graph_beta, graph_v1.0
        self.sdk_chains = set()  # Raw SDK method chains
        self.api_paths = set()    # Converted API paths
        self.operations = set()   # CRUD operations
        self.graph_resources = set()  # Resource type names
        
    def add_sdk_chain(self, chain: str, operation: str):
        """Add an SDK method chain and convert it to an API path."""
        self.sdk_chains.add(chain)
        self.operations.add(operation)
        
        # Convert SDK chain to API path
        api_path = self._sdk_chain_to_api_path(chain)
        if api_path:
            self.api_paths.add(api_path)
            
            # Extract resource types from path
            self._extract_resources_from_path(api_path)
    
    def _sdk_chain_to_api_path(self, chain: str) -> str:
        """
        Convert SDK method chain to Graph API path.
        
        Examples:
          DeviceManagement().DeviceHealthScripts() 
            → /deviceManagement/deviceHealthScripts
          
          DeviceManagement().DeviceHealthScripts().ByDeviceHealthScriptId(id)
            → /deviceManagement/deviceHealthScripts/{id}
        """
        # Remove common prefixes
        chain = chain.replace('client.', '').replace('r.client.', '')
        
        # Split by method calls
        methods = re.findall(r'([A-Z][a-zA-Z0-9]*)\(\)', chain)
        
        if not methods:
            return None
        
        path_parts = []
        for method in methods:
            # Handle By{Type}Id() pattern - these are ID placeholders
            if method.startswith('By') and method.endswith('Id'):
                path_parts.append('{id}')
            else:
                # Convert PascalCase to camelCase for first segment
                if len(path_parts) == 0:
                    # First segment: PascalCase → camelCase (DeviceManagement → deviceManagement)
                    camel = method[0].lower() + method[1:]
                    path_parts.append(camel)
                else:
                    # Subsequent segments: keep camelCase (DeviceHealthScripts → deviceHealthScripts)
                    camel = method[0].lower() + method[1:]
                    path_parts.append(camel)
        
        return '/' + '/'.join(path_parts)
    
    def _extract_resources_from_path(self, path: str):
        """Extract resource type names from API path."""
        # Remove leading slash and split
        parts = path.lstrip('/').split('/')
        
        for part in parts:
            # Skip ID placeholders
            if part == '{id}' or '{' in part:
                continue
            # Add non-trivial path segments as resources
            if len(part) > 2 and part not in ['api', 'v1.0', 'beta']:
                self.graph_resources.add(part)
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            'resource_name': self.resource_name,
            'file_path': self.file_path,
            'service_type': self.service_type,
            'service_domain': self.service_domain,
            'api_version': self.api_version,
            'api_paths': sorted(list(self.api_paths)),
            'graph_resources': sorted(list(self.graph_resources)),
            'operations': sorted(list(self.operations)),
            'sdk_chains': sorted(list(self.sdk_chains))[:10]  # Limit to 10 examples
        }


def find_go_files(base_path: str) -> List[str]:
    """Find all Go files in resource directories."""
    go_files = []
    
    resource_dirs = [
        'internal/resources',
        'internal/services',
    ]
    
    for resource_dir in resource_dirs:
        full_path = os.path.join(base_path, resource_dir)
        if not os.path.exists(full_path):
            continue
        
        for root, dirs, files in os.walk(full_path):
            for file in files:
                if file.endswith('.go') and not file.endswith('_test.go'):
                    go_files.append(os.path.join(root, file))
    
    return go_files


def extract_sdk_calls(file_path: str) -> List[Tuple[str, str]]:
    """
    Extract Graph SDK method chains and their operations from a Go file.
    
    Returns list of (sdk_chain, operation) tuples.
    """
    sdk_calls = []
    
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Pattern to match multi-line SDK method chains
        # Matches: r.client.\n\t\tDeviceManagement().\n\t\tDeviceHealthScripts().\n\t\tPost(...)
        
        # First, normalize whitespace in method chains
        # Find patterns like: r.client. or client. followed by chained methods
        pattern = r'(?:r\.client|client)\s*\.\s*([A-Z][a-zA-Z0-9]*\(\)[.\s]*)+([A-Z][a-z]+)\s*\('
        
        matches = re.finditer(pattern, content, re.MULTILINE)
        
        for match in matches:
            full_match = match.group(0)
            
            # Extract the method chain (without final operation)
            # Remove whitespace and newlines
            chain = re.sub(r'\s+', '', full_match)
            
            # Extract operation (last method call before final parenthesis)
            operation_match = re.search(r'\.([A-Z][a-z]+)\(', chain)
            if operation_match:
                operation = operation_match.group(1).lower()
                
                # Map operation to CRUD
                if operation in ['post', 'create']:
                    operation = 'create'
                elif operation in ['get', 'list']:
                    operation = 'read'
                elif operation in ['patch', 'put', 'update']:
                    operation = 'update'
                elif operation in ['delete']:
                    operation = 'delete'
                
                # Remove the final operation from chain for path extraction
                chain_without_op = re.sub(r'\.[A-Z][a-z]+\($', '', chain)
                
                sdk_calls.append((chain_without_op, operation))
    
    except Exception as e:
        print(f"Warning: Error reading {file_path}: {e}", file=sys.stderr)
    
    return sdk_calls


def scan_provider(base_path: str, verbose: bool = False) -> List[GraphAPIEndpoint]:
    """Scan the provider codebase and extract Graph API endpoints from SDK calls."""
    print(f"Scanning provider at: {base_path}")
    
    go_files = find_go_files(base_path)
    print(f"Found {len(go_files)} Go files")
    
    # Group by full resource path (action/service_domain/version/resource)
    resources_map = {}
    
    for file_path in go_files:
        # Extract resource metadata from path
        resource_name, service_type, service_domain, api_version = extract_resource_metadata_from_path(file_path, base_path)
        
        # Create unique key for grouping
        resource_key = f"{service_type}/{service_domain}/{api_version}/{resource_name}"
        
        if resource_key not in resources_map:
            resources_map[resource_key] = GraphAPIEndpoint(
                resource_name,
                os.path.relpath(os.path.dirname(file_path), base_path),
                service_type,
                service_domain,
                api_version
            )
        
        endpoint = resources_map[resource_key]
        
        # Extract SDK calls from file
        sdk_calls = extract_sdk_calls(file_path)
        
        if verbose and sdk_calls:
            print(f"\n  [{service_type}] {service_domain}/{api_version}/{resource_name}:")
            for chain, operation in sdk_calls[:3]:
                path = endpoint._sdk_chain_to_api_path(chain)
                print(f"    {operation}: {path}")
        
        for chain, operation in sdk_calls:
            endpoint.add_sdk_chain(chain, operation)
    
    # Filter to only resources that have API paths
    endpoints = [ep for ep in resources_map.values() if ep.api_paths]
    
    print(f"Extracted {len(endpoints)} provider resources with Graph API endpoints")
    
    return endpoints


def extract_resource_metadata_from_path(file_path: str, base_path: str) -> Tuple[str, str, str, str]:
    """
    Extract resource metadata from file path matching provider structure.
    
    Returns: (resource_name, service_type, service_domain, api_version)
    
    Example path:
      internal/services/resources/device_management/graph_beta/windows_remediation_script/crud.go
      → ('windows_remediation_script', 'resources', 'device_management', 'graph_beta')
    
    Provider structure:
      internal/services/{service_type}/{service_domain}/{api_version}/{resource_name}/
      
      service_type: resources, datasources, ephemeral, actions
      service_domain: device_management, identity_and_access, m365_admin, windows_365
      api_version: graph_beta, graph_v1.0
      resource_name: the specific resource
    """
    rel_path = os.path.relpath(file_path, base_path)
    parts = rel_path.split('/')
    
    service_type = 'unknown'
    service_domain = 'unknown'
    api_version = 'unknown'
    resource_name = 'unknown'
    
    try:
        # Find 'services' index
        if 'services' in parts:
            services_idx = parts.index('services')
            
            # Service type is next: resources, datasources, ephemeral, actions
            if services_idx + 1 < len(parts):
                service_type = parts[services_idx + 1]
            
            # Service domain is next: device_management, identity_and_access, etc.
            if services_idx + 2 < len(parts):
                service_domain = parts[services_idx + 2]
            
            # API version is next: graph_beta, graph_v1.0
            if services_idx + 3 < len(parts):
                api_version = parts[services_idx + 3]
            
            # Resource name is the directory containing the file
            if services_idx + 4 < len(parts):
                resource_name = parts[services_idx + 4]
    
    except (ValueError, IndexError):
        # Fallback for non-standard paths
        resource_name = os.path.basename(os.path.dirname(file_path))
    
    return resource_name, service_type, service_domain, api_version


def main():
    parser = argparse.ArgumentParser(
        description="Scan Terraform provider for Graph API endpoints from SDK calls"
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
        help='Base path of provider codebase (default: current directory)'
    )
    parser.add_argument(
        '--verbose',
        action='store_true',
        help='Enable verbose output with SDK call examples'
    )
    
    args = parser.parse_args()
    
    try:
        endpoints = scan_provider(args.base_path, verbose=args.verbose)
        
        # Build lookup structures
        endpoint_lookup = {
            'api_paths': {},
            'graph_resources': {}
        }
        
        for ep in endpoints:
            # Index by API paths
            for path in ep.api_paths:
                if path not in endpoint_lookup['api_paths']:
                    endpoint_lookup['api_paths'][path] = []
                endpoint_lookup['api_paths'][path].append(ep.resource_name)
            
            # Index by Graph resource types
            for resource in ep.graph_resources:
                if resource not in endpoint_lookup['graph_resources']:
                    endpoint_lookup['graph_resources'][resource] = {
                        'resources': [],
                        'operations': set()
                    }
                endpoint_lookup['graph_resources'][resource]['resources'].append(ep.resource_name)
                endpoint_lookup['graph_resources'][resource]['operations'].update(ep.operations)
        
        # Convert sets to lists for JSON
        for resource in endpoint_lookup['graph_resources'].values():
            resource['operations'] = sorted(list(resource['operations']))
        
        # Create output
        output_data = {
            'generated_at': json.dumps(sys.maxsize),  # Placeholder
            'extraction_method': 'sdk_call_parsing',
            'total_resources': len(endpoints),
            'total_api_paths': len(endpoint_lookup['api_paths']),
            'total_graph_resources': len(endpoint_lookup['graph_resources']),
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
        print("SUMMARY (SDK-Based Extraction)")
        print("="*60)
        print(f"Total resources: {len(endpoints)}")
        print(f"Total API paths: {len(endpoint_lookup['api_paths'])}")
        print(f"Total Graph resource types: {len(endpoint_lookup['graph_resources'])}")
        
        print(f"\nTop Graph resource types:")
        sorted_resources = sorted(
            endpoint_lookup['graph_resources'].items(),
            key=lambda x: len(x[1]['resources']),
            reverse=True
        )[:15]
        for resource, data in sorted_resources:
            ops = ', '.join(data['operations']) if data['operations'] else 'none'
            print(f"  - {resource}: {len(data['resources'])} resources [{ops}]")
        
        if args.verbose:
            print(f"\nSample API paths:")
            for path in sorted(list(endpoint_lookup['api_paths'].keys()))[:10]:
                resources = endpoint_lookup['api_paths'][path]
                print(f"  {path}")
                print(f"    Used by: {', '.join(resources[:3])}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()

