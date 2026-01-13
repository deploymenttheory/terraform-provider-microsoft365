resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "transition" {
  display_name       = "acc-test-005-lifecycle-maximal-${random_string.test_suffix.result}"
  description        = "acc-test-005-lifecycle-maximal"
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]
}
