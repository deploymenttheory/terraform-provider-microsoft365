data "microsoft365_graph_beta_device_and_app_management_mobile_app" "odata_filter" {
  odata_query = "startswith(publisher, 'Microsoft')"
}

