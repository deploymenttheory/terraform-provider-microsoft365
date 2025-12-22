resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_005" {
  display_name   = "unit-test-windows-platform-script-005-assignments-minimal"
  file_name      = "test-script-005.ps1"
  script_content = "Write-Host 'Hello World'"
  run_as_account = "system"

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

