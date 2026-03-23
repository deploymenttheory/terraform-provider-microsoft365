# ==============================================================================
# Example 5: All Supported Resource Types
#
# The sharder can query five different Microsoft Graph collections. Each has
# distinct use cases and optional filtering behaviour.
#
# resource_type options:
#   "users"              → GET /users          — MFA, CA, phased policy rollouts
#   "devices"            → GET /devices        — Windows Update rings, compliance
#   "applications"       → GET /applications   — App registration distribution
#   "service_principals" → GET /servicePrincipals — App-based CA policies
#   "group_members"      → GET /groups/{id}/members — Split an existing group
#
# All types support odata_filter (except group_members which filters server-side
# via the group membership itself). group_members additionally requires group_id.
# ==============================================================================

# Users — active member accounts only (excludes guests and disabled accounts)
resource "microsoft365_utility_guid_list_sharder" "active_users" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4
  strategy                = "round-robin"
  seed                    = "users-ring-2026"
  recalculate_on_next_run = false
}

# Devices — Azure AD joined Windows devices only (Intune-managed fleet)
resource "microsoft365_utility_guid_list_sharder" "windows_devices" {
  resource_type           = "devices"
  odata_filter            = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages       = [5, 15, 30, 50]
  strategy                = "percentage"
  seed                    = "windows-rings-2026"
  recalculate_on_next_run = false
}

# Applications — all app registrations in the tenant
# Use this when distributing app registrations across policy rings is needed.
resource "microsoft365_utility_guid_list_sharder" "app_registrations" {
  resource_type           = "applications"
  shard_percentages       = [10, 90]
  strategy                = "percentage"
  seed                    = "app-reg-pilot-2026"
  recalculate_on_next_run = false
}

# Service Principals — enterprise app instances used in CA application conditions.
# Filter to Microsoft-published apps (common for targeting M365 workloads).
resource "microsoft365_utility_guid_list_sharder" "microsoft_enterprise_apps" {
  resource_type           = "service_principals"
  odata_filter            = "startswith(displayName, 'Microsoft')"
  shard_count             = 3
  strategy                = "round-robin"
  seed                    = "sp-mfa-2026"
  recalculate_on_next_run = false
}

# Service Principals — all enterprise apps; no filter gives the complete tenant list
resource "microsoft365_utility_guid_list_sharder" "all_enterprise_apps" {
  resource_type           = "service_principals"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  seed                    = "sp-all-2026"
  recalculate_on_next_run = false
}

# Group Members — split an existing group's membership into sub-rings.
# Useful when you have an established Entra ID group (e.g. "All IT Staff") and
# want to phase a new policy across that group without creating static sub-groups.
# group_id is required when resource_type = "group_members".
resource "microsoft365_utility_guid_list_sharder" "it_dept_rings" {
  resource_type           = "group_members"
  group_id                = "12345678-1234-1234-1234-123456789abc" # Replace with real group object ID
  shard_count             = 3
  strategy                = "round-robin"
  seed                    = "it-dept-2026"
  recalculate_on_next_run = false
}
