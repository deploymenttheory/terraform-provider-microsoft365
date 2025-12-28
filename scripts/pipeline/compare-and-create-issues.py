#!/usr/bin/env python3
"""
Compare Graph API changelog with provider implementation and create GitHub issues for gaps.

This script compares the API changes from the Microsoft Graph changelog with
the current provider implementation and creates GitHub issues for identified gaps.
"""

import argparse
import json
import os
import re
import subprocess
import sys
from datetime import datetime
from typing import Dict, List, Optional


class Gap:
    """Represents a gap between Graph API and provider implementation."""
    
    def __init__(self, change: Dict, gap_type: str):
        self.change = change
        self.gap_type = gap_type  # 'new_resource', 'updated_resource', 'missing_operation'
        self.priority = self._calculate_priority()
        self.title = self._generate_title()
        self.body = self._generate_body()
        self.labels = self._generate_labels()
    
    def _calculate_priority(self) -> str:
        """Calculate priority based on change characteristics."""
        # High priority: New resources with full CRUD, v1.0 API
        if (self.gap_type == 'new_resource' and 
            self.change.get('supports_crud_or_minimal', False) and 
            self.change.get('api_version') == 'v1.0'):
            return 'high'
        
        # Medium priority: Updates to existing resources, beta with CRUD
        if (self.gap_type == 'updated_resource' or 
            (self.change.get('supports_crud_or_minimal', False) and 
             self.change.get('api_version') == 'beta')):
            return 'medium'
        
        # Low priority: Everything else
        return 'low'
    
    def _generate_title(self) -> str:
        """Generate a GitHub issue title."""
        change_title = self.change.get('title', 'Unknown')
        api_version = self.change.get('api_version', 'unknown')
        
        if self.gap_type == 'new_resource':
            return f"[{api_version}] Implement new Graph API resource: {change_title}"
        elif self.gap_type == 'updated_resource':
            return f"[{api_version}] Update Graph API resource: {change_title}"
        else:
            return f"[{api_version}] Add missing operations: {change_title}"
    
    def _generate_body(self) -> str:
        """Generate GitHub issue body with details."""
        lines = []
        
        lines.append("## Summary")
        lines.append("")
        lines.append(f"The Microsoft Graph API changelog indicates a {self.gap_type.replace('_', ' ')} ")
        lines.append(f"that is not currently implemented in the Terraform provider.")
        lines.append("")
        
        lines.append("## API Change Details")
        lines.append("")
        lines.append(f"- **Category**: {', '.join(self.change.get('categories', ['Unknown']))}")
        lines.append(f"- **API Version**: `{self.change.get('api_version', 'unknown')}`")
        lines.append(f"- **Published**: {self.change.get('pub_date', 'Unknown')[:10]}")
        lines.append(f"- **Change Type**: {self.change.get('change_type', 'unknown').title()}")
        lines.append("")
        
        if self.change.get('resources'):
            lines.append("### Resources")
            for resource in self.change['resources']:
                lines.append(f"- `{resource}`")
            lines.append("")
        
        if self.change.get('methods'):
            lines.append("### Methods")
            for method in self.change['methods']:
                lines.append(f"- `{method}`")
            lines.append("")
        
        if self.change.get('endpoints'):
            lines.append("### Endpoints")
            for endpoint in self.change['endpoints']:
                lines.append(f"- `{endpoint}`")
            lines.append("")
        
        if self.change.get('properties'):
            lines.append("### Properties")
            for prop in self.change['properties']:
                lines.append(f"- `{prop}`")
            lines.append("")
        
        lines.append("## Description")
        lines.append("")
        # Clean up HTML from description
        desc = self.change.get('description', 'No description available')
        desc = re.sub(r'<[^>]+>', '', desc)  # Remove HTML tags
        desc = re.sub(r'&lt;', '<', desc)
        desc = re.sub(r'&gt;', '>', desc)
        desc = re.sub(r'&#xD;', '', desc)
        lines.append(desc[:500] + ('...' if len(desc) > 500 else ''))
        lines.append("")
        
        lines.append("## Implementation Checklist")
        lines.append("")
        
        if self.gap_type == 'new_resource':
            lines.append("- [ ] Review API documentation")
            lines.append("- [ ] Design Terraform resource schema")
            lines.append("- [ ] Implement CRUD operations")
            lines.append("- [ ] Add unit tests")
            lines.append("- [ ] Add acceptance tests")
            lines.append("- [ ] Add documentation")
            lines.append("- [ ] Add examples")
        elif self.gap_type == 'updated_resource':
            lines.append("- [ ] Review API changes")
            lines.append("- [ ] Update resource schema if needed")
            lines.append("- [ ] Update CRUD operations")
            lines.append("- [ ] Update tests")
            lines.append("- [ ] Update documentation")
        else:
            lines.append("- [ ] Identify missing operations")
            lines.append("- [ ] Implement missing operations")
            lines.append("- [ ] Add/update tests")
            lines.append("- [ ] Update documentation")
        
        lines.append("")
        lines.append("## Additional Context")
        lines.append("")
        lines.append(f"- **GUID**: `{self.change.get('guid', 'N/A')}`")
        lines.append(f"- **Supports CRUD/Minimal**: {self.change.get('supports_crud_or_minimal', False)}")
        lines.append(f"- **Priority**: {self.priority.upper()}")
        lines.append("")
        lines.append("---")
        lines.append("*This issue was automatically created by the API changes monitor workflow.*")
        
        return '\n'.join(lines)
    
    def _generate_labels(self) -> List[str]:
        """Generate appropriate labels for the issue."""
        labels = ['enhancement', 'api-change']
        
        # Add priority label
        labels.append(f'priority-{self.priority}')
        
        # Add API version label
        if self.change.get('api_version'):
            labels.append(f"graph-{self.change['api_version']}")
        
        # Add category labels
        categories = self.change.get('categories', [])
        for cat in categories:
            if 'Device' in cat or 'device' in cat:
                labels.append('device-management')
            if 'Identity' in cat or 'identity' in cat:
                labels.append('identity')
            if 'Security' in cat or 'security' in cat:
                labels.append('security')
        
        # Add gap type label
        labels.append(self.gap_type.replace('_', '-'))
        
        return labels
    
    def to_dict(self) -> Dict:
        """Convert to dictionary for JSON serialization."""
        return {
            'title': self.title,
            'body': self.body,
            'labels': self.labels,
            'priority': self.priority,
            'gap_type': self.gap_type,
            'api_version': self.change.get('api_version'),
            'change_guid': self.change.get('guid')
        }


