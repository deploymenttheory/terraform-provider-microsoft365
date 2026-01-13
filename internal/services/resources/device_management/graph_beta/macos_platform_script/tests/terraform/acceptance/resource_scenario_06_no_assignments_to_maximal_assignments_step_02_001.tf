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
  description      = "Test group 1 for macOS platform script assignment lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_2" {
  display_name     = "acc-test-group-006-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-006-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for macOS platform script assignment lifecycle"
  hard_delete      = true
}

# Scenario 6: Step 2 - Resource with all 4 assignment types
resource "microsoft365_graph_beta_device_management_macos_platform_script" "add_maximal_assignments" {
  display_name   = "acc-test-add-maximal-assignments"
  file_name      = "test_maximal_assignments.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_2.id
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
