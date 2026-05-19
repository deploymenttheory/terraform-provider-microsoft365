resource "microsoft365_graph_beta_device_and_app_management_win32_app" "notepad_plus_plus" {
  allow_available_uninstall = true

  app_installer = {
    installer_file_path_source = "/path/to/notepad++_8.9.5.exe_psadt.intunewin"
  }

  app_icon = {
    icon_url_source = "https://upload.wikimedia.org/wikipedia/commons/f/f5/Notepad_plus_plus.png"
  }

  categories = [
    "Business",
    "Productivity",
  ]

  description     = "Notepad++ v8.9.5 x64 - Free source code editor"
  publisher       = "Don Ho"
  developer       = "Don Ho"
  display_name    = "Notepad++ v8.9.5"
  display_version = "8.9.5"
  file_name       = "notepad++_8.9.5.exe_psadt.intunewin"
  information_url = "https://notepad-plus-plus.org/"

  owner = "IT"
  notes = "Deployed via PowerShell App Deployment Toolkit (PSADT). msi_information is not required for EXE-based installers."

  allowed_architectures             = ["x64"]
  minimum_supported_windows_release = "Windows10_22H2"

  install_experience = {
    device_restart_behavior = "allow"
    max_run_time_in_minutes = 60
    run_as_account          = "system"
  }

  setup_file_path        = "Invoke-AppDeployToolkit.exe"
  install_command_line   = "Invoke-AppDeployToolkit.exe -DeploymentType Install -DeployMode Silent"
  uninstall_command_line = "Invoke-AppDeployToolkit.exe -DeploymentType Uninstall -DeployMode Silent"

  # msi_information is omitted - EXE-based installers do not require it.
  # Only populate msi_information when deploying MSI-based packages.

  rules = [
    {
      rule_type                  = "detection"
      rule_sub_type              = "file_system"
      path                       = "C:\\Program Files\\Notepad++"
      file_or_folder_name        = "notepad++.exe"
      check_32_bit_on_64_system  = false
      file_system_operation_type = "version"
      lob_app_rule_operator      = "greaterThanOrEqual"
      comparison_value           = "8.9.5"
    },
  ]

  return_codes = [
    {
      return_code = 0
      type        = "success"
    },
    {
      return_code = 1707
      type        = "success"
    },
    {
      return_code = 3010
      type        = "softReboot"
    },
    {
      return_code = 1641
      type        = "hardReboot"
    },
    {
      return_code = 1618
      type        = "retry"
    },
  ]
}
