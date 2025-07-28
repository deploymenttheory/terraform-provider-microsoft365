# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_ios_device_compliance_policy" "minimal" {
  display_name = "iOS Minimal Compliance Policy"
  description  = "Minimal iOS device compliance policy with basic security requirements"

  # Basic security requirements
  passcode_required                 = true
  security_block_jailbroken_devices = true

  # Scheduled actions for rules (required)
  scheduled_actions_for_rule = [
    {
      rule_name = "PasscodeRequired"
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 0
        }
      ]
    }
  ]
}

# Example with advanced security settings
resource "microsoft365_graph_beta_device_management_ios_device_compliance_policy" "advanced" {
  display_name = "iOS Advanced Compliance Policy"
  description  = "Advanced iOS device compliance policy with strict security requirements"

  # Strict passcode requirements
  passcode_required                                    = true
  passcode_block_simple                                = true
  passcode_minimum_length                              = 8
  passcode_minimum_character_set_count                 = 3
  passcode_required_type                               = "alphanumeric"
  passcode_expiration_days                             = 30
  passcode_previous_passcode_block_count               = 5
  passcode_minutes_of_inactivity_before_lock           = 2
  passcode_minutes_of_inactivity_before_screen_timeout = 1

  # Strict OS version requirements
  os_minimum_version       = "16.0"
  os_minimum_build_version = "20A362"

  # Security settings
  security_block_jailbroken_devices                  = true
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "high"
  advanced_threat_protection_required_security_level = "secured"
  managed_email_profile_required                     = true

  # Restricted apps
  restricted_apps = [
    {
      name          = "Prohibited App"
      publisher     = "Prohibited Publisher"
      app_id        = "com.prohibited.app"
      app_store_url = "https://apps.apple.com/app/prohibited-app/id123456789"
    }
  ]

  # Scheduled actions for rules
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000",
            "00000000-0000-0000-0000-000000000000"
          ]
        },
        {
          action_type              = "notification"
          grace_period_hours       = 24
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000",
            "00000000-0000-0000-0000-000000000000"
          ]
        },
        {
          action_type              = "remoteLock"
          grace_period_hours       = 72
          notification_template_id = ""
        },
        {
          action_type              = "retire"
          grace_period_hours       = 120
          notification_template_id = ""
        },
        {
          action_type              = "block"
          grace_period_hours       = 0
          notification_template_id = ""
        }
      ]
    }
  ]

  # Assignments
  assignments = [
    # Assignment targeting all devices with an include filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Assignment targeting all licensed users with an exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    }
  ]
}