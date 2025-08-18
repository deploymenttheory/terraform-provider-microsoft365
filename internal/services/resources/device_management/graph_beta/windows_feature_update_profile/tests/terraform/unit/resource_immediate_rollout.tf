resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "immediate_rollout" {
  display_name                                            = "Test Immediate Rollout Windows Feature Update Profile - Unique"
  description                                             = "Windows Feature Update Profile for testing immediate rollout scenarios"
  feature_update_version                                  = "Windows 10, version 22H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

  role_scope_tag_ids = ["8", "9"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}