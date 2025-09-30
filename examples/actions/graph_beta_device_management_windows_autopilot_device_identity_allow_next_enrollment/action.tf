# Allow next enrollment for an Autopilot device
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_allow_next_enrollment" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789abc"

  timeouts = {
    create = "5m"
  }
}