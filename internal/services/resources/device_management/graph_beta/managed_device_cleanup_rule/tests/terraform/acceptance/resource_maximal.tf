resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "test" {
  display_name                                = "Test Acceptance Managed Device Cleanup Rule - Updated"
  description                                 = "Updated description for acceptance testing"
  device_cleanup_rule_platform_type           = "windows"
  device_inactivity_before_retirement_in_days = 90
}


