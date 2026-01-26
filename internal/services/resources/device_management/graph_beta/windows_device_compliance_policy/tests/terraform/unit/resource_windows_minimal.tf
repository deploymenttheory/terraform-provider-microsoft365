resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "minimal" {
  display_name       = "unit-test-wdcp-minimal"
  description        = "unit-test-wdcp-minimal"
  role_scope_tag_ids = ["0"]

  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
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

  # Assignments
  assignments = [
    # Optional: Assignment targeting all devices without filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    },
    # Optional: Assignment targeting all licensed users without filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
    },
    # Optional: Assignment targeting a specific group without filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000006"
      filter_type = "none"
    },
    # Optional: Assignment targeting a specific group without filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000008"
      filter_type = "none"
    },
    # Optional: Exclusion group assignments
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