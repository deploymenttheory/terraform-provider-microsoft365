Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations
Request Method
POST

request body

{"@odata.type":"#microsoft.graph.windows10EnrollmentCompletionPageConfiguration","id":"39f2b3af-f4b2-41d3-adc3-4f898f31c058","displayName":"ttest","description":"test","showInstallationProgress":true,"blockDeviceSetupRetryByUser":false,"allowDeviceResetOnInstallFailure":true,"allowLogCollectionOnInstallFailure":true,"customErrorMessage":"Setup could not be completed. Please try again or contact your support person for help.","installProgressTimeoutInMinutes":60,"allowDeviceUseOnInstallFailure":true,"selectedMobileAppIds":["e83d36e1-3ff2-4567-90d9-940919184ad5","3e17a039-a8c6-4cac-b892-ee5b9d64b2b6","e2dd1944-5854-4b0d-951c-7dbd8edab3a8"],"trackInstallProgressForAutopilotOnly":true,"disableUserStatusTrackingAfterFirstUser":true,"roleScopeTagIds":["0"],"allowNonBlockingAppInstallation":true,"installQualityUpdates":true}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations/$entity",
    "@odata.type": "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration",
    "id": "4998ae20-fa77-4f5c-9f7f-a5956b42dd6c_Windows10EnrollmentCompletionPageConfiguration",
    "displayName": "ttest",
    "description": "test",
    "priority": 3,
    "createdDateTime": "2025-08-21T04:44:51.8433388Z",
    "lastModifiedDateTime": "2025-08-21T04:44:51.8433388Z",
    "version": 1,
    "roleScopeTagIds": [
        "0"
    ],
    "deviceEnrollmentConfigurationType": "windows10EnrollmentCompletionPageConfiguration",
    "showInstallationProgress": true,
    "blockDeviceSetupRetryByUser": false,
    "allowDeviceResetOnInstallFailure": true,
    "allowLogCollectionOnInstallFailure": true,
    "customErrorMessage": "Setup could not be completed. Please try again or contact your support person for help.",
    "installProgressTimeoutInMinutes": 60,
    "allowDeviceUseOnInstallFailure": true,
    "selectedMobileAppIds": [
        "e83d36e1-3ff2-4567-90d9-940919184ad5",
        "3e17a039-a8c6-4cac-b892-ee5b9d64b2b6",
        "e2dd1944-5854-4b0d-951c-7dbd8edab3a8"
    ],
    "allowNonBlockingAppInstallation": true,
    "installQualityUpdates": false,
    "trackInstallProgressForAutopilotOnly": true,
    "disableUserStatusTrackingAfterFirstUser": true
}

get by id

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations/4998ae20-fa77-4f5c-9f7f-a5956b42dd6c_Windows10EnrollmentCompletionPageConfiguration?$expand=assignments
Request Method
GET

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations(assignments())/$entity",
    "@odata.type": "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration",
    "id": "4998ae20-fa77-4f5c-9f7f-a5956b42dd6c_Windows10EnrollmentCompletionPageConfiguration",
    "displayName": "ttest",
    "description": "test",
    "priority": 3,
    "createdDateTime": "2025-08-21T04:44:51.8433388Z",
    "lastModifiedDateTime": "2025-08-21T04:44:51.8433388Z",
    "version": 1,
    "roleScopeTagIds": [
        "0"
    ],
    "deviceEnrollmentConfigurationType": "windows10EnrollmentCompletionPageConfiguration",
    "showInstallationProgress": true,
    "blockDeviceSetupRetryByUser": false,
    "allowDeviceResetOnInstallFailure": true,
    "allowLogCollectionOnInstallFailure": true,
    "customErrorMessage": "Setup could not be completed. Please try again or contact your support person for help.",
    "installProgressTimeoutInMinutes": 60,
    "allowDeviceUseOnInstallFailure": true,
    "selectedMobileAppIds": [
        "e83d36e1-3ff2-4567-90d9-940919184ad5",
        "3e17a039-a8c6-4cac-b892-ee5b9d64b2b6",
        "e2dd1944-5854-4b0d-951c-7dbd8edab3a8"
    ],
    "allowNonBlockingAppInstallation": true,
    "installQualityUpdates": false,
    "trackInstallProgressForAutopilotOnly": true,
    "disableUserStatusTrackingAfterFirstUser": true,
    "assignments@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations('4998ae20-fa77-4f5c-9f7f-a5956b42dd6c_Windows10EnrollmentCompletionPageConfiguration')/microsoft.graph.windows10EnrollmentCompletionPageConfiguration/assignments",
    "assignments": []
}

