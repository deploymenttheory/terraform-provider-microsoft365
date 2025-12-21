
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_1" {
  display_name     = "acc-test-group-006-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-006-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows quality update policy assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_2" {
  display_name     = "acc-test-group-006-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-006-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows quality update policy assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_3" {
  display_name     = "acc-test-group-006-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-006-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows quality update policy exclusion assignments"
  hard_delete      = true
}

# ==============================================================================
# Windows Quality Update Policy Resource with Maximal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_006" {
  display_name     = "acc-test-windows-quality-update-policy-006-assignments-maximal-${random_string.test_suffix.result}"
  hotpatch_enabled = false

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_006_1,
    microsoft365_graph_beta_groups_group.acc_test_group_006_2,
    microsoft365_graph_beta_groups_group.acc_test_group_006_3
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_3.id
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

