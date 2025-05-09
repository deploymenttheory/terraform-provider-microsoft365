resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "example" {
  package_identifier              = "xpfftq037jwmhs"
  automatically_generate_metadata = true

  # Optional metadata fields (will be auto-populated if automatically_generate_metadata = true)
  # display_name                  = "Microsoft Edge"
  # description                   = "Microsoft Edge browser"
  # publisher                     = "Microsoft Corporation"

  # Optional app information
  is_featured             = true
  privacy_information_url = "https://privacy.microsoft.com/en-us/privacystatement"
  information_url         = "https://www.microsoft.com/en-us/edge"
  owner                   = "IT Department"
  developer               = "Microsoft"
  notes                   = "Default browser for all corporate devices"

  # Required install experience settings
  install_experience = {
    run_as_account = "system" # Allowed values: "system" or "user"
  }

  # Optional role scope tag IDs
  role_scope_tag_ids = ["0a129961-8d6a-4496-8add-068fe16b13aa"]

  # App assignments
  assignments = {
    # Required fields
    intent = "required" # Allowed values: "available", "required", "uninstall", "availableWithoutEnrollment"
    source = "direct"   # Possible values: "direct", "policySets"

    # Target configuration (required)
    target = {
      target_type = "allDevices" # Target all devices in the tenant

      # Alternative target types (uncomment only one):
      # target_type = "groupAssignment"
      # group_id = "5df60fc9-54c9-4245-8c0f-a4082f2249c5" # Entra ID group ID

      # Optional assignment filter (if used, uncomment both lines)
      # device_and_app_management_assignment_filter_id = "21e35af7-5c85-4305-9c62-f96e5cf1f2b5"
      # device_and_app_management_assignment_filter_type = "include" # "include", "exclude", or "none" (default)
    }

    # WinGet specific settings (optional)
    settings = {
      win_get = {
        # Installation timing
        install_time_settings = {
          use_local_time     = true
          deadline_date_time = "2025-06-01T18:00:00Z"
        }

        # Notification options
        notifications = "showAll" # Allowed values: "showAll", "showReboot", "hideAll"

        # Restart settings
        restart_settings = {
          grace_period_in_minutes                         = 240 # 4 hours
          countdown_display_before_restart_in_minutes     = 30
          restart_notification_snooze_duration_in_minutes = 60
        }
      }
    }
  }

  # Optional timeouts
  timeouts = {
    create = "3m"
    update = "3m"
    read   = "3m"
    delete = "3m"
  }
}