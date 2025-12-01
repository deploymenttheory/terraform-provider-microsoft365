#!/usr/bin/env python3
"""
Fetch Microsoft 365 endpoint reference data from the official Microsoft API.

This script:
1. Downloads official M365 endpoints from Microsoft's JSON API for all clouds
2. Includes all service areas: MEM, Exchange, Skype, SharePoint, Common
3. Exports to JSON, Go constants, and CSV for each cloud

Usage:
    python fetch_microsoft365_endpoints.py [--output-dir DIR] [--format FORMAT] [--verbose]
"""

import argparse
import json
import logging
import re
import sys
import uuid
from dataclasses import dataclass, field
from pathlib import Path
from typing import List, Dict
from urllib.request import urlopen, Request
from urllib.error import URLError, HTTPError


# ============================================================================
# Configuration
# ============================================================================

# Official Microsoft 365 Endpoints API
# Reference: https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-ip-web-service
M365_CLOUDS = {
    'worldwide': 'https://endpoints.office.com/endpoints/Worldwide',
    'china': 'https://endpoints.office.com/endpoints/China',
    'usgov-dod': 'https://endpoints.office.com/endpoints/USGOVDoD',
    'usgov-gcc-high': 'https://endpoints.office.com/endpoints/USGOVGCCHigh'
}


# ============================================================================
# Setup Logging
# ============================================================================

def setup_logging(verbose: bool = False):
    """Configure logging for the application."""
    level = logging.DEBUG if verbose else logging.INFO
    logging.basicConfig(
        level=level,
        format='%(asctime)s - %(levelname)s - %(message)s',
        datefmt='%H:%M:%S'
    )


# ============================================================================
# Data Classes
# ============================================================================

@dataclass
class M365Endpoint:
    """Represents a Microsoft 365 endpoint (matches official API format)."""
    id: int
    service_area: str
    service_area_display_name: str = ""
    urls: List[str] = field(default_factory=list)
    ips: List[str] = field(default_factory=list)
    tcp_ports: str = ""
    udp_ports: str = ""
    express_route: bool = False
    category: str = ""
    required: bool = True
    notes: str = ""
    
    def to_dict(self) -> Dict:
        """Convert to dictionary matching official API format."""
        result = {'id': self.id, 'serviceArea': self.service_area}
        if self.service_area_display_name:
            result['serviceAreaDisplayName'] = self.service_area_display_name
        if self.urls:
            result['urls'] = self.urls
        if self.ips:
            result['ips'] = self.ips
        if self.tcp_ports:
            result['tcpPorts'] = self.tcp_ports
        if self.udp_ports:
            result['udpPorts'] = self.udp_ports
        # Always include expressRoute field
        result['expressRoute'] = self.express_route
        if self.category:
            result['category'] = self.category
        result['required'] = self.required
        if self.notes:
            result['notes'] = self.notes
        return result


# ============================================================================
# Network Fetching
# ============================================================================

def get_url(url: str) -> str:
    """Get content from a URL."""
    logging.debug("Fetching: %s", url)
    
    headers = {
        'User-Agent': 'Mozilla/5.0 (compatible; TerraformProvider/1.0)'
    }
    
    req = Request(url, headers=headers)
    
    try:
        with urlopen(req, timeout=30) as response:
            content = response.read().decode('utf-8')
            logging.debug("  Fetched %d bytes", len(content))
            return content
    except (HTTPError, URLError) as e:
        raise RuntimeError(f"Failed to fetch {url}: {e}") from e


def get_m365_endpoints(cloud: str, client_request_id: str) -> List[Dict]:
    """Get official M365 endpoints from API with all service areas."""
    # Query with all service areas to get maximum coverage
    service_areas = "MEM,Exchange,Skype,SharePoint,Common"
    url = f"{M365_CLOUDS[cloud]}?ServiceAreas={service_areas}&ClientRequestId={client_request_id}"
    logging.info("Fetching official %s endpoints from API (all service areas)...", cloud)
    
    content = get_url(url)
    data = json.loads(content)
    
    # Count by service area
    service_counts = {}
    for endpoint in data:
        sa = endpoint.get('serviceArea', 'N/A')
        service_counts[sa] = service_counts.get(sa, 0) + 1
    
    counts_str = ', '.join(f"{k}:{v}" for k, v in sorted(service_counts.items()))
    logging.info("  Retrieved %d official endpoints (%s)", len(data), counts_str)
    return data


