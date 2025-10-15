validate configuration categories

Request URL
https://graph.microsoft.com/beta/deviceManagement/configurationCategories?&$filter=(platforms%20has%20%27none%27)%20and%20(technologies%20has%20%27mobileApplicationManagement%27)%20and%20(rootCategoryId%20eq%20%27a25a7a02-4bac-411b-9d02-10cb3297cb17%27)
Request Method
GET

resp 

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationCategories",
    "value": [
        {
            "id": "ef6a4e8c-07b2-4f55-9e94-5701cb2268b1",
            "description": "Microsoft Edge Edge Workspaces settings",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Edge Workspaces settings",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "66615d2a-fec9-47f1-8eaf-9813e30cc023",
            "description": "Microsoft Edge\\Extensions",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Extensions",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8",
            "description": "Microsoft Edge\\ Uncategorized",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Uncategorized",
            "platforms": "macOS,windows10",
            "technologies": "mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "45a89c1f-0a34-4f78-b28f-d30b623fa423",
            "description": "Microsoft Edge\\ Identity and sign-in",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Identity and sign-in",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "3ba8106d-4b2f-4775-939d-1cc8703a41dc",
            "description": "Microsoft Edge\\Password manager and protection",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Password manager and protection",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "fe845e81-5993-4a65-b22a-decfc5928c65",
            "description": "Microsoft Edge\\Proxy server",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Proxy server",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "08c5f391-e156-4a72-bbb9-3670f2f63a56",
            "description": "Microsoft Edge\\SmartScreen settings",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "SmartScreen settings",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "description": "Microsoft Edge",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Microsoft Edge",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "00000000-0000-0000-0000-000000000000",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": [
                "3abaf4c2-d5db-4b3b-a461-b1a208231b36",
                "ef6a4e8c-07b2-4f55-9e94-5701cb2268b1",
                "66615d2a-fec9-47f1-8eaf-9813e30cc023",
                "6d529e48-5477-4ceb-8ff7-c6e959a0e24f",
                "ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8",
                "45a89c1f-0a34-4f78-b28f-d30b623fa423",
                "526e363a-84db-4256-a13c-e01c8c646e26",
                "3ba8106d-4b2f-4775-939d-1cc8703a41dc",
                "fe845e81-5993-4a65-b22a-decfc5928c65",
                "08c5f391-e156-4a72-bbb9-3670f2f63a56",
                "92d69c43-75ac-49b1-a3ef-9350079eef86",
                "ef8760ac-a77c-4055-a812-a95bfbf9c00a",
                "05811ceb-2954-426c-8afa-2a53f02480cc",
                "16ea64a1-563e-43cc-b34a-728c8e7cd13c",
                "00d7396c-cadc-4d29-86ba-fe4df2ecb110",
                "08677354-6f67-455e-a430-4d8d2fbabe84",
                "5e8e9c7f-1988-45cd-b5ca-78d939e3d49e",
                "76e34834-6d47-4e06-b14c-aa2888cdce27",
                "d17b08e6-de3b-445b-ab14-1d47e62efdcf",
                "3fbd3b29-bafd-4adf-89e4-3be612dee275",
                "b3c8c6d9-28bb-475a-9353-4a0e657b33c7",
                "1043e7ed-8651-44b2-b918-7230c0b75a6c",
                "fddc444c-3591-4a50-865b-d8993b798e12",
                "eb6409fc-fb52-413d-ae4b-eff017b52b30",
                "8bcf8b08-35a3-49b7-8760-5fe3b767d6a6",
                "d9678af8-c0c7-401a-a0a5-3e7f5b1253ce",
                "2e241b46-e5ae-41d7-a559-fe40819e86f4",
                "8aa3383a-efac-4ec4-841d-06e3e18646d8",
                "3edb2860-b77b-4240-af16-fb34d45d6ba1",
                "ae78ab75-2d0d-418c-be6f-9e64642de4e2",
                "43057320-7058-46d5-86f9-a56c80bbf8b9",
                "5bd0eaf1-1818-44e8-9168-fc75c5739cc8",
                "81c518f1-522e-4957-b850-e8a66d2ab215",
                "dfab5866-1712-4bbf-8edf-5b080b315b9b",
                "fb1e99d0-b921-4b19-9842-17e3e7987528",
                "c6099521-a05f-480a-8562-7e71318e2cda"
            ]
        },
        {
            "id": "92d69c43-75ac-49b1-a3ef-9350079eef86",
            "description": "Microsoft Edge\\Content settings",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Content settings",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "ef8760ac-a77c-4055-a812-a95bfbf9c00a",
            "description": "Microsoft Edge\\Native Messaging",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Native Messaging",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "05811ceb-2954-426c-8afa-2a53f02480cc",
            "description": "Microsoft Edge\\ Permit or deny screen capture",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Permit or deny screen capture",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "00d7396c-cadc-4d29-86ba-fe4df2ecb110",
            "description": "Microsoft Edge\\Startup, home page and new tab page",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Startup, home page and new tab page",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "fddc444c-3591-4a50-865b-d8993b798e12",
            "description": "Microsoft Edge\\Cast",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Cast",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "8bcf8b08-35a3-49b7-8760-5fe3b767d6a6",
            "description": "Microsoft Edge Immersive Reader settings",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Immersive Reader settings",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "8aa3383a-efac-4ec4-841d-06e3e18646d8",
            "description": "Microsoft Edge\\Default search provider",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Default search provider",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "3edb2860-b77b-4240-af16-fb34d45d6ba1",
            "description": "Microsoft Edge\\Performance",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Performance",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "ae78ab75-2d0d-418c-be6f-9e64642de4e2",
            "description": "Microsoft Edge\\Sleeping Tabs settings",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Sleeping Tabs settings",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "43057320-7058-46d5-86f9-a56c80bbf8b9",
            "description": "Microsoft Edge\\ Private Network Request Settings",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Private Network Request Settings",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        },
        {
            "id": "c6099521-a05f-480a-8562-7e71318e2cda",
            "description": "Microsoft Edge\\Printing",
            "categoryDescription": null,
            "helpText": null,
            "name": null,
            "displayName": "Printing",
            "platforms": "macOS,windows10",
            "technologies": "mdm,mobileApplicationManagement",
            "settingUsage": "configuration",
            "parentCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "rootCategoryId": "a25a7a02-4bac-411b-9d02-10cb3297cb17",
            "childCategoryIds": []
        }
    ]
}

