# Unit Test: Invalid negative minimum_open_hours (should fail validation)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Dependency Negative Minimum"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 0
    minimum_open_hours                = -24
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should not reach here - validation should fail"
}
