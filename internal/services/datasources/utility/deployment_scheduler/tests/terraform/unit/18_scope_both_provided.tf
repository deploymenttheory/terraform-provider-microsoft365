# Unit Test: Invalid - both scope_id and scope_ids provided (should fail validation)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Scope Both Provided"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"
  scope_ids             = ["87654321-4321-4321-4321-cba987654321"]

  time_condition = {
    delay_start_time_by = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should not reach here - validation should fail"
}
