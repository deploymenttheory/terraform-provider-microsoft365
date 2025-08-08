resource "microsoft365_graph_beta_device_management_windows_remediation_script" "all_devices_assignment" {
  display_name                = "Test All Devices Assignment Windows Remediation Script - Unique"
  description                 = ""
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with all devices assignment\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with all devices assignment\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
      daily_schedule = {
        interval = 1
        time     = "02:00:00"
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