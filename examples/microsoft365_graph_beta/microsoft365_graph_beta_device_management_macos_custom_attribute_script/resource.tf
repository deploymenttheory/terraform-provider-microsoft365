resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "example" {
  custom_attribute_name = "ExampleCustomAttribute"
  custom_attribute_type = "string"
  display_name          = "Example macOS Custom Attribute Script"
  description           = "Example description for custom attribute script."
  script_content        = "#!/bin/bash\necho 'Hello World'"
  run_as_account        = "system"
  file_name             = "example-script.sh"

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

  timeouts = {
    create = "30m"
    update = "30m"
    read   = "30m"
    delete = "30m"
  }
} 