get all mobile app management settings

Request URL
https://graph.microsoft.com/beta/deviceManagement/configurationSettings?&$filter=categoryId%20eq%20%27ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8%27%20and%20visibility%20has%20%27settingsCatalog%27%20and%20(applicability/platform%20has%20%27none%27)%20and%20(applicability/technologies%20has%20%27mobileApplicationManagement%27)
Request Method
GET

resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationSettings",
    "value": [
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition",
            "accessTypes": "none",
            "keywords": [
                "Uncategorized"
            ],
            "infoUrls": [],
            "occurrence": null,
            "baseUri": "",
            "offsetUri": "com.microsoft.edge.ImportSearchEngine_recommended",
            "rootDefinitionId": "com.microsoft.edge.mamedgeappconfigsettings.importsearchengine_recommended",
            "categoryId": "ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8",
            "settingUsage": "configuration",
            "uxBehavior": "toggle",
            "visibility": "settingsCatalog",
            "riskLevel": "low",
            "id": "com.microsoft.edge.mamedgeappconfigsettings.importsearchengine_recommended",
            "description": "Allows users to import search engine settings from another browser into Microsoft Edge.\n\nIf you enable, this policy, the option to import search engine settings is automatically selected.\n\nIf you disable this policy, search engine settings aren't imported at first run, and users can't import them manually.\n\nIf you don't configure this policy, search engine settings are imported at first run, and users can choose whether to import this data manually during later browsing sessions.\n\nYou can set this policy as a recommendation. This means that Microsoft Edge imports search engine settings on first run, but users can select or clear the **search engine** option during manual import.\n\n**Note**: This policy currently manages importing from Internet Explorer (on Windows 7, 8, and 10).",
            "helpText": null,
            "name": "ImportSearchEngine_recommended",
            "displayName": "Allow importing of search engine settings (users can override)",
            "version": "638951450276131607",
            "defaultOptionId": null,
            "applicability": {
                "@odata.type": "#microsoft.graph.deviceManagementConfigurationApplicationSettingApplicability",
                "description": null,
                "platform": "macOS,windows10",
                "deviceMode": "none",
                "technologies": "mobileApplicationManagement"
            },
            "referredSettingInformationList": [],
            "options": [
                {
                    "itemId": "com.microsoft.edge.mamedgeappconfigsettings.importsearchengine_recommended_false",
                    "description": null,
                    "helpText": null,
                    "name": "Disabled",
                    "displayName": "Disabled",
                    "optionValue": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                        "settingValueTemplateReference": null,
                        "value": "false"
                    },
                    "dependentOn": [],
                    "dependedOnBy": []
                },
                {
                    "itemId": "com.microsoft.edge.mamedgeappconfigsettings.importsearchengine_recommended_true",
                    "description": null,
                    "helpText": null,
                    "name": "Enabled",
                    "displayName": "Enabled",
                    "optionValue": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                        "settingValueTemplateReference": null,
                        "value": "true"
                    },
                    "dependentOn": [],
                    "dependedOnBy": []
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingDefinition",
            "accessTypes": "none",
            "keywords": [
                "Uncategorized"
            ],
            "infoUrls": [],
            "occurrence": null,
            "baseUri": "",
            "offsetUri": "com.microsoft.edge.InternetExplorerIntegrationWindowOpenWidthAdjustment",
            "rootDefinitionId": "com.microsoft.edge.mamedgeappconfigsettings.internetexplorerintegrationwindowopenwidthadjustment",
            "categoryId": "ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8",
            "settingUsage": "configuration",
            "uxBehavior": "smallTextBox",
            "visibility": "settingsCatalog",
            "riskLevel": "low",
            "id": "com.microsoft.edge.mamedgeappconfigsettings.internetexplorerintegrationwindowopenwidthadjustment",
            "description": "This setting lets you specify a custom adjustment to the width of popup windows generated via window.open from the Internet Explorer mode site.\n\nIf you configure this policy, Microsoft Edge will add the adjustment value to the width, in pixels. The exact difference depends on the UI configuration of both IE and Edge, but a typical difference is 4.\n\nIf you disable or don't configure this policy, Microsoft Edge will treat IE mode window.open the same as Edge mode window.open in window width calculations.",
            "helpText": null,
            "name": "InternetExplorerIntegrationWindowOpenWidthAdjustment",
            "displayName": "Configure the pixel adjustment between window.open widths sourced from IE mode pages vs. Edge mode pages",
            "version": "638951450276131607",
            "applicability": {
                "@odata.type": "#microsoft.graph.deviceManagementConfigurationApplicationSettingApplicability",
                "description": null,
                "platform": "windows10",
                "deviceMode": "none",
                "technologies": "mobileApplicationManagement"
            },
            "referredSettingInformationList": [],
            "valueDefinition": {
                "@odata.type": "#microsoft.graph.deviceManagementConfigurationIntegerSettingValueDefinition",
                "maximumValue": 2147483647,
                "minimumValue": -2147483648
            },
            "defaultValue": {
                "@odata.type": "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                "settingValueTemplateReference": null,
                "value": 0
            },
            "dependentOn": [],
            "dependedOnBy": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition",
            "accessTypes": "none",
            "keywords": [
                "Uncategorized"
            ],
            "infoUrls": [],
            "occurrence": null,
            "baseUri": "",
            "offsetUri": "com.microsoft.edge.ForceSync",
            "rootDefinitionId": "com.microsoft.edge.mamedgeappconfigsettings.forcesync",
            "categoryId": "ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8",
            "settingUsage": "configuration",
            "uxBehavior": "toggle",
            "visibility": "settingsCatalog",
            "riskLevel": "low",
            "id": "com.microsoft.edge.mamedgeappconfigsettings.forcesync",
            "description": "Forces data synchronization in Microsoft Edge. This policy also prevents the user from turning sync off.\n\nIf you don't configure this policy, users will be able to turn sync on or off. If you enable this policy, users will not be able to turn sync off.\n\nFor this policy to work as intended,\n\"BrowserSignin\" policy must not be configured, or must be set to enabled. If \"BrowserSignin\" is set to disabled, then \"ForceSync\" will not take affect.\n\n\"SyncDisabled\" must not be configured or must be set to False. If this is set to True, \"ForceSync\" will not take affect. If you wish to ensure specific datatypes sync or do not sync, use the \"ForceSyncTypes\" policy and \"SyncTypesListDisabled\" policy.\n\n0 = Do not automatically start sync and show the sync consent (default)\n1 = Force sync to be turned on for Azure AD/Azure AD-Degraded user profile and do not show the sync consent prompt",
            "helpText": null,
            "name": "ForceSync",
            "displayName": "Force synchronization of browser data and do not show the sync consent prompt",
            "version": "638951450276131607",
            "defaultOptionId": null,
            "applicability": {
                "@odata.type": "#microsoft.graph.deviceManagementConfigurationApplicationSettingApplicability",
                "description": null,
                "platform": "macOS,windows10",
                "deviceMode": "none",
                "technologies": "mobileApplicationManagement"
            },
            "referredSettingInformationList": [],
            "options": [
                {
                    "itemId": "com.microsoft.edge.mamedgeappconfigsettings.forcesync_false",
                    "description": null,
                    "helpText": null,
                    "name": "Disabled",
                    "displayName": "Disabled",
                    "optionValue": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                        "settingValueTemplateReference": null,
                        "value": "false"
                    },
                    "dependentOn": [],
                    "dependedOnBy": []
                },
                {
                    "itemId": "com.microsoft.edge.mamedgeappconfigsettings.forcesync_true",
                    "description": null,
                    "helpText": null,
                    "name": "Enabled",
                    "displayName": "Enabled",
                    "optionValue": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                        "settingValueTemplateReference": null,
                        "value": "true"
                    },
                    "dependentOn": [],
                    "dependedOnBy": []
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingDefinition",
            "accessTypes": "none",
            "keywords": [
                "Uncategorized"
            ],
            "infoUrls": [],
            "occurrence": null,
            "baseUri": "",
            "offsetUri": "com.microsoft.edge.SearchSuggestEnabled_recommended",
            "rootDefinitionId": "com.microsoft.edge.mamedgeappconfigsettings.searchsuggestenabled_recommended",
            "categoryId": "ce3ba7d1-c101-4ab6-bf2a-0c683e921bf8",
            "settingUsage": "configuration",
            "uxBehavior": "toggle",
            "visibility": "settingsCatalog",
            "riskLevel": "low",
            "id": "com.microsoft.edge.mamedgeappconfigsettings.searchsuggestenabled_recommended",
            "description": "Enables web search suggestions in Microsoft Edge's Address Bar and Auto-Suggest List and prevents users from changing this policy.\n\nIf you enable this policy, web search suggestions are used.\n\nIf you disable this policy, web search suggestions are never used, however local history and local favorites suggestions still appear. If you disable this policy, neither the typed characters, nor the URLs visited will be included in telemetry to Microsoft.\n\nIf this policy is left not set, search suggestions are enabled but the user can change that.",
            "helpText": null,
            "name": "SearchSuggestEnabled_recommended",
            "displayName": "Enable search suggestions (users can override)",
            "version": "638951450276131607",
            "defaultOptionId": null,
            "applicability": {
                "@odata.type": "#microsoft.graph.deviceManagementConfigurationApplicationSettingApplicability",
                "description": null,
                "platform": "macOS,windows10",
                "deviceMode": "none",
                "technologies": "mobileApplicationManagement"
            },
            "referredSettingInformationList": [],
            "options": [
                {
                    "itemId": "com.microsoft.edge.mamedgeappconfigsettings.searchsuggestenabled_recommended_false",
                    "description": null,
                    "helpText": null,
                    "name": "Disabled",
                    "displayName": "Disabled",
                    "optionValue": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                        "settingValueTemplateReference": null,
                        "value": "false"
                    },
                    "dependentOn": [],
                    "dependedOnBy": []
                },
                {
                    "itemId": "com.microsoft.edge.mamedgeappconfigsettings.searchsuggestenabled_recommended_true",
                    "description": null,
                    "helpText": null,
                    "name": "Enabled",
                    "displayName": "Enabled",
                    "optionValue": {
                        "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                        "settingValueTemplateReference": null,
                        "value": "true"
                    },
                    "dependentOn": [],
                    "dependedOnBy": []
                }
            ]
        },

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/targetedManagedAppConfigurations
Request Method
POST

