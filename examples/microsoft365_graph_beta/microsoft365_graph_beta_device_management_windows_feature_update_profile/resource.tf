resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example" {
  display_name                                            = "Windows 11 22H2 Deployment x"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // rollout_settings = Make update available gradually
  rollout_settings = {
    offer_start_date_time_in_utc = "2025-05-01T00:00:00Z"
    offer_end_date_time_in_utc   = "2025-06-30T23:59:59Z"
    offer_interval_in_days       = 7
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

  # Optional timeout block
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }

}

resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example_2" {
  display_name                                            = "Windows 11 22H2 Deployment y"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // rollout_settings = Make update available on a specific date
  rollout_settings = {
    offer_start_date_time_in_utc = "2025-08-01T00:00:00Z"
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

  # Optional timeout block
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example_3" {
  display_name                                            = "Windows 11 22H2 Deployment z"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = false
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // include no rollout_settings block to make Make update available as soon as possible

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

  # Optional timeout block
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

