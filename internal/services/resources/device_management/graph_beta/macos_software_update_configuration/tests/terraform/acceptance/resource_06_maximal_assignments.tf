# ==============================================================================
# Test 06: Minimal Resource with Maximal Assignments
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Test groups for assignments
resource "microsoft365_graph_beta_groups_group" "test_group_06_include" {
  display_name     = "acc-test-06-group-include-${random_string.suffix.result}"
  mail_nickname    = "acc-test-06-group-include-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test include group for macOS software update configuration assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "test_group_06_exclude" {
  display_name     = "acc-test-06-group-exclude-${random_string.suffix.result}"
  mail_nickname    = "acc-test-06-group-exclude-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test exclude group for macOS software update configuration assignments"
  hard_delete      = true
}

# Wait for group replication
resource "time_sleep" "test_06_wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.test_group_06_include,
    microsoft365_graph_beta_groups_group.test_group_06_exclude
  ]
  create_duration = "10s"
}

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_06_max_assignments" {
  display_name                             = "acc-test-06-max-assignments-${random_string.suffix.result}"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_06_include.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_06_exclude.id
    }
  ]

  depends_on = [time_sleep.test_06_wait_for_groups]
}
