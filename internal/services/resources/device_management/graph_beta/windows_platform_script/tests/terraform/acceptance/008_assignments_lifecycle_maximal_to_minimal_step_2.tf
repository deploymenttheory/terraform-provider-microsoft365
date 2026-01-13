
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
  description      = "Test group 1 for windows platform script assignments lifecycle"
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_008" {
  display_name            = "acc-test-windows-platform-script-008-${random_string.test_suffix.result}"
  description             = "Maximal assignments lifecycle test"
  file_name               = "test-script-008.ps1"
  script_content          = "Write-Host 'Maximal script'"
  run_as_account          = "user"
  enforce_signature_check = true
  run_as_32_bit           = true
  role_scope_tag_ids      = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_008_1.id
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

