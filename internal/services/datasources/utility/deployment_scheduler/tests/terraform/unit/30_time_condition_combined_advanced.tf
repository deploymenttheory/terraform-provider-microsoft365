# Unit Test: Time condition with multiple advanced constraints
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Time Condition Combined Advanced"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by    = 0
    absolute_earliest_time = "2025-01-01T00:00:00Z"
    absolute_latest_time   = "2027-12-31T23:59:59Z"
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - current time (2026) is between absolute_earliest (2025) and absolute_latest (2027)"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID"
}

