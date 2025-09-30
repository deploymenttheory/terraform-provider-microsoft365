resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "example" {
  custom_attribute_name = "ExampleCustomAttribute"
  custom_attribute_type = "string"
  display_name          = "Example macOS Custom Attribute Script"
  description           = "Example description for custom attribute script."
  script_content        = "#!/bin/bash\necho 'Hello World'"
  run_as_account        = "system"
  file_name             = "example-script.sh"

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

  timeouts = {
    create = "30m"
    update = "30m"
    read   = "30m"
    delete = "30m"
  }
} 