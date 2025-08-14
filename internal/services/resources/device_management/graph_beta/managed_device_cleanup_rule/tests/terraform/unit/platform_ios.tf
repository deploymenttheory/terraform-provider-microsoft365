resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "ios" {
  display_name                                = "iOS Cleanup"
  device_cleanup_rule_platform_type           = "ios"
  device_inactivity_before_retirement_in_days = 60
}


