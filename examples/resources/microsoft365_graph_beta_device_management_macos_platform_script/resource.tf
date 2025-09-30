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

  # Optional: Assignments block
  assignments = [
    # Optional: inclusion group assignments
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]

  # Timeouts configuration (optional)
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}