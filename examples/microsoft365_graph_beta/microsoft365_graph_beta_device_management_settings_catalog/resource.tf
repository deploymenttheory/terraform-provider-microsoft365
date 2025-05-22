resource "microsoft365_graph_beta_device_management_settings_catalog" "windows_hello_for_business" {
  name               = "pure_hcl_test"
  description        = "18.01.2025\nContext: User\n\nWindows Hello for Business is Set here rather than in enrollment blade globally to allow for targetting specific groups or users.\n\nRef: https://deviceadvice.io/2020/06/22/how-to-set-up-windows-hello-for-business-for-cloud-only-devices/"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  configuration_policy = {
    settings = [
      {
        id = "0"
        setting_instance = {
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id               = "device_vendor_msft_passportforwork_biometrics_usebiometrics"
          setting_instance_template_reference = null
          choice_setting_value = {
            value                            = "device_vendor_msft_passportforwork_biometrics_usebiometrics_true"
            setting_value_template_reference = null
            children                         = []
          }
        }
      },
      {
        id = "1"
        setting_instance = {
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id               = "device_vendor_msft_passportforwork_biometrics_facialfeaturesuseenhancedantispoofing"
          setting_instance_template_reference = null
          choice_setting_value = {
            value                            = "device_vendor_msft_passportforwork_biometrics_facialfeaturesuseenhancedantispoofing_true"
            setting_value_template_reference = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        setting_instance = {
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}"
          setting_instance_template_reference = null
          group_setting_collection_value = [
            {
              setting_value_template_reference = null
              children = [
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_digits"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    value                            = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_digits_0"
                    setting_value_template_reference = null
                    children                         = []
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_enablepinrecovery"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    value                            = "user_vendor_msft_passportforwork_{tenantid}_policies_enablepinrecovery_true"
                    setting_value_template_reference = null
                    children                         = []
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_expiration"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                    setting_value_template_reference = null
                    value                            = 0
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_lowercaseletters"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    value                            = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_lowercaseletters_0"
                    setting_value_template_reference = null
                    children                         = []
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_maximumpinlength"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                    setting_value_template_reference = null
                    value                            = 12
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_minimumpinlength"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                    setting_value_template_reference = null
                    value                            = 6
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_history"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                    setting_value_template_reference = null
                    value                            = 2
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_requiresecuritydevice"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    value                            = "user_vendor_msft_passportforwork_{tenantid}_policies_requiresecuritydevice_true"
                    setting_value_template_reference = null
                    children                         = []
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_specialcharacters"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    value                            = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_specialcharacters_2"
                    setting_value_template_reference = null
                    children                         = []
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_uppercaseletters"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    value                            = "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_uppercaseletters_0"
                    setting_value_template_reference = null
                    children                         = []
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  }

  assignments = {
    all_devices = false
    # all_devices_filter_type = "exclude"
    # all_devices_filter_id   = "11111111-2222-3333-4444-555555555555"

    all_users = false
    # all_users_filter_type = "include"
    # all_users_filter_id   = "11111111-2222-3333-4444-555555555555"

    include_groups = [
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
      },
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
      },
    ]

    exclude_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555",
    ]
  }

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}