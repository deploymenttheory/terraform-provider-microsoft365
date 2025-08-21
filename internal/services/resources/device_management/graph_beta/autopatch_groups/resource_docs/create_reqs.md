Request URL
https://services.autopatch.microsoft.com/device/v2/autopatchGroups
Request Method
POST

request

{
  "name":"auto-patch-group",
  "description":"",
  "globalUserManagedAadGroups":[],
  "deploymentGroups":[
    {"aadId":"00000000-0000-0000-0000-000000000000",
    "name":"auto-patch-group - Test",
    "userManagedAadGroups":[
      {"name":"[Azure]-[ConditonalAccess]-[Prod]-[CAD003-PolicyExclude]-[UG]",
      "id":"410a28bd-9c9f-403f-b1b2-4a0bd04e98d9","type":0}
      ],
      "failedPreRequisiteCheckCount":0,
      "deploymentGroupPolicySettings":
        {
          "aadGroupName":"auto-patch-group - Test",
          "isUpdateSettingsModified":false,
          "deviceConfigurationSetting":
            {"policyId":"000",
            "updateBehavior":"AutoInstallAndRestart",
            "notificationSetting":"DefaultNotifications",
            "qualityDeploymentSettings":
            {
              "deadline":1,
              "deferral":0,
              "gracePeriod":0
            },
            "featureDeploymentSettings":
            {
              "deadline":5,
              "deferral":0
            },
            "updateFrequencyUI":null,
            "installDays":null,
            "installTime":null,
            "activeHourEndTime":null,
            "activeHourStartTime":null
            },
            "dnfUpdateCloudSetting":
            {
              "policyId":"000",
              "approvalType":"Automatic",
              "deploymentDeferralInDays":0
            },
            "officeDCv2Setting":
            {
              "policyId":"000",
              "deadline":1,
              "deferral":0,"
              hideUpdateNotifications":false,
              "targetChannel":"MonthlyEnterprise"
              },
              "edgeDCv2Setting":
              {"policyId":"000",
              "targetChannel":"Beta"
              },
              "featureUpdateAnchorCloudSetting":
              {
                "targetOSVersion":"Windows 11, version 24H2",
              "installLatestWindows10OnWindows11IneligibleDevice":true}}},
              {"name":"auto-patch-group - Ring1",
              "userManagedAadGroups":[
                {"name":"[Azure]-[ConditonalAccess]-[Prod]-[CAD002-PolicyExclude]-[UG]",
                "id":"35d09841-af73-43e6-a59f-024fef1b6b95","type":0}],
                "aadId":"00000000-0000-0000-0000-000000000000",
                "deploymentGroupPolicySettings":
                {
                  "aadGroupName":"auto-patch-group - Ring1",
                  "isUpdateSettingsModified":false,
                  "deviceConfigurationSetting":
                  {
                    "policyId":"000",
                    "updateBehavior":"AutoInstallAndRestart",
                    "notificationSetting":"DefaultNotifications",
                    "qualityDeploymentSettings":
                    {
                      "deadline":2,
                      "deferral":1,
                      "gracePeriod":2},
                      "featureDeploymentSettings":
                      {
                        "deadline":5,"deferral":0},
                        "updateFrequencyUI":null,
                        "installDays":null,
                        "installTime":null,
                        "activeHourEndTime":null,
                        "activeHourStartTime":null},
                        "dnfUpdateCloudSetting":
                        {
                          "policyId":"000",
                          "approvalType":"Automatic",
                          "deploymentDeferralInDays":1},
                          "officeDCv2Setting":
                          {
                            "policyId":"000",
                            "deadline":2,
                            "deferral":1,
                            "hideUpdateNotifications":false,
                            "targetChannel":"MonthlyEnterprise"},
                            "edgeDCv2Setting":
                            {
                              "policyId":"000",
                              "targetChannel":"Stable"},
                              "featureUpdateAnchorCloudSetting":
                              {
                                "targetOSVersion":"Windows 11, version 24H2",
                                "installLatestWindows10OnWindows11IneligibleDevice":true
                                }
                              }
                            },
                            {"aadId":"00000000-0000-0000-0000-000000000000",
                            "name":"auto-patch-group - Last","userManagedAadGroups":[
                              {"name":"[Azure]-[ConditonalAccess]-[Prod]-[CAD005-PolicyExclude]-[UG]",
                              "id":"48fe6d79-f045-448a-bd74-716db27f0783","type":0
                              }
                            ],
                            "failedPreRequisiteCheckCount":0,"deploymentGroupPolicySettings":{"aadGroupName":"auto-patch-group - Last","isUpdateSettingsModified":false,"deviceConfigurationSetting":{"policyId":"000","updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deadline":3,"deferral":5,"gracePeriod":2},"featureDeploymentSettings":{"deadline":5,"deferral":0},"updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"dnfUpdateCloudSetting":{"policyId":"000","approvalType":"Automatic","deploymentDeferralInDays":5},"officeDCv2Setting":{"policyId":"000","deadline":3,"deferral":5,"hideUpdateNotifications":false,"targetChannel":"MonthlyEnterprise"},"edgeDCv2Setting":{"policyId":"000","targetChannel":"Stable"},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 24H2","installLatestWindows10OnWindows11IneligibleDevice":true}}}],"windowsUpdateSettings":[],"status":"Unknown","type":"Unknown","distributionType":"Unknown","driverUpdateSettings":[],"enableDriverUpdate":true,"scopeTags":[0],"enabledContentTypes":31}

