
#microsoft.graph.windows10CompliancePolicy

wsl varient

{"id":"00000000-0000-0000-0000-000000000000","displayName":"test","description":"test","roleScopeTagIds":["0"],
"@odata.type":"#microsoft.graph.windows10CompliancePolicy",
"scheduledActionsForRule":[
  {"ruleName":"PasswordRequired",
  "scheduledActionConfigurations":[
    {
      "actionType":"block",
      "gracePeriodHours":0,
      "notificationTemplateId":"",
      "notificationMessageCCList":[]
    },
    {"actionType":"notification",
    "gracePeriodHours":0,"notificationTemplateId":"6eb92564-c50c-46a8-b9a3-8aa0d3c09ade","notificationMessageCCList":[
      "17ae470d-1d67-4a1d-b55f-c2402b294e6d"]
      },
    {"actionType":"retire",
      "gracePeriodHours":24,
      "notificationTemplateId":"",
      "notificationMessageCCList":[]
    }]}
  ],
  "deviceThreatProtectionRequiredSecurityLevel":"medium","deviceCompliancePolicyScript":null,"passwordRequiredType":"alphanumeric","wslDistributions":[{"distribution":"linux","minimumOSVersion":"1","maximumOSVersion":"2"},{"distribution":"linux2","minimumOSVersion":"1","maximumOSVersion":"2"}],"passwordRequired":true,"passwordBlockSimple":true,"passwordRequiredToUnlockFromIdle":true,"storageRequireEncryption":true,"passwordMinutesOfInactivityBeforeLock":15,"passwordMinimumCharacterSetCount":3,"activeFirewallRequired":true,"tpmRequired":true,"antivirusRequired":true,"antiSpywareRequired":true,"defenderEnabled":true,"signatureOutOfDate":true,"rtpEnabled":true,"defenderVersion":"4.11.0.0","configurationManagerComplianceRequired":true,"osMinimumVersion":"1","osMaximumVersion":"2","mobileOsMinimumVersion":"1","mobileOsMaximumVersion":"2","secureBootEnabled":true,"bitLockerEnabled":true,"codeIntegrityEnabled":true,"deviceThreatProtectionEnabled":true}

Custom Compliance varient

{"id":"00000000-0000-0000-0000-000000000000","displayName":"test","roleScopeTagIds":["0","9","8"],"@odata.type":"#microsoft.graph.windows10CompliancePolicy","scheduledActionsForRule":[{"ruleName":"PasswordRequired","scheduledActionConfigurations":[{"actionType":"block","gracePeriodHours":0,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"notification","gracePeriodHours":240,"notificationTemplateId":"2ba73548-f201-4ed0-8c02-6c6a3d50afad","notificationMessageCCList":["17ae470d-1d67-4a1d-b55f-c2402b294e6d","6d3c36da-fa2c-45f0-97c6-60588faba39c"]},{"actionType":"notification","gracePeriodHours":360,"notificationTemplateId":"cd80867e-0ef4-4eeb-9117-67d2ce5c8bd5","notificationMessageCCList":["aa856a09-cf0c-4b31-a315-cb53251e54d8","a77240dc-2827-47af-8fcb-e209a67e176a"]},{"actionType":"retire","gracePeriodHours":480,"notificationTemplateId":"","notificationMessageCCList":[]}]}],"deviceThreatProtectionRequiredSecurityLevel":"medium","passwordRequiredType":"alphanumeric","deviceCompliancePolicyScript":{"deviceComplianceScriptId":"8c3d2ec3-3e63-4df3-8265-69bbba1e53e5","rulesContent":"ewoiUnVsZXMiOlsgAgICAgICAgIH0KICAgICAgIF0KICAgIH0KIF0KfQ=="},"customComplianceRequired":true,"passwordRequired":true,"passwordBlockSimple":true,"passwordMinimumCharacterSetCount":3,"passwordMinutesOfInactivityBeforeLock":480,"passwordRequiredToUnlockFromIdle":true,"storageRequireEncryption":true,"activeFirewallRequired":true,"tpmRequired":true,"antivirusRequired":true,"antiSpywareRequired":true,"defenderEnabled":true,"signatureOutOfDate":true,"defenderVersion":"4.11.0.0","rtpEnabled":true,"deviceThreatProtectionEnabled":true}