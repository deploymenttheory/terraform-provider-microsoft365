resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "windows" {
  display_name                                = "Windows Cleanup"
  device_cleanup_rule_platform_type           = "windows"
  device_inactivity_before_retirement_in_days = 60
}


