resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "automatic_approval" {
  display_name  = "Test Automatic Approval Windows Driver Update Profile - Unique"
  approval_type = "automatic"
  description   = "Test description for automatic approval driver update profile"

  deployment_deferral_in_days = 5
  role_scope_tag_ids          = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
