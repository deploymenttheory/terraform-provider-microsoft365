resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# The Go acceptance precheck downloads Firefox 140.0.4 from Mozilla over HTTPS,
# packages it on the runner, and substitutes the loopback HTTP package source.
resource "microsoft365_graph_beta_device_and_app_management_win32_app" "test_002" {
  display_name                      = "acc-test-win32-app-002-${random_string.test_suffix.result}"
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

  display_version               = "140.0.4"
  is_featured                   = true
  privacy_information_url       = "https://www.mozilla.org/privacy/firefox/"
  information_url               = "https://www.mozilla.org/firefox/"
  owner                         = "Acceptance Tests"
  developer                     = "Mozilla"
  notes                         = "Maximal Win32 EXE configuration"
  minimum_free_disk_space_in_mb = 1024
  minimum_memory_in_mb          = 2048
  minimum_number_of_processors  = 2
  minimum_cpu_speed_in_mhz      = 1000
  role_scope_tag_ids            = ["0"]

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
    { return_code = 0, type = "success" },
    { return_code = 1707, type = "success" },
    { return_code = 3010, type = "softReboot" },
    { return_code = 1641, type = "hardReboot" },
    { return_code = 1618, type = "retry" }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "5m"
  }
}
