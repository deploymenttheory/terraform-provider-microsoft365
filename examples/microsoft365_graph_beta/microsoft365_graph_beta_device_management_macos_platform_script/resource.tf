// Example: Device Shell Script Resource

resource "microsoft365_graph_beta_device_management_macos_platform_script" "example" {
  # Required fields
  display_name = "MacOS Shell Script"
  description  = "Example shell script for MacOS devices"

  script_content = <<EOT
    #!/bin/bash
    echo "Hello World"
  EOT

  run_as_account = "system" # Possible values: "system" or "user"
  file_name      = "example_script.sh"

  # Optional fields
  block_execution_notifications = false
  execution_frequency           = "P1D" # ISO 8601 duration format (e.g., P1D for 1 day, PT1H for 1 hour)
  retry_count                   = 3

  # Role scope tag IDs (optional)
  role_scope_tag_ids = ["0"]

  # Script assignments (optional)
  assignments = {
    all_devices = false
    all_users   = false

    include_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]

    exclude_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  # Timeouts configuration (optional)
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}