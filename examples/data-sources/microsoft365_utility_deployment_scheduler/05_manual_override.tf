# Manual Override Example
# This example demonstrates using manual_override to immediately release a deployment
# regardless of time conditions. Useful for emergency changes or hotfixes.

# Normal scheduled deployment
data "microsoft365_utility_deployment_scheduler" "security_patch_release" {
  name                  = "critical-security-patch"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "a7b8c9d0-e1f2-3456-a012-3456789abcde"

  # Normal schedule: 7 day delay
  time_condition = {
    delay_start_time_by = 168 # 7 days for normal patches
  }

  # Only deploy during business hours
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "17:00:00"
      }
    ]
  }

  # Manual override bypasses ALL conditions when set to true
  # To activate emergency release: change false to true and apply
  manual_override = false
}

output "patch_deployment_status" {
  value = {
    gate_open      = data.microsoft365_utility_deployment_scheduler.security_patch_release.condition_met
    scope_released = data.microsoft365_utility_deployment_scheduler.security_patch_release.released_scope_id
    status         = data.microsoft365_utility_deployment_scheduler.security_patch_release.status_message
  }
  description = "Monitor patch deployment status"
}

# Usage:
# Normal deployment: Leave manual_override = false (follows schedule)
# Emergency release: Change manual_override = true and run terraform apply
