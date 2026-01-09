# Scenario 3: Update from minimal to maximal (Step 1 - Initial minimal state)
# Tests updating a resource from minimal configuration to maximal
resource "microsoft365_graph_beta_device_management_macos_platform_script" "update_test" {
  display_name   = "acc-test-update-test-script"
  file_name      = "update_test.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
