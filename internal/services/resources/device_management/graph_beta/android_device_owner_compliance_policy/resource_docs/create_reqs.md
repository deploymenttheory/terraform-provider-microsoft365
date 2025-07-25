
	//aospDeviceOwnerCompliancePolicy

{
	"id":"00000000-0000-0000-0000-000000000000",
	"displayName":"testt",
	"description":"test",
	"roleScopeTagIds":["0"],
	"@odata.type":"#microsoft.graph.aospDeviceOwnerCompliancePolicy",
	"scheduledActionsForRule":[
		{"ruleName":"PasswordRequired",
		"scheduledActionConfigurations":[
			{"actionType":"block",
			"gracePeriodHours":0,
			"notificationTemplateId":"",
			"notificationMessageCCList":[]
			},
			{"actionType":
				"notification",
				"gracePeriodHours":24,
				"notificationTemplateId":"2ba73548-f201-4ed0-8c02-6c6a3d50afad",
				"notificationMessageCCList":[
					"aa856a09-cf0c-4b31-a315-cb53251e54d8",
					"a77240dc-2827-47af-8fcb-e209a67e176a"]
				}
			]
		}
	],
	"localActions":[],
	"passwordMinimumLength":null,
	"passwordRequiredType":"required",
	"securityBlockJailbrokenDevices":true,
	"osMinimumVersion":"10",
	"osMaximumVersion":"11",
	"minAndroidSecurityPatchLevel":"2016-10-01",
	"passwordRequired":true,
	"passwordMinutesOfInactivityBeforeLock":10,
	"storageRequireEncryption":true
}

#microsoft.graph.androidDeviceOwnerCompliancePolicy

{
	"id":"00000000-0000-0000-0000-000000000000",
	"displayName":"test",
	"description":"test",
	"roleScopeTagIds":["0"],
	"@odata.type":"#microsoft.graph.androidDeviceOwnerCompliancePolicy",
	"scheduledActionsForRule":
		[
			{"ruleName":"PasswordRequired",
				"scheduledActionConfigurations":[
					{"actionType":"block",
						"gracePeriodHours":0,
						"notificationTemplateId":"",
						"notificationMessageCCList":[]},
					{"actionType":"remoteLock",
						"gracePeriodHours":0,
						"notificationTemplateId":"",
						"notificationMessageCCList":[]},
					{"actionType":"notification",
						"gracePeriodHours":0,
						"notificationTemplateId":"426e6351-c6ff-44d3-910d-8b937ee30bdd",
						"notificationMessageCCList":[]},
					{"actionType":"remoteLock",
						"gracePeriodHours":48,
						"notificationTemplateId":"",
						"notificationMessageCCList":[]},
					{"actionType":"retire",
						"gracePeriodHours":240,
						"notificationTemplateId":"",
						"notificationMessageCCList":[]
					}
				]
			}
		],
	"localActions":[],
	"deviceThreatProtectionRequiredSecurityLevel":"secured","advancedThreatProtectionRequiredSecurityLevel":"medium","passwordRequiredType":"numeric","securityRequiredAndroidSafetyNetEvaluationType":"hardwareBacked","osMinimumVersion":"10",
	"osMaximumVersion":"11",
	"minAndroidSecurityPatchLevel":"2016-10-01",
	"passwordRequired":true,"passwordMinimumLength":16,"passwordMinutesOfInactivityBeforeLock":5,
	"passwordExpirationDays":365,
	"passwordPreviousPasswordCountToBlock":23,
	"storageRequireEncryption":true,
	"securityRequireIntuneAppIntegrity":true,
	"deviceThreatProtectionEnabled":true,"securityRequireSafetyNetAttestationBasicIntegrity":true,"securityRequireSafetyNetAttestationCertifiedDevice":null
}


"#microsoft.graph.iosCompliancePolicy"

