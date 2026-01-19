# Unit Test: Using plural scope_ids
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Scope IDs Plural"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_ids             = ["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-cba987654321"]

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
  description = "Should be null when using scope_ids (plural)"
}

output "released_scope_ids" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_ids
  description = "Should be the list of 2 GUIDs"
}

