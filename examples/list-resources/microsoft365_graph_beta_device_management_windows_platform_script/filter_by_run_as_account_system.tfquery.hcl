provider "microsoft365" {}

# List Windows platform scripts running as system
list "microsoft365_graph_beta_device_management_windows_platform_script" "system_scripts" {
  provider = microsoft365
  config {
    run_as_account_filter = "system"
  }
}
