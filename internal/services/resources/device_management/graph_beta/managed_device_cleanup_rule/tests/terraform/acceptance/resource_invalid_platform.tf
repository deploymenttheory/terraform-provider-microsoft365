resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "test" {
  display_name                                = "Invalid Platform"
  device_cleanup_rule_platform_type           = "invalid_platform"
  device_inactivity_before_retirement_in_days = 60
}


