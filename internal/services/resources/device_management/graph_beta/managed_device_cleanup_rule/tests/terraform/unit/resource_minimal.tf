resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "minimal" {
  display_name                               = "Test Minimal Managed Device Cleanup Rule - Unique"
  device_cleanup_rule_platform_type          = "windows"
  device_inactivity_before_retirement_in_days = 30

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


