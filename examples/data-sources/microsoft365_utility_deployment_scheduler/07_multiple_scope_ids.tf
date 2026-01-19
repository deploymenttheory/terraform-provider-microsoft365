# Multiple Scope IDs Example
# This example demonstrates releasing multiple group IDs simultaneously.
# Useful when you want to apply changes to several groups at once with the same timing.

locals {
  regional_offices = [
    "c9d0e1f2-a3b4-5678-c234-56789abcdef0", # US East office
    "d0e1f2a3-b4c5-6789-d345-6789abcdef01", # US West office
    "e1f2a3b4-c5d6-7890-e456-789abcdef012", # EMEA office
    "f2a3b4c5-d6e7-8901-f567-89abcdef0123", # APAC office
  ]
}

data "microsoft365_utility_deployment_scheduler" "regional_rollout" {
  name                  = "regional-offices-unified-rollout"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_ids             = local.regional_offices # Use scope_ids for multiple groups

  time_condition = {
    delay_start_time_by = 72 # 3 day delay for all regions
  }

  # Deploy only during business hours (respects all timezones effectively)
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
        time_of_day_start = "10:00:00" # Conservative window to catch multiple timezones
        time_of_day_end   = "15:00:00"
      }
    ]
  }
}

output "regional_rollout_status" {
  value = {
    gate_open          = data.microsoft365_utility_deployment_scheduler.regional_rollout.condition_met
    number_of_regions  = length(local.regional_offices)
    released_group_ids = data.microsoft365_utility_deployment_scheduler.regional_rollout.released_scope_ids
    status             = data.microsoft365_utility_deployment_scheduler.regional_rollout.status_message
  }
  description = "Status of unified regional rollout"
}