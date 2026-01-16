# Example: Phased Windows Update Rollout using Deployment Scheduler

# Define the deployment campaign start time (all phases calculated from this)
# This works similar to Unix epoch time - a single reference point for all calculations
locals {
  update_deployment_start = "2024-01-20T00:00:00Z"
}

# ============================================================================
# Example 1: Using scope_id (singular) for deploying to a single group
# ============================================================================

# Phase 1: Pilot Group - Immediate deployment (single group)
data "microsoft365_utility_deployment_scheduler" "phase1_pilot" {
  name                   = "Phase 1 - Pilot Group"
  deployment_start_time  = local.update_deployment_start
  scope_id               = "pilot-group-abc-123"  # Single group ID

  # No time_condition = immediate release at deployment start time
}

# Phase 2: Early Adopters - Deploy 48 hours after deployment start (single group)
data "microsoft365_utility_deployment_scheduler" "phase2_early_adopters" {
  name                   = "Phase 2 - Early Adopters"
  deployment_start_time  = local.update_deployment_start
  scope_id               = "early-adopters-def-456"  # Single group ID

  time_condition {
    offset_hours = 48  # Gate opens 48h after deployment_start_time
  }
}

# Phase 3: Production - Deploy 72 hours after deployment start (single group)
data "microsoft365_utility_deployment_scheduler" "phase3_production" {
  name                   = "Phase 3 - Production"
  deployment_start_time  = local.update_deployment_start
  scope_id               = "production-ghi-789"  # Single group ID

  time_condition {
    offset_hours = 72  # Gate opens 72h after deployment_start_time
  }
}

# Windows Update Policy using the phased deployment
resource "microsoft365_windows_update_policy" "january_updates" {
  name        = "January 2024 Security Updates"
  description = "Phased rollout of January security updates"

  # Policy configuration...
  # (add your update policy settings here)

  # Assignment using deployment scheduler
  # compact() removes null values when gates are closed
  assignment {
    target {
      group_ids = compact([
        data.microsoft365_utility_deployment_scheduler.phase1_pilot.released_scope_id,
        data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.released_scope_id,
        data.microsoft365_utility_deployment_scheduler.phase3_production.released_scope_id,
      ])
    }
  }
}

# Monitor deployment status
output "deployment_status" {
  description = "Current status of phased rollout"
  value = {
    phase1_pilot = {
      ready         = data.microsoft365_utility_deployment_scheduler.phase1_pilot.condition_met
      status        = data.microsoft365_utility_deployment_scheduler.phase1_pilot.status_message
      scope_id      = data.microsoft365_utility_deployment_scheduler.phase1_pilot.released_scope_id
    }
    phase2_early_adopters = {
      ready         = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.condition_met
      status        = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.status_message
      scope_id      = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.released_scope_id
      hours_elapsed = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.conditions_detail.time_condition_detail.hours_elapsed
    }
    phase3_production = {
      ready         = data.microsoft365_utility_deployment_scheduler.phase3_production.condition_met
      status        = data.microsoft365_utility_deployment_scheduler.phase3_production.status_message
      scope_id      = data.microsoft365_utility_deployment_scheduler.phase3_production.released_scope_id
      hours_elapsed = data.microsoft365_utility_deployment_scheduler.phase3_production.conditions_detail.time_condition_detail.hours_elapsed
    }
  }
}

# ============================================================================
# Example 2: Using scope_ids (plural) for deploying to multiple entities
# ============================================================================

# Deploy to multiple specific devices after 24 hours
data "microsoft365_utility_deployment_scheduler" "device_rollout" {
  name                   = "Specific Devices Rollout"
  deployment_start_time  = local.update_deployment_start
  scope_ids              = [
    "device-guid-001",
    "device-guid-002",
    "device-guid-003",
  ]

  time_condition {
    offset_hours = 24
  }
}

# Use with a policy
resource "microsoft365_device_configuration_policy" "device_specific_config" {
  name = "Device Specific Configuration"

  # When gate is open: released_scope_ids contains the device GUIDs
  # When gate is closed: released_scope_ids is null
  assignment {
    target {
      device_ids = data.microsoft365_utility_deployment_scheduler.device_rollout.released_scope_ids
    }
  }
}

# ============================================================================
# Example 3: Without deployment_start_time (not recommended for time conditions)
# ============================================================================

# This will use current time on each evaluation
# Time-based conditions won't work as expected across multiple applies
data "microsoft365_utility_deployment_scheduler" "adhoc_deployment" {
  name      = "Ad-hoc Deployment"
  scope_id  = "adhoc-group-xyz"

  # No deployment_start_time = uses current time each time Terraform runs
  # This means the gate will always show 0 hours elapsed
  time_condition {
    offset_hours = 24
  }
}

# ============================================================================
# Example 4: Multiple independent deployment campaigns
# ============================================================================

locals {
  security_patch_start  = "2024-01-15T00:00:00Z"
  feature_update_start  = "2024-02-01T00:00:00Z"
}

# Security patch campaign
data "microsoft365_utility_deployment_scheduler" "security_patch_phase1" {
  name                   = "Security Patches - Phase 1"
  deployment_start_time  = local.security_patch_start
  scope_id               = "security-pilot-123"

  time_condition {
    offset_hours = 0  # Immediate at deployment start
  }
}

# Feature update campaign (different start time)
data "microsoft365_utility_deployment_scheduler" "feature_update_phase1" {
  name                   = "Feature Updates - Phase 1"
  deployment_start_time  = local.feature_update_start  # Different campaign
  scope_id               = "feature-pilot-456"

  time_condition {
    offset_hours = 0
  }
}

# ============================================================================
# Example 5: No time condition (immediate deployment)
# ============================================================================

# This gate opens immediately - useful for organizing configs
data "microsoft365_utility_deployment_scheduler" "immediate_deployment" {
  name      = "Immediate Deployment"
  scope_id  = "always-on-group-123"

  # No time_condition = no waiting, immediate release
  # deployment_start_time is optional when no time condition is used
}

# ============================================================================
# What happens on successive terraform applies:
# ============================================================================

# terraform apply at 2024-01-20 00:00:00Z (deployment start)
#   phase1_pilot:
#     - status: "No time condition specified (immediate release)"
#     - released_scope_id: "pilot-group-abc-123" ✓
#   phase2_early_adopters:
#     - status: "Waiting: Time condition NOT met (0.0h/48h required)"
#     - released_scope_id: null
#   phase3_production:
#     - status: "Waiting: Time condition NOT met (0.0h/72h required)"
#     - released_scope_id: null

# terraform apply at 2024-01-22 02:00:00Z (50 hours after start)
#   phase1_pilot:
#     - released_scope_id: "pilot-group-abc-123" ✓
#   phase2_early_adopters:
#     - status: "Conditions met: Time condition met (50.0h/48h required)"
#     - released_scope_id: "early-adopters-def-456" ✓
#   phase3_production:
#     - status: "Waiting: Time condition NOT met (50.0h/72h required)"
#     - released_scope_id: null

# terraform apply at 2024-01-23 02:00:00Z (74 hours after start)
#   phase1_pilot:
#     - released_scope_id: "pilot-group-abc-123" ✓
#   phase2_early_adopters:
#     - released_scope_id: "early-adopters-def-456" ✓
#   phase3_production:
#     - status: "Conditions met: Time condition met (74.0h/72h required)"
#     - released_scope_id: "production-ghi-789" ✓
