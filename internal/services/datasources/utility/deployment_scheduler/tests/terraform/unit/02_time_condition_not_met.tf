# Unit Test: Time condition not met (require 48 hours, but condition not yet met)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Time Condition Not Met"
  deployment_start_time = "2099-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 48
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be false - deployment_start_time is in future"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be null when gate is closed"
}

