package enterpriseappcatalog

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MobileAppDataSourceModel defines the data source model
type MobileAppDataSourceModel struct {
	FilterType       types.String                   `tfsdk:"filter_type"`        // Required field to specify how to filter
	FilterValue      types.String                   `tfsdk:"filter_value"`       // Value to filter by (not used for "all" or "odata")
	ODataFilter      types.String                   `tfsdk:"odata_filter"`       // OData filter parameter
	ODataTop         types.Int32                    `tfsdk:"odata_top"`          // OData top parameter for limiting results
	ODataSkip        types.Int32                    `tfsdk:"odata_skip"`         // OData skip parameter for pagination
	ODataSelect      types.String                   `tfsdk:"odata_select"`       // OData select parameter for field selection
	ODataOrderBy     types.String                   `tfsdk:"odata_orderby"`      // OData orderby parameter for sorting
	AppTypeFilter    types.String                   `tfsdk:"app_type_filter"`    // Optional filter to filter by odata mobile app type
	IncludeAppConfig types.Bool                     `tfsdk:"include_app_config"` // Whether to include detailed app configuration
	Items            []MobileAppCatalogPackageModel `tfsdk:"items"`              // List of mobile apps that match the filters
	Timeouts         timeouts.Value                 `tfsdk:"timeouts"`
}

// MobileAppCatalogPackageModel represents a single mobile app catalog package
type MobileAppCatalogPackageModel struct {
	// Package reference fields
	Id                       types.String   `tfsdk:"id"`
	ProductId                types.String   `tfsdk:"product_id"`
	ProductDisplayName       types.String   `tfsdk:"product_display_name"`
	PublisherDisplayName     types.String   `tfsdk:"publisher_display_name"`
	VersionDisplayName       types.String   `tfsdk:"version_display_name"`
	BranchDisplayName        types.String   `tfsdk:"branch_display_name"`
	ApplicableArchitectures  types.String   `tfsdk:"applicable_architectures"`
	Locales                  []types.String `tfsdk:"locales"`
	PackageAutoUpdateCapable types.Bool     `tfsdk:"package_auto_update_capable"`

	// Detailed app configuration
	AppConfig *MobileAppCatalogPackageConfigurationModel `tfsdk:"app_config"`
}

// MobileAppCatalogPackageConfigurationModel represents the detailed app configuration from convertFromMobileAppCatalogPackage
type MobileAppCatalogPackageConfigurationModel struct {
	ODataType               types.String            `tfsdk:"odata_type"`
	DisplayName             types.String            `tfsdk:"display_name"`
	Description             types.String            `tfsdk:"description"`
	Publisher               types.String            `tfsdk:"publisher"`
	Developer               types.String            `tfsdk:"developer"`
	PrivacyInformationUrl   types.String            `tfsdk:"privacy_information_url"`
	InformationUrl          types.String            `tfsdk:"information_url"`
	FileName                types.String            `tfsdk:"file_name"`
	Size                    types.Int64             `tfsdk:"size"`
	InstallCommandLine      types.String            `tfsdk:"install_command_line"`
	UninstallCommandLine    types.String            `tfsdk:"uninstall_command_line"`
	ApplicableArchitectures types.String            `tfsdk:"applicable_architectures"`
	AllowedArchitectures    types.String            `tfsdk:"allowed_architectures"`
	SetupFilePath           types.String            `tfsdk:"setup_file_path"`
	MinSupportedWinRelease  types.String            `tfsdk:"min_supported_windows_release"`
	DisplayVersion          types.String            `tfsdk:"display_version"`
	AllowAvailableUninstall types.Bool              `tfsdk:"allow_available_uninstall"`
	Rules                   []RuleModel             `tfsdk:"rules"`
	InstallExperience       *InstallExperienceModel `tfsdk:"install_experience"`
	ReturnCodes             []ReturnCodeModel       `tfsdk:"return_codes"`
}

// RuleModel represents a detection or requirement rule
type RuleModel struct {
	ODataType            types.String `tfsdk:"odata_type"`
	RuleType             types.String `tfsdk:"rule_type"`
	Path                 types.String `tfsdk:"path"`
	FileOrFolderName     types.String `tfsdk:"file_or_folder_name"`
	Check32BitOn64System types.Bool   `tfsdk:"check_32bit_on_64_system"`
	OperationType        types.String `tfsdk:"operation_type"`
	Operator             types.String `tfsdk:"operator"`
	ComparisonValue      types.String `tfsdk:"comparison_value"`
	KeyPath              types.String `tfsdk:"key_path"`
	ValueName            types.String `tfsdk:"value_name"`
}

// InstallExperienceModel represents the install experience configuration
type InstallExperienceModel struct {
	RunAsAccount          types.String `tfsdk:"run_as_account"`
	MaxRunTimeInMinutes   types.Int64  `tfsdk:"max_run_time_in_minutes"`
	DeviceRestartBehavior types.String `tfsdk:"device_restart_behavior"`
}

// ReturnCodeModel represents a return code configuration
type ReturnCodeModel struct {
	ReturnCode types.Int64  `tfsdk:"return_code"`
	Type       types.String `tfsdk:"type"`
}
