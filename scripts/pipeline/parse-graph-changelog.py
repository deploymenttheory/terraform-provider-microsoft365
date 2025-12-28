#!/usr/bin/env python3
"""
Parse Microsoft Graph API Changelog RSS feed using URL-based matching (v2).

This version extracts API resources and methods from documentation URLs in the
description field, providing precise matching instead of keyword-based filtering.
"""

import argparse
import json
import re
import sys
import traceback
from datetime import datetime, timedelta
from typing import Dict, List, Set
from urllib.parse import urlparse, parse_qs

try:
    import feedparser
    from bs4 import BeautifulSoup
except ImportError:
    print("Error: Required packages not installed. Run: pip install feedparser beautifulsoup4 lxml")
    sys.exit(1)


RSS_FEED_URL = "https://developer.microsoft.com/en-us/graph/changelog/rss"


class GraphAPIChange:
    """Represents a single API change from the changelog."""
    
    def __init__(self, guid: str, title: str, description: str, pub_date: str, 
                 categories: List[str], api_version: str):
        self.guid = guid
        self.title = title
        self.description = description
        self.pub_date = pub_date
        self.categories = categories
        self.api_version = api_version
        
        # Extracted from description URLs
        self.resources = []  # Resource type names
        self.methods = []    # Method names
        self.doc_urls = []   # Documentation URLs
        self.endpoints = []  # API endpoint patterns
        
        # Metadata
        self.change_type = None  # 'added', 'removed', 'deprecated', 'updated'
        self.change_actions = []  # List of specific actions (added resource, added method, etc.)
        
    def parse_description(self, debug: bool = False):
        """Parse the HTML description to extract API details from URLs and change actions."""
        soup = BeautifulSoup(self.description, 'html.parser')
        
        # Extract text content for pattern matching
        desc_text = soup.get_text()
        desc_lower = desc_text.lower()
        
        # Parse each div/line to extract specific change actions
        change_patterns = [
            # New resource
            (r'added the <(?:a|b).*?>([\w]+)<\/(?:a|b)> resource', 'added_resource'),
            # New property/field
            (r'added the <(?:a|b).*?>([\w]+)<\/(?:a|b)> (?:property|member|enumeration) to the <(?:a|b).*?>([\w]+)<\/(?:a|b)> (?:resource|enumeration)', 'added_property'),
            # New method
            (r'added the <a.*?>([\w-]+)<\/a> method to the <(?:a|b).*?>([\w]+)<\/(?:a|b)> resource', 'added_method'),
            # New relationship
            (r'added the <(?:a|b).*?>([\w]+)<\/(?:a|b)> relationship to the <(?:a|b).*?>([\w]+)<\/(?:a|b)> resource', 'added_relationship'),
            # Removed property
            (r'removed the <(?:a|b).*?>([\w]+)<\/(?:a|b)> (?:property|member) from the <(?:a|b).*?>([\w]+)<\/(?:a|b)> resource', 'removed_property'),
            # Deprecated
            (r'deprecated the <(?:a|b).*?>([\w]+)<\/(?:a|b)>', 'deprecated'),
        ]
        
        for pattern, action_type in change_patterns:
            for match in re.finditer(pattern, self.description, re.IGNORECASE):
                groups = match.groups()
                if action_type == 'added_resource':
                    self.change_actions.append({
                        'type': 'added_resource',
                        'resource': groups[0],
                        'impact': 'new_provider_resource'
                    })
                elif action_type == 'added_property':
                    self.change_actions.append({
                        'type': 'added_property',
                        'property': groups[0],
                        'resource': groups[1],
                        'impact': 'schema_update'
                    })
                elif action_type == 'added_method':
                    self.change_actions.append({
                        'type': 'added_method',
                        'method': groups[0],
                        'resource': groups[1],
                        'impact': 'new_operation'
                    })
                elif action_type == 'added_relationship':
                    self.change_actions.append({
                        'type': 'added_relationship',
                        'relationship': groups[0],
                        'resource': groups[1],
                        'impact': 'schema_update'
                    })
                elif action_type == 'removed_property':
                    self.change_actions.append({
                        'type': 'removed_property',
                        'property': groups[0],
                        'resource': groups[1],
                        'impact': 'breaking_change'
                    })
                elif action_type == 'deprecated':
                    self.change_actions.append({
                        'type': 'deprecated',
                        'item': groups[0],
                        'impact': 'deprecation_warning'
                    })
        
        # Determine overall change type from parsed actions
        if self.change_actions:
            if any(a['type'].startswith('added') for a in self.change_actions):
                self.change_type = 'added'
            elif any(a['type'].startswith('removed') for a in self.change_actions):
                self.change_type = 'removed'
            elif any(a['type'] == 'deprecated' for a in self.change_actions):
                self.change_type = 'deprecated'
        else:
            # Fallback to simple text matching
            if 'added' in desc_lower:
                self.change_type = 'added'
            elif 'removed' in desc_lower:
                self.change_type = 'removed'
            elif 'deprecated' in desc_lower:
                self.change_type = 'deprecated'
            elif 'updated' in desc_lower or 'changed' in desc_lower:
                self.change_type = 'updated'
        
        # Extract all links from description
        for link in soup.find_all('a'):
            href = link.get('href')
            if not href or 'learn.microsoft.com/en-us/graph/api' not in href:
                continue
            
            self.doc_urls.append(href)
            
            # Parse URL to extract resource/method information
            parsed_url = urlparse(href)
            path = parsed_url.path
            
            # Extract API version from query params
            query_params = parse_qs(parsed_url.query)
            url_api_version = None
            if 'view' in query_params:
                view = query_params['view'][0]
                if 'graph-rest-beta' in view:
                    url_api_version = 'beta'
                elif 'graph-rest-1.0' in view:
                    url_api_version = 'v1.0'
            
            # Extract resource/method from path
            # Patterns:
            # /en-us/graph/api/resources/RESOURCE_NAME
            # /en-us/graph/api/RESOURCE-METHOD
            # /en-us/graph/api/METHOD
            
            path_parts = path.split('/')
            
            if 'resources' in path_parts:
                # This is a resource definition
                idx = path_parts.index('resources')
                if idx + 1 < len(path_parts):
                    resource_name = path_parts[idx + 1]
                    # Remove query params if they leaked into the path
                    resource_name = resource_name.split('?')[0]
                    self.resources.append(resource_name)
                    
                    # Generate endpoint pattern
                    endpoint = f"resources/{resource_name}"
                    self.endpoints.append(endpoint)
                    
                    if debug:
                        print(f"      Found resource: {resource_name} ({url_api_version})")
            
            elif '/api/' in path:
                # This is a method/operation
                # Extract the part after /api/
                api_idx = path.rfind('/api/')
                if api_idx != -1:
                    method_part = path[api_idx + 5:].split('?')[0]
                    
                    # Parse method name and resource
                    # Format: resourceName-methodName or just methodName
                    if '-' in method_part:
                        resource_name, method_name = method_part.split('-', 1)
                        self.resources.append(resource_name)
                        self.methods.append(method_name)
                        endpoint = f"{resource_name}/{method_name}"
                        self.endpoints.append(endpoint)
                        
                        if debug:
                            print(f"      Found method: {resource_name}.{method_name} ({url_api_version})")
                    else:
                        # Just a method name
                        self.methods.append(method_part)
                        self.endpoints.append(f"api/{method_part}")
                        
                        if debug:
                            print(f"      Found operation: {method_part} ({url_api_version})")
        
        # Remove duplicates
        self.resources = list(set(self.resources))
        self.methods = list(set(self.methods))
        self.endpoints = list(set(self.endpoints))
        self.doc_urls = list(set(self.doc_urls))
    
    def is_relevant_to_provider(self, provider_resources: Set[str], debug: bool = False) -> bool:
        """
        Determine if this change is relevant by checking if any extracted resources
        match what's implemented in the provider.
        """
        # Always include if we have resources or methods (we'll compare later)
        if not self.resources and not self.methods:
            if debug:
                print("    No resources or methods extracted from URLs")
            return False
        
        # Check if any of our resources match provider resources
        matches = []
        for resource in self.resources:
            resource_lower = resource.lower()
            for provider_resource in provider_resources:
                provider_lower = provider_resource.lower()
                # Check for exact match or substring match
                if resource_lower == provider_lower or resource_lower in provider_lower or provider_lower in resource_lower:
                    matches.append(f"{resource} ↔ {provider_resource}")
        
        if matches and debug:
            print(f"    Matches provider resources: {matches}")
        
        # For now, include ALL changes that have resources/methods
        # The comparison script will determine if it's a gap
        return True
    
    def supports_crud_or_minimal(self) -> bool:
        """Check if this change involves CRUD operations."""
        if not self.methods:
            return False
        
        # Look for CRUD-related method names
        crud_keywords = ['create', 'get', 'list', 'update', 'patch', 'delete', 'post', 'put']
        return any(keyword in method.lower() for method in self.methods for keyword in crud_keywords)
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            'guid': self.guid,
            'title': self.title,
            'description': self.description[:500],
            'pub_date': self.pub_date,
            'categories': self.categories,
            'api_version': self.api_version,
            'resources': self.resources,
            'methods': self.methods,
            'endpoints': self.endpoints,
            'doc_urls': self.doc_urls,
            'change_type': self.change_type,
            'change_actions': self.change_actions,
            'supports_crud_or_minimal': self.supports_crud_or_minimal()
        }


