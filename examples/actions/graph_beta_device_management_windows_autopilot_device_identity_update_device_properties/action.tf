# Update device properties action
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_update_device_properties" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"
  user_principal_name                  = "user@contoso.com"
  addressable_user_name                = "John Doe"
  group_tag                            = "Finance"
  display_name                         = "John's Laptop"

  timeouts = {
    create = "5m"
  }
}