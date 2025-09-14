resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "minimal" {
  name                  = "acc-test-windows-autopilot-device-preparation-policy-minimal"
  description           = "acc-test-windows-autopilot-device-preparation-policy-minimal"
  role_scope_tag_ids    = ["0"]
  device_security_group = data.azuread_group.test_autopilot_security_group.object_id

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_0" # Standard mode
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    join_type       = "enrollment_autopilot_dpp_jointype_0"       # Entra ID joined
    account_type    = "enrollment_autopilot_dpp_accountype_0"     # Standard User
  }

  oobe_settings = {
    timeout_in_minutes   = 60
    custom_error_message = "Contact your organization's support person for help."
    allow_skip           = false
    allow_diagnostics    = false
  }

  assignments = {
    include_group_ids = [
      data.azuread_group.test_group1.object_id,
      data.azuread_group.test_group2.object_id
    ]
  }
}