# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "minimal" {
  display_name = "macOS Minimal Compliance Policy"
  description  = "Minimal macOS device compliance policy with basic security requirements"

  # Basic security requirements
  password_required          = true
  storage_require_encryption = true
  firewall_enabled           = true

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
resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "advanced" {
  display_name = "macOS Advanced Compliance Policy"
  description  = "Advanced macOS device compliance policy with strict security requirements"

  # Strict password requirements
  password_required                          = true
  password_block_simple                      = true
  password_minimum_length                    = 12
  password_minimum_character_set_count       = 4
  password_required_type                     = "alphanumeric"
  password_expiration_days                   = 60
  password_previous_password_block_count     = 10
  password_minutes_of_inactivity_before_lock = 5

  # Strict OS version requirements
  os_minimum_version       = "14.0"
  os_minimum_build_version = "23A344"

  # Maximum security settings
  system_integrity_protection_enabled                = true
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "high"
  advanced_threat_protection_required_security_level = "high"
  storage_require_encryption                         = true
  gatekeeper_allowed_app_source                      = "macAppStore"

  # Strict firewall settings
  firewall_enabled             = true
  firewall_block_all_incoming  = true
  firewall_enable_stealth_mode = true

  # Scheduled actions with aggressive enforcement
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "426e6351-c6ff-44d3-910d-8b937ee30bdd"
        },
        {
          action_type        = "block"
          grace_period_hours = 1
        },
        {
          action_type        = "retire"
          grace_period_hours = 24
        }
      ]
    }
  ]

  # Target specific high-security groups
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "high-security-macos-devices-group-id"
    }
  ]
}

resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "comprehensive" {
  display_name = "macOS Comprehensive Compliance Policy"
  description  = "Comprehensive macOS device compliance policy with all available security settings"

  # Password requirements
  password_required                          = true
  password_block_simple                      = true
  password_minimum_length                    = 8
  password_minimum_character_set_count       = 3
  password_required_type                     = "alphanumeric"
  password_expiration_days                   = 90
  password_previous_password_block_count     = 5
  password_minutes_of_inactivity_before_lock = 15

  # OS version requirements
  os_minimum_version       = "13.0"
  os_maximum_version       = "14.0"
  os_minimum_build_version = "22A380"
  os_maximum_build_version = "23A344"

  # Security requirements
  system_integrity_protection_enabled                = true
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "medium"
  advanced_threat_protection_required_security_level = "low"
  storage_require_encryption                         = true
  gatekeeper_allowed_app_source                      = "macAppStoreAndIdentifiedDevelopers"

  # Firewall requirements
  firewall_enabled             = true
  firewall_block_all_incoming  = false
  firewall_enable_stealth_mode = true

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Scheduled actions for rules
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "426e6351-c6ff-44d3-910d-8b937ee30bdd"
          notification_message_cc_list = [
            "aa856a09-cf0c-4b31-a315-cb53251e54d8",
            "a77240dc-2827-47af-8fcb-e209a67e176a"
          ]
        },
        {
          action_type              = "notification"
          grace_period_hours       = 24
          notification_template_id = "bbf43ceb-5e68-428b-8ad3-00c9efb54210"
          notification_message_cc_list = [
            "91710c72-1358-4438-b0b2-70eb32b542dd",
            "aa856a09-cf0c-4b31-a315-cb53251e54d8"
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
      filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
      filter_type = "include"
    },
    # Assignment targeting all licensed users with an exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
      filter_type = "exclude"
    },
    # Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
      filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
      filter_type = "include"
    },
    # Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
      filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
      filter_type = "exclude"
    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "b8c661c2-fa9a-4351-af86-adc1729c343f"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
    }
  ]
}