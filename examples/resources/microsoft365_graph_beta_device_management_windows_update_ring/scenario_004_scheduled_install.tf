# Scenario 4: Auto Install and Restart at Scheduled Time
# This configuration automatically installs updates and restarts devices at a specific scheduled
# time and day. Use this for predictable maintenance windows.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "scheduled_install" {
  display_name                            = "Windows Update Ring - Scheduled Install and Restart"
  description                             = "Automatically install and restart at scheduled time"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAndRebootAtScheduledTime"
  scheduled_install_day                   = "everyday"
  scheduled_install_time                  = "03:00:00"
  user_pause_access                       = "enabled"
  user_windows_update_scan_access         = "enabled"
  update_notification_level               = "restartWarningsOnly"
  update_weeks                            = "everyWeek"
  feature_updates_rollback_window_in_days = 10

  deadline_settings = {
    deadline_for_feature_updates_in_days = 5
    deadline_for_quality_updates_in_days = 7
    deadline_grace_period_in_days        = 7
    postpone_reboot_until_after_deadline = false
  }
}

