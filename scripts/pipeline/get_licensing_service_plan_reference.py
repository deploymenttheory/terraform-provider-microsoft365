#!/usr/bin/env python3
"""
Script to fetch and parse Microsoft licensing service plan reference data.

This script fetches the licensing service plan reference page from Microsoft Learn,
parses the product table, and generates Go constants for use in the Terraform provider.

Usage:
    python get_licensing_service_plan_reference.py [--output FILE] [--format go|csv|json]

Arguments:
    --output, -o    Output file path (default: stdout)
    --format, -f    Output format: go, csv, or json (default: go)
    --verbose, -v   Enable verbose logging

Reference:
    https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference
"""

import argparse
import json
import sys
import re
from typing import List, Dict, Optional
from urllib.request import urlopen, Request
from urllib.error import URLError, HTTPError
from html.parser import HTMLParser


class LicenseTableParser(HTMLParser):
    """HTML parser to extract licensing data from Microsoft Learn documentation."""

    def __init__(self):
        super().__init__()
        self.in_table = False
        self.in_row = False
        self.in_cell = False
        self.current_row = []
        self.current_cell = []
        self.rows = []
        self.cell_count = 0

    def handle_starttag(self, tag, attrs):
        if tag == 'table':
            self.in_table = True
        elif tag == 'tr' and self.in_table:
            self.in_row = True
            self.current_row = []
            self.cell_count = 0
        elif (tag == 'td' or tag == 'th') and self.in_row:
            self.in_cell = True
            self.current_cell = []

    def handle_endtag(self, tag):
        if tag == 'table':
            self.in_table = False
        elif tag == 'tr' and self.in_row:
            self.in_row = False
            # Only add rows with at least 3 cells (Product name, String ID, GUID)
            if len(self.current_row) >= 3:
                self.rows.append(self.current_row)
        elif (tag == 'td' or tag == 'th') and self.in_cell:
            self.in_cell = False
            cell_text = ''.join(self.current_cell).strip()
            self.current_row.append(cell_text)
            self.cell_count += 1

    def handle_data(self, data):
        if self.in_cell:
            self.current_cell.append(data)


def fetch_page_content(url: str, verbose: bool = False) -> str:
    """
    Fetch the HTML content from the Microsoft Learn page.

    Args:
        url: URL to fetch
        verbose: Enable verbose logging

    Returns:
        HTML content as string

    Raises:
        HTTPError: If HTTP request fails
        URLError: If connection fails
    """
    if verbose:
        print(f"Fetching content from: {url}", file=sys.stderr)

    headers = {
        'User-Agent': 'Mozilla/5.0 (compatible; TerraformProvider/1.0; +https://github.com/deploymenttheory/terraform-provider-microsoft365)'
    }
    
    req = Request(url, headers=headers)
    
    try:
        with urlopen(req, timeout=30) as response:
            html_content = response.read().decode('utf-8')
            if verbose:
                print(f"Successfully fetched {len(html_content)} bytes", file=sys.stderr)
            return html_content
    except HTTPError as e:
        print(f"HTTP Error {e.code}: {e.reason}", file=sys.stderr)
        raise
    except URLError as e:
        print(f"URL Error: {e.reason}", file=sys.stderr)
        raise


def parse_service_plans(service_plans_text: str) -> List[Dict[str, str]]:
    """
    Parse service plans from the table cell text.
    
    Expected format: "PLAN_ID (guid)" with multiple entries separated by newlines or spaces
    
    Args:
        service_plans_text: Raw text containing service plan IDs and GUIDs
        
    Returns:
        List of dictionaries with 'id' and 'guid' keys
    """
    plans = []
    
    # Pattern to match: PLAN_ID (guid)
    pattern = r'([A-Z0-9_]+)\s*\(([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\)'
    
    matches = re.finditer(pattern, service_plans_text, re.IGNORECASE)
    
    for match in matches:
        plans.append({
            'id': match.group(1).strip(),
            'guid': match.group(2).lower()
        })
    
    return plans


def parse_service_plan_friendly_names(service_plans_text: str) -> List[Dict[str, str]]:
    """
    Parse service plan friendly names from the table cell text.
    
    Expected format: "Friendly Name (guid)" with multiple entries
    The friendly name can itself contain parentheses, e.g., "App Name (Version A) (guid)"
    
    Args:
        service_plans_text: Raw text containing service plan friendly names and GUIDs
        
    Returns:
        List of dictionaries with 'name' and 'guid' keys
    """
    plans = []
    
    # First, find all GUIDs with their positions
    guid_pattern = r'\(([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\)'
    guid_matches = list(re.finditer(guid_pattern, service_plans_text, re.IGNORECASE))
    
    for i, match in enumerate(guid_matches):
        guid = match.group(1).lower()
        
        # Determine the start position for this entry
        if i == 0:
            start_pos = 0
        else:
            # Start after the previous GUID's closing parenthesis
            start_pos = guid_matches[i-1].end()
        
        # The name is everything between start_pos and this GUID's opening parenthesis
        name_end_pos = match.start()
        name = service_plans_text[start_pos:name_end_pos].strip()
        
        # Remove any leading/trailing whitespace and newlines
        name = ' '.join(name.split())
        
        # Skip if name is empty or looks like a service plan ID (all caps with underscores)
        if name and not re.match(r'^[A-Z0-9_]+$', name):
            plans.append({
                'name': name,
                'guid': guid
            })
    
    return plans


