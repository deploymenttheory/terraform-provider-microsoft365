# Example 1: Assign user to Autopilot device
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_assign_user_to_device" "assign_user" {
  config {
    windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789abc"
    user_principal_name                  = "user@contoso.com"
    addressable_user_name                = "John Doe"
  }
}

# Example 2: Assign user to Autopilot device with extended timeout
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_assign_user_to_device" "assign_user_extended" {
  config {
    windows_autopilot_device_identity_id = "87654321-4321-4321-4321-ba9876543210"
    user_principal_name                  = "jane.smith@contoso.com"
    addressable_user_name                = "Jane Smith"

    timeouts = {
      invoke = "10m"
    }
  }
}
