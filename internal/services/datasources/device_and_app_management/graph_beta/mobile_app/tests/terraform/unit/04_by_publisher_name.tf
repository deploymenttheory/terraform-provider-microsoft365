data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_publisher" {
  filter_type  = "publisher_name"
  filter_value = "Microsoft"
}

