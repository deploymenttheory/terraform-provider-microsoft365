resource "random_uuid" "test_assignments" {
}

resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "test_assignments" {
  display_name  = "Acceptance - Windows Driver Update Profile with Assignments"
  approval_type = "manual"
  description   = "Test description for driver update profile with assignments"

  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    }
  ]

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
