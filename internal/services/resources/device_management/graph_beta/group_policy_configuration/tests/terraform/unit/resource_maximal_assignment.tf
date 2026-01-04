resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal_assignment" {
  display_name       = "Maximal Assignment Group Policy Configuration"
  description        = "Configuration with comprehensive assignments"
  role_scope_tag_ids = ["0", "1"]

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
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "exclude"
    }
  ]
}