{
    "displayName":"test",
    "description":"test",
    "apps":[],
    "assignments":[
        {"target":{
            "groupId":"35d09841-af73-43e6-a59f-024fef1b6b95",
            "@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget"
        }
        },
        {"target":
        {"groupId":"410a28bd-9c9f-403f-b1b2-4a0bd04e98d9",
        "@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget"
        }},{"target":{"groupId":"8993d466-1e7c-44e5-8246-99675956a27f","@odata.type":"#microsoft.graph.groupAssignmentTarget"}},{"target":{"groupId":"09cb5968-2722-4ed4-aa86-2aeaeeab8e1f","@odata.type":"#microsoft.graph.groupAssignmentTarget"}}],
        "roleScopeTagIds":["0"],
        "customSettings":[
            {
                "name":"thing","value":"thing"
                },
                {
                    "name":"thing2","value":"thing2"},
                {"name":"com.microsoft.intune.mam.managedbrowser.AppProxyRedirection","value":"true"
                },
                {
                    "name":"com.microsoft.intune.mam.managedbrowser.homepage","value":"http://thing"},{"name":"com.microsoft.intune.mam.managedbrowser.bookmarks","value":"thing|thing||thing2|thing2"},{"name":"com.microsoft.intune.mam.managedbrowser.AllowListURLs","value":"url1|url2"},{"name":"com.microsoft.intune.mam.managedbrowser.AllowTransitionOnBlock","value":"true"},{"name":"com.microsoft.outlook.Mail.BlockExternalImagesEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.BlockExternalImagesEnabled.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.DefaultSignatureEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.ExternalRecipientsToolTipEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.FocusedInbox","value":"true"},{"name":"com.microsoft.intune.mam.areWearablesAllowed","value":"true"},{"name":"com.microsoft.outlook.Auth.Biometric","value":"true"},{"name":"com.microsoft.outlook.Auth.Biometric.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Contacts.LocalSyncEnabled","value":"true"},{"name":"com.microsoft.outlook.Contacts.LocalSyncEnabled.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.SuggestedRepliesEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.SuggestedRepliesEnabled.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Calendar.Notifications.IntuneMAMOnly","value":"1"},{"name":"com.microsoft.outlook.AddinsAvailable.IntuneMAMOnly","value":"true"},{"name":"com.microsoft.outlook.Mail.OfficeFeedEnabled","value":"true"},{"name":"com.microsoft.outlook.Calendar.NativeSyncEnabled","value":"true"},{"name":"com.microsoft.outlook.Calendar.NativeSyncEnabled.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Calendar.NativeSyncAvailable.IntuneMAMOnly","value":"true"},{"name":"com.microsoft.outlook.Mail.OrganizeByThreadEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.PlayMyEmailsEnabled","value":"true"},{"name":"com.microsoft.outlook.ContactSync.AddressAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.BirthdayAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.CompanyAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.DepartmentAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.EmailAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.InstantMessageAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.JobTitleAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.NicknameAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.NotesAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhoneHomeAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhoneHomeFaxAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhoneMobileAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhoneOtherAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhonePagerAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhoneWorkAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.PhoneWorkFaxAllowed","value":"true"},{"name":"com.microsoft.outlook.ContactSync.SuffixAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.TextPredictionsEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.TextPredictionsEnabled.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled.EncryptAllMail","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled.EncryptAllMail.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled.SignAllMail","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled.SignAllMail.UserChangeAllowed","value":"true"},{"name":"com.microsoft.outlook.Mail.SMIMEEnabled.LDAPHostName","value":"url"}],"settings":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.enablemediarouter","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.enablemediarouter_true","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.showcasticonintoolbar","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.showcasticonintoolbar_true","children":[]}}}],"appGroupType":"allApps","targetedAppManagementLevels":"unspecified"}

resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations/$entity",
    "displayName": "test",
    "description": "test",
    "createdDateTime": "2025-10-13T05:00:50.975446Z",
    "lastModifiedDateTime": "2025-10-13T05:00:50.975446Z",
    "roleScopeTagIds": [
        "0"
    ],
    "id": "A_970b3da6-cd39-4bdd-8fa5-edb2971c8269",
    "version": "\"cb02153a-0000-0d00-0000-68ec87830000\"",
    "deployedAppCount": 0,
    "isAssigned": false,
    "targetedAppManagementLevels": "unspecified",
    "appGroupType": "allApps",
    "customSettings": [
        {
            "name": "thing",
            "value": "thing"
        },
        {
            "name": "thing2",
            "value": "thing2"
        },
        {
            "name": "com.microsoft.intune.mam.managedbrowser.AppProxyRedirection",
            "value": "true"
        },
        {
            "name": "com.microsoft.intune.mam.managedbrowser.homepage",
            "value": "http://thing"
        },
        {
            "name": "com.microsoft.intune.mam.managedbrowser.bookmarks",
            "value": "thing|thing||thing2|thing2"
        },
        {
            "name": "com.microsoft.intune.mam.managedbrowser.AllowListURLs",
            "value": "url1|url2"
        },
        {
            "name": "com.microsoft.intune.mam.managedbrowser.AllowTransitionOnBlock",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.BlockExternalImagesEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.BlockExternalImagesEnabled.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.DefaultSignatureEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.ExternalRecipientsToolTipEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.FocusedInbox",
            "value": "true"
        },
        {
            "name": "com.microsoft.intune.mam.areWearablesAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Auth.Biometric",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Auth.Biometric.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Contacts.LocalSyncEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Contacts.LocalSyncEnabled.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SuggestedRepliesEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SuggestedRepliesEnabled.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Calendar.Notifications.IntuneMAMOnly",
            "value": "1"
        },
        {
            "name": "com.microsoft.outlook.AddinsAvailable.IntuneMAMOnly",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.OfficeFeedEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Calendar.NativeSyncEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Calendar.NativeSyncEnabled.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Calendar.NativeSyncAvailable.IntuneMAMOnly",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.OrganizeByThreadEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.PlayMyEmailsEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.AddressAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.BirthdayAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.CompanyAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.DepartmentAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.EmailAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.InstantMessageAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.JobTitleAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.NicknameAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.NotesAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhoneHomeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhoneHomeFaxAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhoneMobileAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhoneOtherAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhonePagerAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhoneWorkAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.PhoneWorkFaxAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.ContactSync.SuffixAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.TextPredictionsEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.TextPredictionsEnabled.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled.EncryptAllMail",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled.EncryptAllMail.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled.SignAllMail",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled.SignAllMail.UserChangeAllowed",
            "value": "true"
        },
        {
            "name": "com.microsoft.outlook.Mail.SMIMEEnabled.LDAPHostName",
            "value": "url"
        }
    ]
}

