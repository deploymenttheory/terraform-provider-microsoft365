// Example: Device Management Script Resource

resource "microsoft365_graph_beta_device_management_windows_platform_script" "example" {
  display_name       = "Example Device Management Script"
  description        = "This is an example device management script"
  role_scope_tag_ids = ["0"]

  # The actual script content (should be PowerShell script content only)
  script_content = <<EOT
    # Your PowerShell script here
    Write-Host "Hello from device management script!"
    # Add more PowerShell commands as needed
  EOT

  run_as_account          = "system" # Can be "system" or "user"
  enforce_signature_check = false
  file_name               = "example_script.ps1"
  run_as_32_bit           = false

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
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}