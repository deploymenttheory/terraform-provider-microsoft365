resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  custom_attribute_type = "string"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}