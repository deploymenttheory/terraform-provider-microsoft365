resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_004" {
  display_name               = "unit-test-windows-remediation-script-004-downgrade"
  description                = "Downgrade Step 2: Downgraded to minimal configuration"
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

