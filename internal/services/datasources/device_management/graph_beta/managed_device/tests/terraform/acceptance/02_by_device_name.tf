data "microsoft365_graph_beta_device_management_managed_device" "by_device_name" {
  filter_type  = "device_name"
  filter_value = "DESKTOP"
}

