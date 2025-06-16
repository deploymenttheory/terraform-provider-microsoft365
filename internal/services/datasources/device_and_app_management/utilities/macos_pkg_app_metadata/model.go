// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta

package utilityMacOSPKGAppMetadata

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacOSPKGAppMetadataDataSourceModel defines the data source model
type MacOSPKGAppMetadataDataSourceModel struct {
	InstallerFilePathSource types.String             `tfsdk:"installer_file_path_source"` // Path to a local PKG file
	InstallerURLSource      types.String             `tfsdk:"installer_url_source"`       // URL to a PKG file
	Metadata                *MetadataDataSourceModel `tfsdk:"metadata"`                   // Extracted metadata
	Timeouts                timeouts.Value           `tfsdk:"timeouts"`
}

// MetadataDataSourceModel represents detailed metadata for a macOS PKG app
type MetadataDataSourceModel struct {
	// Core bundle metadata (as requested)
	CFBundleIdentifier         types.String `tfsdk:"cf_bundle_identifier"`
	CFBundleShortVersionString types.String `tfsdk:"cf_bundle_short_version_string"`

	// Additional metadata extracted from PKG/XAR
	Name            types.String `tfsdk:"name"`
	PackageIDs      types.List   `tfsdk:"package_ids"`
	InstallLocation types.String `tfsdk:"install_location"`
	AppPaths        types.List   `tfsdk:"app_paths"`
	MinOSVersion    types.String `tfsdk:"min_os_version"`
	SizeMB          types.Int64  `tfsdk:"size_mb"`
	SHA256Checksum  types.String `tfsdk:"sha256_checksum"`
	MD5Checksum     types.String `tfsdk:"md5_checksum"`

	// Included bundles inside the PKG
	IncludedBundles types.List `tfsdk:"included_bundles"`
}

// BundleInfoModel represents information about an included bundle in the PKG
type BundleInfoModel struct {
	BundleID        types.String `tfsdk:"bundle_id"`
	Version         types.String `tfsdk:"version"`
	Path            types.String `tfsdk:"path"`
	CFBundleVersion types.String `tfsdk:"cf_bundle_version"`
}
