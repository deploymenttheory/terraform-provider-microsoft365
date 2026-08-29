resource "microsoft365_graph_beta_device_and_app_management_win32_app" "test" {
  display_name                      = {{NAME}}
  description                       = "Win32 content lifecycle acceptance test"
  publisher                         = "Terraform Provider Test"
  file_name                         = "acceptance.intunewin"
  allow_available_uninstall         = true
  allowed_architectures             = ["x64"]
  minimum_supported_windows_release = "Windows10_22H2"
  setup_file_path                   = "setup.cmd"
  install_command_line              = "cmd.exe /c setup.cmd"
  uninstall_command_line            = "cmd.exe /c rmdir /s /q \"%ProgramData%\\TerraformWin32Acceptance\""
  {{SOURCE_BLOCK}} = {
    installer_file_path_source = {{SOURCE_PATH}}
  }
  install_experience = {
    run_as_account          = "system"
    device_restart_behavior = "suppress"
    max_run_time_in_minutes = 10
  }
  rules = [{
    rule_type                  = "detection"
    rule_sub_type              = "file_system"
    path                       = "C:\\ProgramData\\TerraformWin32Acceptance"
    file_or_folder_name        = "version.txt"
    check_32_bit_on_64_system  = false
    file_system_operation_type = "exists"
  }]
  return_codes = [{ return_code = 0, type = "success" }]
  timeouts     = { create = "15m", update = "15m", read = "5m", delete = "5m" }
}

# Use an empty, dedicated test group. Never target all users or all devices.
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win32_app.test.id
  intent        = "available"
  target = {
    target_type = "groupAssignment"
    group_id    = {{GROUP_ID}}
  }
}
