resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "exclusion_assignment" {
  display_name  = "Test Exclusion Assignment Windows Driver Update Profile - Unique"
  approval_type = "manual"
  description   = "Test description for driver update profile with exclusion assignment"

  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000003"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
