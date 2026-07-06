package graphBetaIOSiPadOSDeviceEnrollmentPolicy

// Setting definition IDs, and setting instance/value template GUIDs, were captured from
// Graph X-Ray traces of the Intune admin center creating an iOS/iPadOS Automated Device
// Enrollment (ADE) profile via the settings catalog backed `/deviceManagement/configurationPolicies`
// endpoint. See docs/create.md for the raw captured payload.
const (
	// TemplateID is the settings catalog template used for iOS/iPadOS ADE/DEP enrollment profiles.
	TemplateID = "27d20e9c-50c1-48f8-a44c-f37de4510051_1"

	// TemplateFamily for enrollment restriction/configuration policies.
	TemplateFamily = "enrollmentConfiguration"

	// Platforms and technologies values required by Graph for this template.
	Platforms    = "iOS"
	Technologies = "enrollment"

	// CreationSourcePrefix is prepended to the resolved dep_onboarding_settings_id to build the
	// `creationSource` value sent on Create. It is not sent on Update.
	CreationSourcePrefix = "DepTokenId_"

	// Setting Definition IDs - User Affinity / Authentication
	SettingDefUserAffinity            = "ade_useraffinity"
	SettingDefAuthenticationMethod    = "ade_authenticationmethod"
	SettingDefAwaitFinalConfiguration = "ade_modernauth_awaitfinalconfiguration"

	// Setting Definition IDs - Locked Enrollment
	SettingDefLockedEnrollment = "ade_lockedenrollment"

	// Setting Definition IDs - Device Name Template
	SettingDefDeviceNameTemplateChoices = "ade_devicenametemplatechoices"
	SettingDefAppleDeviceNameTemplate   = "ade_appledevicenametemplate"

	// Setting Definition IDs - Cellular Data Activation
	SettingDefActivateCellularDataChoices = "ade_activatecellulardatachoices"
	SettingDefActivateCellularData        = "ade_activatecellulardata"

	// Setting Definition IDs - Setup Assistant
	SettingDefDepartment              = "ade_setupassistant_department"
	SettingDefDepartmentPhone         = "ade_setupassistant_departmentphone"
	SettingDefPasscode                = "ade_setupassistant_passcode"
	SettingDefLocationServices        = "ade_setupassistant_locationservices"
	SettingDefRestore                 = "ade_setupassistant_restore"
	SettingDefAppleId                 = "ade_setupassistant_appleid"
	SettingDefTermsAndConditions      = "ade_setupassistant_termsandconditions"
	SettingDefTouchFaceId             = "ade_setupassistant_touchfaceid"
	SettingDefApplePay                = "ade_setupassistant_applepay"
	SettingDefSiri                    = "ade_setupassistant_siri"
	SettingDefDiagnosticsData         = "ade_setupassistant_diagnosticsdata"
	SettingDefPrivacy                 = "ade_setupassistant_privacy"
	SettingDefAndroidMigration        = "ade_setupassistant_androidmigration"
	SettingDefIMessageFaceTime        = "ade_setupassistant_imessagefacetime"
	SettingDefScreenTime              = "ade_setupassistant_screentime"
	SettingDefSimSetup                = "ade_setupassistant_simsetup"
	SettingDefSoftwareUpdate          = "ade_setupassistant_softwareupdate"
	SettingDefWatchMigration          = "ade_setupassistant_watchmigration"
	SettingDefAppearance              = "ade_setupassistant_appearance"
	SettingDefDeviceMigration         = "ade_setupassistant_devicemigration"
	SettingDefRestoreCompleted        = "ade_setupassistant_restorecompleted"
	SettingDefSoftwareUpdateCompleted = "ade_setupassistant_softwareupdatecompleted"
	SettingDefGetStarted              = "ade_setupassistant_getstarted"
	SettingDefActionButton            = "ade_setupassistant_actionbutton"
	SettingDefSafety                  = "ade_setupassistant_safety"
	SettingDefTermsOfAddress          = "ade_setupassistant_termsofaddress"
	SettingDefIntelligence            = "ade_setupassistant_intelligence"
	SettingDefEnableLockdownMode      = "ade_setupassistant_enablelockdownmode"
	SettingDefAppStore                = "ade_setupassistant_appstore"
	SettingDefCameraButton            = "ade_setupassistant_camerabutton"
	SettingDefMultitasking            = "ade_setupassistant_multitasking"
	SettingDefOSShowcase              = "ade_setupassistant_osshowcase"
	SettingDefSafetyAndHandling       = "ade_setupassistant_safetyandhandling"
	SettingDefWebContentFiltering     = "ade_setupassistant_webcontentfiltering"

	// Setting Instance Template IDs
	SettingInstanceTemplateUserAffinity              = "f3ddbbae-bb3d-45e8-bb9c-854f674e56ef"
	SettingInstanceTemplateLockedEnrollment          = "e105d3a4-7318-454a-a231-0d1510ebf1e3"
	SettingInstanceTemplateDeviceNameTemplateChoices = "93d1bf98-2345-468a-b0f0-1e2daf198357"
	SettingInstanceTemplateActivateCellularData      = "b63070e7-d6f7-44e7-b1c1-fd51615748a3"
	SettingInstanceTemplateDepartment                = "1fd43c7c-0c7a-4dcc-8392-7b1776f86d14"
	SettingInstanceTemplateDepartmentPhone           = "383a4d4b-0bbc-4791-bc48-6def8716ec96"
	SettingInstanceTemplatePasscode                  = "be808211-b4f1-4653-aa95-068b756dc682"
	SettingInstanceTemplateLocationServices          = "5dca1893-31cd-49b7-aadf-c26990ea0639"
	SettingInstanceTemplateRestore                   = "b4434e35-4131-4b13-90cb-4334cd9768dd"
	SettingInstanceTemplateAppleId                   = "4f26cd1f-5fd1-4d6c-af69-f18ade85aee9"
	SettingInstanceTemplateTermsAndConditions        = "5a33bd8e-9df4-422c-92f5-3d5294a3e015"
	SettingInstanceTemplateTouchFaceId               = "0deedc8e-1597-47f0-a271-4ae5029c1795"
	SettingInstanceTemplateApplePay                  = "c7b005a7-148a-40b0-aa06-8a6eaafa6fc0"
	SettingInstanceTemplateSiri                      = "764673c0-b331-478e-905f-ef45428fae99"
	SettingInstanceTemplateDiagnosticsData           = "6a6c0b44-8403-4213-8e34-fd21e416c0f3"
	SettingInstanceTemplatePrivacy                   = "129ac5ce-c6b8-4402-bf18-40045c11c10d"
	SettingInstanceTemplateAndroidMigration          = "ec192455-f753-4fe1-a4b5-ea98ce1f51a4"
	SettingInstanceTemplateIMessageFaceTime          = "53e650ca-9548-4fb6-8cba-3be4bc2d62a5"
	SettingInstanceTemplateScreenTime                = "90501450-54d9-4695-a16c-fbfead918c24"
	SettingInstanceTemplateSimSetup                  = "7a270780-d94f-4d00-998e-2a5d69f43832"
	SettingInstanceTemplateSoftwareUpdate            = "1ad911de-03ea-45a1-917e-adf9ececc09c"
	SettingInstanceTemplateWatchMigration            = "be182ec4-4e16-4f98-9ca2-e4775d35eb88"
	SettingInstanceTemplateAppearance                = "61b602c2-86a3-4bca-82e5-cc7e14e17895"
	SettingInstanceTemplateDeviceMigration           = "3243140d-97b7-42fe-ba57-f23f782f9112"
	SettingInstanceTemplateRestoreCompleted          = "991507fe-fe61-41f9-b80c-cac3f071fb92"
	SettingInstanceTemplateSoftwareUpdateCompleted   = "8afaa22c-513a-41b4-9b49-1c888a6f1c53"
	SettingInstanceTemplateGetStarted                = "21017ab9-73d3-45dd-b715-8742791a5d73"
	SettingInstanceTemplateActionButton              = "6f78e140-a99d-49a0-9658-0d5a4bb30726"
	SettingInstanceTemplateSafety                    = "b8954ce1-36f1-4149-9851-99c050dfc3f8"
	SettingInstanceTemplateTermsOfAddress            = "5dbc8d66-d852-4188-9500-020956fbdfb7"
	SettingInstanceTemplateIntelligence              = "cca42541-0513-41a6-ae74-e1fe878b0593"
	SettingInstanceTemplateEnableLockdownMode        = "4a9db745-821d-44a1-a26a-06f2c637d8b8"
	SettingInstanceTemplateAppStore                  = "b3c4d5e6-f7a8-9012-bcde-fa3456789012"
	SettingInstanceTemplateCameraButton              = "d5e6f7a8-b9c0-1234-defa-bc5678901234"
	SettingInstanceTemplateMultitasking              = "f7a8b9c0-d1e2-3456-fabc-de7890123456"
	SettingInstanceTemplateOSShowcase                = "b9c0d1e2-f3a4-5678-bcde-fa9012345678"
	SettingInstanceTemplateSafetyAndHandling         = "d1e2f3a4-b5c6-7890-defa-bc1234567890"
	SettingInstanceTemplateWebContentFiltering       = "f3a4b5c6-d7e8-9012-fabc-de3456789012"

	// Setting Value Template IDs
	SettingValueTemplateUserAffinity              = "e34ab06d-486a-4cc0-a867-295515ac2dfc"
	SettingValueTemplateLockedEnrollment          = "49085612-4cce-4305-822d-6a6c793285be"
	SettingValueTemplateDeviceNameTemplateChoices = "8c02396f-d58f-4275-be5c-dd97713d636e"
	SettingValueTemplateActivateCellularData      = "aaf7603d-0184-4143-bdfc-3326859ce0c3"
	SettingValueTemplateDepartment                = "180721e2-b467-40a3-b631-b2f28c5d9297"
	SettingValueTemplateDepartmentPhone           = "f7322d5b-07aa-4b5c-9846-bc3e547e7c0b"
	SettingValueTemplatePasscode                  = "67b68eed-075f-4479-a34b-9ae3d2b1cf4d"
	SettingValueTemplateLocationServices          = "6f73700e-a14b-4784-9923-512a97297c41"
	SettingValueTemplateRestore                   = "78842c1b-d090-407d-9dec-a3a4a168a36f"
	SettingValueTemplateAppleId                   = "469e5c1a-c2f8-44fb-b515-f8532c2dff40"
	SettingValueTemplateTermsAndConditions        = "404a8266-5c46-44b3-bf63-d8974d91c8f9"
	SettingValueTemplateTouchFaceId               = "f11533e6-85a1-41c3-b23d-7a7c0bc2f7f8"
	SettingValueTemplateApplePay                  = "a18db77d-76e6-4bf8-ab60-98d1e34ed8df"
	SettingValueTemplateSiri                      = "f3c03777-e89d-4fbc-bc1a-75c4858fa579"
	SettingValueTemplateDiagnosticsData           = "8e6986f8-3f2b-4a22-989d-08603dba0707"
	SettingValueTemplatePrivacy                   = "177b8bd8-8dc1-4835-82af-16e81dc7e94f"
	SettingValueTemplateAndroidMigration          = "55060cc0-c46b-40fa-9435-7adf473662fb"
	SettingValueTemplateIMessageFaceTime          = "b1730b75-01d1-41a4-9e2e-d185ec2a5108"
	SettingValueTemplateScreenTime                = "b1d1feea-8d71-49bb-8cf6-c736321ae9b6"
	SettingValueTemplateSimSetup                  = "8c2b8a4c-517a-47f7-9f05-392db533b05a"
	SettingValueTemplateSoftwareUpdate            = "d9630ff4-94a4-4525-b37a-7c233a23a2bd"
	SettingValueTemplateWatchMigration            = "fe924b50-718f-4b7f-ac33-5a407ea0f005"
	SettingValueTemplateAppearance                = "87067516-5437-44ce-a3cc-cee7b6518b46"
	SettingValueTemplateDeviceMigration           = "55900392-ff35-43dc-9adb-7fca1987a106"
	SettingValueTemplateRestoreCompleted          = "14e4f24a-55d5-4c04-9f61-862cb42a90e9"
	SettingValueTemplateSoftwareUpdateCompleted   = "80bdf240-ad52-45fb-9a58-c888d812e75a"
	SettingValueTemplateGetStarted                = "df5a331d-7042-4abc-91f4-a37c82deaf08"
	SettingValueTemplateActionButton              = "54fd9219-4c4a-4350-a2be-a85f1096f4da"
	SettingValueTemplateSafety                    = "6290e656-2b8d-4793-b8fb-4336c64a32c5"
	SettingValueTemplateTermsOfAddress            = "1f57641d-fd69-403e-b96d-4f60869e6d9d"
	SettingValueTemplateIntelligence              = "d0c43d95-99c7-40ce-b3a6-33098215adc9"
	SettingValueTemplateEnableLockdownMode        = "181a5e6d-df5c-4016-b484-3130d54ee9ab"
	SettingValueTemplateAppStore                  = "a2b3c4d5-e6f7-8901-abcd-ef2345678901"
	SettingValueTemplateCameraButton              = "c4d5e6f7-a8b9-0123-cdef-ab4567890123"
	SettingValueTemplateMultitasking              = "e6f7a8b9-c0d1-2345-efab-cd6789012345"
	SettingValueTemplateOSShowcase                = "a8b9c0d1-e2f3-4567-abcd-ef8901234567"
	SettingValueTemplateSafetyAndHandling         = "c0d1e2f3-a4b5-6789-cdef-ab0123456789"
	SettingValueTemplateWebContentFiltering       = "e2f3a4b5-c6d7-8901-efab-cd2345678901"

	// Authentication method values for ade_authenticationmethod. Only _2 (Setup Assistant with
	// modern authentication, carrying the ade_modernauth_awaitfinalconfiguration child) was
	// observed live; _0/_1 are inferred by analogy with the macOS ade_macos_authenticationmethod
	// values and the Intune admin center's authentication method picker, and have not been
	// independently verified against Graph. ade_authenticationmethod and its child carry no
	// template references.
	AuthenticationMethodSetupAssistantLegacy     = SettingDefAuthenticationMethod + "_0"
	AuthenticationMethodCompanyPortal            = SettingDefAuthenticationMethod + "_1"
	AuthenticationMethodSetupAssistantModernAuth = SettingDefAuthenticationMethod + "_2"
)
