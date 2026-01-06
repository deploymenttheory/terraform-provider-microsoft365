provider "microsoft365" {}

# List policies from the baseline template family
list "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "by_template_family" {
  provider = microsoft365
  config {
    template_family_filter = "baseline"
  }
}
