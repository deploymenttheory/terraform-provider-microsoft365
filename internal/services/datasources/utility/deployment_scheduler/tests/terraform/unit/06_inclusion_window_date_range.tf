# Unit Test: Inclusion window - date range
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Inclusion Window Date Range"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  inclusion_time_windows = {
    window = [
      {
        date_start = "2026-01-01T00:00:00Z"
        date_end   = "2026-12-31T23:59:59Z"
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - current date (2026-01-16) is within range"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID when within date range"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show gate open with inclusion window satisfied"
}
