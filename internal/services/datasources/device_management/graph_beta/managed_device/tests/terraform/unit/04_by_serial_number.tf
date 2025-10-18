data "microsoft365_graph_beta_device_management_managed_device" "by_serial_number" {
  filter_type  = "serial_number"
  filter_value = "SN-WIN-001"
}

