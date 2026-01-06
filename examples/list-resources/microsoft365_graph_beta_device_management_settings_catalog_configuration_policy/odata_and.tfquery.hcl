# Use OData AND operator to combine multiple conditions
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "defender_windows10" {
  provider = microsoft365
  config {
    odata_filter = "contains(name, 'Defender') and platforms eq 'windows10'"
  }
}