def parse_license_table(html_content: str, verbose: bool = False) -> List[Dict[str, str]]:
    """
    Parse the HTML content and extract license data from the table.

    Args:
        html_content: HTML content as string
        verbose: Enable verbose logging

    Returns:
        List of dictionaries containing product_name, string_id, guid, and service_plans_included
        Each service_plan_included entry contains: id, name, and guid
    """
    parser = LicenseTableParser()
    parser.feed(html_content)

    if verbose:
        print(f"Found {len(parser.rows)} total rows in tables", file=sys.stderr)

    licenses = []
    
    for row in parser.rows:
        if len(row) < 3:
            continue
            
        product_name = row[0].strip()
        string_id = row[1].strip()
        guid = row[2].strip()

        # Skip header rows
        if product_name.lower() in ['product name', 'product', '']:
            continue
            
        # Validate GUID format (basic check)
        guid_pattern = r'^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
        if not re.match(guid_pattern, guid.lower()):
            if verbose:
                print(f"Skipping invalid GUID format: {guid}", file=sys.stderr)
            continue

        license_data = {
            'product_name': product_name,
            'string_id': string_id,
            'guid': guid.lower(),
            'service_plans_included': []
        }

        # Parse service plans included if column exists (column index 3)
        service_plan_ids = []
        if len(row) > 3 and row[3].strip():
            service_plan_ids = parse_service_plans(row[3])
        
        # Parse service plans included friendly names if column exists (column index 4)
        service_plan_names = []
        if len(row) > 4 and row[4].strip():
            service_plan_names = parse_service_plan_friendly_names(row[4])
        
        # Merge service plan IDs with their friendly names by matching GUIDs
        service_plans_map = {}
        
        # First, add all IDs
        for sp in service_plan_ids:
            service_plans_map[sp['guid']] = {
                'id': sp['id'],
                'guid': sp['guid']
            }
        
        # Then, merge in the friendly names
        for sp in service_plan_names:
            if sp['guid'] in service_plans_map:
                service_plans_map[sp['guid']]['name'] = sp['name']
            else:
                # If we have a name but no ID, still include it
                service_plans_map[sp['guid']] = {
                    'id': '',  # No ID found
                    'name': sp['name'],
                    'guid': sp['guid']
                }
        
        # Convert map back to list, ensuring id, name, guid order
        license_data['service_plans_included'] = [
            {
                'id': sp.get('id', ''),
                'name': sp.get('name', ''),
                'guid': sp['guid']
            }
            for sp in service_plans_map.values()
        ]

        licenses.append(license_data)

    if verbose:
        print(f"Extracted {len(licenses)} valid license entries", file=sys.stderr)
        total_service_plans = sum(len(lic['service_plans_included']) for lic in licenses)
        print(f"Extracted {total_service_plans} total service plan mappings", file=sys.stderr)

    return licenses


def sanitize_const_name(text: str) -> str:
    """
    Convert product name to a valid Go constant name.

    Args:
        text: Product name text

    Returns:
        Sanitized constant name
    """
    # Remove special characters and convert to title case
    # Example: "Microsoft 365 E3 (no Teams)" -> "M365E3NoTeams"
    
    # Replace common patterns
    text = text.replace('Microsoft 365', 'M365')
    text = text.replace('Office 365', 'O365')
    text = text.replace('Microsoft', 'MS')
    text = text.replace('Enterprise Mobility + Security', 'EMS')
    text = text.replace('Microsoft Entra ID', 'EntraID')
    text = text.replace('Azure Active Directory', 'AAD')
    
    # Remove parentheses content but keep meaningful parts
    text = re.sub(r'\(no Teams\)', 'NoTeams', text)
    text = re.sub(r'\(.*?\)', '', text)
    
    # Remove special characters and spaces
    text = re.sub(r'[^a-zA-Z0-9]', '', text)
    
    # Ensure it starts with a letter
    if text and text[0].isdigit():
        text = 'Sku' + text
    
    return text


