provider "microsoft365" {}

# List Windows platform scripts filtered by display name
list "microsoft365_graph_beta_device_management_windows_platform_script" "by_display_name" {
  provider = microsoft365
  config {
    display_name_filter = "Baseline"
  }
}
