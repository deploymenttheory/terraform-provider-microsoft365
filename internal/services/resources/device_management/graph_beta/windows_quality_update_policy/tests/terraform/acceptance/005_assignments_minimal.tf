
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependency
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_005_1" {
  display_name     = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows quality update policy minimal assignments"
  hard_delete      = true
}

# ==============================================================================
# Windows Quality Update Policy Resource with Minimal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_005" {
  display_name     = "acc-test-windows-quality-update-policy-005-assignments-minimal-${random_string.test_suffix.result}"
  hotpatch_enabled = false

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_005_1
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_005_1.id
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

