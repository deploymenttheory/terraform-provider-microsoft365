data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_packages" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "productDisplayName eq '7-Zip'"
  odata_count   = true
  odata_orderby = "productDisplayName"
  odata_search  = "\"productDisplayName:7-Zip\""
  odata_select  = "id,productId,productDisplayName,publisherDisplayName,versionDisplayName"
  odata_top     = 5
  odata_skip    = 0
}
