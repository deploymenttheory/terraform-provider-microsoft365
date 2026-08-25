resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_012" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "available"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    win32_catalog = {
      notifications                  = "showAll"
      delivery_optimization_priority = "foreground"

      auto_update_settings = {
        auto_update_superseded_apps_state = "enabled"
      }

      install_time_settings = {
        use_local_time     = true
        deadline_date_time = "2027-01-01T12:00:00Z"
        start_date_time    = "2026-12-01T12:00:00Z"
      }

      restart_settings = {
        grace_period_in_minutes                         = 60
        countdown_display_before_restart_in_minutes     = 15
        restart_notification_snooze_duration_in_minutes = 10
      }
    }
  }
}
