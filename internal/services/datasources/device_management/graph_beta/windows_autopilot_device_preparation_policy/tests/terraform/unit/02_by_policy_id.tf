# Lookup a Windows Autopilot device preparation policy by policy id
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  policy_id = "11111111-1111-1111-1111-111111111111"
}
