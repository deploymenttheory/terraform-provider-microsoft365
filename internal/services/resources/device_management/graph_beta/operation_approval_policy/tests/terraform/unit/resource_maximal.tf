resource "microsoft365_graph_beta_device_management_operation_approval_policy" "maximal" {
  display_name    = "Test Maximal Operation Approval Policy"
  description     = "Maximal operation approval policy for testing with all features"
  policy_type     = "script"
  policy_platform = "windows10AndLater"

  policy_set = {
    policy_type     = "script"
    policy_platform = "windows10AndLater"
  }

  approver_group_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
