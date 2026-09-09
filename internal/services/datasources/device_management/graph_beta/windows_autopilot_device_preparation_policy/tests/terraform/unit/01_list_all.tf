# List all Windows Autopilot device preparation policies
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  list_all = true
}
