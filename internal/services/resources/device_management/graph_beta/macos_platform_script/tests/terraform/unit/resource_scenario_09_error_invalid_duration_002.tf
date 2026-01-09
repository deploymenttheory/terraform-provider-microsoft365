# Scenario 9: Error case - Invalid execution_frequency format
resource "microsoft365_graph_beta_device_management_macos_platform_script" "error_invalid_duration" {
  display_name        = "unit-test-error-invalid-duration"
  file_name           = "error_test.sh"
  script_content      = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account      = "system"
  execution_frequency = "INVALID_DURATION"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