then

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/targetedManagedAppConfigurations('A_d5501407-b29c-4a53-bbd6-27b9b129882e')?$expand=apps,assignments,settings
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations(apps(),assignments(),settings())/$entity",
    "displayName": "tet",
    "description": "tset",
    "createdDateTime": "2025-09-05T13:38:54.5357837Z",
    "lastModifiedDateTime": "2025-09-05T13:38:54Z",
    "roleScopeTagIds": [
        "0"
    ],
    "id": "A_d5501407-b29c-4a53-bbd6-27b9b129882e",
    "version": "\"6f05ce19-0000-0d00-0000-68bae7ee0000\"",
    "deployedAppCount": 12,
    "isAssigned": true,
    "targetedAppManagementLevels": "unspecified",
    "appGroupType": "selectedPublicApps",
    "customSettings": [
        {
            "name": "thing",
            "value": "value"
        },
        {
            "name": "thing2",
            "value": "value"
        },
        {
            "name": "com.microsoft.tunnel.connection_type",
            "value": "MicrosoftProtect"
        },
        {
            "name": "com.microsoft.tunnel.connection_name",
            "value": "some-test-name"
        },
        {
            "name": "com.microsoft.tunnel.proxy_pacurl",
            "value": "test"
        },
        {
            "name": "com.microsoft.tunnel.proxy_address",
            "value": "10.10.10.10"
        },
        {
            "name": "com.microsoft.tunnel.proxy_port",
            "value": "9"
        }
    ],
    "apps@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations('A_d5501407-b29c-4a53-bbd6-27b9b129882e')/apps",
    "apps": [
        {
            "id": "com.microsoft.exchange.bookings.android",
            "version": "-1296747733",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
                "packageId": "com.microsoft.exchange.bookings"
            }
        },
        {
            "id": "com.microsoft.office.excel.android",
            "version": "-1789826587",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
                "packageId": "com.microsoft.office.excel"
            }
        },
        {
            "id": "com.microsoft.office.excel.ios",
            "version": "-1255026913",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
                "bundleId": "com.microsoft.office.excel"
            }
        },
        {
            "id": "com.microsoft.office.officelens.android",
            "version": "-641720584",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
                "packageId": "com.microsoft.office.officelens"
            }
        },
        {
            "id": "com.microsoft.office365booker.ios",
            "version": "-169796186",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
                "bundleId": "com.microsoft.office365booker"
            }
        },
        {
            "id": "com.microsoft.officelens.ios",
            "version": "1265495298",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
                "bundleId": "com.microsoft.officelens"
            }
        },
        {
            "id": "com.microsoft.ramobile.android",
            "version": "359972835",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
                "packageId": "com.microsoft.ramobile"
            }
        },
        {
            "id": "com.microsoft.ramobile.ios",
            "version": "2006965645",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
                "bundleId": "com.microsoft.ramobile"
            }
        },
        {
            "id": "com.microsoft.stream.android",
            "version": "648128536",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
                "packageId": "com.microsoft.stream"
            }
        },
        {
            "id": "com.microsoft.stream.ios",
            "version": "126860698",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
                "bundleId": "com.microsoft.stream"
            }
        },
        {
            "id": "com.microsoft.visio.ios",
            "version": "-387333446",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.iosMobileAppIdentifier",
                "bundleId": "com.microsoft.visio"
            }
        },
        {
            "id": "com.servicenow.onboarding.mam.intune.android",
            "version": "556600564",
            "mobileAppIdentifier": {
                "@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
                "packageId": "com.servicenow.onboarding.mam.intune"
            }
        }
    ],
    "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations('A_d5501407-b29c-4a53-bbd6-27b9b129882e')/assignments",
    "assignments": [
        {
            "id": "35d09841-af73-43e6-a59f-024fef1b6b95_incl",
            "source": "direct",
            "sourceId": "00000000-0000-0000-0000-000000000000",
            "target": {
                "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none",
                "groupId": "35d09841-af73-43e6-a59f-024fef1b6b95"
            }
        },
        {
            "id": "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2_incl",
            "source": "direct",
            "sourceId": "00000000-0000-0000-0000-000000000000",
            "target": {
                "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none",
                "groupId": "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
            }
        },
        {
            "id": "0fdc63eb-c85d-4b82-8c2f-5abaaf18ff30_excl",
            "source": "direct",
            "sourceId": "00000000-0000-0000-0000-000000000000",
            "target": {
                "@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none",
                "groupId": "0fdc63eb-c85d-4b82-8c2f-5abaaf18ff30"
            }
        },
        {
            "id": "8993d466-1e7c-44e5-8246-99675956a27f_excl",
            "source": "direct",
            "sourceId": "00000000-0000-0000-0000-000000000000",
            "target": {
                "@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
                "deviceAndAppManagementAssignmentFilterId": null,
                "deviceAndAppManagementAssignmentFilterType": "none",
                "groupId": "8993d466-1e7c-44e5-8246-99675956a27f"
            }
        }
    ],
    "settings@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations('A_d5501407-b29c-4a53-bbd6-27b9b129882e')/settings",
    "settings": []
}

