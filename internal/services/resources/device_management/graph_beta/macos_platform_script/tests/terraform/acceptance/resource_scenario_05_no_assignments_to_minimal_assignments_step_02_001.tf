# Dependancies

resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_005_1" {
  display_name     = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for macOS platform script assignment lifecycle"
  hard_delete      = true

}

# Scenario 5: Step 2 - Resource with single assignment (base resource unchanged)
resource "microsoft365_graph_beta_device_management_macos_platform_script" "add_minimal_assignment" {
  display_name   = "acc-test-add-minimal-assignment"
  file_name      = "test_minimal_assignment.sh"
  script_content = "#!/bin/bash\necho 'Test'\nexit 0"
  run_as_account = "system"

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