def compare_changes_with_provider(changelog_data: Dict, provider_data: Dict) -> List[Gap]:
    """Compare API changes with provider implementation and identify gaps."""
    gaps = []
    
    changes = changelog_data.get('changes', [])
    provider_lookup = provider_data.get('lookup', {})
    
    print(f"Comparing {len(changes)} API changes with provider implementation...")
    
    for change in changes:
        # Skip if not relevant
        if not change.get('is_relevant', False):
            continue
        
        # Check if this is about a new resource
        resources = change.get('resources', [])
        endpoints = change.get('endpoints', [])
        change_type = change.get('change_type', '')
        
        # Case 1: New resource added
        if change_type == 'added' and resources:
            # Check if any of these resources are implemented
            is_implemented = False
            for resource in resources:
                # Check in provider's Graph resources
                if resource in provider_lookup.get('operations', {}):
                    is_implemented = True
                    break
                
                # Check if resource name appears in endpoints
                for endpoint in provider_lookup.get('endpoints', {}).keys():
                    if resource.lower() in endpoint.lower():
                        is_implemented = True
                        break
            
            if not is_implemented and change.get('supports_crud_or_minimal', False):
                gap = Gap(change, 'new_resource')
                gaps.append(gap)
        
        # Case 2: Updated resource (new methods/properties added)
        elif change_type in ['added', 'updated'] and (resources or endpoints):
            # Check if the resource exists but might be missing new functionality
            is_fully_implemented = False
            is_partially_implemented = False
            
            for resource in resources:
                if resource in provider_lookup.get('operations', {}):
                    is_partially_implemented = True
                    # Check if all operations are supported
                    provider_ops = set(provider_lookup['operations'][resource].get('operations', []))
                    # If we have methods mentioned in the change, check them
                    if change.get('methods'):
                        # This is an update, might need attention
                        pass
                    else:
                        is_fully_implemented = True
            
            # If it's partially implemented and has new methods, it might need updating
            if is_partially_implemented and not is_fully_implemented and change.get('methods'):
                gap = Gap(change, 'updated_resource')
                gaps.append(gap)
        
        # Case 3: Deprecated resources (for awareness)
        elif change_type == 'deprecated':
            # Check if we're still using this deprecated resource
            for resource in resources:
                if resource in provider_lookup.get('operations', {}):
                    # We're using a deprecated resource - might need a gap
                    gap = Gap(change, 'updated_resource')
                    gap.priority = 'high'  # Deprecations are high priority
                    gaps.append(gap)
    
    print(f"Identified {len(gaps)} gaps")
    
    # Sort gaps by priority
    priority_order = {'high': 0, 'medium': 1, 'low': 2}
    gaps.sort(key=lambda g: (priority_order.get(g.priority, 3), g.title))
    
    return gaps


