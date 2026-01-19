# Phased Rollout with Dependency Chaining
# This example demonstrates a multi-phase rollout where each phase depends on the previous.
# Phase 1: Pilot users (opens after 24h)
# Phase 2: Early adopters (opens after pilot has been open for 48h)
# Phase 3: Production users (opens after early adopters has been open for 72h)

locals {
  deployment_start = "2026-01-20T08:00:00Z"
  pilot_group_id   = "d4e5f6a7-b8c9-0123-def0-123456789abc"
  early_adopter_id = "e5f6a7b8-c9d0-1234-ef01-23456789abcd"
  production_id    = "f6a7b8c9-d0e1-2345-f012-3456789abcde"
}

# Phase 1: Pilot Group (opens after 24h delay)
data "microsoft365_utility_deployment_scheduler" "phase1_pilot" {
  name                  = "mfa-policy-phase1-pilot"
  deployment_start_time = local.deployment_start
  scope_id              = local.pilot_group_id

  time_condition = {
    delay_start_time_by = 24 # 24h soak time for initial testing
  }
}

# Phase 2: Early Adopters (opens after pilot has been open for 48h)
data "microsoft365_utility_deployment_scheduler" "phase2_early_adopters" {
  name                  = "mfa-policy-phase2-early-adopters"
  deployment_start_time = local.deployment_start
  scope_id              = local.early_adopter_id

  time_condition = {
    delay_start_time_by = 24 # Same base delay as pilot
  }

  # Wait for pilot phase to be open for 48h before opening this gate
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 24 # Pilot's delay
    minimum_open_hours               = 48 # Wait 48h after pilot opens
  }
}

# Phase 3: Production (opens after early adopters has been open for 72h)
data "microsoft365_utility_deployment_scheduler" "phase3_production" {
  name                  = "mfa-policy-phase3-production"
  deployment_start_time = local.deployment_start
  scope_id              = local.production_id

  time_condition = {
    delay_start_time_by = 72 # Base delay to align with early adopters (24 + 48)
  }

  # Wait for early adopters phase to be open for 72h before opening this gate
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 72 # Early adopters effective delay (24 + 48)
    minimum_open_hours               = 72 # Wait 72h after early adopters opens
  }
}

# Monitoring outputs
output "rollout_status" {
  value = {
    phase1_pilot = {
      open   = data.microsoft365_utility_deployment_scheduler.phase1_pilot.condition_met
      scope  = data.microsoft365_utility_deployment_scheduler.phase1_pilot.released_scope_id
      status = data.microsoft365_utility_deployment_scheduler.phase1_pilot.status_message
    }
    phase2_early_adopters = {
      open   = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.condition_met
      scope  = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.released_scope_id
      status = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.status_message
    }
    phase3_production = {
      open   = data.microsoft365_utility_deployment_scheduler.phase3_production.condition_met
      scope  = data.microsoft365_utility_deployment_scheduler.phase3_production.released_scope_id
      status = data.microsoft365_utility_deployment_scheduler.phase3_production.status_message
    }
  }
  description = "Current status of all rollout phases"
}
