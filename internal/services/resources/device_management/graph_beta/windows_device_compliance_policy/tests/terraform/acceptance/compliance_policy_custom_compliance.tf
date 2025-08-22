resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "custom_compliance" {
  display_name    = "tf-test-custom-compliance"
  description     = "tf-test-custom-compliance"
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