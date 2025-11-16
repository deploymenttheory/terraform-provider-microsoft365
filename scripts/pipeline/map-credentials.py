#!/usr/bin/env python3
"""
Maps service-specific credentials to standard M365 env vars.
Usage: ./map-credentials.py <service>
"""

import os
import sys
from typing import Tuple, Optional


# Service to credential variable mapping
CREDENTIAL_MAP = {
    "applications": ("M365_CLIENT_ID_APPLICATIONS", "M365_CLIENT_SECRET_APPLICATIONS"),
    "backup_storage": ("M365_CLIENT_ID_BACKUP_STORAGE", "M365_CLIENT_SECRET_BACKUP_STORAGE"),
    "device_and_app_management": ("M365_CLIENT_ID_DEVICE_AND_APP_MGMT", "M365_CLIENT_SECRET_DEVICE_AND_APP_MGMT"),
    "device_management": ("M365_CLIENT_ID_DEVICE_MGMT", "M365_CLIENT_SECRET_DEVICE_MGMT"),
    "groups": ("M365_CLIENT_ID_GROUPS", "M365_CLIENT_SECRET_GROUPS"),
    "identity_and_access": ("M365_CLIENT_ID_IDENTITY_ACCESS", "M365_CLIENT_SECRET_IDENTITY_ACCESS"),
    "m365_admin": ("M365_CLIENT_ID_M365_ADMIN", "M365_CLIENT_SECRET_M365_ADMIN"),
    "multitenant_management": ("M365_CLIENT_ID_MULTITENANT_MGMT", "M365_CLIENT_SECRET_MULTITENANT_MGMT"),
    "users": ("M365_CLIENT_ID_USERS", "M365_CLIENT_SECRET_USERS"),
    "utility": ("M365_CLIENT_ID_UTILITY", "M365_CLIENT_SECRET_UTILITY"),
    "windows_365": ("M365_CLIENT_ID_WINDOWS_365", "M365_CLIENT_SECRET_WINDOWS_365"),
}


def get_credentials(service: str) -> Tuple[Optional[str], Optional[str]]:
    """Get credentials for a service from environment."""
    if service not in CREDENTIAL_MAP:
        print(f"Error: unknown service: {service}", file=sys.stderr)
        sys.exit(1)
    
    client_id_var, client_secret_var = CREDENTIAL_MAP[service]
    client_id = os.environ.get(client_id_var, "")
    client_secret = os.environ.get(client_secret_var, "")
    
    return client_id, client_secret


def export_to_github_env(key: str, value: str) -> None:
    """Export variable to GITHUB_ENV if available."""
    github_env = os.environ.get("GITHUB_ENV")
    if github_env:
        with open(github_env, 'a') as f:
            f.write(f"{key}={value}\n")


def map_credentials(service: str) -> None:
    """Map service-specific credentials to standard env vars."""
    client_id, client_secret = get_credentials(service)
    
    # Check if credentials are available
    if not client_id or not client_secret:
        print(f"⚠️  No credentials found for {service} - tests will be skipped")
        export_to_github_env("SKIP_TESTS", "true")
        return
    
    # Export to GITHUB_ENV if available (for GitHub Actions)
    github_env = os.environ.get("GITHUB_ENV")
    if github_env:
        export_to_github_env("M365_CLIENT_ID", client_id)
        export_to_github_env("M365_CLIENT_SECRET", client_secret)
        export_to_github_env("SKIP_TESTS", "false")
    
    print(f"✅ Credentials configured for {service}")


def main():
    if len(sys.argv) < 2:
        print("Usage: map-credentials.py <service>", file=sys.stderr)
        sys.exit(1)
    
    service = sys.argv[1]
    map_credentials(service)


if __name__ == "__main__":
    main()

