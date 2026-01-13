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
  description      = "Test group 1 for windows feature update policy assignments lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_2" {
  display_name     = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows feature update policy assignments lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_3" {
  display_name     = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows feature update policy exclusion assignments lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_007" {
  display_name                                            = "acc-test-windows-feature-update-policy-007-${random_string.test_suffix.result}"
  feature_update_version                                  = "Windows 11, version 23H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

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
