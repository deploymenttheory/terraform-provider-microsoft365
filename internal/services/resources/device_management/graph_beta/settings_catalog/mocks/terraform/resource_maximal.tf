resource "microsoft365_graph_beta_device_management_settings_catalog" "maximal" {
  name               = "Maximal Settings Catalog"
  description        = "Maximal settings catalog policy with all options"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  configuration_policy = {
    settings = [
      {
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
    # all_devices_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"

    all_users = false
    # all_users_filter_type = "include"
    # all_users_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"

    include_groups = [
      {
        group_id                   = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
      },
      {
        group_id                   = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
      },
    ]

    exclude_group_ids = [
      "b8c661c2-fa9a-4351-af86-adc1729c343f",
      "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f",
    ]
  }

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}