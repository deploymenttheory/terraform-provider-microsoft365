# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy" "minimal" {
  display_name = "Android Device Owner Minimal Compliance Policy"
  description  = "Minimal Android device owner compliance policy with basic security requirements"

  # Basic password requirements
  password_required       = true
  password_minimum_length = 6

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

# Example with comprehensive security settings
resource "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy" "comprehensive" {
  display_name = "Android Device Owner Comprehensive Compliance Policy"
  description  = "Comprehensive Android device owner compliance policy with advanced security requirements"

  # Threat protection settings
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "medium"
  advanced_threat_protection_required_security_level = "high"

  # Security settings
  security_block_jailbroken_devices                        = true
  security_require_safety_net_attestation_basic_integrity  = true
  security_require_safety_net_attestation_certified_device = true
  security_require_intune_app_integrity                    = true
  require_no_pending_system_updates                        = true
  security_required_android_safety_net_evaluation_type     = "hardwareBacked"

  # OS version requirements
  os_minimum_version               = "14"
  os_maximum_version               = "15"
  min_android_security_patch_level = "February 1, 2025"

  # Comprehensive password requirements
  password_required                          = true
  password_minimum_length                    = 12
  password_minimum_letter_characters         = 2
  password_minimum_lower_case_characters     = 1
  password_minimum_upper_case_characters     = 1
  password_minimum_numeric_characters        = 2
  password_minimum_symbol_characters         = 1
  password_minimum_non_letter_characters     = 3
  password_required_type                     = "alphanumericWithSymbols"
  password_minutes_of_inactivity_before_lock = 5
  password_expiration_days                   = 90
  password_previous_password_count_to_block  = 5

  # Storage settings
  storage_require_encryption = true

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
          action_type        = "block"
          grace_period_hours = 72
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

# Example with moderate security settings for enterprise use
resource "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy" "enterprise" {
  display_name = "Android Device Owner Enterprise Compliance Policy"
  description  = "Enterprise Android device owner compliance policy balancing security and usability"

  # Threat protection settings
  device_threat_protection_enabled                 = true
  device_threat_protection_required_security_level = "low"

  # Security settings
  security_block_jailbroken_devices                       = true
  security_require_safety_net_attestation_basic_integrity = true
  security_require_intune_app_integrity                   = true
  security_required_android_safety_net_evaluation_type    = "basic"

  # OS version requirements - allowing broader range for compatibility
  os_minimum_version               = "13"
  min_android_security_patch_level = "January 1, 2024"

  # Balanced password requirements
  password_required                          = true
  password_minimum_length                    = 8
  password_minimum_letter_characters         = 1
  password_minimum_numeric_characters        = 1
  password_required_type                     = "alphanumeric"
  password_minutes_of_inactivity_before_lock = 15
  password_expiration_days                   = 180
  password_previous_password_count_to_block  = 3

  # Storage encryption required
  storage_require_encryption = true

  # Scheduled actions with grace periods for user adaptation
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type        = "notification"
          grace_period_hours = 48
        },
        {
          action_type        = "block"
          grace_period_hours = 168 # 7 days
        }
      ]
    }
  ]

  # Assignment to enterprise device group
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000" # Replace with actual enterprise device group ID
    }
  ]
}