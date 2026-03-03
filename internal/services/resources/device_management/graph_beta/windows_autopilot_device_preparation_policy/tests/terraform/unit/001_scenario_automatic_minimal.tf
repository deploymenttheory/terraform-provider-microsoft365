resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "auto_minimal" {
  name               = "unit-test-autopilot-dpp-auto-minimal"
  description        = "Unit test for automatic mode minimal configuration"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1"
  }

  allowed_apps = [
    {
      app_id   = "00000000-0000-0000-0000-000000000001"
      app_type = "winGetApp"
    }
  ]
}
