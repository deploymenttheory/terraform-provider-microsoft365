data "microsoft365_graph_beta_device_and_app_management_application_category" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "startswith(displayName, 'Business')"
  timeouts = {
    read = "10s"
  }
}

