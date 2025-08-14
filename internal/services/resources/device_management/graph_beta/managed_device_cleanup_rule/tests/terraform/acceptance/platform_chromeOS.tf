resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "chromeOS" {
  display_name                                = "Acc - ChromeOS Cleanup"
  device_cleanup_rule_platform_type           = "chromeOS"
  device_inactivity_before_retirement_in_days = 60
}


