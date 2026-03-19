data "microsoft365_graph_beta_device_management_managed_device" "test" {
  odata_query = "operatingSystem eq 'Windows'"
}
