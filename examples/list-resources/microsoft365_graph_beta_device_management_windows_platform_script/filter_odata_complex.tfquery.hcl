provider "microsoft365" {}

# List Windows platform scripts with complex OData filter
list "microsoft365_graph_beta_device_management_windows_platform_script" "odata_complex" {
  provider = microsoft365
  config {
    odata_filter = "runAsAccount eq 'system' and contains(displayName, 'Baseline') and contains(fileName, 'ps1')"
  }
}
