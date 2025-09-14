resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "maximal" {
  name                  = "unit-test-windows-autopilot-device-preparation-policy-maximal"
  description           = "unit-test-windows-autopilot-device-preparation-policy-maximal"
  role_scope_tag_ids    = ["0"]
  device_security_group = "00000000-0000-0000-0000-000000000001"

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1" # Enhanced mode
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1" # Self-deploying
    join_type       = "enrollment_autopilot_dpp_jointype_1"       # Entra ID hybrid joined
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
      app_id   = "00000000-0000-0000-0000-000000000003"
      app_type = "win32LobApp"
    },
    {
      app_id   = "00000000-0000-0000-0000-000000000004"
      app_type = "winGetApp"
    },
    {
      app_id   = "00000000-0000-0000-0000-000000000005"
      app_type = "officeSuiteApp"
    }
  ]

  allowed_scripts = [
    "00000000-0000-0000-0000-000000000006",
    "00000000-0000-0000-0000-000000000007"
  ]

  assignments = {
    include_group_ids = [
      "00000000-0000-0000-0000-000000000001",
      "00000000-0000-0000-0000-000000000002",
      "00000000-0000-0000-0000-000000000003"
    ]
  }
}