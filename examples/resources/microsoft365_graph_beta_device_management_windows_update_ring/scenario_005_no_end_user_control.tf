# Scenario 5: Auto Install and Reboot Without End User Control
# This configuration provides the most aggressive update policy, automatically installing and
# restarting devices without user interaction or the ability to postpone updates.
# Use with caution as it provides no end-user control.

resource "microsoft365_graph_beta_device_management_windows_update_ring" "no_end_user_control" {
  display_name                            = "Windows Update Ring - No End User Control"
  description                             = "Automatically install and reboot without end user control"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAndRebootWithoutEndUserControl"
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
}

