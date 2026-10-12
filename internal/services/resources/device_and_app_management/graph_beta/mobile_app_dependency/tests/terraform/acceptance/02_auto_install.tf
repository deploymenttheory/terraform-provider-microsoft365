# The acceptance precheck downloads and packages the pinned Mozilla installer.
variable "installer_zip" {
  type = string
}

resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Terraform uploads and publishes each Win32 app before creating dependencies.
resource "microsoft365_graph_beta_device_and_app_management_win32_app" "apps" {
  for_each = toset(["main", "runtime"])

  display_name                      = "acc-test-app-dependency-${each.key}-${random_string.test_suffix.result}"
  description                       = "Win32 app for the mobile app dependency acceptance test"
  publisher                         = "Mozilla"
  allow_available_uninstall         = true
  display_version                   = "140.0.2"
  file_name                         = "firefox-140.0.2.intunewin"
  setup_file_path                   = "setup.exe"
  install_command_line              = "setup.exe /S"
  uninstall_command_line            = "\"C:\\Program Files\\Mozilla Firefox\\uninstall\\helper.exe\" /S"
  minimum_supported_windows_release = "Windows10_22H2"
  allowed_architectures             = ["x64"]

  app_installer_zip = {
    installer_file_path_source = var.installer_zip
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
    comparison_value           = "140.0.2"
  }]

  return_codes = [{ return_code = 0, type = "success" }]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "5m"
  }
}

# Allow the uploaded apps to propagate before adding their relationships.
resource "time_sleep" "wait_for_apps" {
  depends_on      = [microsoft365_graph_beta_device_and_app_management_win32_app.apps]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "test" {
  depends_on = [time_sleep.wait_for_apps]

  source_id       = microsoft365_graph_beta_device_and_app_management_win32_app.apps["main"].id
  target_id       = microsoft365_graph_beta_device_and_app_management_win32_app.apps["runtime"].id
  dependency_type = "autoInstall"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
