// Example 1: Custom Compliance Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "custom_compliance" {
  display_name       = "tf-reg-example-windows-device-compliance-policy-custom-compliance"
  description        = "tf-reg-example-windows-device-compliance-policy-custom-compliance"
  role_scope_tag_ids = ["0"]

  # Custom compliance script
  custom_compliance_required = true
  device_compliance_policy_script = {
    device_compliance_script_id = microsoft365_graph_beta_device_management_windows_device_compliance_script.example.id
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
            },
            {
              "Language" : "de_DE",
              "Title" : "BIOS-Version muss auf mindestens 2.3 aktualisiert werden. Der erkannte Wert lautet {ActualValue}.",
              "Description" : "BIOS muss aktualisiert werden. Bitte beziehen Sie sich auf den obigen Link"
            }
          ]
        },
        {
          "SettingName" : "TPMChipPresent",
          "Operator" : "IsEquals",
          "DataType" : "Boolean",
          "Operand" : true,
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "TPM chip must be enabled.",
              "Description" : "TPM chip must be enabled. Please refer to the link above"
            },
            {
              "Language" : "de_DE",
              "Title" : "TPM-Chip muss aktiviert sein.",
              "Description" : "TPM-Chip muss aktiviert sein. Bitte beziehen Sie sich auf den obigen Link"
            }
          ]
        },
        {
          "SettingName" : "Manufacturer",
          "Operator" : "IsEquals",
          "DataType" : "String",
          "Operand" : "Microsoft Corporation",
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "Only Microsoft devices are supported.",
              "Description" : "You are not currently using a Microsoft device."
            },
            {
              "Language" : "de_DE",
              "Title" : "Nur Microsoft-Geräte werden unterstützt.",
              "Description" : "Sie verwenden derzeit kein Microsoft-Gerät."
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

// Example 2: Device Health Policy
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

// Example 3: Device Properties Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "device_properties" {
  display_name       = "tf-reg-windows-device-compliance-policy-device-properties"
  description        = "tf-reg-windows-device-compliance-policy-device-properties"
  role_scope_tag_ids = ["0"]

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

// Example 4: Microsoft Defender for Endpoint Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "microsoft_defender_for_endpoint" {
  display_name       = "tf-reg-windows-device-compliance-policy-microsoft-defender-for-endpoint"
  description        = "tf-reg-windows-device-compliance-policy-microsoft-defender-for-endpoint"
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

// Example 5: System Security Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "system_security" {
  display_name       = "tf-reg-windows-device-compliance-policy-system-security"
  description        = "tf-reg-windows-device-compliance-policy-system-security"
  role_scope_tag_ids = ["0"]

  # System Security Settings
  system_security = {
    active_firewall_required                         = true
    anti_spyware_required                            = true
    antivirus_required                               = true
    configuration_manager_compliance_required        = true
    defender_enabled                                 = true
    defender_version                                 = "1.0.0.0"
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
    password_block_simple                            = true
    password_minimum_character_set_count             = 4
    password_required                                = true
    password_required_to_unlock_from_idle            = true
    password_required_type                           = "alphanumeric"
    rtp_enabled                                      = true
    signature_out_of_date                            = true
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

// Example 6: WSL Policy with assignments
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "wsl_assignments" {
  display_name       = "tf-reg-windows-device-compliance-policy-wsl-assignments"
  description        = "tf-reg-windows-device-compliance-policy-wsl-assignments"
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
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
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