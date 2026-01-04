# Example 1: Update Autopilot device properties
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_update_device_properties" "update_properties" {
  config {
    windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789abc"
    user_principal_name                  = "user@contoso.com"
    addressable_user_name                = "John Doe"
    group_tag                            = "Finance"
    display_name                         = "John's Laptop"

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Update Autopilot device properties with minimal fields
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_update_device_properties" "update_minimal" {
  config {
    windows_autopilot_device_identity_id = "87654321-4321-4321-4321-ba9876543210"
    display_name                         = "IT Department Laptop"

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 3: Update multiple Autopilot device properties
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_update_device_properties" "update_extended" {
  config {
    windows_autopilot_device_identity_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    user_principal_name                  = "jane.smith@contoso.com"
    addressable_user_name                = "Jane Smith"
    group_tag                            = "Marketing"
    display_name                         = "Jane's Surface Pro"

    timeouts = {
      invoke = "10m"
    }
  }
}
