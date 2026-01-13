example 2

Request URL
https://services.autopatch.microsoft.com/device/v2/autopatchGroups
Request Method
POST

{"name":"test","description":"test2","globalUserManagedAadGroups":[{"name":"[Intune]-[ManagedDevices]-[Prod]-[Corp]-[All_1909_Devices]-[W10]-[DG]","id":"550dba96-abd7-4ef0-9cf1-25be705f676c","type":0},{"name":"[Intune]-[ManagedDevices]-[Prod]-[Corp]-[All_2009_Devices]-[W10]-[DG]","id":"6a08c3a0-1693-4089-80cb-f9c2f8063a3b","type":0}],"deploymentGroups":[{"aadId":"00000000-0000-0000-0000-000000000000","name":"test - Test","userManagedAadGroups":[{"name":"[Azure]-[ConditonalAccess]-[Prod]-[CAP001-PolicyExclude]-[UG]","id":"5a0832d3-19b9-4f78-9b50-906774ac4d49","type":0}],"failedPreRequisiteCheckCount":0,"deploymentGroupPolicySettings":{"aadGroupName":"test - Test","isUpdateSettingsModified":false,"deviceConfigurationSetting":{"policyId":"000","updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deadline":1,"deferral":0,"gracePeriod":0},"featureDeploymentSettings":{"deadline":5,"deferral":0},"updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"dnfUpdateCloudSetting":{"policyId":"000","approvalType":"Automatic","deploymentDeferralInDays":0},"officeDCv2Setting":{"policyId":"000","deadline":1,"deferral":0,"hideUpdateNotifications":false,"targetChannel":"MonthlyEnterprise"},"edgeDCv2Setting":{"policyId":"000","targetChannel":"Beta"},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true}}},{"name":"test - Ring1","distribution":75,"userManagedAadGroups":[{"name":"[Azure]-[ConditonalAccess]-[Prod]-[CAU010-PolicyExclude]-[UG]","id":"0fe3d2cb-62ae-4fa4-858f-a122061ada62","type":0}],"aadId":"00000000-0000-0000-0000-000000000000","deploymentGroupPolicySettings":{"aadGroupName":"test - Ring1","isUpdateSettingsModified":false,"deviceConfigurationSetting":{"policyId":"000","updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deadline":2,"deferral":1,"gracePeriod":2},"featureDeploymentSettings":{"deadline":5,"deferral":0},"updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"dnfUpdateCloudSetting":{"approvalType":"Automatic","policyId":"000","deploymentDeferralInDays":1},"officeDCv2Setting":{"policyId":"000","deadline":2,"deferral":1,"hideUpdateNotifications":false,"targetChannel":"MonthlyEnterprise"},"edgeDCv2Setting":{"policyId":"000","targetChannel":"Stable"},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true}}},{"name":"test - Ring2","distribution":25,"userManagedAadGroups":[{"name":"[Azure]-[RBAC]-[Intune-Administrator]-[UG]","id":"20f4b274-0ca2-406a-ae20-12cf99730d62","type":0}],"aadId":"00000000-0000-0000-0000-000000000000","deploymentGroupPolicySettings":{"aadGroupName":"test - Ring2","isUpdateSettingsModified":false,"deviceConfigurationSetting":{"policyId":"000","updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deadline":3,"deferral":5,"gracePeriod":2},"featureDeploymentSettings":{"deadline":5,"deferral":0},"updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"dnfUpdateCloudSetting":{"approvalType":"Manual","policyId":"000","deploymentDeferralInDays":null},"officeDCv2Setting":{"policyId":"000","deadline":3,"deferral":5,"hideUpdateNotifications":false,"targetChannel":"MonthlyEnterprise"},"edgeDCv2Setting":{"policyId":"000","targetChannel":"Stable"},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true}}},{"aadId":"00000000-0000-0000-0000-000000000000","name":"test - Last","userManagedAadGroups":[{"name":"[Azure]-[RBAC]-[Azure-DevOps-Administrator]-[UG]","id":"6c57971c-4369-4569-9d4b-29c9351665e1","type":0}],"failedPreRequisiteCheckCount":0,"deploymentGroupPolicySettings":{"aadGroupName":"test - Last","isUpdateSettingsModified":false,"deviceConfigurationSetting":{"policyId":"000","updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deadline":5,"deferral":9,"gracePeriod":2},"featureDeploymentSettings":{"deadline":5,"deferral":0},"updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"dnfUpdateCloudSetting":{"approvalType":"Manual","policyId":"000","deploymentDeferralInDays":null},"officeDCv2Setting":{"policyId":"000","deadline":5,"deferral":9,"hideUpdateNotifications":false,"targetChannel":"MonthlyEnterprise"},"edgeDCv2Setting":{"policyId":"000","targetChannel":"Stable"},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true}}}],"windowsUpdateSettings":[],"status":"Unknown","type":"Unknown","distributionType":"Unknown","driverUpdateSettings":[],"enableDriverUpdate":true,"scopeTags":[0,1232,1234],"enabledContentTypes":31}

