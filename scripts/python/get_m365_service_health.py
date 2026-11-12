#!/usr/bin/env python3
"""Microsoft 365 Service Health Monitor.

A comprehensive monitoring tool for Microsoft 365 service health status,
incidents, and service degradations using Microsoft Graph APIs.

Usage:
    Basic health check:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client"
        --client-secret "your-secret"

    Include resolved issues and export:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client"
        --client-secret "your-secret"
                                     --include-resolved-issues --export-path "health-report.json"

    Monitor specific issues:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client"
        --client-secret "your-secret"
                                     --watch-issues "MP1085955" "SP1092170"

    Show only services with problems:
        python m365_health_monitor.py --tenant-id "your-tenant" --client-id "your-client"
        --client-secret "your-secret"
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
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from typing import Any, Dict, List, Optional

import requests
from azure.core.exceptions import ClientAuthenticationError
from azure.identity import ClientSecretCredential
from colorama import Fore, init
from graph_authentication import ClientCredentials, GraphAuthentication, add_credential_arguments
from service_health_display import (
    HEALTH_OPERATION_ERRORS,
    SERVICE_NAMES,
    STATUS_COLORS,
    ServiceHealthDisplay,
)

__all__ = ["M365ServiceHealthMonitor", "main", "STATUS_COLORS", "SERVICE_NAMES"]

init(autoreset=True)


@dataclass
# Explicit data fields mirror the configuration or external metadata contract.
class HealthMonitorOptions:  # pylint: disable=too-many-instance-attributes
    """Credentials and selection options for a service health report."""

    tenant_id: str
    client_id: str
    client_secret: str
    include_resolved: bool = False
    export_path: Optional[str] = None
    problems_only: bool = False
    watch_issues: Optional[List[str]] = None
    required_permissions: Optional[List[str]] = None


class M365ServiceHealthMonitor(GraphAuthentication, ServiceHealthDisplay):
    """Retrieve health data and render it with the selected reporting options."""

    def __init__(
        self, tenant_id: str, client_id: str, client_secret: str, *options, **named_options
    ):
        self.options = HealthMonitorOptions(
            tenant_id, client_id, client_secret, *options, **named_options
        )
        self.options.watch_issues = self.options.watch_issues or []
        GraphAuthentication.__init__(self, ClientCredentials(tenant_id, client_id, client_secret))
        ServiceHealthDisplay.__init__(
            self, self.options.include_resolved, self.options.problems_only
        )
        self.required_permissions = self.options.required_permissions or [
            "ServiceHealth.Read.All",
            "ServiceMessage.Read.All",
        ]

    def _make_graph_request(self, url: str, params: dict = None) -> dict:
        """Retrieve Graph JSON with the shared bounded refresh/retry behavior."""
        return self._request_graph("GET", url, params=params)

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
            sp_filter = f"appId eq '{self.options.client_id}'"
            encoded_filter = urllib.parse.quote(sp_filter)
            sp_url = f"https://graph.microsoft.com/v1.0/servicePrincipals?$filter={encoded_filter}"

            sp_data = self._make_graph_request(sp_url)

            if not sp_data.get("value"):
                print(
                    f"{Fore.RED}❌ Service principal not found for application:"
                    f" {self.options.client_id}"
                )
                return False

            service_principal = sp_data["value"][0]
            sp_id = service_principal["id"]

            print(f"   {Fore.GREEN}✅ Found service principal: {service_principal['displayName']}")

            # Get Microsoft Graph service principal (resource)
            graph_sp_filter = "appId eq '00000003-0000-0000-c000-000000000000'"
            encoded_graph_filter = urllib.parse.quote(graph_sp_filter)
            graph_sp_url = (
                f"https://graph.microsoft.com/v1.0/servicePrincipals?$filter={encoded_graph_filter}"
            )

            graph_sp_data = self._make_graph_request(graph_sp_url)
            graph_service_principal = graph_sp_data["value"][0]

            # Get app role assignments for the service principal
            assignments_url = (
                f"https://graph.microsoft.com/v1.0/servicePrincipals/{sp_id}/appRoleAssignments"
            )
            assignments_data = self._make_graph_request(assignments_url)

            # Build a map of role names to role IDs for Microsoft Graph
            role_map = {}
            for role in graph_service_principal["appRoles"]:
                role_map[role["value"]] = role["id"]

            print(f"\n   {Fore.CYAN}📋 Checking required permissions:")

            all_permissions_present = True
            missing_permissions = []

            for permission in self.required_permissions:
                role_id = role_map.get(permission)

                if not role_id:
                    print(
                        f"   {Fore.RED}❌ {Fore.WHITE}{permission}{Fore.RED} - Unknown permission"
                    )
                    all_permissions_present = False
                    missing_permissions.append(permission)
                    continue

                # Check if this permission is assigned
                assignment = next(
                    (
                        a
                        for a in assignments_data["value"]
                        if a["resourceId"] == graph_service_principal["id"]
                        and a["appRoleId"] == role_id
                    ),
                    None,
                )

                if assignment:
                    print(f"   {Fore.GREEN}✅ {Fore.WHITE}{permission}{Fore.GREEN} - Granted")
                else:
                    print(f"   {Fore.RED}❌ {Fore.WHITE}{permission}{Fore.RED} - Not granted")
                    all_permissions_present = False
                    missing_permissions.append(permission)

            if all_permissions_present:
                print(f"\n   {Fore.GREEN}✅ All required permissions are present")
                return True
            missing_perms_str = ", ".join(missing_permissions)
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

                if data.get("value"):
                    all_services.extend(data["value"])

                url = data.get("@odata.nextLink")

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
            seven_days_ago = (datetime.now(timezone.utc) - timedelta(days=7)).strftime(
                "%Y-%m-%dT%H:%M:%SZ"
            )
            filter_clause = f"lastModifiedDateTime ge {seven_days_ago}"

            if not self.include_resolved:
                resolved_statuses = [
                    "ServiceRestored",
                    "FalsePositive",
                    "Resolved",
                    "Mitigated",
                    "MitigatedExternal",
                    "ResolvedExternal",
                    "PostIncidentReviewPublished",
                ]
                status_filter = " and ".join(
                    [f"status ne '{status}'" for status in resolved_statuses]
                )
                filter_clause += f" and {status_filter}"

            encoded_filter = urllib.parse.quote(filter_clause)
            url = (
                "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/issues?"
                + f"$filter={encoded_filter}&$orderby=lastModifiedDateTime"
            ) + " desc"

            while url:
                data = self._make_graph_request(url)

                if data.get("value"):
                    all_issues.extend(data["value"])

                url = data.get("@odata.nextLink")

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
                    url = (
                        "https://graph.microsoft.com/v1.0/admin/serviceAnnouncement/"
                        + f"issues/{issue_id}"
                    )
                    issue = self._make_graph_request(url)
                    specific_issues.append(issue)
                    print(f"   {Fore.GREEN}✅ Retrieved issue: {issue_id}")
                except requests.exceptions.RequestException as e:
                    print(f"   {Fore.RED}❌ Failed to retrieve issue: {issue_id} - {e}")

            print(
                f"{Fore.GREEN}✅ Retrieved {len(specific_issues)} of {len(issue_ids)} requested"
                " issues"
            )
            return specific_issues

        except HEALTH_OPERATION_ERRORS as e:
            print(f"{Fore.RED}❌ Failed to retrieve specific service issues: {e}")
            raise

    def export_detailed_report(
        self, health_overview: List[Dict], issues: List[Dict], watched_issues: List[Dict] = None
    ):
        """Export comprehensive health data to JSON file.

        Creates a detailed JSON report containing service health overview, issues,
        watched issues, and summary statistics for external analysis or archival.

        Args:
            health_overview (List[Dict]): Service health status data
            issues (List[Dict]): Service issues and incidents data
            watched_issues (List[Dict], optional): Specific issues being monitored. Defaults to
            None.

        Raises:
            IOError: If file write operation fails
        """
        if watched_issues is None:
            watched_issues = []

        try:
            print(f"\n{Fore.YELLOW}📄 Exporting detailed report...")

            resolved_statuses = [
                "ServiceRestored",
                "FalsePositive",
                "Resolved",
                "Mitigated",
                "MitigatedExternal",
                "ResolvedExternal",
                "PostIncidentReviewPublished",
            ]

            report = {
                "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "summary": {
                    "totalServices": len(health_overview),
                    "healthyServices": sum(
                        1 for s in health_overview if s["status"] == "ServiceOperational"
                    ),
                    "degradedServices": sum(
                        1
                        for s in health_overview
                        if s["status"]
                        in [
                            "ServiceDegradation",
                            "Investigating",
                            "RestoringService",
                            "VerifyingService",
                            "ExtendedRecovery",
                            "Confirmed",
                            "Reported",
                        ]
                    ),
                    "downServices": sum(
                        1 for s in health_overview if s["status"] == "ServiceInterruption"
                    ),
                    "activeIncidents": sum(
                        1 for i in issues if i["status"] not in resolved_statuses
                    ),
                    "resolvedIncidents": sum(1 for i in issues if i["status"] in resolved_statuses),
                    "watchedIssues": len(watched_issues),
                    "activeWatchedIssues": sum(
                        1 for i in watched_issues if i["status"] not in resolved_statuses
                    ),
                },
                "serviceHealth": health_overview,
                "issues": issues,
                "watchedIssues": watched_issues,
            }

            with open(self.options.export_path, "w", encoding="utf-8") as f:
                json.dump(report, f, indent=2, ensure_ascii=False)

            print(f"{Fore.GREEN}✅ Report exported to: {self.options.export_path}")

        except HEALTH_OPERATION_ERRORS as e:
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
            if self.options.watch_issues:
                watched_issues = self.get_specific_service_issues(self.options.watch_issues)

            # Display dashboard
            self.show_service_health_dashboard(health_overview, issues, watched_issues)

            # Export if requested
            if self.options.export_path:
                self.export_detailed_report(health_overview, issues, watched_issues)

            print(f"\n{Fore.GREEN}🎉 Service health check completed successfully!")

            # Return summary for programmatic use
            resolved_statuses = [
                "ServiceRestored",
                "FalsePositive",
                "Resolved",
                "Mitigated",
                "MitigatedExternal",
                "ResolvedExternal",
                "PostIncidentReviewPublished",
            ]

            active_incidents = sum(1 for i in issues if i["status"] not in resolved_statuses)
            degraded_services = sum(
                1
                for s in health_overview
                if s["status"]
                in [
                    "ServiceDegradation",
                    "ServiceInterruption",
                    "Investigating",
                    "RestoringService",
                    "VerifyingService",
                    "ExtendedRecovery",
                    "Confirmed",
                    "Reported",
                ]
            )
            active_watched_issues = sum(
                1 for i in watched_issues if i["status"] not in resolved_statuses
            )

            return {
                "has_issues": (
                    active_incidents > 0 or degraded_services > 0 or active_watched_issues > 0
                ),
                "active_incidents": active_incidents,
                "degraded_services": degraded_services,
                "total_services": len(health_overview),
                "healthy_services": sum(
                    1 for s in health_overview if s["status"] == "ServiceOperational"
                ),
                "watched_issues": len(watched_issues),
                "active_watched_issues": active_watched_issues,
            }

        except HEALTH_OPERATION_ERRORS as e:
            print(f"\n{Fore.RED}❌ Service health check failed: {e}")
            raise

    def get_entra_id_health(self) -> Dict[str, Any]:
        """Get Entra ID service health alerts"""
        try:
            url = "https://graph.microsoft.com/beta/reports/healthMonitoring/alerts"
            response = self._make_graph_request("GET", url)
            return response
        except HEALTH_OPERATION_ERRORS as e:
            print(f"{Fore.RED}❌ Error getting Entra ID health: {e}")
            return {"value": []}

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
                "entra_id_alerts": entra_id_health,
            }

        except HEALTH_OPERATION_ERRORS as e:
            print(f"{Fore.RED}❌ Error getting service health: {e}")
            return {
                "overview": {"value": []},
                "issues": {"value": []},
                "entra_id_alerts": {"value": []},
            }


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
        description=(
            "Microsoft 365 Service Health Monitor - Following exact patterns from PowerShell"
            " version"
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )

    add_credential_arguments(parser)

    parser.add_argument(
        "--include-resolved-issues",
        action="store_true",
        help="Include recently resolved issues from the last 7 days",
    )

    parser.add_argument("--export-path", help="Export detailed report to JSON file")

    parser.add_argument(
        "--problems-only", action="store_true", help="Show only services with issues"
    )

    parser.add_argument("--watch-issues", nargs="*", help="Specific issue IDs to monitor")

    parser.add_argument(
        "--required-permissions",
        nargs="*",
        default=["ServiceHealth.Read.All", "ServiceMessage.Read.All"],
        help="Required Microsoft Graph application permissions to validate",
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
        required_permissions=args.required_permissions,
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
                print(
                    "   • Watched Issues:"
                    f" {result['active_watched_issues']}/{result['watched_issues']} active"
                )
            sys.exit(1)
        else:
            status_message = (
                "All Microsoft 365 services are healthy"
                f" ({result['healthy_services']}/{result['total_services']})"
            )
            if result["watched_issues"] > 0:
                status_message += f" | Watched Issues: {result['watched_issues']} (all resolved)"
            print(f"\n{Fore.GREEN}✅ {status_message}")
            sys.exit(0)

    except HEALTH_OPERATION_ERRORS as e:
        print(f"\n{Fore.RED}❌ Script execution failed: {e}")
        sys.exit(2)


if __name__ == "__main__":
    main()
