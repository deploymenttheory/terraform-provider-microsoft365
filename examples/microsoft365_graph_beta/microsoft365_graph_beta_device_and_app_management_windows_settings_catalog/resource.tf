resource "microsoft365_graph_beta_device_and_app_management_windows_settings_catalog" "example_policy" {
  display_name = "Test Settings Catalog Profile"
  description  = "Test settings catalog profile"
  platforms    = "windows10"
  role_scope_tag_ids = [
    "0"
  ]
  settings = [
    {
      odata_type = "#microsoft.graph.deviceManagementConfigurationSetting"
      setting_instance = {
        odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
        setting_definition_id = "device_vendor_msft_bitlocker_allowwarningforotherdiskencryption"
        choice_setting_value = {
          odata_type   = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
          string_value = "device_vendor_msft_bitlocker_allowwarningforotherdiskencryption_1"
        }
      }
    },
    {
      odata_type = "#microsoft.graph.deviceManagementConfigurationSetting"
      setting_instance = {
        odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
        setting_definition_id = "device_vendor_msft_bitlocker_configurerecoverypasswordrotation"
        choice_setting_value = {
          odata_type   = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
          string_value = "device_vendor_msft_bitlocker_configurerecoverypasswordrotation_1"
        }
      }
    },
    {
      odata_type = "#microsoft.graph.deviceManagementConfigurationSetting"
      setting_instance = {
        odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
        setting_definition_id = "device_vendor_msft_policy_config_bits_bandwidththrottlingendtime"
        choice_setting_value = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
          int_value  = 17
        }
      }
    },
    {
      odata_type = "#microsoft.graph.deviceManagementConfigurationSetting"
      setting_instance = {
        odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
        setting_definition_id = "device_vendor_msft_bitlocker_removabledrivesexcludedfromencryption"
        choice_setting_value = {
          odata_type = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
          children = [
            {
              choice_setting_value = {
                string_value = "D:\\"
              }
            },
            {
              choice_setting_value = {
                string_value = "E:\\"
              }
            }
          ]
        }
      }
    },
  ]

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