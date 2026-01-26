# Scenario 4: Service Principals (Enterprise Apps) → Conditional Access Policy
# Use case: Roll out stricter authentication requirements to enterprise applications in phases

# Distribute all service principals (enterprise apps) into 3 deployment rings
# Note: This queries service principals (/servicePrincipals), which are the app instances in your tenant
# These are the IDs used in Conditional Access policies
data "microsoft365_utility_guid_list_sharder" "app_rollout" {
  resource_type     = "service_principals"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "app-auth-policy-2026"
}

# Alternative: Filter to Microsoft apps only
data "microsoft365_utility_guid_list_sharder" "microsoft_apps" {
  resource_type     = "service_principals"
  odata_filter      = "startswith(displayName, 'Microsoft')"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "microsoft-app-policy-2026"
}

# Create CA policy targeting Ring 0 (10% pilot) applications
resource "microsoft365_graph_beta_conditional_access_policy" "app_auth_ring_0" {
  display_name = "Enhanced Auth Requirements - App Ring 0 (Pilot 10%)"
  state        = "enabledForReportingButNotEnforced" # Start in report-only mode

  conditions {
    users {
      include_users = ["All"]
    }
    applications {
      # Target the 10% pilot shard of applications
      include_applications = data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_0"]
    }
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]
  }

  grant_controls {
    operator = "AND"
    built_in_controls = [
      "mfa",
      "compliantDevice"
    ]
  }
}

# Create CA policy targeting Ring 1 (30% broader rollout)
resource "microsoft365_graph_beta_conditional_access_policy" "app_auth_ring_1" {
  display_name = "Enhanced Auth Requirements - App Ring 1 (Broader 30%)"
  state        = "disabled" # Enable after Ring 0 validation

  conditions {
    users {
      include_users = ["All"]
    }
    applications {
      # Target the 30% broader shard of applications
      include_applications = data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_1"]
    }
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]
  }

  grant_controls {
    operator = "AND"
    built_in_controls = [
      "mfa",
      "compliantDevice"
    ]
  }
}

# Create CA policy targeting Ring 2 (60% full rollout)
resource "microsoft365_graph_beta_conditional_access_policy" "app_auth_ring_2" {
  display_name = "Enhanced Auth Requirements - App Ring 2 (Full 60%)"
  state        = "disabled" # Enable after Ring 1 validation

  conditions {
    users {
      include_users = ["All"]
    }
    applications {
      # Target the 60% final shard of applications
      include_applications = data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_2"]
    }
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]
  }

  grant_controls {
    operator = "AND"
    built_in_controls = [
      "mfa",
      "compliantDevice"
    ]
  }
}

# Output application distribution for verification
output "app_distribution" {
  value = {
    ring_0_pilot_count   = length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_0"])
    ring_1_broader_count = length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_1"])
    ring_2_full_count    = length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_2"])
    total_apps           = length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_2"])
  }
}

# Output first few application IDs from pilot ring for manual verification
output "pilot_app_sample" {
  value       = slice(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_0"], 0, min(5, length(data.microsoft365_utility_guid_list_sharder.app_rollout.shards["shard_0"])))
  description = "Sample of application IDs in pilot ring (first 5)"
}
