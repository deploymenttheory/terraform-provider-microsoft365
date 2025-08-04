resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Acceptance Role Scope Tag - Updated"
  description  = "Updated description for acceptance testing"

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

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}