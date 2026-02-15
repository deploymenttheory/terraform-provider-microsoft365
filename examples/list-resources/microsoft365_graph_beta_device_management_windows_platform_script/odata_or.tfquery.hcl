provider "microsoft365" {}

# Use OData filter with OR logic
list "microsoft365_graph_beta_device_management_windows_platform_script" "odata_or" {
  provider = microsoft365
  config {
    odata_filter = "contains(displayName, 'Baseline') or contains(displayName, 'Security')"
  }
}
