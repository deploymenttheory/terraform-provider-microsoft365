# Unit Test: Dependency gate - prerequisite hasn't opened yet
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Dependency Prerequisite Not Open"
  deployment_start_time = "2099-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 168
    minimum_open_hours               = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be false - prerequisite hasn't opened (needs 168h, deployment_start_time in future)"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be null - dependency not satisfied"
}

