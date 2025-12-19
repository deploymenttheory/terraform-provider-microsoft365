resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_008" {
  display_name               = "unit-test-windows-remediation-script-008-assignments-downgrade"
  description                = "Scenario 8 Step 2: Downgraded to minimal assignments"
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

