resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "available_ios_lob" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "available"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    ios_lob = {
      prevent_managed_app_backup = true
    }
  }
}
