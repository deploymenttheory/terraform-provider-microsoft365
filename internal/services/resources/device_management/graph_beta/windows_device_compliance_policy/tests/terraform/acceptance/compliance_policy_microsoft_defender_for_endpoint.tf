resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "microsoft_defender_for_endpoint" {
  display_name       = "tf-test-microsoft-defender-for-endpoint"
  description        = "tf-test-microsoft-defender-for-endpoint"
  role_scope_tag_ids = ["0"]

  # Microsoft Defender for Endpoint Settings
  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
  }

  # Scheduled Actions for Rule
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
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = ["a77240dc-2827-47af-8fcb-e209a67e176a"]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}