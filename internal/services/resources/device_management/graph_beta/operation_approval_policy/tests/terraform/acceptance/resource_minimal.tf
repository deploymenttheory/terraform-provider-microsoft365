resource "random_id" "test" {
  byte_length = 4
}

resource "microsoft365_graph_beta_groups_group" "approvers" {
  display_name     = "acc-test-oap-approvers-${random_id.test.hex}"
  mail_nickname    = "acctestoapapprovers${random_id.test.hex}"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name = "Test Acceptance Operation Approval Policy"

  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }

  approver_group_ids = [microsoft365_graph_beta_groups_group.approvers.id]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
