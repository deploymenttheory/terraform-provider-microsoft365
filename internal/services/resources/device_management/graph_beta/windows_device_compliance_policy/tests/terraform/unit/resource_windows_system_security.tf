resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "system_security" {
  display_name       = "unit-test-wdcp-system-security"
  description        = "unit-test-wdcp-system-security"
  role_scope_tag_ids = ["0"]

  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
  }

  system_security = {
    active_firewall_required                         = true
    anti_spyware_required                            = true
    antivirus_required                               = true
    configuration_manager_compliance_required        = false
    defender_enabled                                 = true
    password_block_simple                            = true
    password_minimum_character_set_count             = 3
    password_minutes_of_inactivity_before_lock       = 15
    password_required                                = true
    password_required_to_unlock_from_idle            = true
    password_required_type                           = "alphanumeric"
    rtp_enabled                                      = true
    signature_out_of_date                            = false
    storage_require_encryption                       = true
    tpm_required                                     = true
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = "00000000-0000-0000-0000-000000000001"
          notification_message_cc_list = ["00000000-0000-0000-0000-000000000002", "00000000-0000-0000-0000-000000000003"]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000006"
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000008"
      filter_type = "none"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000010"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000011"
    },
  ]
}
