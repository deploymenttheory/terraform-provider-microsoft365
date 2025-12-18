
resource "random_string" "test_011" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_011_1" {
  display_name     = "acc-test-group-011-1-${random_string.test_011.result}"
  mail_nickname    = "acc-test-group-011-1-${random_string.test_011.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows update ring assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_011_2" {
  display_name     = "acc-test-group-011-2-${random_string.test_011.result}"
  mail_nickname    = "acc-test-group-011-2-${random_string.test_011.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows update ring assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_011_3" {
  display_name     = "acc-test-group-011-3-${random_string.test_011.result}"
  mail_nickname    = "acc-test-group-011-3-${random_string.test_011.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows update ring exclusion assignments"
  hard_delete      = true
}

# ==============================================================================
# Windows Update Ring Resource with Maximal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_update_ring" "test_011" {
  display_name                            = "acc-test-windows-update-ring-011-assignments-downgrade-${random_string.test_011.result}"
  description                             = "Assignments Downgrade Step 1: Maximal"
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

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_011_1.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_011_2.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_011_3.id
      filter_type = "none"
      filter_id   = "00000000-0000-0000-0000-000000000000"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

