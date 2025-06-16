#!/usr/bin/env python3
"""Microsoft 365 Service Health Monitor.

A comprehensive monitoring tool for Microsoft 365 service health status,
incidents, and service degradations using Microsoft Graph APIs.

Usage:
    Basic health check:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client" --client-secret "your-secret"
    
    Include resolved issues and export:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client" --client-secret "your-secret" 
                                     --include-resolved-issues --export-path "health-report.json"
    
    Monitor specific issues:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client" --client-secret "your-secret" 
                                     --watch-issues "MP1085955" "SP1092170"
    
    Show only services with problems:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client" --client-secret "your-secret" 
                                     --problems-only

Requirements:
    - Azure AD application with ServiceHealth.Read.All and ServiceMessage.Read.All permissions
    - Admin consent granted for the application permissions
    - Python packages: requests, colorama, azure-identity
    
Exit Codes:
    0: All services healthy
    1: Issues detected (degraded services or active incidents)
    2: Authentication or permission errors
"""

import argparse
import json
import sys
import urllib.parse
from datetime import datetime, timedelta, UTC
from typing import Dict, List, Optional, Any

import requests
from colorama import init, Fore, Style
from azure.identity import ClientSecretCredential
from azure.core.exceptions import ClientAuthenticationError

# Initialize colorama for cross-platform color support
init(autoreset=True)

# Service status color mapping - Based on actual API responses
STATUS_COLORS = {
    "ServiceOperational": Fore.GREEN,
    "Investigating": Fore.YELLOW,
    "RestoringService": Fore.CYAN,
    "VerifyingService": Fore.CYAN,
    "ServiceRestored": Fore.GREEN,
    "PostIncidentReviewPublished": Fore.GREEN,
    "ServiceDegradation": Fore.YELLOW,
    "ServiceInterruption": Fore.RED,
    "ExtendedRecovery": Fore.YELLOW,
    "FalsePositive": Fore.WHITE,
    "InvestigationSuspended": Fore.MAGENTA,
    "Resolved": Fore.GREEN,
    "MitigatedExternal": Fore.GREEN,
    "Mitigated": Fore.GREEN,
    "ResolvedExternal": Fore.GREEN,
    "Confirmed": Fore.YELLOW,
    "Reported": Fore.YELLOW,
    # Issue classification types
    "Advisory": Fore.BLUE,
    "Incident": Fore.RED,
}

# Service ID to friendly name mapping - Based on actual API responses
SERVICE_NAMES = {
    "Exchange": "Exchange Online",
    "SharePoint": "SharePoint Online",
    "OneDriveForBusiness": "OneDrive for Business",
    "MicrosoftTeams": "Microsoft Teams",
    "PowerBICloud": "Power BI",
    "Dynamics365": "Dynamics 365",
    "OSDPPlatform": "Microsoft 365 suite",
    "OrgLiveID": "Identity Service",
    "Intune": "Microsoft Intune",
    "PowerAppsM365": "Power Apps",
    "PowerAutomateM365": "Power Automate",
    "WindowsVirtualDesktop": "Azure Virtual Desktop",
    "MicrosoftGraphConnectivity": "Microsoft Graph",
    "PowerVirtualAgents": "Power Virtual Agents",
    "Viva": "Microsoft Viva",
    "UniversalPrint": "Universal Print",
    "MicrosoftBookings": "Microsoft Bookings",
    "MicrosoftForms": "Microsoft Forms",
    "MicrosoftStream": "Microsoft Stream",
    "Yammer": "Yammer Enterprise",
    "MicrosoftDefenderforOffice365": "Microsoft Defender for Office 365",
}


