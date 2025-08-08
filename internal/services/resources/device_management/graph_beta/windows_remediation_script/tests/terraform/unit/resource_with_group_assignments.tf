resource "microsoft365_graph_beta_device_management_windows_remediation_script" "group_assignments" {
  display_name                = "Test Group Assignments Windows Remediation Script - Unique"
  description                 = ""
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with group assignments\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with group assignments\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
      daily_schedule = {
        interval = 1
        time     = "09:00:00"
        use_utc  = false
      }
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
      daily_schedule = {
        interval = 1
        time     = "15:00:00"
        use_utc  = false
      }
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}