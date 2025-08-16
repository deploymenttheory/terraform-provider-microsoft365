resource "random_uuid" "test_manual" {
}

resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "test_manual" {
  display_name  = "Acceptance - Windows Driver Update Profile Manual"
  approval_type = "manual"
  description   = "Test description for manual approval driver update profile"
  
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
