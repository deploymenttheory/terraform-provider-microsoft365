# Creates three published Win32 apps with a dependency (main -> runtime) and a
# supersedence (legacy -> main), then looks the relationships up every supported way.
# The main app's collection therefore holds an outgoing dependency and an incoming supersedence.

# The acceptance precheck downloads and packages the pinned Mozilla installer.
variable "installer_zip" {
  type = string
}

resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Relationships can only be added between published apps, so content is uploaded.
# Apps upload one at a time; concurrent uploads exhaust the per-block upload timeout.
resource "microsoft365_graph_beta_device_and_app_management_win32_app" "main" {

  display_name                      = "acc-test-app-relationship-ds-main-${random_string.test_suffix.result}"
  description                       = "Win32 app for the mobile app relationship data source acceptance test"
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

resource "microsoft365_graph_beta_device_and_app_management_win32_app" "runtime" {
  depends_on = [microsoft365_graph_beta_device_and_app_management_win32_app.main]


  display_name                      = "acc-test-app-relationship-ds-runtime-${random_string.test_suffix.result}"
  description                       = "Win32 app for the mobile app relationship data source acceptance test"
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

resource "microsoft365_graph_beta_device_and_app_management_win32_app" "legacy" {
  depends_on = [microsoft365_graph_beta_device_and_app_management_win32_app.runtime]


  display_name                      = "acc-test-app-relationship-ds-legacy-${random_string.test_suffix.result}"
  description                       = "Win32 app for the mobile app relationship data source acceptance test"
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
  depends_on      = [microsoft365_graph_beta_device_and_app_management_win32_app.legacy]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "test" {
  depends_on = [time_sleep.wait_for_apps]

  source_id       = microsoft365_graph_beta_device_and_app_management_win32_app.main.id
  target_id       = microsoft365_graph_beta_device_and_app_management_win32_app.runtime.id
  dependency_type = "autoInstall"
}

# The legacy app supersedes the main app, so the main app sees it as a parent relationship.
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "test" {
  depends_on = [microsoft365_graph_beta_device_and_app_management_mobile_app_dependency.test]

  source_id         = microsoft365_graph_beta_device_and_app_management_win32_app.legacy.id
  target_id         = microsoft365_graph_beta_device_and_app_management_win32_app.main.id
  supersedence_type = "update"
}

locals {
  main_app_id    = microsoft365_graph_beta_device_and_app_management_win32_app.main.id
  runtime_app_id = microsoft365_graph_beta_device_and_app_management_win32_app.runtime.id
}

data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_relationship_id" {
  relationship_id = microsoft365_graph_beta_device_and_app_management_mobile_app_dependency.test.id

  depends_on = [microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.test]
}

data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_id" {
  source_id = local.main_app_id
  target_id = microsoft365_graph_beta_device_and_app_management_win32_app.legacy.id

  depends_on = [microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.test]
}

data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_display_name" {
  source_id           = local.main_app_id
  target_display_name = microsoft365_graph_beta_device_and_app_management_win32_app.runtime.display_name

  depends_on = [microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.test]
}

# The main app holds two relationships; only Graph's $filter narrows this to one.
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_odata_query" {
  source_id   = local.main_app_id
  odata_query = "isof('microsoft.graph.mobileAppDependency')"

  depends_on = [microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.test]
}

# From the runtime app's side the dependency is incoming, so its target is the parent.
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "parent_by_odata_query" {
  source_id   = local.runtime_app_id
  odata_query = "targetType eq 'parent'"

  depends_on = [microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.test]
}
