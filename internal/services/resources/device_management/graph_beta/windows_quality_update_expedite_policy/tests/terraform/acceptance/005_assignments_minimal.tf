
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_005_1" {
  display_name     = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows quality update expedite policy assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_after_groups" {
  create_duration = "15s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_005_1,
  ]
}

resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_005" {
  display_name = "acc-test-expedite-policy-005-${random_string.test_suffix.result}"

  expedited_update_settings = {
    quality_update_release   = "2025-12-09T00:00:00Z"
    days_until_forced_reboot = 2
  }

  depends_on = [
    time_sleep.wait_after_groups,
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


