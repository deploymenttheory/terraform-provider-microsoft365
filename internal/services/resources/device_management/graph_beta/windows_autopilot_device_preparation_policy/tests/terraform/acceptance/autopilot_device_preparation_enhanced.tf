resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "enhanced" {
  name                  = "acc-test-windows-autopilot-device-preparation-policy-enhanced"
  description           = "acc-test-windows-autopilot-device-preparation-policy-enhanced with all features"
  role_scope_tag_ids    = ["0"]
  device_security_group = data.azuread_group.test_autopilot_security_group.object_id

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1" # Enhanced mode
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1" # Self-deploying
    join_type       = "enrollment_autopilot_dpp_jointype_0"       # Entra ID joined
    account_type    = "enrollment_autopilot_dpp_accountype_1"     # Administrator
  }

  oobe_settings = {
    timeout_in_minutes   = 120
    custom_error_message = "Please contact your IT administrator for assistance with device setup."
    allow_skip           = true
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = microsoft365_graph_beta_device_management_win32_lob_app.test_app.id
      app_type = "win32LobApp"
    }
  ]

  allowed_scripts = [
    microsoft365_graph_beta_device_management_device_shell_script.test_script.id
  ]

  assignments = {
    include_group_ids = [
      data.azuread_group.test_group1.object_id,
      data.azuread_group.test_group2.object_id,
      data.azuread_group.test_group3.object_id
    ]
  }
}