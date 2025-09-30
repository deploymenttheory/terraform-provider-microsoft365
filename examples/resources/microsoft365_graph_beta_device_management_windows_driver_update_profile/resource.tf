resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "manual_example" {
  display_name       = "Windows Driver Updates - Production x"
  description        = "Driver update profile for production machines"
  approval_type      = "manual"
  role_scope_tag_ids = [8, 9]

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

resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "automatic_example" {
  display_name                = "Windows Driver Updates - Production y"
  description                 = "Driver update profile for production machines"
  approval_type               = "automatic"
  deployment_deferral_in_days = 14
  role_scope_tag_ids          = [8, 9]

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