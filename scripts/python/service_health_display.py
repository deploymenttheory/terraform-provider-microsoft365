"""Format and display Microsoft 365 health responses."""

import re
from datetime import datetime, timezone
from typing import Any, Dict, List

from azure.core.exceptions import AzureError
from colorama import Fore, Style

HEALTH_OPERATION_ERRORS = (
    OSError,
    ValueError,
    TypeError,
    LookupError,
    AttributeError,
    RuntimeError,
    AzureError,
)

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


class ServiceHealthDisplay:
    """Render health responses independently of credentials and Graph retrieval."""

    def __init__(self, include_resolved: bool = False, problems_only: bool = False):
        self.include_resolved = include_resolved
        self.problems_only = problems_only

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

        return {"emoji": emoji, "text": text, "color": color}

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

            friendly_name = re.sub(r"([a-z])([A-Z])", r"\1 \2", service_key)
        return friendly_name

    def show_service_health_dashboard(
        self, health_overview: List[Dict], issues: List[Dict], watched_issues: List[Dict] = None
    ):
        """Display comprehensive service health dashboard with colored output.

        Renders a detailed console dashboard showing service status overview,
        active incidents, watched issues, and optionally resolved issues.
        Uses colors and emoji for enhanced visual presentation.

        Args:
            health_overview (List[Dict]): Service health status data
            issues (List[Dict]): Service issues and incidents data
            watched_issues (List[Dict], optional): Specific issues being monitored. Defaults to
            None.
        """
        if watched_issues is None:
            watched_issues = []

        print(f"\n{Fore.CYAN}🌐 Microsoft 365 Service Health Dashboard")
        print(f"{Fore.CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        current_time = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S UTC")
        print(f"{Style.DIM}🕒 Last Updated: {current_time}")

        # Filter services for display
        services_to_show = health_overview
        if self.problems_only:
            services_to_show = [s for s in health_overview if s["status"] != "ServiceOperational"]

        # Service Status Overview
        if services_to_show:
            if self.problems_only:
                print(f"\n{Fore.RED}🚨 Services with Issues:")
            else:
                print(f"\n{Fore.CYAN}📊 Service Status Overview:")

            healthy_count = 0
            degraded_count = 0
            down_count = 0

            for service in sorted(services_to_show, key=lambda x: x["service"]):
                service_name = service["service"]
                status = self.format_service_status(service["status"])

                print(
                    f"   {status['emoji']} {Fore.WHITE}{service_name}{Fore.RESET}:"
                    f" {status['color']}{status['text']}"
                )

                # Count services based on actual API status values
                if service["status"].lower() in ["serviceoperational", "serviceoperational"]:
                    healthy_count += 1
                elif service["status"].lower() in [
                    "servicedegradation",
                    "investigating",
                    "restoringservice",
                    "verifyingservice",
                    "extendedrecovery",
                    "confirmed",
                    "reported",
                    "servicedegradation",
                ]:
                    degraded_count += 1
                elif service["status"].lower() == "serviceinterruption":
                    down_count += 1

            # Summary stats for all services (even if filtered view)
            if not self.problems_only:
                total_healthy = sum(
                    1
                    for s in health_overview
                    if s["status"].lower() in ["serviceoperational", "serviceoperational"]
                )
                total_degraded = sum(
                    1
                    for s in health_overview
                    if s["status"].lower()
                    in [
                        "servicedegradation",
                        "investigating",
                        "restoringservice",
                        "verifyingservice",
                        "extendedrecovery",
                        "confirmed",
                        "reported",
                        "servicedegradation",
                    ]
                )
                total_down = sum(
                    1 for s in health_overview if s["status"].lower() == "serviceinterruption"
                )

                print(f"\n{Fore.CYAN}📈 Summary:")
                print(f"   • {Fore.GREEN}Healthy: {total_healthy}")
                print(f"   • {Fore.YELLOW}Degraded: {total_degraded}")
                print(f"   • {Fore.RED}Down: {total_down}")
        else:
            print(f"\n{Fore.GREEN}✅ All services are healthy")

        # Watched Issues
        if watched_issues:
            print(f"\n{Fore.CYAN}👁️  Watched Issues ({len(watched_issues)}):")

            for issue in sorted(
                watched_issues, key=lambda x: x["lastModifiedDateTime"], reverse=True
            ):
                status = self.format_service_status(issue["status"])
                classification = self.format_service_status(issue["classification"])
                last_updated = datetime.fromisoformat(
                    issue["lastModifiedDateTime"].replace("Z", "+00:00")
                ).strftime("%m/%d %H:%M")

                # Determine if this is an active or resolved issue
                resolved_statuses = [
                    "ServiceRestored",
                    "FalsePositive",
                    "Resolved",
                    "Mitigated",
                    "MitigatedExternal",
                    "ResolvedExternal",
                    "PostIncidentReviewPublished",
                ]
                is_active = issue["status"] not in resolved_statuses
                status_icon = "🔴" if is_active else "🟢"

                print(
                    f"   {status_icon} {Fore.YELLOW}{issue['id']}{Fore.RESET}:"
                    f" {Fore.WHITE}{issue['title']}"
                )
                print(
                    f"     Status: {status['color']}{status['text']}{Fore.RESET} | Type:"
                    f" {classification['color']}{classification['text']}{Fore.RESET} | Updated:"
                    f" {Style.DIM}{last_updated}"
                )

                # Show affected service
                if issue.get("service") and issue["service"].strip():
                    print(f"     Affected: {Style.DIM}{issue['service']}")

                # Show feature information if available
                if issue.get("featureGroup") and issue["featureGroup"].strip():
                    feature_info = issue["featureGroup"]
                    if issue.get("feature") and issue["feature"].strip():
                        feature_info += f" - {issue['feature']}"
                    print(f"     Feature: {Style.DIM}{feature_info}")

                # Show latest post if available
                if issue.get("posts") and issue["posts"]:
                    latest_post = max(issue["posts"], key=lambda x: x["createdDateTime"])
                    if latest_post.get("description", {}).get("content"):

                        post_content = re.sub(r"<[^>]+>", "", latest_post["description"]["content"])
                        post_content = post_content.replace("\n", " ")
                        trimmed_post = (
                            post_content[:150] + "..." if len(post_content) > 150 else post_content
                        )
                        print(f"     Latest: {Style.DIM}{trimmed_post}")

        # Active Issues
        resolved_statuses = [
            "ServiceRestored",
            "FalsePositive",
            "Resolved",
            "Mitigated",
            "MitigatedExternal",
            "ResolvedExternal",
            "PostIncidentReviewPublished",
        ]
        active_issues = [i for i in issues if i["status"] not in resolved_statuses]

        if active_issues:
            print(f"\n{Fore.RED}🚨 Active Incidents ({len(active_issues)}):")

            for issue in sorted(
                active_issues, key=lambda x: x["lastModifiedDateTime"], reverse=True
            ):
                status = self.format_service_status(issue["status"])
                classification = self.format_service_status(issue["classification"])
                last_updated = datetime.fromisoformat(
                    issue["lastModifiedDateTime"].replace("Z", "+00:00")
                ).strftime("%m/%d %H:%M")

                print(f"   • {Fore.YELLOW}{issue['id']}{Fore.RESET}: {Fore.WHITE}{issue['title']}")
                print(
                    f"     Status: {status['color']}{status['text']}{Fore.RESET} | Type:"
                    f" {classification['color']}{classification['text']}{Fore.RESET} | Updated:"
                    f" {Style.DIM}{last_updated}"
                )

                # Show affected service
                if issue.get("service") and issue["service"].strip():
                    print(f"     Affected: {Style.DIM}{issue['service']}")

                # Show feature information if available
                if issue.get("featureGroup") and issue["featureGroup"].strip():
                    feature_info = issue["featureGroup"]
                    if issue.get("feature") and issue["feature"].strip():
                        feature_info += f" - {issue['feature']}"
                    print(f"     Feature: {Style.DIM}{feature_info}")

                if issue.get("impactDescription") and issue["impactDescription"].strip():
                    trimmed_impact = (
                        issue["impactDescription"][:100] + "..."
                        if len(issue["impactDescription"]) > 100
                        else issue["impactDescription"]
                    )
                    print(f"     Impact: {Style.DIM}{trimmed_impact}")

                # Show latest post if available
                if issue.get("posts") and issue["posts"]:
                    latest_post = max(issue["posts"], key=lambda x: x["createdDateTime"])
                    if latest_post.get("description", {}).get("content"):

                        post_content = re.sub(r"<[^>]+>", "", latest_post["description"]["content"])
                        post_content = post_content.replace("\n", " ")
                        trimmed_post = (
                            post_content[:120] + "..." if len(post_content) > 120 else post_content
                        )
                        print(f"     Latest: {Style.DIM}{trimmed_post}")
        else:
            print(f"\n{Fore.GREEN}✅ No Active Incidents")

        # Recently Resolved Issues
        if self.include_resolved:
            resolved_issues = [i for i in issues if i["status"] in resolved_statuses]

            if resolved_issues:
                print(f"\n{Fore.GREEN}🔧 Recently Resolved Issues ({len(resolved_issues)}):")

                for issue in sorted(
                    resolved_issues, key=lambda x: x["lastModifiedDateTime"], reverse=True
                )[:5]:
                    resolved_time = datetime.fromisoformat(
                        issue["lastModifiedDateTime"].replace("Z", "+00:00")
                    ).strftime("%m/%d %H:%M")
                    classification = self.format_service_status(issue["classification"])

                    print(
                        f"   • {Fore.YELLOW}{issue['id']}{Fore.RESET}: {Fore.WHITE}{issue['title']}"
                    )
                    print(
                        f"     Type: {classification['color']}{classification['text']}{Fore.RESET}"
                        f" | Resolved: {Style.DIM}{resolved_time}"
                    )

    def format_entra_id_alert(self, alert: Dict[str, Any]) -> str:
        """Format Entra ID alert for display"""
        try:
            created_time = datetime.fromisoformat(alert["createdDateTime"].replace("Z", "+00:00"))
            time_ago = datetime.now(timezone.utc) - created_time

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
            user_impact = next(
                (
                    imp
                    for imp in alert.get("enrichment", {}).get("impacts", [])
                    if imp.get("@odata.type")
                    == "#microsoft.graph.healthMonitoring.userImpactSummary"
                ),
                {},
            )
            app_impact = next(
                (
                    imp
                    for imp in alert.get("enrichment", {}).get("impacts", [])
                    if imp.get("@odata.type")
                    == "#microsoft.graph.healthMonitoring.applicationImpactSummary"
                ),
                {},
            )

            user_count = user_impact.get("impactedCount", 0)
            app_count = app_impact.get("impactedCount", 0)
            user_limit_exceeded = user_impact.get("impactedCountLimitExceeded", False)
            app_limit_exceeded = app_impact.get("impactedCountLimitExceeded", False)

            # Format alert type for display
            alert_type = alert["alertType"].replace("_", " ").title()

            # Get state emoji
            state_emoji = {
                "active": "🔴",
                "resolved": "✅",
                "investigating": "🔍",
                "mitigated": "🟡",
            }.get(alert.get("state", "").lower(), "❓")

            # Get enrichment state emoji
            enrichment_state = alert.get("enrichment", {}).get("state", "").lower()
            enrichment_emoji = {"enriched": "📊", "enriching": "🔄", "failed": "❌"}.get(
                enrichment_state, ""
            )

            # Build the alert string
            alert_str = [
                f"{state_emoji} {alert_type} ({time_str})",
                f"   • Category: {alert.get('category', 'N/A')}",
                f"   • Scenario: {alert.get('scenario', 'N/A')}",
                f"   • State: {alert.get('state', 'N/A')}",
                (
                    f"   • Enrichment: {enrichment_emoji}"
                    f" {enrichment_state.title() if enrichment_state else 'N/A'}"
                ),
            ]

            # Add impact information
            if user_count > 0 or app_count > 0:
                alert_str.append("   • Impact:")
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
            supporting_data = alert.get("enrichment", {}).get("supportingData", {})
            if supporting_data:
                alert_str.append("   • Supporting Data:")
                if "signIns" in supporting_data:
                    alert_str.append(f"     - Sign-ins: {supporting_data['signIns']}")
                if "audits" in supporting_data:
                    alert_str.append(f"     - Audits: {supporting_data['audits']}")

            # Add signals if available
            signals = alert.get("signals", {})
            if signals:
                alert_str.append("   • Signals:")
                for signal_type, signal_url in signals.items():
                    alert_str.append(f"     - {signal_type}: {signal_url}")

            # Add documentation if available
            docs = alert.get("documentation", {})
            if docs and "troubleshootingGuide" in docs:
                alert_str.append(f"   • Troubleshooting Guide: {docs['troubleshootingGuide']}")

            # Add creation time
            alert_str.append(f"   • Created: {created_time.strftime('%Y-%m-%d %H:%M:%S')} UTC")

            return "\n".join(alert_str)

        except HEALTH_OPERATION_ERRORS as e:
            print(f"{Fore.RED}❌ Error formatting Entra ID alert: {e}")
            return str(alert)

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
            status_counts = {"healthy": 0, "degraded": 0, "down": 0}

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

        except HEALTH_OPERATION_ERRORS as e:
            print(f"{Fore.RED}❌ Error displaying service health: {e}")

    def format_issue(self, issue: Dict[str, Any]) -> str:
        """Format service issue for display"""
        try:
            # Get start time
            start_time = datetime.fromisoformat(issue["startDateTime"].replace("Z", "+00:00"))
            time_ago = datetime.now(timezone.utc) - start_time

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
                "serviceOperational": "✅",
                "investigating": "🔍",
                "serviceDegradation": "⚠️",
                "serviceInterruption": "🔴",
                "restoringService": "🔄",
                "extendedRecovery": "⏳",
                "falsePositive": "✅",
                "investigationSuspended": "⏸️",
                "resolved": "✅",
                "mitigated": "🟡",
                "postIncidentReviewPublished": "📋",
            }.get(issue.get("status", "").lower(), "❓")

            # Get classification emoji

            return (
                f"{status_emoji} {issue['title']} ({time_str})\n"
                f"   • Status: {issue.get('status', 'N/A')}\n"
                f"   • Classification: {issue.get('classification', 'N/A')}\n"
                f"   • Impact: {issue.get('impactDescription', 'N/A')}\n"
                f"   • Start Time: {start_time.strftime('%Y-%m-%d %H:%M:%S')} UTC"
            )

        except HEALTH_OPERATION_ERRORS as e:
            print(f"{Fore.RED}❌ Error formatting issue: {e}")
            return str(issue)
