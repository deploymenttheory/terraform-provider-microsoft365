resource "microsoft365_graph_beta_device_and_app_management_settings_catalog" "test_macOS" {
  name               = "Test Settings Catalog Profile - macOS"
  description        = ""
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({

    "settingsDetails" : [
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

resource "microsoft365_graph_beta_device_and_app_management_settings_catalog" "test_fxlogix" {
  name               = "Windows Multisession Settings Catalog Profile - fxlogix"
  description        = ""
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({

    "settingsDetails" : [
      {
        "id" : "0",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix_cleanupinvalidsessions_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix_cleanupinvalidsessions",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "1",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [
              {
                "simpleSettingValue" : {
                  "value" : "thing3",
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                  "settingValueTemplateReference" : null
                },
                "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdcachedirectory_ccdcachedirectory",
                "settingInstanceTemplateReference" : null
              }
            ],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdcachedirectory_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdcachedirectory",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "2",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [
              {
                "simpleSettingValue" : {
                  "value" : "thing2",
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                  "settingValueTemplateReference" : null
                },
                "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdproxydirectory_ccdproxydirectory",
                "settingInstanceTemplateReference" : null
              }
            ],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdproxydirectory_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdproxydirectory",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "3",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [
              {
                "simpleSettingValue" : {
                  "value" : "thing",
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                  "settingValueTemplateReference" : null
                },
                "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdwritecachedirectory_ccdwritecachedirectory",
                "settingInstanceTemplateReference" : null
              }
            ],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdwritecachedirectory_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdwritecachedirectory",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "4",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingadcomputergroupprocess_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingadcomputergroupprocess",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "5",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingdriverinterface_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingdriverinterface",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "6",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingfrxlauncher_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingfrxlauncher",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "7",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingieplugin_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingieplugin",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "8",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingjavaruleeditor_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingjavaruleeditor",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "9",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingprofileconfigurationtool_0",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingprofileconfigurationtool",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "10",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingrulecompilation_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingrulecompilation",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "11",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingruleeditor_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingruleeditor",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "12",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingsearchplugin_0",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingsearchplugin",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "13",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingsearchroaming_0",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingsearchroaming",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "14",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingservices_0",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~logging_loggingservices",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "15",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [
              {
                "simpleSettingValue" : {
                  "value" : "",
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                  "settingValueTemplateReference" : null
                },
                "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdlocations_profilesccdlocations",
                "settingInstanceTemplateReference" : null
              }
            ],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdlocations_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdlocations",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "16",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdmaxcachesizeinmbs_0",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdmaxcachesizeinmbs",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "17",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [
              {
                "simpleSettingValue" : {
                  "value" : 1,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                  "settingValueTemplateReference" : null
                },
                "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdunregistertimeout_profilesccdunregistertimeout",
                "settingInstanceTemplateReference" : null
              }
            ],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdunregistertimeout_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdunregistertimeout",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "18",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesclearcacheonforcedunregister_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesclearcacheonforcedunregister",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "19",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesclearcacheonlogoff_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesclearcacheonlogoff",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "20",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profileshealthyprovidersrequiredforregister_0",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profileshealthyprovidersrequiredforregister",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "21",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [
              {
                "simpleSettingValue" : {
                  "value" : 1,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                  "settingValueTemplateReference" : null
                },
                "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profileshealthyprovidersrequiredforunregister_profileshealthyprovidersrequiredforunregister",
                "settingInstanceTemplateReference" : null
              }
            ],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profileshealthyprovidersrequiredforunregister_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profileshealthyprovidersrequiredforunregister",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "22",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix_roamrecyclebin_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix_roamrecyclebin",
          "settingInstanceTemplateReference" : null
        }
      },
      {
        "id" : "23",
        "settingInstance" : {
          "choiceSettingValue" : {
            "children" : [],
            "value" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix_vhdcompactdisk_1",
            "settingValueTemplateReference" : null
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_config_fslogixv1~policy~fslogix_vhdcompactdisk",
          "settingInstanceTemplateReference" : null
        }
      }
    ]
  })

  assignments = {
    all_devices = false
    # all_devices_filter_type = "exclude"
    # all_devices_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"

    all_users = true
    all_users_filter_type = "include"
    all_users_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"

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