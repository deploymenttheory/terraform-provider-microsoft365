# Assign user to device action
action "microsoft365_graph_beta_device_management_assign_user_to_device" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"
  user_principal_name                  = "user@contoso.com"
  addressable_user_name                = "John Doe"

  timeouts = {
    create = "5m"
  }
}