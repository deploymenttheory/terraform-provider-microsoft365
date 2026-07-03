Request URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/59a48111-7836-4822-8bc9-99fe2cfdbd08/enrollmentProfiles/59a48111-7836-4822-8bc9-99fe2cfdbd08_34867b1a-9246-448b-bf8a-cad609162fcd/setDefaultProfile
Request Method
POST
Status Code
200

Request URL
https://graph.microsoft.com/beta/deviceManagement/configurationPolicies?$select=id,name,description,platforms,lastModifiedDateTime,technologies,settingCount,isAssigned,templateReference,creationSource%20&$top=500%20&$filter=(technologies%20has%20%27enrollment%27)%20and%20((((((platforms%20eq%20%27ios%27)%20and%20(TemplateReference/templateId%20eq%20%2727d20e9c-50c1-48f8-a44c-f37de4510051_1%27)))%20or%20((platforms%20eq%20%27macOS%27)%20and%20(TemplateReference/templateId%20eq%20%272e29557d-70fc-405a-8082-d1e5b6be2b8c_1%27)))%20or%20((platforms%20eq%20%27visionOS%27)%20and%20(TemplateReference/templateId%20eq%20%27b974d292-62c4-4853-bdab-166cf42df51c_1%27)))%20or%20((platforms%20eq%20%27tvOS%27)%20and%20(TemplateReference/templateId%20eq%20%2746af9f2a-ea44-41f5-9a93-baeebd93776e_1%27)))%20and%20(creationSource%20eq%20%27DepTokenId_59a48111-7836-4822-8bc9-99fe2cfdbd08%27)%20%20
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies(id,name,description,platforms,lastModifiedDateTime,technologies,settingCount,isAssigned,templateReference,creationSource)",
    "@odata.count": 1,
    "value": [
        {
            "id": "34867b1a-9246-448b-bf8a-cad609162fcd",
            "name": "acc-test-macos-device-enrollment-policy-maximal",
            "description": "macOS ADE enrollment policy exercising the full settings tree",
            "platforms": "macOS",
            "lastModifiedDateTime": "2026-07-03T12:50:58.4181194Z",
            "technologies": "enrollment",
            "settingCount": 28,
            "isAssigned": false,
            "creationSource": "DepTokenId_59a48111-7836-4822-8bc9-99fe2cfdbd08",
            "templateReference": {
                "templateId": "2e29557d-70fc-405a-8082-d1e5b6be2b8c_1",
                "templateFamily": "enrollmentConfiguration",
                "templateDisplayName": "Apple macOS Device Enrollment Policy",
                "templateDisplayVersion": "Version 1"
            }
        }
    ]
}

equest URL
https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings/59a48111-7836-4822-8bc9-99fe2cfdbd08/?$expand=defaultiosenrollmentprofile,defaultmacosenrollmentprofile,defaulttvosenrollmentprofile,defaultvisionosenrollmentprofile
Request Method
GET

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/depOnboardingSettings(defaultIosEnrollmentProfile(),defaultMacOsEnrollmentProfile(),defaultTvOSEnrollmentProfile(),defaultVisionOSEnrollmentProfile())/$entity",
    "id": "59a48111-7836-4822-8bc9-99fe2cfdbd08",
    "appleIdentifier": "dafydd.watkins@bankofscotland.appleaccount.com",
    "tokenExpirationDateTime": "2027-03-27T08:57:42Z",
    "lastModifiedDateTime": "2026-03-27T08:57:58.9351782Z",
    "lastSuccessfulSyncDateTime": "2026-04-14T07:06:10.2510832Z",
    "lastSyncTriggeredDateTime": "2026-03-27T08:57:59.9985264Z",
    "shareTokenWithSchoolDataSyncService": false,
    "lastSyncErrorCode": 3,
    "tokenType": "dep",
    "tokenName": "lab_deploymenttheory_intune",
    "syncedDeviceCount": 0,
    "dataSharingConsentGranted": true,
    "roleScopeTagIds": [
        "0"
    ],
    "defaultIosEnrollmentProfile": null,
    "defaultMacOsEnrollmentProfile": {
        "id": "ECV2_59a48111-7836-4822-8bc9-99fe2cfdbd08_34867b1a-9246-448b-bf8a-cad609162fcd",
        "displayName": "acc-test-macos-device-enrollment-policy-maximal",
        "description": "macOS ADE enrollment policy exercising the full settings tree",
        "requiresUserAuthentication": true,
        "configurationEndpointUrl": "https://appleconfigurator2.manage.microsoft.com/EnrollmentServer/MDMServiceConfig?id=54fac284-7866-43e5-860a-9c8e10fa3d7d&AADTenantId=2fd6bb84-ad40-4ec5-9369-a215b25c9952",
        "enableAuthenticationViaCompanyPortal": false,
        "requireCompanyPortalOnSetupAssistantEnrolledDevices": false,
        "isDefault": false,
        "supervisedModeEnabled": true,
        "supportDepartment": "IT Support",
        "isMandatory": true,
        "locationDisabled": false,
        "supportPhoneNumber": "+1-555-0100",
        "profileRemovalDisabled": true,
        "restoreBlocked": true,
        "appleIdDisabled": true,
        "termsAndConditionsDisabled": false,
        "touchIdDisabled": false,
        "applePayDisabled": true,
        "siriDisabled": true,
        "diagnosticsDisabled": true,
        "displayToneSetupDisabled": false,
        "privacyPaneDisabled": true,
        "screenTimeScreenDisabled": true,
        "deviceNameTemplate": "",
        "configurationWebUrl": true,
        "enabledSkipKeys": [
            "Restore",
            "AppleID",
            "Payment",
            "Siri",
            "Diagnostics",
            "iCloudDiagnostics",
            "iCloudStorage",
            "ScreenTime",
            "Privacy",
            "UnlockWithWatch",
            "EnableLockdownMode",
            "UpdateCompleted",
            "TermsOfAddress",
            "OSShowcase"
        ],
        "enrollmentTimeAzureAdGroupIds": [],
        "waitForDeviceConfiguredConfirmation": true,
        "registrationDisabled": false,
        "fileVaultDisabled": false,
        "iCloudDiagnosticsDisabled": true,
        "passCodeDisabled": false,
        "zoomDisabled": false,
        "iCloudStorageDisabled": true,
        "chooseYourLockScreenDisabled": false,
        "accessibilityScreenDisabled": false,
        "autoUnlockWithWatchDisabled": true,
        "skipPrimarySetupAccountCreation": false,
        "setPrimarySetupAccountAsRegularUser": true,
        "dontAutoPopulatePrimaryAccountInfo": false,
        "primaryAccountFullName": "Primary User",
        "primaryAccountUserName": "primaryuser",
        "enableRestrictEditing": true,
        "adminAccountUserName": "localadmin",
        "adminAccountFullName": "Local Administrator",
        "adminAccountPassword": null,
        "hideAdminAccount": true,
        "requestRequiresNetworkTether": false,
        "autoAdvanceSetupEnabled": false,
        "usePlatformSSODuringSetupAssistant": false,
        "depProfileAdminAccountPasswordRotationSetting": {
            "autoRotationPeriodInDays": 90,
            "depProfileDelayAutoRotationSetting": null
        }
    },
    "defaultTvOSEnrollmentProfile": null,
    "defaultVisionOSEnrollmentProfile": null
}