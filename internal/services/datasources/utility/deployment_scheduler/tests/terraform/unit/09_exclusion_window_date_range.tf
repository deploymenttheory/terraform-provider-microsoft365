# Unit Test: Exclusion window - date range (holiday freeze)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Exclusion Window Date Range"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  exclusion_time_windows = {
    window = [
      {
        date_start = "2026-12-20T00:00:00Z"
        date_end   = "2027-01-05T23:59:59Z"
      }
    ]
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "Should be true - current date (2026-01-16) is NOT in exclusion range"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "Should be the GUID - not blocked by exclusion"
}

output "status_message" {
  value       = data.microsoft365_utility_deployment_scheduler.test.status_message
  description = "Should show gate open, no exclusion active"
}
