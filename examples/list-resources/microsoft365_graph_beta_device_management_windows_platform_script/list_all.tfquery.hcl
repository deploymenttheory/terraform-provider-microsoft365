provider "microsoft365" {}

# List all Windows platform scripts
list "microsoft365_graph_beta_device_management_windows_platform_script" "all" {
  provider = microsoft365
  config {}
}