GET - but resource hasnt reached active state yet

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
    "status": "Active",
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
        "aadId": "e33d75c2-5533-4700-a54c-63b248fbcb42",
        "name": "auto-patch-group - Test"
      },
      {
        "userManagedAadGroups": [
          {
            "id": "35d09841-af73-43e6-a59f-024fef1b6b95",
            "type": "Device"
          }
        ],
        "aadId": "d6aeec90-ea6b-42e5-9ed0-989098097bbb",
        "name": "auto-patch-group - Ring1"
      },
      {
        "userManagedAadGroups": [
          {
            "id": "48fe6d79-f045-448a-bd74-716db27f0783",
            "type": "Device"
          }
        ],
        "aadId": "0473d745-c16d-45e2-986c-65ac5dbc7c10",
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
    "flowStatus": "Succeeded",
    "umbrellaGroupId": "8cf21f14-cb3a-4212-b670-7032cde105d0",
    "userHasAllScopeTag": true
  },
  {
    "id": "6733a2ee-d2a9-4930-ba83-403514923e8a",
    "name": "test",
    "tenantId": "2fd6bb84-ad40-4ec5-9369-a215b25c9952",
    "description": "test2",
    "type": "User",
    "status": "Creating",
    "isLockedByPolicy": false,
    "distributionType": "Mixed",
    "deploymentGroups": [
      {
        "userManagedAadGroups": [
          {
            "id": "5a0832d3-19b9-4f78-9b50-906774ac4d49",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "test - Test"
      },
      {
        "distribution": 75,
        "userManagedAadGroups": [
          {
            "id": "0fe3d2cb-62ae-4fa4-858f-a122061ada62",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "test - Ring1"
      },
      {
        "distribution": 25,
        "userManagedAadGroups": [
          {
            "id": "20f4b274-0ca2-406a-ae20-12cf99730d62",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "test - Ring2"
      },
      {
        "userManagedAadGroups": [
          {
            "id": "6c57971c-4369-4569-9d4b-29c9351665e1",
            "type": "Device"
          }
        ],
        "aadId": "00000000-0000-0000-0000-000000000000",
        "name": "test - Last"
      }
    ],
    "policyBasedDeploymentGroups": [],
    "globalUserManagedAadGroups": [
      {
        "id": "550dba96-abd7-4ef0-9cf1-25be705f676c",
        "type": "None"
      },
      {
        "id": "6a08c3a0-1693-4089-80cb-f9c2f8063a3b",
        "type": "None"
      }
    ],
    "numberOfRegisteredDevices": 0,
    "readOnly": false,
    "scopeTags": [
      0,
      1232,
      1234
    ],
    "flowId": "7bc711b5-db2a-4cd2-9bde-efec61b8ea7f",
    "flowType": constants.TfOperationCreate,
    "flowStatus": "InProgress",
    "userHasAllScopeTag": true
  }
]

GET - resource has reached active state

{
  "autopatchGroupId": "6733a2ee-d2a9-4930-ba83-403514923e8a",
  "enabledFeatureTypes": 31,
  "deploymentGroupPolicies": [
    {
      "aadGroupId": "eae6f5b3-9ee9-49f8-991e-81d780f89de9",
      "deploymentGroupPolicySettings": {
        "deviceConfigurationSetting": {
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deferral": 0,
            "deadline": 1,
            "gracePeriod": 0
          },
          "featureDeploymentSettings": {
            "deferral": 0,
            "deadline": 5
          },
          "policyId": "b97bff48-1138-4ad2-b2b2-9f51f98df847"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 25H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true,
          "policyId": "297070c5-0748-4ec0-9452-b51d03fe1d3a"
        },
        "dnfUpdateCloudSetting": {
          "approvalType": "Automatic",
          "deploymentDeferralInDays": 0,
          "policyId": "6e638095-13ee-407d-8eb8-2c46e47461b4"
        },
        "edgeDCv2Setting": {
          "targetChannel": "Beta",
          "policyId": "1d6ff71e-38c6-45cc-bb86-f61dd19c4022"
        },
        "officeDCv2Setting": {
          "targetChannel": "MonthlyEnterprise",
          "deferral": 0,
          "deadline": 1,
          "hideUpdateNotifications": false,
          "enableAutomaticUpdate": true,
          "hideEnableDisableUpdate": true,
          "enableOfficeMgmt": false,
          "updatePath": "http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6",
          "policyId": "bae5a6d7-12bf-48d0-96ad-08f462f5e7e1"
        }
      }
    },
    {
      "aadGroupId": "5d06e413-696f-41ae-a199-fbf39f8b0eda",
      "deploymentGroupPolicySettings": {
        "deviceConfigurationSetting": {
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deferral": 1,
            "deadline": 2,
            "gracePeriod": 2
          },
          "featureDeploymentSettings": {
            "deferral": 0,
            "deadline": 5
          },
          "policyId": "94eae107-6752-438c-8cc9-a256596df957"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 25H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true,
          "policyId": "297070c5-0748-4ec0-9452-b51d03fe1d3a"
        },
        "dnfUpdateCloudSetting": {
          "approvalType": "Automatic",
          "deploymentDeferralInDays": 1,
          "policyId": "4ae0a5ef-7387-404f-bfe2-ca51f664ce80"
        },
        "edgeDCv2Setting": {
          "targetChannel": "Stable",
          "policyId": "3db7f07e-cc1b-46fb-8b19-b5b62b103dce"
        },
        "officeDCv2Setting": {
          "targetChannel": "MonthlyEnterprise",
          "deferral": 1,
          "deadline": 2,
          "hideUpdateNotifications": false,
          "enableAutomaticUpdate": true,
          "hideEnableDisableUpdate": true,
          "enableOfficeMgmt": false,
          "updatePath": "http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6",
          "policyId": "361fd3f5-ce5d-4578-8452-624527606d16"
        }
      }
    },
    {
      "aadGroupId": "edb5086a-6c50-4708-b19c-65e1ca8c15b3",
      "deploymentGroupPolicySettings": {
        "deviceConfigurationSetting": {
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deferral": 5,
            "deadline": 3,
            "gracePeriod": 2
          },
          "featureDeploymentSettings": {
            "deferral": 0,
            "deadline": 5
          },
          "policyId": "0b6de189-60cd-4067-830c-73529a6d4b20"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 25H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true,
          "policyId": "297070c5-0748-4ec0-9452-b51d03fe1d3a"
        },
        "dnfUpdateCloudSetting": {
          "approvalType": "Manual",
          "policyId": "9148fd63-251d-4960-aedc-fbba80fc1005"
        },
        "edgeDCv2Setting": {
          "targetChannel": "Stable",
          "policyId": "fce5a64c-0ed2-4209-9fc2-ab0a34bcb118"
        },
        "officeDCv2Setting": {
          "targetChannel": "MonthlyEnterprise",
          "deferral": 5,
          "deadline": 3,
          "hideUpdateNotifications": false,
          "enableAutomaticUpdate": true,
          "hideEnableDisableUpdate": true,
          "enableOfficeMgmt": false,
          "updatePath": "http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6",
          "policyId": "373f5a73-b8e5-4411-98b5-a05ea3da497c"
        }
      }
    },
    {
      "aadGroupId": "8e142eb4-e791-4553-af7d-a4f45af2e89d",
      "deploymentGroupPolicySettings": {
        "deviceConfigurationSetting": {
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deferral": 9,
            "deadline": 5,
            "gracePeriod": 2
          },
          "featureDeploymentSettings": {
            "deferral": 0,
            "deadline": 5
          },
          "policyId": "fdea4e0c-1042-4ed7-946d-b84e87cbdd8e"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 25H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true,
          "policyId": "297070c5-0748-4ec0-9452-b51d03fe1d3a"
        },
        "dnfUpdateCloudSetting": {
          "approvalType": "Manual",
          "policyId": "5d1c5352-b847-4ddd-8746-f2441114e944"
        },
        "edgeDCv2Setting": {
          "targetChannel": "Stable",
          "policyId": "ccff1a87-a3bb-4532-a06e-cce3d3e48d1f"
        },
        "officeDCv2Setting": {
          "targetChannel": "MonthlyEnterprise",
          "deferral": 9,
          "deadline": 5,
          "hideUpdateNotifications": false,
          "enableAutomaticUpdate": true,
          "hideEnableDisableUpdate": true,
          "enableOfficeMgmt": false,
          "updatePath": "http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6",
          "policyId": "12d2c372-2b69-41dc-b09e-dff2fcced5f0"
        }
      }
    }
  ]
}

update

Request URL
https://services.autopatch.microsoft.com/device/v2/autopatchGroups
Request Method
PUT

{"id":"6733a2ee-d2a9-4930-ba83-403514923e8a","name":"test","description":"test - update","globalUserManagedAadGroups":[{"id":"550dba96-abd7-4ef0-9cf1-25be705f676c","type":"None"},{"id":"6a08c3a0-1693-4089-80cb-f9c2f8063a3b","type":"None"}],"deploymentGroups":[{"userManagedAadGroups":[{"id":"5a0832d3-19b9-4f78-9b50-906774ac4d49","type":"Device"}],"aadId":"eae6f5b3-9ee9-49f8-991e-81d780f89de9","name":"test - Test","deploymentGroupPolicySettings":{"aadGroupName":"test - Test","deviceConfigurationSetting":{"updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deferral":0,"deadline":1,"gracePeriod":0},"featureDeploymentSettings":{"deferral":0,"deadline":5},"policyId":"b97bff48-1138-4ad2-b2b2-9f51f98df847","updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true},"dnfUpdateCloudSetting":{"approvalType":"Automatic","deploymentDeferralInDays":0,"policyId":"6e638095-13ee-407d-8eb8-2c46e47461b4"},"edgeDCv2Setting":{"targetChannel":"Beta","policyId":"1d6ff71e-38c6-45cc-bb86-f61dd19c4022"},"officeDCv2Setting":{"targetChannel":"MonthlyEnterprise","deferral":0,"deadline":1,"hideUpdateNotifications":false,"enableAutomaticUpdate":true,"hideEnableDisableUpdate":true,"enableOfficeMgmt":false,"updatePath":"http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6","policyId":"bae5a6d7-12bf-48d0-96ad-08f462f5e7e1"}}},{"distribution":75,"userManagedAadGroups":[{"id":"0fe3d2cb-62ae-4fa4-858f-a122061ada62","type":"Device"}],"aadId":"5d06e413-696f-41ae-a199-fbf39f8b0eda","name":"test - Ring1","deploymentGroupPolicySettings":{"aadGroupName":"test - Ring1","deviceConfigurationSetting":{"updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deferral":1,"deadline":2,"gracePeriod":2},"featureDeploymentSettings":{"deferral":0,"deadline":5},"policyId":"94eae107-6752-438c-8cc9-a256596df957","updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true},"dnfUpdateCloudSetting":{"approvalType":"Automatic","deploymentDeferralInDays":1,"policyId":"4ae0a5ef-7387-404f-bfe2-ca51f664ce80"},"edgeDCv2Setting":{"targetChannel":"Stable","policyId":"3db7f07e-cc1b-46fb-8b19-b5b62b103dce"},"officeDCv2Setting":{"targetChannel":"MonthlyEnterprise","deferral":1,"deadline":2,"hideUpdateNotifications":false,"enableAutomaticUpdate":true,"hideEnableDisableUpdate":true,"enableOfficeMgmt":false,"updatePath":"http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6","policyId":"361fd3f5-ce5d-4578-8452-624527606d16"}}},{"distribution":25,"userManagedAadGroups":[{"id":"20f4b274-0ca2-406a-ae20-12cf99730d62","type":"Device"}],"aadId":"edb5086a-6c50-4708-b19c-65e1ca8c15b3","name":"test - Ring2","deploymentGroupPolicySettings":{"aadGroupName":"test - Ring2","deviceConfigurationSetting":{"updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deferral":5,"deadline":3,"gracePeriod":2},"featureDeploymentSettings":{"deferral":0,"deadline":5},"policyId":"0b6de189-60cd-4067-830c-73529a6d4b20","updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true},"dnfUpdateCloudSetting":{"approvalType":"Manual","policyId":"9148fd63-251d-4960-aedc-fbba80fc1005","deploymentDeferralInDays":null},"edgeDCv2Setting":{"targetChannel":"Stable","policyId":"fce5a64c-0ed2-4209-9fc2-ab0a34bcb118"},"officeDCv2Setting":{"targetChannel":"MonthlyEnterprise","deferral":5,"deadline":3,"hideUpdateNotifications":false,"enableAutomaticUpdate":true,"hideEnableDisableUpdate":true,"enableOfficeMgmt":false,"updatePath":"http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6","policyId":"373f5a73-b8e5-4411-98b5-a05ea3da497c"}}},{"userManagedAadGroups":[{"id":"6c57971c-4369-4569-9d4b-29c9351665e1","type":"Device"}],"aadId":"8e142eb4-e791-4553-af7d-a4f45af2e89d","name":"test - Last","deploymentGroupPolicySettings":{"aadGroupName":"test - Last","deviceConfigurationSetting":{"updateBehavior":"AutoInstallAndRestart","notificationSetting":"DefaultNotifications","qualityDeploymentSettings":{"deferral":9,"deadline":5,"gracePeriod":2},"featureDeploymentSettings":{"deferral":0,"deadline":5},"policyId":"fdea4e0c-1042-4ed7-946d-b84e87cbdd8e","updateFrequencyUI":null,"installDays":null,"installTime":null,"activeHourEndTime":null,"activeHourStartTime":null},"featureUpdateAnchorCloudSetting":{"targetOSVersion":"Windows 11, version 25H2","installLatestWindows10OnWindows11IneligibleDevice":true},"dnfUpdateCloudSetting":{"approvalType":"Manual","policyId":"5d1c5352-b847-4ddd-8746-f2441114e944","deploymentDeferralInDays":null},"edgeDCv2Setting":{"targetChannel":"Stable","policyId":"ccff1a87-a3bb-4532-a06e-cce3d3e48d1f"},"officeDCv2Setting":{"targetChannel":"MonthlyEnterprise","deferral":9,"deadline":5,"hideUpdateNotifications":false,"enableAutomaticUpdate":true,"hideEnableDisableUpdate":true,"enableOfficeMgmt":false,"updatePath":"http://officecdn.microsoft.com/pr/55336b82-a18d-4dd6-b5f6-9e5095c314a6","policyId":"12d2c372-2b69-41dc-b09e-dff2fcced5f0"}}}],"status":"Unknown","type":"User","distributionType":"Unknown","windowsUpdateSettings":[],"driverUpdateSettings":[],"enableDriverUpdate":true,"scopeTags":[0,1232,1234],"enabledContentTypes":31}