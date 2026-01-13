resource "microsoft365_graph_beta_device_management_group_policy_configuration" "transition" {
  display_name       = "unit-test-005-lifecycle-maximal"
  description        = "unit-test-005-lifecycle-maximal"
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
    }
  ]
}
