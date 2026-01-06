# Use custom OData filter for exact name match
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "exact_match" {
  provider = microsoft365
  config {
    odata_filter = "name eq '[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0'"
  }
}

