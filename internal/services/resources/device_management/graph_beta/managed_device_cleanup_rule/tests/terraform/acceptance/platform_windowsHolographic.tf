resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "windowsHolographic" {
  display_name                                = "Acc - Windows Holographic Cleanup"
  device_cleanup_rule_platform_type           = "windowsHolographic"
  device_inactivity_before_retirement_in_days = 60
}


