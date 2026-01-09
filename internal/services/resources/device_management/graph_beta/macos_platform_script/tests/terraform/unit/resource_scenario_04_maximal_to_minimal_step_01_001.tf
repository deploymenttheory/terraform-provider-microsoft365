# Scenario 4 Step 1: Maximal configuration (before downgrade)
resource "microsoft365_graph_beta_device_management_macos_platform_script" "downgrade_test" {
  display_name                  = "unit-test-downgrade-test-script"
  description                   = "Initial maximal configuration for downgrade testing"
  file_name                     = "downgrade_test.sh"
  script_content                = "#!/bin/bash\necho 'Max Test'\nexit 0"
  run_as_account                = "user"
  role_scope_tag_ids            = ["0", "1", "2"]
  block_execution_notifications = true
  execution_frequency           = "P1D"
  retry_count                   = 5

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
