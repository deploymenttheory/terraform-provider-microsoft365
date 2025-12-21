resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "quality_update_policy_example" {
  display_name       = "Windows Quality Update Policy"
  description        = "Monthly quality updates for Windows devices"
  hotpatch_enabled   = true
  role_scope_tag_ids = ["9", "8"]

  # Assignments
  assignments = [
    # Assignment targeting a specific group
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    # Additional assignment targeting a specific group
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]

  // Optional timeout block
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}