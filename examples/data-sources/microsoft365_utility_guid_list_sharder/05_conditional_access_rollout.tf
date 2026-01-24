# Real-World Example: Conditional Access MFA Rollout with distributed pilot burden
# Demonstrates using different seeds for different initiatives

# MFA Rollout: User A ends up in 10% pilot
data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-rollout-2026" # Unique seed per initiative
}

# Windows Updates Rollout: Same User A ends up in 80% final wave
data "microsoft365_utility_guid_list_sharder" "windows_rollout" {
  resource_type     = "devices"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  seed              = "windows-updates-2026" # Different seed = different distribution
}

# MFA Phase 1: Pilot (10%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_phase_1" {
  display_name = "Require MFA - Phase 1 (10% Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# MFA Phase 2: Broader (30%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_phase_2" {
  display_name = "Require MFA - Phase 2 (30% Broader)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# MFA Phase 3: Full (60%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_phase_3" {
  display_name = "Require MFA - Phase 3 (60% Full)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Windows Updates Ring 0 (5% pilot devices)
resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_ring_0" {
  display_name = "Windows Updates - Ring 0 (5% Pilot)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  quality_update_deferral_period_days = 0
  feature_update_deferral_period_days = 0

  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_0"]
    }
  }
}

# Windows Updates Ring 1 (15% validation devices)
resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_ring_1" {
  display_name = "Windows Updates - Ring 1 (15% Validation)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  quality_update_deferral_period_days = 7
  feature_update_deferral_period_days = 14

  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_1"]
    }
  }
}

# Windows Updates Ring 2 (80% broad devices)
resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_ring_2" {
  display_name = "Windows Updates - Ring 2 (80% Broad)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  quality_update_deferral_period_days = 14
  feature_update_deferral_period_days = 30

  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_2"]
    }
  }
}

# Demonstrate distributed pilot burden
output "pilot_burden_distribution" {
  value = {
    mfa_pilot_users       = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"])
    mfa_broader_users     = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"])
    mfa_full_users        = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"])
    windows_pilot_devices = length(data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_0"])
    windows_broad_devices = length(data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_2"])
  }
  description = "Different seeds ensure User A isn't always in pilot groups across all initiatives"
}
