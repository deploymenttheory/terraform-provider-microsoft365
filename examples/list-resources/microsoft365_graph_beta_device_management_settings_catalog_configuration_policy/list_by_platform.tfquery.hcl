# List policies for Windows 10 platform
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "by_platform" {
  provider = microsoft365
  config {
    platform_filter = ["windows10"]
  }
}

