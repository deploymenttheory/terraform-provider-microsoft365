# Unit Test: Inclusion window - time of day (00:00-23:59 = always)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Inclusion Window Time of Day"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  inclusion_time_windows = {
    window = [
      {
        time_of_day_start = "00:00:00"
        time_of_day_end   = "23:59:59"
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should always be true - window covers entire day"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should always be the GUID"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show gate open"
}
