# Real-World Settings Catalog Deployment
# This example demonstrates a production-ready Settings Catalog policy deployment
# with proper phasing, timing, and safety controls for Windows security settings.

locals {
  policy_name = "Windows Security Baseline v2.0"

  # Define deployment groups
  pilot_users_group       = "a1a2a3a4-b5b6-c7c8-d9d0-e1e2e3e4e5e6" # 5-10 pilot devices
  it_department_group     = "b2b3b4b5-c6c7-d8d9-e0e1-f2f3f4f5f6f7" # IT team for validation
  production_workstations = "c3c4c5c6-d7d8-e9e0-f1f2-a3a4a5a6a7a8" # All production Windows devices
}

# Phase 1: Pilot Users - Opens after 24h, business hours only
data "microsoft365_utility_deployment_scheduler" "pilot_phase" {
  name                  = "${local.policy_name}-pilot"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = local.pilot_users_group

  time_condition = {
    delay_start_time_by = 24 # 24h delay for safety
  }

  # Business hours only for pilot
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

# Phase 2: IT Department - Opens 72h after pilot opens
data "microsoft365_utility_deployment_scheduler" "it_phase" {
  name                  = "${local.policy_name}-it-dept"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = local.it_department_group

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

# Phase 3: Production - Opens 1 week after IT phase opens
data "microsoft365_utility_deployment_scheduler" "production_phase" {
  name                  = "${local.policy_name}-production"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = local.production_workstations

  time_condition = {
    delay_start_time_by     = 192                    # 24 + 72 + 96 (total: 8 days)
    absolute_latest         = "2026-03-01T23:59:59Z" # Must complete by end of Feb (compliance)
    max_open_duration_hours = 720                    # Auto-close after 30 days
  }

  # Wait for IT phase to be open for 1 week
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

  exclusion_time_windows = {
    window = [
      {
        # No Friday deployments to prod
        days_of_week = ["friday"]
      },
      {
        # Month-end freeze
        date_start = "2026-01-29T00:00:00Z"
        date_end   = "2026-02-02T23:59:59Z"
      }
    ]
  }
}

# Monitoring and Reporting
output "deployment_dashboard" {
  value = {
    policy_name = local.policy_name

    phase_1_pilot = {
      status        = data.microsoft365_utility_deployment_scheduler.pilot_phase.condition_met ? "OPEN" : "CLOSED"
      group_id      = data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id
      status_detail = data.microsoft365_utility_deployment_scheduler.pilot_phase.status_message
    }

    phase_2_it = {
      status        = data.microsoft365_utility_deployment_scheduler.it_phase.condition_met ? "OPEN" : "CLOSED"
      group_id      = data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id
      status_detail = data.microsoft365_utility_deployment_scheduler.it_phase.status_message
    }

    phase_3_production = {
      status        = data.microsoft365_utility_deployment_scheduler.production_phase.condition_met ? "OPEN" : "CLOSED"
      group_id      = data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id
      status_detail = data.microsoft365_utility_deployment_scheduler.production_phase.status_message
    }

    total_groups_deployed = length(compact([
      data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id,
    ]))
  }

  description = "Complete rollout status dashboard for Settings Catalog policy deployment"
}

# Example of how you'd use this in a real Settings Catalog policy
# (Schema shown is illustrative - adjust to match actual resource)
#
# resource "microsoft365_graph_beta_device_management_configuration_policy" "security_baseline" {
#   name         = local.policy_name
#   description  = "Windows Security Baseline settings for corporate devices"
#   platforms    = "windows10"
#   technologies = ["mdm"]
#   
#   # Assignments based on which gates are open
#   assignments = compact([
#     data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id != null ? {
#       target = {
#         group_id = data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id
#       }
#     } : null,
#     data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id != null ? {
#       target = {
#         group_id = data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id
#       }
#     } : null,
#     data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id != null ? {
#       target = {
#         group_id = data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id
#       }
#     } : null,
#   ])
#   
#   # ... settings configuration ...
# }
