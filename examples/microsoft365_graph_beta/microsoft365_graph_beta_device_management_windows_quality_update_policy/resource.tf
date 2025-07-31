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
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    # Assignment targeting a specific group
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"

    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
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