# Lookup a Windows Autopilot device preparation policy by name
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  name = "Autopilot Device Preparation - User Driven"
}
