# List only policies that do not have assignments
# Note: This checks actual assignments via API calls and may take 20-30 seconds for large tenants
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "unassigned_only" {
  provider = microsoft365
  config {
    is_assigned_filter = false
  }
}

