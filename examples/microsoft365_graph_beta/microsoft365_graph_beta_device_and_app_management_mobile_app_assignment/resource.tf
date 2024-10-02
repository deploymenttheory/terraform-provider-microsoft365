resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "example" {
  source_id = "app-123456"

  target = {
    type = "group"
    device_and_app_management_assignment_filter_id = "filter-789012"
    device_and_app_management_assignment_filter_type = "include"
    group_id = "group-345678"
  }

  intent = "available"

  settings = {
    notifications = "showAll"

    restart_settings ={
      grace_period_in_minutes = 60
      countdown_display_before_restart_in_minutes = 15
      restart_notification_snooze_duration_in_minutes = 5
    }

    install_time_settings ={
      use_local_time = true
      deadline_date_time = "2023-12-31T23:59:59Z"
    }
  }

  source = "direct"

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}