resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_003" {
  display_name               = "unit-test-windows-remediation-script-003-lifecycle"
  description                = "Lifecycle Step 1: Starting with minimal configuration"
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

