# Use OData startsWith function to find policies by name prefix
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "prod_policies" {
  provider = microsoft365
  config {
    odata_filter = "startsWith(name, '[Base] Prod')"
  }
}

