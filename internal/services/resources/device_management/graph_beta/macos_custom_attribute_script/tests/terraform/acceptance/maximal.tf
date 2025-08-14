terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_maximal" {
  display_name          = "Test Maximal macOS Custom Attribute Script - Acceptance"
  description           = "Comprehensive acceptance test configuration with maximal settings"
  custom_attribute_type = "string"
  file_name             = "test_maximal_acceptance.sh"
  script_content        = "#!/bin/bash\necho 'Maximal Acceptance Test Value'\ndate\necho $USER\nexit 0"
  run_as_account        = "user"
  role_scope_tag_ids    = ["0", "1"]
}