// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsmsiapp?view=graph-rest-beta

package utilityWindowsMSIAppMetadata

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsMSIAppMetadataDataSourceModel defines the data source model
type WindowsMSIAppMetadataDataSourceModel struct {
	ID                      types.String             `tfsdk:"id"`                         // Unique identifier for the data source
	InstallerFilePathSource types.String             `tfsdk:"installer_file_path_source"` // Path to a local MSI file
	InstallerURLSource      types.String             `tfsdk:"installer_url_source"`       // URL to an MSI file
	Metadata                *MetadataDataSourceModel `tfsdk:"metadata"`                   // Extracted metadata
	Timeouts                timeouts.Value           `tfsdk:"timeouts"`
}

// MetadataDataSourceModel represents detailed metadata for a Windows MSI app
type MetadataDataSourceModel struct {
	// Core MSI properties
	ProductCode    types.String `tfsdk:"product_code"`    // MSI Product Code (GUID)
	ProductVersion types.String `tfsdk:"product_version"` // MSI Product Version
	ProductName    types.String `tfsdk:"product_name"`    // MSI Product Name
	Publisher      types.String `tfsdk:"publisher"`       // MSI Manufacturer

	// Additional metadata
	UpgradeCode      types.String  `tfsdk:"upgrade_code"`      // MSI Upgrade Code (GUID)
	Language         types.String  `tfsdk:"language"`          // MSI Language
	PackageType      types.String  `tfsdk:"package_type"`      // MSI Package Type (e.g., Application, Patch)
	InstallLocation  types.String  `tfsdk:"install_location"`  // Default installation location
	InstallCommand   types.String  `tfsdk:"install_command"`   // Default install command
	UninstallCommand types.String  `tfsdk:"uninstall_command"` // Default uninstall command
	TransformPaths   types.List    `tfsdk:"transform_paths"`   // MST transform paths
	SizeMB           types.Float64 `tfsdk:"size_mb"`           // Size in MB
	SHA256Checksum   types.String  `tfsdk:"sha256_checksum"`   // SHA256 hash
	MD5Checksum      types.String  `tfsdk:"md5_checksum"`      // MD5 hash

	// MSI tables and properties
	Properties       types.Map  `tfsdk:"properties"`        // All MSI properties
	RequiredFeatures types.List `tfsdk:"required_features"` // Required features
	Files            types.List `tfsdk:"files"`             // Files included in the MSI

	// System requirements
	MinOSVersion types.String `tfsdk:"min_os_version"` // Minimum OS version
	Architecture types.String `tfsdk:"architecture"`   // Target architecture (x86, x64, etc.)
}

// MSIFileModel represents information about a file in the MSI
type MSIFileModel struct {
	FileName    types.String `tfsdk:"file_name"`
	FilePath    types.String `tfsdk:"file_path"`
	FileSize    types.Int64  `tfsdk:"file_size"`
	Version     types.String `tfsdk:"version"`
	Description types.String `tfsdk:"description"`
}

// MSIFeatureModel represents information about a feature in the MSI
type MSIFeatureModel struct {
	FeatureName types.String `tfsdk:"feature_name"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Required    types.Bool   `tfsdk:"required"`
}
