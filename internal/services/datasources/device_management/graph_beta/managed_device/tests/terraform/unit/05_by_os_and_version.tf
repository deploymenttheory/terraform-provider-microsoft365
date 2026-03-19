data "microsoft365_graph_beta_device_management_managed_device" "test" {
  operating_system = "Windows"
  os_version       = "10.0.19045"
}
