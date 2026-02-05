resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_1" {
  display_name     = "acc-test-group-007-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows feature update policy assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_2" {
  display_name     = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows feature update policy assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_3" {
  display_name     = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows feature update policy exclusion assignments"
  hard_delete      = true
}

# ==============================================================================
# Time Wait
# ==============================================================================

resource "time_sleep" "wait_15_seconds" {
  create_duration = "15s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_007_1,
    microsoft365_graph_beta_groups_group.acc_test_group_007_2,
    microsoft365_graph_beta_groups_group.acc_test_group_007_3
  ]
}

# ==============================================================================
# Test Case: Assignments Maximal
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_006" {
  depends_on = [time_sleep.wait_15_seconds]
  display_name                                            = "acc-test-007-assignments-maximal-${random_string.test_suffix.result}"
  description                                             = "Maximal test with multiple assignments"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  role_scope_tag_ids = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_3.id
    }
  ]
}
