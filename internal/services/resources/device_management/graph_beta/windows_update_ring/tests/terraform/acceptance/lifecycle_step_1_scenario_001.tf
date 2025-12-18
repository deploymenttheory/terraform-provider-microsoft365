
resource "random_string" "test_007" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "test_007" {
  display_name                            = "acc-test-windows-update-ring-007-lifecycle-${random_string.test_007.result}"
  description                             = "Lifecycle Step 1: Notify Download"
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

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

