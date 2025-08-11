resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "assignments" {
  name         = "Test All Assignment Types Settings Catalog Policy"
  description  = "Settings catalog configuration policy with comprehensive assignments for acceptance testing"
  platforms    = "macOS"
  technologies = ["mdm", "appleRemoteManagement"]

  role_scope_tag_ids = [
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id,
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id
  ]


  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id               = "com.apple.mcx.filevault2_com.apple.mcx.filevault2"
          setting_instance_template_reference = null
          group_setting_collection_value = [
            {
              setting_value_template_reference = null
              children = [
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_defer"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_defer_true"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_deferdontaskatuserlogout"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_deferdontaskatuserlogout_false"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_deferforceatuserloginmaxbypassattempts"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                    setting_value_template_reference = null
                    value                            = 0
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_enable"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_enable_0"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_forceenableinsetupassistant"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_forceenableinsetupassistant_false"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_outputpath"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    setting_value_template_reference = null
                    value                            = "/output/path"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_password"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
                    setting_value_template_reference = null
                    value_state                      = "notEncrypted"
                    value                            = "3669d68b-ea40-4682-abc9-9445f3f6fc7e"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_recoverykeyrotationinmonths"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_recoverykeyrotationinmonths_10"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_showrecoverykey"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_showrecoverykey_true"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_usekeychain"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_usekeychain_false"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_userecoverykey"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_userecoverykey_true"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_userentersmissinginfo"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx.filevault2_userentersmissinginfo_false"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "com.apple.mcx.filevault2_username"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    setting_value_template_reference = null
                    value                            = "username"
                  }
                }
              ]
            }
          ]
        }
        id = "0"
      },
      {
        setting_instance = {
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id               = "com.apple.mcx_com.apple.mcx-fdefilevaultoptions"
          setting_instance_template_reference = null
          group_setting_collection_value = [
            {
              setting_value_template_reference = null
              children = [
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx_dontallowfdedisable"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx_dontallowfdedisable_true"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id               = "com.apple.mcx_dontallowfdeenable"
                  setting_instance_template_reference = null
                  choice_setting_value = {
                    setting_value_template_reference = null
                    children                         = []
                    value                            = "com.apple.mcx_dontallowfdeenable_false"
                  }
                }
              ]
            }
          ]
        }
        id = "1"
      },
      {
        setting_instance = {
          odata_type                          = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id               = "com.apple.security.fderecoverykeyescrow_com.apple.security.fderecoverykeyescrow"
          setting_instance_template_reference = null
          group_setting_collection_value = [
            {
              setting_value_template_reference = null
              children = [
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "com.apple.security.fderecoverykeyescrow_devicekey"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    setting_value_template_reference = null
                    value                            = "some-device-key"
                  }
                },
                {
                  odata_type                          = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                  setting_definition_id               = "com.apple.security.fderecoverykeyescrow_location"
                  setting_instance_template_reference = null
                  simple_setting_value = {
                    odata_type                       = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    setting_value_template_reference = null
                    value                            = "location"
                  }
                }
              ]
            }
          ]
        }
        id = "2"
      }
    ]
  }


  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_3.id
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
    }
  ]


  timeouts = {
    create = "60s"  
    read   = "30s"  
    update = "30s"  
    delete = "20s"
  }
}

