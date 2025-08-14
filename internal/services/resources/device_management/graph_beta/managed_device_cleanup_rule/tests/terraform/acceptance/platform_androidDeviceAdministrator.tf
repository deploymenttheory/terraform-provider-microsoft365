resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "androidDeviceAdministrator" {
  display_name                                = "Acc - Android Device Admin Cleanup"
  device_cleanup_rule_platform_type           = "androidDeviceAdministrator"
  device_inactivity_before_retirement_in_days = 60
}


