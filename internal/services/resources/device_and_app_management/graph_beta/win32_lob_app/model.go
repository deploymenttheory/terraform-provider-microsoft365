// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapp?view=graph-rest-beta
package graphBetaWin32LobApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Win32LobAppResourceModel struct {
	ID                              types.String                                     `tfsdk:"id"`
	DisplayName                     types.String                                     `tfsdk:"display_name"`
	Description                     types.String                                     `tfsdk:"description"`
	Publisher                       types.String                                     `tfsdk:"publisher"`
	LargeIcon                       sharedmodels.MimeContentResourceModel            `tfsdk:"large_icon"`
	CreatedDateTime                 types.String                                     `tfsdk:"created_date_time"`
	LastModifiedDateTime            types.String                                     `tfsdk:"last_modified_date_time"`
	IsFeatured                      types.Bool                                       `tfsdk:"is_featured"`
	PrivacyInformationUrl           types.String                                     `tfsdk:"privacy_information_url"`
	InformationUrl                  types.String                                     `tfsdk:"information_url"`
	Owner                           types.String                                     `tfsdk:"owner"`
	Developer                       types.String                                     `tfsdk:"developer"`
	Notes                           types.String                                     `tfsdk:"notes"`
	UploadState                     types.Int32                                      `tfsdk:"upload_state"`
	PublishingState                 types.String                                     `tfsdk:"publishing_state"`
	IsAssigned                      types.Bool                                       `tfsdk:"is_assigned"`
	RoleScopeTagIds                 types.Set                                        `tfsdk:"role_scope_tag_ids"`
	DependentAppCount               types.Int32                                      `tfsdk:"dependent_app_count"`
	SupersedingAppCount             types.Int32                                      `tfsdk:"superseding_app_count"`
	SupersededAppCount              types.Int32                                      `tfsdk:"superseded_app_count"`
	CommittedContentVersion         types.String                                     `tfsdk:"committed_content_version"`
	FileName                        types.String                                     `tfsdk:"file_name"`
	Size                            types.Int64                                      `tfsdk:"size"`
	InstallCommandLine              types.String                                     `tfsdk:"install_command_line"`
	UninstallCommandLine            types.String                                     `tfsdk:"uninstall_command_line"`
	ApplicableArchitectures         types.String                                     `tfsdk:"applicable_architectures"`
	MinimumSupportedOperatingSystem WindowsMinimumOperatingSystemResourceModel       `tfsdk:"minimum_supported_operating_system"`
	MinimumFreeDiskSpaceInMB        types.Int32                                      `tfsdk:"minimum_free_disk_space_in_mb"`
	MinimumMemoryInMB               types.Int32                                      `tfsdk:"minimum_memory_in_mb"`
	MinimumNumberOfProcessors       types.Int32                                      `tfsdk:"minimum_number_of_processors"`
	MinimumCpuSpeedInMHz            types.Int32                                      `tfsdk:"minimum_cpu_speed_in_mhz"`
	DetectionRules                  []Win32LobAppRegistryDetectionRulesResourceModel `tfsdk:"detection_rules"`
	RequirementRules                []Win32LobAppRegistryRequirementResourceModel    `tfsdk:"requirement_rules"`
	Rules                           []Win32LobAppRegistryRuleResourceModel           `tfsdk:"rules"`
	InstallExperience               Win32LobAppInstallExperienceResourceModel        `tfsdk:"install_experience"`
	ReturnCodes                     []Win32LobAppReturnCodeResourceModel             `tfsdk:"return_codes"`
	MsiInformation                  Win32LobAppMsiInformationResourceModel           `tfsdk:"msi_information"`
	SetupFilePath                   types.String                                     `tfsdk:"setup_file_path"`
	MinimumSupportedWindowsRelease  types.String                                     `tfsdk:"minimum_supported_windows_release"`
	DisplayVersion                  types.String                                     `tfsdk:"display_version"`
	AllowAvailableUninstall         types.Bool                                       `tfsdk:"allow_available_uninstall"`
	Timeouts                        timeouts.Value                                   `tfsdk:"timeouts"`
}

