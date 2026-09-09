# A name that does not match any device preparation policy
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  name = "Policy That Does Not Exist"
}
