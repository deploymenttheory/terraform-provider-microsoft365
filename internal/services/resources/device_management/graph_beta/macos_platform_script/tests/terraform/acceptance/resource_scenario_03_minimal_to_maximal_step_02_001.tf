# Scenario 3 Step 2: Update from minimal to maximal configuration
resource "microsoft365_graph_beta_device_management_macos_platform_script" "update_test" {
  display_name                  = "acc-test-update-test-script-updated"
  description                   = "Updated to maximal configuration"
  file_name                     = "update_test_maximal.sh"
  script_content                = "#!/bin/bash\necho 'Max Test'\nexit 0"
  run_as_account                = "user"
  role_scope_tag_ids            = ["0", "1"]
  block_execution_notifications = true
  execution_frequency           = "PT12H"
  retry_count                   = 2

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
