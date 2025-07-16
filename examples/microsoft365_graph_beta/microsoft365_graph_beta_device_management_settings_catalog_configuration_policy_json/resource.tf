resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "test_macOS" {
  name               = "Test Settings Catalog Profile - macOS"
  description        = ""
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({

    "settings" : [
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.mcx_disableguestaccount_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.mcx_disableguestaccount"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.mcx_enableguestaccount_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.mcx_enableguestaccount"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.mcx_com.apple.mcx-accounts"
        },
        "id" : "0"
      },
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavaccountdescription"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavhostname"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "valueState" : "notEncrypted",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value" : "test-password"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavpassword"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "value" : 1
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavport"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavprincipalurl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.caldav.account_caldavusessl_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavusessl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "{{USERNAME}}"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavusername"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.caldav.account_com.apple.caldav.account"
        },
        "id" : "1"
      },
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavaccountdescription"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavhostname"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "valueState" : "notEncrypted",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value" : "e7776185-0499-4e47-bdf5-1b3bc42ba965"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavpassword"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "value" : 1
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavport"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.carddav.account_carddavusessl_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavusessl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "{{USERNAME}}"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavusername"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.carddav.account_com.apple.carddav.account"
        },
        "id" : "2"
      },
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountdescription"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccounthostname"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "valueState" : "notEncrypted",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value" : "762b8bea-3715-449e-b4cd-abc0cb5e16ad"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountpassword"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.ldap.account_ldapaccountusessl_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountusessl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "{{USERNAME}}"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountusername"
                },
                {
                  "groupSettingCollectionValue" : [
                    {
                      "settingValueTemplateReference" : null,
                      "children" : [
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingdescription"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                          "choiceSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "value" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope_2",
                            "children" : []
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingsearchbase"
                        }
                      ]
                    },
                    {
                      "settingValueTemplateReference" : null,
                      "children" : [
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingdescription"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                          "choiceSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "value" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope_2",
                            "children" : []
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingsearchbase"
                        }
                      ]
                    }
                  ],
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
                  "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.ldap.account_com.apple.ldap.account"
        },
        "id" : "3"
      }
    ]
  })

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
