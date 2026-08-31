resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# The Go acceptance precheck downloads Firefox 140.0.4 from Mozilla over HTTPS,
# packages it on the runner, and substitutes the loopback HTTP package source.
resource "microsoft365_graph_beta_device_and_app_management_win32_app" "test_003" {
  display_name                      = "acc-test-win32-app-003-${random_string.test_suffix.result}"
  description                       = "Mozilla Firefox acceptance test"
  publisher                         = "Mozilla"
  file_name                         = "firefox-140.0.4.intunewin"
  allow_available_uninstall         = true
  minimum_supported_windows_release = "Windows10_22H2"
  allowed_architectures             = ["x64"]
  setup_file_path                   = "setup.exe"
  install_command_line              = "setup.exe /S"
  uninstall_command_line            = "\"C:\\Program Files\\Mozilla Firefox\\uninstall\\helper.exe\" /S"

  app_installer = {
    installer_url_source = "http://win32-packages.test/firefox-140.0.4.intunewin"
  }

  install_experience = {
    run_as_account          = "system"
    device_restart_behavior = "suppress"
    max_run_time_in_minutes = 60
  }

  rules = [{
    rule_type                  = "detection"
    rule_sub_type              = "file_system"
    path                       = "C:\\Program Files\\Mozilla Firefox"
    file_or_folder_name        = "firefox.exe"
    check_32_bit_on_64_system  = false
    file_system_operation_type = "version"
    lob_app_rule_operator      = "greaterThanOrEqual"
    comparison_value           = "140.0.4"
  }]

  return_codes = [
    { return_code = 0, type = "success" }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "5m"
  }
}

# The test owns this empty group; no pre-existing group or device is targeted.
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-win32-group-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-win32-group-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win32_app.test_003.id
  intent        = "available"
  source        = "direct"
  target = {
    target_type = "groupAssignment"
    group_id    = microsoft365_graph_beta_groups_group.test.id
  }
}
