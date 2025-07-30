resource "microsoft365_graph_beta_device_management_macos_platform_script" "maximal" {
  display_name    = "Maximal macOS Script"
  description     = "This is a comprehensive script with all fields populated"
  script_content  = "#!/bin/bash\necho 'Maximal Script Configuration'"
  run_as_account  = "user"
  file_name       = "maximal-script.sh"
  block_execution_notifications = true
  execution_frequency = "P4W"
  retry_count     = 10
  role_scope_tag_ids = ["0", "1"]

  assignments = {
    all_devices = true
    all_users   = false
  }
} 