resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_007" {
  display_name               = "unit-test-windows-remediation-script-007-assignments-lifecycle"
  description                = "Scenario 7 Step 1: Starting with minimal assignments"
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

