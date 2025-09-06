Request URL
https://graph.microsoft.com/beta/deviceAppManagement/targetedManagedAppConfigurations
Request Method
POST

{"displayName":"tet","description":"tset","apps":[{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.androidMobileAppIdentifier","packageId":"com.microsoft.office.excel"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.iosMobileAppIdentifier","bundleId":"com.microsoft.office.excel"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.androidMobileAppIdentifier","packageId":"com.microsoft.office.officelens"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.iosMobileAppIdentifier","bundleId":"com.microsoft.officelens"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.androidMobileAppIdentifier","packageId":"com.servicenow.onboarding.mam.intune"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.iosMobileAppIdentifier","bundleId":"com.microsoft.visio"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.iosMobileAppIdentifier","bundleId":"com.microsoft.stream"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.androidMobileAppIdentifier","packageId":"com.microsoft.stream"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.iosMobileAppIdentifier","bundleId":"com.microsoft.ramobile"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.androidMobileAppIdentifier","packageId":"com.microsoft.ramobile"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.iosMobileAppIdentifier","bundleId":"com.microsoft.office365booker"}},{"mobileAppIdentifier":{"@odata.type":"#microsoft.graph.androidMobileAppIdentifier","packageId":"com.microsoft.exchange.bookings"}}],"assignments":[],"roleScopeTagIds":["0"],"customSettings":[{"name":"thing","value":"value"},{"name":"thing2","value":"value"},{"name":"com.microsoft.tunnel.connection_type","value":"MicrosoftProtect"},{"name":"com.microsoft.tunnel.connection_name","value":"some-test-name"},{"name":"com.microsoft.tunnel.proxy_pacurl","value":"test"},{"name":"com.microsoft.tunnel.proxy_address","value":"10.10.10.10"},{"name":"com.microsoft.tunnel.proxy_port","value":"9"}],"settings":[],"appGroupType":"selectedPublicApps","targetedAppManagementLevels":"unspecified"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations/$entity",
    "displayName": "tet",
    "description": "tset",
    "createdDateTime": "2025-09-05T13:38:54.4989757Z",
    "lastModifiedDateTime": "2025-09-05T13:38:54.4989757Z",
    "roleScopeTagIds": [
        "0"
    ],
    "id": "A_d5501407-b29c-4a53-bbd6-27b9b129882e",
    "version": "\"6f05ce19-0000-0d00-0000-68bae7ee0000\"",
    "deployedAppCount": 12,
    "isAssigned": false,
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
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations?&$filter=microsoft.graph.androidManagedStoreAppConfiguration/appSupportsOemConfig%20eq%20false%20or%20isof(%27microsoft.graph.androidManagedStoreAppConfiguration%27)%20eq%20false&$top=500&$count=true
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileAppConfigurations",
    "@odata.count": 7,
    "value": [
        {
            "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
            "id": "5e18ec15-2587-42f8-a974-e2bdb2cef334",
            "targetedMobileApps": [
                "33029352-f792-4507-b963-ab2441a0c5f0"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-04-21T08:59:22.4273305Z",
            "description": "Limit access to only allowed organization user accounts and block personal accounts on enrolled devices for Teams\n\nRef: https://docs.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android#allow-only-configured-organization-accounts-in-apps\n\n21.04.2022",
            "lastModifiedDateTime": "2022-04-21T08:59:22.4273305Z",
            "displayName": "[MDM] ACP | Microsoft Teams [Android Enterprise - Fully Managed]",
            "version": 1,
            "packageId": "app:com.microsoft.teams",
            "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC50ZWFtcyIsIm1hbmFnZWRQcm9wZXJ0eSI6W3sia2V5IjoiY29tLm1pY3Jvc29mdC5pbnR1bmUubWFtLkFsbG93ZWRBY2NvdW50VVBOcyIsInZhbHVlU3RyaW5nIjoie3t1c2VycHJpbmNpcGFsbmFtZX19In1dfQ==",
            "appSupportsOemConfig": false,
            "profileApplicability": "androidDeviceOwner",
            "connectedAppsEnabled": false,
            "permissionActions": []
        },
        {
            "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
            "id": "787c21fa-68f1-4c41-aadb-c56cf0da50f2",
            "targetedMobileApps": [
                "f400d6e7-08de-4267-8db2-035253751022"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-04-21T08:51:19.3301921Z",
            "description": "Limit access to only allowed organization user accounts and block personal accounts on enrolled devices for OneNote\n\nRef: https://docs.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android#allow-only-configured-organization-accounts-in-apps",
            "lastModifiedDateTime": "2022-04-21T08:51:19.3301921Z",
            "displayName": "[MDM] ACP | Microsoft OneNote [Android Enterprise - Fully Managed]",
            "version": 1,
            "packageId": "app:com.microsoft.office.onenote",
            "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5vZmZpY2Uub25lbm90ZSIsIm1hbmFnZWRQcm9wZXJ0eSI6W3sia2V5IjoiY29tLm1pY3Jvc29mdC5pbnR1bmUubWFtLkFsbG93ZWRBY2NvdW50VVBOcyIsInZhbHVlU3RyaW5nIjoie3t1c2VycHJpbmNpcGFsbmFtZX19In0seyJrZXkiOiJjb20ubWljcm9zb2Z0Lm9mZmljZS5Ob3Rlc0NyZWF0aW9uRW5hYmxlZCIsInZhbHVlQm9vbCI6dHJ1ZX1dfQ==",
            "appSupportsOemConfig": false,
            "profileApplicability": "androidDeviceOwner",
            "connectedAppsEnabled": false,
            "permissionActions": []
        },
        {
            "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
            "id": "863a383f-3144-4a76-9338-621fa9b3592c",
            "targetedMobileApps": [
                "57eb943a-1542-4c9b-b220-ed09b95a50a2"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-04-21T07:20:10.3847099Z",
            "description": "Limit access to only allowed organization user accounts and block personal accounts on enrolled devices for Edge\n\nRef: https://docs.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android#allow-only-configured-organization-accounts-in-apps",
            "lastModifiedDateTime": "2022-04-21T07:20:10.3847099Z",
            "displayName": "[MDM] ACP | Microsoft Edge [Android Enterprise - Fully Managed]",
            "version": 1,
            "packageId": "app:com.microsoft.emmx",
            "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5lbW14IiwibWFuYWdlZFByb3BlcnR5IjpbeyJrZXkiOiJjb20ubWljcm9zb2Z0LmludHVuZS5tYW0uQWxsb3dlZEFjY291bnRVUE5zIiwidmFsdWVTdHJpbmciOiJ7e3VzZXJwcmluY2lwbGVuYW1lfX0ifV19",
            "appSupportsOemConfig": false,
            "profileApplicability": "androidDeviceOwner",
            "connectedAppsEnabled": false,
            "permissionActions": []
        },
        {
            "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
            "id": "89550771-b65b-4c11-90f5-c08b870db1fa",
            "targetedMobileApps": [
                "9711516a-f6f8-4953-ad1f-45920ef34dda"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-04-21T07:29:53.5133968Z",
            "description": "Limit access to only allowed organization user accounts and block personal accounts on enrolled devices\n\nRef: https://docs.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android#allow-only-configured-organization-accounts-in-apps",
            "lastModifiedDateTime": "2022-04-21T07:29:53.5133968Z",
            "displayName": "[MDM] ACP | Microsoft Office [Android Enterprise - Fully Managed]",
            "version": 1,
            "packageId": "app:com.microsoft.office.officehubrow",
            "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5vZmZpY2Uub2ZmaWNlaHVicm93IiwibWFuYWdlZFByb3BlcnR5IjpbeyJrZXkiOiJjb20ubWljcm9zb2Z0Lm9mZmljZS5Ob3Rlc0NyZWF0aW9uRW5hYmxlZCIsInZhbHVlQm9vbCI6dHJ1ZX0seyJrZXkiOiJjb20ubWljcm9zb2Z0LmludHVuZS5tYW0uQWxsb3dlZEFjY291bnRVUE5zIiwidmFsdWVTdHJpbmciOiJ7e3VzZXJwcmluY2lwYWxuYW1lfX0ifV19",
            "appSupportsOemConfig": false,
            "profileApplicability": "androidDeviceOwner",
            "connectedAppsEnabled": false,
            "permissionActions": []
        },
        {
            "@odata.type": "#microsoft.graph.iosMobileAppConfiguration",
            "id": "91425042-1035-4e41-bcc5-be8a1f8d2430",
            "targetedMobileApps": [
                "5c4b0c8a-4016-449b-a4ed-cdfbf39ab0ff"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-02-28T12:57:33.6260799Z",
            "description": "21.04.2022",
            "lastModifiedDateTime": "2022-04-21T07:07:43.6890195Z",
            "displayName": "[MDM] ACP | Microsoft Edge [iOS]",
            "version": 5,
            "encodedSettingXml": null,
            "settings": [
                {
                    "appConfigKey": "Intune Allowed Accounts",
                    "appConfigKeyType": "stringType",
                    "appConfigKeyValue": "{{userprinciplename}}"
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
            "id": "9c7a451d-f1c3-45ab-9cda-391d302bec90",
            "targetedMobileApps": [
                "df3baafe-df9e-43c7-9bda-8c59f0e9c2ed"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-04-21T08:55:29.3905261Z",
            "description": "Limit access to only allowed organization user accounts and block personal accounts on enrolled devices for Outlook\n\nRef: https://docs.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android#allow-only-configured-organization-accounts-in-apps",
            "lastModifiedDateTime": "2022-04-21T08:55:51.800111Z",
            "displayName": "[MDM] ACP | Microsoft Outlook [Android Enterprise - Fully Managed]",
            "version": 2,
            "packageId": "app:com.microsoft.office.outlook",
            "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5vZmZpY2Uub3V0bG9vayIsIm1hbmFnZWRQcm9wZXJ0eSI6W3sia2V5IjoiY29tLm1pY3Jvc29mdC5pbnR1bmUubWFtLkFsbG93ZWRBY2NvdW50VVBOcyIsInZhbHVlU3RyaW5nIjoie3t1c2VycHJpbmNpcGFsbmFtZX19In0seyJrZXkiOiJjb20ubWljcm9zb2Z0Lm91dGxvb2suRW1haWxQcm9maWxlLkFjY291bnRUeXBlIiwidmFsdWVTdHJpbmciOiJNb2Rlcm5BdXRoIn0seyJrZXkiOiJjb20ubWljcm9zb2Z0Lm91dGxvb2suRW1haWxQcm9maWxlLkVtYWlsVVBOIiwidmFsdWVTdHJpbmciOiJ7e3VzZXJwcmluY2lwYWxuYW1lfX0ifSx7ImtleSI6ImNvbS5taWNyb3NvZnQub3V0bG9vay5FbWFpbFByb2ZpbGUuRW1haWxBZGRyZXNzIiwidmFsdWVTdHJpbmciOiJ7e21haWx9fSJ9LHsia2V5IjoiSW50dW5lTUFNQWxsb3dlZEFjY291bnRzT25seSIsInZhbHVlU3RyaW5nIjoiRW5hYmxlZCJ9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vdXRsb29rLk1haWwuRm9jdXNlZEluYm94IiwidmFsdWVCb29sIjpmYWxzZX0seyJrZXkiOiJjb20ubWljcm9zb2Z0Lm91dGxvb2suQ29udGFjdHMuTG9jYWxTeW5jRW5hYmxlZCIsInZhbHVlQm9vbCI6dHJ1ZX0seyJrZXkiOiJjb20ubWljcm9zb2Z0Lm91dGxvb2suTWFpbC5vZmZpY2VGZWVkRW5hYmxlZCIsInZhbHVlQm9vbCI6ZmFsc2V9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vdXRsb29rLk1haWwuU3VnZ2VzdGVkUmVwbGllc0VuYWJsZWQiLCJ2YWx1ZUJvb2wiOnRydWV9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vdXRsb29rLk1haWwuRXh0ZXJuYWxSZWNpcGllbnRzVG9vbFRpcEVuYWJsZWQiLCJ2YWx1ZUJvb2wiOnRydWV9LHsia2V5IjoiY29tLm1pY3Jvc29mdC5vdXRsb29rLk1haWwuRGVmYXVsdFNpZ25hdHVyZUVuYWJsZWQiLCJ2YWx1ZUJvb2wiOmZhbHNlfSx7ImtleSI6ImNvbS5taWNyb3NvZnQub3V0bG9vay5NYWlsLkJsb2NrRXh0ZXJuYWxJbWFnZXNFbmFibGVkIiwidmFsdWVCb29sIjpmYWxzZX0seyJrZXkiOiJjb20ubWljcm9zb2Z0Lm91dGxvb2suQ2FsZW5kYXIuTmF0aXZlU3luY0VuYWJsZWQiLCJ2YWx1ZUJvb2wiOnRydWV9XX0=",
            "appSupportsOemConfig": false,
            "profileApplicability": "androidDeviceOwner",
            "connectedAppsEnabled": false,
            "permissionActions": []
        },
        {
            "@odata.type": "#microsoft.graph.androidManagedStoreAppConfiguration",
            "id": "bb828edb-6e98-404b-9ddf-0aa753f08bf4",
            "targetedMobileApps": [
                "970c9b4a-4879-4b6b-985e-693167bff8f6"
            ],
            "roleScopeTagIds": [
                "0"
            ],
            "createdDateTime": "2022-04-21T07:32:27.674601Z",
            "description": "Limit access to only allowed organization user accounts and block personal accounts on enrolled devices for OneDrive\n\nRef: https://docs.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android#allow-only-configured-organization-accounts-in-apps\n\n21.04.2022",
            "lastModifiedDateTime": "2022-04-21T08:56:30.6256456Z",
            "displayName": "[MDM] ACP | Microsoft OneDrive [Android Enterprise - Fully Managed]",
            "version": 2,
            "packageId": "app:com.microsoft.skydrive",
            "payloadJson": "eyJraW5kIjoiYW5kcm9pZGVudGVycHJpc2UjbWFuYWdlZENvbmZpZ3VyYXRpb24iLCJwcm9kdWN0SWQiOiJhcHA6Y29tLm1pY3Jvc29mdC5za3lkcml2ZSIsIm1hbmFnZWRQcm9wZXJ0eSI6W3sia2V5IjoiY29tLm1pY3Jvc29mdC5pbnR1bmUubWFtLkFsbG93ZWRBY2NvdW50VVBOcyIsInZhbHVlU3RyaW5nIjoie3t1c2VycHJpbmNpcGFsbmFtZX19In1dfQ==",
            "appSupportsOemConfig": false,
            "profileApplicability": "androidDeviceOwner",
            "connectedAppsEnabled": false,
            "permissionActions": []
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/targetedManagedAppConfigurations?$count=true
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/targetedManagedAppConfigurations",
    "@odata.count": 4,
    "value": [
        {
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
            "isAssigned": false,
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
            ]
        },
        {
            "displayName": "test",
            "description": "test",
            "createdDateTime": "2025-07-10T14:53:14.4387858Z",
            "lastModifiedDateTime": "2025-07-10T14:53:15Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "A_936da198-54cb-4c9c-886a-0947e71d63b9",
            "version": "\"2d04f7c3-0000-0d00-0000-686fd3db0000\"",
            "deployedAppCount": 310,
            "isAssigned": false,
            "targetedAppManagementLevels": "unspecified",
            "appGroupType": "selectedPublicApps",
            "customSettings": [
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.AppProxyRedirection",
                    "value": "false"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.AllowTransitionOnBlock",
                    "value": "true"
                }
            ]
        },
        {
            "displayName": "[MAM] ACP | Microsoft Edge [iOS + Android]",
            "description": "- Enforces Microsoft Edge as the default browser for MAM managed applications on iOS and Android\n- Enable switching between personal and inprivate context within Edge \n- Account.syncDisabled\n21.04.2022",
            "createdDateTime": "2022-04-21T06:55:10.8178973Z",
            "lastModifiedDateTime": "2022-04-21T06:55:52Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "A_909d2947-9f9a-4f4f-8be0-41f3079e86b6",
            "version": "\"ab01f92b-0000-0d00-0000-6260fff80000\"",
            "deployedAppCount": 2,
            "isAssigned": true,
            "targetedAppManagementLevels": "unspecified",
            "appGroupType": "selectedPublicApps",
            "customSettings": [
                {
                    "name": "com.microsoft.intune.useEdge",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.openinPrivateIfBlocked",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.account.syncDisabled",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.AppProxyRedirection",
                    "value": "false"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.homepage",
                    "value": "https://www.bbc.co.uk/"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.bookmarks",
                    "value": "BBC|https://www.bbc.co.uk/"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.BlockListURLs",
                    "value": "www.pornhub.com"
                },
                {
                    "name": "com.microsoft.intune.mam.managedbrowser.AllowTransitionOnBlock",
                    "value": "true"
                }
            ]
        },
        {
            "displayName": "[MAM] ACP | Microsoft Outlook [iOS + Android]",
            "description": "- Controls which contacts fields can sync to native contacts app\n- Blocks two way calender sync\n21.04.2022",
            "createdDateTime": "2022-04-21T07:04:32.683661Z",
            "lastModifiedDateTime": "2022-04-21T07:04:32Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "A_cba67ad3-5c0b-4030-97a5-9d55d15ebb7f",
            "version": "\"ac01925d-0000-0d00-0000-626102000000\"",
            "deployedAppCount": 2,
            "isAssigned": true,
            "targetedAppManagementLevels": "unspecified",
            "appGroupType": "selectedPublicApps",
            "customSettings": [
                {
                    "name": "com.microsoft.outlook.Mail.BlockExternalImagesEnabled",
                    "value": "false"
                },
                {
                    "name": "com.microsoft.outlook.Mail.BlockExternalImagesEnabled.UserChangeAllowed",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.outlook.Mail.DefaultSignatureEnabled",
                    "value": "false"
                },
                {
                    "name": "com.microsoft.outlook.Mail.FocusedInbox",
                    "value": "false"
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
                    "name": "com.microsoft.outlook.Mail.SuggestedRepliesEnabled",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.outlook.Mail.SuggestedRepliesEnabled.UserChangeAllowed",
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
                    "name": "com.microsoft.outlook.ContactSync.PhoneHomeAllowed",
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
                    "name": "com.microsoft.outlook.ContactSync.PhoneWorkAllowed",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.outlook.Mail.TextPredictionsEnabled",
                    "value": "true"
                },
                {
                    "name": "com.microsoft.outlook.Mail.TextPredictionsEnabled.UserChangeAllowed",
                    "value": "true"
                }
            ]
        }
    ]
}