def parse_rss_feed(url: str, lookback_days: int = 30, debug: bool = False, 
                   verbose: bool = False, provider_resources: Set[str] = None) -> List[GraphAPIChange]:
    """Parse the RSS feed and return a list of API changes."""
    print(f"Fetching RSS feed from {url}...")
    
    feed = feedparser.parse(url)
    
    if feed.bozo:
        print(f"Warning: Feed parsing encountered an error: {feed.bozo_exception}")
    
    print(f"Found {len(feed.entries)} entries in feed")
    
    if debug:
        print(f"\nDEBUG: Cutoff date: {datetime.now() - timedelta(days=lookback_days)}")
        print("DEBUG: Using URL-based extraction from documentation links")
        if provider_resources:
            print(f"DEBUG: Comparing against {len(provider_resources)} provider resources\n")
    
    changes = []
    cutoff_date = datetime.now() - timedelta(days=lookback_days)
    
    for entry in feed.entries:
        # Parse publication date
        pub_date = None
        if hasattr(entry, 'published_parsed'):
            pub_date = datetime(*entry.published_parsed[:6])
        elif hasattr(entry, 'published'):
            try:
                pub_date = datetime.strptime(entry.published, '%a, %d %b %Y %H:%M:%S %Z')
            except ValueError:
                pass
        
        # Skip entries older than lookback period
        if pub_date and pub_date < cutoff_date:
            continue
        
        # Extract categories and API version
        categories = []
        api_version = 'unknown'
        
        if hasattr(entry, 'tags'):
            for tag in entry.tags:
                if hasattr(tag, 'term'):
                    term = tag.term
                    if term in ['v1.0', 'beta']:
                        api_version = term
                    else:
                        categories.append(term)
        
        # Create change object
        change = GraphAPIChange(
            guid=entry.get('id', ''),
            title=entry.get('title', ''),
            description=entry.get('description', ''),
            pub_date=pub_date.isoformat() if pub_date else '',
            categories=categories,
            api_version=api_version
        )
        
        # Parse description to extract details from URLs
        if verbose and len(changes) < 3:
            print(f"\n  Parsing entry {len(changes) + 1}:")
            print(f"    Title: {change.title}")
            change.parse_description(debug=True)
        else:
            change.parse_description(debug=False)
        
        changes.append(change)
    
    print(f"Parsed {len(changes)} changes within the last {lookback_days} days")
    
    # Filter to relevant changes
    if provider_resources:
        relevant_changes = [c for c in changes if c.is_relevant_to_provider(provider_resources, debug=debug)]
        print(f"Found {len(relevant_changes)} changes with extractable resources/methods")
    else:
        # Without provider resources, include all that have resources/methods
        relevant_changes = [c for c in changes if c.resources or c.methods]
        print(f"Found {len(relevant_changes)} changes with extractable resources/methods")
    
    if debug and len(relevant_changes) > 0:
        print("\nDEBUG: Sample of relevant changes (first 5):")
        for i, change in enumerate(relevant_changes[:5], 1):
            print(f"\n  {i}. Title: {change.title}")
            print(f"     API Version: {change.api_version}")
            print(f"     Resources: {change.resources}")
            print(f"     Methods: {change.methods}")
            print(f"     Endpoints: {change.endpoints}")
            print(f"     Doc URLs: {len(change.doc_urls)} links")
    
    return relevant_changes


