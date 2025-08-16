resource "random_uuid" "test_automatic" {
}

resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "test_automatic" {
  display_name  = "Acceptance - Windows Driver Update Profile Automatic"
  approval_type = "automatic"
  description   = "Test description for automatic approval driver update profile"
  
  deployment_deferral_in_days = 7
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

  lifecycle {
    ignore_changes = [
      role_scope_tag_ids
    ]
  }
}
