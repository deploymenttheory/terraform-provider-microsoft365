resource "microsoft365_graph_beta_device_management_windows_platform_script" "with_assignments" {
  display_name  = "Test Windows Platform Script with Assignments - Unique"
  description   = "Test description for script with assignments"
  file_name     = "test-assignments-script.ps1"
  script_content = "Write-Host 'Script with assignments'"
  run_as_account = "system"
  
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}