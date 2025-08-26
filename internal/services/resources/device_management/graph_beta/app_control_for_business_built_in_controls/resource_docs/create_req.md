Request URL
https://graph.microsoft.com/beta/deviceManagement/configurationPolicies
Request Method
POST

audit mode

{
    "creationSource":null,"name":"unit-test-app-control-for-business-built-in-controls-maximal","description":"unit-test-app-control-for-business-built-in-controls-maximal","platforms":"windows10",
    "technologies":"mdm",
    "roleScopeTagIds":["0","1","2"],
    "settings":[
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions",
                "choiceSettingValue":
                {
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                    "value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_selected",
                    "children":[
                        {
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
                            "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls",
                            "groupSettingCollectionValue":[
                                {
                                    "children":[
                                        {
                                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                                            "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control",
                                            "choiceSettingValue":
                                            {
                                                "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                                                "value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_1",
                                                "children":[]
                                            }
                                        },
                                        {
                                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance",
                                            "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps",
                                            "choiceSettingCollectionValue":[
                                                {
                                                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                                                    "value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_1","children":[]
                                                },
                                                {
                                                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                                                    "value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_0","children":[]
                                                }
                                            ]
                                        }
                                    ]
                                }
                            ]
                        }
                    ],
                    "settingValueTemplateReference":{
                        "settingValueTemplateId":"b28c7dc4-c7b2-4ce2-8f51-6ebfd3ea69d3"
                    }
                },
                "settingInstanceTemplateReference":
                {
                    "settingInstanceTemplateId":"1de98212-6949-42dc-a89c-e0ff6e5da04b"
                }
            }
        }
    ],
    "templateReference":
    {
        "templateId":"4321b946-b76b-4450-8afd-769c08b16ffc_1",
        "templateFamily":"endpointSecurityApplicationControl",
        "templateDisplayName":"App Control for Business",
        "templateDisplayVersion":"Version 1"
    }
}

enforce mode

{
    "creationSource":null,
    "name":"unit-test-app-control-for-business-built-in-controls-maximal",
    "description":"unit-test-app-control-for-business-built-in-controls-maximal",
    "platforms":"windows10",
    "technologies":"mdm",
    "roleScopeTagIds":["0","1","2"],
    "settings":[
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions",
                "choiceSettingValue":
                {
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                    "value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_selected",
                    "children":
                    [
                        {
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
                            "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls",
                            "groupSettingCollectionValue":
                            [
                                {
                                    "children":[
                                    {
                                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_0","children":[]
                                    }
                                },
                                {
                                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance",
                                    "settingDefinitionId":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps",
                                    "choiceSettingCollectionValue":[
                                        {
                                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_1","children":[]
                                        },
                                        {
                                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                                            "value":"device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_0","children":[]
                                        }
                                    ]
                                }
                            ]
                        }
                    ]
                }
            ],
            "settingValueTemplateReference":
            {
                "settingValueTemplateId":"b28c7dc4-c7b2-4ce2-8f51-6ebfd3ea69d3"
                }
            },
            "settingInstanceTemplateReference":
            {
                "settingInstanceTemplateId":"1de98212-6949-42dc-a89c-e0ff6e5da04b"
                }
            }
        }
    ],
    "templateReference":
    {
        "templateId":"4321b946-b76b-4450-8afd-769c08b16ffc_1","templateFamily":"endpointSecurityApplicationControl",
        "templateDisplayName":"App Control for Business",
        "templateDisplayVersion":"Version 1"
    }
}

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp
Request Method
GET

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
    "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d",
    "availableVersion": "1.93.102.0",
    "managedInstaller": "disabled",
    "managedInstallerConfiguredDateTime": null
}

then to enable it

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller
Request Method
POST

empty body.

then get again

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp
Request Method
GET

now it's enabled

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
    "id": "54fac284-7866-43e5-860a-9c8e10fa3d7d",
    "availableVersion": "1.93.102.0",
    "managedInstaller": "enabled",
    "managedInstallerConfiguredDateTime": "8/23/2025 7:51:54 AM +00:00"
}
post again to disable it. 

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller
Request Method
POST

what it all means - https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/configure-authorized-apps-deployed-with-a-managed-installer