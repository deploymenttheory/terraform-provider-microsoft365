resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Acceptance macOS Custom Attribute Script"
  description           = ""
  custom_attribute_type = "string"
  file_name             = "test_acceptance.sh"
  script_content        = "#!/bin/bash\necho 'Acceptance test value'\nexit 0"
  run_as_account        = "system"
}