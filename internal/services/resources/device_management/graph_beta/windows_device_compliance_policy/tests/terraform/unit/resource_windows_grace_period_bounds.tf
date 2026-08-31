resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "grace_period_bounds" {
  display_name       = "unit-test-wdcp-grace-period-bounds"
  description        = "unit-test-wdcp-grace-period-bounds"
  role_scope_tag_ids = ["0"]

  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        # Lower bound
        {
          action_type        = "block"
          grace_period_hours = 0
        },
        # Upper bound: 8760 hours == 365 days
        {
          action_type        = "retire"
          grace_period_hours = 8760
        },
      ]
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }

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
