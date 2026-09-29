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
