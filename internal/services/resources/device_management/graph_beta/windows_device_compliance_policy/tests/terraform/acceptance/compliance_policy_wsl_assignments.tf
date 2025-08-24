resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "wsl" {
  display_name       = "tf-test-wsl"
  description        = "tf-test-wsl"
  role_scope_tag_ids = ["0"]

  wsl_distributions = [
    {
      distribution       = "Ubuntu"
      minimum_os_version = "1.0"
      maximum_os_version = "1.0"
    },
    {
      distribution       = "redhat"
      minimum_os_version = "1.0"
      maximum_os_version = "1.0"
    }
  ]

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

  # Assignments
  assignments = [
    # Optional: Assignment targeting all devices with a daily schedule
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with an hourly schedule
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_2.id
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_3.id
      filter_type = "include"

    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_4.id
      filter_type = "exclude"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    },
  ]

}