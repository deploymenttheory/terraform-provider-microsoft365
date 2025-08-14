resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_type = "string"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}