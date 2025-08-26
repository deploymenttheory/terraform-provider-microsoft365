resource "microsoft365_graph_beta_device_management_app_control_for_business_managed_installer" "maximal" {
  intune_management_extension_as_managed_installer = "Enabled"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}