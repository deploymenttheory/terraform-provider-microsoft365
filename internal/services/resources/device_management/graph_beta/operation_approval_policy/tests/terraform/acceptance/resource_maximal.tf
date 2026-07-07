resource "random_id" "test" {
  byte_length = 4
}

resource "microsoft365_graph_beta_groups_group" "approvers" {
  display_name     = "acc-test-oap-approvers-${random_id.test.hex}"
  mail_nickname    = "acctestoapapprovers${random_id.test.hex}"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_groups_group" "approvers_secondary" {
  display_name     = "acc-test-oap-approvers2-${random_id.test.hex}"
  mail_nickname    = "acctestoapapprovers2${random_id.test.hex}"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name    = "Test Acceptance Operation Approval Policy - Updated"
  description     = "Updated description for acceptance testing"
  policy_type     = "script"
  policy_platform = "windows10AndLater"

  policy_set = {
    policy_type     = "script"
    policy_platform = "windows10AndLater"
  }

  approver_group_ids = [
    microsoft365_graph_beta_groups_group.approvers.id,
    microsoft365_graph_beta_groups_group.approvers_secondary.id,
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
