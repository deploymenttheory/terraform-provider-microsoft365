resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal_assignment" {
  display_name       = "unit-test-004-maximal-assignment"
  description        = "unit-test-004-maximal-assignment"
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]
}
