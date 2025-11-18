resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "acc-test-role-scope-tag-maximal-${random_string.suffix.result}"
  description  = "acc-test-role-scope-tag-maximal-${random_string.suffix.result}"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    }
  ]

  # Only depend on groups that are actually assigned
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_3,
    microsoft365_graph_beta_groups_group.acc_test_group_4
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}