// Example: macOS Platform Script Resource

resource "microsoft365_graph_beta_device_and_app_management_macos_platform_script" "example" {
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
      "51a96cdd-4b9b-4849-b416-8c94a6d88797",
      "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
    ]

    exclude_group_ids = [
      "b8c661c2-fa9a-4351-af86-adc1729c343f",
      "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
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