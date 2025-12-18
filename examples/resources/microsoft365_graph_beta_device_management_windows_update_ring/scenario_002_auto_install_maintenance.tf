# Scenario 2: Auto Install at Maintenance Time
# This configuration automatically installs updates outside of active hours but requires
# user interaction to restart. Updates install during maintenance windows (outside active hours).

resource "microsoft365_graph_beta_device_management_windows_update_ring" "auto_install_maintenance" {
  display_name                            = "Windows Update Ring - Auto Install at Maintenance Time"
  description                             = "Automatically install updates at maintenance time"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 30
  feature_updates_deferral_period_in_days = 30
  allow_windows11_upgrade                 = true
  quality_updates_paused                  = false
  feature_updates_paused                  = false
  business_ready_updates_only             = "windowsInsiderBuildRelease"
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAtMaintenanceTime"
  active_hours_start                      = "08:00:00"
  active_hours_end                        = "17:00:00"
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

