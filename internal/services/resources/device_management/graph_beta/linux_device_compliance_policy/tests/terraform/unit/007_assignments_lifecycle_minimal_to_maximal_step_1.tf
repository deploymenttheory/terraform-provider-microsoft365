resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_007" {
  name                       = "unit-test-linux-compliance-007"
  device_encryption_required = false
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
  ]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
