# Scenario 7: Maximal Assignments
# This configuration demonstrates how to assign a Windows Update Ring to multiple groups
# and built-in targets, including group assignments, all licensed users, all devices,
# and exclusion groups.

# Example groups for assignment (you would use your actual group IDs)
resource "microsoft365_graph_beta_groups_group" "update_ring_group_1" {
  display_name     = "Windows Update Ring - Group 1"
  mail_nickname    = "windows-update-ring-group-1"
  mail_enabled     = false
  security_enabled = true
  description      = "First group for windows update ring assignments"
}

resource "microsoft365_graph_beta_groups_group" "update_ring_group_2" {
  display_name     = "Windows Update Ring - Group 2"
  mail_nickname    = "windows-update-ring-group-2"
  mail_enabled     = false
  security_enabled = true
  description      = "Second group for windows update ring assignments"
}

resource "microsoft365_graph_beta_groups_group" "update_ring_exclusion_group" {
  display_name     = "Windows Update Ring - Exclusion Group"
  mail_nickname    = "windows-update-ring-exclusion-group"
  mail_enabled     = false
  security_enabled = true
  description      = "Exclusion group for windows update ring assignments"
}

# Windows Update Ring with comprehensive assignments
resource "microsoft365_graph_beta_device_management_windows_update_ring" "maximal_assignments" {
  display_name                            = "Windows Update Ring - Maximal Assignments"
  description                             = "Demonstrates multiple assignment types"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "notifyDownload"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }

  # Multiple assignment types
  assignments = [
    # Assign to specific group 1
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.update_ring_group_1.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Assign to specific group 2
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.update_ring_group_2.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Assign to all licensed users
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Assign to all devices
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    # Exclude a specific group
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.update_ring_exclusion_group.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    }
  ]
}

