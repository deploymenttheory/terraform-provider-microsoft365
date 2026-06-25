resource "microsoft365_graph_beta_device_management_operation_approval_policy" "minimal" {
  display_name = "Test Minimal Operation Approval Policy"

  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }

  approver_group_ids = ["11111111-1111-1111-1111-111111111111"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
