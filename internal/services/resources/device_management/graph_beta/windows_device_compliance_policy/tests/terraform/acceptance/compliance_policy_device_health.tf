resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "device_health" {
  display_name       = "tf-test-device-health"
  description        = "tf-test-device-health"
  role_scope_tag_ids = ["0"]

  # Device Health Settings
  device_health = {
    bit_locker_enabled     = true
    secure_boot_enabled    = true
    code_integrity_enabled = true
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