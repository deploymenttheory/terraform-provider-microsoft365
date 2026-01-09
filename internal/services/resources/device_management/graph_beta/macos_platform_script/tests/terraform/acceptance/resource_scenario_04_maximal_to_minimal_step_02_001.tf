# Scenario 4: Update from maximal to minimal (Step 2 - Downgraded to minimal state)
# Tests updating a resource from maximal configuration to minimal
resource "microsoft365_graph_beta_device_management_macos_platform_script" "downgrade_test" {
  display_name   = "acc-test-downgrade-test-script-minimal"
  file_name      = "downgrade_test_minimal.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
