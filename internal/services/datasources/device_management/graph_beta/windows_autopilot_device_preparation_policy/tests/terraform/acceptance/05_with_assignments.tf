# Lookup a device preparation policy and its assignments
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  list_all         = true
  list_assignments = true
}
