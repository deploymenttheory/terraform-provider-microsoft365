resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_max_assign" {
  name                  = "unit-test-autopilot-dpp-ud-max-assign"
  description           = "Unit test for user-driven mode with maximal assignments"
  role_scope_tag_ids    = ["0"]
  device_security_group = "00000000-0000-0000-0000-000000000001"

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1"
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0"
    account_type    = "enrollment_autopilot_dpp_accountype_1"
  }

  oobe_settings = {
    timeout_in_minutes   = 90
    custom_error_message = "Please contact your IT administrator for assistance."
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
    }
  ]

  allowed_scripts = [
    "00000000-0000-0000-0000-000000000004"
  ]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000005"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000006"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000007"
    }
  ]
}
