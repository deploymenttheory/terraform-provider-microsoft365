resource "microsoft365_graph_beta_device_and_app_management_win32_lob_app" "example" {
  display_name           = "Example App"
  description            = "This is an example Win32 LOB app"
  publisher              = "Example Publisher"
  file_name              = "example_app.msi"
  install_command_line   = "/install /quiet"
  uninstall_command_line = "/uninstall /quiet"

  # Minimum supported operating system
  minimum_supported_operating_system = {
    v10_0 = true  # Windows 10 or later
  }

  # Detection rules - must be a list. Select one of the detection rule examples below
  detection_rules = [
    {
      detection_type                = "registry"
      key_path                      = "HKEY_LOCAL_MACHINE\\SOFTWARE\\ExampleApp"
      value_name                    = "Version"
      check_32_bit_on_64_system     = false
      registry_detection_type       = "version"
      registry_detection_operator   = "greaterThan"
      detection_value               = "1.0"
    }
  ]

  # Alternative detection rules examples (uncomment to use):
  # detection_rules = [
  #   {
  #     detection_type               = "msi_information"
  #     product_code                 = "1234-5678-ABCD-EFGH"
  #     product_version              = "1.0.0"
  #     product_version_operator     = "greaterThanOrEqual"
  #   }
  # ]

  # detection_rules = [
  #   {
  #     detection_type                = "file_system"
  #     file_path                     = "C:\\Program Files\\ExampleApp"
  #     file_or_folder_name           = "app.exe"
  #     filesystem_detection_type     = "exists"
  #     filesystem_detection_operator = "equal"
  #     check_32_bit_on_64_system     = false
  #   }
  # ]

  # detection_rules = [
  #   {
  #     detection_type          = "powershell_script"
  #     script_content          = <<EOT
  # Get-ItemProperty -Path "HKLM:\\Software\\ExampleApp" -Name "Version" | Select-Object -ExpandProperty Version
  # EOT
  #     enforce_signature_check = false
  #     run_as_32_bit           = false
  #   }
  # ]

  # Requirement rules - must be a list
  requirement_rules = [
    {
      requirement_type    = "file"
      key_path            = "C:\\Program Files\\ExampleApp\\app.exe"
      file_or_folder_name = "app.exe"
      operator            = "equal"
      detection_value     = "exists"
    }
  ]

  # MSI Information - single nested attribute (use = not block)
  msi_information = {
    product_code    = "1234-5678-ABCD-EFGH"
    product_version = "1.0.0"
    upgrade_code    = "9876-5432-HGFE-DCBA"
    requires_reboot = false
    package_type    = "perMachine"
  }

  # Install experience - single nested attribute (use = not block)
  install_experience = {
    run_as_account          = "system"
    device_restart_behavior = "suppress"
  }

  # Return codes - must be a list
  return_codes = [
    {
      return_code = 0
      type        = "success"
    },
    {
      return_code = 3010
      type        = "softReboot"
    }
  ]

  # Timeouts
  timeouts = {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
