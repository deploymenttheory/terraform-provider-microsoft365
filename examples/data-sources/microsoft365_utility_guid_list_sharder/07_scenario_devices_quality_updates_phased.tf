# Scenario 3: Devices → Groups → Deployment Scheduler → Windows Quality Update Policy
# Use case: Phased rollout of Windows quality updates with automated timing gates

# Step 1: Shard Windows devices into 3 deployment rings (10%, 30%, 60%)
data "microsoft365_utility_guid_list_sharder" "quality_update_rings" {
  resource_type     = "devices"
  odata_query       = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "quality-updates-2024"
}

# Step 2: Create Entra ID groups for each deployment ring
resource "microsoft365_graph_beta_group" "ring_0_pilot" {
  display_name     = "Quality Updates - Ring 0 (10% Pilot)"
  mail_nickname    = "quality-updates-ring-0"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "ring_1_broad" {
  display_name     = "Quality Updates - Ring 1 (30% Broad)"
  mail_nickname    = "quality-updates-ring-1"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "ring_2_production" {
  display_name     = "Quality Updates - Ring 2 (60% Production)"
  mail_nickname    = "quality-updates-ring-2"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_2"]
}

# Step 3: Define deployment timing gates for phased rollout
locals {
  deployment_start = "2026-01-20T00:00:00Z"
}

# Phase 1: Pilot ring opens after 24h
data "microsoft365_utility_deployment_scheduler" "ring_0_gate" {
  name                  = "quality-updates-ring-0-pilot"
  deployment_start_time = local.deployment_start
  scope_id              = microsoft365_graph_beta_group.ring_0_pilot.id

  time_condition = {
    delay_start_time_by = 24 # Open after 24 hours
  }

  # Only deploy during business hours
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "17:00:00"
      }
    ]
  }
}

# Phase 2: Broad ring opens 72h after pilot opens
data "microsoft365_utility_deployment_scheduler" "ring_1_gate" {
  name                  = "quality-updates-ring-1-broad"
  deployment_start_time = local.deployment_start
  scope_id              = microsoft365_graph_beta_group.ring_1_broad.id

  time_condition = {
    delay_start_time_by = 24
  }

  # Wait for pilot to be open for 72 hours
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 24
    minimum_open_hours               = 72
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "08:00:00"
        time_of_day_end   = "18:00:00"
      }
    ]
  }
}

# Phase 3: Production ring opens 1 week after broad opens
data "microsoft365_utility_deployment_scheduler" "ring_2_gate" {
  name                  = "quality-updates-ring-2-production"
  deployment_start_time = local.deployment_start
  scope_id              = microsoft365_graph_beta_group.ring_2_production.id

  time_condition = {
    delay_start_time_by = 192 # 24 + 72 + 96 = 192 hours total
  }

  # Wait for broad ring to be open for 1 week
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 96  # 24 + 72
    minimum_open_hours               = 168 # 1 week
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
        time_of_day_start = "08:00:00"
        time_of_day_end   = "18:00:00"
      }
    ]
  }

  # Avoid Friday deployments to production
  exclusion_time_windows = {
    window = [
      {
        days_of_week = ["friday"]
      }
    ]
  }
}

# Step 4: Create Windows Quality Update Policy with conditional assignments
resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "phased_quality_updates" {
  display_name     = "Windows Quality Updates - Phased Rollout"
  description      = "Monthly quality updates deployed in phases with automated timing"
  hotpatch_enabled = true

  # Conditional assignments based on deployment scheduler gates
  # Only assign groups when their gates are open (released_scope_id != null)
  assignments = compact([
    data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id != null ? {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id
    } : null,
    data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id != null ? {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id
    } : null,
    data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id != null ? {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id
    } : null,
  ])
}

# Monitoring outputs
output "deployment_dashboard" {
  value = {
    device_distribution = {
      ring_0_pilot      = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_0"])
      ring_1_broad      = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_1"])
      ring_2_production = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_2"])
      total_devices     = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_2"])
    }

    deployment_gates = {
      ring_0_pilot = {
        status   = data.microsoft365_utility_deployment_scheduler.ring_0_gate.condition_met ? "OPEN" : "CLOSED"
        message  = data.microsoft365_utility_deployment_scheduler.ring_0_gate.status_message
        group_id = data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id
      }
      ring_1_broad = {
        status   = data.microsoft365_utility_deployment_scheduler.ring_1_gate.condition_met ? "OPEN" : "CLOSED"
        message  = data.microsoft365_utility_deployment_scheduler.ring_1_gate.status_message
        group_id = data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id
      }
      ring_2_production = {
        status   = data.microsoft365_utility_deployment_scheduler.ring_2_gate.condition_met ? "OPEN" : "CLOSED"
        message  = data.microsoft365_utility_deployment_scheduler.ring_2_gate.status_message
        group_id = data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id
      }
    }

    active_assignments = length(compact([
      data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id,
    ]))
  }
  description = "Comprehensive view of device distribution, gate status, and active policy assignments"
}
