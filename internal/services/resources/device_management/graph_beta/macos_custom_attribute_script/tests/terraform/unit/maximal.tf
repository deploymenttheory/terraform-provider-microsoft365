resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_maximal" {
  display_name          = "Test Maximal macOS Custom Attribute Script - Unit"
  description           = "Comprehensive test configuration with maximal settings"
  custom_attribute_type = "string"
  file_name             = "test_maximal.sh"
  script_content        = "#!/bin/bash\necho 'Maximal Test Value'\ndate\necho $USER\nexit 0"
  run_as_account        = "user"
  role_scope_tag_ids    = ["0", "1"]
}