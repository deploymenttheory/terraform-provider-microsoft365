resource "microsoft365_graph_beta_device_management_windows_remediation_script" "exclusion_assignment" {
  display_name               = "Test Exclusion Assignment Windows Remediation Script - Unique"
  description                = ""
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with exclusion assignment\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with exclusion assignment\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}