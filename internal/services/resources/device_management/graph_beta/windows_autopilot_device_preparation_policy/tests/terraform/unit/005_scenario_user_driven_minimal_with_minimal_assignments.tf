resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_min_assign" {
  name                  = "unit-test-autopilot-dpp-ud-min-assign"
  description           = "Unit test for user-driven mode with minimal assignments"
  role_scope_tag_ids    = ["0"]
  device_security_group = "00000000-0000-0000-0000-000000000001"

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_0"
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0"
    account_type    = "enrollment_autopilot_dpp_accountype_0"
  }

  oobe_settings = {
    timeout_in_minutes   = 60
    custom_error_message = "Contact your organization's support person for help."
    allow_skip           = false
    allow_diagnostics    = false
  }

  allowed_apps = [
    {
      app_id   = "00000000-0000-0000-0000-000000000002"
      app_type = "win32LobApp"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000003"
    }
  ]
}