create assignments

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/targetedManagedAppConfigurations('A_d5501407-b29c-4a53-bbd6-27b9b129882e')/assign
Request Method
POST

{"assignments":[{"target":{"groupId":"b15228f4-9d49-41ed-9b4f-0e7c721fd9c2","@odata.type":"#microsoft.graph.groupAssignmentTarget"}},{"target":{"groupId":"35d09841-af73-43e6-a59f-024fef1b6b95","@odata.type":"#microsoft.graph.groupAssignmentTarget"}},{"target":{"groupId":"0fdc63eb-c85d-4b82-8c2f-5abaaf18ff30","@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget"}},{"target":{"groupId":"8993d466-1e7c-44e5-8246-99675956a27f","@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget"}}]}



Request URL
https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations
Request Method
POST

{"@odata.type":"#microsoft.graph.iosMobileAppConfiguration","displayName":"managed_device_test","description":"thing","roleScopeTagIds":["0"],"targetedMobileApps":["a5204f11-fe5b-4ac5-b379-945d79889188"],"encodedSettingXml":"","settings":[{"appConfigKey":"thing","appConfigKeyType":"IntegerType","appConfigKeyValue":"1"},{"appConfigKey":"thing2","appConfigKeyType":"RealType","appConfigKeyValue":"1.1"},{"appConfigKey":"thing3","appConfigKeyType":"StringType","appConfigKeyValue":"some_string"},{"appConfigKey":"thing4","appConfigKeyType":"BooleanType","appConfigKeyValue":"true"}],"id":"00000000-0000-0000-0000-000000000000"}

