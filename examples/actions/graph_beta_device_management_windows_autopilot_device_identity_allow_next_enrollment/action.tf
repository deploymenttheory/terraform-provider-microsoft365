# Example 1: Allow next enrollment for an Autopilot device
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_allow_next_enrollment" "allow_enrollment" {
  config {
    windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789abc"

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Allow next enrollment with extended timeout
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_allow_next_enrollment" "allow_enrollment_extended" {
  config {
    windows_autopilot_device_identity_id = "87654321-4321-4321-4321-ba9876543210"

    timeouts = {
      invoke = "10m"
    }
  }
}
