# Unit Test: Invalid - neither scope_id nor scope_ids provided (should fail validation)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Scope Neither Provided"
  deployment_start_time = "2024-01-01T00:00:00Z"

  time_condition = {
    delay_start_time_by = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should not reach here - validation should fail"
}
