resource "microsoft365_graph_beta_device_management_app_control_for_business_managed_installer" "enabled" {
  intune_management_extension_as_managed_installer = "Enabled"

  timeouts = {
    create = "300s"
    read   = "300s"
    update = "300s"
    delete = "300s"
  }
}