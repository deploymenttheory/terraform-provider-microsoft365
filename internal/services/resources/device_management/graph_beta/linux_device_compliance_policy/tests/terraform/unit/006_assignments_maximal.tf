resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_006" {
  name                       = "unit-test-linux-compliance-006"
  device_encryption_required = false
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "22222222-2222-2222-2222-222222222222"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
  ]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