assignments supports - include groups, all users, all devices
- no filters
- noexclusions


for app look up

Request URL
https://graph.microsoft.com/beta/deviceAppManagement/mobileApps?$filter=isof(%27microsoft.graph.windowsAppX%27)%20or%20isof(%27microsoft.graph.windowsMobileMSI%27)%20or%20isof(%27microsoft.graph.windowsUniversalAppX%27)%20or%20isof(%27microsoft.graph.officeSuiteApp%27)%20or%20isof(%27microsoft.graph.windowsMicrosoftEdgeApp%27)%20or%20isof(%27microsoft.graph.winGetApp%27)%20or%20isof(%27microsoft.graph.win32LobApp%27)%20or%20isof(%27microsoft.graph.win32CatalogApp%27)&$top=250&$orderBy=displayname
Request Method
GET

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps",
    "@odata.count": 0,
    "value": [
        {
            "@odata.type": "#microsoft.graph.winGetApp",
            "id": "e4938228-aab3-493b-a9d5-8250aa8e9d55",
            "displayName": "3CX",
            "description": "The 3CX Softphone notifies users in real-time of incoming calls and enables easy call management from their desktop. Users can set their status, view CRM & phonebook entries as well as access their BLF panel for fast calling. \n\nGetting Started:\n\n   1. Download the 3CX App and open it.\n   2. Enter your 3CX URL - found in your “Welcome to 3CX” email.\n   3. Enter your email or extension number along with your password and click on “Login”.\n   4. The softphone will auto-provision.\n\nNote: If 2FA is enabled for your extension, you will be requested to enter your security PIN.\n\nSign in using your Google or Microsoft 365 account by:\n\n   1. Open the softphone and enter your 3CX URL -  found in your “Welcome to 3CX” email.\n   2. Click on the Google or M365 buttons.\n   3. Your 3CX Web Client will open, log in. The Windows Tile dialog box will appear.\n   4. Click on the “Provision” button.\n   5. The softphone will auto-provision.\n\nLearn more: https://www.3cx.com/user-manual/windows-softphone-app/",
            "publisher": "3CX Software DMCC",
            "largeIcon": null,
            "createdDateTime": "2024-09-30T13:47:10Z",
            "lastModifiedDateTime": "2024-09-30T13:47:10Z",
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
            "manifestHash": null,
            "packageIdentifier": "9NW77489NGJ0",
            "installExperience": {
                "runAsAccount": "system"
            }
        },
        {
            "@odata.type": "#microsoft.graph.win32LobApp",
            "id": "e83d36e1-3ff2-4567-90d9-940919184ad5",
            "displayName": "Azure Cli",
            "description": "The Azure Command-Line Interface (CLI) is a cross-platform command-line tool that can be installed locally on Windows computers. You can use the Azure CLI for Windows to connect to Azure and execute administrative commands on Azure resources. The Azure CLI for Windows can also be used from a browser through the Azure Cloud Shell or run from inside a Docker container.\n\nRef: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-windows?tabs=azure-powershell",
            "publisher": "Microsoft",
            "largeIcon": null,
            "createdDateTime": "2022-02-04T14:31:32Z",
            "lastModifiedDateTime": "2022-02-04T15:11:39Z",
            "isFeatured": false,
            "privacyInformationUrl": null,
            "informationUrl": null,
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
            "fileName": "Azure CLI - Latest.intunewin",
            "size": 208864,
            "installCommandLine": "Deploy-Application.exe",
            "uninstallCommandLine": "Deploy-Application.exe - Uninstall",
            "applicableArchitectures": "x64",
            "allowedArchitectures": "x64,arm64",
            "minimumFreeDiskSpaceInMB": null,
            "minimumMemoryInMB": null,
            "minimumNumberOfProcessors": null,
            "minimumCpuSpeedInMHz": null,
            "msiInformation": null,
            "setupFilePath": "Deploy-Application.exe",
            "minimumSupportedWindowsRelease": "1903",
            "displayVersion": "Latest",
            "allowAvailableUninstall": false,
            "minimumSupportedOperatingSystem": {
                "v8_0": false,
                "v8_1": false,
                "v10_0": false,
                "v10_1607": false,
                "v10_1703": false,
                "v10_1709": false,
                "v10_1803": false,
                "v10_1809": false,
                "v10_1903": true,
                "v10_1909": false,
                "v10_2004": false,
                "v10_2H20": false,
                "v10_21H1": false
            },
            "detectionRules": [
                {
                    "@odata.type": "#microsoft.graph.win32LobAppFileSystemDetection",
                    "path": "C:\\Program Files (x86)\\Microsoft SDKs\\Azure\\CLI2\\wbin",
                    "fileOrFolderName": "az.cmd",
                    "check32BitOn64System": false,
                    "detectionType": "exists",
                    "operator": "notConfigured",
                    "detectionValue": null
                }
            ],
            "requirementRules": [],
            "rules": [
                {
                    "@odata.type": "#microsoft.graph.win32LobAppFileSystemRule",
                    "ruleType": "detection",
                    "path": "C:\\Program Files (x86)\\Microsoft SDKs\\Azure\\CLI2\\wbin",
                    "fileOrFolderName": "az.cmd",
                    "check32BitOn64System": false,
                    "operationType": "exists",
                    "operator": "notConfigured",
                    "comparisonValue": null
                }
            ],
            "installExperience": {
                "runAsAccount": "system",
                "maxRunTimeInMinutes": 60,
                "deviceRestartBehavior": "allow"
            },
            "returnCodes": [
                {
                    "returnCode": 0,
                    "type": "success"
                },
                {
                    "returnCode": 1707,
                    "type": "success"
                },
                {
                    "returnCode": 3010,
                    "type": "softReboot"
                },
                {
                    "returnCode": 1641,
                    "type": "hardReboot"
                },
                {
                    "returnCode": 1618,
                    "type": "retry"
                }
            ]
        },
        {
            "@odata.type": "#microsoft.graph.winGetApp",
            "id": "3e17a039-a8c6-4cac-b892-ee5b9d64b2b6",
            "displayName": "Azure Virtual Desktop Preview",
            "description": "The Azure Virtual Desktop store app is no longer available for download or installation.\n \nTo ensure a seamless experience and avoid any disruption, users are encouraged to download the Windows App. Windows App is the gateway to securely connect to any devices or apps across Azure Virtual Desktop, Windows 365, and Microsoft Dev Box.\n \nAzure Virtual Desktop is a cloud VDI service that delivers secure remote desktop and app experiences from virtually anywhere. It provides the flexibility and control organizations need with exclusive support for Windows 11 and Windows 10 multi-session cost-savings capabilities and the built-in security and reliability of Azure. Company data is safe and secure because it lives in the cloud and not on your personal devices.\n \nLearn more: aka.ms/AVDstoreapp",
            "publisher": "Microsoft Corporation",
            "largeIcon": null,
            "createdDateTime": "2024-10-02T05:26:45Z",
            "lastModifiedDateTime": "2024-10-02T05:26:45Z",
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
            "manifestHash": null,
            "packageIdentifier": "9NZSG2H7MS6B",
            "installExperience": {
                "runAsAccount": "user"
            }
        },