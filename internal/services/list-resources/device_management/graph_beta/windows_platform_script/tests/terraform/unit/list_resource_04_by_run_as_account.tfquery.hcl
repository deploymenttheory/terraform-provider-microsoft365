provider "microsoft365" {}

# List Windows platform scripts filtered by run as account
list "microsoft365_graph_beta_device_management_windows_platform_script" "filtered" {
  provider = microsoft365
  config {
    run_as_account_filter = "system"
  }
}
