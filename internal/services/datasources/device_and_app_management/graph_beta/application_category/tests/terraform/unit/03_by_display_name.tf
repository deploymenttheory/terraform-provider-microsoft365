data "microsoft365_graph_beta_device_and_app_management_application_category" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Business"
  timeouts = {
    read = "10s"
  }
}

