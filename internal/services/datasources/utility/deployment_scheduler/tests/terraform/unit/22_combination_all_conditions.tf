# Unit Test: Combination - all conditions (time, inclusion, exclusion, dependency)
data "microsoft365_utility_deployment_scheduler" "test" {
  name                  = "Test - Combination All Conditions"
  deployment_start_time = "2024-01-01T00:00:00Z"
  scope_id              = "12345678-1234-1234-1234-123456789abc"

  time_condition = {
    delay_start_time_by = 0
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "17:00:00"
      }
    ]
  }

  exclusion_time_windows = {
    window = [
      {
        date_start = "2026-12-20T00:00:00Z"
        date_end   = "2027-01-05T23:59:59Z"
      }
    ]
  }

  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 0
    minimum_open_hours               = 0
  }
}

output "condition_met" {
  value       = data.microsoft365_utility_deployment_scheduler.test.condition_met
  description = "True if: weekday AND 09:00-17:00 AND NOT in holiday freeze (Dec 20 - Jan 5)"
}

output "released_scope_id" {
  value       = data.microsoft365_utility_deployment_scheduler.test.released_scope_id
  description = "GUID when all conditions pass, null otherwise"
}

