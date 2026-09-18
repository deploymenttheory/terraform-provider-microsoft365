resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_004" {
  name                       = "acc-test-linux-compliance-004-${random_string.test_suffix.result}"
  device_encryption_required = false
  description                = ""
  role_scope_tag_ids         = ["0"]
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "5m"
  }
}
