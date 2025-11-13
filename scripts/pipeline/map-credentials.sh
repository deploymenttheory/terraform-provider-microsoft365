#!/bin/bash
set -euo pipefail

# Maps service-specific credentials to standard M365 env vars
# Usage: ./map-credentials.sh <service>

SERVICE="${1:-}"

if [[ -z "$SERVICE" ]]; then
    echo "Usage: $0 <service>"
    exit 1
fi

map_credentials() {
    local service="$1"
    local client_id_var=""
    local client_secret_var=""

    case "$service" in
        applications)
            client_id_var="M365_CLIENT_ID_APPLICATIONS"
            client_secret_var="M365_CLIENT_SECRET_APPLICATIONS"
            ;;
        backup_storage)
            client_id_var="M365_CLIENT_ID_BACKUP_STORAGE"
            client_secret_var="M365_CLIENT_SECRET_BACKUP_STORAGE"
            ;;
        device_and_app_management)
            client_id_var="M365_CLIENT_ID_DEVICE_AND_APP_MGMT"
            client_secret_var="M365_CLIENT_SECRET_DEVICE_AND_APP_MGMT"
            ;;
        device_management)
            client_id_var="M365_CLIENT_ID_DEVICE_MGMT"
            client_secret_var="M365_CLIENT_SECRET_DEVICE_MGMT"
            ;;
        groups)
            client_id_var="M365_CLIENT_ID_GROUPS"
            client_secret_var="M365_CLIENT_SECRET_GROUPS"
            ;;
        identity_and_access)
            client_id_var="M365_CLIENT_ID_IDENTITY_ACCESS"
            client_secret_var="M365_CLIENT_SECRET_IDENTITY_ACCESS"
            ;;
        m365_admin)
            client_id_var="M365_CLIENT_ID_M365_ADMIN"
            client_secret_var="M365_CLIENT_SECRET_M365_ADMIN"
            ;;
        multitenant_management)
            client_id_var="M365_CLIENT_ID_MULTITENANT_MGMT"
            client_secret_var="M365_CLIENT_SECRET_MULTITENANT_MGMT"
            ;;
        users)
            client_id_var="M365_CLIENT_ID_USERS"
            client_secret_var="M365_CLIENT_SECRET_USERS"
            ;;
        utility)
            client_id_var="M365_CLIENT_ID_UTILITY"
            client_secret_var="M365_CLIENT_SECRET_UTILITY"
            ;;
        windows_365)
            client_id_var="M365_CLIENT_ID_WINDOWS_365"
            client_secret_var="M365_CLIENT_SECRET_WINDOWS_365"
            ;;
        *)
            echo "Error: unknown service: $service"
            exit 1
            ;;
    esac

    local client_id="${!client_id_var:-}"
    local client_secret="${!client_secret_var:-}"

    # Check if credentials are available
    if [[ -z "$client_id" || -z "$client_secret" ]]; then
        echo "⚠️  No credentials found for ${service} - tests will be skipped"
        # Export skip flag to GITHUB_ENV
        if [[ -n "${GITHUB_ENV:-}" ]]; then
            echo "SKIP_TESTS=true" >> "$GITHUB_ENV"
        fi
        exit 0
    fi

    # Export to GITHUB_ENV if available (for GitHub Actions)
    if [[ -n "${GITHUB_ENV:-}" ]]; then
        echo "M365_CLIENT_ID=${client_id}" >> "$GITHUB_ENV"
        echo "M365_CLIENT_SECRET=${client_secret}" >> "$GITHUB_ENV"
        echo "SKIP_TESTS=false" >> "$GITHUB_ENV"
    else
        # For local testing, just print
        echo "export M365_CLIENT_ID=${client_id}"
        echo "export M365_CLIENT_SECRET=${client_secret}"
    fi

    echo "✅ Credentials configured for ${service}"
}

map_credentials "$SERVICE"
