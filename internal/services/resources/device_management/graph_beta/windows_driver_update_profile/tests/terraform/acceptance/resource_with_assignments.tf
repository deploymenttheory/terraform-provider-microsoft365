resource "random_uuid" "test_assignments" {
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-driver-include-${random_uuid.test_assignments.result}"
  mail_nickname    = "acc-driver-include-${random_uuid.test_assignments.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-driver-exclude-${random_uuid.test_assignments.result}"
  mail_nickname    = "acc-driver-exclude-${random_uuid.test_assignments.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2,
  ]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "test_assignments" {
  display_name  = "Acceptance - Windows Driver Update Profile with Assignments"
  approval_type = "manual"
  description   = "Test description for driver update profile with assignments"

  role_scope_tag_ids = ["0"]

  depends_on = [
    time_sleep.wait_for_groups
  ]

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
