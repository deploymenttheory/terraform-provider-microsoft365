# Scenario 9: Error case - Invalid run_as_account value
resource "microsoft365_graph_beta_device_management_macos_platform_script" "error_invalid_run_as" {
  display_name   = "unit-test-error-invalid-run-as"
  file_name      = "error_test.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "invalid_value"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
