resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_004" {
  name                       = "unit-test-linux-compliance-004"
  device_encryption_required = false
  description                = ""
  role_scope_tag_ids         = ["0"]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
