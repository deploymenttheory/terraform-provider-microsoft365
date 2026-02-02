resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "all_duplicate" {
  display_name                                = "All Platforms Cleanup - Duplicate"
  device_cleanup_rule_platform_type           = "all"
  device_inactivity_before_retirement_in_days = 90
}