class M365ServiceHealthMonitor:
    """Microsoft 365 Service Health Monitor.
    
    A comprehensive monitoring class that interfaces with Microsoft Graph APIs to retrieve
    and display Microsoft 365 service health information, including service status,
    active incidents, and service degradations. Uses Azure Identity SDK for robust
    authentication with automatic token refresh.
    
    Args:
        tenant_id (str): Azure AD tenant ID (Directory ID)
        client_id (str): Application (client) ID of the Azure AD app registration
        client_secret (str): Client secret of the Azure AD app registration
        include_resolved (bool, optional): Include recently resolved issues from last 7 days. Defaults to False.
        export_path (str, optional): Path to export detailed JSON report. Defaults to None.
        problems_only (bool, optional): Show only services with issues. Defaults to False.
        watch_issues (List[str], optional): Specific issue IDs to monitor. Defaults to None.
        required_permissions (List[str], optional): Microsoft Graph permissions to validate. 
            Defaults to ["ServiceHealth.Read.All", "ServiceMessage.Read.All"].
    
    Attributes:
        credential (ClientSecretCredential): Azure Identity credential for authentication
        headers (dict): HTTP headers with authorization for API requests
    
    Example:
        >>> monitor = M365ServiceHealthMonitor(
        ...     tenant_id="your-tenant-id",
        ...     client_id="your-client-id", 
        ...     client_secret="your-secret",
        ...     include_resolved=True,
        ...     export_path="health-report.json"
        ... )
        >>> monitor.get_access_token()
        >>> result = monitor.run()
        >>> print(f"Issues detected: {result['has_issues']}")
    
    Raises:
        ClientAuthenticationError: For Azure AD authentication failures
        requests.exceptions.RequestException: For API communication errors
        ValueError: For invalid parameter values
        PermissionError: For insufficient Microsoft Graph permissions
    """
    
    def __init__(self, tenant_id: str, client_id: str, client_secret: str, 
                 include_resolved: bool = False, export_path: Optional[str] = None,
                 problems_only: bool = False, watch_issues: Optional[List[str]] = None,
                 required_permissions: Optional[List[str]] = None):
        self.tenant_id = tenant_id
        self.client_id = client_id
        self.client_secret = client_secret
        self.include_resolved = include_resolved
        self.export_path = export_path
        self.problems_only = problems_only
        self.watch_issues = watch_issues or []
        self.required_permissions = required_permissions or ["ServiceHealth.Read.All", "ServiceMessage.Read.All"]
        self.credential = None
        self.headers = None
        
    def get_access_token(self) -> bool:
        """Authenticate with Microsoft Graph using Azure Identity SDK.
        
        Uses ClientSecretCredential from azure-identity for robust authentication
        with automatic token refresh and proper error handling. This is the recommended
        approach for service principal authentication with Microsoft Graph.
        
        Returns:
            bool: True if authentication successful, False otherwise
            
        Raises:
            ClientAuthenticationError: If authentication with Azure AD fails
            Exception: For other authentication-related errors
        """
        print(f"{Fore.CYAN}Connecting to Microsoft Graph...")
        
        try:
            # Create credential using Azure Identity SDK
            self.credential = ClientSecretCredential(
                tenant_id=self.tenant_id,
                client_id=self.client_id,
                client_secret=self.client_secret
            )
            
            # Get token for Microsoft Graph scope
            token = self.credential.get_token("https://graph.microsoft.com/.default")
            
            # Set up headers for requests
            self.headers = {
                'Authorization': f'Bearer {token.token}',
                'Content-Type': 'application/json'
            }
            
            print(f"{Fore.GREEN}✅ Connected to Microsoft Graph successfully")
            return True
            
        except ClientAuthenticationError as e:
            print(f"{Fore.RED}❌ Authentication failed: {e}")
            return False
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to connect to Microsoft Graph: {e}")
            return False
    
    def _refresh_token_if_needed(self):
        """Refresh access token if needed using Azure Identity SDK.
        
        The Azure Identity SDK handles token refresh automatically,
        but this method can be called to get a fresh token explicitly.
        """
        try:
            if self.credential:
                token = self.credential.get_token("https://graph.microsoft.com/.default")
                self.headers['Authorization'] = f'Bearer {token.token}'
        except Exception as e:
            print(f"{Fore.YELLOW}⚠️ Token refresh warning: {e}")
    
    def _make_graph_request(self, url: str, params: dict = None) -> dict:
        """Make a Microsoft Graph API request with automatic token refresh.
        
        Args:
            url (str): The Graph API endpoint URL
            params (dict, optional): Query parameters
            
        Returns:
            dict: JSON response from the API
            
        Raises:
            requests.exceptions.RequestException: If the request fails
        """
        try:
            response = requests.get(url, headers=self.headers, params=params)
            
            # If token expired, refresh and retry once
            if response.status_code == 401:
                self._refresh_token_if_needed()
                response = requests.get(url, headers=self.headers, params=params)
            
            response.raise_for_status()
            return response.json()
            
        except requests.exceptions.RequestException:
            # Try refreshing token once more for any request errors
            self._refresh_token_if_needed()
            response = requests.get(url, headers=self.headers, params=params)
            response.raise_for_status()
            return response.json()
    
    def test_service_principal_permissions(self) -> bool:
        """Validate service principal has required Microsoft Graph permissions.
        
        Performs a preflight check to ensure the authenticated service principal
        has all required application permissions and that admin consent has been granted.
        
        Returns:
            bool: True if all required permissions are granted, False otherwise
            
        Raises:
            requests.exceptions.RequestException: If permission validation API calls fail
        """
        try:
            print(f"\n{Fore.YELLOW}🔐 Validating service principal permissions...")
            
            # Get the service principal for this application
            sp_filter = f"appId eq '{self.client_id}'"
            encoded_filter = urllib.parse.quote(sp_filter)
            sp_url = f"https://graph.microsoft.com/v1.0/servicePrincipals?$filter={encoded_filter}"
            
            sp_data = self._make_graph_request(sp_url)
            
            if not sp_data.get('value'):
                print(f"{Fore.RED}❌ Service principal not found for application: {self.client_id}")
                return False
            
            service_principal = sp_data['value'][0]
            sp_id = service_principal['id']
            
            print(f"   {Fore.GREEN}✅ Found service principal: {service_principal['displayName']}")
            
            # Get Microsoft Graph service principal (resource)
            graph_sp_filter = "appId eq '00000003-0000-0000-c000-000000000000'"
            encoded_graph_filter = urllib.parse.quote(graph_sp_filter)
            graph_sp_url = f"https://graph.microsoft.com/v1.0/servicePrincipals?$filter={encoded_graph_filter}"
            
            graph_sp_data = self._make_graph_request(graph_sp_url)
            graph_service_principal = graph_sp_data['value'][0]
            
            # Get app role assignments for the service principal
            assignments_url = f"https://graph.microsoft.com/v1.0/servicePrincipals/{sp_id}/appRoleAssignments"
            assignments_data = self._make_graph_request(assignments_url)
            
            # Build a map of role names to role IDs for Microsoft Graph
            role_map = {}
            for role in graph_service_principal['appRoles']:
                role_map[role['value']] = role['id']
            
            print(f"\n   {Fore.CYAN}📋 Checking required permissions:")
            
            all_permissions_present = True
            missing_permissions = []
            
            for permission in self.required_permissions:
                role_id = role_map.get(permission)
                
                if not role_id:
                    print(f"   {Fore.RED}❌ {Fore.WHITE}{permission}{Fore.RED} - Unknown permission")
                    all_permissions_present = False
                    missing_permissions.append(permission)
                    continue
                
                # Check if this permission is assigned
                assignment = next((a for a in assignments_data['value'] 
                                 if a['resourceId'] == graph_service_principal['id'] and a['appRoleId'] == role_id), None)
                
                if assignment:
                    print(f"   {Fore.GREEN}✅ {Fore.WHITE}{permission}{Fore.GREEN} - Granted")
                else:
                    print(f"   {Fore.RED}❌ {Fore.WHITE}{permission}{Fore.RED} - Not granted")
                    all_permissions_present = False
                    missing_permissions.append(permission)
            
            if all_permissions_present:
                print(f"\n   {Fore.GREEN}✅ All required permissions are present")
                return True
            else:
                missing_perms_str = ', '.join(missing_permissions)
                print(f"\n   {Fore.RED}❌ Missing required permissions: {missing_perms_str}")
                return False
                
        except requests.exceptions.RequestException as e:
            print(f"{Fore.RED}❌ Failed to validate service principal permissions: {e}")
            return False
    
    def get_service_health_overview(self) -> List[Dict[str, Any]]:
        """Retrieve Microsoft 365 service health overview with pagination support.
        
        Fetches current status for all Microsoft 365 services including Exchange Online,
        SharePoint Online, Teams, etc. Handles pagination automatically to retrieve
        all available services.
        
        Returns:
            List[Dict[str, Any]]: List of service health objects containing service name and status
            
        Raises:
            requests.exceptions.RequestException: If API request fails
        """
        try:
            print(f"{Fore.YELLOW}📊 Retrieving service health overview...")
            
            all_services = []
            url = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/healthOverviews"
            
            while url:
                data = self._make_graph_request(url)
                
                if data.get('value'):
                    all_services.extend(data['value'])
                
                url = data.get('@odata.nextLink')
            
            print(f"{Fore.GREEN}✅ Retrieved {len(all_services)} services")
            return all_services
            
        except requests.exceptions.RequestException as e:
            print(f"{Fore.RED}❌ Failed to retrieve service health overview: {e}")
            raise
    
    def get_current_service_issues(self) -> List[Dict[str, Any]]:
        """Retrieve current Microsoft 365 service issues with filtering and pagination.
        
        Fetches service incidents and advisories from the last 7 days. Optionally filters
        out resolved issues based on the include_resolved setting. Handles pagination
        automatically to retrieve all matching issues.
        
        Returns:
            List[Dict[str, Any]]: List of service issue objects with details, status, and posts
            
        Raises:
            requests.exceptions.RequestException: If API request fails
        """
        try:
            print(f"{Fore.YELLOW}🚨 Retrieving current service issues...")
            
            all_issues = []
            
            # Get issues from last 7 days that are not resolved
            seven_days_ago = (datetime.now(UTC) - timedelta(days=7)).strftime("%Y-%m-%dT%H:%M:%SZ")
            filter_clause = f"lastModifiedDateTime ge {seven_days_ago}"
            
            if not self.include_resolved:
                resolved_statuses = [
                    "ServiceRestored", "FalsePositive", "Resolved", "Mitigated", 
                    "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished"
                ]
                status_filter = " and ".join([f"status ne '{status}'" for status in resolved_statuses])
                filter_clause += f" and {status_filter}"
            
            encoded_filter = urllib.parse.quote(filter_clause)
            url = f"https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/issues?$filter={encoded_filter}&$orderby=lastModifiedDateTime desc"
            
            while url:
                data = self._make_graph_request(url)
                
                if data.get('value'):
                    all_issues.extend(data['value'])
                
                url = data.get('@odata.nextLink')
            
            print(f"{Fore.GREEN}✅ Retrieved {len(all_issues)} service issues")
            return all_issues
            
        except requests.exceptions.RequestException as e:
            print(f"{Fore.RED}❌ Failed to retrieve service issues: {e}")
            raise
    
    def get_specific_service_issues(self, issue_ids: List[str]) -> List[Dict[str, Any]]:
        """Retrieve specific service issues by their IDs.
        
        Fetches detailed information for specific service issues being monitored.
        Useful for tracking particular incidents or advisories of interest.
        
        Args:
            issue_ids (List[str]): List of service issue IDs to retrieve
            
        Returns:
            List[Dict[str, Any]]: List of successfully retrieved service issue objects
            
        Raises:
            requests.exceptions.RequestException: If any API request fails
        """
        try:
            print(f"{Fore.YELLOW}🔍 Retrieving specific service issues...")
            
            specific_issues = []
            
            for issue_id in issue_ids:
                try:
                    url = f"https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/issues/{issue_id}"
                    issue = self._make_graph_request(url)
                    specific_issues.append(issue)
                    print(f"   {Fore.GREEN}✅ Retrieved issue: {issue_id}")
                except requests.exceptions.RequestException as e:
                    print(f"   {Fore.RED}❌ Failed to retrieve issue: {issue_id} - {e}")
            
            print(f"{Fore.GREEN}✅ Retrieved {len(specific_issues)} of {len(issue_ids)} requested issues")
            return specific_issues
            
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to retrieve specific service issues: {e}")
            raise
    
    def format_service_status(self, status: str) -> Dict[str, str]:
        """Format service status with appropriate colors and emoji icons.
        
        Maps Microsoft Graph API status values to user-friendly text, color codes,
        and emoji representations for enhanced visual display.
        
        Args:
            status (str): Raw status value from Microsoft Graph API
            
        Returns:
            Dict[str, str]: Dictionary containing 'emoji', 'text', and 'color' keys
            
        Example:
            >>> monitor.format_service_status("ServiceDegradation")
            {'emoji': '⚠️ ', 'text': 'Service Degradation', 'color': '\x1b[33m'}
        """
        color = STATUS_COLORS.get(status, Fore.WHITE)
        
        emoji_map = {
            "ServiceOperational": "✅",
            "Investigating": "🔍",
            "RestoringService": "🔄",
            "VerifyingService": "🔍",
            "ServiceRestored": "✅",
            "PostIncidentReviewPublished": "📋",
            "ServiceDegradation": "⚠️ ",
            "ServiceInterruption": "❌",
            "ExtendedRecovery": "🔄",
            "FalsePositive": "✅",
            "InvestigationSuspended": "⏸️ ",
            "Resolved": "✅",
            "Mitigated": "✅",
            "MitigatedExternal": "✅",
            "ResolvedExternal": "✅",
            "Confirmed": "⚠️ ",
            "Reported": "🔍",
            "Advisory": "💡",
            "Incident": "🚨",
            "serviceOperational": "✅",
            "serviceDegradation": "⚠️ ",
            "extendedRecovery": "🔄",
        }
        
        status_text_map = {
            "ServiceOperational": "Healthy",
            "Investigating": "Investigating",
            "RestoringService": "Restoring Service",
            "VerifyingService": "Verifying Service",
            "ServiceRestored": "Service Restored",
            "PostIncidentReviewPublished": "Post-Incident Review Published",
            "ServiceDegradation": "Service Degradation",
            "ServiceInterruption": "Service Interruption",
            "ExtendedRecovery": "Extended Recovery",
            "FalsePositive": "False Positive",
            "InvestigationSuspended": "Investigation Suspended",
            "Resolved": "Resolved",
            "Mitigated": "Mitigated",
            "MitigatedExternal": "Mitigated (External)",
            "ResolvedExternal": "Resolved (External)",
            "Confirmed": "Confirmed",
            "Reported": "Reported",
            "Advisory": "Advisory",
            "Incident": "Incident",
            "serviceOperational": "Healthy",
            "serviceDegradation": "Service Degradation",
            "extendedRecovery": "Extended Recovery",
        }
        
        emoji = emoji_map.get(status, "❓")
        text = status_text_map.get(status, status)
        
        return {
            "emoji": emoji,
            "text": text,
            "color": color
        }
    
    def get_friendly_service_name(self, service_key: str) -> str:
        """Convert service key to human-readable service name.
        
        Maps internal service identifiers to user-friendly display names.
        Falls back to converting CamelCase to spaced words for unknown services.
        
        Args:
            service_key (str): Internal service identifier from API
            
        Returns:
            str: Human-readable service name
            
        Example:
            >>> monitor.get_friendly_service_name("OneDriveForBusiness")
            'OneDrive for Business'
        """
        friendly_name = SERVICE_NAMES.get(service_key)
        if not friendly_name:
            # Convert CamelCase to spaced words
            import re
            friendly_name = re.sub(r'([a-z])([A-Z])', r'\1 \2', service_key)
        return friendly_name
    
    def show_service_health_dashboard(self, health_overview: List[Dict], issues: List[Dict], watched_issues: List[Dict] = None):
        """Display comprehensive service health dashboard with colored output.
        
        Renders a detailed console dashboard showing service status overview, 
        active incidents, watched issues, and optionally resolved issues.
        Uses colors and emoji for enhanced visual presentation.
        
        Args:
            health_overview (List[Dict]): Service health status data
            issues (List[Dict]): Service issues and incidents data  
            watched_issues (List[Dict], optional): Specific issues being monitored. Defaults to None.
        """
        if watched_issues is None:
            watched_issues = []
        
        print(f"\n{Fore.CYAN}🌐 Microsoft 365 Service Health Dashboard")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        current_time = datetime.now(UTC).strftime('%Y-%m-%d %H:%M:%S UTC')
        print(f"{Style.DIM}🕒 Last Updated: {current_time}")
        
        # Filter services for display
        services_to_show = health_overview
        if self.problems_only:
            services_to_show = [s for s in health_overview if s['status'] != "ServiceOperational"]
        
        # Service Status Overview
        if services_to_show:
            if self.problems_only:
                print(f"\n{Fore.RED}🚨 Services with Issues:")
            else:
                print(f"\n{Fore.CYAN}📊 Service Status Overview:")
            
            healthy_count = 0
            degraded_count = 0
            down_count = 0
            
            for service in sorted(services_to_show, key=lambda x: x['service']):
                service_name = service['service']
                status = self.format_service_status(service['status'])
                
                print(f"   {status['emoji']} {Fore.WHITE}{service_name}{Fore.RESET}: {status['color']}{status['text']}")
                
                # Count services based on actual API status values
                if service['status'].lower() in ["serviceoperational", "serviceoperational"]:
                    healthy_count += 1
                elif service['status'].lower() in ["servicedegradation", "investigating", "restoringservice", 
                                          "verifyingservice", "extendedrecovery", "confirmed", "reported", "servicedegradation"]:
                    degraded_count += 1
                elif service['status'].lower() == "serviceinterruption":
                    down_count += 1
            
            # Summary stats for all services (even if filtered view)
            if not self.problems_only:
                total_healthy = sum(1 for s in health_overview if s['status'].lower() in ["serviceoperational", "serviceoperational"])
                total_degraded = sum(1 for s in health_overview if s['status'].lower() in [
                    "servicedegradation", "investigating", "restoringservice", 
                    "verifyingservice", "extendedrecovery", "confirmed", "reported", "servicedegradation"
                ])
                total_down = sum(1 for s in health_overview if s['status'].lower() == "serviceinterruption")
                
                print(f"\n{Fore.CYAN}📈 Summary:")
                print(f"   • {Fore.GREEN}Healthy: {total_healthy}")
                print(f"   • {Fore.YELLOW}Degraded: {total_degraded}")
                print(f"   • {Fore.RED}Down: {total_down}")
        else:
            print(f"\n{Fore.GREEN}✅ All services are healthy")
        
        # Watched Issues
        if watched_issues:
            print(f"\n{Fore.CYAN}👁️  Watched Issues ({len(watched_issues)}):")
            
            for issue in sorted(watched_issues, key=lambda x: x['lastModifiedDateTime'], reverse=True):
                status = self.format_service_status(issue['status'])
                classification = self.format_service_status(issue['classification'])
                last_updated = datetime.fromisoformat(issue['lastModifiedDateTime'].replace('Z', '+00:00')).strftime("%m/%d %H:%M")
                
                # Determine if this is an active or resolved issue
                resolved_statuses = ["ServiceRestored", "FalsePositive", "Resolved", "Mitigated", 
                                  "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished"]
                is_active = issue['status'] not in resolved_statuses
                status_icon = "🔴" if is_active else "🟢"
                
                print(f"   {status_icon} {Fore.YELLOW}{issue['id']}{Fore.RESET}: {Fore.WHITE}{issue['title']}")
                print(f"     Status: {status['color']}{status['text']}{Fore.RESET} | Type: {classification['color']}{classification['text']}{Fore.RESET} | Updated: {Style.DIM}{last_updated}")
                
                # Show affected service
                if issue.get('service') and issue['service'].strip():
                    print(f"     Affected: {Style.DIM}{issue['service']}")
                
                # Show feature information if available
                if issue.get('featureGroup') and issue['featureGroup'].strip():
                    feature_info = issue['featureGroup']
                    if issue.get('feature') and issue['feature'].strip():
                        feature_info += f" - {issue['feature']}"
                    print(f"     Feature: {Style.DIM}{feature_info}")
                
                # Show latest post if available
                if issue.get('posts') and issue['posts']:
                    latest_post = max(issue['posts'], key=lambda x: x['createdDateTime'])
                    if latest_post.get('description', {}).get('content'):
                        import re
                        post_content = re.sub(r'<[^>]+>', '', latest_post['description']['content'])
                        post_content = post_content.replace('\n', ' ')
                        trimmed_post = post_content[:150] + "..." if len(post_content) > 150 else post_content
                        print(f"     Latest: {Style.DIM}{trimmed_post}")
        
        # Active Issues
        resolved_statuses = ["ServiceRestored", "FalsePositive", "Resolved", "Mitigated", 
                          "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished"]
        active_issues = [i for i in issues if i['status'] not in resolved_statuses]
        
        if active_issues:
            print(f"\n{Fore.RED}🚨 Active Incidents ({len(active_issues)}):")
            
            for issue in sorted(active_issues, key=lambda x: x['lastModifiedDateTime'], reverse=True):
                status = self.format_service_status(issue['status'])
                classification = self.format_service_status(issue['classification'])
                last_updated = datetime.fromisoformat(issue['lastModifiedDateTime'].replace('Z', '+00:00')).strftime("%m/%d %H:%M")
                
                print(f"   • {Fore.YELLOW}{issue['id']}{Fore.RESET}: {Fore.WHITE}{issue['title']}")
                print(f"     Status: {status['color']}{status['text']}{Fore.RESET} | Type: {classification['color']}{classification['text']}{Fore.RESET} | Updated: {Style.DIM}{last_updated}")
                
                # Show affected service
                if issue.get('service') and issue['service'].strip():
                    print(f"     Affected: {Style.DIM}{issue['service']}")
                
                # Show feature information if available
                if issue.get('featureGroup') and issue['featureGroup'].strip():
                    feature_info = issue['featureGroup']
                    if issue.get('feature') and issue['feature'].strip():
                        feature_info += f" - {issue['feature']}"
                    print(f"     Feature: {Style.DIM}{feature_info}")
                
                if issue.get('impactDescription') and issue['impactDescription'].strip():
                    trimmed_impact = issue['impactDescription'][:100] + "..." if len(issue['impactDescription']) > 100 else issue['impactDescription']
                    print(f"     Impact: {Style.DIM}{trimmed_impact}")
                
                # Show latest post if available
                if issue.get('posts') and issue['posts']:
                    latest_post = max(issue['posts'], key=lambda x: x['createdDateTime'])
                    if latest_post.get('description', {}).get('content'):
                        import re
                        post_content = re.sub(r'<[^>]+>', '', latest_post['description']['content'])
                        post_content = post_content.replace('\n', ' ')
                        trimmed_post = post_content[:120] + "..." if len(post_content) > 120 else post_content
                        print(f"     Latest: {Style.DIM}{trimmed_post}")
        else:
            print(f"\n{Fore.GREEN}✅ No Active Incidents")
        
        # Recently Resolved Issues
        if self.include_resolved:
            resolved_issues = [i for i in issues if i['status'] in resolved_statuses]
            
            if resolved_issues:
                print(f"\n{Fore.GREEN}🔧 Recently Resolved Issues ({len(resolved_issues)}):")
                
                for issue in sorted(resolved_issues, key=lambda x: x['lastModifiedDateTime'], reverse=True)[:5]:
                    resolved_time = datetime.fromisoformat(issue['lastModifiedDateTime'].replace('Z', '+00:00')).strftime("%m/%d %H:%M")
                    classification = self.format_service_status(issue['classification'])
                    
                    print(f"   • {Fore.YELLOW}{issue['id']}{Fore.RESET}: {Fore.WHITE}{issue['title']}")
                    print(f"     Type: {classification['color']}{classification['text']}{Fore.RESET} | Resolved: {Style.DIM}{resolved_time}")
    
    def export_detailed_report(self, health_overview: List[Dict], issues: List[Dict], watched_issues: List[Dict] = None):
        """Export comprehensive health data to JSON file.
        
        Creates a detailed JSON report containing service health overview, issues,
        watched issues, and summary statistics for external analysis or archival.
        
        Args:
            health_overview (List[Dict]): Service health status data
            issues (List[Dict]): Service issues and incidents data
            watched_issues (List[Dict], optional): Specific issues being monitored. Defaults to None.
            
        Raises:
            IOError: If file write operation fails
        """
        if watched_issues is None:
            watched_issues = []
        
        try:
            print(f"\n{Fore.YELLOW}📄 Exporting detailed report...")
            
            resolved_statuses = ["ServiceRestored", "FalsePositive", "Resolved", "Mitigated", 
                               "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished"]
            
            report = {
                "generatedAt": datetime.now(UTC).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "summary": {
                    "totalServices": len(health_overview),
                    "healthyServices": sum(1 for s in health_overview if s['status'] == "ServiceOperational"),
                    "degradedServices": sum(1 for s in health_overview if s['status'] in [
                        "ServiceDegradation", "Investigating", "RestoringService", 
                        "VerifyingService", "ExtendedRecovery", "Confirmed", "Reported"
                    ]),
                    "downServices": sum(1 for s in health_overview if s['status'] == "ServiceInterruption"),
                    "activeIncidents": sum(1 for i in issues if i['status'] not in resolved_statuses),
                    "resolvedIncidents": sum(1 for i in issues if i['status'] in resolved_statuses),
                    "watchedIssues": len(watched_issues),
                    "activeWatchedIssues": sum(1 for i in watched_issues if i['status'] not in resolved_statuses)
                },
                "serviceHealth": health_overview,
                "issues": issues,
                "watchedIssues": watched_issues
            }
            
            with open(self.export_path, 'w', encoding='utf-8') as f:
                json.dump(report, f, indent=2, ensure_ascii=False)
            
            print(f"{Fore.GREEN}✅ Report exported to: {self.export_path}")
            
        except Exception as e:
            print(f"{Fore.RED}❌ Failed to export report: {e}")
    
    def run(self) -> Dict[str, Any]:
        """Execute the complete service health monitoring workflow.
        
        Orchestrates the entire monitoring process: retrieves service health data,
        fetches current issues, processes watched issues, displays dashboard,
        and optionally exports detailed report.
        
        Returns:
            Dict[str, Any]: Summary dictionary containing:
                - has_issues (bool): Whether any issues were detected
                - active_incidents (int): Number of active service incidents
                - degraded_services (int): Number of services with degraded status
                - total_services (int): Total number of monitored services
                - healthy_services (int): Number of healthy services
                - watched_issues (int): Number of issues being watched
                - active_watched_issues (int): Number of active watched issues
                
        Raises:
            requests.exceptions.RequestException: If API requests fail
            Exception: For other processing errors
        """
        try:
            print(f"\n{Fore.CYAN}🔍 Starting Microsoft 365 Service Health Check...")
            print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
            
            # Get service health data
            health_overview = self.get_service_health_overview()
            issues = self.get_current_service_issues()
            
            # Get specific watched issues if requested
            watched_issues = []
            if self.watch_issues:
                watched_issues = self.get_specific_service_issues(self.watch_issues)
            
            # Display dashboard
            self.show_service_health_dashboard(health_overview, issues, watched_issues)
            
            # Export if requested
            if self.export_path:
                self.export_detailed_report(health_overview, issues, watched_issues)
            
            print(f"\n{Fore.GREEN}🎉 Service health check completed successfully!")
            
            # Return summary for programmatic use
            resolved_statuses = ["ServiceRestored", "FalsePositive", "Resolved", "Mitigated", 
                               "MitigatedExternal", "ResolvedExternal", "PostIncidentReviewPublished"]
            
            active_incidents = sum(1 for i in issues if i['status'] not in resolved_statuses)
            degraded_services = sum(1 for s in health_overview if s['status'] in [
                "ServiceDegradation", "ServiceInterruption", "Investigating", "RestoringService", 
                "VerifyingService", "ExtendedRecovery", "Confirmed", "Reported"
            ])
            active_watched_issues = sum(1 for i in watched_issues if i['status'] not in resolved_statuses)
            
            return {
                "has_issues": (active_incidents > 0 or degraded_services > 0 or active_watched_issues > 0),
                "active_incidents": active_incidents,
                "degraded_services": degraded_services,
                "total_services": len(health_overview),
                "healthy_services": sum(1 for s in health_overview if s['status'] == "ServiceOperational"),
                "watched_issues": len(watched_issues),
                "active_watched_issues": active_watched_issues
            }
            
        except Exception as e:
            print(f"\n{Fore.RED}❌ Service health check failed: {e}")
            raise

    def get_entra_id_health(self) -> Dict[str, Any]:
        """Get Entra ID service health alerts"""
        try:
            url = "https://graph.microsoft.com/beta/reports/healthMonitoring/alerts"
            response = self._make_graph_request("GET", url)
            return response
        except Exception as e:
            print(f"{Fore.RED}❌ Error getting Entra ID health: {e}")
            return {"value": []}

    def format_entra_id_alert(self, alert: Dict[str, Any]) -> str:
        """Format Entra ID alert for display"""
        try:
            created_time = datetime.fromisoformat(alert['createdDateTime'].replace('Z', '+00:00'))
            time_ago = datetime.now(UTC) - created_time
            
            # Format time ago
            if time_ago.days > 0:
                time_str = f"{time_ago.days}d ago"
            elif time_ago.seconds >= 3600:
                time_str = f"{time_ago.seconds // 3600}h ago"
            elif time_ago.seconds >= 60:
                time_str = f"{time_ago.seconds // 60}m ago"
            else:
                time_str = f"{time_ago.seconds}s ago"
            
            # Get impact counts and limits
            user_impact = next((imp for imp in alert.get('enrichment', {}).get('impacts', []) 
                              if imp.get('@odata.type') == '#microsoft.graph.healthMonitoring.userImpactSummary'), {})
            app_impact = next((imp for imp in alert.get('enrichment', {}).get('impacts', []) 
                             if imp.get('@odata.type') == '#microsoft.graph.healthMonitoring.applicationImpactSummary'), {})
            
            user_count = user_impact.get('impactedCount', 0)
            app_count = app_impact.get('impactedCount', 0)
            user_limit_exceeded = user_impact.get('impactedCountLimitExceeded', False)
            app_limit_exceeded = app_impact.get('impactedCountLimitExceeded', False)
            
            # Format alert type for display
            alert_type = alert['alertType'].replace('_', ' ').title()
            
            # Get state emoji
            state_emoji = {
                'active': '🔴',
                'resolved': '✅',
                'investigating': '🔍',
                'mitigated': '🟡'
            }.get(alert.get('state', '').lower(), '❓')
            
            # Get enrichment state emoji
            enrichment_state = alert.get('enrichment', {}).get('state', '').lower()
            enrichment_emoji = {
                'enriched': '📊',
                'enriching': '🔄',
                'failed': '❌'
            }.get(enrichment_state, '')
            
            # Build the alert string
            alert_str = [
                f"{state_emoji} {alert_type} ({time_str})",
                f"   • Category: {alert.get('category', 'N/A')}",
                f"   • Scenario: {alert.get('scenario', 'N/A')}",
                f"   • State: {alert.get('state', 'N/A')}",
                f"   • Enrichment: {enrichment_emoji} {enrichment_state.title() if enrichment_state else 'N/A'}"
            ]
            
            # Add impact information
            if user_count > 0 or app_count > 0:
                alert_str.append(f"   • Impact:")
                if user_count > 0:
                    user_str = f"{user_count:,} users"
                    if user_limit_exceeded:
                        user_str += " (limit exceeded)"
                    alert_str.append(f"     - {user_str}")
                if app_count > 0:
                    app_str = f"{app_count:,} applications"
                    if app_limit_exceeded:
                        app_str += " (limit exceeded)"
                    alert_str.append(f"     - {app_str}")
            
            # Add supporting data links if available
            supporting_data = alert.get('enrichment', {}).get('supportingData', {})
            if supporting_data:
                alert_str.append(f"   • Supporting Data:")
                if 'signIns' in supporting_data:
                    alert_str.append(f"     - Sign-ins: {supporting_data['signIns']}")
                if 'audits' in supporting_data:
                    alert_str.append(f"     - Audits: {supporting_data['audits']}")
            
            # Add signals if available
            signals = alert.get('signals', {})
            if signals:
                alert_str.append(f"   • Signals:")
                for signal_type, signal_url in signals.items():
                    alert_str.append(f"     - {signal_type}: {signal_url}")
            
            # Add documentation if available
            docs = alert.get('documentation', {})
            if docs and 'troubleshootingGuide' in docs:
                alert_str.append(f"   • Troubleshooting Guide: {docs['troubleshootingGuide']}")
            
            # Add creation time
            alert_str.append(f"   • Created: {created_time.strftime('%Y-%m-%d %H:%M:%S')} UTC")
            
            return '\n'.join(alert_str)
                   
        except Exception as e:
            print(f"{Fore.RED}❌ Error formatting Entra ID alert: {e}")
            return str(alert)

    def get_service_health(self) -> Dict[str, Any]:
        """Get Microsoft 365 service health"""
        try:
            # Get service health overview
            url = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/healthOverviews"
            response = self._make_graph_request("GET", url)
            
            # Get current issues
            issues_url = "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/issues"
            issues_response = self._make_graph_request("GET", issues_url)
            
            # Get Entra ID health alerts
            entra_id_health = self.get_entra_id_health()
            
            return {
                "overview": response,
                "issues": issues_response,
                "entra_id_alerts": entra_id_health
            }
            
        except Exception as e:
            print(f"{Fore.RED}❌ Error getting service health: {e}")
            return {
                "overview": {"value": []},
                "issues": {"value": []},
                "entra_id_alerts": {"value": []}
            }

    def display_service_health(self, health_data: Dict[str, Any]):
        """Display Microsoft 365 service health"""
        try:
            # Display service health overview
            print(f"\n{Fore.CYAN}📊 Microsoft 365 Service Health Overview")
            print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
            
            services = health_data["overview"].get("value", [])
            if not services:
                print(f"{Fore.YELLOW}⚠️ No service health data available")
                return
            
            # Count services by status
            status_counts = {
                "healthy": 0,
                "degraded": 0,
                "down": 0
            }
            
            # Display each service
            for service in services:
                status = service.get("status", "").lower()
                status_counts[status] = status_counts.get(status, 0) + 1
                
                print(self.format_service_status(service))
            
            # Display summary
            print(f"\n{Fore.CYAN}📈 Summary:")
            print(f"   • Healthy Services: {status_counts['healthy']}")
            print(f"   • Degraded Services: {status_counts['degraded']}")
            print(f"   • Down Services: {status_counts['down']}")
            
            # Display current issues
            issues = health_data["issues"].get("value", [])
            if issues:
                print(f"\n{Fore.CYAN}⚠️ Current Issues")
                print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
                for issue in issues:
                    print(self.format_issue(issue))
            
            # Display Entra ID health alerts
            entra_id_alerts = health_data["entra_id_alerts"].get("value", [])
            if entra_id_alerts:
                print(f"\n{Fore.CYAN}🔐 Entra ID Health Alerts")
                print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
                for alert in entra_id_alerts:
                    print(self.format_entra_id_alert(alert))
            
        except Exception as e:
            print(f"{Fore.RED}❌ Error displaying service health: {e}")

    def format_issue(self, issue: Dict[str, Any]) -> str:
        """Format service issue for display"""
        try:
            # Get start time
            start_time = datetime.fromisoformat(issue['startDateTime'].replace('Z', '+00:00'))
            time_ago = datetime.now(UTC) - start_time
            
            # Format time ago
            if time_ago.days > 0:
                time_str = f"{time_ago.days}d ago"
            elif time_ago.seconds >= 3600:
                time_str = f"{time_ago.seconds // 3600}h ago"
            elif time_ago.seconds >= 60:
                time_str = f"{time_ago.seconds // 60}m ago"
            else:
                time_str = f"{time_ago.seconds}s ago"
            
            # Get status emoji
            status_emoji = {
                'serviceOperational': '✅',
                'investigating': '🔍',
                'serviceDegradation': '⚠️',
                'serviceInterruption': '🔴',
                'restoringService': '🔄',
                'extendedRecovery': '⏳',
                'falsePositive': '✅',
                'investigationSuspended': '⏸️',
                'resolved': '✅',
                'mitigated': '🟡',
                'postIncidentReviewPublished': '📋'
            }.get(issue.get('status', '').lower(), '❓')
            
            # Get classification emoji
            classification_emoji = {
                'advisory': 'ℹ️',
                'incident': '🚨',
                'message': '📢'
            }.get(issue.get('classification', '').lower(), '❓')
            
            return (f"{status_emoji} {issue['title']} ({time_str})\n"
                   f"   • Status: {issue.get('status', 'N/A')}\n"
                   f"   • Classification: {issue.get('classification', 'N/A')}\n"
                   f"   • Impact: {issue.get('impactDescription', 'N/A')}\n"
                   f"   • Start Time: {start_time.strftime('%Y-%m-%d %H:%M:%S')} UTC")
                   
        except Exception as e:
            print(f"{Fore.RED}❌ Error formatting issue: {e}")
            return str(issue)


