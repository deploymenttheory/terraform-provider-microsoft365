# Rendezvous Hashing: Stable when ring count changes
# Use case: Start with 3 rings, later expand to 4 without massive user disruption

# Initial deployment: 3 rings
data "microsoft365_utility_guid_list_sharder" "stable_deployment" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 3 # Change to 4 later - only ~25% of users will move
  strategy      = "rendezvous"
  seed          = "stable-deployment-2026" # Required for rendezvous
}

# Ring 0: Early Adopters
resource "microsoft365_graph_beta_group" "ring_0_early" {
  display_name     = "Windows Updates - Ring 0 (Early Adopters)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_0"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_0" {
  display_name = "Windows Updates - Ring 0 (Early Adopters)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 0
  feature_update_deferral_period_days = 0

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.ring_0_early.id
    }
  }
}

# Ring 1: Broad Deployment
resource "microsoft365_graph_beta_group" "ring_1_broad" {
  display_name     = "Windows Updates - Ring 1 (Broad)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_1"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_1" {
  display_name = "Windows Updates - Ring 1 (Broad)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 7
  feature_update_deferral_period_days = 14

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.ring_1_broad.id
    }
  }
}

# Ring 2: Production (Conservative)
resource "microsoft365_graph_beta_group" "ring_2_production" {
  display_name     = "Windows Updates - Ring 2 (Production)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_2"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_2" {
  display_name = "Windows Updates - Ring 2 (Production)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 14
  feature_update_deferral_period_days = 30

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.ring_2_production.id
    }
  }
}

# Monitor ring distribution
output "ring_distribution" {
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_2"])
    total_users  = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_2"])
  }
  description = "When changing shard_count from 3 to 4, only ~25% of users move (vs ~75% with other strategies)"
}
