# Unit Test: Dependency gate - prerequisite open but insufficient time
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Dependency Insufficient Time"
  deployment_start_time = "2099-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 0
    minimum_open_hours               = 168
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be false - prerequisite opened but hasn't been open for 168h (deployment_start_time in future)"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be null - minimum open hours not met"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show dependency not satisfied"
}
