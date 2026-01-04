# Example 1: Unassign user from Autopilot device
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_unassign_user_from_device" "unassign_user" {
  config {
    windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789abc"
  }
}

# Example 2: Unassign user from Autopilot device with extended timeout
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_unassign_user_from_device" "unassign_user_extended" {
  config {
    windows_autopilot_device_identity_id = "87654321-4321-4321-4321-ba9876543210"

    timeouts = {
      invoke = "10m"
    }
  }
}