# ============================================================================
# Text Sanitization
# ============================================================================

def sanitize_text(text: str) -> str:
    """Remove HTML tags and clean up text."""
    if not text:
        return text
    
    # Remove HTML tags (case-insensitive)
    text = re.sub(r'<br\s*/?>', ' ', text, flags=re.IGNORECASE)
    text = re.sub(r'<[^>]+>', '', text)
    
    # Clean up whitespace
    text = ' '.join(text.split())
    
    return text.strip()


# ============================================================================
# Cloud Data Manager
# ============================================================================

class CloudEndpointManager:
    """Manages endpoints for a specific cloud."""
    
    def __init__(self, cloud_name: str):
        self.cloud_name = cloud_name
        self.endpoints: List[Dict] = []
        self.max_id = 0
    
    def load_official_endpoints(self, client_request_id: str):
        """Load official M365 endpoints from API."""
        try:
            self.endpoints = get_m365_endpoints(self.cloud_name, client_request_id)
            self.max_id = max((e.get('id', 0) for e in self.endpoints), default=0)
            logging.info("[%s] Loaded %d official endpoints (max ID: %d)",
                        self.cloud_name.upper(), len(self.endpoints), self.max_id)
        except RuntimeError as e:
            logging.error("[%s] Failed to load official endpoints: %s",
                        self.cloud_name.upper(), e)
            raise
    
    def export_json(self, output_path: Path):
        """Export endpoints to JSON file."""
        # Sort by ID ascending
        sorted_endpoints = sorted(self.endpoints, key=lambda x: x.get('id', 0))
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(sorted_endpoints, f, indent=2)
        logging.info("[%s] Exported JSON: %s", self.cloud_name.upper(), output_path)
    
    def export_csv(self, output_path: Path):
        """Export endpoints to CSV file."""
        lines = ["ID,ServiceArea,URLs,IPs,TCPPorts,UDPPorts,ExpressRoute,Category,Required,Notes"]
        
        # Sort by ID ascending
        sorted_endpoints = sorted(self.endpoints, key=lambda x: x.get('id', 0))
        
        for endpoint in sorted_endpoints:
            row = [
                str(endpoint.get('id', '')),
                endpoint.get('serviceArea', ''),
                '; '.join(endpoint.get('urls', [])),
                '; '.join(endpoint.get('ips', [])),
                endpoint.get('tcpPorts', ''),
                endpoint.get('udpPorts', ''),
                'Yes' if endpoint.get('expressRoute', False) else 'No',
                endpoint.get('category', ''),
                'Yes' if endpoint.get('required', True) else 'No',
                endpoint.get('notes', '').replace(',', ';')
            ]
            lines.append(','.join(f'"{cell}"' if ',' in cell or ';' in cell else cell for cell in row))
        
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write('\n'.join(lines))
        logging.info("[%s] Exported CSV: %s", self.cloud_name.upper(), output_path)
    
    def export_go(self, output_path: Path):
        """Export endpoints to Go constants file."""
        # Sort by ID ascending
        sorted_endpoints = sorted(self.endpoints, key=lambda x: x.get('id', 0))
        
        lines = [
            "// Code generated by fetch_microsoft365_endpoints.py; DO NOT EDIT.",
            f"// Cloud: {self.cloud_name}",
            f"// Total endpoints: {len(sorted_endpoints)}",
            "",
            "package constants",
            "",
            "// M365Endpoint represents a Microsoft 365 network endpoint",
            "type M365Endpoint struct {",
            "\tID                  int      `json:\"id\"`",
            "\tServiceArea         string   `json:\"serviceArea\"`",
            "\tServiceAreaDisplay  string   `json:\"serviceAreaDisplayName,omitempty\"`",
            "\tURLs                []string `json:\"urls,omitempty\"`",
            "\tIPs                 []string `json:\"ips,omitempty\"`",
            "\tTCPPorts            string   `json:\"tcpPorts,omitempty\"`",
            "\tUDPPorts            string   `json:\"udpPorts,omitempty\"`",
            "\tExpressRoute        bool     `json:\"expressRoute\"`",
            "\tCategory            string   `json:\"category,omitempty\"`",
            "\tRequired            bool     `json:\"required\"`",
            "\tNotes               string   `json:\"notes,omitempty\"`",
            "}",
            "",
            f"// M365Endpoints{self.cloud_name.replace('-', '').title()} contains all endpoints for {self.cloud_name}",
            f"var M365Endpoints{self.cloud_name.replace('-', '').title()} = []M365Endpoint{{",
        ]
        
        for endpoint in sorted_endpoints:
            lines.append("\t{")
            lines.append(f"\t\tID: {endpoint.get('id', 0)},")
            lines.append(f"\t\tServiceArea: \"{endpoint.get('serviceArea', '')}\",")
            if endpoint.get('serviceAreaDisplayName'):
                lines.append(f"\t\tServiceAreaDisplay: \"{endpoint.get('serviceAreaDisplayName')}\",")
            if endpoint.get('urls'):
                urls_str = '", "'.join(endpoint['urls'])
                lines.append(f"\t\tURLs: []string{{\"{urls_str}\"}},")
            if endpoint.get('ips'):
                ips_str = '", "'.join(endpoint['ips'])
                lines.append(f"\t\tIPs: []string{{\"{ips_str}\"}},")
            if endpoint.get('tcpPorts'):
                lines.append(f"\t\tTCPPorts: \"{endpoint['tcpPorts']}\",")
            if endpoint.get('udpPorts'):
                lines.append(f"\t\tUDPPorts: \"{endpoint['udpPorts']}\",")
            lines.append(f"\t\tExpressRoute: {str(endpoint.get('expressRoute', False)).lower()},")
            if endpoint.get('category'):
                lines.append(f"\t\tCategory: \"{endpoint['category']}\",")
            lines.append(f"\t\tRequired: {str(endpoint.get('required', True)).lower()},")
            if endpoint.get('notes'):
                notes = endpoint['notes'].replace('"', '\\"')
                lines.append(f"\t\tNotes: \"{notes}\",")
            lines.append("\t},")
        
        lines.append("}")
        
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write('\n'.join(lines))
        logging.info("[%s] Exported Go: %s", self.cloud_name.upper(), output_path)


