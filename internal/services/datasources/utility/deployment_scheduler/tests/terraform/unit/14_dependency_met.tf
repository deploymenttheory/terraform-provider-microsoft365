# Unit Test: Dependency gate - satisfied
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Dependency Met"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 0
    minimum_open_hours               = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - prerequisite opened immediately and minimum_open_hours is 0"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID - dependency satisfied"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show dependency satisfied"
}
