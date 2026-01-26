# Basic Service Principals: Query and shard enterprise applications
# Use for phased rollout of app-based conditional access policies

# All service principals with percentage-based sharding
data "microsoft365_utility_guid_list_sharder" "all_apps" {
  resource_type     = "service_principals"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "app-rollout-2026"
}

# Filter to Microsoft applications (common for targeting Microsoft 365 apps)
data "microsoft365_utility_guid_list_sharder" "microsoft_apps" {
  resource_type     = "service_principals"
  odata_filter      = "startswith(displayName, 'Microsoft')"
  shard_percentages = [20, 80]
  strategy          = "percentage"
  seed              = "microsoft-apps-2026"
}

# Filter to specific apps by AppId (useful for targeting known applications)
data "microsoft365_utility_guid_list_sharder" "office_apps" {
  resource_type = "service_principals"
  odata_filter  = "appId eq '00000003-0000-0000-c000-000000000000'"
  shard_count   = 3
  strategy      = "round-robin"
  seed          = "office-apps-2026"
}

# Filter to agentic service principals (AI Copilot agents)
# Use advanced OData type filter to target Microsoft Graph agentIdentityBlueprintPrincipal
data "microsoft365_utility_guid_list_sharder" "agentic_principals" {
  resource_type = "service_principals"
  odata_filter  = "isof('microsoft.graph.agentIdentityBlueprintPrincipal')"
  shard_count   = 2
  strategy      = "round-robin"
  seed          = "agentic-apps-2026"
}

# Output application distribution
output "all_apps_distribution" {
  value = {
    pilot_10pct   = length(data.microsoft365_utility_guid_list_sharder.all_apps.shards["shard_0"])
    broader_30pct = length(data.microsoft365_utility_guid_list_sharder.all_apps.shards["shard_1"])
    full_60pct    = length(data.microsoft365_utility_guid_list_sharder.all_apps.shards["shard_2"])
  }
}

output "microsoft_apps_distribution" {
  value = {
    pilot_20pct = length(data.microsoft365_utility_guid_list_sharder.microsoft_apps.shards["shard_0"])
    prod_80pct  = length(data.microsoft365_utility_guid_list_sharder.microsoft_apps.shards["shard_1"])
  }
}

output "agentic_principals_distribution" {
  value = {
    shard_0 = length(data.microsoft365_utility_guid_list_sharder.agentic_principals.shards["shard_0"])
    shard_1 = length(data.microsoft365_utility_guid_list_sharder.agentic_principals.shards["shard_1"])
    total   = length(data.microsoft365_utility_guid_list_sharder.agentic_principals.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.agentic_principals.shards["shard_1"])
  }
}
