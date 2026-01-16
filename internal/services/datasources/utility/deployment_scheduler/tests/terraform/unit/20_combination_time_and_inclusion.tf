# Unit Test: Combination - time condition (met) AND inclusion window (weekdays)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Combination Time And Inclusion"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Depends on current day - true on weekdays, false on weekends"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "GUID on weekdays, null on weekends"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Shows both time and inclusion window status"
}
