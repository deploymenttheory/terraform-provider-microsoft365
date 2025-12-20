resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_009" {
  display_name               = "unit-test-windows-remediation-script-009-validation"
  description                = "Validation: Invalid run_as_account value"
  publisher                  = "Terraform Provider Test"
  run_as_account             = "invalid"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}

