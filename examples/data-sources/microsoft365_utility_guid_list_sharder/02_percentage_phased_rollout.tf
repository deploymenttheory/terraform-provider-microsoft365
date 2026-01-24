# Percentage-Based Distribution: Standard phased rollout pattern
# Use case: 10% pilot, 30% broader rollout, 60% full deployment

data "microsoft365_utility_guid_list_sharder" "ca_phases" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "ca-policies-2026" # Optional: ensures same users in same phases
}

# Phase 1: Pilot (10% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "block_legacy_auth_pilot" {
  display_name = "Block Legacy Auth - Phase 1 (10% Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
    client_app_types = ["exchangeActiveSync", "other"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}

# Phase 2: Broader Rollout (30% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "block_legacy_auth_broader" {
  display_name = "Block Legacy Auth - Phase 2 (30% Broader)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
    client_app_types = ["exchangeActiveSync", "other"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}

# Phase 3: Full Deployment (60% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "block_legacy_auth_full" {
  display_name = "Block Legacy Auth - Phase 3 (60% Full)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
    client_app_types = ["exchangeActiveSync", "other"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}

# Monitor phase distribution
output "phase_distribution" {
  value = {
    pilot_count   = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_0"])
    broader_count = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_1"])
    full_count    = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_2"])
    total_users   = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_2"])
  }
  description = "Phase counts (should be approximately 10%, 30%, 60% of total)"
}