response - n/a

Request URL
https://services.autopatch.microsoft.com/device/v2/autopatchGroups
Request Method
GET

[
  {
    "id": "4aa9b805-9494-4eed-a04b-ed51ec9e631e",
    "name": "Windows Autopatch",
    "tenantId": "2fd6bb84-ad40-4ec5-9369-a215b25c9952",
    "description": "Windows Autopatch, the default Autopatch Group",
    "type": "Default",
    "status": "Active",
    "isLockedByPolicy": false,
    "distributionType": "Mixed",
    "deploymentGroups": [
      {
        "userManagedAadGroups": [],
        "aadId": "e7acdbca-81e8-4d69-84a5-c927b82dbf2b",
        "name": "Windows Autopatch - Test"
      },
      {
        "distribution": 1,
        "userManagedAadGroups": [],
        "aadId": "3e923eeb-9a6e-45ec-a39a-122848c80ae9",
        "name": "Windows Autopatch - Ring1"
      },
      {
        "distribution": 9,
        "userManagedAadGroups": [],
        "aadId": "bb1eaf1c-4271-4ee3-953d-9e8448ac7787",
        "name": "Windows Autopatch - Ring2"
      },
      {
        "distribution": 90,
        "userManagedAadGroups": [],
        "aadId": "efff769e-1bd7-4038-ae77-ab6378c05fcb",
        "name": "Windows Autopatch - Ring3"
      },
      {
        "userManagedAadGroups": [],
        "aadId": "20a81ae9-ac3e-48d7-aa04-b3501f80e46e",
        "name": "Windows Autopatch - Last"
      }
    ],
    "policyBasedDeploymentGroups": [],
    "globalUserManagedAadGroups": [
      {
        "id": "24b51abd-ce88-42dd-9058-753b80011a57",
        "type": "Device"
      }
    ],
    "numberOfRegisteredDevices": 0,
    "readOnly": false,
    "scopeTags": [
      0
    ],
    "flowStatus": "Succeeded",
    "umbrellaGroupId": "a5db54fb-62da-478f-90ec-01f149eaa1fb",
    "userHasAllScopeTag": true
  },
  {
    "id": "66ad5d6f-3494-409d-b8cb-d39521200474",
    "name": "auto-patch-group",
    "tenantId": "2fd6bb84-ad40-4ec5-9369-a215b25c9952",
    "description": "",
    "type": "User",
    "status": "Creating",
    "isLockedByPolicy": false,
    "distributionType": "AdminAssigned",
    "deploymentGroups": [
      {
        "userManagedAadGroups": [
          {
            "id": "410a28bd-9c9f-403f-b1b2-4a0bd04e98d9",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "auto-patch-group - Test"
      },
      {
        "userManagedAadGroups": [
          {
            "id": "35d09841-af73-43e6-a59f-024fef1b6b95",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "auto-patch-group - Ring1"
      },
      {
        "userManagedAadGroups": [
          {
            "id": "48fe6d79-f045-448a-bd74-716db27f0783",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "auto-patch-group - Last"
      }
    ],
    "policyBasedDeploymentGroups": [],
    "globalUserManagedAadGroups": [],
    "numberOfRegisteredDevices": 0,
    "readOnly": false,
    "scopeTags": [
      0
    ],
    "flowId": "15064391-d437-450e-ac93-e7865df76398",
    "flowType": "Create",
    "flowStatus": "InProgress",
    "userHasAllScopeTag": true
  }
]