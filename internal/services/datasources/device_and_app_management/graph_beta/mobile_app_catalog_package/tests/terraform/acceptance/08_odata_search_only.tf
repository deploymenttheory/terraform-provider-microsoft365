data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_packages" "odata_search_only" {
  filter_type  = "odata"
  odata_search = "\"productDisplayName:Microsoft\""
}
