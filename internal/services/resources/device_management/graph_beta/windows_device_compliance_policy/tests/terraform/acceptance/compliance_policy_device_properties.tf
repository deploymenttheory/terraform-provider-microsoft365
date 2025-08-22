resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "device_properties" {
  display_name    = "tf-test-device-properties"
  description     = "tf-test-device-properties"
  role_scope_tag_ids = ["0"]

  device_properties = {
    os_minimum_version = "10.0.22631.5768"
    os_maximum_version = "10.0.26100.9999"
    mobile_os_minimum_version = "10.0.22631.5768"
    mobile_os_maximum_version = "10.0.26100.9999"
    valid_operating_system_build_ranges = [
      {
        description = "Windows 11 24H2"
        low_os_version = "10.0.26100.4946"
        high_os_version = "10.0.26100.9999"
      },
      {
        description = "Windows 11 23H2"
        low_os_version = "10.0.22631.5768"
        high_os_version = "10.0.22631.9999"
      },
      
    ]
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