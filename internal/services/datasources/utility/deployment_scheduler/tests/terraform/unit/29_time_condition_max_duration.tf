# Unit Test: Time condition with max_open_duration_hours
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Time Condition Max Duration"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by      = 0
    max_open_duration_hours = 17520  # 2 years in hours
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - gate has been open for ~2 years, max is 2 years"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show gate open, not exceeded max duration"
}
