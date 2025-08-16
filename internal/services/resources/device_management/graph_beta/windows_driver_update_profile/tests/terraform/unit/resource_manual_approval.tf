resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "manual_approval" {
  display_name  = "Test Manual Approval Windows Driver Update Profile - Unique"
  approval_type = "manual"
  description   = "Test description for manual approval driver update profile"
  
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
