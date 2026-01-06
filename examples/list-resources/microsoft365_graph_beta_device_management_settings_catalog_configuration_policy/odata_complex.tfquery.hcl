# Use complex OData query with grouping and mixed AND/OR operators
# This finds Windows 10 policies that contain either "Edge" or "Defender" in the name
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "edge_or_defender_windows10" {
  provider = microsoft365
  config {
    odata_filter = "(contains(name, 'Edge') or contains(name, 'Defender')) and platforms eq 'windows10'"
  }
}

