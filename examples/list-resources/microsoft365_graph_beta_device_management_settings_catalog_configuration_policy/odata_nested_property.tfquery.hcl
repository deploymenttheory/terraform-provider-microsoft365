# Use OData to filter by nested templateReference properties
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "baseline_templates" {
  provider = microsoft365
  config {
    odata_filter = "templateReference/templateFamily eq 'baseline'"
  }
}

