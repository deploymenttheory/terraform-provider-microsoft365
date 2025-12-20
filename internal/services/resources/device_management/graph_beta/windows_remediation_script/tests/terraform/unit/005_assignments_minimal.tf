resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_005" {
  display_name               = "unit-test-windows-remediation-script-005-assignments-minimal"
  description                = "Scenario 5: Minimal assignments"
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

