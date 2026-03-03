resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "auto_maximal" {
  name               = "unit-test-autopilot-dpp-auto-maximal"
  description        = "Unit test for automatic mode maximal configuration"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1"
  }

  allowed_apps = [
    {
      app_id   = "00000000-0000-0000-0000-000000000001"
      app_type = "win32LobApp"
    },
    {
      app_id   = "00000000-0000-0000-0000-000000000002"
      app_type = "winGetApp"
    },
    {
      app_id   = "00000000-0000-0000-0000-000000000003"
      app_type = "officeSuiteApp"
    }
  ]

  allowed_scripts = [
    "00000000-0000-0000-0000-000000000004",
    "00000000-0000-0000-0000-000000000005"
  ]
}
