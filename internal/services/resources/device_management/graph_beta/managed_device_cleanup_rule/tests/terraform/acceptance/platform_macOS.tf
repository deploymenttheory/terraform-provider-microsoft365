resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "macOS" {
  display_name                                = "Acc - macOS Cleanup"
  device_cleanup_rule_platform_type           = "macOS"
  device_inactivity_before_retirement_in_days = 60
}


