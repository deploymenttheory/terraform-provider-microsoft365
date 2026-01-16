# Unit Test: Invalid negative delay_start_time_by (should fail validation)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Time Condition Negative"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = -24
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should not reach here - validation should fail"
}
