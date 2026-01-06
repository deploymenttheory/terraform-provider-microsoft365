# Use OData OR operator to match multiple specific policies
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "kerberos_or_licensing" {
  provider = microsoft365
  config {
    odata_filter = "name eq '[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0' or name eq '[Base] Prod | Windows - Settings Catalog | Licensing ver1.0'"
  }
}

