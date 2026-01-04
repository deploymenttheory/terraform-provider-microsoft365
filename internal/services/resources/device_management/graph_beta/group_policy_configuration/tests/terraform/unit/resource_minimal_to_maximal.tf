resource "microsoft365_graph_beta_device_management_group_policy_configuration" "transition" {
  display_name       = "Transition Group Policy Configuration"
  description        = "Configuration that transitions from minimal to maximal"
  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_id   = "22222222-2222-2222-2222-222222222222"
      filter_type = "include"
    }
  ]
}

