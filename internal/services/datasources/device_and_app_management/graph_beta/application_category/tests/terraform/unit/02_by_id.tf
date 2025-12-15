data "microsoft365_graph_beta_device_and_app_management_application_category" "by_id" {
  filter_type  = "id"
  filter_value = "5b0e1e8d-7a5c-4f3a-9c2d-1e4f5a6b7c8d"
  timeouts = {
    read = "10s"
  }
}

