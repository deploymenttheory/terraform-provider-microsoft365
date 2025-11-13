data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_advanced" {
  filter_type  = "odata"
  odata_select = "id,productId,productDisplayName,publisherDisplayName"
  odata_top    = 10
  odata_skip   = 0
}
