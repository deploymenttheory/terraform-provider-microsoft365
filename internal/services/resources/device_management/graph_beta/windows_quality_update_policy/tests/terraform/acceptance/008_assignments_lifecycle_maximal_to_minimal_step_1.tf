
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_1" {
  display_name     = "acc-test-group-008-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-008-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows quality update policy lifecycle assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_2" {
  display_name     = "acc-test-group-008-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-008-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows quality update policy lifecycle assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_3" {
  display_name     = "acc-test-group-008-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-008-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows quality update policy lifecycle exclusion assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_after_groups" {
  create_duration = "15s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_008_1,
    microsoft365_graph_beta_groups_group.acc_test_group_008_2,
    microsoft365_graph_beta_groups_group.acc_test_group_008_3,
  ]
}

# ==============================================================================
# Windows Quality Update Policy Resource - Step 1: Maximal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_008" {
  display_name     = "acc-test-windows-quality-update-policy-008-lifecycle-${random_string.test_suffix.result}"
  hotpatch_enabled = false

  depends_on = [
    time_sleep.wait_after_groups,
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_008_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_008_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_008_3.id
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

