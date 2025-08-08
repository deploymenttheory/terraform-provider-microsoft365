resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name               = "Test Acceptance Windows Remediation Script"
  description                = ""
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}