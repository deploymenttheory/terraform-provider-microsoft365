provider "microsoft365" {}

# Use OData filter with exact match on display name
list "microsoft365_graph_beta_device_management_windows_platform_script" "exact_match" {
  provider = microsoft365
  config {
    odata_filter = "displayName eq 'Windows Baseline Setup'"
  }
}
