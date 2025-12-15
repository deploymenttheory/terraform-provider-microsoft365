data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "all" {
  filter_type = "all"
  timeouts = {
    read = "10s"
  }
}

