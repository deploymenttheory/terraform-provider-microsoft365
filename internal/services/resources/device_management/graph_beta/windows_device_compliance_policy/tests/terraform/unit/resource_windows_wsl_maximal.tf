
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "wsl" {
  display_name = "Windows 10/11 - WSL Compliance Policy"
  description  = "Windows device compliance policy with WSL distribution requirements"

  # Password requirements
  password_required                     = true
  password_block_simple                 = true
  password_required_to_unlock_from_idle = true
  password_minimum_length               = 8
  password_minimum_character_set_count  = 3
  password_required_type                = "alphanumeric"

  # Security requirements
  bit_locker_enabled  = true
  secure_boot_enabled = true
  tpm_required        = true

  # WSL distributions
  wsl_distributions = [
    {
      distribution       = "Ubuntu"
      minimum_os_version = "20.04"
      maximum_os_version = "22.04"
    },
    {
      distribution       = "Debian"
      minimum_os_version = "11.0"
      maximum_os_version = "12.0"
    }
  ]

  # Role scope tags
  role_scope_tag_ids = ["0"]

  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "retire"
          grace_period_hours       = 1440
          notification_template_id = ""
        },
        {
          action_type              = "notification"
          grace_period_hours       = 120
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = ["00000000-0000-0000-0000-000000000000",
          "00000000-0000-0000-0000-000000000000"]
        },
        {
          action_type              = "block"
          grace_period_hours       = 1152
          notification_template_id = "00000000-0000-0000-0000-000000000000"
        },
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000",
          "00000000-0000-0000-0000-000000000000"]
        }
      ]
    }
  ]

  # Assignments
  assignments = [
    # Optional: Assignment targeting all devices with a daily schedule
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with an hourly schedule
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]
} 