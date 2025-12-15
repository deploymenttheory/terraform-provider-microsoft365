data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Microsoft"
}

