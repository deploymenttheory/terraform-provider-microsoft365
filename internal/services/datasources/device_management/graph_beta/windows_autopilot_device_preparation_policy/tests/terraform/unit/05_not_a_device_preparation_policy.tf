# A settings catalog policy that is not created from an Autopilot device preparation template
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  policy_id = "33333333-3333-3333-3333-333333333333"
}