def main():
    """Main entry point for the changelog parser."""
    parser = argparse.ArgumentParser(
        description="Parse Microsoft Graph API Changelog RSS feed using URL-based extraction"
    )
    parser.add_argument(
        '--output',
        type=str,
        default='changelog-data.json',
        help='Output JSON file path'
    )
    parser.add_argument(
        '--lookback-days',
        type=int,
        default=30,
        help='Number of days to look back for changes (default: 30)'
    )
    parser.add_argument(
        '--url',
        type=str,
        default=RSS_FEED_URL,
        help='RSS feed URL (default: Microsoft Graph changelog)'
    )
    parser.add_argument(
        '--debug',
        action='store_true',
        help='Enable debug output'
    )
    parser.add_argument(
        '--verbose',
        action='store_true',
        help='Enable verbose parsing output'
    )
    parser.add_argument(
        '--provider-resources',
        type=str,
        help='Path to provider endpoints JSON for comparison filtering'
    )
    
    args = parser.parse_args()
    
    try:
        # Load provider resources if provided
        provider_resources = None
        if args.provider_resources:
            with open(args.provider_resources, 'r', encoding='utf-8') as f:
                provider_data = json.load(f)
                # Extract all Graph resource types from provider
                provider_resources = set(provider_data.get('lookup', {}).get('operations', {}).keys())
                print(f"Loaded {len(provider_resources)} provider resources for comparison")
        
        # Parse RSS feed
        changes = parse_rss_feed(
            args.url, 
            args.lookback_days, 
            debug=args.debug, 
            verbose=args.verbose,
            provider_resources=provider_resources
        )
        
        # Convert to JSON
        output_data = {
            'generated_at': datetime.now().isoformat(),
            'lookback_days': args.lookback_days,
            'total_changes': len(changes),
            'extraction_method': 'url_based',
            'changes': [c.to_dict() for c in changes]
        }
        
        # Write output
        with open(args.output, 'w', encoding='utf-8') as f:
            json.dump(output_data, f, indent=2)
        
        print(f"\n✓ Successfully wrote {len(changes)} changes to {args.output}")
        
        # Print summary
        print("\n" + "="*60)
        print("SUMMARY")
        print("="*60)
        print(f"Total changes: {len(changes)}")
        print("Extraction method: URL-based (from documentation links)")
        
        # Count unique resources
        all_resources = set()
        all_methods = set()
        for change in changes:
            all_resources.update(change.resources)
            all_methods.update(change.methods)
        
        print(f"\nUnique resources extracted: {len(all_resources)}")
        print(f"Unique methods extracted: {len(all_methods)}")
        
        print("\nBy API version:")
        version_counts = {}
        for change in changes:
            version_counts[change.api_version] = version_counts.get(change.api_version, 0) + 1
        for version, count in sorted(version_counts.items()):
            print(f"  - {version}: {count}")
        
        print("\nBy change type:")
        type_counts = {}
        for change in changes:
            change_type = change.change_type or 'unknown'
            type_counts[change_type] = type_counts.get(change_type, 0) + 1
        for change_type, count in sorted(type_counts.items(), key=lambda x: x[1], reverse=True):
            print(f"  - {change_type}: {count}")
        
        print("\nBy impact type:")
        impact_counts = {
            'new_provider_resource': 0,
            'schema_update': 0,
            'new_operation': 0,
            'breaking_change': 0,
            'deprecation_warning': 0,
            'uncategorized': 0
        }
        for change in changes:
            if change.change_actions:
                for action in change.change_actions:
                    impact = action.get('impact', 'uncategorized')
                    impact_counts[impact] = impact_counts.get(impact, 0) + 1
            else:
                impact_counts['uncategorized'] += 1
        
        for impact, count in sorted(impact_counts.items(), key=lambda x: x[1], reverse=True):
            if count > 0:
                print(f"  - {impact}: {count}")
        
        print(f"\nWith CRUD/minimal operations: {sum(1 for c in changes if c.supports_crud_or_minimal())}")
    
    except (OSError, json.JSONDecodeError, KeyError, ValueError) as e:
        print(f"Error: {e}", file=sys.stderr)
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()

