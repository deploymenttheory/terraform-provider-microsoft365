resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_015" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    android_managed_store = {
      auto_update_mode                    = "priority"
      android_managed_store_app_track_ids = ["production"]
    }
  }
}
