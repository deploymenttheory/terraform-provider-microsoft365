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
  description      = "Test group 1 for windows feature update policy maximal to minimal lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_2" {
  display_name     = "acc-test-group-008-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-008-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows feature update policy maximal to minimal lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_3" {
  display_name     = "acc-test-group-008-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-008-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows feature update policy maximal to minimal exclusion"
  hard_delete      = true
}

resource "time_sleep" "wait_15_seconds" {
  create_duration = "15s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_008_1,
    microsoft365_graph_beta_groups_group.acc_test_group_008_2,
    microsoft365_graph_beta_groups_group.acc_test_group_008_3
  ]
}

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_008" {
  depends_on                                              = [time_sleep.wait_15_seconds]
  display_name                                            = "acc-test-009-assignments-lifecycle-maximal-to-minimal-${random_string.test_suffix.result}"
  description                                             = "Maximal assignments lifecycle test"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  role_scope_tag_ids = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_008_1.id
    }
  ]
}
