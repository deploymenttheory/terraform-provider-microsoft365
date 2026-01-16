# Unit Test: Manual override bypasses all conditions
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Manual Override Bypasses All"
  deployment_start_time = "2099-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"
  manual_override       = true

  time_condition = {
    delay_start_time_by = 168
  }

  exclusion_time_windows = {
    window = [
      {
        days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - manual override bypasses time and exclusion"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID - all conditions bypassed"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show manual override enabled"
}
