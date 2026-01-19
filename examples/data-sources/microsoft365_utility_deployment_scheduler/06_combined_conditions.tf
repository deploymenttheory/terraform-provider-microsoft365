# Combined Conditions Example
# This example demonstrates combining multiple scheduling constraints:
# - Time delay with absolute deadlines
# - Business hours only
# - Exclude maintenance windows
# Useful for complex deployment scenarios with multiple requirements.

data "microsoft365_utility_deployment_scheduler" "complex_release" {
  name                  = "compliance-policy-complex-schedule"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "b8c9d0e1-f2a3-4567-b123-456789abcdef"

  # Time-based constraints with deadlines
  time_condition = {
    delay_start_time_by     = 48                     # Wait 48 hours minimum
    absolute_earliest       = "2026-01-22T00:00:00Z" # Don't release before this date
    absolute_latest         = "2026-02-15T23:59:59Z" # Must release before this date (compliance deadline)
    max_open_duration_hours = 168                    # Auto-close gate after 7 days (1 week deployment window)
  }

  # Only deploy during business hours on weekdays
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "18:00:00"
      }
    ]
  }

  # Exclude sensitive periods
  exclusion_time_windows = {
    window = [
      {
        # Exclude weekends
        days_of_week = ["saturday", "sunday"]
      },
      {
        # Exclude nightly maintenance
        time_of_day_start = "22:00:00"
        time_of_day_end   = "06:00:00"
      },
      {
        # Exclude specific high-activity period (month-end)
        date_start = "2026-01-29T00:00:00Z"
        date_end   = "2026-02-02T23:59:59Z"
      }
    ]
  }

  # All conditions must be met
  require_all_conditions = true
}

# Detailed monitoring output
output "complex_schedule_status" {
  value = {
    gate_status = {
      open           = data.microsoft365_utility_deployment_scheduler.complex_release.condition_met
      status_message = data.microsoft365_utility_deployment_scheduler.complex_release.status_message
      scope_released = data.microsoft365_utility_deployment_scheduler.complex_release.released_scope_id
    }

    conditions_detail = data.microsoft365_utility_deployment_scheduler.complex_release.conditions_detail

    timing = {
      deployment_start = data.microsoft365_utility_deployment_scheduler.complex_release.deployment_start_time
      window_closes    = "2026-02-15T23:59:59Z" # absolute_latest
    }
  }
  description = "Comprehensive status of the complex deployment schedule"
}
