provider "microsoft365" {}

# List Windows platform scripts with combined filters
list "microsoft365_graph_beta_device_management_windows_platform_script" "filtered" {
  provider = microsoft365
  config {
    display_name_filter    = "Script"
    run_as_account_filter  = "system"
  }
}
