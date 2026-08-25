resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_019" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    windows_app_x = {
      use_device_context = true
    }
  }
}
