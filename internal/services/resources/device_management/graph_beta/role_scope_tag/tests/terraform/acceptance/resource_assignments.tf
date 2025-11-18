resource "microsoft365_graph_beta_device_management_role_scope_tag" "assignments" {
  display_name = "acc-test-role-scope-tag-assignments-${random_string.suffix.result}"
  description  = "acc-test-role-scope-tag-assignments-${random_string.suffix.result}"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    }
  ]

  # Only depend on groups that are actually assigned
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}