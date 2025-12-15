data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_id" {
  filter_type  = "id"
  filter_value = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  timeouts = {
    read = "10s"
  }
}

