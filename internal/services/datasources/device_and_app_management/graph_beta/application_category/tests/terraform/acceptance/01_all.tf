data "microsoft365_graph_beta_device_and_app_management_application_category" "all" {
  filter_type = "all"
  timeouts = {
    read = "30s"
  }
}

