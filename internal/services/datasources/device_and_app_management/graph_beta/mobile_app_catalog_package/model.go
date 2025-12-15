// REF: undocumented api endpoint
//
// IMPORTANT: Microsoft Graph API Limitations for mobileAppCatalogPackages endpoint
// This endpoint has significant limitations with OData query parameters that affect the model design.
//
// WORKING OData Parameters:
//   - odata_filter: Only works with startswith() function (e.g., "startswith(publisherDisplayName, 'Microsoft')")
//   - odata_top: Works for limiting results
//
// NOT WORKING/PROBLEMATIC OData Parameters:
//   - odata_skip: Causes 500 Internal Server errors and timeouts - DO NOT USE
//   - odata_select: Causes 500 Internal Server errors and timeouts - DO NOT USE
//   - odata_orderby: Returns no results when combined with odata_filter - DO NOT COMBINE
//   - odata_count: Returns no results when combined with odata_filter - DO NOT COMBINE
//   - odata_search: Not reliably supported by this endpoint - AVOID
//   - eq operator in filters: Not reliable, use startswith() instead
//
// RECOMMENDED USAGE:
//   - Use simple filter_type values: "all", "id", "product_name", "publisher_name"
//   - If using OData, only use odata_filter (with startswith()) and odata_top
//   - Avoid combining multiple OData parameters as they cause empty results or errors

package graphBetaMobileAppCatalogPackage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppCatalogPackageDataSourceModel defines the data source model
type MobileAppCatalogPackageDataSourceModel struct {
	FilterType   types.String                   `tfsdk:"filter_type"`   // Required field to specify how to filter
	FilterValue  types.String                   `tfsdk:"filter_value"`  // Value to filter by (not used for "all" or "odata")
	ODataFilter  types.String                   `tfsdk:"odata_filter"`  // OData filter parameter
	ODataTop     types.Int32                    `tfsdk:"odata_top"`     // OData top parameter for limiting results
	ODataSkip    types.Int32                    `tfsdk:"odata_skip"`    // OData skip parameter for pagination
	ODataSelect  types.String                   `tfsdk:"odata_select"`  // OData select parameter for field selection
	ODataOrderBy types.String                   `tfsdk:"odata_orderby"` // OData orderby parameter for sorting
	ODataCount   types.Bool                     `tfsdk:"odata_count"`   // OData count parameter
	ODataSearch  types.String                   `tfsdk:"odata_search"`  // OData search parameter
	ODataExpand  types.String                   `tfsdk:"odata_expand"`  // OData expand parameter
	Items        []MobileAppCatalogPackageModel `tfsdk:"items"`         // List of mobile app catalog packages that match the filters
	Timeouts     timeouts.Value                 `tfsdk:"timeouts"`
}

