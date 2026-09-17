resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_008" {
  name                       = "unit-test-linux-compliance-008"
  device_encryption_required = false
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