{"id":"00000000-0000-0000-0000-000000000000","displayName":"test","description":"test","roleScopeTagIds":["0"],"@odata.type":"#microsoft.graph.iosCompliancePolicy","scheduledActionsForRule":[{"ruleName":"PasswordRequired","scheduledActionConfigurations":[{"actionType":"block","gracePeriodHours":0,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"pushNotification","gracePeriodHours":0,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"retire","gracePeriodHours":240,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"remoteLock","gracePeriodHours":480,"notificationTemplateId":"","notificationMessageCCList":[]}]}],"deviceThreatProtectionRequiredSecurityLevel":"secured","advancedThreatProtectionRequiredSecurityLevel":"medium","passcodeRequiredType":"alphanumeric","managedEmailProfileRequired":true,"securityBlockJailbrokenDevices":true,"osMinimumVersion":"10","osMaximumVersion":"11","osMinimumBuildVersion":"20E772520a","osMaximumBuildVersion":"20E772520a","passcodeRequired":true,"passcodeMinimumCharacterSetCount":1,"passcodeMinutesOfInactivityBeforeLock":1,"passcodeMinutesOfInactivityBeforeScreenTimeout":4,"passcodeExpirationDays":730,"passcodePreviousPasscodeBlockCount":24,"restrictedApps":[{"name":"thing","appId":"com.thing.id"},{"name":"thing2","appId":"com.thing.id2"}],"deviceThreatProtectionEnabled":true}

#microsoft.graph.windows10CompliancePolicy

{"id":"00000000-0000-0000-0000-000000000000","displayName":"test","description":"test","roleScopeTagIds":["0"],"@odata.type":"#microsoft.graph.windows10CompliancePolicy","scheduledActionsForRule":[{"ruleName":"PasswordRequired","scheduledActionConfigurations":[{"actionType":"block","gracePeriodHours":0,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"notification","gracePeriodHours":0,"notificationTemplateId":"6eb92564-c50c-46a8-b9a3-8aa0d3c09ade","notificationMessageCCList":["17ae470d-1d67-4a1d-b55f-c2402b294e6d"]},{"actionType":"retire","gracePeriodHours":24,"notificationTemplateId":"","notificationMessageCCList":[]}]}],"deviceThreatProtectionRequiredSecurityLevel":"medium","deviceCompliancePolicyScript":null,"passwordRequiredType":"alphanumeric","wslDistributions":[{"distribution":"linux","minimumOSVersion":"1","maximumOSVersion":"2"},{"distribution":"linux2","minimumOSVersion":"1","maximumOSVersion":"2"}],"passwordRequired":true,"passwordBlockSimple":true,"passwordRequiredToUnlockFromIdle":true,"storageRequireEncryption":true,"passwordMinutesOfInactivityBeforeLock":15,"passwordMinimumCharacterSetCount":3,"activeFirewallRequired":true,"tpmRequired":true,"antivirusRequired":true,"antiSpywareRequired":true,"defenderEnabled":true,"signatureOutOfDate":true,"rtpEnabled":true,"defenderVersion":"4.11.0.0","configurationManagerComplianceRequired":true,"osMinimumVersion":"1","osMaximumVersion":"2","mobileOsMinimumVersion":"1","mobileOsMaximumVersion":"2","secureBootEnabled":true,"bitLockerEnabled":true,"codeIntegrityEnabled":true,"deviceThreatProtectionEnabled":true}


"#microsoft.graph.macOSCompliancePolicy"

{"id":"00000000-0000-0000-0000-000000000000","displayName":"test","description":"test","roleScopeTagIds":["0"],"@odata.type":"#microsoft.graph.macOSCompliancePolicy","scheduledActionsForRule":[{"ruleName":"PasswordRequired","scheduledActionConfigurations":[{"actionType":"block","gracePeriodHours":720,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"notification","gracePeriodHours":480,"notificationTemplateId":"49cc1559-4424-4ebe-9a61-13ca10458de1","notificationMessageCCList":[]},{"actionType":"remoteLock","gracePeriodHours":240,"notificationTemplateId":"","notificationMessageCCList":[]},{"actionType":"retire","gracePeriodHours":0,"notificationTemplateId":"","notificationMessageCCList":[]}]}],"gatekeeperAllowedAppSource":"macAppStoreAndIdentifiedDevelopers","passwordRequiredType":"alphanumeric","systemIntegrityProtectionEnabled":true,"osMinimumVersion":"10","osMaximumVersion":"11","osMinimumBuildVersion":"20E772520a","osMaximumBuildVersion":"20E772520a","passwordRequired":true,"passwordBlockSimple":true,"passwordMinimumCharacterSetCount":2,"passwordMinutesOfInactivityBeforeLock":15,"storageRequireEncryption":true,"firewallEnabled":true,"firewallBlockAllIncoming":true,"firewallEnableStealthMode":true}