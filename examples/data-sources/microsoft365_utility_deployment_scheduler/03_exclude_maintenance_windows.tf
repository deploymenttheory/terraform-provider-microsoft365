# Exclude Maintenance Windows
# This example demonstrates blocking deployments during maintenance windows.
# Prevents changes during critical business periods or scheduled maintenance.

data "microsoft365_utility_deployment_scheduler" "avoid_maintenance" {
  name                  = "app-deployment-avoid-maintenance"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "c3d4e5f6-a7b8-9012-cdef-123456789012" # Device group

  time_condition = {
    delay_start_time_by = 0 # Deploy immediately when conditions allow
  }

  # Block deployments during maintenance windows
  exclusion_time_windows = {
    window = [
      {
        # Block during weekend maintenance
        days_of_week = ["saturday", "sunday"]
      },
      {
        # Block during year-end freeze (absolute dates)
        date_start = "2026-12-20T00:00:00Z"
        date_end   = "2027-01-05T23:59:59Z"
      },
      {
        # Block during nightly maintenance window
        time_of_day_start = "22:00:00" # 10 PM UTC
        time_of_day_end   = "06:00:00" # 6 AM UTC
      }
    ]
  }
}

output "maintenance_window_check" {
  value = {
    deployment_allowed = data.microsoft365_utility_deployment_scheduler.avoid_maintenance.condition_met
    status             = data.microsoft365_utility_deployment_scheduler.avoid_maintenance.status_message
    released_scope     = data.microsoft365_utility_deployment_scheduler.avoid_maintenance.released_scope_id
  }
  description = "Indicates if we're outside maintenance windows"
}
