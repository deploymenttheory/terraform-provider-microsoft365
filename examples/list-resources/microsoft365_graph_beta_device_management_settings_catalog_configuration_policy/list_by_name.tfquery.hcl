# List policies with "Kerberos" in the name
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "by_name" {
  provider = microsoft365
  config {
    name_filter = "Kerberos"
  }
}

