provider "microsoft365" {}

# Use OData filter with AND logic
list "microsoft365_graph_beta_device_management_windows_platform_script" "odata_and" {
  provider = microsoft365
  config {
    odata_filter = "runAsAccount eq 'system' and contains(fileName, '.ps1')"
  }
}
