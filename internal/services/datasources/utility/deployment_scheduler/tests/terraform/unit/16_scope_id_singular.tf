# Unit Test: Using singular scope_id
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Scope ID Singular"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - time condition met"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the singular GUID"
}

output "released_scope_ids" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_ids
  description = "Should be null when using scope_id (singular)"
}

