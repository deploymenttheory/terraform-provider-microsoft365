data "microsoft365_utility_deployment_scheduler" "test" {
  name                   = "Test Invalid Day of Week"
  deployment_start_time  = "2024-01-01T00:00:00Z"
  scope_id               = "12345678-1234-1234-1234-123456789abc"
  
  time_condition = {
    delay_start_time_by = 0
  }

  # This should fail validation - "Monday" should be "monday"
  inclusion_time_windows = {
    window = [
      {
        days_of_week = ["Monday", "Tuesday"]
      }
    ]
  }
}

output "condition_met" {
  value = data.microsoft365_utility_deployment_scheduler.test.condition_met
}

output "released_scope_id" {
  value = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
}

output "status_message" {
  value = data.microsoft365_utility_deployment_scheduler.test.status_message
}
