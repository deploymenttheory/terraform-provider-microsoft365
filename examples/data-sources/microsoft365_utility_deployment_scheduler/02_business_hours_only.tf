# Business Hours Only Deployment
# This example demonstrates releasing changes only during business hours.
# Useful for changes that might require user interaction or support availability.

data "microsoft365_utility_deployment_scheduler" "business_hours_release" {
  name                  = "settings-catalog-business-hours"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "b2c3d4e5-f6a7-8901-bcde-f12345678901" # Entra ID Group for workstations

  time_condition = {
    delay_start_time_by = 0 # No delay, but constrained by inclusion windows
  }

  # Only allow deployment during business hours (UTC)
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00" # 9 AM UTC
        time_of_day_end   = "17:00:00" # 5 PM UTC
      }
    ]
  }
}

output "business_hours_status" {
  value = {
    gate_open          = data.microsoft365_utility_deployment_scheduler.business_hours_release.condition_met
    current_time_valid = data.microsoft365_utility_deployment_scheduler.business_hours_release.status_message
    released_scope     = data.microsoft365_utility_deployment_scheduler.business_hours_release.released_scope_id
  }
  description = "Shows if we're currently in business hours deployment window"
}
