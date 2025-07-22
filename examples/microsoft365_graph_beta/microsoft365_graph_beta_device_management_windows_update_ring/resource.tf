# Example 1: Basic Windows Update Ring Configuration
resource "microsoft365_graph_beta_device_management_windows_update_ring" "basic_update_ring" {
  display_name       = "Standard Windows Update Ring"
  description        = "Default update ring for standard workstations"
  role_scope_tag_ids = ["0"]

  microsoft_update_service_allowed             = true
  drivers_excluded                             = false
  quality_updates_deferral_period_in_days      = 7
  feature_updates_deferral_period_in_days      = 14
  allow_windows11_upgrade                      = false
  quality_updates_paused                       = false
  feature_updates_paused                       = false
  business_ready_updates_only                  = "businessReadyOnly"
  automatic_update_mode                        = "autoInstallAtMaintenanceTime"
  active_hours_start                           = "08:00:00"
  active_hours_end                             = "17:00:00"
  user_pause_access                            = "enabled"
  user_windows_update_scan_access              = "enabled"
  update_notification_level                    = "defaultNotifications"
  deadline_for_feature_updates_in_days         = 7
  deadline_for_quality_updates_in_days         = 3
  deadline_grace_period_in_days                = 2
  skip_checks_before_restart                   = false
  postpone_reboot_until_after_deadline         = true
  engaged_restart_deadline_in_days             = 7
  engaged_restart_snooze_schedule_in_days      = 2
  engaged_restart_transition_schedule_in_days  = 7
  auto_restart_notification_dismissal          = "notConfigured"
  schedule_restart_warning_in_hours            = 4
  schedule_imminent_restart_warning_in_minutes = 30
  delivery_optimization_mode                   = "httpWithPeeringNat"
  prerelease_features                          = "notAllowed"
  update_weeks                                 = "everyWeek"
  feature_updates_rollback_window_in_days      = 10

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}

# Example 2: Minimal Windows Update Ring Configuration
resource "microsoft365_graph_beta_device_management_windows_update_ring" "minimal_update_ring" {
  display_name = "Minimal Windows Update Ring"
  description  = "Basic update ring with minimal configuration"

  # Only required fields
  automatic_update_mode = "autoInstallAndRebootAtScheduledTime"
}

# Example 3: Advanced Windows Update Ring Configuration
resource "microsoft365_graph_beta_device_management_windows_update_ring" "advanced_update_ring" {
  display_name = "Advanced Windows Update Ring"
  description  = "Advanced update ring for specialized workstations"

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Update service configuration
  microsoft_update_service_allowed = true
  drivers_excluded                 = true

  # Update deferral configuration
  quality_updates_deferral_period_in_days = 14
  feature_updates_deferral_period_in_days = 30

  # Windows 11 upgrade settings
  allow_windows11_upgrade = true

  # Update pause configuration
  quality_updates_paused = false
  feature_updates_paused = false

  # Update branch configuration
  business_ready_updates_only = "businessReadyOnly"

  # Automatic update mode
  automatic_update_mode = "autoInstallAndRebootAtScheduledTime"

  # Active hours configuration
  active_hours_start = "07:00:00"
  active_hours_end   = "19:00:00"

  # User control settings
  user_pause_access               = "disabled"
  user_windows_update_scan_access = "enabled"
  update_notification_level       = "restartWarningsOnly"

  # Deadline configuration
  deadline_for_feature_updates_in_days = 14
  deadline_for_quality_updates_in_days = 7
  deadline_grace_period_in_days        = 3

  # Restart behavior
  skip_checks_before_restart           = false
  postpone_reboot_until_after_deadline = true

  # Engaged restart settings
  engaged_restart_deadline_in_days            = 14
  engaged_restart_snooze_schedule_in_days     = 3
  engaged_restart_transition_schedule_in_days = 14

  # Restart notifications
  auto_restart_notification_dismissal          = "automatic"
  schedule_restart_warning_in_hours            = 8
  schedule_imminent_restart_warning_in_minutes = 60

  # Delivery optimization
  delivery_optimization_mode = "httpWithInternetPeering"

  # Feature management
  prerelease_features = "notAllowed"
  update_weeks        = "firstWeek"

  # Rollback settings
  feature_updates_rollback_window_in_days = 20
}
