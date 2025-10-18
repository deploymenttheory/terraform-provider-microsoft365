data "microsoft365_graph_beta_device_management_managed_device" "by_id" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000001"
}

