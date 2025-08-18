resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "group_assignments" {
  display_name  = "Test Group Assignments Windows Driver Update Profile - Unique"
  approval_type = "manual"
  description   = "Test description for driver update profile with group assignments"

  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
