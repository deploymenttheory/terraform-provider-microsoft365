resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "minimal" {
  display_name  = "Test Minimal Windows Driver Update Profile - Unique"
  approval_type = "manual"
  description   = "Test description for minimal driver update profile"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
