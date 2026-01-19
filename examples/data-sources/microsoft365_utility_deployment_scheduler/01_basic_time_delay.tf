# Basic Time Delay Example
# This example demonstrates a simple time-based delay before releasing a scope ID.
# The gate opens 24 hours after deployment_start_time.

data "microsoft365_utility_deployment_scheduler" "pilot_group_release" {
  name                  = "pilot-group-24h-delay"
  deployment_start_time = "2026-01-20T08:00:00Z"
  scope_id              = "a1b2c3d4-e5f6-7890-abcd-ef1234567890" # Entra ID Group ID for pilot users

  time_condition = {
    delay_start_time_by = 24 # Wait 24 hours before releasing
  }
}

# Use the released scope ID in your resource
# Example: Conditional Access policy is only assigned when gate opens
output "pilot_gate_status" {
  value = {
    condition_met  = data.microsoft365_utility_deployment_scheduler.pilot_group_release.condition_met
    released_scope = data.microsoft365_utility_deployment_scheduler.pilot_group_release.released_scope_id
    status_message = data.microsoft365_utility_deployment_scheduler.pilot_group_release.status_message
  }
  description = "Monitor gate status and timing"
}
