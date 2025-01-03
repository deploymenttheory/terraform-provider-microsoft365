resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "whatsapp" {
  package_identifier              = "9NKSQGP7F2NH" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = true

  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "visual_studio_code" {
  package_identifier              = "XP9KHM4BK9FZ7Q" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = true
  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10s"
    update = "10s"
    delete = "10s"
  }

  # App assignments configuration
  assignments = [
    {
      intent = "available"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "612233b1-55ca-4815-a6b9-5c4aa5a4ac87"
        device_and_app_management_assignment_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "allDevices"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "allLicensedUsers"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "exclude"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "required"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
        device_and_app_management_assignment_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
        device_and_app_management_assignment_filter_type = "include"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 100
            countdown_display_before_restart_in_minutes     = 15
            restart_notification_snooze_duration_in_minutes = 42
          }
        }
      }
    },
    {
      intent = "uninstall"
      source = "direct"
      target = {
        target_type                                      = "groupAssignment"
        group_id                                         = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
        device_and_app_management_assignment_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        device_and_app_management_assignment_filter_type = "exclude"
      }
      settings = {
        win_get = {
          notifications = "showAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 200
            countdown_display_before_restart_in_minutes     = 50
            restart_notification_snooze_duration_in_minutes = 1
          }
        }
      }
    },
  ]
}
