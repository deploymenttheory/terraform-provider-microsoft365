resource "microsoft365_graph_beta_device_management_windows_update_ring" "test" {
  display_name                                 = "Test Acceptance Windows Update Ring - Updated"
  description                                  = "Updated description for acceptance testing"
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
  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id, microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id]

  
  deadline_settings = {
    deadline_for_feature_updates_in_days = 7
    deadline_for_quality_updates_in_days = 2
    deadline_grace_period_in_days        = 1
    postpone_reboot_until_after_deadline = true
  }
}