#!/usr/bin/env python3
"""
Parse Microsoft Graph API Changelog RSS feed and extract API changes.

This script fetches the Microsoft Graph changelog RSS feed, parses it,
and extracts relevant API changes including new resources, methods, and endpoints
across all Microsoft Graph API service areas (Device Management, Identity & Access,
Applications, Security, Groups, Users, Conditional Access, etc.).
"""

import argparse
import json
import re
import sys
from datetime import datetime, timedelta
from typing import Dict, List, Optional

try:
    import feedparser
    from bs4 import BeautifulSoup
except ImportError:
    print("Error: Required packages not installed. Run: pip install feedparser beautifulsoup4 lxml")
    sys.exit(1)


RSS_FEED_URL = "https://developer.microsoft.com/en-us/graph/changelog/rss"

# API categories relevant to the Microsoft 365 Terraform Provider
# Covers all major service areas: Device Management, Identity, Security, 
# Applications, Groups, Conditional Access, and more
RELEVANT_CATEGORIES = [
    "Device and app management",
    "Devices and app management",
    "Identity and access",
    "Security",
    "Agents",
    "Cloud communications",
    "Teamwork and communications",
    "Teamwork",
    "Calendar",
    "Files",
    "Applications",
    "Users",
    "Groups",
]


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
        self.resources = []
        self.methods = []
        self.endpoints = []
        self.properties = []
        self.change_type = None  # 'added', 'removed', 'deprecated', 'updated'
        
    def parse_description(self):
        """Parse the HTML description to extract API details."""
        soup = BeautifulSoup(self.description, 'html.parser')
        
        for div in soup.find_all('div'):
            text = div.get_text(strip=True)
            
            # Detect change type
            if 'Added the' in text or 'added the' in text:
                self.change_type = 'added'
            elif 'Removed the' in text or 'removed the' in text:
                self.change_type = 'removed'
            elif 'Deprecated the' in text or 'deprecated the' in text:
                self.change_type = 'deprecated'
            elif 'Updated the' in text or 'updated the' in text:
                self.change_type = 'updated'
            
            # Extract resources
            resource_match = re.search(r'the\s+<a[^>]*>([a-zA-Z0-9_]+)</a>\s+resource', text)
            if resource_match:
                resource_name = resource_match.group(1)
                
                # Add ALL resources - we'll filter for relevance later in is_relevant()
                self.resources.append(resource_name)
                
                # Extract link to documentation
                link_tag = div.find('a')
                if link_tag and link_tag.get('href'):
                    doc_url = link_tag['href']
                    # Extract endpoint from documentation URL
                    endpoint = self._extract_endpoint_from_url(doc_url)
                    if endpoint:
                        self.endpoints.append(endpoint)
            
            # Extract methods
            method_match = re.search(r'the\s+<a[^>]*>([a-zA-Z0-9_]+)</a>\s+method', text)
            if method_match:
                method_name = method_match.group(1)
                self.methods.append(method_name)
                
                # Extract associated resource
                resource_match = re.search(r'to the\s+<a[^>]*>([a-zA-Z0-9_]+)</a>\s+resource', text)
                if resource_match:
                    resource_name = resource_match.group(1)
                    # Add all resource/method combinations
                    endpoint = f"{resource_name}/{method_name}"
                    self.endpoints.append(endpoint)
            
            # Extract properties
            property_match = re.search(r'the\s+<b>([a-zA-Z0-9_]+)</b>\s+property', text)
            if property_match:
                property_name = property_match.group(1)
                self.properties.append(property_name)
    
    def _extract_endpoint_from_url(self, url: str) -> Optional[str]:
        """Extract API endpoint pattern from documentation URL."""
        # Example: https://learn.microsoft.com/en-us/graph/api/resources/intune-...
        # or: https://learn.microsoft.com/en-us/graph/api/device-post-...
        
        if '/api/resources/' in url:
            # Resource documentation
            parts = url.split('/api/resources/')
            if len(parts) > 1:
                resource = parts[1].split('?')[0].split('#')[0]
                return f"resources/{resource}"
        
        elif '/api/' in url:
            # Method documentation
            parts = url.split('/api/')
            if len(parts) > 1:
                method = parts[1].split('?')[0].split('#')[0]
                return f"api/{method}"
        
        return None
    
    def is_relevant(self, debug: bool = False) -> bool:
        """Determine if this change is relevant to the Terraform provider."""
        # Check if it's in a relevant category
        if not any(cat in self.categories for cat in RELEVANT_CATEGORIES):
            if debug:
                print(f"    DEBUG: Filtered - categories {self.categories} not in {RELEVANT_CATEGORIES}")
            return False
        
        # Check if it's a beta-only change (we might want to skip these)
        if self.api_version == 'beta':
            # For now, include beta changes but mark them
            pass
        
        # Expanded keywords covering all service areas the provider implements
        relevant_keywords = [
            # Device Management
            'device', 'intune', 'management', 'mdm', 'mam', 'enrollment',
            
            # Policies & Configuration
            'policy', 'policies', 'configuration', 'compliance', 'conditional',
            'assignment', 'remediation', 'setting',
            
            # Applications
            'app', 'application', 'mobile', 'mobileapp', 'appmanagement',
            
            # Identity & Access
            'identity', 'user', 'group', 'role', 'permission', 'authentication',
            'authorization', 'access', 'entitlement', 'governance',
            
            # Security
            'security', 'threat', 'protection', 'defender', 'vulnerability',
            'attack', 'risk', 'incident', 'alert',
            
            # Directory & Users
            'directory', 'azuread', 'entra', 'tenant', 'domain',
            
            # Service Principals & Apps
            'serviceprincipal', 'oauth', 'consent', 'api', 'permission',
            
            # Conditional Access
            'conditionalaccess', 'ca', 'mfa', 'authentication',
            
            # Microsoft 365 Services
            'teams', 'sharepoint', 'onedrive', 'exchange', 'calendar',
            
            # Administrative Units & Management
            'administrativeunit', 'organization', 'subscription',
        ]
        
        full_text = f"{self.title} {self.description} {' '.join(self.resources)} {' '.join(self.endpoints)}".lower()
        
        # If any keyword matches, it's relevant
        if any(keyword in full_text for keyword in relevant_keywords):
            return True
        
        # Also consider it relevant if it's in a relevant category and has resources or endpoints
        # This catches new resources even if they don't match keywords yet
        if (self.resources or self.endpoints) and self.change_type == 'added':
            return True
        
        return False
    
    def supports_crud_or_minimal(self) -> bool:
        """Check if this change supports full CRUD or at least update+get."""
        # This is a heuristic - we look for methods that suggest CRUD operations
        crud_methods = ['create', 'get', 'list', 'update', 'patch', 'delete', 'post', 'put']
        
        found_methods = set()
        for method in self.methods:
            method_lower = method.lower()
            for crud in crud_methods:
                if crud in method_lower:
                    found_methods.add(crud)
        
        # Check description for HTTP methods
        if 'POST' in self.description or 'post' in self.description.lower():
            found_methods.add('create')
        if 'GET' in self.description or 'get' in self.description.lower():
            found_methods.add('get')
        if 'PATCH' in self.description or 'patch' in self.description.lower():
            found_methods.add('update')
        if 'PUT' in self.description or 'put' in self.description.lower():
            found_methods.add('update')
        if 'DELETE' in self.description or 'delete' in self.description.lower():
            found_methods.add('delete')
        
        # Full CRUD: create, read (get), update, delete
        has_full_crud = all(op in found_methods for op in ['create', 'get', 'update', 'delete'])
        
        # Minimal: update and get
        has_minimal = 'get' in found_methods and 'update' in found_methods
        
        return has_full_crud or has_minimal
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            'guid': self.guid,
            'title': self.title,
            'description': self.description[:500],  # Truncate for readability
            'pub_date': self.pub_date,
            'categories': self.categories,
            'api_version': self.api_version,
            'resources': self.resources,
            'methods': self.methods,
            'endpoints': self.endpoints,
            'properties': self.properties,
            'change_type': self.change_type,
            'is_relevant': self.is_relevant(),
            'supports_crud_or_minimal': self.supports_crud_or_minimal()
        }


