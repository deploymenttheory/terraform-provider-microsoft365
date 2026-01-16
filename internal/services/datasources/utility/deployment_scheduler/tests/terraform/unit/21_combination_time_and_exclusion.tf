# Unit Test: Combination - time condition and exclusion window
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Combination Time And Exclusion"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  exclusion_time_windows = {
    window = [
      {
        days_of_week = ["saturday", "sunday"]
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Depends on current day - false on weekends (blocked), true on weekdays"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Null on weekends, GUID on weekdays"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Shows time met but blocked by exclusion on weekends"
}
