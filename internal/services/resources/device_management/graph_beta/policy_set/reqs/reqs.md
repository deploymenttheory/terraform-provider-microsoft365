Request URL
https://graph.microsoft.com/beta/deviceAppManagement/mobileApps?$filter=(isof(%27microsoft.graph.iosStoreApp%27)%20or%20isof(%27microsoft.graph.iosLobApp%27)%20or%20(isof(%27microsoft.graph.managedIOSStoreApp%27)%20and%20microsoft.graph.managedApp/appAvailability%20eq%20microsoft.graph.managedAppAvailability%27global%27)%20or%20isof(%27microsoft.graph.managedIOSLobApp%27)%20or%20isof(%27microsoft.graph.androidStoreApp%27)%20or%20isof(%27microsoft.graph.androidLobApp%27)%20or%20(isof(%27microsoft.graph.managedAndroidStoreApp%27)%20and%20microsoft.graph.managedApp/appAvailability%20eq%20microsoft.graph.managedAppAvailability%27global%27)%20or%20isof(%27microsoft.graph.managedAndroidLobApp%27)%20or%20isof(%27microsoft.graph.officeSuiteApp%27)%20or%20isof(%27microsoft.graph.webApp%27)%20or%20isof(%27microsoft.graph.windowsMobileMSI%27)%20or%20isof(%27microsoft.graph.windowsMicrosoftEdgeApp%27)%20or%20isof(%27microsoft.graph.macOSOfficeSuiteApp%27)%20or%20isof(%27microsoft.graph.macOSLobApp%27)%20or%20isof(%27microsoft.graph.macOSMicrosoftEdgeApp%27)%20or%20isof(%27microsoft.graph.macOSMicrosoftDefenderApp%27)%20or%20(isof(%27microsoft.graph.managedAndroidStoreApp%27)%20and%20microsoft.graph.managedApp/appAvailability%20eq%20microsoft.graph.managedAppAvailability%27lineOfBusiness%27)%20or%20(isof(%27microsoft.graph.managedIOSStoreApp%27)%20and%20microsoft.graph.managedApp/appAvailability%20eq%20microsoft.graph.managedAppAvailability%27lineOfBusiness%27))%20and%20(microsoft.graph.managedApp/appAvailability%20eq%20null%20or%20microsoft.graph.managedApp/appAvailability%20eq%20%27lineOfBusiness%27%20or%20isAssigned%20eq%20true)&$orderby=displayName&
Request Method
GET

resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps",
    "@odata.count": 0,
    "value": [
        {
            "@odata.type": "#microsoft.graph.windowsMobileMSI",
            "id": "3658ae29-7dea-4d37-8d61-1c8424742696",
            "displayName": "Configuration Manager Support Center OneTrace",
            "description": "Configuration Manager Support Center OneTrace [2002]",
            "publisher": "Microsoft",
            "largeIcon": null,
            "createdDateTime": "2020-10-29T16:55:48Z",
            "lastModifiedDateTime": "2021-05-13T07:02:22Z",
            "isFeatured": true,
            "privacyInformationUrl": null,
            "informationUrl": "https://docs.microsoft.com/en-us/mem/configmgr/core/support/support-center",
            "owner": "",
            "developer": "",
            "notes": "",
            "uploadState": 1,
            "publishingState": "published",
            "isAssigned": true,
            "roleScopeTagIds": [],
            "dependentAppCount": 0,
            "supersedingAppCount": 0,
            "supersededAppCount": 0,
            "committedContentVersion": "1",
            "fileName": "SupportCenterInstaller.msi",
            "size": 5230656,
            "commandLine": "",
            "productCode": "{79C523EC-042E-420F-83FB-4CDB6B85CE15}",
            "productVersion": "5.2002.1084.1000",
            "ignoreVersionDetection": false,
            "identityVersion": "5.2002.1084.1000",
            "useDeviceContext": false
        },
        {
            "@odata.type": "#microsoft.graph.iosStoreApp",
            "id": "a5204f11-fe5b-4ac5-b379-945d79889188",
            "displayName": "Microsoft Intune Company Portal",
            "description": "Microsoft Intune Company Portal used to create user relationship after enrollment.",
            "publisher": "Microsoft Corporation",
            "largeIcon": null,
            "createdDateTime": "2025-10-06T13:28:24Z",
            "lastModifiedDateTime": "0001-01-01T00:00:00Z",
            "isFeatured": false,
            "privacyInformationUrl": null,
            "informationUrl": null,
            "owner": null,
            "developer": null,
            "notes": null,
            "uploadState": 1,
            "publishingState": "published",
            "isAssigned": false,
            "roleScopeTagIds": [],
            "dependentAppCount": 0,
            "supersedingAppCount": 0,
            "supersededAppCount": 0,
            "bundleId": "com.microsoft.CompanyPortal",
            "appStoreUrl": "https://itunes.apple.com/us/app/microsoft-intune-company-portal/id719171358?mt=8",
            "applicableDeviceType": {
                "iPad": true,
                "iPhoneAndIPod": true
            },
            "minimumSupportedOperatingSystem": {
                "v8_0": false,
                "v9_0": true,
                "v10_0": false,
                "v11_0": false,
                "v12_0": false,
                "v13_0": false,
                "v14_0": false,
                "v15_0": false,
                "v16_0": false,
                "v17_0": false,
                "v18_0": false,
                "v26_0": false
            }
        },
        {
            "@odata.type": "#microsoft.graph.officeSuiteApp",
            "id": "b200f17f-1f5b-4806-8985-9d2a910a47f9",
            "displayName": "Visio Microsoft 365 Apps for Windows 10",
            "description": "Visio Microsoft 365 Apps for Windows 10",
            "publisher": "Microsoft",
            "largeIcon": null,
            "createdDateTime": "2021-05-27T15:39:22Z",
            "lastModifiedDateTime": "2021-10-22T09:15:44Z",
            "isFeatured": true,
            "privacyInformationUrl": "https://privacy.microsoft.com/en-US/privacystatement",
            "informationUrl": "https://products.office.com/en-us/explore-office-for-home",
            "owner": "Microsoft",
            "developer": "Microsoft",
            "notes": "",
            "uploadState": 1,
            "publishingState": "published",
            "isAssigned": true,
            "roleScopeTagIds": [],
            "dependentAppCount": 0,
            "supersedingAppCount": 0,
            "supersededAppCount": 0,
            "autoAcceptEula": false,
            "productIds": [],
            "excludedApps": null,
            "useSharedComputerActivation": false,
            "updateChannel": "none",
            "officeSuiteAppDefaultFileFormat": "notConfigured",
            "officePlatformArchitecture": "none",
            "localesToInstall": [],
            "installProgressDisplayLevel": "none",
            "shouldUninstallOlderVersionsOfOffice": false,
            "targetVersion": null,
            "updateVersion": null,
            "officeConfigurationXml": "PENvbmZpZ3VyYXRpb24gSUQ9IjY4MTJhNGFjLTY0YWEtNGY2ZS05MTg0LTBmYmM5OWM0Yjc5NiI+CiAgPEluZm8gRGVzY3JpcHRpb249IlZpc2lvIE0zNjUgQXBwcyBmb3IgRW50ZXJwcmlzZSAyNy4wNS4yMDIxIC0gU3RhbmRhcmQgQnVpbGQiIC8+CiAgPEFkZCBPZmZpY2VDbGllbnRFZGl0aW9uPSI2NCIgQ2hhbm5lbD0iTW9udGhseUVudGVycHJpc2UiIE1pZ3JhdGVBcmNoPSJUUlVFIj4KICAgIDxQcm9kdWN0IElEPSJPMzY1UHJvUGx1c1JldGFpbCI+CiAgICAgIDxMYW5ndWFnZSBJRD0iTWF0Y2hPUyIgLz4KICAgICAgPExhbmd1YWdlIElEPSJlbi11cyIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9IkFjY2VzcyIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9IkV4Y2VsIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iR3Jvb3ZlIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iTHluYyIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9Ik9uZURyaXZlIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iT25lTm90ZSIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9Ik91dGxvb2siIC8+CiAgICAgIDxFeGNsdWRlQXBwIElEPSJQb3dlclBvaW50IiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iUHVibGlzaGVyIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iVGVhbXMiIC8+CiAgICAgIDxFeGNsdWRlQXBwIElEPSJXb3JkIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iQmluZyIgLz4KICAgIDwvUHJvZHVjdD4KICAgIDxQcm9kdWN0IElEPSJWaXNpb1Byb1JldGFpbCI+CiAgICAgIDxMYW5ndWFnZSBJRD0iTWF0Y2hPUyIgLz4KICAgICAgPExhbmd1YWdlIElEPSJlbi11cyIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9IkFjY2VzcyIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9IkV4Y2VsIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iR3Jvb3ZlIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iTHluYyIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9Ik9uZURyaXZlIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iT25lTm90ZSIgLz4KICAgICAgPEV4Y2x1ZGVBcHAgSUQ9Ik91dGxvb2siIC8+CiAgICAgIDxFeGNsdWRlQXBwIElEPSJQb3dlclBvaW50IiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iUHVibGlzaGVyIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iVGVhbXMiIC8+CiAgICAgIDxFeGNsdWRlQXBwIElEPSJXb3JkIiAvPgogICAgICA8RXhjbHVkZUFwcCBJRD0iQmluZyIgLz4KICAgIDwvUHJvZHVjdD4KICA8L0FkZD4KICA8UHJvcGVydHkgTmFtZT0iU2hhcmVkQ29tcHV0ZXJMaWNlbnNpbmciIFZhbHVlPSIwIiAvPgogIDxQcm9wZXJ0eSBOYW1lPSJTQ0xDYWNoZU92ZXJyaWRlIiBWYWx1ZT0iMCIgLz4KICA8UHJvcGVydHkgTmFtZT0iQVVUT0FDVElWQVRFIiBWYWx1ZT0iMCIgLz4KICA8UHJvcGVydHkgTmFtZT0iRk9SQ0VBUFBTSFVURE9XTiIgVmFsdWU9IlRSVUUiIC8+CiAgPFByb3BlcnR5IE5hbWU9IkRldmljZUJhc2VkTGljZW5zaW5nIiBWYWx1ZT0iMCIgLz4KICA8VXBkYXRlcyBFbmFibGVkPSJUUlVFIiAvPgogIDxSZW1vdmVNU0kgLz4KICA8QXBwU2V0dGluZ3M+CiAgICA8U2V0dXAgTmFtZT0iQ29tcGFueSIgVmFsdWU9IkRlcGxveW1lbnQgVGhlb3J5IiAvPgogICAgPFVzZXIgS2V5PSJzb2Z0d2FyZVxtaWNyb3NvZnRcb2ZmaWNlXDE2LjBcZXhjZWxcb3B0aW9ucyIgTmFtZT0iZGVmYXVsdGZvcm1hdCIgVmFsdWU9IjUxIiBUeXBlPSJSRUdfRFdPUkQiIEFwcD0iZXhjZWwxNiIgSWQ9IkxfU2F2ZUV4Y2VsZmlsZXNhcyIgLz4KICAgIDxVc2VyIEtleT0ic29mdHdhcmVcbWljcm9zb2Z0XG9mZmljZVwxNi4wXHBvd2VycG9pbnRcb3B0aW9ucyIgTmFtZT0iZGVmYXVsdGZvcm1hdCIgVmFsdWU9IjI3IiBUeXBlPSJSRUdfRFdPUkQiIEFwcD0icHB0MTYiIElkPSJMX1NhdmVQb3dlclBvaW50ZmlsZXNhcyIgLz4KICAgIDxVc2VyIEtleT0ic29mdHdhcmVcbWljcm9zb2Z0XG9mZmljZVwxNi4wXHdvcmRcb3B0aW9ucyIgTmFtZT0iZGVmYXVsdGZvcm1hdCIgVmFsdWU9IiIgVHlwZT0iUkVHX1NaIiBBcHA9IndvcmQxNiIgSWQ9IkxfU2F2ZVdvcmRmaWxlc2FzIiAvPgogIDwvQXBwU2V0dGluZ3M+CiAgPERpc3BsYXkgTGV2ZWw9IkZ1bGwiIEFjY2VwdEVVTEE9IlRSVUUiIC8+CjwvQ29uZmlndXJhdGlvbj4="
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/managedAppPolicies?_=1760087442283
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/managedAppPolicies",
    "value": [
        {
            "@odata.type": "#microsoft.graph.androidManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-Android | UnManaged Devices v1.10 - Printing Disabled",
            "description": "MAM WE Device Types\n\n30.05.2022",
            "createdDateTime": "2021-06-04T08:01:31.7845006Z",
            "lastModifiedDateTime": "2022-05-30T14:46:34Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "T_b149f8ef-4244-4ed0-8536-4fd975eda462",
            "version": "\"b5002b09-0000-0d00-0000-6294d8ca0000\"",
            "periodOfflineBeforeAccessCheck": "PT12H",
            "periodOnlineBeforeAccessCheck": "PT30M",
            "allowedInboundDataTransferSources": "managedApps",
            "allowedOutboundDataTransferDestinations": "managedApps",
            "organizationalCredentialsRequired": false,
            "allowedOutboundClipboardSharingLevel": "managedAppsWithPasteIn",
            "dataBackupBlocked": true,
            "deviceComplianceRequired": true,
            "managedBrowserToOpenLinksRequired": true,
            "saveAsBlocked": true,
            "periodOfflineBeforeWipeIsEnforced": "P90D",
            "pinRequired": true,
            "maximumPinRetries": 5,
            "simplePinBlocked": false,
            "minimumPinLength": 6,
            "pinCharacterSet": "numeric",
            "periodBeforePinReset": "PT0S",
            "allowedDataStorageLocations": [
                "oneDriveForBusiness",
                "sharePoint"
            ],
            "contactSyncBlocked": true,
            "printBlocked": true,
            "fingerprintBlocked": false,
            "disableAppPinIfDevicePinIsSet": false,
            "maximumRequiredOsVersion": null,
            "maximumWarningOsVersion": null,
            "maximumWipeOsVersion": null,
            "minimumRequiredOsVersion": null,
            "minimumWarningOsVersion": null,
            "minimumRequiredAppVersion": null,
            "minimumWarningAppVersion": null,
            "minimumWipeOsVersion": null,
            "minimumWipeAppVersion": null,
            "appActionIfDeviceComplianceRequired": "block",
            "appActionIfMaximumPinRetriesExceeded": "block",
            "pinRequiredInsteadOfBiometricTimeout": "PT30M",
            "allowedOutboundClipboardSharingExceptionLength": 0,
            "notificationRestriction": "allow",
            "previousPinBlockCount": 0,
            "managedBrowser": "microsoftEdge",
            "maximumAllowedDeviceThreatLevel": "medium",
            "mobileThreatDefenseRemediationAction": "block",
            "mobileThreatDefensePartnerPriority": null,
            "blockDataIngestionIntoOrganizationDocuments": true,
            "allowedDataIngestionLocations": [
                "oneDriveForBusiness",
                "sharePoint",
                "camera"
            ],
            "appActionIfUnableToAuthenticateUser": "wipe",
            "dialerRestrictionLevel": "allApps",
            "gracePeriodToBlockAppsDuringOffClockHours": null,
            "protectedMessagingRedirectAppType": "anyApp",
            "isAssigned": true,
            "targetedAppManagementLevels": "unmanaged",
            "appGroupType": "selectedPublicApps",
            "screenCaptureBlocked": true,
            "disableAppEncryptionIfDeviceEncryptionIsEnabled": false,
            "encryptAppData": true,
            "deployedAppCount": 39,
            "minimumRequiredPatchVersion": "0000-00-00",
            "minimumWarningPatchVersion": "0000-00-00",
            "minimumWipePatchVersion": "0000-00-00",
            "allowedAndroidDeviceManufacturers": null,
            "appActionIfAndroidDeviceManufacturerNotAllowed": "block",
            "appActionIfAccountIsClockedOut": null,
            "appActionIfSamsungKnoxAttestationRequired": null,
            "requiredAndroidSafetyNetDeviceAttestationType": "basicIntegrityAndDeviceCertification",
            "appActionIfAndroidSafetyNetDeviceAttestationFailed": "block",
            "requiredAndroidSafetyNetAppsVerificationType": "enabled",
            "appActionIfAndroidSafetyNetAppsVerificationFailed": "block",
            "customBrowserPackageId": "",
            "customBrowserDisplayName": "",
            "minimumRequiredCompanyPortalVersion": null,
            "minimumWarningCompanyPortalVersion": null,
            "minimumWipeCompanyPortalVersion": null,
            "keyboardsRestricted": true,
            "allowedAndroidDeviceModels": [],
            "appActionIfAndroidDeviceModelNotAllowed": "block",
            "customDialerAppPackageId": "",
            "customDialerAppDisplayName": "",
            "biometricAuthenticationBlocked": false,
            "requiredAndroidSafetyNetEvaluationType": "basic",
            "blockAfterCompanyPortalUpdateDeferralInDays": 0,
            "warnAfterCompanyPortalUpdateDeferralInDays": 0,
            "wipeAfterCompanyPortalUpdateDeferralInDays": 0,
            "deviceLockRequired": true,
            "appActionIfDeviceLockNotSet": "block",
            "connectToVpnOnLaunch": false,
            "appActionIfDevicePasscodeComplexityLessThanLow": null,
            "appActionIfDevicePasscodeComplexityLessThanMedium": null,
            "appActionIfDevicePasscodeComplexityLessThanHigh": null,
            "requireClass3Biometrics": false,
            "requirePinAfterBiometricChange": false,
            "fingerprintAndBiometricEnabled": null,
            "messagingRedirectAppPackageId": null,
            "messagingRedirectAppDisplayName": null,
            "exemptedAppPackages": [
                {
                    "name": "Google Maps",
                    "value": "com.google.maps"
                },
                {
                    "name": "Google Earth",
                    "value": "com.google.earth"
                },
                {
                    "name": "Google Play Music",
                    "value": "com.google.android.music"
                },
                {
                    "name": "Cisco WebEx",
                    "value": "com.cisco.webex.meetings"
                },
                {
                    "name": "Android SMS",
                    "value": "com.google.android.apps.messaging"
                },
                {
                    "name": "Android MMS ",
                    "value": "com.android.mms"
                },
                {
                    "name": "Samsung SMS",
                    "value": "com.Samsung.android.messaging"
                },
                {
                    "name": "Certificate Installer",
                    "value": "com.android.certinstaller"
                },
                {
                    "name": "Microsoft Authenticator",
                    "value": "com.microsoft.authenticator"
                }
            ],
            "approvedKeyboards": [
                {
                    "name": "com.google.android.inputmethod.latin",
                    "value": "Gboard - the Google Keyboard"
                },
                {
                    "name": "com.touchtype.swiftkey",
                    "value": "SwiftKey Keyboard"
                },
                {
                    "name": "com.sec.android.inputmethod",
                    "value": "Samsung Keyboard"
                },
                {
                    "name": "com.google.android.apps.inputmethod.hindi",
                    "value": "Google Indic Keyboard"
                },
                {
                    "name": "com.google.android.inputmethod.pinyin",
                    "value": "Google Pinyin Input"
                },
                {
                    "name": "com.google.android.inputmethod.japanese",
                    "value": "Google Japanese Input"
                },
                {
                    "name": "com.google.android.inputmethod.korean",
                    "value": "Google Korean Input"
                },
                {
                    "name": "com.google.android.apps.handwriting.ime",
                    "value": "Google Handwriting Input"
                },
                {
                    "name": "com.google.android.googlequicksearchbox",
                    "value": "Google voice typing"
                },
                {
                    "name": "com.samsung.android.svoiceime",
                    "value": "Samsung voice input"
                },
                {
                    "name": "com.samsung.android.honeyboard",
                    "value": "Samsung Keyboard"
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-iOS | UnManaged Devices v1.00 - Printing Disabled",
            "description": "08.10.2021",
            "createdDateTime": "2021-06-04T08:57:57.775389Z",
            "lastModifiedDateTime": "2022-01-05T00:39:51Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "T_7f2f1c83-0f32-41c0-aee5-ec44e1e83a83",
            "version": "\"4900bdbc-0000-0d00-0000-61d4e8d70000\"",
            "periodOfflineBeforeAccessCheck": "PT12H",
            "periodOnlineBeforeAccessCheck": "PT5M",
            "allowedInboundDataTransferSources": "managedApps",
            "allowedOutboundDataTransferDestinations": "managedApps",
            "organizationalCredentialsRequired": false,
            "allowedOutboundClipboardSharingLevel": "managedAppsWithPasteIn",
            "dataBackupBlocked": true,
            "deviceComplianceRequired": true,
            "managedBrowserToOpenLinksRequired": true,
            "saveAsBlocked": true,
            "periodOfflineBeforeWipeIsEnforced": "P90D",
            "pinRequired": true,
            "maximumPinRetries": 5,
            "simplePinBlocked": false,
            "minimumPinLength": 6,
            "pinCharacterSet": "numeric",
            "periodBeforePinReset": "PT0S",
            "allowedDataStorageLocations": [
                "oneDriveForBusiness",
                "sharePoint"
            ],
            "contactSyncBlocked": true,
            "printBlocked": true,
            "fingerprintBlocked": false,
            "disableAppPinIfDevicePinIsSet": false,
            "maximumRequiredOsVersion": null,
            "maximumWarningOsVersion": null,
            "maximumWipeOsVersion": null,
            "minimumRequiredOsVersion": "14.8",
            "minimumWarningOsVersion": null,
            "minimumRequiredAppVersion": null,
            "minimumWarningAppVersion": null,
            "minimumWipeOsVersion": null,
            "minimumWipeAppVersion": null,
            "appActionIfDeviceComplianceRequired": "block",
            "appActionIfMaximumPinRetriesExceeded": "block",
            "pinRequiredInsteadOfBiometricTimeout": "PT30M",
            "allowedOutboundClipboardSharingExceptionLength": 0,
            "notificationRestriction": "allow",
            "previousPinBlockCount": 0,
            "managedBrowser": "microsoftEdge",
            "maximumAllowedDeviceThreatLevel": "medium",
            "mobileThreatDefenseRemediationAction": "block",
            "mobileThreatDefensePartnerPriority": null,
            "blockDataIngestionIntoOrganizationDocuments": true,
            "allowedDataIngestionLocations": [
                "oneDriveForBusiness",
                "sharePoint",
                "camera"
            ],
            "appActionIfUnableToAuthenticateUser": "wipe",
            "dialerRestrictionLevel": "allApps",
            "gracePeriodToBlockAppsDuringOffClockHours": null,
            "protectedMessagingRedirectAppType": "anyApp",
            "isAssigned": true,
            "targetedAppManagementLevels": "unmanaged",
            "appGroupType": "selectedPublicApps",
            "genmojiConfigurationState": null,
            "screenCaptureConfigurationState": null,
            "writingToolsConfigurationState": null,
            "appDataEncryptionType": "whenDeviceLocked",
            "minimumRequiredSdkVersion": null,
            "deployedAppCount": 31,
            "faceIdBlocked": false,
            "allowWidgetContentSync": false,
            "minimumWipeSdkVersion": null,
            "allowedIosDeviceModels": null,
            "appActionIfIosDeviceModelNotAllowed": "block",
            "appActionIfAccountIsClockedOut": null,
            "thirdPartyKeyboardsBlocked": true,
            "filterOpenInToOnlyManagedApps": true,
            "disableProtectionOfManagedOutboundOpenInData": false,
            "protectInboundDataFromUnknownSources": false,
            "customBrowserProtocol": "",
            "customDialerAppProtocol": "",
            "managedUniversalLinks": [
                "http://*.sharepoint.com/*",
                "http://*.sharepoint-df.com/*",
                "http://*.yammer.com/*",
                "http://*.onedrive.com/*",
                "http://tasks.office.com/*",
                "http://to-do.microsoft.com/sharing*",
                "http://web.microsoftstream.com/video/*",
                "http://msit.microsoftstream.com/video/*",
                "http://*.powerbi.com/*",
                "http://app.powerbi.cn/*",
                "http://app.powerbigov.us/*",
                "http://app.powerbi.de/*",
                "http://*.service-now.com/*",
                "http://*.appsplatform.us/*",
                "http://*.powerapps.cn/*",
                "http://*.powerapps.com/*",
                "http://*.powerapps.us/*",
                "http://*teams.microsoft.com/l/*",
                "http://*devspaces.skype.com/l/*",
                "http://*teams.live.com/l/*",
                "http://*collab.apps.mil/l/*",
                "http://*teams.microsoft.us/l/*",
                "http://*teams-fl.microsoft.com/l/*",
                "http://*.zoom.us/*",
                "http://zoom.us/*",
                "https://*.sharepoint.com/*",
                "https://*.sharepoint-df.com/*",
                "https://*.yammer.com/*",
                "https://*.onedrive.com/*",
                "https://tasks.office.com/*",
                "https://to-do.microsoft.com/sharing*",
                "https://web.microsoftstream.com/video/*",
                "https://msit.microsoftstream.com/video/*",
                "https://*.powerbi.com/*",
                "https://app.powerbi.cn/*",
                "https://app.powerbigov.us/*",
                "https://app.powerbi.de/*",
                "https://*.service-now.com/*",
                "https://*.appsplatform.us/*",
                "https://*.powerapps.cn/*",
                "https://*.powerapps.com/*",
                "https://*.powerapps.us/*",
                "https://*teams.microsoft.com/l/*",
                "https://*devspaces.skype.com/l/*",
                "https://*teams.live.com/l/*",
                "https://*collab.apps.mil/l/*",
                "https://*teams.microsoft.us/l/*",
                "https://*teams-fl.microsoft.com/l/*",
                "https://*.zoom.us/*",
                "https://zoom.us/*"
            ],
            "exemptedUniversalLinks": [
                "http://maps.apple.com",
                "https://maps.apple.com",
                "http://facetime.apple.com",
                "https://facetime.apple.com"
            ],
            "minimumWarningSdkVersion": null,
            "messagingRedirectAppUrlScheme": null,
            "exemptedAppProtocols": [
                {
                    "name": "Default",
                    "value": "skype;app-settings;calshow;itms;itmss;itms-apps;itms-appss;itms-services;"
                },
                {
                    "name": "Apple Maps",
                    "value": "com.apple.Maps"
                },
                {
                    "name": "Google Maps",
                    "value": "com.google.maps"
                },
                {
                    "name": "WebEx",
                    "value": "wbx"
                },
                {
                    "name": "Apple Messages",
                    "value": "com.apple.MobileSMS"
                },
                {
                    "name": "Zoom",
                    "value": "zoomus"
                },
                {
                    "name": "Google Meet",
                    "value": "gmeet"
                },
                {
                    "name": "BlueJeans",
                    "value": "bjn"
                },
                {
                    "name": "BlueJeans",
                    "value": "bjn-intunemam"
                },
                {
                    "name": "BlueJeans",
                    "value": "bjn-a2m"
                },
                {
                    "name": "Duo Security",
                    "value": "otpauth"
                },
                {
                    "name": "Mimecast",
                    "value": "com.mimecast.mobile.saml"
                },
                {
                    "name": "PDF Expert",
                    "value": "pdfe-callback"
                },
                {
                    "name": "Salesforce",
                    "value": "salesforce1"
                },
                {
                    "name": " Go To Meeting",
                    "value": "gotomeeting"
                },
                {
                    "name": "AutoCAD DWG Viewer and Edito",
                    "value": "autocad"
                },
                {
                    "name": "Slack",
                    "value": "slack"
                },
                {
                    "name": "Docusign",
                    "value": "Docusignit"
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.targetedManagedAppConfiguration",
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
            "@odata.type": "#microsoft.graph.targetedManagedAppConfiguration",
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
            "@odata.type": "#microsoft.graph.iosManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-iOS | Managed Devices v1.00 - Printing Enabled",
            "description": "08.10.2021",
            "createdDateTime": "2021-06-04T08:50:57.1285003Z",
            "lastModifiedDateTime": "2022-01-05T00:39:51Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "T_443f240c-14de-4538-a5a7-410625a9cbdc",
            "version": "\"4900bfbc-0000-0d00-0000-61d4e8d70000\"",
            "periodOfflineBeforeAccessCheck": "PT12H",
            "periodOnlineBeforeAccessCheck": "PT5M",
            "allowedInboundDataTransferSources": "managedApps",
            "allowedOutboundDataTransferDestinations": "managedApps",
            "organizationalCredentialsRequired": false,
            "allowedOutboundClipboardSharingLevel": "managedAppsWithPasteIn",
            "dataBackupBlocked": true,
            "deviceComplianceRequired": true,
            "managedBrowserToOpenLinksRequired": true,
            "saveAsBlocked": true,
            "periodOfflineBeforeWipeIsEnforced": "P90D",
            "pinRequired": true,
            "maximumPinRetries": 5,
            "simplePinBlocked": false,
            "minimumPinLength": 6,
            "pinCharacterSet": "numeric",
            "periodBeforePinReset": "PT0S",
            "allowedDataStorageLocations": [
                "oneDriveForBusiness",
                "sharePoint"
            ],
            "contactSyncBlocked": true,
            "printBlocked": false,
            "fingerprintBlocked": false,
            "disableAppPinIfDevicePinIsSet": false,
            "maximumRequiredOsVersion": null,
            "maximumWarningOsVersion": null,
            "maximumWipeOsVersion": null,
            "minimumRequiredOsVersion": "14.8",
            "minimumWarningOsVersion": null,
            "minimumRequiredAppVersion": null,
            "minimumWarningAppVersion": null,
            "minimumWipeOsVersion": null,
            "minimumWipeAppVersion": null,
            "appActionIfDeviceComplianceRequired": "block",
            "appActionIfMaximumPinRetriesExceeded": "block",
            "pinRequiredInsteadOfBiometricTimeout": "PT30M",
            "allowedOutboundClipboardSharingExceptionLength": 0,
            "notificationRestriction": "allow",
            "previousPinBlockCount": 0,
            "managedBrowser": "microsoftEdge",
            "maximumAllowedDeviceThreatLevel": "medium",
            "mobileThreatDefenseRemediationAction": "block",
            "mobileThreatDefensePartnerPriority": null,
            "blockDataIngestionIntoOrganizationDocuments": true,
            "allowedDataIngestionLocations": [
                "oneDriveForBusiness",
                "sharePoint",
                "camera"
            ],
            "appActionIfUnableToAuthenticateUser": "block",
            "dialerRestrictionLevel": "allApps",
            "gracePeriodToBlockAppsDuringOffClockHours": null,
            "protectedMessagingRedirectAppType": "anyApp",
            "isAssigned": true,
            "targetedAppManagementLevels": "mdm",
            "appGroupType": "selectedPublicApps",
            "genmojiConfigurationState": null,
            "screenCaptureConfigurationState": null,
            "writingToolsConfigurationState": null,
            "appDataEncryptionType": "whenDeviceLocked",
            "minimumRequiredSdkVersion": null,
            "deployedAppCount": 31,
            "faceIdBlocked": false,
            "allowWidgetContentSync": false,
            "minimumWipeSdkVersion": null,
            "allowedIosDeviceModels": null,
            "appActionIfIosDeviceModelNotAllowed": "block",
            "appActionIfAccountIsClockedOut": null,
            "thirdPartyKeyboardsBlocked": true,
            "filterOpenInToOnlyManagedApps": true,
            "disableProtectionOfManagedOutboundOpenInData": false,
            "protectInboundDataFromUnknownSources": false,
            "customBrowserProtocol": "",
            "customDialerAppProtocol": "",
            "managedUniversalLinks": [
                "http://*.sharepoint.com/*",
                "http://*.sharepoint-df.com/*",
                "http://*.yammer.com/*",
                "http://*.onedrive.com/*",
                "http://tasks.office.com/*",
                "http://to-do.microsoft.com/sharing*",
                "http://web.microsoftstream.com/video/*",
                "http://msit.microsoftstream.com/video/*",
                "http://*.powerbi.com/*",
                "http://app.powerbi.cn/*",
                "http://app.powerbigov.us/*",
                "http://app.powerbi.de/*",
                "http://*.service-now.com/*",
                "http://*.appsplatform.us/*",
                "http://*.powerapps.cn/*",
                "http://*.powerapps.com/*",
                "http://*.powerapps.us/*",
                "http://*teams.microsoft.com/l/*",
                "http://*devspaces.skype.com/l/*",
                "http://*teams.live.com/l/*",
                "http://*collab.apps.mil/l/*",
                "http://*teams.microsoft.us/l/*",
                "http://*teams-fl.microsoft.com/l/*",
                "http://*.zoom.us/*",
                "http://zoom.us/*",
                "https://*.sharepoint.com/*",
                "https://*.sharepoint-df.com/*",
                "https://*.yammer.com/*",
                "https://*.onedrive.com/*",
                "https://tasks.office.com/*",
                "https://to-do.microsoft.com/sharing*",
                "https://web.microsoftstream.com/video/*",
                "https://msit.microsoftstream.com/video/*",
                "https://*.powerbi.com/*",
                "https://app.powerbi.cn/*",
                "https://app.powerbigov.us/*",
                "https://app.powerbi.de/*",
                "https://*.service-now.com/*",
                "https://*.appsplatform.us/*",
                "https://*.powerapps.cn/*",
                "https://*.powerapps.com/*",
                "https://*.powerapps.us/*",
                "https://*teams.microsoft.com/l/*",
                "https://*devspaces.skype.com/l/*",
                "https://*teams.live.com/l/*",
                "https://*collab.apps.mil/l/*",
                "https://*teams.microsoft.us/l/*",
                "https://*teams-fl.microsoft.com/l/*",
                "https://*.zoom.us/*",
                "https://zoom.us/*"
            ],
            "exemptedUniversalLinks": [
                "http://maps.apple.com",
                "https://maps.apple.com",
                "http://facetime.apple.com",
                "https://facetime.apple.com"
            ],
            "minimumWarningSdkVersion": null,
            "messagingRedirectAppUrlScheme": null,
            "exemptedAppProtocols": [
                {
                    "name": "Default",
                    "value": "skype;app-settings;calshow;itms;itmss;itms-apps;itms-appss;itms-services;"
                },
                {
                    "name": "Apple Maps",
                    "value": "com.apple.Maps"
                },
                {
                    "name": "Apple Messages",
                    "value": "com.apple.MobileSMS"
                },
                {
                    "name": "Google Maps",
                    "value": "comgooglemaps"
                },
                {
                    "name": "WebEx",
                    "value": "wbx"
                },
                {
                    "name": "Zoom",
                    "value": "zoomus"
                },
                {
                    "name": "Google Meet",
                    "value": "gmeet"
                },
                {
                    "name": "iOS BlueJeans",
                    "value": "bjn"
                },
                {
                    "name": "iOS BlueJeans",
                    "value": "bjn-intunemam"
                },
                {
                    "name": "iOS BlueJeans",
                    "value": "bjn-a2m"
                },
                {
                    "name": "iOS Duo Security",
                    "value": "otpauth"
                },
                {
                    "name": "Mimecast",
                    "value": "com.mimecast.mobile.saml"
                },
                {
                    "name": "PDF Expert",
                    "value": "pdfe-callback"
                },
                {
                    "name": "Salesforce",
                    "value": "salesforce1"
                },
                {
                    "name": "Go To Meeting",
                    "value": "gotomeeting"
                },
                {
                    "name": "AutoCAD DWG Viewer and Editor",
                    "value": "autocad"
                },
                {
                    "name": "Slack",
                    "value": "slack"
                },
                {
                    "name": "Docusign",
                    "value": "Docusignit"
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.androidManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-Android | Managed Devices v1.00 - Printing Enabled",
            "description": "08.10.2021",
            "createdDateTime": "2021-05-25T06:29:37.7011924Z",
            "lastModifiedDateTime": "2022-03-02T18:48:25Z",
            "roleScopeTagIds": [
                "0"
            ],
            "id": "T_764579ab-93b2-4061-8f21-59cd221cbfe4",
            "version": "\"1a0383b9-0000-0d00-0000-621fbbf90000\"",
            "periodOfflineBeforeAccessCheck": "PT12H",
            "periodOnlineBeforeAccessCheck": "PT30M",
            "allowedInboundDataTransferSources": "managedApps",
            "allowedOutboundDataTransferDestinations": "managedApps",
            "organizationalCredentialsRequired": false,
            "allowedOutboundClipboardSharingLevel": "managedAppsWithPasteIn",
            "dataBackupBlocked": true,
            "deviceComplianceRequired": true,
            "managedBrowserToOpenLinksRequired": true,
            "saveAsBlocked": true,
            "periodOfflineBeforeWipeIsEnforced": "P90D",
            "pinRequired": true,
            "maximumPinRetries": 5,
            "simplePinBlocked": false,
            "minimumPinLength": 6,
            "pinCharacterSet": "numeric",
            "periodBeforePinReset": "PT0S",
            "allowedDataStorageLocations": [
                "oneDriveForBusiness",
                "sharePoint"
            ],
            "contactSyncBlocked": true,
            "printBlocked": false,
            "fingerprintBlocked": false,
            "disableAppPinIfDevicePinIsSet": false,
            "maximumRequiredOsVersion": null,
            "maximumWarningOsVersion": null,
            "maximumWipeOsVersion": null,
            "minimumRequiredOsVersion": null,
            "minimumWarningOsVersion": null,
            "minimumRequiredAppVersion": null,
            "minimumWarningAppVersion": null,
            "minimumWipeOsVersion": null,
            "minimumWipeAppVersion": null,
            "appActionIfDeviceComplianceRequired": "block",
            "appActionIfMaximumPinRetriesExceeded": "block",
            "pinRequiredInsteadOfBiometricTimeout": "PT30M",
            "allowedOutboundClipboardSharingExceptionLength": 250,
            "notificationRestriction": "allow",
            "previousPinBlockCount": 0,
            "managedBrowser": "microsoftEdge",
            "maximumAllowedDeviceThreatLevel": "medium",
            "mobileThreatDefenseRemediationAction": "block",
            "mobileThreatDefensePartnerPriority": null,
            "blockDataIngestionIntoOrganizationDocuments": true,
            "allowedDataIngestionLocations": [
                "oneDriveForBusiness",
                "sharePoint",
                "camera"
            ],
            "appActionIfUnableToAuthenticateUser": "block",
            "dialerRestrictionLevel": "allApps",
            "gracePeriodToBlockAppsDuringOffClockHours": null,
            "protectedMessagingRedirectAppType": "anyApp",
            "isAssigned": false,
            "targetedAppManagementLevels": "mdm,androidEnterprise",
            "appGroupType": "selectedPublicApps",
            "screenCaptureBlocked": true,
            "disableAppEncryptionIfDeviceEncryptionIsEnabled": false,
            "encryptAppData": true,
            "deployedAppCount": 35,
            "minimumRequiredPatchVersion": "0000-00-00",
            "minimumWarningPatchVersion": "0000-00-00",
            "minimumWipePatchVersion": "0000-00-00",
            "allowedAndroidDeviceManufacturers": null,
            "appActionIfAndroidDeviceManufacturerNotAllowed": "block",
            "appActionIfAccountIsClockedOut": null,
            "appActionIfSamsungKnoxAttestationRequired": null,
            "requiredAndroidSafetyNetDeviceAttestationType": "basicIntegrityAndDeviceCertification",
            "appActionIfAndroidSafetyNetDeviceAttestationFailed": "block",
            "requiredAndroidSafetyNetAppsVerificationType": "none",
            "appActionIfAndroidSafetyNetAppsVerificationFailed": "block",
            "customBrowserPackageId": "",
            "customBrowserDisplayName": "",
            "minimumRequiredCompanyPortalVersion": null,
            "minimumWarningCompanyPortalVersion": null,
            "minimumWipeCompanyPortalVersion": null,
            "keyboardsRestricted": true,
            "allowedAndroidDeviceModels": [],
            "appActionIfAndroidDeviceModelNotAllowed": "block",
            "customDialerAppPackageId": "",
            "customDialerAppDisplayName": "",
            "biometricAuthenticationBlocked": false,
            "requiredAndroidSafetyNetEvaluationType": "basic",
            "blockAfterCompanyPortalUpdateDeferralInDays": 0,
            "warnAfterCompanyPortalUpdateDeferralInDays": 0,
            "wipeAfterCompanyPortalUpdateDeferralInDays": 0,
            "deviceLockRequired": false,
            "appActionIfDeviceLockNotSet": "block",
            "connectToVpnOnLaunch": false,
            "appActionIfDevicePasscodeComplexityLessThanLow": null,
            "appActionIfDevicePasscodeComplexityLessThanMedium": null,
            "appActionIfDevicePasscodeComplexityLessThanHigh": null,
            "requireClass3Biometrics": false,
            "requirePinAfterBiometricChange": false,
            "fingerprintAndBiometricEnabled": null,
            "messagingRedirectAppPackageId": null,
            "messagingRedirectAppDisplayName": null,
            "exemptedAppPackages": [
                {
                    "name": "Google Maps",
                    "value": "com.google.maps"
                },
                {
                    "name": "Google Earth ",
                    "value": "com.google.earth"
                },
                {
                    "name": "Google Play Music",
                    "value": "com.google.android.music"
                },
                {
                    "name": "Cisco WebEx",
                    "value": "com.cisco.webex.meetings"
                },
                {
                    "name": "Android SMS",
                    "value": "com.google.android.apps.messaging"
                },
                {
                    "name": "Android MMS ",
                    "value": "com.android.mms"
                },
                {
                    "name": "Samsung SMS",
                    "value": "com.Samsung.android.messaging"
                },
                {
                    "name": "Certificate Installer",
                    "value": "com.android.certinstaller"
                }
            ],
            "approvedKeyboards": [
                {
                    "name": "com.google.android.inputmethod.latin",
                    "value": "Gboard - the Google Keyboard"
                },
                {
                    "name": "com.touchtype.swiftkey",
                    "value": "SwiftKey Keyboard"
                },
                {
                    "name": "com.sec.android.inputmethod",
                    "value": "Samsung Keyboard"
                },
                {
                    "name": "com.google.android.apps.inputmethod.hindi",
                    "value": "Google Indic Keyboard"
                },
                {
                    "name": "com.google.android.inputmethod.pinyin",
                    "value": "Google Pinyin Input"
                },
                {
                    "name": "com.google.android.inputmethod.japanese",
                    "value": "Google Japanese Input"
                },
                {
                    "name": "com.google.android.inputmethod.korean",
                    "value": "Google Korean Input"
                },
                {
                    "name": "com.google.android.apps.handwriting.ime",
                    "value": "Google Handwriting Input"
                },
                {
                    "name": "com.google.android.googlequicksearchbox",
                    "value": "Google voice typing"
                },
                {
                    "name": "com.samsung.android.svoiceime",
                    "value": "Samsung voice input"
                },
                {
                    "name": "com.samsung.android.honeyboard",
                    "value": "Samsung Keyboard"
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.targetedManagedAppConfiguration",
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
        },
        {
            "@odata.type": "#microsoft.graph.targetedManagedAppConfiguration",
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
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations?$select=id,displayName,lastModifiedDateTime,roleScopeTagIds,microsoft.graph.unsupportedDeviceConfiguration/originalEntityTypeName&&top=1000
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations(id,displayName,lastModifiedDateTime,roleScopeTagIds,microsoft.graph.unsupportedDeviceConfiguration/originalEntityTypeName)",
    "value": [
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "3fda1896-0033-4632-b3bf-4b1bbc6f2bf7",
            "displayName": "[Base] Deprecated | Windows  - Custom | Windows Telemetry ver1.0",
            "lastModifiedDateTime": "2023-02-23T08:24:06.1478533Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "17f5f062-fb52-403a-b7c6-2d8c5625c42f",
            "displayName": "[Base] Deprecated | Windows - Custom | DeliveryOptimization - M365",
            "lastModifiedDateTime": "2023-02-23T08:24:55.5010421Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "a3f7fa2d-da19-435e-9251-fe39efebeb01",
            "displayName": "[Base] Deprecated | Windows - Custom | Diagnostics Collection ver1.0",
            "lastModifiedDateTime": "2023-02-23T08:25:34.1994276Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "4763926a-629b-455d-9c46-91203071176b",
            "displayName": "[Base] Deprecated | Windows - Custom | Internet Explorer 11 ver1.0",
            "lastModifiedDateTime": "2023-02-23T08:24:38.5961925Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "720cb9c9-8872-4e28-a60c-de28819dc878",
            "displayName": "[Base] Deprecated | Windows - Custom | Microsoft Edge - AppAssociations ver1.0",
            "lastModifiedDateTime": "2023-02-23T08:25:46.0780916Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "54b8f613-3820-4e99-b11d-1ef31193f33e",
            "displayName": "[Base] Deprecated | Windows - Custom | Security Audit Settings ver1.0",
            "lastModifiedDateTime": "2023-02-23T08:24:23.698057Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "a9ee3ba9-2eff-4241-a29a-52b846a9a7fe",
            "displayName": "[Base] Dev | Windows - Custom | Device Control - Printer Protection [User Level] ver0.1",
            "lastModifiedDateTime": "2022-05-18T13:24:39.1853055Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "10dc6ef6-1881-417a-9bb1-d71606d5824f",
            "displayName": "[Base] Dev | Windows - Custom | Microsoft Desktop App Installer [Winget] ADMX ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:31:07.1038815Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "26188faf-d631-4aa9-8073-a66f4f215f8d",
            "displayName": "[Base] Dev | Windows - Custom | Mozilla Firefox ver0.1",
            "lastModifiedDateTime": "2022-05-18T13:32:49.5311821Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10DeviceFirmwareConfigurationInterface",
            "id": "2b71bd8b-0125-4c80-8f73-e1fe154f6420",
            "displayName": "[Base] Dev | Windows - Device Firmware Configuration Interface |  ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:19:09.7026966Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.editionUpgradeConfiguration",
            "id": "15aabf7f-2eb5-4cbd-8094-ec07e092c328",
            "displayName": "[Base] Dev | Windows - Edition Upgrade and Mode Switch | Windows 10 Pro To Enterprise ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:28:59.1349477Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10GeneralConfiguration",
            "id": "1892a735-2af0-466a-8b1d-6ca96ead05b3",
            "displayName": "[Base] DoNotUse | Windows - Device Restrictions | Password Complexity ver1.0",
            "lastModifiedDateTime": "2022-09-19T11:12:35.1236778Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "5ddce99c-e84c-4b08-9030-f2c3995c36e4",
            "displayName": "[Base] Physical | Windows  - Custom | WiFi Connectivity: PLUSNET-M2MJ ver1.1",
            "lastModifiedDateTime": "2022-06-02T18:50:10.7827551Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "797154f2-06c6-4195-bd3c-42d18b562d0b",
            "displayName": "[Base] Physical | Windows - Custom | WiFi Connectivity: DT-Home_5G-2 ver1.0",
            "lastModifiedDateTime": "2023-01-15T17:27:11.4594542Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "f220b238-4d27-4687-a0de-d448ec503854",
            "displayName": "[Base] Physical | Windows - Custom | WiFi Connectivity: DT-IoT ver1.0",
            "lastModifiedDateTime": "2023-01-15T17:28:43.4419134Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "c0e47cdf-18b6-4f47-9da7-a097dbbc94ae",
            "displayName": "[Base] Physical | Windows - Custom | WiFi Connectivity: DT-Office ver1.0",
            "lastModifiedDateTime": "2023-01-15T17:28:15.6760473Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "170b0a81-16ae-4642-bd2b-146a9d8346a0",
            "displayName": "[Base] Prod | Windows  - Custom | Windows Defender Application Control ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:55:21.5462232Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "fd8f88d9-e0f1-4c26-b158-8a095e9d9dc0",
            "displayName": "[Base] Prod | Windows - Custom | Applocker Rules - .appx ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:15:19.9679712Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "5077888c-ad7d-4382-8d2f-94990c2464d8",
            "displayName": "[Base] Prod | Windows - Custom | Applocker Rules - .dll ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:15:43.0663057Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "d026dbad-cada-46eb-a0bd-b62d40a11fc4",
            "displayName": "[Base] Prod | Windows - Custom | Applocker Rules - .exe ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:16:23.9889192Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "1c40de2a-2ea6-487b-8b3b-ca3e5323c252",
            "displayName": "[Base] Prod | Windows - Custom | Applocker Rules - .msi ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:17:06.8838761Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "6690f2d1-993d-4550-80ad-28462db05a18",
            "displayName": "[Base] Prod | Windows - Custom | Applocker Rules - Script ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:21:01.1519162Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "29669b6e-652d-460c-9a88-54a6e5c96c8e",
            "displayName": "[Base] Prod | Windows - Custom | Block IE11 App Access ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:21:49.3156023Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "31c3693a-29dd-4ac0-839b-420facfde682",
            "displayName": "[Base] Prod | Windows - Custom | Device Control - Printer Protection [Device Level] ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:24:09.886293Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "81194756-83ce-4584-a9fb-f868bfcc77b0",
            "displayName": "[Base] Prod | Windows - Custom | Enrollment Status Page-SkipAccountSetup ver1.0",
            "lastModifiedDateTime": "2022-05-23T21:55:58.3462435Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "f374718b-7017-474f-9d26-d24bf7958a26",
            "displayName": "[Base] Prod | Windows - Custom | Google Chrome ADMX-v2r2-STIG ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:28:00.3354992Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "8be07e90-1e87-4823-9b12-e730bd2611ed",
            "displayName": "[Base] Prod | Windows - Custom | Lenovo System Update for Windows 10 ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:29:22.3910186Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "7534b708-350b-40cc-a140-2cdc64d29040",
            "displayName": "[Base] Prod | Windows - Custom | Manage Endpoint Administrators ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:30:13.4912279Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "42d2a033-c8fd-4d26-90c0-e28688363158",
            "displayName": "[Base] Prod | Windows - Custom | Outlook 2016 Silent Config ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:33:11.737728Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "0fc6a704-662d-4334-a233-a272529aa0f9",
            "displayName": "[Base] Prod | Windows - Custom | System - Location & Sensors ver1.0",
            "lastModifiedDateTime": "2022-09-19T10:27:07.8436368Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "82c9f71e-5404-4c37-a91d-44dcc9ad5f53",
            "displayName": "[Base] Prod | Windows - Custom | Time Language Settings - GMT ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:38:29.6076559Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "81e5fbfc-6c6c-4959-a889-e7960356e8f3",
            "displayName": "[Base] Prod | Windows - Custom | Update Compliance ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:41:12.3779542Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "5fc52bc3-04a1-442f-9b69-6649f58708e5",
            "displayName": "[Base] Prod | Windows - Custom | Windows Components - Ink ver1.0",
            "lastModifiedDateTime": "2022-09-19T10:39:03.9387856Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "6f511f91-33ba-471a-a2da-6c467c0874cd",
            "displayName": "[Base] Prod | Windows - Custom | Windows Components - Messaging ver1.0",
            "lastModifiedDateTime": "2022-09-19T10:26:08.3444475Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10GeneralConfiguration",
            "id": "7c69eade-e71d-4c06-853d-9f798f917c81",
            "displayName": "[Base] Prod | Windows - Device Restrictions | General ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:19:54.5327916Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10GeneralConfiguration",
            "id": "ee5ec69f-f441-4c3f-8e4e-e6dd37d4c02f",
            "displayName": "[Base] Prod | Windows - Device Restrictions | Privacy ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:25:43.7558749Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10EasEmailProfileConfiguration",
            "id": "b0079b61-617d-4b3f-8222-09e641c75644",
            "displayName": "[Base] Prod | Windows - Email | Deployment Theory Mailbox ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:31:53.5798558Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10EndpointProtectionConfiguration",
            "id": "536b33d9-a8d2-4a62-891a-012a1194bc48",
            "displayName": "[Base] Prod | Windows - Endpoint Protection | Microsoft Defender Firewall - Teams Firewall Rules ver1.0",
            "lastModifiedDateTime": "2022-05-18T15:45:45.5704908Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsHealthMonitoringConfiguration",
            "id": "71ae573c-fa24-4918-bb5a-5474d4beaae2",
            "displayName": "[Base] Prod | Windows - WindowsHealthMonitoring | ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:10:15.3629544Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "54504b24-68a7-49f2-a458-1bd25b3179a7",
            "displayName": "[Base] Test | Windows - Custom | Advanced Audit Policy Configuration ver1.0",
            "lastModifiedDateTime": "2022-09-19T11:17:30.9304151Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "6114b8ed-e4ad-4c34-b58b-ccb43aff4366",
            "displayName": "[Base] Test | Windows - Custom | Device Naming Standard ver0.1",
            "lastModifiedDateTime": "2022-09-19T11:56:01.7705567Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "a10056cb-5a3d-46bf-9290-79ea13ae45f5",
            "displayName": "[Base] Test | Windows - Custom | Local Policies - Security Options ver0.1",
            "lastModifiedDateTime": "2022-09-19T11:21:23.1740633Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "43cc766d-70eb-4937-b800-b2bfd23af32b",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | App Store, Doc Viewing, Gaming",
            "lastModifiedDateTime": "2021-08-23T07:59:52.823524Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "8d42d785-b186-4c09-9022-ac4eb4b15521",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Autonomous Single App Mode",
            "lastModifiedDateTime": "2021-08-23T08:01:24.8460058Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "82d6e737-b2b9-43b2-8f27-fb9624b5131b",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Built-in Apps",
            "lastModifiedDateTime": "2021-08-23T08:04:40.2951449Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "27453c0d-d417-4d24-ab04-1b99cf39f795",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Cloud and Storage",
            "lastModifiedDateTime": "2021-08-23T08:07:51.0916633Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "c39fbc36-ea9a-4e7a-b994-c850c6190ac2",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Connected Devices",
            "lastModifiedDateTime": "2021-08-23T08:12:45.2130265Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "8ded8ee8-d161-402a-b989-8f57b8b44a70",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Domains",
            "lastModifiedDateTime": "2021-08-23T08:15:26.0786779Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "864d6eed-5feb-44cf-8640-eb0240f6a495",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | General",
            "lastModifiedDateTime": "2021-08-23T08:21:58.5421583Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "aaebcbde-1e2c-4a02-a2f0-36c08fbedee4",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Kiosk Mode With MS Edge",
            "lastModifiedDateTime": "2021-08-23T08:26:15.0429852Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "11b2267e-6ffd-4e0e-b838-e953f3bc60fb",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Locked Screen Experience",
            "lastModifiedDateTime": "2021-08-23T08:27:57.0132622Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "33d1c79b-0477-43e3-b9cc-f015f82dcf37",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Password [Standard Users]",
            "lastModifiedDateTime": "2021-08-23T08:40:31.0332239Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "61ccc106-a045-4fa7-b2a8-37e0eec18f28",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Restricted Apps",
            "lastModifiedDateTime": "2021-08-23T09:26:21.4379006Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "id": "b350fdfe-a7a9-4bdf-b8f9-f2281a384212",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Wireless",
            "lastModifiedDateTime": "2021-08-23T09:05:02.5603413Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosDeviceFeaturesConfiguration",
            "id": "e5eabd6a-2208-4bcb-b6c6-45afec6e9640",
            "displayName": "[Global] iOS/iPadOS - DeviceFeatures | App Notifications",
            "lastModifiedDateTime": "2021-08-23T14:43:57.259432Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosDeviceFeaturesConfiguration",
            "id": "433c3715-0f30-4407-a005-1e6d747d9259",
            "displayName": "[Global] iOS/iPadOS - DeviceFeatures | Enable and configure Microsoft Enterprise SSO plug-in",
            "lastModifiedDateTime": "2021-08-23T09:10:24.8329418Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosDeviceFeaturesConfiguration",
            "id": "624bcbee-be11-47fa-8084-d582e151a165",
            "displayName": "[Global] iOS/iPadOS - DeviceFeatures | Lock Screen Message",
            "lastModifiedDateTime": "2021-08-23T14:46:38.3984344Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosDeviceFeaturesConfiguration",
            "id": "1d68932d-0433-4d1a-bfdf-da3da0f93c03",
            "displayName": "[Global] iOS/iPadOS - DeviceFeatures | Single sign-on",
            "lastModifiedDateTime": "2021-08-23T14:51:56.1833239Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosEasEmailProfileConfiguration",
            "id": "036bc109-7d39-4e9b-b402-6775b7aa4611",
            "displayName": "[Global] iOS/iPadOS - Email | Deployment Theory",
            "lastModifiedDateTime": "2021-08-23T09:08:29.8859051Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosWiFiConfiguration",
            "id": "9e0561b9-673c-40da-8359-22a4ed0d0c13",
            "displayName": "[Global] iOS/iPadOS - WiFi | Orbi",
            "lastModifiedDateTime": "2021-08-23T09:07:25.7426924Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "0b73079a-c6ae-4a78-9fa4-1e74c6f805aa",
            "displayName": "[Kiosk] Dev | Windows - Custom | Enable Shared PC Mode ver0.1",
            "lastModifiedDateTime": "2022-05-18T14:03:30.1698118Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "be678d9b-f00e-4914-8b14-d7e811e353c5",
            "displayName": "[Physical] Prod | Windows  - Custom | WiFi Connectivity: ORBI92 ver1.1",
            "lastModifiedDateTime": "2022-06-02T18:49:19.8280799Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "6f2e9b58-9ee7-4cc7-a129-734083283c85",
            "displayName": "[Physical] Prod | Windows  - Custom | Windows Hello Dynamic Lock – Phone Proximity ver1.0",
            "lastModifiedDateTime": "2022-06-04T06:09:41.1596995Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "2392ce8b-34f8-4b96-b4cd-ed1a66f16bc9",
            "displayName": "[Physical] Prod | Windows  - Custom | Windows Hello Multifactor Unlock – 1st Unlock Factor ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:57:07.2052038Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "6f4f4f02-396c-4507-8c4b-84fa940b5cd7",
            "displayName": "[Physical] Prod | Windows  - Custom | Windows Hello Multifactor Unlock – Unlock Signal Rules ver1.0",
            "lastModifiedDateTime": "2022-05-18T13:58:05.5489203Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "a5b089c5-2be9-4234-b2a6-5f57a15ac804",
            "displayName": "[Physical] Prod | Windows  - Custom | Windows PIN Reset From Logon Screen ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:01:36.556653Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "1d168bd6-bdeb-4e14-911d-586b7ec8cb85",
            "displayName": "[Physical] Prod | Windows  - Custom | Windows PIN SignIn Allowed Urls ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:04:15.5716975Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsDeliveryOptimizationConfiguration",
            "id": "3a21df03-bc42-4443-af2b-911b7ecf798b",
            "displayName": "[Physical] Prod | Windows  - DeliveryOptimization | Windows OS ver1.0",
            "lastModifiedDateTime": "2022-05-18T14:06:46.7382441Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "46489c2b-b717-4c4b-9bb7-25e4fb392dc0",
            "displayName": "[Physical] Prod | Windows - Custom | Windows Hello Multifactor Unlock – 2nd Unlock Factor ver1.0",
            "lastModifiedDateTime": "2022-06-04T06:08:42.302879Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "a98f9bf8-e974-4cd3-9468-a56e6864b877",
            "displayName": "[SurfaceHub] Prod | Windows - Custom | Quality of Service for MS Teams ver1.0",
            "lastModifiedDateTime": "2022-05-19T10:36:14.258894Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CustomConfiguration",
            "id": "0139380f-338d-42c1-9161-5af1416bc000",
            "displayName": "[SurfaceHub] Prod | Windows - Custom | Quality of Service for Skype for Business ver1.0",
            "lastModifiedDateTime": "2022-05-19T10:36:02.7977176Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10TeamGeneralConfiguration",
            "id": "07d915d6-e4c9-4213-b265-f927512bcd6b",
            "displayName": "[SurfaceHub] Prod | Windows - Device restrictions (Windows 10 Team) ver1.0",
            "lastModifiedDateTime": "2022-05-19T10:35:45.3554909Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.macOSSoftwareUpdateConfiguration",
            "id": "5a50ba96-90d4-402f-9189-e9a7edbbb8d0",
            "displayName": "test",
            "lastModifiedDateTime": "2025-08-19T07:51:44.7540006Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosUpdateConfiguration",
            "id": "f859e6b9-8878-4c81-9ae5-1dbabbb91bb7",
            "displayName": "test",
            "lastModifiedDateTime": "2025-01-15T14:41:21.3903464Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "0691e9d5-81d6-493e-a698-6ff6c320234c",
            "displayName": "Win10/11 - Ring 00 | [TestDeployment]-[Engineering]-[WindowsInsider]-[DG]",
            "lastModifiedDateTime": "2022-04-16T08:34:07.5855185Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "fdedb0ab-62b8-4752-a8d4-25f7eb47240b",
            "displayName": "Win10/11 - Ring 01 | [TestDeployment]-[Engineering]-[DG]",
            "lastModifiedDateTime": "2022-04-16T08:34:25.4608799Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "273f2291-3c9f-44b4-9ec9-c8a36ab16184",
            "displayName": "Win10/11 - Ring 02 | [TestDeployment]-[Broad_Engineering_Pilot]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:40:53.9796824Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "9400b44b-28f1-4272-9e4f-150e9cb75670",
            "displayName": "Win10/11 - Ring 03 | [UATDeployment]-[Champions]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:41:09.6221966Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "bb65de46-a3b2-4eb0-95d0-2747fa401a84",
            "displayName": "Win10/11 - Ring 04 | [UATDeployment]-[Broad_User_Pilot]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:41:37.811352Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "d725e342-a2ad-4687-9b72-d8c293083b4d",
            "displayName": "Win10/11 - Ring 05 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:42:09.0057482Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "ae918ad4-eaa7-4a07-a4bf-8f8e193b4819",
            "displayName": "Win10/11 - Ring 06 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:42:29.9959949Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "29d7692f-70e9-4538-b30e-d88306ce0d42",
            "displayName": "Win10/11 - Ring 07 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:42:48.1838357Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "9043b678-6afd-4952-b3e7-9b4e00522adb",
            "displayName": "Win10/11 - Ring 08 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:43:03.4326219Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "f54bb312-0500-40e8-9ceb-1cbfa1758d32",
            "displayName": "Win10/11 - Ring 09 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-19T12:43:16.5081971Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "781f3b1d-3686-462f-a3c5-5834dddc0951",
            "displayName": "Win10/11 - Ring 10 | [ProdDeployment]-[Deferrals]-[DG]",
            "lastModifiedDateTime": "2022-04-16T08:36:56.659062Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "f28abe6a-b43e-41fa-8721-b7f91e91ea3e",
            "displayName": "Win10/11 Multisession - Ring 01 | [TestDeployment]-[Engineering]-[Windows]-[DG]",
            "lastModifiedDateTime": "2022-04-24T21:05:26.5341071Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "d225b398-2200-44c7-a3fc-a2a5c5dc1031",
            "displayName": "Win10/11 Multisession - Ring 02 | [UATDeployment]-[Broad_User_Pilot]-[DG]",
            "lastModifiedDateTime": "2022-04-24T21:07:44.5960048Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "18900079-f55a-4c0d-bc36-bfa292231714",
            "displayName": "Win10/11 Multisession - Ring 03 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-24T21:09:49.5605405Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "c0176aa5-d36d-4eb5-b1cc-b60507ccec3a",
            "displayName": "Win10/11 Multisession - Ring 04 | [ProdDeployment]-[Broad_Deployment]-[DG]",
            "lastModifiedDateTime": "2022-04-24T21:11:26.9298787Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "a4511a6c-4156-4119-8dfa-a98132446f51",
            "displayName": "Windows Autopatch Update Policy - auto-patch-group - Last",
            "lastModifiedDateTime": "2025-08-21T14:47:36.4651846Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "296ce441-5e93-40f3-a5c3-a7d17bba14f4",
            "displayName": "Windows Autopatch Update Policy - auto-patch-group - Ring1",
            "lastModifiedDateTime": "2025-08-21T14:47:36.4877527Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "10cadad0-feb3-41e3-94d2-fab3028d24c7",
            "displayName": "Windows Autopatch Update Policy - auto-patch-group - Test",
            "lastModifiedDateTime": "2025-08-21T14:47:36.4799179Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "14acfdfa-5d07-4db9-8ae4-2d6473e759d7",
            "displayName": "Windows Autopatch Update Policy - Windows Autopatch - Last",
            "lastModifiedDateTime": "2025-04-29T07:13:59.2472796Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "1e9c1c38-e3fa-4f8e-8b65-d0428cc3cbae",
            "displayName": "Windows Autopatch Update Policy - Windows Autopatch - Ring1",
            "lastModifiedDateTime": "2025-04-29T07:13:57.963483Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "75c8dde6-a03f-4d2f-92f7-de2377efc8a7",
            "displayName": "Windows Autopatch Update Policy - Windows Autopatch - Ring2",
            "lastModifiedDateTime": "2025-04-29T07:13:57.9646718Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "dda956bc-9a2b-4feb-bca5-d5a410237e0c",
            "displayName": "Windows Autopatch Update Policy - Windows Autopatch - Ring3",
            "lastModifiedDateTime": "2025-04-29T07:13:58.4573694Z",
            "roleScopeTagIds": [
                "0"
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windowsUpdateForBusinessConfiguration",
            "id": "2d389cfc-60a7-4d95-8193-05e86a97ab73",
            "displayName": "Windows Autopatch Update Policy - Windows Autopatch - Test",
            "lastModifiedDateTime": "2025-04-29T07:13:57.9333973Z",
            "roleScopeTagIds": [
                "0"
            ]
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations?top=1000
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations",
    "value": [
        {
            "createdDateTime": "2020-12-03T16:34:21.4745022Z",
            "displayName": "[Base] Prod | Windows - AdministrativeTemplates | Microsoft Office 2016 ver1.0",
            "description": "Includes settings for:\nMicrosoft Office 2016 - Computer Configuration\nMicrosoft Office 2016 - User Configuration\n",
            "roleScopeTagIds": [
                "0"
            ],
            "policyConfigurationIngestionType": "builtIn",
            "id": "3d211a29-07ee-46a7-b178-80e2db565711",
            "lastModifiedDateTime": "2022-05-18T13:12:16.8230629Z"
        },
        {
            "createdDateTime": "2021-10-21T12:41:28.6466666Z",
            "displayName": "[Base] Dev | Windows - AdministrativeTemplates | AVD RDS - [Device] ver0.1 ",
            "description": "21.10.2021 - test",
            "roleScopeTagIds": [
                "0"
            ],
            "policyConfigurationIngestionType": "builtIn",
            "id": "6f9ba788-f719-46a7-b7c5-d566963d5999",
            "lastModifiedDateTime": "2022-05-18T13:09:09.6048241Z"
        },
        {
            "createdDateTime": "2020-09-22T12:48:48.6482984Z",
            "displayName": "[Base] Deprecated | Windows - AdministrativeTemplates | Microsoft Edge ver1.0",
            "description": "Includes settings for\nMicrosoft Edge - Computer Configuration\nMicrosoft Edge Update - Computer Configuration\nMicrosoft Edge - Default Settings (users can override) - Computer Configuration",
            "roleScopeTagIds": [
                "0"
            ],
            "policyConfigurationIngestionType": "builtIn",
            "id": "7666cb9a-16eb-4493-a159-c3713876d8a3",
            "lastModifiedDateTime": "2023-02-23T08:25:20.1757921Z"
        },
        {
            "createdDateTime": "2020-09-22T11:15:04.4615855Z",
            "displayName": "[Base] Prod | Windows - AdministrativeTemplates | OneDrive ver1.0",
            "description": "23.08.2021",
            "roleScopeTagIds": [
                "0"
            ],
            "policyConfigurationIngestionType": "builtIn",
            "id": "7f774f0f-2f2d-4dc3-a76f-6d45af51019e",
            "lastModifiedDateTime": "2022-05-18T13:13:47.4710049Z"
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations?$filter=microsoft.graph.androidManagedStoreAppConfiguration/appSupportsOemConfig%20eq%20true
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileAppConfigurations",
    "value": []
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/configurationPolicies?$select=id,name,description,platforms,technologies,lastModifiedDateTime,settingCount,roleScopeTagIds,isAssigned&$top=100&$filter=(platforms%20eq%20%27windows10%27%20or%20platforms%20eq%20%27macOS%27%20or%20platforms%20eq%20%27iOS%27%20or%20platforms%20eq%20%27aosp%27%20or%20platforms%20eq%20%27androidEnterprise%27)%20and%20(technologies%20eq%20%27mdm%27%20or%20technologies%20eq%20%27windows10XManagement%27%20or%20technologies%20eq%20%27appleRemoteManagement%27%20or%20technologies%20eq%20%27mdm,appleRemoteManagement%27%20or%20technologies%20eq%20%27android%27)%20and%20(templateReference/templateFamily%20eq%20%27none%27)
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies(id,name,description,platforms,technologies,lastModifiedDateTime,settingCount,roleScopeTagIds,isAssigned)",
    "@odata.count": 100,
    "@odata.nextLink": "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies?$select=id%2cname%2cdescription%2cplatforms%2ctechnologies%2clastModifiedDateTime%2csettingCount%2croleScopeTagIds%2cisAssigned&$top=100&$filter=(platforms+eq+%27windows10%27+or+platforms+eq+%27macOS%27+or+platforms+eq+%27iOS%27+or+platforms+eq+%27aosp%27+or+platforms+eq+%27androidEnterprise%27)+and+(technologies+eq+%27mdm%27+or+technologies+eq+%27windows10XManagement%27+or+technologies+eq+%27appleRemoteManagement%27+or+technologies+eq+%27mdm%2cappleRemoteManagement%27+or+technologies+eq+%27android%27)+and+(templateReference%2ftemplateFamily+eq+%27none%27)&$skiptoken=%255Bcosmosdb%255D%255B%257B%2522compositeToken%2522%253A%257B%2522token%2522%253Anull%252C%2522range%2522%253A%257B%2522min%2522%253A%2522114FC998234D2B737DB891EED734C675%2522%252C%2522max%2522%253A%252219F7AE6434F3C12D3C94DAE642CF29B0%2522%257D%257D%252C%2522resumeValues%2522%253A%255B%2522%255Bbase%255D%2520test%2520%257C%2520windows%2520-%2520settings%2520catalog%2520%257C%2520network%2520connectivity%2520status%2520indicator%2520%2520ver0.1%2522%255D%252C%2522rid%2522%253A%2522EWA-AIIdN8xRnx8AAAAACA%253D%253D%2522%252C%2522skipCount%2522%253A1%257D%255D",
    "value": [
        {
            "id": "15767346-7388-46dc-b6a5-a53c962d0cea",
            "name": "[Base] Deprecated | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 - BitLocker [Device] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2024-10-15T16:18:43.6053582Z",
            "settingCount": 5,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "3e7973ab-c518-4daf-91e6-43bc995187e8",
            "name": "[Base] Deprecated | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 - Computer [Device] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog\nExperience Settings to moved to deadicated experience config profile. due to user/device level config mismatch",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2023-02-23T08:20:50.0120577Z",
            "settingCount": 122,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "bf4f86fa-9cd8-4638-ae1b-a27de0cbb6b7",
            "name": "[Base] Deprecated | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 - Defender Antivirus [Device] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2023-02-23T08:42:56.8765333Z",
            "settingCount": 8,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "5edf7e62-7ff6-4820-9d81-db7d6a8219ca",
            "name": "[Base] Dev | Windows - Settings Catalog | Control Panel > Personalization [user] ver0.1",
            "description": "context: user level\n15.06.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-15T13:10:04.4010906Z",
            "settingCount": 24,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "c8eb03e1-3160-4885-888b-4c8b2854817e",
            "name": "[Base] Dev | Windows - Settings Catalog | Credential Providers ver0.1",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:11:14.094408Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "d4b3906c-05e3-42c4-845e-7528d697808c",
            "name": "[Base] Dev | Windows - Settings Catalog | Default search provider [User]- ver0.1",
            "description": "",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-04T18:58:08.2483649Z",
            "settingCount": 9,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "17436f8b-a93c-45d6-a204-6a80d3d43155",
            "name": "[Base] Dev | Windows - Settings Catalog | Delivery Optimization ver0.1",
            "description": "30.05.2022\nContext: Device\n\nConfig half done",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-30T15:13:21.0533046Z",
            "settingCount": 28,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "47a4d644-f80b-4378-affb-699e92033d85",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Edge - Identity and sign-in [User]- ver0.1",
            "description": "04.10.2022\nContext: User",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-04T20:28:51.8534756Z",
            "settingCount": 5,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "df287451-37e4-4e45-8cc8-6aa35638cb51",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Edge - Update ver0.1",
            "description": "Update configuration for microsoft Edge standard version. Not Dev or Canary\n15.06.2022\nContext: Device",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-04T16:28:21.4949528Z",
            "settingCount": 14,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "8c2be96d-4bbc-4735-a00f-8e8114fc6a6a",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Edge [User]- ver0.1",
            "description": "15.06.2022\nContext: User Level",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-04T18:57:56.0888531Z",
            "settingCount": 266,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "ff2ba044-11f9-4d06-b0b4-6a21e9033f8a",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Edge Web View2 [User] - ver0.1",
            "description": "",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T08:03:21.9396047Z",
            "settingCount": 3,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "ee3e8eae-593f-4c20-b620-5d7c99accc32",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Edge Version 107 [Device] ver1.0",
            "description": "Scope: Device\n\nImported GPReport from - https://learn.microsoft.com/en-us/windows/security/threat-protection/windows-security-configuration-framework/security-compliance-toolkit-10",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2023-01-27T14:58:42.9329378Z",
            "settingCount": 21,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "ab933397-e557-42ea-a41d-f6127dc81d9d",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Edge Version 98 [Device] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nMicrosoft Security Baseline for MSFT Edge Version 98 based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-04T20:35:33.170308Z",
            "settingCount": 20,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "557b0ef1-b2ae-4dcb-a2aa-16ca33708827",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Internet Explorer 11 [Device] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T05:12:30.86793Z",
            "settingCount": 116,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "728bd1cc-d710-4078-845a-dfa084a6a494",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Internet Explorer 11 [User] ver1.0",
            "description": "19.09.2022\nContext: User Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T05:12:08.7938513Z",
            "settingCount": 117,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "7c6d51e8-fc24-47be-bbba-59bb36442660",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Microsoft 365 Apps v2206 - DDE Block [User] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nMicrosoft 365 Apps Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T05:20:29.1954704Z",
            "settingCount": 18,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "43f92fb2-c476-4006-98d1-7e3782f82af4",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Microsoft 365 Apps v2206 - Require Macro Signing [User] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nMicrosoft 365 Apps Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T05:33:29.0514525Z",
            "settingCount": 8,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "440687a0-589c-4994-9d45-308cdbaab052",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Microsoft 365 Apps v2206 [User] ver1.0",
            "description": "19.09.2022\nContext: User Level\nMicrosoft 365 Apps Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T05:35:50.4675787Z",
            "settingCount": 94,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "053966ad-9cc6-4fca-9ee1-aecec1d7ff84",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 - Domain Security [Device] ver1.0",
            "description": "19.09.2022\nContext: Device Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog\n\nDevice settings have been switched off as it breaks device compliance rules, due to the admin account also requiring to be password changed",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T13:23:03.3563587Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "f3fdd8e1-05aa-4b9f-ae41-44426d935161",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 22H2 - Computer [Device] ver1.0",
            "description": "23-02-2023\nContext: Device Level\n\nConfiguration from https://learn.microsoft.com/en-us/windows/security/threat-protection/windows-security-configuration-framework/security-compliance-toolkit-10\n\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\n\nGPO migration to Settings Catalog\nExperience Settings to moved to deadicated experience config profile. due to user/device level config mismatch",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2023-02-23T08:22:12.4493985Z",
            "settingCount": 123,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "077c120f-0f87-4efc-857d-d44a9ff61d92",
            "name": "[Base] Dev | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 [User] ver1.0",
            "description": "19.09.2022\nContext: User Level\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T06:06:19.0067453Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "a4777a1e-516b-45a4-b808-46654e826e30",
            "name": "[Base] Dev | Windows - Settings Catalog | Start Menu ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:27:18.5728184Z",
            "settingCount": 27,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "8c7c02ce-6f43-47a8-9e8d-080216dcd7a1",
            "name": "[Base] Dev | Windows - Settings Catalog | Text Input (Japanese + Chinese) ver1.0",
            "description": "22.08.2021",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:30:56.8849564Z",
            "settingCount": 26,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "6440fff5-4ed0-4cac-a9b3-37aa8f4d93b8",
            "name": "[Base] NotInUse | Windows - CIS-Framework | Network - Network Provider ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004\n\nUsing Microsoft Secuirty Baseline for windows 11",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T07:33:25.9611192Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "277d30af-c598-434d-8566-6d9609eb1b78",
            "name": "[Base] Physical | Windows - Settings Catalog | WiFi Settings",
            "description": "19.09.2022\nContext Device\nCIS Baseline - windows 10 2004 + DT Config",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T11:01:19.9412753Z",
            "settingCount": 6,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "5b0c3c05-d7e7-48e3-9e77-174d8ab095be",
            "name": "[Base] Prod | Windows - Settings Catalog  | System ver1.0",
            "description": "Migrated to settings catalogue from Administrative Templates\n18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-19T15:15:24.9370047Z",
            "settingCount": 26,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "8d1ff18c-92f6-4a7e-8cb2-be69418a7a34",
            "name": "[Base] Prod | Windows - Settings Catalog | Above Lock ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:30:54.4682726Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "2bb5f88b-98f5-4fe5-8f49-ee678e817ee8",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Control Panel - Personalization ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:22:32.2625237Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "6aa986c3-9f73-4156-a95f-8a4706691f31",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Microsoft Defender Smartscreen ver1.0",
            "description": "17.09.2022\nContext Device\nCIS Baseline - windows 10 2004 + DT config",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T11:49:28.2753416Z",
            "settingCount": 11,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "97bcc98d-0638-4e28-8301-4b924b9fbb35",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - MSS (Legacy) ver1.2",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004\nAdded \"Enable Safe DLL search mode (recommended)\"\n\nAdded \"MSS: (DisableIPSourceRouting) IP source routing protection level (protects against packet spoofing) to Enabled\" for ipv4 and ipv6 as part of MSFT security baseline. All other settings have not been configured\n\nAdded \"MSS: (NoNameReleaseOnDemand) Allow the computer to ignore NetBIOS name release\" to follow CIS baseline 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:33:18.698355Z",
            "settingCount": 22,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "3e7a79a4-1f72-48da-ac11-5ebef736a553",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Network - Network Connections ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:43:32.3279414Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "1809ca98-908f-4de0-b0bf-b9c4b54b60ae",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Network - Windows Connection Manager ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:43:15.6654964Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "77f9aa2c-4863-4814-bf28-81b15ba08f54",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Start Menu and Taskbar - Notifications ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:42:35.1872707Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "fa6abde5-dbaa-47e4-95da-f6290d734978",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Credentials Delegation ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:42:18.2298647Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "0cf16d74-2a8d-4830-ab02-edbe9f20afa9",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Device Installation Restrictions ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:44:19.2910697Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "da43104b-4a3d-4e51-8577-1f53f824d039",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Early Launch Antimalware ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:42:05.0360268Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "44235e3a-cc09-4b6f-910e-09899bd86f3e",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Internet Communication settings ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:41:49.6497819Z",
            "settingCount": 4,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "7b2e2a93-914d-43cb-836f-365e408070c0",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Logon ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:29:06.7980796Z",
            "settingCount": 4,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "955a700c-66a3-42c3-9c00-5e2bd9d06253",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Remote Assistance ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:41:32.413726Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "aacbe987-672e-4a96-afd1-3bc4c068b3ea",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Remote Procedure Call (RPC) ver1.0",
            "description": "13.09.2022\nContext Device\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:41:14.9268521Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "6959eac3-6d12-4f01-9e2a-3a371f3e2015",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Windows Components - App Privacy ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:24:59.811485Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "b1c905e7-c5d9-4bdf-82a0-f1595f47a74a",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Windows Components - App Runtime ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:24:45.6034057Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "2610f89e-0fc3-440c-a8bb-42a711ec1341",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Windows Components - Attachment Manager ver1.0",
            "description": "14.09.2022\nProfile Context: User Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:24:31.5703992Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "8e9578b2-bd58-4d3b-b925-de6c17951756",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Windows Components - AutoPlay Policies ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:24:11.5942922Z",
            "settingCount": 3,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "3de41308-4460-4008-b2e8-5de1d6b1a148",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Windows Components - Credential User Interface ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:23:52.5617042Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "5682a398-71ed-415f-89a1-29faadd5a4e3",
            "name": "[Base] Prod | Windows - Settings Catalog | Admin Templates - Windows Components - File Explorer ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:23:30.3189035Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "89053047-a324-4322-bf82-384cb9012bd1",
            "name": "[Base] Prod | Windows - Settings Catalog | Application Defaults ver1.0",
            "description": "Guide: https://techcommunity.microsoft.com/t5/ask-the-performance-team/how-to-configure-file-associations-for-it-pros/ba-p/1313151\n\nbase64Encode: https://www.base64encode.org/\n",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T14:58:55.9065312Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "6aa24a8e-878f-4046-961f-b2f9e99c3a3f",
            "name": "[Base] Prod | Windows - Settings Catalog | Authentication ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T14:59:22.530491Z",
            "settingCount": 7,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "c6b6f00b-cb55-4be3-b3f1-a55f800db0c0",
            "name": "[Base] Prod | Windows - Settings Catalog | Bitlocker ver1.0",
            "description": "23.05.2022\n\nSets:\nThe encryption method for fixed data drives: AES-CBC 256-bit\nThe encryption method for operating system drives: XTS-AES 256-bit\nThe encryption method for removable data drives: XTS-AES 256-bit",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T12:08:10.6041999Z",
            "settingCount": 16,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "718a0086-2b94-41ab-bd4e-1f2d2d855a28",
            "name": "[Base] Prod | Windows - Settings Catalog | BITS ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:01:05.2677736Z",
            "settingCount": 6,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "b1e17084-c32c-4b15-9124-3422809d2db5",
            "name": "[Base] Prod | Windows - Settings Catalog | Camera ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-23T21:58:10.5475773Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "fc30a96e-a0aa-4fd9-91a7-3c14f6be687d",
            "name": "[Base] Prod | Windows - Settings Catalog | Control Policy Conflict (MDM Wins Over GP) ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:10:29.120237Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "d18e84d5-a4fa-41bc-ba1c-e21c88cd908f",
            "name": "[Base] Prod | Windows - Settings Catalog | Cryptography ver1.0",
            "description": "\nRef: https://docs.microsoft.com/en-us/answers/questions/166697/disable-weak-cipher-suits-with-windows-server-2016.html\n\nDisabled:\nTLS_DHE_RSA_WITH_AES_256_CBC_SHA\nTLS_DHE_RSA_WITH_AES_128_CBC_SHA\nTLS_RSA_WITH_AES_256_GCM_SHA384\nTLS_RSA_WITH_AES_128_GCM_SHA256\nTLS_RSA_WITH_AES_256_CBC_SHA256\nTLS_RSA_WITH_AES_128_CBC_SHA256\nTLS_RSA_WITH_AES_256_CBC_SHA\nTLS_RSA_WITH_AES_128_CBC_SHA\nTLS_RSA_WITH_3DES_EDE_CBC_SHA\nTLS_DHE_DSS_WITH_AES_256_CBC_SHA256\nTLS_DHE_DSS_WITH_AES_128_CBC_SHA256\nTLS_DHE_DSS_WITH_AES_256_CBC_SHA\nTLS_DHE_DSS_WITH_AES_128_CBC_SHA\nTLS_DHE_DSS_WITH_3DES_EDE_CBC_SHA\nTLS_RSA_WITH_RC4_128_SHA\nTLS_RSA_WITH_RC4_128_MD5\nTLS_RSA_WITH_NULL_SHA256\nTLS_RSA_WITH_NULL_SHA\nTLS_PSK_WITH_AES_256_GCM_SHA384\nTLS_PSK_WITH_AES_128_GCM_SHA256\nTLS_PSK_WITH_AES_256_CBC_SHA384\nTLS_PSK_WITH_AES_128_CBC_SHA256\nTLS_PSK_WITH_NULL_SHA384\nTLS_PSK_WITH_NULL_SHA256",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:14:51.2037598Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "133a0505-b6a1-4a08-812d-64d422567092",
            "name": "[Base] Prod | Windows - Settings Catalog | Device Lock ver1.0",
            "description": "Only set Max Inactivity Time Device Lock. All other setting break password policy and therefore device compliance.\n\nhttps://smsagent.blog/2021/04/27/a-case-of-the-unexplained-intune-password-policy-and-forced-local-account-password-changes/\n\nWhen you enable a password policy some default values will get set for password length and complexity, and these polices will require that a local administrator account change its password at next logon. The reason is that the policy doesn’t know if the currently set password meets the requirements of the policy. The only way it can be sure it complies is to force you to change it, and the new password must meet the policy requirements. In this way it can truthfully report whether the device is compliant to the policy.\n\nSounds reasonable, but forcing the password change breaks your LAPS solution as you cannot programmatically change the password, or even remove the requirement to change the password at next logon",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-04T20:53:32.1921788Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "48165c84-3891-4aaf-9f1a-b0232acd43ca",
            "name": "[Base] Prod | Windows - Settings Catalog | Event Log Service ver1.1",
            "description": "23.05.2022\nContext: Device Level\nConfigured to match CIS baseline 2004\n\nSets: \nApplication Logging\nSecurity Logging\nSetup Logging\nSystem Logging",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-17T06:11:57.6247186Z",
            "settingCount": 21,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "d6c524a0-6dea-4188-a0c0-8cfc18eb12ea",
            "name": "[Base] Prod | Windows - Settings Catalog | Experience ver1.0",
            "description": "18.05.2022\nContext: User Level\nMatches MSFT Security Baseline",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T08:36:39.9204201Z",
            "settingCount": 12,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "8701af03-5b3f-4a0b-bbf0-287ec76f4f41",
            "name": "[Base] Prod | Windows - Settings Catalog | FIDO Security Keys Enablement for Windows Sign-In [User] ver1.0",
            "description": "21.05.2022\n\nEnables targetted FIDO Security Keys to be used during Windows Sign In. Virtual devices are included since FIDO 2 keys present as keyboards to vm's.\n\nRef: https://docs.microsoft.com/en-us/azure/active-directory/authentication/howto-authentication-passwordless-security-key-windows#targeted-intune-deployment",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-17T07:22:55.9584233Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "5aaf0ede-c20f-4eb1-8cdf-5421b09338bf",
            "name": "[Base] Prod | Windows - Settings Catalog | Games ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:21:02.6746341Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "9af36d64-2294-4c46-8f78-f30bf4d17061",
            "name": "[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0",
            "description": "30.05.2022\nContext: Device\n\nThis parameter adds a list of domains that an Azure Active Directory joined device should attempt to contact if it is otherwise unable to resolve a UPN to a principal.",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-04T18:54:51.4799413Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "4a126db0-997f-4156-98ed-07e4449730b8",
            "name": "[Base] Prod | Windows - Settings Catalog | Licensing ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:22:58.71261Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "819d4e88-5214-45a8-a44a-e034b173ef70",
            "name": "[Base] Prod | Windows - Settings Catalog | Local Policies - Security Options ver1.1.3",
            "description": "04.11.2021\nver1.1.3\nContext: Device\n\nInteractive Logon Machine Inactivity Limit set to 900. Setting is in seconds. Had on 15 thinking minutes and it was painful :)\n\nUpdated LAN Manager authentication level to \"Send NTLMv2 response only. Refuse LM & NTLM\"\n\nAdded the inactivity time (in seconds) of a logon session from the device lock config to LPSO as the screen saver part is housed here\n",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-12T17:34:59.7314052Z",
            "settingCount": 41,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "6ca8f216-471b-478b-bb14-62b92fe33c54",
            "name": "[Base] Prod | Windows - Settings Catalog | Lock Screen User Experience ver1.1",
            "description": "05.06.2022\nver1.0\nContext: User\n\nConfigures lock screen behaviour for the End User\n\nPolicy Updated to support windows CIS baseline 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-14T11:36:56.4532818Z",
            "settingCount": 9,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "25469cfa-5303-4cff-9b61-1bfa6e029ed9",
            "name": "[Base] Prod | Windows - Settings Catalog | Memory Dump ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:24:50.3991409Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "9750f3c7-41a7-441b-86a2-b9d71cdf56ea",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft App Store & App Install Behaviour ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-18T05:29:55.7644498Z",
            "settingCount": 13,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "9e326b81-0094-4c27-8155-a9532ac7db68",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Application Guard for Office 365 ver1.0",
            "description": "30.05.2022\ncontext: User\n\nApplies to: Word, Excel, and PowerPoint for Microsoft 365, Windows 10 Enterprise, Windows 11 Enterprise\n\nMicrosoft Defender Application Guard for Office (Application Guard for Office) helps prevent untrusted files from accessing trusted resources, keeping your enterprise safe from new and emerging attacks.\n\nhttps://docs.microsoft.com/en-us/microsoft-365/security/office-365-security/install-app-guard?view=o365-worldwide\n\nTo reset (clean up) a container and clear persistent data inside the container:\n\nhttps://techcommunity.microsoft.com/t5/core-infrastructure-and-security/windows-10-all-things-about-application-guard/ba-p/2455596",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-04T18:54:24.5160236Z",
            "settingCount": 8,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "7b19810a-97ac-48dc-a17b-e9930d267463",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Attack Surface Reduction ver1.1",
            "description": "23.05.2022\nver 1.1\ncontext: device\n\nSet \"Block process creations originating from PSExec and WMI commands\" to audit, rather than block. As WmiPrvSE.exe isn't allowed to interact with C:\\Windows\\System32\\Dism.exe on W365 devices.",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-02T17:49:51.7451662Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "49d932ba-b537-4d19-b5d6-ffd79c1db7d0",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Controlled Folder Access  - Base Part 1 ver1.1.8",
            "description": "ver1.1.8\n23.10.2024\nContext: Device\nRansomware Protection Controls\n\nThe following controlled folder access events appear in Windows Event Viewer under Application and Service Logs -> Microsoft -> Windows -> Windows Defender -> Operational. <- use to find path for white listing\n\n\nhttps://github.com/jdgregson/Windows-10-Exploit-Protection-Settings\n\nBaseline - https://gist.github.com/ag-michael/f90751782090f8a92ce6ccc3629bccfc with edits and additions from testing.\n\nexclusions : https://docs.microsoft.com/en-us/microsoft-365/security/defender-endpoint/configure-extension-file-exclusions-microsoft-defender-antivirus?view=o365-worldwide",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2024-10-23T10:02:50.8259348Z",
            "settingCount": 3,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "b78de4cb-da8e-4289-b047-385649904a86",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Controlled Folder Access - Base Part 2 ver1.0.5",
            "description": "24.12.2022\nContext: all devices\nPart 2 continues the rules due to the limts of the number of enteries being reached for part 1",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2024-10-23T09:57:24.6258112Z",
            "settingCount": 3,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "ef7c53b9-324d-4fdc-83be-a44fc8823772",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Controlled Folder Access - Developer ver1.0.1",
            "description": "context: device\n23.10.2024\n\nA set of Defender controlled folder access rules for developer tools.\nScoped to all devices, but out of lab should be scoped to a targetted scope.\n\nApplication and Service Logs -> Microsoft -> Windows -> Windows Defender -> Operational",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2024-10-23T10:36:48.6011653Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": false
        },
        {
            "id": "10340be1-7edf-4d64-893d-96197e0c04af",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Device Guard ver1.0",
            "description": "13.09.2022\nContext: Device Level\nCIS Baseline - windows 10 2004\n\nIncludes Credential Guard",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T11:05:18.818068Z",
            "settingCount": 4,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "2eaa9a65-ab7a-4099-a29d-8305c4aecfb6",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender DMA Guard ver1.0",
            "description": "21.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-21T10:35:41.3324491Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "dc45c07d-458f-480d-b2c4-65e09ca11493",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender For Endpoint [AV] ver1.0",
            "description": "19.05.2022\n\nDescription\n\nExclusions Ref: https://support.microsoft.com/en-us/help/822158/virus-scanning-recommendations-for-enterprise-computers\n\nPlus exclusions from https://github.com/ukncsc/Device-Security-Guidance-Configuration-Packs",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T06:26:22.6703478Z",
            "settingCount": 37,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "bbebee95-df04-47e5-b163-5ccad2af34e8",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender For Endpoint Reporting ver1.0",
            "description": "23.05.2022\nContext: Device",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-30T15:51:46.427988Z",
            "settingCount": 8,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "bd189ee5-6d6c-4c54-be39-2931d06fe191",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Security Center Access ver1.1",
            "description": "23.05.2022\nContext: Device Level\n\nUpdated to match CIS baseline 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-17T06:07:28.527702Z",
            "settingCount": 22,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "a021f1db-b653-41ce-aac2-de0c5f95b9b8",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Defender Security Intelligence Updates ver1.0",
            "description": "Microsoft Defender Antivirus uses cloud-delivered protection (also called the Microsoft Advanced Protection Service or MAPS) and periodically downloads dynamic security intelligence updates to provide additional protection. These dynamic updates don't take the place of regular security intelligence updates via security intelligence update KB2267602.\n\nApplying to user level as device level causes 65000 errors on unknown accounts",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-24T06:54:21.6739785Z",
            "settingCount": 10,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "969110d2-876d-4b39-9ef1-fc6faf8c4525",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Edge - Default Settings (users can override) [User]- ver1.1",
            "description": "context: user\n15.06.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-24T07:15:18.8340627Z",
            "settingCount": 82,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "27a390f8-45a7-45a2-ba8f-1955321fc051",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Edge - Extensions [User]- ver1.0",
            "description": "19.09.2022\nContext: User",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-24T07:10:37.565775Z",
            "settingCount": 7,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "8106a98d-3478-42bc-8e5f-c41cd803910f",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 22H2 - BitLocker [Device] ver1.0",
            "description": "23-02-2023\nContext: Device Level\n\nConfiguration from https://learn.microsoft.com/en-us/windows/security/threat-protection/windows-security-configuration-framework/security-compliance-toolkit-10\n\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2023-02-23T08:36:50.6617072Z",
            "settingCount": 5,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "f04404aa-1809-4ede-8e44-2d6746538a8d",
            "name": "[Base] Prod | Windows - Settings Catalog | Microsoft Security Baseline | MSFT Windows 11 22H2 - Defender Antivirus [Device] ver1.0",
            "description": "23.02.2023\nContext: Device Level\nhttps://learn.microsoft.com/en-us/windows/security/threat-protection/windows-security-configuration-framework/security-compliance-toolkit-10\n\nWindows 11 Security baseline based upon the Microsoft Security Compliance Toolkit 1.0 - Ref: https://www.microsoft.com/en-us/download/details.aspx?id=55319\nGPO migration to Settings Catalog",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2023-02-23T08:43:56.636764Z",
            "settingCount": 9,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "7b4b9a36-6542-4857-94ae-9ae53dea4ee0",
            "name": "[Base] Prod | Windows - Settings Catalog | MS Security Guide ver1.1",
            "description": "04.06.2022\nContext: Device\n\nAdded \"Configure SMB v1 client driver\" to follow CIS baseline 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-14T11:23:39.7792349Z",
            "settingCount": 5,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "d0546eda-adcb-4c42-b4c7-1ad6c872403c",
            "name": "[Base] Prod | Windows - Settings Catalog | Network Connections ver1.0",
            "description": "04.06.2022\nContext: User",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-05T05:54:29.8918595Z",
            "settingCount": 28,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "2fcb17fa-d834-49ac-ae97-03353e245505",
            "name": "[Base] Prod | Windows - Settings Catalog | OneDrive ver1.0",
            "description": "19.09.2022\nContext: User Level",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T07:03:53.2164107Z",
            "settingCount": 40,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "a0d2ea67-15f2-4f5f-a3d1-8ddf8bab369f",
            "name": "[Base] Prod | Windows - Settings Catalog | Remote Desktop AVD URL ver1.0",
            "description": "23.05.2022\nContext: User Level\n\nSets the AVD URL within the MSFT Remote Desktop App\n\nhttps://rozemuller.com/configure-subscribe-to-avd-in-rdp-client-automated/",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-30T14:54:32.5957451Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "00dfc9cf-70ef-4a1c-b5a0-751fad397edd",
            "name": "[Base] Prod | Windows - Settings Catalog | Remote Desktop Services - Session Time Limits ver1.0",
            "description": "23.09.2022\nContext: Device\n\nSets maximum session time and sets time limits",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-23T07:29:13.4030752Z",
            "settingCount": 4,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "72065043-f74a-4c69-a0dd-80129aee8917",
            "name": "[Base] Prod | Windows - Settings Catalog | Remote Desktop Services ~ Security ver1.0",
            "description": "04.06.2022\nContext: Device\n\nSets the security configuration for RDP connections on all device types\n\nUsed in conjuction with:\n\n[Physical] Prod | Windows - Settings Catalog | Remote Desktop Services ~ Jumphost/User device Connections ver1.0",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-08-12T15:06:58.7243325Z",
            "settingCount": 6,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "d3aa73b6-e3f6-4063-8eec-4e1e98b86ffc",
            "name": "[Base] Prod | Windows - Settings Catalog | Search ver1.0",
            "description": "13.09.2022\nContext Device\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T11:08:23.8605363Z",
            "settingCount": 14,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "3e2ed8d1-5d39-449c-9f03-ff527a4da209",
            "name": "[Base] Prod | Windows - Settings Catalog | Storage Sense ver1.0",
            "description": "Configures Storage Sense of windows devices\n\nTurn on Storage Sense\nConfigure Storage Sense to run on a weekly schedule\nDelete temporary files\nDelete files in recycle bin that has been present for over 30 days\nDelete files in my Downloads folder if present for over 90 days\nRevert OneDrive synced files to revert to cloud-only if not accessed in over 30 days\n\nRef: https://letsconfigmgr.com/storage-sense-settings-catalog-msintune/",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:34:38.4561578Z",
            "settingCount": 8,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "c06fcc51-7396-4a02-91f1-2592489724e6",
            "name": "[Base] Prod | Windows - Settings Catalog | System Services - Xbox ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:40:47.2331158Z",
            "settingCount": 4,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "e24ccf91-75f8-446e-980a-0b3e9af02fc1",
            "name": "[Base] Prod | Windows - Settings Catalog | Task Manager ver1.0",
            "description": "22.08.2021",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:27:51.3132291Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "8a66ebfa-ce5c-4a01-9f7b-66928ceecf67",
            "name": "[Base] Prod | Windows - Settings Catalog | Troubleshooting ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:32:05.1002131Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "056bedd4-6a3e-4e0f-ad88-4c31be2032f0",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows Components - RSS Feeds ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:46:55.2668641Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "a0d6d692-fb4f-4729-9c96-4d2ec14ccede",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows Components - Windows Logon Options ver1.0",
            "description": "14.09.2022\nProfile Context: User Level\nCIS Baseline - windows 10 2004",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T10:13:19.9749975Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "9a8d6572-6f44-4b0c-aaf0-026c41258637",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows Connection Manager ver1.0",
            "description": "17.09.2021\n\nThis policy setting prevents computers from connecting to both a domain-based network and a non-domain-based network at the same time.\n\nIf this policy setting is enabled, the computer responds to automatic and manual network connection attempts based on the following circumstances:\n\nTesting with user based deployment as device based causes lots of 65000 errors",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-23T21:50:14.01939Z",
            "settingCount": 4,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "0838f743-8533-4948-af44-090fc18a7d2b",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows Installer ver1.0",
            "description": "05.06.2022\nver0.1\nContext: Device\n",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-05T05:52:37.4330307Z",
            "settingCount": 24,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "f27f0b73-a9ef-45d0-895d-03c561144c0d",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows Logon ver1.0",
            "description": "18.05.2022",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:31:33.1738186Z",
            "settingCount": 2,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "b1dbc936-3731-4619-9d5a-f91f53333cca",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows Remote Management (WinRM) ver1.1",
            "description": "05.06.2022\nver1.0\nContext: Device\nStarting to migrate settings across from admin template\nUpdated to match windows CIS baseline 2004\n\nAllow Basic authentication and Turn On Compatibility HTTP Listener both fail on W365 devices. all other settings apply.",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T06:59:36.8189206Z",
            "settingCount": 17,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "fd0d702b-48fc-4b45-adb7-d41cfa8c097f",
            "name": "[Base] Prod | Windows - Settings Catalog | Windows SSPR AAD Password Reset From Login Screen ver1.0",
            "description": "18.10.2021\n\nconfigures a Windows device for SSPR at the sign-in screen\n\nRef: https://docs.microsoft.com/en-us/azure/active-directory/authentication/howto-sspr-windows",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-05-18T15:36:13.5389154Z",
            "settingCount": 1,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "d06f9e56-82c2-4975-a654-3e2bdd971ae1",
            "name": "[Base] Test | Windows - Settings Catalog | Microsoft Defender Application Guard [User] ver0.1",
            "description": "30.05.2022\nContext: User\n\nTo not apply to an azure VDI unless the underlying SKU can support nested virtualization. e.g W365 will not work. AVD can work so long as the correct SKU is selected.",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-09-19T06:55:33.1366552Z",
            "settingCount": 6,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "4b196a09-803f-48dc-a877-6279579f0fc5",
            "name": "[Base] Test | Windows - Settings Catalog | Microsoft Defender Firewall : State ver1.0",
            "description": "14.09.2022\nProfile Context: Device Level\nCIS Baseline - windows 10 2004\n\nNeed to migrate this to another profile",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-10-12T11:43:33.3699119Z",
            "settingCount": 3,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        },
        {
            "id": "a768c8eb-3b81-4e34-bfb5-ba34a338e1b8",
            "name": "[Base] test | Windows - Settings Catalog | Network Connectivity Status Indicator  ver0.1",
            "description": "04.06.2022\ncontext: device\n\nUsed for determining the network state for a given windows device. Where line of sight to a beacon infers that a device is on prem, whereas the lack thereof denotes that the device is off site. The status detected by this policy is leveraged by other applications such as the local windows firewall to apply different firewall policies based on network state.",
            "platforms": "windows10",
            "technologies": "mdm",
            "lastModifiedDateTime": "2022-06-04T10:07:56.0093305Z",
            "settingCount": 7,
            "roleScopeTagIds": [
                "0"
            ],
            "isAssigned": true
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies?$select=id,displayName,lastModifiedDateTime,roleScopeTagIds,microsoft.graph.androidCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.androidWorkProfileCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.iosCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.windows10CompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.iosCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,microsoft.graph.androidWorkProfileCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,microsoft.graph.androidDeviceOwnerCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,microsoft.graph.androidDeviceOwnerCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.androidCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel&$expand=assignments&top=1000
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies(id,displayName,lastModifiedDateTime,roleScopeTagIds,microsoft.graph.androidCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.androidWorkProfileCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.iosCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.windows10CompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.iosCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,microsoft.graph.androidWorkProfileCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,microsoft.graph.androidDeviceOwnerCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,microsoft.graph.androidDeviceOwnerCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel,microsoft.graph.androidCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel,assignments())",
    "value": [
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "04ab1704-76f9-4306-b9a6-4dd5056bd466",
            "displayName": "Windows - Device Security | Applies to: [Virtual Corporate Device]\u00A0ver1.0",
            "lastModifiedDateTime": "2022-07-11T09:56:06.170645Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('04ab1704-76f9-4306-b9a6-4dd5056bd466')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "04ab1704-76f9-4306-b9a6-4dd5056bd466_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "04ab1704-76f9-4306-b9a6-4dd5056bd466",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "2983b1c2-8ec2-45d3-84ed-deca619d2c04",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "110b7be5-1ca8-4504-adb1-017f5812c18c",
            "displayName": "iOS/iPadOS - Password | [MAM-WE] | Standard Users",
            "lastModifiedDateTime": "2025-07-23T12:30:03.3266083Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('110b7be5-1ca8-4504-adb1-017f5812c18c')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "110b7be5-1ca8-4504-adb1-017f5812c18c_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "110b7be5-1ca8-4504-adb1-017f5812c18c",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "28b767ca-654c-4605-9371-f1ea044f4207",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.androidDeviceOwnerCompliancePolicy",
            "id": "1e42d8de-62e2-4ba1-927e-c0716f744f8a",
            "displayName": "Android Enterprise - Defender Risk Score | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:39:49.0009789Z",
            "roleScopeTagIds": [
                "0"
            ],
            "advancedThreatProtectionRequiredSecurityLevel": "low",
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('1e42d8de-62e2-4ba1-927e-c0716f744f8a')/microsoft.graph.androidDeviceOwnerCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "1e42d8de-62e2-4ba1-927e-c0716f744f8a_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "1e42d8de-62e2-4ba1-927e-c0716f744f8a",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": null,
                        "deviceAndAppManagementAssignmentFilterType": "none",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "229bab1c-dfde-4433-9e57-02c081e8f2c7",
            "displayName": "Windows - Device Properties | Applies to: [Physical Autopilot Device] & [Corporate Owned Specialty Devices]\u00A0ver1.0",
            "lastModifiedDateTime": "2022-07-11T09:54:54.2870745Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('229bab1c-dfde-4433-9e57-02c081e8f2c7')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "229bab1c-dfde-4433-9e57-02c081e8f2c7_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "229bab1c-dfde-4433-9e57-02c081e8f2c7",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665",
                        "deviceAndAppManagementAssignmentFilterType": "exclude",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "3aa6e3d1-76bf-4ee4-8a31-aad8977f3538",
            "displayName": "Windows - Device Properties | Applies to: [Corporate Owned Specialty Devices] ver1.0",
            "lastModifiedDateTime": "2022-07-11T09:57:06.137178Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('3aa6e3d1-76bf-4ee4-8a31-aad8977f3538')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "3aa6e3d1-76bf-4ee4-8a31-aad8977f3538_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "3aa6e3d1-76bf-4ee4-8a31-aad8977f3538",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "465712ec-13fe-4c99-a918-0496461ef884",
            "displayName": "iOS/iPadOS - Device Properties | [ADE] [BYOD] [MAM-WE] | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:46:04.977932Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('465712ec-13fe-4c99-a918-0496461ef884')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "465712ec-13fe-4c99-a918-0496461ef884_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "465712ec-13fe-4c99-a918-0496461ef884",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "6df418a0-c390-41ba-9fe4-927eb56f6ce7",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.androidCompliancePolicy",
            "id": "4d964631-ff9f-4da2-b59d-162c01ee1e09",
            "displayName": "Android Device Administrator - Full Compliance Policy [Specialty Devices] | MS Teams",
            "lastModifiedDateTime": "2022-07-11T10:00:54.2155028Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "secured",
            "advancedThreatProtectionRequiredSecurityLevel": "secured",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('4d964631-ff9f-4da2-b59d-162c01ee1e09')/microsoft.graph.androidCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "4d964631-ff9f-4da2-b59d-162c01ee1e09_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "4d964631-ff9f-4da2-b59d-162c01ee1e09",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": null,
                        "deviceAndAppManagementAssignmentFilterType": "none",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "500f7039-1ecd-4b2c-bd33-ec12310050f3",
            "displayName": "Check-AvEnabled",
            "lastModifiedDateTime": "2024-01-02T13:19:57.9768596Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('500f7039-1ecd-4b2c-bd33-ec12310050f3')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "500f7039-1ecd-4b2c-bd33-ec12310050f3_ea8e2fb8-e909-44e6-bae7-56757cf6f347",
                    "source": "direct",
                    "sourceId": "500f7039-1ecd-4b2c-bd33-ec12310050f3",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": null,
                        "deviceAndAppManagementAssignmentFilterType": "none",
                        "groupId": "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "51104e06-f84e-4ca4-9adb-d7691ee1a1f8",
            "displayName": "Windows - Custom Compliance | Verify installed applications | [Corporate Owned Specialty Devices] ver1.0",
            "lastModifiedDateTime": "2022-07-11T09:59:52.4602017Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('51104e06-f84e-4ca4-9adb-d7691ee1a1f8')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": []
        },
        {
            "@odata.type": "#microsoft.graph.androidDeviceOwnerCompliancePolicy",
            "id": "52274b00-3964-4ad3-a8c6-a3931fe145df",
            "displayName": "Android Enterprise - Device Health | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:39:00.6168696Z",
            "roleScopeTagIds": [
                "0"
            ],
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "deviceThreatProtectionRequiredSecurityLevel": "low",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('52274b00-3964-4ad3-a8c6-a3931fe145df')/microsoft.graph.androidDeviceOwnerCompliancePolicy/assignments",
            "assignments": []
        },
        {
            "@odata.type": "#microsoft.graph.androidDeviceOwnerCompliancePolicy",
            "id": "5eea7ff0-14b3-487e-a579-76f9c3216bc6",
            "displayName": "Android Enterprise - Device Properties | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:41:43.7979961Z",
            "roleScopeTagIds": [
                "0"
            ],
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('5eea7ff0-14b3-487e-a579-76f9c3216bc6')/microsoft.graph.androidDeviceOwnerCompliancePolicy/assignments",
            "assignments": []
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "5f37000d-8479-43d2-882b-e0ccfcf8bac6",
            "displayName": "iOS/iPadOS - Defender Risk Score | [ADE] [BYOD] [MAM-WE] | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:44:35.7956145Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "advancedThreatProtectionRequiredSecurityLevel": "low",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('5f37000d-8479-43d2-882b-e0ccfcf8bac6')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "5f37000d-8479-43d2-882b-e0ccfcf8bac6_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "5f37000d-8479-43d2-882b-e0ccfcf8bac6",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": null,
                        "deviceAndAppManagementAssignmentFilterType": "none",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.androidDeviceOwnerCompliancePolicy",
            "id": "6cdb97fa-8dbe-4b9d-a736-a3623c37e848",
            "displayName": "Android Enterprise - System Security | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:42:28.4762558Z",
            "roleScopeTagIds": [
                "0"
            ],
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('6cdb97fa-8dbe-4b9d-a736-a3623c37e848')/microsoft.graph.androidDeviceOwnerCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "6cdb97fa-8dbe-4b9d-a736-a3623c37e848_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "6cdb97fa-8dbe-4b9d-a736-a3623c37e848",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": null,
                        "deviceAndAppManagementAssignmentFilterType": "none",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.androidWorkProfileCompliancePolicy",
            "id": "8143c78d-5d4b-4826-ad1c-6fab0cca2a2f",
            "displayName": "Android Enterprise - Full Compliance Policy | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:42:07.646211Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "low",
            "advancedThreatProtectionRequiredSecurityLevel": "low",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('8143c78d-5d4b-4826-ad1c-6fab0cca2a2f')/microsoft.graph.androidWorkProfileCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "8143c78d-5d4b-4826-ad1c-6fab0cca2a2f_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "8143c78d-5d4b-4826-ad1c-6fab0cca2a2f",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": null,
                        "deviceAndAppManagementAssignmentFilterType": "none",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "85bbefa1-300c-420b-9877-0ff2030e7b6b",
            "displayName": "iOS/iPadOS - App Restrictions | [ADE] [BYOD] | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:43:41.7426505Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('85bbefa1-300c-420b-9877-0ff2030e7b6b')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "85bbefa1-300c-420b-9877-0ff2030e7b6b_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "85bbefa1-300c-420b-9877-0ff2030e7b6b",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "6df418a0-c390-41ba-9fe4-927eb56f6ce7",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "87df96d7-9bdb-4999-af25-31258e8a3be0",
            "displayName": "Windows - Device Health | Applies to: [Physical Autopilot Device] & [Corporate Owned Specialty Devices]\u00A0ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:00:27.1327172Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('87df96d7-9bdb-4999-af25-31258e8a3be0')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "87df96d7-9bdb-4999-af25-31258e8a3be0_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "87df96d7-9bdb-4999-af25-31258e8a3be0",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "dc20e791-31c9-47d1-8e74-ae7995cabb09",
                        "deviceAndAppManagementAssignmentFilterType": "exclude",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "a8a3be58-c8b1-47a3-8979-d23b27028bda",
            "displayName": "iOS/iPadOS - Password | [ADE] [BYOD] | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:47:44.324861Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('a8a3be58-c8b1-47a3-8979-d23b27028bda')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "a8a3be58-c8b1-47a3-8979-d23b27028bda_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "a8a3be58-c8b1-47a3-8979-d23b27028bda",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "6df418a0-c390-41ba-9fe4-927eb56f6ce7",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "b0c5f8ba-1214-406b-a9d8-31d3641888d9",
            "displayName": "Windows - Defender | Applies to: [All] ver1.1",
            "lastModifiedDateTime": "2022-07-11T09:59:23.9598648Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('b0c5f8ba-1214-406b-a9d8-31d3641888d9')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "b0c5f8ba-1214-406b-a9d8-31d3641888d9_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "b0c5f8ba-1214-406b-a9d8-31d3641888d9",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "43cb3789-2d36-4fb6-aa4d-0c3678b064e7",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "b73ab77c-9dea-410e-9a91-4996432f7348",
            "displayName": "Windows - Encryption | Applies to: [Physical Autopilot Device] & [Corporate Owned Specialty Devices] ver1.0",
            "lastModifiedDateTime": "2022-07-11T09:55:23.480291Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('b73ab77c-9dea-410e-9a91-4996432f7348')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "b73ab77c-9dea-410e-9a91-4996432f7348_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "b73ab77c-9dea-410e-9a91-4996432f7348",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "2983b1c2-8ec2-45d3-84ed-deca619d2c04",
                        "deviceAndAppManagementAssignmentFilterType": "exclude",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.macOSCompliancePolicy",
            "id": "cd1e38d4-434f-42f2-b7fa-2f75097bcd2d",
            "displayName": "macOS - Firewall [Automated Device Enrollment] [Device Enrollment] ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:02:54.7139802Z",
            "roleScopeTagIds": [
                "0"
            ],
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('cd1e38d4-434f-42f2-b7fa-2f75097bcd2d')/microsoft.graph.macOSCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "cd1e38d4-434f-42f2-b7fa-2f75097bcd2d_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "cd1e38d4-434f-42f2-b7fa-2f75097bcd2d",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "9feee0f6-92e9-4285-864c-08a1dee2b747",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "e1515187-45a9-40f3-9ef6-39d36b98ef4f",
            "displayName": "Windows - Password | Applies to: [Physical Autopilot Device] & [Virtual Corporate Device] ver1.0",
            "lastModifiedDateTime": "2022-10-04T18:18:12.7194915Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('e1515187-45a9-40f3-9ef6-39d36b98ef4f')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "e1515187-45a9-40f3-9ef6-39d36b98ef4f_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "e1515187-45a9-40f3-9ef6-39d36b98ef4f",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665",
                        "deviceAndAppManagementAssignmentFilterType": "exclude",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.macOSCompliancePolicy",
            "id": "e3b0f942-b905-421c-802b-5b521eee6543",
            "displayName": "macOS - Device Properties | [Automated Device Enrollment] [Device Enrollment] ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:01:47.321702Z",
            "roleScopeTagIds": [
                "0"
            ],
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('e3b0f942-b905-421c-802b-5b521eee6543')/microsoft.graph.macOSCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "e3b0f942-b905-421c-802b-5b521eee6543_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "e3b0f942-b905-421c-802b-5b521eee6543",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "9feee0f6-92e9-4285-864c-08a1dee2b747",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.macOSCompliancePolicy",
            "id": "e3c2cc94-cc25-474a-a655-1d9ad073ae38",
            "displayName": "macOS - Device Health | [Automated Device Enrollment] [Device Enrollment] ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:01:22.3244001Z",
            "roleScopeTagIds": [
                "0"
            ],
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('e3c2cc94-cc25-474a-a655-1d9ad073ae38')/microsoft.graph.macOSCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "e3c2cc94-cc25-474a-a655-1d9ad073ae38_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "e3c2cc94-cc25-474a-a655-1d9ad073ae38",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "9feee0f6-92e9-4285-864c-08a1dee2b747",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.macOSCompliancePolicy",
            "id": "ef8a2f10-4799-41d2-8986-b1ada11a04ff",
            "displayName": "macOS - Encryption | [Automated Device Enrollment] [Device Enrollment] ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:02:14.9366907Z",
            "roleScopeTagIds": [
                "0"
            ],
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('ef8a2f10-4799-41d2-8986-b1ada11a04ff')/microsoft.graph.macOSCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "ef8a2f10-4799-41d2-8986-b1ada11a04ff_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "ef8a2f10-4799-41d2-8986-b1ada11a04ff",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "9feee0f6-92e9-4285-864c-08a1dee2b747",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "f3e8421a-c8ea-4ca0-beba-139b69650fe3",
            "displayName": "iOS/iPadOS - Device Health | [ADE] [BYOD] [MAM-WE] | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:45:37.2030842Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "low",
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('f3e8421a-c8ea-4ca0-beba-139b69650fe3')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "f3e8421a-c8ea-4ca0-beba-139b69650fe3_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "f3e8421a-c8ea-4ca0-beba-139b69650fe3",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "6df418a0-c390-41ba-9fe4-927eb56f6ce7",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.macOSCompliancePolicy",
            "id": "f41e7c2e-558a-4efd-bd84-c08379ce531c",
            "displayName": "macOS - Gatekeeper | [Automated Device Enrollment] [Device Enrollment] ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:03:13.2568551Z",
            "roleScopeTagIds": [
                "0"
            ],
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('f41e7c2e-558a-4efd-bd84-c08379ce531c')/microsoft.graph.macOSCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "f41e7c2e-558a-4efd-bd84-c08379ce531c_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "f41e7c2e-558a-4efd-bd84-c08379ce531c",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "9feee0f6-92e9-4285-864c-08a1dee2b747",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "f5a4677b-b65a-4480-9dac-25b3250828e6",
            "displayName": "Windows - Defender Risk Score | Applies to: [Physical Autopilot Device] & [Virtual Corporate Device] ver1.0",
            "lastModifiedDateTime": "2023-05-16T12:05:01.7379733Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "secured",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('f5a4677b-b65a-4480-9dac-25b3250828e6')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "f5a4677b-b65a-4480-9dac-25b3250828e6_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "f5a4677b-b65a-4480-9dac-25b3250828e6",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665",
                        "deviceAndAppManagementAssignmentFilterType": "exclude",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.iosCompliancePolicy",
            "id": "f6645903-bd42-4776-a5fc-d753ecaf12d7",
            "displayName": "iOS/iPadOS - Email | [ADE] [BYOD] [MAM-WE] | Standard Users",
            "lastModifiedDateTime": "2022-07-11T09:48:27.4719035Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "advancedThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('f6645903-bd42-4776-a5fc-d753ecaf12d7')/microsoft.graph.iosCompliancePolicy/assignments",
            "assignments": []
        },
        {
            "@odata.type": "#microsoft.graph.macOSCompliancePolicy",
            "id": "f88090d3-a720-4ff6-bf1a-56bde7bab8b0",
            "displayName": "macOS - Password | [Automated Device Enrollment] [Device Enrollment] ver1.0",
            "lastModifiedDateTime": "2022-07-11T10:02:33.065643Z",
            "roleScopeTagIds": [
                "0"
            ],
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('f88090d3-a720-4ff6-bf1a-56bde7bab8b0')/microsoft.graph.macOSCompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "f88090d3-a720-4ff6-bf1a-56bde7bab8b0_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "f88090d3-a720-4ff6-bf1a-56bde7bab8b0",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "9feee0f6-92e9-4285-864c-08a1dee2b747",
                        "deviceAndAppManagementAssignmentFilterType": "include",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.windows10CompliancePolicy",
            "id": "fb389ebc-4a74-414c-b115-caa3b5f7fbdf",
            "displayName": "Windows - Device Security | Applies to: [Physical Autopilot Device] [Corporate Owned Specialty Devices] ver1.0",
            "lastModifiedDateTime": "2022-07-11T09:56:37.1895976Z",
            "roleScopeTagIds": [
                "0"
            ],
            "deviceThreatProtectionRequiredSecurityLevel": "unavailable",
            "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('fb389ebc-4a74-414c-b115-caa3b5f7fbdf')/microsoft.graph.windows10CompliancePolicy/assignments",
            "assignments": [
                {
                    "id": "fb389ebc-4a74-414c-b115-caa3b5f7fbdf_192584fa-43ba-4a94-93b7-d3cd4a1cc631",
                    "source": "direct",
                    "sourceId": "fb389ebc-4a74-414c-b115-caa3b5f7fbdf",
                    "target": {
                        "@odata.type": "#microsoft.graph.groupAssignmentTarget",
                        "deviceAndAppManagementAssignmentFilterId": "dc20e791-31c9-47d1-8e74-ae7995cabb09",
                        "deviceAndAppManagementAssignmentFilterType": "exclude",
                        "groupId": "192584fa-43ba-4a94-93b7-d3cd4a1cc631"
                    }
                }
            ]
        }
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsAutopilotDeploymentProfiles",
    "value": [
        {
            "@odata.type": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "id": "fc00e267-e7b3-4df0-81c1-f98aab2aaafb",
            "displayName": "Windows Autopilot Deployment Profile | User driven enrollment with AADJ",
            "description": "06.01.2022",
            "language": "os-default",
            "locale": "os-default",
            "createdDateTime": "2022-01-06T09:18:57.4835392Z",
            "lastModifiedDateTime": "2022-01-06T09:18:57.4835392Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": false,
            "hardwareHashExtractionEnabled": false,
            "deviceNameTemplate": "DT-%SERIAL%",
            "deviceType": "windowsPc",
            "enableWhiteGlove": true,
            "preprovisioningAllowed": true,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile",
            "id": "1dcd709c-7b8b-4040-847e-d4a15ced948f",
            "displayName": "Windows Autopilot Deployment Profile | User driven enrollment with HAADJ",
            "description": "Hybrid Domain Join for Autopilot\n17.03.2022",
            "language": "os-default",
            "locale": "os-default",
            "createdDateTime": "2022-03-17T21:11:52.8760856Z",
            "lastModifiedDateTime": "2022-03-17T21:18:04.8739451Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": false,
            "hardwareHashExtractionEnabled": false,
            "deviceNameTemplate": "",
            "deviceType": "windowsPc",
            "enableWhiteGlove": true,
            "preprovisioningAllowed": true,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "hybridAzureADJoinSkipConnectivityCheck": false,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "id": "6b894d4a-f95a-4992-b158-d13b8c4d7133",
            "displayName": "acc test user driven autopilot profile with os default locale",
            "description": "user driven autopilot profile with os default locale",
            "language": "os-default",
            "locale": "os-default",
            "createdDateTime": "2025-09-25T10:47:19.9485881Z",
            "lastModifiedDateTime": "2025-09-25T10:47:19.9485881Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": true,
            "hardwareHashExtractionEnabled": true,
            "deviceNameTemplate": "thing-%RAND:5%",
            "deviceType": "windowsPc",
            "enableWhiteGlove": true,
            "preprovisioningAllowed": true,
            "roleScopeTagIds": [
                "0",
                "1"
            ],
            "managementServiceAppId": null,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "id": "8273ffff-92a2-469e-95a4-ffed1c164017",
            "displayName": "acc_test_hololens_with_all_device_assignment",
            "description": "hololens autopilot profile with hk locale and all device assignment",
            "language": "zh-HK",
            "locale": "zh-HK",
            "createdDateTime": "2025-09-25T10:38:49.5120205Z",
            "lastModifiedDateTime": "2025-09-25T10:38:49.5120205Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": false,
            "hardwareHashExtractionEnabled": false,
            "deviceNameTemplate": "thing-%RAND:2%",
            "deviceType": "holoLens",
            "enableWhiteGlove": false,
            "preprovisioningAllowed": false,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "shared",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "shared",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "id": "a36aab37-ba2d-4522-9515-0a9013ccdeab",
            "displayName": "acc_test_hololens_with_all_device_assignment",
            "description": "hololens autopilot profile with hk locale and all device assignment",
            "language": "zh-HK",
            "locale": "zh-HK",
            "createdDateTime": "2025-09-25T10:39:59.8316374Z",
            "lastModifiedDateTime": "2025-09-25T10:39:59.8316374Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": false,
            "hardwareHashExtractionEnabled": false,
            "deviceNameTemplate": "thing-%RAND:2%",
            "deviceType": "holoLens",
            "enableWhiteGlove": false,
            "preprovisioningAllowed": false,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "shared",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "shared",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile",
            "id": "608440c0-6fd0-46ba-8510-dfa3c9c1d15e",
            "displayName": "acc_test_user_driven_japanese_preprovisioned",
            "description": "user driven autopilot profile with japanese locale and allow pre provisioned deployment",
            "language": "ja-JP",
            "locale": "ja-JP",
            "createdDateTime": "2025-09-25T09:28:47.4377149Z",
            "lastModifiedDateTime": "2025-09-25T09:28:47.4377149Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": true,
            "hardwareHashExtractionEnabled": true,
            "deviceNameTemplate": "",
            "deviceType": "windowsPc",
            "enableWhiteGlove": true,
            "preprovisioningAllowed": true,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "hybridAzureADJoinSkipConnectivityCheck": true,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile",
            "id": "04f8c68f-00a8-4ded-8c17-3aad2a5172e7",
            "displayName": "acc_test_user_driven_japanese_preprovisioned",
            "description": "user driven autopilot profile with japanese locale and allow pre provisioned deployment",
            "language": "ja-JP",
            "locale": "ja-JP",
            "createdDateTime": "2025-09-25T10:38:36.4532998Z",
            "lastModifiedDateTime": "2025-09-25T10:38:36.4532998Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": true,
            "hardwareHashExtractionEnabled": true,
            "deviceNameTemplate": "",
            "deviceType": "windowsPc",
            "enableWhiteGlove": true,
            "preprovisioningAllowed": true,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "hybridAzureADJoinSkipConnectivityCheck": true,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile",
            "id": "3bbaebc7-9b8c-4a1c-951f-5248b2f7f9bb",
            "displayName": "acc_test_user_driven_japanese_preprovisioned",
            "description": "user driven autopilot profile with japanese locale and allow pre provisioned deployment",
            "language": "ja-JP",
            "locale": "ja-JP",
            "createdDateTime": "2025-09-25T10:39:47.3175507Z",
            "lastModifiedDateTime": "2025-09-25T10:39:47.3175507Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": true,
            "hardwareHashExtractionEnabled": true,
            "deviceNameTemplate": "",
            "deviceType": "windowsPc",
            "enableWhiteGlove": true,
            "preprovisioningAllowed": true,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "hybridAzureADJoinSkipConnectivityCheck": true,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "singleUser",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        },
        {
            "@odata.type": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "id": "5233e587-f83c-424b-a10b-bf6ee4f2c96b",
            "displayName": "hololens",
            "description": "",
            "language": "zh-HK",
            "locale": "zh-HK",
            "createdDateTime": "2025-09-24T20:57:08.2474365Z",
            "lastModifiedDateTime": "2025-09-24T20:57:08.2474365Z",
            "enrollmentStatusScreenSettings": null,
            "extractHardwareHash": false,
            "hardwareHashExtractionEnabled": false,
            "deviceNameTemplate": "",
            "deviceType": "holoLens",
            "enableWhiteGlove": false,
            "preprovisioningAllowed": false,
            "roleScopeTagIds": [
                "0"
            ],
            "managementServiceAppId": null,
            "outOfBoxExperienceSettings": {
                "hidePrivacySettings": true,
                "hideEULA": true,
                "userType": "standard",
                "deviceUsageType": "shared",
                "skipKeyboardSelectionPage": true,
                "hideEscapeLink": true
            },
            "outOfBoxExperienceSetting": {
                "privacySettingsHidden": true,
                "eulaHidden": true,
                "userType": "standard",
                "deviceUsageType": "shared",
                "keyboardSelectionPageSkipped": true,
                "escapeLinkHidden": true
            }
        }
    ]
}


Request URL
https://graph.microsoft.com/beta/deviceAppManagement/policySets
Request Method
POST

{
  "displayName":"test",
  "description":"test",
  "assignments":[],
  "items":[
    {"@odata.type":"#microsoft.graph.mobileAppPolicySetItem",
    "payloadId":"3658ae29-7dea-4d37-8d61-1c8424742696",
    "intent":"required",
    "settings":null
    },{
      "@odata.type":"#microsoft.graph.mobileAppPolicySetItem",
      "payloadId":"a5204f11-fe5b-4ac5-b379-945d79889188",
      "intent":"required",
      "settings":
      {
        "@odata.type":"#microsoft.graph.iosStoreAppAssignmentSettings",
        "vpnConfigurationId":null,
        "uninstallOnDeviceRemoval":true,"isRemovable":true
        }
      },
      {
        "@odata.type":"#microsoft.graph.targetedManagedAppConfigurationPolicySetItem","payloadId":"A_909d2947-9f9a-4f4f-8be0-41f3079e86b6","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.targetedManagedAppConfigurationPolicySetItem","payloadId":"A_cba67ad3-5c0b-4030-97a5-9d55d15ebb7f","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.managedAppProtectionPolicySetItem","payloadId":"T_764579ab-93b2-4061-8f21-59cd221cbfe4","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.managedAppProtectionPolicySetItem","payloadId":"T_b149f8ef-4244-4ed0-8536-4fd975eda462","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.managedAppProtectionPolicySetItem","payloadId":"T_443f240c-14de-4538-a5a7-410625a9cbdc","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.managedAppProtectionPolicySetItem","payloadId":"T_7f2f1c83-0f32-41c0-aee5-ec44e1e83a83","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceConfigurationPolicySetItem","payloadId":"a9ee3ba9-2eff-4241-a29a-52b846a9a7fe","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem","payloadId":"df287451-37e4-4e45-8cc8-6aa35638cb51","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceConfigurationPolicySetItem","payloadId":"0fc6a704-662d-4334-a233-a272529aa0f9","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem","payloadId":"aacbe987-672e-4a96-afd1-3bc4c068b3ea","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem","payloadId":"6ca8f216-471b-478b-bb14-62b92fe33c54","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceConfigurationPolicySetItem","payloadId":"33d1c79b-0477-43e3-b9cc-f015f82dcf37","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"4d964631-ff9f-4da2-b59d-162c01ee1e09","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"500f7039-1ecd-4b2c-bd33-ec12310050f3","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"465712ec-13fe-4c99-a918-0496461ef884","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"110b7be5-1ca8-4504-adb1-017f5812c18c","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"3aa6e3d1-76bf-4ee4-8a31-aad8977f3538","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"229bab1c-dfde-4433-9e57-02c081e8f2c7","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"04ab1704-76f9-4306-b9a6-4dd5056bd466","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem","payloadId":"6b894d4a-f95a-4992-b158-d13b8c4d7133","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem","payloadId":"a36aab37-ba2d-4522-9515-0a9013ccdeab","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem","payloadId":"608440c0-6fd0-46ba-8510-dfa3c9c1d15e","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem","payloadId":"04f8c68f-00a8-4ded-8c17-3aad2a5172e7","guidedDeploymentTags":[]}],"guidedDeploymentTags":[],"roleScopeTags":["0"]}


resp

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/policySets/$entity",
    "id": "3c4e884a-3fe3-4b9f-b841-4dab926c5bb7",
    "createdDateTime": "2025-10-10T09:17:21.3070757Z",
    "lastModifiedDateTime": "2025-10-10T09:17:21.5258351Z",
    "displayName": "test",
    "description": "test",
    "status": "notAssigned",
    "errorCode": "noError",
    "guidedDeploymentTags": [],
    "roleScopeTags": [
        "0"
    ]
}

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/policySets/f6249df2-a80e-4497-8eb5-3898e26fe025/?$expand=assignments,items
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/policySets(assignments(),items())/$entity",
    "id": "f6249df2-a80e-4497-8eb5-3898e26fe025",
    "createdDateTime": "2025-10-10T09:55:00.1178827Z",
    "lastModifiedDateTime": "2025-10-10T09:55:00Z",
    "displayName": "unit-test",
    "description": "unit-test",
    "status": "notAssigned",
    "errorCode": "noError",
    "guidedDeploymentTags": [],
    "roleScopeTags": [
        "0"
    ],
    "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/policySets('f6249df2-a80e-4497-8eb5-3898e26fe025')/assignments",
    "assignments": [],
    "items@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/policySets('f6249df2-a80e-4497-8eb5-3898e26fe025')/items",
    "items": [
        {
            "@odata.type": "#microsoft.graph.mobileAppPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_3658ae29-7dea-4d37-8d61-1c8424742696_MobileApp_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "3658ae29-7dea-4d37-8d61-1c8424742696",
            "itemType": null,
            "displayName": "Configuration Manager Support Center OneTrace",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": [],
            "intent": "required",
            "settings": null
        },
        {
            "@odata.type": "#microsoft.graph.mobileAppPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_a5204f11-fe5b-4ac5-b379-945d79889188_MobileApp_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "a5204f11-fe5b-4ac5-b379-945d79889188",
            "itemType": null,
            "displayName": "Microsoft Intune Company Portal",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": [],
            "intent": "required",
            "settings": {
                "@odata.type": "#microsoft.graph.iosStoreAppAssignmentSettings",
                "vpnConfigurationId": null,
                "uninstallOnDeviceRemoval": true,
                "isRemovable": true,
                "preventManagedAppBackup": null
            }
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_04ab1704-76f9-4306-b9a6-4dd5056bd466_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "04ab1704-76f9-4306-b9a6-4dd5056bd466",
            "itemType": "#microsoft.graph.windows10CompliancePolicy",
            "displayName": "Windows - Device Security | Applies to: [Virtual Corporate Device]\u00A0ver1.0",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_110b7be5-1ca8-4504-adb1-017f5812c18c_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "110b7be5-1ca8-4504-adb1-017f5812c18c",
            "itemType": "#microsoft.graph.iosCompliancePolicy",
            "displayName": "iOS/iPadOS - Password | [MAM-WE] | Standard Users",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_229bab1c-dfde-4433-9e57-02c081e8f2c7_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "229bab1c-dfde-4433-9e57-02c081e8f2c7",
            "itemType": "#microsoft.graph.windows10CompliancePolicy",
            "displayName": "Windows - Device Properties | Applies to: [Physical Autopilot Device] & [Corporate Owned Specialty Devices]\u00A0ver1.0",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_3aa6e3d1-76bf-4ee4-8a31-aad8977f3538_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "3aa6e3d1-76bf-4ee4-8a31-aad8977f3538",
            "itemType": "#microsoft.graph.windows10CompliancePolicy",
            "displayName": "Windows - Device Properties | Applies to: [Corporate Owned Specialty Devices] ver1.0",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_465712ec-13fe-4c99-a918-0496461ef884_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "465712ec-13fe-4c99-a918-0496461ef884",
            "itemType": "#microsoft.graph.iosCompliancePolicy",
            "displayName": "iOS/iPadOS - Device Properties | [ADE] [BYOD] [MAM-WE] | Standard Users",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_4d964631-ff9f-4da2-b59d-162c01ee1e09_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "4d964631-ff9f-4da2-b59d-162c01ee1e09",
            "itemType": "#microsoft.graph.androidCompliancePolicy",
            "displayName": "Android Device Administrator - Full Compliance Policy [Specialty Devices] | MS Teams",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceCompliancePolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_500f7039-1ecd-4b2c-bd33-ec12310050f3_DeviceCompliance_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "500f7039-1ecd-4b2c-bd33-ec12310050f3",
            "itemType": "#microsoft.graph.windows10CompliancePolicy",
            "displayName": "Check-AvEnabled",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceConfigurationPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_0fc6a704-662d-4334-a233-a272529aa0f9_DeviceConfigurationPolicy_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "0fc6a704-662d-4334-a233-a272529aa0f9",
            "itemType": "#microsoft.graph.windows10CustomConfiguration",
            "displayName": "[Base] Prod | Windows - Custom | System - Location & Sensors ver1.0",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceConfigurationPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_33d1c79b-0477-43e3-b9cc-f015f82dcf37_DeviceConfigurationPolicy_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "33d1c79b-0477-43e3-b9cc-f015f82dcf37",
            "itemType": "#microsoft.graph.iosGeneralDeviceConfiguration",
            "displayName": "[Global] iOS/iPadOS - Device Restrictions | Password [Standard Users]",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceConfigurationPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_a9ee3ba9-2eff-4241-a29a-52b846a9a7fe_DeviceConfigurationPolicy_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "a9ee3ba9-2eff-4241-a29a-52b846a9a7fe",
            "itemType": "#microsoft.graph.windows10CustomConfiguration",
            "displayName": "[Base] Dev | Windows - Custom | Device Control - Printer Protection [User Level] ver0.1",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_6ca8f216-471b-478b-bb14-62b92fe33c54_DeviceManagementConfigurationPolicy_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "6ca8f216-471b-478b-bb14-62b92fe33c54",
            "itemType": "settingsCatalogWindows10",
            "displayName": "[Base] Prod | Windows - Settings Catalog | Lock Screen User Experience ver1.1",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_aacbe987-672e-4a96-afd1-3bc4c068b3ea_DeviceManagementConfigurationPolicy_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "aacbe987-672e-4a96-afd1-3bc4c068b3ea",
            "itemType": "settingsCatalogWindows10",
            "displayName": "[Base] Prod | Windows - Settings Catalog | Admin Templates - System - Remote Procedure Call (RPC) ver1.0",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_df287451-37e4-4e45-8cc8-6aa35638cb51_DeviceManagementConfigurationPolicy_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.1178827Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "df287451-37e4-4e45-8cc8-6aa35638cb51",
            "itemType": "settingsCatalogWindows10",
            "displayName": "[Base] Dev | Windows - Settings Catalog | Microsoft Edge - Update ver0.1",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.managedAppProtectionPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_T_443f240c-14de-4538-a5a7-410625a9cbdc_ManagedAppProtection_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "T_443f240c-14de-4538-a5a7-410625a9cbdc",
            "itemType": "#microsoft.graph.iosManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-iOS | Managed Devices v1.00 - Printing Enabled",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": [],
            "targetedAppManagementLevels": "MDM"
        },
        {
            "@odata.type": "#microsoft.graph.managedAppProtectionPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_T_764579ab-93b2-4061-8f21-59cd221cbfe4_ManagedAppProtection_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "T_764579ab-93b2-4061-8f21-59cd221cbfe4",
            "itemType": "#microsoft.graph.androidManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-Android | Managed Devices v1.00 - Printing Enabled",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": [],
            "targetedAppManagementLevels": "MDM, AndroidEnt"
        },
        {
            "@odata.type": "#microsoft.graph.managedAppProtectionPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_T_7f2f1c83-0f32-41c0-aee5-ec44e1e83a83_ManagedAppProtection_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "T_7f2f1c83-0f32-41c0-aee5-ec44e1e83a83",
            "itemType": "#microsoft.graph.iosManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-iOS | UnManaged Devices v1.00 - Printing Disabled",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": [],
            "targetedAppManagementLevels": "Unmanaged"
        },
        {
            "@odata.type": "#microsoft.graph.managedAppProtectionPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_T_b149f8ef-4244-4ed0-8536-4fd975eda462_ManagedAppProtection_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "T_b149f8ef-4244-4ed0-8536-4fd975eda462",
            "itemType": "#microsoft.graph.androidManagedAppProtection",
            "displayName": "[Global] MAM AppProtection-Android | UnManaged Devices v1.10 - Printing Disabled",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": [],
            "targetedAppManagementLevels": "Unmanaged"
        },
        {
            "@odata.type": "#microsoft.graph.targetedManagedAppConfigurationPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_A_909d2947-9f9a-4f4f-8be0-41f3079e86b6_TargetedManagedAppConfiguration_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "A_909d2947-9f9a-4f4f-8be0-41f3079e86b6",
            "itemType": null,
            "displayName": "[MAM] ACP | Microsoft Edge [iOS + Android]",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.targetedManagedAppConfigurationPolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_A_cba67ad3-5c0b-4030-97a5-9d55d15ebb7f_TargetedManagedAppConfiguration_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "A_cba67ad3-5c0b-4030-97a5-9d55d15ebb7f",
            "itemType": null,
            "displayName": "[MAM] ACP | Microsoft Outlook [iOS + Android]",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_04f8c68f-00a8-4ded-8c17-3aad2a5172e7_WindowsAutopilotDeploymentProfile_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "04f8c68f-00a8-4ded-8c17-3aad2a5172e7",
            "itemType": "#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile",
            "displayName": "acc_test_user_driven_japanese_preprovisioned",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_608440c0-6fd0-46ba-8510-dfa3c9c1d15e_WindowsAutopilotDeploymentProfile_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "608440c0-6fd0-46ba-8510-dfa3c9c1d15e",
            "itemType": "#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile",
            "displayName": "acc_test_user_driven_japanese_preprovisioned",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_6b894d4a-f95a-4992-b158-d13b8c4d7133_WindowsAutopilotDeploymentProfile_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "6b894d4a-f95a-4992-b158-d13b8c4d7133",
            "itemType": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "displayName": "acc test user driven autopilot profile with os default locale",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        },
        {
            "@odata.type": "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem",
            "id": "f6249df2-a80e-4497-8eb5-3898e26fe025_a36aab37-ba2d-4522-9515-0a9013ccdeab_WindowsAutopilotDeploymentProfile_Parcel",
            "createdDateTime": "2025-10-10T09:55:00.133509Z",
            "lastModifiedDateTime": "2025-10-10T09:55:00Z",
            "payloadId": "a36aab37-ba2d-4522-9515-0a9013ccdeab",
            "itemType": "#microsoft.graph.azureADWindowsAutopilotDeploymentProfile",
            "displayName": "acc_test_hololens_with_all_device_assignment",
            "status": "notAssigned",
            "errorCode": "noError",
            "guidedDeploymentTags": []
        }
    ]
}

update base resource

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/policySets/3976eca6-b1c0-4c5d-99c7-08814fca68bf
Request Method
PATCH

{"displayName":"test-2","description":"test-2"}

update main resource

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/policySets/3976eca6-b1c0-4c5d-99c7-08814fca68bf/update
Request Method
POST

{
    "addedPolicySetItems":[
        {"@odata.type":"#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem","payloadId":"4b196a09-803f-48dc-a877-6279579f0fc5","guidedDeploymentTags":[]
    },
    {
        "@odata.type":"#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem","payloadId":"a768c8eb-3b81-4e34-bfb5-ba34a338e1b8","guidedDeploymentTags":[]
    },
    {"@odata.type":"#microsoft.graph.windows10EnrollmentCompletionPageConfigurationPolicySetItem","payloadId":"acdf7778-98be-4086-8a43-f5d89b305229_Windows10EnrollmentCompletionPageConfiguration","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"f6645903-bd42-4776-a5fc-d753ecaf12d7","guidedDeploymentTags":[]},{"@odata.type":"#microsoft.graph.deviceCompliancePolicyPolicySetItem","payloadId":"f88090d3-a720-4ff6-bf1a-56bde7bab8b0","guidedDeploymentTags":[]}],"updatedPolicySetItems":[{"@odata.type":"#microsoft.graph.mobileAppPolicySetItem","payloadId":"b200f17f-1f5b-4806-8985-9d2a910a47f9","intent":"required","settings":null,"guidedDeploymentTags":[]}],"deletedPolicySetItems":["3976eca6-b1c0-4c5d-99c7-08814fca68bf_04ab1704-76f9-4306-b9a6-4dd5056bd466_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_04f8c68f-00a8-4ded-8c17-3aad2a5172e7_WindowsAutopilotDeploymentProfile_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_0fc6a704-662d-4334-a233-a272529aa0f9_DeviceConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_110b7be5-1ca8-4504-adb1-017f5812c18c_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_229bab1c-dfde-4433-9e57-02c081e8f2c7_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_33d1c79b-0477-43e3-b9cc-f015f82dcf37_DeviceConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_3aa6e3d1-76bf-4ee4-8a31-aad8977f3538_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_3bbaebc7-9b8c-4a1c-951f-5248b2f7f9bb_WindowsAutopilotDeploymentProfile_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_465712ec-13fe-4c99-a918-0496461ef884_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_4d964631-ff9f-4da2-b59d-162c01ee1e09_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_500f7039-1ecd-4b2c-bd33-ec12310050f3_DeviceCompliance_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_608440c0-6fd0-46ba-8510-dfa3c9c1d15e_WindowsAutopilotDeploymentProfile_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_6b894d4a-f95a-4992-b158-d13b8c4d7133_WindowsAutopilotDeploymentProfile_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_6ca8f216-471b-478b-bb14-62b92fe33c54_DeviceManagementConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_A_936da198-54cb-4c9c-886a-0947e71d63b9_TargetedManagedAppConfiguration_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_a36aab37-ba2d-4522-9515-0a9013ccdeab_WindowsAutopilotDeploymentProfile_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_a5204f11-fe5b-4ac5-b379-945d79889188_MobileApp_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_a9ee3ba9-2eff-4241-a29a-52b846a9a7fe_DeviceConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_aacbe987-672e-4a96-afd1-3bc4c068b3ea_DeviceManagementConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_b0079b61-617d-4b3f-8222-09e641c75644_DeviceConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_df287451-37e4-4e45-8cc8-6aa35638cb51_DeviceManagementConfigurationPolicy_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_T_443f240c-14de-4538-a5a7-410625a9cbdc_ManagedAppProtection_Parcel","3976eca6-b1c0-4c5d-99c7-08814fca68bf_T_b149f8ef-4244-4ed0-8536-4fd975eda462_ManagedAppProtection_Parcel"]}

assignments

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/policySets/3976eca6-b1c0-4c5d-99c7-08814fca68bf/update
Request Method
POST

{"assignments":[{"target":{"@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget","groupId":"d7e6ac2e-4b31-481f-81ca-3e0cab841674"}},{"target":{"@odata.type":"#microsoft.graph.exclusionGroupAssignmentTarget","groupId":"2e3819e7-d935-4419-a5c0-b256777b27bc"}},{"target":{"@odata.type":"#microsoft.graph.groupAssignmentTarget","groupId":"fdf755b0-cce8-43ed-b047-17dac8c905af"}},{"target":{"@odata.type":"#microsoft.graph.groupAssignmentTarget","groupId":"dd13e083-30e9-451b-9df3-7f3880e39fa3"}}]}