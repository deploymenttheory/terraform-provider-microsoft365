resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "minimal" {
  display_name           = "Test Minimal macOS Custom Attribute Script - Unique"
  custom_attribute_type  = "string"
  file_name              = "test_minimal.sh"
  script_content         = "#!/bin/bash\necho 'Test Value'\nexit 0"
  run_as_account         = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}