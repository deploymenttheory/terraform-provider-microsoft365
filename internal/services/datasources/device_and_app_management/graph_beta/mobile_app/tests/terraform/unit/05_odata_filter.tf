data "microsoft365_graph_beta_device_and_app_management_mobile_app" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "startswith(publisher, 'Microsoft')"
  odata_top    = 10
}

