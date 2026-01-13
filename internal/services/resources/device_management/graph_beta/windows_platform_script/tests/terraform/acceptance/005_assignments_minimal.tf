
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
  description      = "Test group 1 for windows platform script assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_005" {
  display_name   = "acc-test-windows-platform-script-005-${random_string.test_suffix.result}"
  file_name      = "test-script-005.ps1"
  script_content = "Write-Host 'Hello World'"
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

