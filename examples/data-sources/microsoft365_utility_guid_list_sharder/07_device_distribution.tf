# Device Distribution: Distribute managed devices for Windows Updates rollout
# Use case: Controlled Windows Update deployment across device population

data "microsoft365_utility_guid_list_sharder" "windows_update_rings" {
  resource_type     = "devices"
  odata_query       = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages = [5, 15, 30, 50]
  strategy          = "percentage"
  seed              = "windows-updates-2026"
}

# Ring 0: Validation (5% of devices)
resource "microsoft365_graph_beta_group" "update_ring_0" {
  display_name     = "Windows Updates - Ring 0 (5% Validation)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_0_validation" {
  display_name = "Windows Updates - Ring 0 (5% Validation)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 0
  feature_update_deferral_period_days = 0
  deadline_for_quality_updates_days   = 7
  deadline_for_feature_updates_days   = 14

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_0.id
    }
  }
}

# Ring 1: Pilot (15% of devices)
resource "microsoft365_graph_beta_group" "update_ring_1" {
  display_name     = "Windows Updates - Ring 1 (15% Pilot)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_1_pilot" {
  display_name = "Windows Updates - Ring 1 (15% Pilot)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 3
  feature_update_deferral_period_days = 7
  deadline_for_quality_updates_days   = 10
  deadline_for_feature_updates_days   = 21

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_1.id
    }
  }
}

# Ring 2: Broad (30% of devices)
resource "microsoft365_graph_beta_group" "update_ring_2" {
  display_name     = "Windows Updates - Ring 2 (30% Broad)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_2_broad" {
  display_name = "Windows Updates - Ring 2 (30% Broad)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 7
  feature_update_deferral_period_days = 14
  deadline_for_quality_updates_days   = 14
  deadline_for_feature_updates_days   = 28

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_2.id
    }
  }
}

# Ring 3: Production (50% of devices)
resource "microsoft365_graph_beta_group" "update_ring_3" {
  display_name     = "Windows Updates - Ring 3 (50% Production)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_3"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_3_production" {
  display_name = "Windows Updates - Ring 3 (50% Production)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 14
  feature_update_deferral_period_days = 30
  deadline_for_quality_updates_days   = 21
  deadline_for_feature_updates_days   = 45

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_3.id
    }
  }
}

# Monitor device ring distribution
output "device_ring_distribution" {
  value = {
    ring_0_validation_count = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"])
    ring_1_pilot_count      = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"])
    ring_2_broad_count      = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"])
    ring_3_production_count = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_3"])
    total_windows_devices   = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"]) + length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_3"])
  }
  description = "Device counts per ring (5%, 15%, 30%, 50% distribution)"
}
