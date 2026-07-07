package graphBetaMacOSDeviceEnrollmentPolicy

// Setting definition IDs, and setting instance/value template GUIDs, were captured from
// Graph X-Ray traces of the Intune admin center creating/updating a macOS Automated Device
// Enrollment (ADE) profile via the settings catalog backed `/deviceManagement/configurationPolicies`
// endpoint. See docs/create.md, docs/create_2.md and docs/update.md for the raw captured payloads.
const (
	// TemplateID is the settings catalog template used for macOS ADE/DEP enrollment profiles.
	TemplateID = "2e29557d-70fc-405a-8082-d1e5b6be2b8c_1"

	// TemplateFamily for enrollment restriction/configuration policies.
	TemplateFamily = "enrollmentConfiguration"

	// Platforms and technologies values required by Graph for this template.
	Platforms    = "macOS"
	Technologies = "enrollment"

	// CreationSourcePrefix is prepended to the resolved dep_onboarding_settings_id to build the
	// `creationSource` value sent on Create. It is not sent on Update.
	CreationSourcePrefix = "DepTokenId_"

	// Setting Definition IDs - User Affinity / Authentication
	SettingDefUserAffinity         = "ade_macos_useraffinity"
	SettingDefAuthenticationMethod = "ade_macos_authenticationmethod"

	// Setting Definition IDs - Await Configuration / Account Settings
	SettingDefAwaitConfiguration           = "ade_macos_awaitconfiguration"
	SettingDefCreateLocalAdmin             = "ade_accountsettings_createlocaladmin"
	SettingDefAdminAccountName             = "ade_accountsettings_adminaccountname"
	SettingDefAdminAccountFullName         = "ade_accountsettings_adminaccountfullname"
	SettingDefHideUsersGroups              = "ade_accountsettings_hideusersgroups"
	SettingDefAdminAccountPasswordRotation = "ade_accountsettings_adminaccountpasswordrotation"
	SettingDefCreateLocalPrimary           = "ade_accountsettings_createlocalprimary"
	SettingDefPrefillAccountInfo           = "ade_accountsettings_prefillaccountinfo"
	SettingDefRestrictEditing              = "ade_accountsettings_restrictediting"
	SettingDefPrimaryAccountFullName       = "ade_accountsettings_primaryaccountfullname"
	SettingDefPrimaryAccountName           = "ade_accountsettings_primaryaccountname"

	// Setting Definition IDs - Locked Enrollment
	SettingDefLockedEnrollment = "ade_lockedenrollment"

	// Setting Definition IDs - Setup Assistant
	SettingDefDepartment              = "ade_setupassistant_department"
	SettingDefDepartmentPhone         = "ade_setupassistant_departmentphone"
	SettingDefLocationServices        = "ade_setupassistant_locationservices"
	SettingDefRestore                 = "ade_setupassistant_restore"
	SettingDefAppleId                 = "ade_setupassistant_appleid"
	SettingDefTermsAndConditions      = "ade_setupassistant_termsandconditions"
	SettingDefTouchFaceId             = "ade_setupassistant_touchfaceid"
	SettingDefApplePay                = "ade_setupassistant_applepay"
	SettingDefSiri                    = "ade_setupassistant_siri"
	SettingDefDiagnosticsData         = "ade_setupassistant_diagnosticsdata"
	SettingDefFileVault               = "ade_setupassistant_filevault"
	SettingDefICloudDiagnostics       = "ade_setupassistant_iclouddiagnostics"
	SettingDefICloudStorage           = "ade_setupassistant_icloudstorage"
	SettingDefAppearance              = "ade_setupassistant_appearance"
	SettingDefScreenTime              = "ade_setupassistant_screentime"
	SettingDefPrivacy                 = "ade_setupassistant_privacy"
	SettingDefAccessibility           = "ade_setupassistant_accessibility"
	SettingDefUnlockWithWatch         = "ade_setupassistant_unlockwithwatch"
	SettingDefEnableLockdownMode      = "ade_setupassistant_enablelockdownmode"
	SettingDefSoftwareUpdate          = "ade_setupassistant_softwareupdate"
	SettingDefSoftwareUpdateCompleted = "ade_setupassistant_softwareupdatecompleted"
	SettingDefTermsOfAddress          = "ade_setupassistant_termsofaddress"
	SettingDefIntelligence            = "ade_setupassistant_intelligence"
	SettingDefOSShowcase              = "ade_setupassistant_osshowcase"
	SettingDefAppStore                = "ade_setupassistant_appstore"

	// Setting Instance Template IDs
	SettingInstanceTemplateUserAffinity            = "f3ddbbae-bb3d-45e8-bb9c-854f674e56ef"
	SettingInstanceTemplateAwaitConfiguration      = "7e4d9c8a-2f5b-4a3c-8d9e-1b6a4f7c2e5d"
	SettingInstanceTemplateLockedEnrollment        = "e105d3a4-7318-454a-a231-0d1510ebf1e3"
	SettingInstanceTemplateDepartment              = "1fd43c7c-0c7a-4dcc-8392-7b1776f86d14"
	SettingInstanceTemplateDepartmentPhone         = "383a4d4b-0bbc-4791-bc48-6def8716ec96"
	SettingInstanceTemplateLocationServices        = "5dca1893-31cd-49b7-aadf-c26990ea0639"
	SettingInstanceTemplateRestore                 = "b4434e35-4131-4b13-90cb-4334cd9768dd"
	SettingInstanceTemplateAppleId                 = "4f26cd1f-5fd1-4d6c-af69-f18ade85aee9"
	SettingInstanceTemplateTermsAndConditions      = "5a33bd8e-9df4-422c-92f5-3d5294a3e015"
	SettingInstanceTemplateTouchFaceId             = "0deedc8e-1597-47f0-a271-4ae5029c1795"
	SettingInstanceTemplateApplePay                = "c7b005a7-148a-40b0-aa06-8a6eaafa6fc0"
	SettingInstanceTemplateSiri                    = "764673c0-b331-478e-905f-ef45428fae99"
	SettingInstanceTemplateDiagnosticsData         = "6a6c0b44-8403-4213-8e34-fd21e416c0f3"
	SettingInstanceTemplateFileVault               = "7d8e9c44-3b2a-4f13-8a5e-2c9d1b4f5a3e"
	SettingInstanceTemplateICloudDiagnostics       = "4c8a7d9e-2f5b-4a3c-9d8e-1b6a4f7c2e9d"
	SettingInstanceTemplateICloudStorage           = "8f7d3c9a-4e2b-4d1a-9c8e-6b5a3d7f2e4c"
	SettingInstanceTemplateAppearance              = "61b602c2-86a3-4bca-82e5-cc7e14e17895"
	SettingInstanceTemplateScreenTime              = "90501450-54d9-4695-a16c-fbfead918c24"
	SettingInstanceTemplatePrivacy                 = "129ac5ce-c6b8-4402-bf18-40045c11c10d"
	SettingInstanceTemplateAccessibility           = "9c7e4d8a-3f2b-4e1a-8d9c-5b7a3e6f2c4d"
	SettingInstanceTemplateUnlockWithWatch         = "5c8e9d7a-4f3b-4d2a-8e9c-7b6a5d3f2e1c"
	SettingInstanceTemplateEnableLockdownMode      = "4a9db745-821d-44a1-a26a-06f2c637d8b8"
	SettingInstanceTemplateSoftwareUpdate          = "1ad911de-03ea-45a1-917e-adf9ececc09c"
	SettingInstanceTemplateSoftwareUpdateCompleted = "8afaa22c-513a-41b4-9b49-1c888a6f1c53"
	SettingInstanceTemplateTermsOfAddress          = "21017ab9-73d3-45dd-b715-8742791a5d73"
	SettingInstanceTemplateIntelligence            = "cca42541-0513-41a6-ae74-e1fe878b0593"
	SettingInstanceTemplateOSShowcase              = "5d9c8e7a-3f2b-4d1a-8e9c-6b5a7d3f2e4c"
	SettingInstanceTemplateAppStore                = "4c8e7d9a-3f2b-4d1a-9c8e-5b7a6d3f2e1c"

	// Setting Value Template IDs
	SettingValueTemplateUserAffinity            = "e34ab06d-486a-4cc0-a867-295515ac2dfc"
	SettingValueTemplateAwaitConfiguration      = "3a7c9d2e-8f4b-4e1a-9c5d-6b8a3f7e2d4c"
	SettingValueTemplateLockedEnrollment        = "49085612-4cce-4305-822d-6a6c793285be"
	SettingValueTemplateDepartment              = "180721e2-b467-40a3-b631-b2f28c5d9297"
	SettingValueTemplateDepartmentPhone         = "f7322d5b-07aa-4b5c-9846-bc3e547e7c0b"
	SettingValueTemplateLocationServices        = "6f73700e-a14b-4784-9923-512a97297c41"
	SettingValueTemplateRestore                 = "78842c1b-d090-407d-9dec-a3a4a168a36f"
	SettingValueTemplateAppleId                 = "469e5c1a-c2f8-44fb-b515-f8532c2dff40"
	SettingValueTemplateTermsAndConditions      = "404a8266-5c46-44b3-bf63-d8974d91c8f9"
	SettingValueTemplateTouchFaceId             = "f11533e6-85a1-41c3-b23d-7a7c0bc2f7f8"
	SettingValueTemplateApplePay                = "a18db77d-76e6-4bf8-ab60-98d1e34ed8df"
	SettingValueTemplateSiri                    = "f3c03777-e89d-4fbc-bc1a-75c4858fa579"
	SettingValueTemplateDiagnosticsData         = "8e6986f8-3f2b-4a22-989d-08603dba0707"
	SettingValueTemplateFileVault               = "2f5c47e9-aa1b-4c35-9e8d-1fa2b8c3d7ea"
	SettingValueTemplateICloudDiagnostics       = "9a7d5c21-4e3b-4f8a-a9c6-7b5d3e1f2c4a"
	SettingValueTemplateICloudStorage           = "5b3e8a7c-9d2f-4a1c-8e5b-3d7a2f9c4e1b"
	SettingValueTemplateAppearance              = "87067516-5437-44ce-a3cc-cee7b6518b46"
	SettingValueTemplateScreenTime              = "b1d1feea-8d71-49bb-8cf6-c736321ae9b6"
	SettingValueTemplatePrivacy                 = "177b8bd8-8dc1-4835-82af-16e81dc7e94f"
	SettingValueTemplateAccessibility           = "3f8d7c9b-5e2a-4d1c-8f7e-2b9a5c6d3e1f"
	SettingValueTemplateUnlockWithWatch         = "7e9d4c8a-2f5b-4e3c-9d8a-6b7c3e5f2a1d"
	SettingValueTemplateEnableLockdownMode      = "181a5e6d-df5c-4016-b484-3130d54ee9ab"
	SettingValueTemplateSoftwareUpdate          = "d9630ff4-94a4-4525-b37a-7c233a23a2bd"
	SettingValueTemplateSoftwareUpdateCompleted = "80bdf240-ad52-45fb-9a58-c888d812e75a"
	SettingValueTemplateTermsOfAddress          = "df5a331d-7042-4abc-91f4-a37c82deaf08"
	SettingValueTemplateIntelligence            = "d0c43d95-99c7-40ce-b3a6-33098215adc9"
	SettingValueTemplateOSShowcase              = "8c2d7f9e-4a3b-4e1c-9d8e-5b7a3c6f2d1e"
	SettingValueTemplateAppStore                = "7e8d9c4a-2f5b-4e3c-8d9a-6b7c3e5f2a1d"

	// SettingValueTemplateAccountSettings is the shared settingValueTemplateId reused by every
	// setting nested under ade_macos_awaitconfiguration's account-settings subtree (createlocaladmin,
	// adminaccountname, adminaccountfullname, hideusersgroups, adminaccountpasswordrotation,
	// createlocalprimary, prefillaccountinfo). ade_accountsettings_restrictediting and
	// ade_macos_authenticationmethod are the only nested settings observed with NO template
	// references at all.
	SettingValueTemplateAccountSettings = "3a7c9d2e-8f4b-4e1a-9c5d-6b8a3f7e2d4c"

	// Authentication method values for ade_macos_authenticationmethod (only observed value is _2;
	// _0/_1 are inferred by analogy with the mutually exclusive company-portal booleans on the
	// legacy macos_dep_enrollment_profile resource and have not been independently verified against
	// Graph).
	AuthenticationMethodBasic                         = SettingDefAuthenticationMethod + "_0"
	AuthenticationMethodCompanyPortal                 = SettingDefAuthenticationMethod + "_1"
	AuthenticationMethodCompanyPortalOnSetupAssistant = SettingDefAuthenticationMethod + "_2"
)
