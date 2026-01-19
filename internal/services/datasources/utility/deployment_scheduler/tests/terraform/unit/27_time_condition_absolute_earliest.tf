# Unit Test: Time condition with absolute_earliest
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Time Condition Absolute Earliest"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by    = 0
    absolute_earliest_time = "2025-01-01T00:00:00Z"
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - current time (2026) is after absolute_earliest (2025)"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show gate open, absolute_earliest constraint met"
}