# ============================================================================
# Main Entry Point
# ============================================================================

def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description='Fetch and merge Microsoft 365 endpoint reference data',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument('--output-dir', '-o', default='output',
                        help='Output directory for generated files (default: output)')
    parser.add_argument('--format', '-f',
                        choices=['json', 'csv', 'go', 'all'],
                        default='all',
                        help='Output format(s) to generate (default: all)')
    parser.add_argument('--verbose', '-v', action='store_true',
                        help='Enable verbose/debug logging')
    
    args = parser.parse_args()
    
    setup_logging(args.verbose)
    
    output_dir = Path(args.output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)
    
    client_request_id = str(uuid.uuid4())
    logging.info("Client Request ID: %s", client_request_id)
    
    try:
        # Process each cloud
        for cloud_name in M365_CLOUDS.keys():
            logging.info("\n%s", "="*60)
            logging.info("Processing cloud: %s", cloud_name.upper())
            logging.info("%s", "="*60)
            
            manager = CloudEndpointManager(cloud_name)
            
            # Load official endpoints
            manager.load_official_endpoints(client_request_id)
            
            # Export to selected format(s)
            cloud_safe = cloud_name.replace('-', '_')
            
            if args.format in ['json', 'all']:
                manager.export_json(output_dir / f"microsoft365_endpoints_{cloud_safe}.json")
            
            if args.format in ['csv', 'all']:
                manager.export_csv(output_dir / f"microsoft365_endpoints_{cloud_safe}.csv")
            
            if args.format in ['go', 'all']:
                manager.export_go(output_dir / f"microsoft365_endpoints_{cloud_safe}.go")
        
        logging.info("\n%s", "="*60)
        logging.info("✓ All clouds processed successfully")
        logging.info("%s", "="*60)
        
    except KeyboardInterrupt:
        logging.warning("\nInterrupted by user")
        sys.exit(130)
    except RuntimeError as e:
        logging.error("Error: %s", e)
        sys.exit(1)


if __name__ == '__main__':
    main()

