provider "microsoft365" {}

# List Windows platform scripts filtered by file name
list "microsoft365_graph_beta_device_management_windows_platform_script" "filtered" {
  provider = microsoft365
  config {
    file_name_filter = "baseline_setup.ps1"
  }
}
