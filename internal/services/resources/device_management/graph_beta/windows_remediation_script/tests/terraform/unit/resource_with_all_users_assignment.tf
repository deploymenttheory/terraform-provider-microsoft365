resource "microsoft365_graph_beta_device_management_windows_remediation_script" "all_users_assignment" {
  display_name                = "Test All Users Assignment Windows Remediation Script - Unique"
  description                 = ""
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with all licensed users assignment\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with all licensed users assignment\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
      daily_schedule = {
        interval = 1
        time     = "12:00:00"
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