def parse_rss_feed(url: str, lookback_days: int = 30, debug: bool = False, verbose: bool = False) -> List[GraphAPIChange]:
    """Parse the RSS feed and return a list of API changes."""
    print(f"Fetching RSS feed from {url}...")
    
    feed = feedparser.parse(url)
    
    if feed.bozo:
        print(f"Warning: Feed parsing encountered an error: {feed.bozo_exception}")
    
    print(f"Found {len(feed.entries)} entries in feed")
    
    if debug:
        print(f"\nDEBUG: Cutoff date: {datetime.now() - timedelta(days=lookback_days)}")
        print(f"DEBUG: Relevant categories: {RELEVANT_CATEGORIES}\n")
    
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
        
        # Parse description to extract details
        change.parse_description()
        
        changes.append(change)
        
        if verbose and len(changes) <= 5:  # Show first 5 in verbose mode
            print(f"\n  Sample entry {len(changes)}:")
            print(f"    Title: {change.title}")
            print(f"    Categories: {change.categories}")
            print(f"    Resources: {change.resources}")
            print(f"    Methods: {change.methods}")
    
    print(f"Parsed {len(changes)} changes within the last {lookback_days} days")
    
    # Filter to relevant changes with debug info
    relevant_changes = []
    filtered_out = []
    
    for change in changes:
        if change.is_relevant():
            relevant_changes.append(change)
        else:
            filtered_out.append(change)
    
    print(f"Found {len(relevant_changes)} relevant changes for Microsoft Graph API")
    
    if debug and len(filtered_out) > 0:
        print(f"\nDEBUG: Filtered out {len(filtered_out)} changes")
        print("\nDEBUG: Sample of filtered changes (first 10):")
        for i, change in enumerate(filtered_out[:10], 1):
            print(f"\n  {i}. Title: {change.title}")
            print(f"     Categories: {change.categories}")
            print(f"     API Version: {change.api_version}")
            print(f"     Resources: {change.resources}")
            print(f"     Endpoints: {change.endpoints}")
            print(f"     Change Type: {change.change_type}")
            
            # Check why it was filtered
            has_category = any(cat in change.categories for cat in RELEVANT_CATEGORIES)
            print(f"     Has relevant category: {has_category}")
            if not has_category:
                print(f"     -> Filtered: Not in relevant categories")
    
    if debug and len(relevant_changes) > 0:
        print(f"\nDEBUG: Sample of relevant changes (first 5):")
        for i, change in enumerate(relevant_changes[:5], 1):
            print(f"\n  {i}. Title: {change.title}")
            print(f"     Categories: {change.categories}")
            print(f"     Resources: {change.resources}")
            print(f"     Endpoints: {change.endpoints}")
    
    return relevant_changes


