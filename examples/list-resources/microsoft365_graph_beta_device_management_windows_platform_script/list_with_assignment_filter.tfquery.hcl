provider "microsoft365" {}

# List assigned Windows platform scripts running as system
list "microsoft365_graph_beta_device_management_windows_platform_script" "assigned_system" {
  provider = microsoft365
  config {
    run_as_account_filter = "system"
    is_assigned_filter    = true
  }
}