settings - update

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/targetedManagedAppConfigurations('A_af5beee3-f460-428a-86e5-71e64e492dda')/changeSettings
Request Method
POST

{"settings":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.enablemediarouter","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.enablemediarouter_true","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.showcasticonintoolbar","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.showcasticonintoolbar_true","children":[]}}}]}


or like this 

{
    "displayName":"test",
    "description":"test",
    "apps":[],
    "assignments":[],
    "roleScopeTagIds":["0"],
    "customSettings":[
        {
            "name":"com.microsoft.intune.mam.managedbrowser.AppProxyRedirection",
            "value":"false"
        },
        {
            "name":"com.microsoft.intune.mam.managedbrowser.AllowTransitionOnBlock",
            "value":"true"
        }
    ],
    "settings":[
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":{
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.cookiesallowedforurls",
                "simpleSettingCollectionValue":[
                    {
                        "value":"thing",
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    },
                    {
                        "value":"thing2",
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    }
                ]
            }
        },
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":{
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.imagesallowedforurls",
                "simpleSettingCollectionValue":[
                    {
                        "value":"thing",
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    },
                    {
                        "value":"thing2",
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    }
                ]
            }
        },
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":{
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.insecurecontentallowedforurls",
                "simpleSettingCollectionValue":[
                    {
                        "value":"thing",
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    },
                    {
                        "value":"thing2",
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    }
                ]
            }
        },
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":{
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.intranetfilelinksenabled",
                "choiceSettingValue":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.intranetfilelinksenabled_true",
                    "children":[]
                    }
                }
            },
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.javascriptallowedforurls",
                    "simpleSettingCollectionValue":[
                        {
                            "value":"thing",
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                        },
                        {
                            "value":"thing2",
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                        }
                    ]
                }
            },
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.javascriptjitallowedforsites",
                    "simpleSettingCollectionValue":[
                        {
                            "value":"thing",
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                        },
                        {
                            "value":"thing2",
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                        }
                    ]
                }
            },
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.automaticdownloadsallowedforurls",
                    "simpleSettingCollectionValue":[
                        {
                            "value":"thing",
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                        }
                    ]
                }
            },
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.notificationsallowedforurls","simpleSettingCollectionValue":[
                        {
                            "value":"thing",
                            "@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                        }
                    ]
                }
            },
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                    "settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.showpdfdefaultrecommendationsenabled",
                    "choiceSettingValue":{
                        "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
                        "value":"com.microsoft.edge.mamedgeappconfigsettings.showpdfdefaultrecommendationsenabled_false",
                        "children":[]
                    }
                }
            },
            {
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.popupsallowedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.filesystemreadaskforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.pluginsallowedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.webhidaskforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.webusbaskforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.filesystemwriteaskforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.javascriptjitblockedforsites","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.javascriptblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.automaticdownloadsblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.notificationsblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.popupsblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.filesystemreadblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.pluginsblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.webusbblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.filesystemwriteblockedforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultcookiessetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultcookiessetting_allowcookies","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultinsecurecontentsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultinsecurecontentsetting_allowexceptionsinsecurecontent","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultjavascriptjitsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultjavascriptjitsetting_blockjavascriptjit","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultfilesystemreadguardsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultfilesystemreadguardsetting_askfilesystemread","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultfilesystemwriteguardsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultfilesystemwriteguardsetting_askfilesystemwrite","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultwebbluetoothguardsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultwebbluetoothguardsetting_askwebbluetooth","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultwebhidguardsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultwebhidguardsetting_askwebhid","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultwebusbguardsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultwebusbguardsetting_askwebusb","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultpluginssetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultpluginssetting_clicktoplay","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultautomaticdownloadssetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultautomaticdownloadssetting_blockautomaticdownloads","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultgeolocationsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultgeolocationsetting_allowgeolocation","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultimagessetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultimagessetting_blockimages","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultjavascriptsetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultjavascriptsetting_blockjavascript","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultnotificationssetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultnotificationssetting_allownotifications","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.defaultpopupssetting","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.defaultpopupssetting_blockpopups","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.legacysamesitecookiebehaviorenabled","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.legacysamesitecookiebehaviorenabled_defaulttosamesitebydefaultcookiebehavior","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.webusballowdevicesforurls","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","value":"1"}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.cookiessessiononlyforurls","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.registeredprotocolhandlers_recommended","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","value":"\"%s\""}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.legacysamesitecookiebehaviorenabledfordomainlist","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.extensioninstallallowlist","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"},{"value":"thing2","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.blockexternalextensions","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.blockexternalextensions_true","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes","choiceSettingCollectionValue":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes_extension","children":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes_theme","children":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes_user_script","children":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes_legacy_packaged_app","children":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes_hosted_app","children":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.extensionallowedtypes_platform_app","children":[]}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.controldefaultstateofallowextensionfromotherstoressettingenabled_recommended","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.controldefaultstateofallowextensionfromotherstoressettingenabled_recommended_true","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.extensioninstallsources","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"},{"value":"thing2","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.extensionsettings","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","value":"thing"}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.extensioninstallforcelist","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"},{"value":"thing2","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.extensioninstallblocklist","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"},{"value":"thing2","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.linkedaccountenabled","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.linkedaccountenabled_true","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.nativemessaginguserlevelhosts","choiceSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.nativemessaginguserlevelhosts_true","children":[]}}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.nativemessagingblocklist","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"},{"value":"thing2","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.nativemessagingallowlist","simpleSettingCollectionValue":[{"value":"thing","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"},{"value":"thing2","@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue"}]}}],"appGroupType":"allApps","targetedAppManagementLevels":"unspecified"}



{
    "displayName":"test",
    "description":"test",
    "apps":[],
    "assignments":[],
    "roleScopeTagIds":["0"],
    "customSettings":[
        {
            "name":"com.microsoft.intune.mam.managedbrowser.AppProxyRedirection",
            "value":"false"
        },
        {
            "name":"com.microsoft.intune.mam.managedbrowser.AllowTransitionOnBlock",
            "value":"true"
        }
    ],
    "settings":[
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":{
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.enablemediarouter","choiceSettingValue":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.enablemediarouter_true","children":[]
                }
            }
        },
        {
            "@odata.type":"#microsoft.graph.deviceManagementConfigurationSetting",
            "settingInstance":{
                "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"com.microsoft.edge.mamedgeappconfigsettings.showcasticonintoolbar",
                "choiceSettingValue":{
                    "@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingValue","value":"com.microsoft.edge.mamedgeappconfigsettings.showcasticonintoolbar_true","children":[]
                }
            }
        }
    ],
    "appGroupType":"allApps",
    "targetedAppManagementLevels":"unspecified"
}