# An uninstall intent accepts no Apple app settings at all: the service rejects
# uninstallOnDeviceRemoval and preventManagedAppBackup for it, as well as isRemovable. An
# empty settings block must therefore produce an empty payload rather than a defaulted one.
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "uninstall_ios_store" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "uninstall"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    ios_store = {}
  }
}