def main():
    """Main entry point for the Microsoft 365 Service Health Monitor.
    
    Parses command line arguments, initializes the monitor, performs authentication
    and permission validation, executes the health check, and returns appropriate
    exit codes based on service status.
    
    Exit Codes:
        0: All Microsoft 365 services are healthy
        1: Issues detected (service degradations or active incidents)  
        2: Authentication, permission, or execution errors
    """
    parser = argparse.ArgumentParser(
        description="Microsoft 365 Service Health Monitor - Following exact patterns from PowerShell version",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument(
        "--tenant-id",
        required=True,
        help="Specify the Entra ID tenant ID (Directory ID) where the application is registered"
    )
    
    parser.add_argument(
        "--client-id", 
        required=True,
        help="Specify the application (client) ID of the Entra ID app registration"
    )
    
    parser.add_argument(
        "--client-secret",
        required=True,
        help="Specify the client secret of the Entra ID app registration"
    )
    
    parser.add_argument(
        "--include-resolved-issues",
        action="store_true",
        help="Include recently resolved issues from the last 7 days"
    )
    
    parser.add_argument(
        "--export-path",
        help="Export detailed report to JSON file"
    )
    
    parser.add_argument(
        "--problems-only",
        action="store_true", 
        help="Show only services with issues"
    )
    
    parser.add_argument(
        "--watch-issues",
        nargs="*",
        help="Specific issue IDs to monitor"
    )
    
    parser.add_argument(
        "--required-permissions",
        nargs="*",
        default=["ServiceHealth.Read.All", "ServiceMessage.Read.All"],
        help="Required Microsoft Graph application permissions to validate"
    )
    
    args = parser.parse_args()
    
    # Create monitor instance
    monitor = M365ServiceHealthMonitor(
        tenant_id=args.tenant_id,
        client_id=args.client_id,
        client_secret=args.client_secret,
        include_resolved=args.include_resolved_issues,
        export_path=args.export_path,
        problems_only=args.problems_only,
        watch_issues=args.watch_issues,
        required_permissions=args.required_permissions
    )
    
    try:
        # Authenticate
        if not monitor.get_access_token():
            sys.exit(2)
        
        # Validate permissions
        if not monitor.test_service_principal_permissions():
            print(f"\n{Fore.RED}❌ Cannot proceed due to missing permissions")
            sys.exit(2)
        
        # Run health check
        result = monitor.run()
        
        # Set exit code based on service health
        if result["has_issues"]:
            print(f"\n{Fore.YELLOW}⚠️ Issues detected in Microsoft 365 services")
            print(f"   • Active Incidents: {result['active_incidents']}")
            print(f"   • Degraded Services: {result['degraded_services']}")
            if result["watched_issues"] > 0:
                print(f"   • Watched Issues: {result['active_watched_issues']}/{result['watched_issues']} active")
            sys.exit(1)
        else:
            status_message = f"All Microsoft 365 services are healthy ({result['healthy_services']}/{result['total_services']})"
            if result["watched_issues"] > 0:
                status_message += f" | Watched Issues: {result['watched_issues']} (all resolved)"
            print(f"\n{Fore.GREEN}✅ {status_message}")
            sys.exit(0)
            
    except Exception as e:
        print(f"\n{Fore.RED}❌ Script execution failed: {e}")
        sys.exit(2)


if __name__ == "__main__":
    main()