package graphBetaWindowsAutopilotDevicePreparationPolicy

const (
	// Deployment type constants (used to determine template and settings)
	DeploymentTypeUserDriven    = "enrollment_autopilot_dpp_deploymenttype_0"
	DeploymentTypeSelfDeploying = "enrollment_autopilot_dpp_deploymenttype_1"

	// Template IDs for different deployment types
	// Self-deploying/automatic mode uses the automatic template
	TemplateIDAutomatic = "a6157a7f-aa00-42d9-ac82-7d2479f545db_1"
	// User-driven mode uses the user-driven template
	TemplateIDUserDriven = "80d33118-b7b4-40d8-b15f-81be745e053f_1"

	// Template family
	TemplateFamily = "enrollmentConfiguration"

	// Setting Definition IDs
	SettingDefDeploymentMode         = "enrollment_autopilot_dpp_deploymentmode"
	SettingDefDeploymentType         = "enrollment_autopilot_dpp_deploymenttype"
	SettingDefJoinType               = "enrollment_autopilot_dpp_jointype"
	SettingDefAccountType            = "enrollment_autopilot_dpp_accountype"
	SettingDefTimeout                = "enrollment_autopilot_dpp_timeout"
	SettingDefCustomErrorMessage     = "enrollment_autopilot_dpp_customerrormessage"
	SettingDefAllowSkip              = "enrollment_autopilot_dpp_allowskip"
	SettingDefAllowDiagnostics       = "enrollment_autopilot_dpp_allowdiagnostics"
	SettingDefDeviceSecurityGroupIDs = "enrollment_autopilot_dpp_devicesecuritygroupids"
	SettingDefAllowedAppIDs          = "enrollment_autopilot_dpp_allowedappids"
	SettingDefAllowedScriptIDs       = "enrollment_autopilot_dpp_allowedscriptids"

	// Setting Instance Template IDs - Deployment Settings
	SettingInstanceTemplateDeploymentMode = "5180aeab-886e-4589-97d4-40855c646315"
	SettingInstanceTemplateDeploymentType = "f4184296-fa9f-4b67-8b12-1723b3f8456b"
	SettingInstanceTemplateJoinType       = "6310e95d-6cfa-4d2f-aae0-1e7af12e2182"
	SettingInstanceTemplateAccountType    = "d4f2a840-86d5-4162-9a08-fa8cc608b94e"

	// Setting Instance Template IDs - OOBE Settings
	SettingInstanceTemplateTimeout             = "6dec0657-dfb8-4906-a7ee-3ac6ee1edecb"
	SettingInstanceTemplateCustomErrorMessage  = "2ddf0619-2b7a-46de-b29b-c6191e9dda6e"
	SettingInstanceTemplateAllowSkip           = "2a71dc89-0f17-4ba9-bb27-af2521d34710"
	SettingInstanceTemplateAllowDiagnostics    = "e2b7a81b-f243-4abd-bce3-c1856345f405"
	SettingInstanceTemplateDeviceSecurityGroup = "a46a50ab-3076-4968-9366-75a40dde950e"

	// Setting Instance Template IDs - Apps (Automatic Mode)
	SettingInstanceTemplateAllowedAppsAutomatic = "a9dedfd6-c3b2-46d9-ae39-91fd0dcb7a20"

	// Setting Instance Template IDs - Apps (User-Driven Mode)
	SettingInstanceTemplateAllowedAppsUserDriven = "70d22a8a-a03c-4f62-b8df-dded3e327639"

	// Setting Instance Template IDs - Scripts (Automatic Mode)
	SettingInstanceTemplateAllowedScriptsAutomatic = "ff20a4a9-a2f4-4a2e-84e0-4cd1dc9bed31"

	// Setting Instance Template IDs - Scripts (User-Driven Mode)
	SettingInstanceTemplateAllowedScriptsUserDriven = "1bc67702-800c-4271-8fd9-609351cc19cf"

	// Setting Value Template IDs - Deployment Settings
	SettingValueTemplateDeploymentMode = "5874c2f6-bcf1-463b-a9eb-bee64e2f2d82"
	SettingValueTemplateDeploymentType = "e0af022f-37f3-4a40-916d-1ab7281c88d9"
	SettingValueTemplateJoinType       = "1fa84eb3-fcfa-4ed6-9687-0f3d486402c4"
	SettingValueTemplateAccountType    = "bf13bb47-69ef-4e06-97c1-50c2859a49c2"

	// Setting Value Template IDs - OOBE Settings
	SettingValueTemplateTimeout             = "0bbcce5b-a55a-4e05-821a-94bf576d6cc8"
	SettingValueTemplateCustomErrorMessage  = "fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97"
	SettingValueTemplateAllowSkip           = "a2323e5e-ac56-4517-8847-b0a6fdb467e7"
	SettingValueTemplateAllowDiagnostics    = "c59d26fd-3460-4b26-b47a-f7e202e7d5a3"
	SettingValueTemplateDeviceSecurityGroup = "5f7d09e1-1a90-44ad-9c9f-ad90ba509e60"
)
