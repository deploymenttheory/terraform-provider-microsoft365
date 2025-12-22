resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_006" {
  display_name            = "unit-test-windows-platform-script-006-assignments-maximal"
  description             = "Maximal test with multiple assignments"
  file_name               = "test-script-006.ps1"
  script_content          = "Write-Host 'Maximal script'"
  run_as_account          = "user"
  enforce_signature_check = true
  run_as_32_bit           = true
  role_scope_tag_ids      = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

