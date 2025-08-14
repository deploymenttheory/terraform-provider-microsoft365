resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_lifecycle" {
  display_name          = "Test Lifecycle macOS Custom Attribute Script - Acceptance"
  custom_attribute_type = "string"
  file_name             = "test_lifecycle.sh"
  script_content        = "#!/bin/bash\necho 'Lifecycle Test Value'\nexit 0"
  run_as_account        = "system"
}