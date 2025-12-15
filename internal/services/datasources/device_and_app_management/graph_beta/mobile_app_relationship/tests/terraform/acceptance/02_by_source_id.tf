data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_source_id" {
  filter_type  = "source_id"
  filter_value = "app-source-test-001"
  timeouts = {
    read = "30s"
  }
}

