# Unit Test: Invalid scope_id format (not a GUID - should fail validation)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Invalid Scope ID Format"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "not-a-valid-guid"

  time_condition = {
    delay_start_time_by = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should not reach here - validation should fail"
}
