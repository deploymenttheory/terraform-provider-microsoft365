resource "microsoft365_graph_beta_device_management_managed_device_cleanup_rule" "maximal" {
  display_name                                = "Test Maximal Managed Device Cleanup Rule - Unique"
  description                                 = "Maximal managed device cleanup rule for testing"
  device_cleanup_rule_platform_type           = "windows"
  device_inactivity_before_retirement_in_days = 180

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


