provider "microsoft365" {}

# List only Windows platform scripts without assignments
list "microsoft365_graph_beta_device_management_windows_platform_script" "unassigned_only" {
  provider = microsoft365
  config {
    is_assigned_filter = false
  }
}
