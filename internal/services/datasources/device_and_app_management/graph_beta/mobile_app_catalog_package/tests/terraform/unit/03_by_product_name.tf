data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_product_name" {
  filter_type  = "product_name"
  filter_value = "7-Zip"
}
