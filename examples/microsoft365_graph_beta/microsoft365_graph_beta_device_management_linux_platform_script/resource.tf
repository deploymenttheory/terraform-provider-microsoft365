resource "microsoft365_graph_beta_device_management_linux_platform_script" "example" {
  name                = "Example Linux Script"
  description         = "Example script to demonstrate Linux platform script configuration"
  execution_context   = "user" // Possible values are user, root
  execution_frequency = "1day" // Possible values are 1day, 3days, 5days, 1week
  execution_retries   = 2      // Can be one of: `15minutes`, `30minutes`, `1hour`, `2hour`, `3hour`, `6hour`, `12hour`, `1day`, or `1week`. Defaults to `15minutes`.

  script_content = <<-EOT
    #!/bin/bash
    echo "Hello from Linux script"
    # Add your script content here 2
  EOT

  role_scope_tag_ids = ["0"]

  assignments = [
    # Optional: Assignment targeting all devices with inlcude filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
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

  # Optional: Custom timeouts
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}