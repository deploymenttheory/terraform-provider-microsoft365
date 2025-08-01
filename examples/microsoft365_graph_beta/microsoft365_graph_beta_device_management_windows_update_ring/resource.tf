resource "microsoft365_graph_beta_device_management_windows_update_ring" "minimal" {
  display_name                            = "Test Minimal Windows Update Ring - Unique"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 0
  feature_updates_deferral_period_in_days = 0
  allow_windows11_upgrade                 = true
  skip_checks_before_restart              = false
  automatic_update_mode                   = "userDefined"
  feature_updates_rollback_window_in_days = 10

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "maximal" {
  display_name                                 = "Test Maximal Windows Update Ring - Unique"
  description                                  = "Maximal Windows update ring for testing with all features"
  microsoft_update_service_allowed             = true
  drivers_excluded                             = false
  quality_updates_deferral_period_in_days      = 7
  feature_updates_deferral_period_in_days      = 14
  allow_windows11_upgrade                      = false
  skip_checks_before_restart                   = true
  automatic_update_mode                        = "autoInstallAndRebootAtScheduledTime"
  business_ready_updates_only                  = "businessReadyOnly"
  delivery_optimization_mode                   = "httpWithPeeringNat"
  prerelease_features                          = "settingsOnly"
  update_weeks                                 = "firstWeek"
  active_hours_start                           = "09:00:00"
  active_hours_end                             = "17:00:00"
  user_pause_access                            = "disabled"
  user_windows_update_scan_access              = "disabled"
  update_notification_level                    = "defaultNotifications"
  feature_updates_rollback_window_in_days      = 10
  engaged_restart_deadline_in_days             = 3
  engaged_restart_snooze_schedule_in_days      = 1
  engaged_restart_transition_schedule_in_days  = 2
  auto_restart_notification_dismissal          = "automatic"
  schedule_restart_warning_in_hours            = 4
  schedule_imminent_restart_warning_in_minutes = 15
  role_scope_tag_ids                           = ["0", "1"]

  uninstall = {
    feature_updates_will_be_rolled_back = true
    quality_updates_will_be_rolled_back = false
  }

  deadline_settings = {
    deadline_for_feature_updates_in_days = 7
    deadline_for_quality_updates_in_days = 2
    deadline_grace_period_in_days        = 1
    postpone_reboot_until_after_deadline = true
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}