def group_licenses_by_category(licenses: List[Dict[str, str]]) -> Dict[str, List[Dict[str, str]]]:
    """
    Group licenses by category based on product names.

    Args:
        licenses: List of license dictionaries

    Returns:
        Dictionary with category names as keys and lists of licenses as values
    """
    categories = {
        'Microsoft 365 Enterprise': [],
        'Microsoft 365 Business': [],
        'Office 365 Enterprise': [],
        'Exchange Online': [],
        'Microsoft Entra ID': [],
        'Enterprise Mobility + Security': [],
        'Microsoft Intune': [],
        'Power Platform': [],
        'Project and Visio': [],
        'Microsoft Defender': [],
        'Other': []
    }

    for license in licenses:
        product = license['product_name']
        
        if 'Microsoft 365 E' in product or 'M365 E' in product:
            categories['Microsoft 365 Enterprise'].append(license)
        elif 'Microsoft 365 Business' in product or 'M365 Business' in product:
            categories['Microsoft 365 Business'].append(license)
        elif 'Office 365 E' in product or 'O365 E' in product:
            categories['Office 365 Enterprise'].append(license)
        elif 'Exchange Online' in product or 'Exchange' in product:
            categories['Exchange Online'].append(license)
        elif 'Entra' in product or 'Azure AD' in product or 'Azure Active Directory' in product:
            categories['Microsoft Entra ID'].append(license)
        elif 'EMS' in product or 'Enterprise Mobility' in product:
            categories['Enterprise Mobility + Security'].append(license)
        elif 'Intune' in product:
            categories['Microsoft Intune'].append(license)
        elif 'Power BI' in product or 'Power Apps' in product or 'Power Automate' in product:
            categories['Power Platform'].append(license)
        elif 'Project' in product or 'Visio' in product:
            categories['Project and Visio'].append(license)
        elif 'Defender' in product or 'Security' in product:
            categories['Microsoft Defender'].append(license)
        else:
            categories['Other'].append(license)

    # Remove empty categories
    return {k: v for k, v in categories.items() if v}


def generate_go_constants(licenses: List[Dict[str, str]], verbose: bool = False) -> str:
    """
    Generate Go constant declarations from license data.

    Args:
        licenses: List of license dictionaries
        verbose: Enable verbose logging

    Returns:
        Go source code as string
    """
    output = []
    
    output.append("// Code generated by get_licensing_service_plan_reference.py; DO NOT EDIT.")
    output.append("// Source: https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference")
    output.append("")
    output.append("package constants")
    output.append("")
    output.append("// ============================================================================")
    output.append("// Microsoft License SKU Constants")
    output.append("// ============================================================================")
    output.append("// This file contains constants for Microsoft license SKUs including:")
    output.append("// - Product names (as displayed in Azure Portal)")
    output.append("// - String IDs (used by PowerShell v1.0 and skuPartNumber in Graph API)")
    output.append("// - GUIDs (used by skuId in Graph API)")
    output.append("//")
    output.append("// Reference: https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference")
    output.append("// ============================================================================")
    output.append("")

    # Group licenses by category
    grouped = group_licenses_by_category(licenses)
    
    const_names_used = set()
    mapping_entries = []

    for category, category_licenses in grouped.items():
        output.append("// ============================================================================")
        output.append(f"// {category} Licenses")
        output.append("// ============================================================================")
        output.append("")

        for license in category_licenses:
            base_name = sanitize_const_name(license['product_name'])
            
            # Handle duplicates by appending a number
            const_name = base_name
            counter = 2
            while const_name in const_names_used:
                const_name = f"{base_name}{counter}"
                counter += 1
            
            const_names_used.add(const_name)

            output.append(f"// {license['product_name']}")
            output.append(f"// Product Name: {license['product_name']}")
            output.append(f"// String ID: {license['string_id']}")
            output.append(f"// GUID: {license['guid']}")
            
            # Add service plans included information if available
            if license.get('service_plans_included'):
                output.append("// Service Plans Included:")
                for sp in license['service_plans_included'][:10]:  # Limit to first 10 to avoid excessive comments
                    if sp.get('name'):
                        output.append(f"//   - {sp['id']} - {sp['name']} ({sp['guid']})")
                    else:
                        output.append(f"//   - {sp['id']} ({sp['guid']})")
                if len(license['service_plans_included']) > 10:
                    output.append(f"//   ... and {len(license['service_plans_included']) - 10} more")
            
            output.append("const (")
            output.append(f"\tSku{const_name}ProductName = \"{license['product_name']}\"")
            output.append(f"\tSku{const_name}StringID    = \"{license['string_id']}\"")
            output.append(f"\tSku{const_name}GUID        = \"{license['guid']}\"")
            output.append(")")
            output.append("")

            # Store for mapping functions
            mapping_entries.append({
                'const_name': const_name,
                'string_id': license['string_id'],
                'guid': license['guid']
            })

    # Generate helper functions
    output.append("// ============================================================================")
    output.append("// Helper Functions")
    output.append("// ============================================================================")
    output.append("")
    output.append("// GetSkuGUIDByStringID returns the GUID for a given SKU String ID")
    output.append("func GetSkuGUIDByStringID(stringID string) string {")
    output.append("\tskuMap := map[string]string{")
    
    for entry in mapping_entries:
        output.append(f"\t\tSku{entry['const_name']}StringID: Sku{entry['const_name']}GUID,")
    
    output.append("\t}")
    output.append("")
    output.append("\treturn skuMap[stringID]")
    output.append("}")
    output.append("")
    output.append("// GetSkuStringIDByGUID returns the String ID for a given SKU GUID")
    output.append("func GetSkuStringIDByGUID(guid string) string {")
    output.append("\tguidMap := map[string]string{")
    
    for entry in mapping_entries:
        output.append(f"\t\tSku{entry['const_name']}GUID: Sku{entry['const_name']}StringID,")
    
    output.append("\t}")
    output.append("")
    output.append("\treturn guidMap[guid]")
    output.append("}")
    output.append("")

    if verbose:
        print(f"Generated {len(mapping_entries)} Go constants", file=sys.stderr)

    return '\n'.join(output)


