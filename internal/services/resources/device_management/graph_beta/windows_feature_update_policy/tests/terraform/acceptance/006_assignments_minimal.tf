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
  description      = "Test group 1 for windows feature update policy minimal assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_15_seconds" {
  create_duration = "15s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_006_1
  ]
}


resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_005" {
  depends_on = [time_sleep.wait_15_seconds]
  display_name                                            = "acc-test-006-assignments-minimal-${random_string.test_suffix.result}"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_1.id
    }
  ]
}