def main():
    parser = argparse.ArgumentParser(
        description="Parse Microsoft Graph API Changelog RSS feed"
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
        help='Enable debug output to see filtering decisions'
    )
    parser.add_argument(
        '--verbose',
        action='store_true',
        help='Enable verbose output with detailed parsing info'
    )
    
    args = parser.parse_args()
    
    try:
        # Parse RSS feed
        changes = parse_rss_feed(args.url, args.lookback_days, debug=args.debug, verbose=args.verbose)
        
        # Convert to JSON
        output_data = {
            'generated_at': datetime.now().isoformat(),
            'lookback_days': args.lookback_days,
            'total_changes': len(changes),
            'changes': [c.to_dict() for c in changes]
        }
        
        # Write output
        with open(args.output, 'w') as f:
            json.dump(output_data, f, indent=2)
        
        print(f"\n✓ Successfully wrote {len(changes)} changes to {args.output}")
        
        # Print summary
        print("\n" + "="*60)
        print("SUMMARY")
        print("="*60)
        print(f"Total changes: {len(changes)}")
        print(f"By category:")
        category_counts = {}
        for change in changes:
            for cat in change.categories:
                category_counts[cat] = category_counts.get(cat, 0) + 1
        for cat, count in sorted(category_counts.items(), key=lambda x: x[1], reverse=True):
            print(f"  - {cat}: {count}")
        
        print(f"\nBy change type:")
        type_counts = {}
        for change in changes:
            change_type = change.change_type or 'unknown'
            type_counts[change_type] = type_counts.get(change_type, 0) + 1
        for change_type, count in sorted(type_counts.items(), key=lambda x: x[1], reverse=True):
            print(f"  - {change_type}: {count}")
        
        print(f"\nSupporting CRUD or minimal operations: {sum(1 for c in changes if c.supports_crud_or_minimal())}")
        
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()