def generate_csv(licenses: List[Dict[str, str]]) -> str:
    """
    Generate CSV output from license data.

    Args:
        licenses: List of license dictionaries

    Returns:
        CSV data as string
    """
    output = ["Product Name,String ID,GUID,Service Plans Included"]
    
    for license in licenses:
        # Escape commas in product names
        product_name = f'"{license["product_name"]}"' if ',' in license['product_name'] else license['product_name']
        
        # Format service plans included as "ID - Name (GUID); ID - Name (GUID)"
        service_plans_list = []
        for sp in license.get('service_plans_included', []):
            if sp.get('name'):
                service_plans_list.append(f"{sp['id']} - {sp['name']} ({sp['guid']})")
            else:
                service_plans_list.append(f"{sp['id']} ({sp['guid']})")
        
        service_plans_str = '; '.join(service_plans_list)
        service_plans_str = f'"{service_plans_str}"' if service_plans_str else ''
        
        output.append(f"{product_name},{license['string_id']},{license['guid']},{service_plans_str}")
    
    return '\n'.join(output)


def generate_json(licenses: List[Dict[str, str]]) -> str:
    """
    Generate JSON output from license data.

    Args:
        licenses: List of license dictionaries

    Returns:
        JSON data as string
    """
    return json.dumps(licenses, indent=2)


def main():
    """Main entry point for the script."""
    parser = argparse.ArgumentParser(
        description='Fetch and parse Microsoft licensing service plan reference data',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
Examples:
  # Generate Go constants to stdout
  python get_licensing_service_plan_reference.py

  # Generate Go constants to file
  python get_licensing_service_plan_reference.py -o licensing_service_plan.go

  # Generate CSV format
  python get_licensing_service_plan_reference.py -f csv -o licenses.csv

  # Generate JSON format
  python get_licensing_service_plan_reference.py -f json -o licenses.json

  # Verbose mode
  python get_licensing_service_plan_reference.py -v -o licensing_service_plan.go
        '''
    )
    
    parser.add_argument(
        '--output', '-o',
        type=str,
        help='Output file path (default: stdout)',
        default=None
    )
    
    parser.add_argument(
        '--format', '-f',
        type=str,
        choices=['go', 'csv', 'json'],
        default='go',
        help='Output format (default: go)'
    )
    
    parser.add_argument(
        '--verbose', '-v',
        action='store_true',
        help='Enable verbose logging'
    )

    parser.add_argument(
        '--url',
        type=str,
        default='https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference',
        help='URL to fetch license data from (default: Microsoft Learn documentation)'
    )

    args = parser.parse_args()

    try:
        # Fetch page content
        html_content = fetch_page_content(args.url, args.verbose)

        # Parse license table
        licenses = parse_license_table(html_content, args.verbose)

        if not licenses:
            print("Error: No license data found in the page", file=sys.stderr)
            sys.exit(1)

        # Generate output based on format
        if args.format == 'go':
            output = generate_go_constants(licenses, args.verbose)
        elif args.format == 'csv':
            output = generate_csv(licenses)
        elif args.format == 'json':
            output = generate_json(licenses)
        else:
            print(f"Error: Unsupported format: {args.format}", file=sys.stderr)
            sys.exit(1)

        # Write output
        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(output)
            if args.verbose:
                print(f"Output written to: {args.output}", file=sys.stderr)
        else:
            print(output)

        if args.verbose:
            print("Script completed successfully", file=sys.stderr)

    except KeyboardInterrupt:
        print("\nScript interrupted by user", file=sys.stderr)
        sys.exit(130)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        if args.verbose:
            import traceback
            traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()

