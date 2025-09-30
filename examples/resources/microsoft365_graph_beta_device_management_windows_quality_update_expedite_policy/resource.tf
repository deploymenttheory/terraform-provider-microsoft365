resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "example" {
  display_name       = "Windows Quality Update expedite policy"
  description        = "Emergency fixes"
  role_scope_tag_ids = ["9", "8"]

  expedited_update_settings = {
    quality_update_release   = "2025-04-22T00:00:00Z"
    days_until_forced_reboot = 1
  }

  // Optional assignment blocks
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

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}