type WindowsMinimumOperatingSystemResourceModel struct {
	V8_0     types.Bool `tfsdk:"v8_0"`
	V8_1     types.Bool `tfsdk:"v8_1"`
	V10_0    types.Bool `tfsdk:"v10_0"`
	V10_1607 types.Bool `tfsdk:"v10_1607"`
	V10_1703 types.Bool `tfsdk:"v10_1703"`
	V10_1709 types.Bool `tfsdk:"v10_1709"`
	V10_1803 types.Bool `tfsdk:"v10_1803"`
	V10_1809 types.Bool `tfsdk:"v10_1809"`
	V10_1903 types.Bool `tfsdk:"v10_1903"`
	V10_1909 types.Bool `tfsdk:"v10_1909"`
	V10_2004 types.Bool `tfsdk:"v10_2004"`
	V10_2H20 types.Bool `tfsdk:"v10_2h20"`
	V10_21H1 types.Bool `tfsdk:"v10_21h1"`
}

type Win32LobAppRegistryDetectionRulesResourceModel struct {
	// Common for multiple detection types
	DetectionType        types.String `tfsdk:"detection_type"`            // registry, msi_information, file_system, powershell_script
	Check32BitOn64System types.Bool   `tfsdk:"check_32_bit_on_64_system"` // Only for registry, file_system, powershell_script
	DetectionValue       types.String `tfsdk:"detection_value"`           // For registry and file_system detection types
	// Registry-specific fields
	RegistryDetectionType     types.String `tfsdk:"registry_detection_type"`
	KeyPath                   types.String `tfsdk:"key_path"`
	ValueName                 types.String `tfsdk:"value_name"`
	RegistryDetectionOperator types.String `tfsdk:"registry_detection_operator"`
	// MSI-specific fields
	ProductCode            types.String `tfsdk:"product_code"`
	ProductVersion         types.String `tfsdk:"product_version"`
	ProductVersionOperator types.String `tfsdk:"product_version_operator"`
	// File detection-specific fields
	FileSystemDetectionType     types.String `tfsdk:"registry_detection_type"`
	FilePath                    types.String `tfsdk:"file_path"`
	FileFolderName              types.String `tfsdk:"file_or_folder_name"`
	FileSystemDetectionOperator types.String `tfsdk:"filesystem_detection_operator"`
	// PowerShell script detection-specific fields
	ScriptContent         types.String `tfsdk:"script_content"`
	EnforceSignatureCheck types.Bool   `tfsdk:"enforce_signature_check"`
	RunAs32Bit            types.Bool   `tfsdk:"run_as_32_bit"`
}

type Win32LobAppRegistryRequirementResourceModel struct {
	Operator             types.String `tfsdk:"operator"`
	DetectionValue       types.String `tfsdk:"detection_value"`
	Check32BitOn64System types.Bool   `tfsdk:"check_32_bit_on_64_system"`
	KeyPath              types.String `tfsdk:"key_path"`
	ValueName            types.String `tfsdk:"value_name"`
	DetectionType        types.String `tfsdk:"detection_type"`
}

type Win32LobAppRegistryRuleResourceModel struct {
	RuleType             types.String `tfsdk:"rule_type"`
	Check32BitOn64System types.Bool   `tfsdk:"check_32_bit_on_64_system"`
	KeyPath              types.String `tfsdk:"key_path"`
	ValueName            types.String `tfsdk:"value_name"`
	OperationType        types.String `tfsdk:"operation_type"`
	Operator             types.String `tfsdk:"operator"`
	ComparisonValue      types.String `tfsdk:"comparison_value"`
}

type Win32LobAppInstallExperienceResourceModel struct {
	RunAsAccount          types.String `tfsdk:"run_as_account"`
	MaxRunTimeInMinutes   types.Int32  `tfsdk:"max_run_time_in_minutes"`
	DeviceRestartBehavior types.String `tfsdk:"device_restart_behavior"`
}

type Win32LobAppReturnCodeResourceModel struct {
	ReturnCode types.Int32  `tfsdk:"return_code"`
	Type       types.String `tfsdk:"type"`
}

type Win32LobAppMsiInformationResourceModel struct {
	ProductCode    types.String `tfsdk:"product_code"`
	ProductVersion types.String `tfsdk:"product_version"`
	UpgradeCode    types.String `tfsdk:"upgrade_code"`
	RequiresReboot types.Bool   `tfsdk:"requires_reboot"`
	PackageType    types.String `tfsdk:"package_type"`
	ProductName    types.String `tfsdk:"product_name"`
	Publisher      types.String `tfsdk:"publisher"`
}
