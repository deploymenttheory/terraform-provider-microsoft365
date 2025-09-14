resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "enhanced" {
  name                  = "acc-test-windows-autopilot-device-preparation-policy-enhanced"
  description           = "acc-test-windows-autopilot-device-preparation-policy-enhanced with all features"
  role_scope_tag_ids    = ["0"]
  device_security_group = microsoft365_graph_beta_groups_group.acc_test_group_1.id

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
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.test_firefox.id
      app_type = "winGetApp"
    },
    {
      app_id   = microsoft365_graph_beta_device_and_app_management_office_suite_app.office_365_config_designer.id
      app_type = "officeSuiteApp"
    }
  ]

  allowed_scripts = [
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_windows_platform_script_01.id
  ]

  assignments = {
    include_group_ids = [
      microsoft365_graph_beta_groups_group.acc_test_group_1.id,
      microsoft365_graph_beta_groups_group.acc_test_group_2.id,
      microsoft365_graph_beta_groups_group.acc_test_group_3.id
    ]
  }
}