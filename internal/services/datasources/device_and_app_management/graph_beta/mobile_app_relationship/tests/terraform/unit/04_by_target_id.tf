data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_id" {
  filter_type  = "target_id"
  filter_value = "app-target-001"
  timeouts = {
    read = "10s"
  }
}

