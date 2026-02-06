# ==============================================================================
# Test 05: Minimal Resource with Minimal Assignments
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Test group for assignments
resource "microsoft365_graph_beta_groups_group" "test_group_05" {
  display_name     = "acc-test-05-group-${random_string.suffix.result}"
  mail_nickname    = "acc-test-05-group-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group for macOS software update configuration assignments"
  hard_delete      = true
}

# Wait for group replication
resource "time_sleep" "test_05_wait_for_groups" {
  depends_on      = [microsoft365_graph_beta_groups_group.test_group_05]
  create_duration = "10s"
}

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_05_min_assignments" {
  display_name                             = "acc-test-05-min-assignments-${random_string.suffix.result}"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_05.id
    }
  ]

  depends_on = [time_sleep.test_05_wait_for_groups]
}
