resource "microsoft365_graph_beta_device_and_app_management_linux_platform_script" "example" {
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

  # Optional: Custom timeouts
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }

  # Optional: Assignments
  assignments = {
    all_devices = false

    all_users = false

    include_groups = [
      {
        group_id                   = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
        include_groups_filter_type = "none"
      },
      {
        group_id                   = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
        include_groups_filter_type = "none"
      },
    ]

    exclude_group_ids = [
      "b8c661c2-fa9a-4351-af86-adc1729c343f",
      "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f",
    ]
  }
}