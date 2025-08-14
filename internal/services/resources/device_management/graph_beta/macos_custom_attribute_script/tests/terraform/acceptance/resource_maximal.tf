resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Acceptance macOS Custom Attribute Script - Updated"
  description           = "Updated description for acceptance testing"
  custom_attribute_type = "string"
  file_name             = "test_acceptance_updated.sh"
  script_content        = "#!/bin/bash\necho 'Updated acceptance test value'\ndate\nexit 0"
  run_as_account        = "user"
  role_scope_tag_ids    = ["0", "1"]
}