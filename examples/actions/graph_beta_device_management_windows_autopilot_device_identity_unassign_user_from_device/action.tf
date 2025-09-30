# Unassign user from device action
action "microsoft365_graph_beta_device_management_unassign_user_from_device" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"

  timeouts = {
    create = "5m"
  }
}