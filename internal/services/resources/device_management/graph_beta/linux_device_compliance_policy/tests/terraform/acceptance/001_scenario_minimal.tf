resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_001" {
  name                       = "acc-test-linux-compliance-001-${random_string.test_suffix.result}"
  device_encryption_required = false
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "5m"
  }
}
