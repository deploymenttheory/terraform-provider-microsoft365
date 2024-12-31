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

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "Adobe_Acrobat_Reader_DC" {
  package_identifier              = "xpdp273c0xhqh2" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = false
  display_name                    = "Adobe Acrobat Reader DC"
  description                     = "Adobe Acrobat Reader DC is the free, trusted standard for viewing, printing, signing, and annotating PDFs. It's the only PDF viewer that can open and interact with all types of PDF content – including forms and multimedia."
  publisher                       = "Adobe Inc."
  large_icon = {
    type  = "image/png"
    value = filebase64("${path.module}/Adobe_Reader_XI_icon.png")
  }
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
  package_identifier = "XP9KHM4BK9FZ7Q" # The unique identifier for the app obtained from msft app store

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

  # App assignments configuration
  assignments = {
    id            = "assignment_id" # Read-only, typically auto-generated
    mobile_app_id = "app_id_value"  # The ID of the app being assigned

    mobile_app_assignments = [
      # 2 Assignments for "available" intent
      {
        intent = "required" # Possible values: available, required, uninstall, availableWithoutEnrollment
        source = "direct"   # Possible values: direct, policySets

        target = {
          target_type                                      = "groupAssignmentTarget" # Possible values: groupAssignmentTarget, allLicensedUsersAssignmentTarget, etc.
          group_id                                         = "group_id_value"
          device_and_app_management_assignment_filter_id   = "filter_id_value"
          device_and_app_management_assignment_filter_type = "include" # Possible values: include, exclude, none
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "showAll" # Possible values: showAll, showReboot, hideAll
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 15
            countdown_display_before_restart_in_minutes     = 5
            restart_notification_snooze_duration_in_minutes = 10
          }
        }
      },
      {
        id        = "assignment_2"
        intent    = "available"
        source    = "policySets"
        source_id = "source_id_2"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_available_2"
          device_and_app_management_assignment_filter_id   = "filter_id_2"
          device_and_app_management_assignment_filter_type = "include"
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "showReboot"
        }
      },

      # 2 Assignments for "required" intent
      {
        id        = "assignment_3"
        intent    = "required"
        source    = "direct"
        source_id = "source_id_3"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_required_1"
          device_and_app_management_assignment_filter_id   = "filter_id_3"
          device_and_app_management_assignment_filter_type = "exclude"
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "hideAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
        }
      },
      {
        id        = "assignment_4"
        intent    = "required"
        source    = "policySets"
        source_id = "source_id_4"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_required_2"
          device_and_app_management_assignment_filter_id   = "filter_id_4"
          device_and_app_management_assignment_filter_type = "include"
          is_exclusion_group                               = true
        }

        settings = {
          notifications = "showAll"
        }
      },

      # 2 Assignments for "uninstall" intent
      {
        id        = "assignment_5"
        intent    = "uninstall"
        source    = "direct"
        source_id = "source_id_5"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_uninstall_1"
          device_and_app_management_assignment_filter_id   = "filter_id_5"
          device_and_app_management_assignment_filter_type = "none"
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "showReboot"
          restart_settings = {
            grace_period_in_minutes                         = 15
            countdown_display_before_restart_in_minutes     = 5
            restart_notification_snooze_duration_in_minutes = 10
          }
        }
      },
      {
        id        = "assignment_6"
        intent    = "uninstall"
        source    = "policySets"
        source_id = "source_id_6"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_uninstall_2"
          device_and_app_management_assignment_filter_id   = "filter_id_6"
          device_and_app_management_assignment_filter_type = "exclude"
          is_exclusion_group                               = true
        }

        settings = {
          notifications = "hideAll"
        }
      }
    ]
  }
}
