data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "startswith(publisherDisplayName, 'Microsoft')"
  odata_top    = 10
}
