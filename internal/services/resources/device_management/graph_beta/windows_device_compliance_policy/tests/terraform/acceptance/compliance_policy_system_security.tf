resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "system_security" {
  display_name    = "tf-test-system-security"
  description     = "tf-test-system-security"
  role_scope_tag_ids = ["0"]

  # System Security Settings
  system_security = {
    active_firewall_required = true
    anti_spyware_required = true
    antivirus_required = true
    configuration_manager_compliance_required = true
    defender_enabled = true
    defender_version = "1.0.0.0"
    device_threat_protection_enabled = true
    device_threat_protection_required_security_level = "medium"
    password_block_simple = true
    password_minimum_character_set_count = 4
    password_required = true
    password_required_to_unlock_from_idle = true
    password_required_type = "alphanumeric"
    rtp_enabled = true
    signature_out_of_date = true
    storage_require_encryption = true
    tpm_required = true
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type = "block"
          grace_period_hours = 12
        },
        {
          action_type = "notification"
          grace_period_hours = 24
          notification_template_id = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = ["a77240dc-2827-47af-8fcb-e209a67e176a"]
        },
        {
          action_type = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}