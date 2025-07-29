# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_aosp_device_owner_compliance_policy" "minimal" {
  display_name = "AOSP Minimal Compliance Policy"
  description  = "Minimal AOSP device owner compliance policy with basic security requirements"

  # Basic password requirements
  passcode_required       = true
  passcode_minimum_length = 6

  # Security settings
  security_block_jailbroken_devices = true
  storage_require_encryption        = true

  # Scheduled actions for rules (required)
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
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
resource "microsoft365_graph_beta_device_management_aosp_device_owner_compliance_policy" "advanced" {
  display_name = "AOSP Advanced Compliance Policy"
  description  = "Advanced AOSP device owner compliance policy with strict security requirements"

  # Password requirements
  passcode_required                          = true
  passcode_minimum_length                    = 8
  passcode_minutes_of_inactivity_before_lock = 5
  passcode_required_type                     = "alphanumeric"

  # OS version requirements
  os_minimum_version               = "12.0"
  os_maximum_version               = "13.0"
  min_android_security_patch_level = "2023-01-01"

  # Security settings
  security_block_jailbroken_devices = true
  storage_require_encryption        = true

  # Role scope tags
  role_scope_tag_ids = ["0", "1"]

  # Scheduled actions for rules
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 24
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000"
          ]
        },
        {
          action_type              = "block"
          grace_period_hours       = 72
          notification_template_id = ""
        }
      ]
    }
  ]

  # Assignments
  assignments = [
    # Assignment targeting all devices
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "none"
    },
    # Assignment targeting a specific group
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "none"
    },
    # Exclusion group assignment
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    }
  ]
} 