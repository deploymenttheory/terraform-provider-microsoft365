// Example: Device Management Script Resource

resource "microsoft365_graph_beta_device_and_app_management_device_management_script" "example" {
  display_name       = "Example Device Management Script 1"
  description        = "This is an example device management script"
  role_scope_tag_ids = ["0"]

  # The actual script content (should be PowerShell script content only)
  script_content = <<EOT
    # Your PowerShell script here
    Write-Host "Hello from device management script!"
    # Add more PowerShell commands as needed
  EOT

  run_as_account         = "system" # Can be "system" or "user"
  enforce_signature_check = false
  file_name              = "example_script.ps1"
  run_as_32_bit          = false

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

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}