def create_github_issue(gap: Gap, repo: str, token: str, dry_run: bool = False) -> Optional[str]:
    """Create a GitHub issue for the gap."""
    if dry_run:
        print(f"[DRY RUN] Would create issue: {gap.title}")
        return None
    
    try:
        # Use GitHub CLI to create issue
        cmd = [
            'gh', 'issue', 'create',
            '--repo', repo,
            '--title', gap.title,
            '--body', gap.body,
            '--label', ','.join(gap.labels)
        ]
        
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            check=True
        )
        
        # Extract issue URL from output
        issue_url = result.stdout.strip()
        print(f"✓ Created issue: {issue_url}")
        return issue_url
    
    except subprocess.CalledProcessError as e:
        print(f"✗ Failed to create issue '{gap.title}': {e.stderr}", file=sys.stderr)
        return None
    except FileNotFoundError:
        print("✗ GitHub CLI (gh) not found. Please install it or set dry-run mode.", file=sys.stderr)
        return None


def check_existing_issues(gap: Gap, repo: str) -> bool:
    """Check if an issue for this gap already exists."""
    try:
        # Search for existing issues with similar title
        cmd = [
            'gh', 'issue', 'list',
            '--repo', repo,
            '--search', f'"{gap.change.get("title", "")}" in:title',
            '--state', 'open',
            '--json', 'title',
            '--limit', '10'
        ]
        
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            check=True
        )
        
        issues = json.loads(result.stdout)
        
        # Check if any existing issue matches
        for issue in issues:
            if gap.change.get('title', '') in issue.get('title', ''):
                return True
        
        return False
    
    except (subprocess.CalledProcessError, FileNotFoundError, json.JSONDecodeError):
        # If we can't check, assume it doesn't exist
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Compare API changes with provider implementation and create issues"
    )
    parser.add_argument(
        '--changelog',
        type=str,
        required=True,
        help='Path to changelog JSON file'
    )
    parser.add_argument(
        '--provider',
        type=str,
        required=True,
        help='Path to provider endpoints JSON file'
    )
    parser.add_argument(
        '--output',
        type=str,
        default='gaps-report.json',
        help='Output JSON file path'
    )
    parser.add_argument(
        '--create-issues',
        type=str,
        default='false',
        choices=['true', 'false'],
        help='Whether to create GitHub issues (default: false)'
    )
    parser.add_argument(
        '--check-existing',
        action='store_true',
        help='Check for existing issues before creating new ones'
    )
    
    args = parser.parse_args()
    
    try:
        # Load changelog data
        with open(args.changelog, 'r') as f:
            changelog_data = json.load(f)
        
        # Load provider data
        with open(args.provider, 'r') as f:
            provider_data = json.load(f)
        
        # Compare and identify gaps
        gaps = compare_changes_with_provider(changelog_data, provider_data)
        
        # Get repo and token from environment
        repo = os.environ.get('GITHUB_REPOSITORY', '')
        token = os.environ.get('GITHUB_TOKEN', '')
        
        # Create issues if requested
        created_issues = []
        should_create = args.create_issues.lower() == 'true'
        
        if should_create and repo and token:
            print(f"\nCreating GitHub issues for {len(gaps)} gaps...")
            
            for gap in gaps:
                # Check if issue already exists
                if args.check_existing and check_existing_issues(gap, repo):
                    print(f"⊘ Skipping duplicate issue: {gap.title}")
                    continue
                
                issue_url = create_github_issue(gap, repo, token, dry_run=False)
                if issue_url:
                    created_issues.append({
                        'url': issue_url,
                        'title': gap.title,
                        'priority': gap.priority
                    })
        elif should_create:
            print("Warning: Cannot create issues - GITHUB_REPOSITORY or GITHUB_TOKEN not set")
        
        # Generate report
        report = {
            'generated_at': datetime.now().isoformat(),
            'total_changes': len(changelog_data.get('changes', [])),
            'relevant_changes': sum(1 for c in changelog_data.get('changes', []) if c.get('is_relevant')),
            'already_implemented': changelog_data.get('total_changes', 0) - len(gaps),
            'gaps_identified': len(gaps),
            'issues_created': len(created_issues),
            'gaps': [g.to_dict() for g in gaps],
            'created_issues': created_issues
        }
        
        # Write report
        with open(args.output, 'w') as f:
            json.dump(report, f, indent=2)
        
        print(f"\n✓ Successfully wrote gaps report to {args.output}")
        
        # Print summary
        print("\n" + "="*60)
        print("GAPS ANALYSIS SUMMARY")
        print("="*60)
        print(f"Total API changes: {report['total_changes']}")
        print(f"Relevant changes: {report['relevant_changes']}")
        print(f"Already implemented: {report['already_implemented']}")
        print(f"Gaps identified: {report['gaps_identified']}")
        print(f"Issues created: {report['issues_created']}")
        
        if gaps:
            print("\nGaps by priority:")
            for priority in ['high', 'medium', 'low']:
                count = sum(1 for g in gaps if g.priority == priority)
                if count > 0:
                    print(f"  - {priority.upper()}: {count}")
            
            print("\nTop 5 gaps:")
            for i, gap in enumerate(gaps[:5], 1):
                priority_emoji = '🔴' if gap.priority == 'high' else '🟡' if gap.priority == 'medium' else '🟢'
                print(f"  {i}. {priority_emoji} {gap.title}")
    
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()

