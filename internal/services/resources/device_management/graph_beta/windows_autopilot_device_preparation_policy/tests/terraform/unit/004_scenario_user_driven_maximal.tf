resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_maximal" {
  name                  = "unit-test-autopilot-dpp-ud-maximal"
  description           = "Unit test for user-driven mode maximal configuration"
  role_scope_tag_ids    = ["0"]
  device_security_group = "00000000-0000-0000-0000-000000000001"

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1"
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0"
    account_type    = "enrollment_autopilot_dpp_accountype_1"
  }

  oobe_settings = {
    timeout_in_minutes   = 120
    custom_error_message = "Please contact your IT administrator for assistance with device setup."
    allow_skip           = true
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = "00000000-0000-0000-0000-000000000002"
      app_type = "win32LobApp"
    },
    {
      app_id   = "00000000-0000-0000-0000-000000000003"
      app_type = "winGetApp"
    },
    {
      app_id   = "00000000-0000-0000-0000-000000000004"
      app_type = "officeSuiteApp"
    }
  ]

  allowed_scripts = [
    "00000000-0000-0000-0000-000000000005",
    "00000000-0000-0000-0000-000000000006"
  ]
}
