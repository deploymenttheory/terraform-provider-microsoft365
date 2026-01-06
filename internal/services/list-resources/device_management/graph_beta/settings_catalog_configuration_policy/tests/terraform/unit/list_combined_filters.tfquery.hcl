provider "microsoft365" {}

# List policies using multiple filters combined (AND logic)
# This example finds Windows 10 policies with "Defender" in the name
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "combined" {
  provider = microsoft365
  config {
    name_filter     = "Defender"
    platform_filter = ["windows10"]
  }
}

