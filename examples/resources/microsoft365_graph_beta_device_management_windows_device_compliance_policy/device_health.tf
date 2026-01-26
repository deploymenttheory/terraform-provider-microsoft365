# Example: Device Health Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "device_health" {
  display_name       = "tf-reg-windows-device-compliance-policy-device-health"
  description        = "tf-reg-windows-device-compliance-policy-device-health"
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
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}
