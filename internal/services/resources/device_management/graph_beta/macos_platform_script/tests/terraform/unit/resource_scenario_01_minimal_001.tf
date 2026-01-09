# Scenario 1: Minimal resource configuration
# Tests the resource with only required fields
resource "microsoft365_graph_beta_device_management_macos_platform_script" "minimal" {
  display_name   = "unit-test-minimal-macos-script"
  file_name      = "minimal_test.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
