resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_014" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    win_get = {
      notifications = "showAll"

      install_time_settings = {
        use_local_time     = true
        deadline_date_time = "2027-01-01T12:00:00Z"
      }

      restart_settings = {
        grace_period_in_minutes                         = 90
        countdown_display_before_restart_in_minutes     = 20
        restart_notification_snooze_duration_in_minutes = 10
      }
    }
  }
}
