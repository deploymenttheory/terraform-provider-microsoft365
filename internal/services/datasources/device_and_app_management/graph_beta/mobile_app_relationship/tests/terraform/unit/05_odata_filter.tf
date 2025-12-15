data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "sourceId eq 'app-source-001'"
  timeouts = {
    read = "10s"
  }
}

