resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "custom_compliance" {
  display_name = "Windows 10/11 - Custom Compliance Policy"
  description  = "Windows device compliance policy with custom compliance script"

  # Password requirements
  password_required                     = true
  password_block_simple                 = true
  password_required_to_unlock_from_idle = true
  password_minimum_length               = 8
  password_required_type                = "alphanumeric"

  # Security requirements
  storage_require_encryption = true
  active_firewall_required   = true
  tpm_required               = true
  antivirus_required         = true
  anti_spyware_required      = true

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

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Non-compliance actions
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "retire"
          grace_period_hours       = 1440
          notification_template_id = ""
        },
        {
          action_type              = "notification"
          grace_period_hours       = 120
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = ["00000000-0000-0000-0000-000000000000",
          "00000000-0000-0000-0000-000000000000"]
        },
        {
          action_type              = "block"
          grace_period_hours       = 1152
          notification_template_id = "00000000-0000-0000-0000-000000000000"
        },
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000",
          "00000000-0000-0000-0000-000000000000"]
        }
      ]
    }
  ]

  # Assignments
  assignments = [
    # Optional: Assignment targeting all devices with a daily schedule
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with an hourly schedule
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]
}
