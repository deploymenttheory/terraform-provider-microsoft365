resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "androidDedicatedAndFullyManagedCorporateOwnedWorkProfile" {
  display_name                                = "Acc - Android COPE/COBO Cleanup"
  device_cleanup_rule_platform_type           = "androidDedicatedAndFullyManagedCorporateOwnedWorkProfile"
  device_inactivity_before_retirement_in_days = 60
}


