package graphBetaVisionOSDeviceEnrollmentPolicy

// Setting definition IDs, and setting instance/value template GUIDs, were captured from
// Graph X-Ray traces of the Intune admin center creating a visionOS Automated Device
// Enrollment (ADE) profile via the settings catalog backed `/deviceManagement/configurationPolicies`
// endpoint. See docs/create.md for the raw captured payload.
const (
	// TemplateID is the settings catalog template used for visionOS ADE/DEP enrollment profiles.
	TemplateID = "b974d292-62c4-4853-bdab-166cf42df51c_1"

	// TemplateFamily for enrollment restriction/configuration policies.
	TemplateFamily = "enrollmentConfiguration"

	// Platforms and technologies values required by Graph for this template.
	Platforms    = "visionOS"
	Technologies = "enrollment"

	// CreationSourcePrefix is prepended to the resolved dep_onboarding_settings_id to build the
	// `creationSource` value sent on Create. It is not sent on Update.
	CreationSourcePrefix = "DepTokenId_"

	// Setting Definition IDs - User Affinity / Await Configuration. visionOS uses the "basic"
	// variants of both settings: unlike iOS/iPadOS and macOS, neither carries nested children
	// (no authentication method choice, no account settings subtree).
	SettingDefUserAffinity       = "ade_useraffinitybasic"
	SettingDefAwaitConfiguration = "ade_awaitconfiguration_basic"

	// Setting Definition IDs - Locked Enrollment
	SettingDefLockedEnrollment = "ade_lockedenrollment"

	// Setting Definition IDs - Setup Assistant
	SettingDefDepartment         = "ade_setupassistant_department"
	SettingDefDepartmentPhone    = "ade_setupassistant_departmentphone"
	SettingDefAppleId            = "ade_setupassistant_appleid"
	SettingDefApplePay           = "ade_setupassistant_applepay"
	SettingDefDiagnosticsData    = "ade_setupassistant_diagnosticsdata"
	SettingDefGetStarted         = "ade_setupassistant_getstarted"
	SettingDefIntelligence       = "ade_setupassistant_intelligence"
	SettingDefLocationServices   = "ade_setupassistant_locationservices"
	SettingDefPasscode           = "ade_setupassistant_passcode"
	SettingDefPrivacy            = "ade_setupassistant_privacy"
	SettingDefScreenTime         = "ade_setupassistant_screentime"
	SettingDefSiri               = "ade_setupassistant_siri"
	SettingDefSoftwareUpdate     = "ade_setupassistant_softwareupdate"
	SettingDefTermsAndConditions = "ade_setupassistant_termsandconditions"
	SettingDefTips               = "ade_setupassistant_tips"
	SettingDefTouchFaceId        = "ade_setupassistant_touchfaceid"

	// Setting Instance Template IDs. ade_lockedenrollment, ade_setupassistant_department and
	// ade_setupassistant_departmentphone carry visionOS-specific template GUIDs that differ from
	// the iOS/iPadOS and macOS templates; the shared Setup Assistant panes reuse the same GUIDs.
	SettingInstanceTemplateUserAffinity       = "24fae708-81c3-450b-bcb0-711f00b0f78e"
	SettingInstanceTemplateAwaitConfiguration = "78bb9ed5-6bae-4e30-a258-c50eefe1afc6"
	SettingInstanceTemplateLockedEnrollment   = "b4df4b6a-b8d7-4adc-a4cf-15a3d2dd8018"
	SettingInstanceTemplateDepartment         = "96c5e30b-2d83-4055-b04a-b6f4a6af3568"
	SettingInstanceTemplateDepartmentPhone    = "d5ba8ce3-91ee-409d-bbbf-10640a197d21"
	SettingInstanceTemplateAppleId            = "4f26cd1f-5fd1-4d6c-af69-f18ade85aee9"
	SettingInstanceTemplateApplePay           = "c7b005a7-148a-40b0-aa06-8a6eaafa6fc0"
	SettingInstanceTemplateDiagnosticsData    = "6a6c0b44-8403-4213-8e34-fd21e416c0f3"
	SettingInstanceTemplateGetStarted         = "21017ab9-73d3-45dd-b715-8742791a5d73"
	SettingInstanceTemplateIntelligence       = "cca42541-0513-41a6-ae74-e1fe878b0593"
	SettingInstanceTemplateLocationServices   = "5dca1893-31cd-49b7-aadf-c26990ea0639"
	SettingInstanceTemplatePasscode           = "be808211-b4f1-4653-aa95-068b756dc682"
	SettingInstanceTemplatePrivacy            = "129ac5ce-c6b8-4402-bf18-40045c11c10d"
	SettingInstanceTemplateScreenTime         = "90501450-54d9-4695-a16c-fbfead918c24"
	SettingInstanceTemplateSiri               = "764673c0-b331-478e-905f-ef45428fae99"
	SettingInstanceTemplateSoftwareUpdate     = "1ad911de-03ea-45a1-917e-adf9ececc09c"
	SettingInstanceTemplateTermsAndConditions = "5a33bd8e-9df4-422c-92f5-3d5294a3e015"
	SettingInstanceTemplateTips               = "b00de12d-5ae8-4b76-8763-c2230ff1018b"
	SettingInstanceTemplateTouchFaceId        = "0deedc8e-1597-47f0-a271-4ae5029c1795"

	// Setting Value Template IDs
	SettingValueTemplateUserAffinity       = "ad35a84e-5803-412a-b590-792f30a77438"
	SettingValueTemplateAwaitConfiguration = "76430454-b68e-4685-bf56-4bffa1b0c696"
	SettingValueTemplateLockedEnrollment   = "8c73f868-3afd-4b99-b978-ca3dc8a773dc"
	SettingValueTemplateDepartment         = "ece47992-26c1-4377-ba25-07136be7f6a1"
	SettingValueTemplateDepartmentPhone    = "b35b460b-2bd6-4dc1-bd01-97fcff821054"
	SettingValueTemplateAppleId            = "469e5c1a-c2f8-44fb-b515-f8532c2dff40"
	SettingValueTemplateApplePay           = "a18db77d-76e6-4bf8-ab60-98d1e34ed8df"
	SettingValueTemplateDiagnosticsData    = "8e6986f8-3f2b-4a22-989d-08603dba0707"
	SettingValueTemplateGetStarted         = "df5a331d-7042-4abc-91f4-a37c82deaf08"
	SettingValueTemplateIntelligence       = "d0c43d95-99c7-40ce-b3a6-33098215adc9"
	SettingValueTemplateLocationServices   = "6f73700e-a14b-4784-9923-512a97297c41"
	SettingValueTemplatePasscode           = "67b68eed-075f-4479-a34b-9ae3d2b1cf4d"
	SettingValueTemplatePrivacy            = "177b8bd8-8dc1-4835-82af-16e81dc7e94f"
	SettingValueTemplateScreenTime         = "b1d1feea-8d71-49bb-8cf6-c736321ae9b6"
	SettingValueTemplateSiri               = "f3c03777-e89d-4fbc-bc1a-75c4858fa579"
	SettingValueTemplateSoftwareUpdate     = "d9630ff4-94a4-4525-b37a-7c233a23a2bd"
	SettingValueTemplateTermsAndConditions = "404a8266-5c46-44b3-bf63-d8974d91c8f9"
	SettingValueTemplateTips               = "7d55896a-705c-4a75-b226-f90ffc9d5524"
	SettingValueTemplateTouchFaceId        = "f11533e6-85a1-41c3-b23d-7a7c0bc2f7f8"
)
