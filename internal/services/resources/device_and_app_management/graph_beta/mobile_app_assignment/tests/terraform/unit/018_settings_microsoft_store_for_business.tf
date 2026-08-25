resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_018" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    microsoft_store_for_business = {
      use_device_context = true
    }
  }
}