// MobileAppCatalogPackageModel represents a win32CatalogApp with full details
type MobileAppCatalogPackageModel struct {
	// Base mobile app fields
	ID                      types.String   `tfsdk:"id"`
	DisplayName             types.String   `tfsdk:"display_name"`
	Description             types.String   `tfsdk:"description"`
	Publisher               types.String   `tfsdk:"publisher"`
	CreatedDateTime         types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime    types.String   `tfsdk:"last_modified_date_time"`
	IsFeatured              types.Bool     `tfsdk:"is_featured"`
	PrivacyInformationUrl   types.String   `tfsdk:"privacy_information_url"`
	InformationUrl          types.String   `tfsdk:"information_url"`
	Owner                   types.String   `tfsdk:"owner"`
	Developer               types.String   `tfsdk:"developer"`
	Notes                   types.String   `tfsdk:"notes"`
	UploadState             types.Int32    `tfsdk:"upload_state"`
	PublishingState         types.String   `tfsdk:"publishing_state"`
	IsAssigned              types.Bool     `tfsdk:"is_assigned"`
	RoleScopeTagIds         []types.String `tfsdk:"role_scope_tag_ids"`
	DependentAppCount       types.Int32    `tfsdk:"dependent_app_count"`
	SupersedingAppCount     types.Int32    `tfsdk:"superseding_app_count"`
	SupersededAppCount      types.Int32    `tfsdk:"superseded_app_count"`
	CommittedContentVersion types.String   `tfsdk:"committed_content_version"`

	// Win32 specific fields
	FileName                       types.String            `tfsdk:"file_name"`
	Size                           types.Int64             `tfsdk:"size"`
	InstallCommandLine             types.String            `tfsdk:"install_command_line"`
	UninstallCommandLine           types.String            `tfsdk:"uninstall_command_line"`
	ApplicableArchitectures        types.String            `tfsdk:"applicable_architectures"`
	AllowedArchitectures           types.String            `tfsdk:"allowed_architectures"`
	MinimumFreeDiskSpaceInMB       types.Int32             `tfsdk:"minimum_free_disk_space_in_mb"`
	MinimumMemoryInMB              types.Int32             `tfsdk:"minimum_memory_in_mb"`
	MinimumNumberOfProcessors      types.Int32             `tfsdk:"minimum_number_of_processors"`
	MinimumCpuSpeedInMHz           types.Int32             `tfsdk:"minimum_cpu_speed_in_mhz"`
	SetupFilePath                  types.String            `tfsdk:"setup_file_path"`
	MinimumSupportedWindowsRelease types.String            `tfsdk:"minimum_supported_windows_release"`
	DisplayVersion                 types.String            `tfsdk:"display_version"`
	AllowAvailableUninstall        types.Bool              `tfsdk:"allow_available_uninstall"`
	MobileAppCatalogPackageId      types.String            `tfsdk:"mobile_app_catalog_package_id"`
	Rules                          []RuleModel             `tfsdk:"rules"`
	InstallExperience              *InstallExperienceModel `tfsdk:"install_experience"`
	ReturnCodes                    []ReturnCodeModel       `tfsdk:"return_codes"`
	MsiInformation                 *MsiInformationModel    `tfsdk:"msi_information"`
}

// RuleModel represents detection and requirement rules
type RuleModel struct {
	ODataType            types.String `tfsdk:"odata_type"`
	RuleType             types.String `tfsdk:"rule_type"`
	Path                 types.String `tfsdk:"path"`
	FileOrFolderName     types.String `tfsdk:"file_or_folder_name"`
	Check32BitOn64System types.Bool   `tfsdk:"check_32bit_on_64system"`
	OperationType        types.String `tfsdk:"operation_type"`
	Operator             types.String `tfsdk:"operator"`
	ComparisonValue      types.String `tfsdk:"comparison_value"`
	KeyPath              types.String `tfsdk:"key_path"`
	ValueName            types.String `tfsdk:"value_name"`
}

// InstallExperienceModel represents installation experience settings
type InstallExperienceModel struct {
	RunAsAccount          types.String `tfsdk:"run_as_account"`
	MaxRunTimeInMinutes   types.Int32  `tfsdk:"max_run_time_in_minutes"`
	DeviceRestartBehavior types.String `tfsdk:"device_restart_behavior"`
}

// ReturnCodeModel represents application return codes
type ReturnCodeModel struct {
	ReturnCode types.Int32  `tfsdk:"return_code"`
	Type       types.String `tfsdk:"type"`
}

// MsiInformationModel represents MSI package information
type MsiInformationModel struct {
	ProductCode    types.String `tfsdk:"product_code"`
	ProductVersion types.String `tfsdk:"product_version"`
	UpgradeCode    types.String `tfsdk:"upgrade_code"`
	RequiresReboot types.Bool   `tfsdk:"requires_reboot"`
	PackageType    types.String `tfsdk:"package_type"`
	ProductName    types.String `tfsdk:"product_name"`
	Publisher      types.String `tfsdk:"publisher"`
}
