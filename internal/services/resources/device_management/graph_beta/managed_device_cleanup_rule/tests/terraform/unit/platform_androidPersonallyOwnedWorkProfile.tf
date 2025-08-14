resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "androidPersonallyOwnedWorkProfile" {
  display_name                                = "Android Work Profile Cleanup"
  device_cleanup_rule_platform_type           = "androidPersonallyOwnedWorkProfile"
  device_inactivity_before_retirement_in_days = 60
}


