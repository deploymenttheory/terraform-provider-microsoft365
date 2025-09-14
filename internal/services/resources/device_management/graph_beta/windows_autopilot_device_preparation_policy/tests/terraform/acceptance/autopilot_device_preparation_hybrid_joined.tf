resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "hybrid_joined" {
  name                  = "acc-test-windows-autopilot-device-preparation-policy-hybrid-joined"
  description           = "acc-test-windows-autopilot-device-preparation-policy-hybrid-joined mode"
  role_scope_tag_ids    = ["0"]
  device_security_group = microsoft365_graph_beta_groups_group.acc_test_group_1.id

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_0" # Standard mode
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    join_type       = "enrollment_autopilot_dpp_jointype_1"       # Entra ID hybrid joined
    account_type    = "enrollment_autopilot_dpp_accountype_0"     # Standard User
  }

  oobe_settings = {
    timeout_in_minutes   = 75
    custom_error_message = "Hybrid join setup failed. Contact your system administrator."
    allow_skip           = false
    allow_diagnostics    = false
  }

  assignments = {
    include_group_ids = [
      microsoft365_graph_beta_groups_group.acc_test_group_1.id,
      microsoft365_graph_beta_groups_group.acc_test_group_2.id
    ]
  }
}