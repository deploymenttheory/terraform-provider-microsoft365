# Scenario 2: Maximal resource configuration
# Tests the resource with all optional fields populated
resource "microsoft365_graph_beta_device_management_macos_platform_script" "maximal" {
  display_name                  = "acc-test-maximal-macos-script"
  description                   = "Comprehensive macOS platform script with all features enabled for unit testing"
  file_name                     = "maximal_test.sh"
  script_content                = "#!/bin/bash\necho 'Max Test'\nexit 0"
  run_as_account                = "user"
  role_scope_tag_ids            = ["0", "1", "2"]
  block_execution_notifications = true
  execution_frequency           = "P1D"
  retry_count                   = 3

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
