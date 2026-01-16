# Unit Test: Inclusion window - multiple windows
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Inclusion Window Multiple"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "wednesday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "12:00:00"
      },
      {
        days_of_week      = ["tuesday", "thursday"]
        time_of_day_start = "14:00:00"
        time_of_day_end   = "17:00:00"
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Depends on current day and time"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "GUID if within any window, null otherwise"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Shows which windows match or don't match"
}
