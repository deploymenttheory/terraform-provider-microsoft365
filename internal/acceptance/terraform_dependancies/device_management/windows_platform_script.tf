resource "random_string" "windows_platform_script_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "acc_test_windows_platform_script_01" {
  display_name       = "acc-test-windows-platform-script-${random_string.windows_platform_script_suffix.result}"
  description        = "acc-test-windows-platform-script"
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

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}


resource "microsoft365_graph_beta_device_management_windows_platform_script" "acc_test_windows_platform_script_02" {
  display_name       = "acc-test-windows-platform-script-${random_string.windows_platform_script_suffix.result}"
  description        = "acc-test-windows-platform-script"
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

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}