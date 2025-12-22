resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_008" {
  display_name            = "unit-test-windows-platform-script-008-assignments-lifecycle"
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
      group_id = "11111111-1111-1111-1111-111111111111"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

