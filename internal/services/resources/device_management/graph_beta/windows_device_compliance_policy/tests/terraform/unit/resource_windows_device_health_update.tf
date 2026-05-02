resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "maximal" {
  display_name       = "unit-test-wdcp-maximal"
  description        = "unit-test-wdcp-maximal"
  role_scope_tag_ids = ["0"]

  device_health = {
    bit_locker_enabled     = false
    secure_boot_enabled    = false
    code_integrity_enabled = false
  }

  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
  }

  device_properties = {
    os_minimum_version        = "10.0.22631.5768"
    os_maximum_version        = "10.0.26100.9999"
    mobile_os_minimum_version = "10.0.22631.5768"
    mobile_os_maximum_version = "10.0.26100.9999"
    valid_operating_system_build_ranges = [
      {
        description     = "Windows 11 24H2"
        low_os_version  = "10.0.26100.4946"
        high_os_version = "10.0.26100.9999"
      },
      {
        description     = "Windows 11 23H2"
        low_os_version  = "10.0.22631.5768"
        high_os_version = "10.0.22631.9999"
      },
    ]
  }

  custom_compliance_required = true
  device_compliance_policy_script = {
    device_compliance_script_id = "00000000-0000-0000-0000-000000000001"
    rules_content = jsonencode({
      "Rules" : [
        {
          "SettingName" : "BiosVersion",
          "Operator" : "GreaterEquals",
          "DataType" : "Version",
          "Operand" : "2.3",
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "BIOS Version needs to be upgraded to at least 2.3. Value discovered was {ActualValue}.",
              "Description" : "BIOS must be updated. Please refer to the link above"
            }
          ]
        }
      ]
    })
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
