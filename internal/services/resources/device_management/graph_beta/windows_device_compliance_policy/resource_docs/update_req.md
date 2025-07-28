{
  "id":"78b936bd-0911-49ce-89b9-70b22b01cdd2",
  "displayName":"Windows 10/11 - Basic Compliance Policy",
  "description":"Basic Windows device compliance policy requiring BitLocker, Secure Boot and a password",
  "roleScopeTagIds":["0"],
  "@odata.type":"#microsoft.graph.windows10CompliancePolicy",
  "passwordRequired":true,
  "passwordBlockSimple":true,
  "passwordRequiredToUnlockFromIdle":true,
  "passwordMinutesOfInactivityBeforeLock":15,
  "passwordMinimumLength":8,
  "passwordMinimumCharacterSetCount":3,
  "passwordRequiredType":"alphanumeric",
  "requireHealthyDeviceReport":false,
  "osMinimumVersion":"10.0.19041.0",
  "earlyLaunchAntiMalwareDriverEnabled":false,
  "bitLockerEnabled":true,
  "secureBootEnabled":true,
  "codeIntegrityEnabled":true,
  "memoryIntegrityEnabled":false,
  "kernelDmaProtectionEnabled":false,
  "virtualizationBasedSecurityEnabled":false,
  "firmwareProtectionEnabled":false,
  "storageRequireEncryption":true,
  "activeFirewallRequired":false,
  "defenderEnabled":true,
  "signatureOutOfDate":false,
  "rtpEnabled":true,
  "antivirusRequired":true,
  "antiSpywareRequired":true,
  "deviceThreatProtectionEnabled":false,
  "deviceThreatProtectionRequiredSecurityLevel":"unavailable",
  "configurationManagerComplianceRequired":false,
  "tpmRequired":false,
  "validOperatingSystemBuildRanges":[],
  "wslDistributions":[],
  "assignments@odata.context":"https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('78b936bd-0911-49ce-89b9-70b22b01cdd2')/microsoft.graph.windows10CompliancePolicy/assignments",
  "scheduledActionsForRule@odata.context":"https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCompliancePolicies('78b936bd-0911-49ce-89b9-70b22b01cdd2')/microsoft.graph.windows10CompliancePolicy/scheduledActionsForRule(scheduledActionConfigurations())"}


// https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies/78b936bd-0911-49ce-89b9-70b22b01cdd2/scheduleActionsForRules

{"deviceComplianceScheduledActionForRules":[
  {"ruleName":"PasswordRequired",
  "scheduledActionConfigurations":[
    {"actionType":"block",
    "gracePeriodHours":1152,
    "notificationTemplateId":"00000000-0000-0000-0000-000000000000",
    "notificationMessageCCList":[]
    },
    {"actionType":"notification",
    "gracePeriodHours":0,
    "notificationTemplateId":"426e6351-c6ff-44d3-910d-8b937ee30bdd",
    "notificationMessageCCList":["aa856a09-cf0c-4b31-a315-cb53251e54d8","a77240dc-2827-47af-8fcb-e209a67e176a"]
    },
    {"actionType":"notification",
    "gracePeriodHours":120,
    "notificationTemplateId":"cd80867e-0ef4-4eeb-9117-67d2ce5c8bd5",
    "notificationMessageCCList":["17ae470d-1d67-4a1d-b55f-c2402b294e6d","6d3c36da-fa2c-45f0-97c6-60588faba39c"]
    },
    {"actionType":"retire",
    "gracePeriodHours":1440,
    "notificationTemplateId":"",
    "notificationMessageCCList":[]
    }]
  }]
}