# List assigned Edge policies (combining name filter with assignment check)
# This is efficient: name filter reduces results first, then assignment check runs on fewer policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "assigned_edge" {
  provider = microsoft365
  config {
    name_filter        = "Edge"
    is_assigned_filter = true
  }
}

