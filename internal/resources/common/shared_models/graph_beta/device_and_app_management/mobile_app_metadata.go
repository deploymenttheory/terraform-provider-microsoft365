package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// MobileAppMetaDataResourceModel contains common metadata relevant to all mobile app types
// extracted during app upload
type MobileAppMetaDataResourceModel struct {
	InstallerFilePathSource types.String `tfsdk:"installer_file_path_source"`
	InstallerURLSource      types.String `tfsdk:"installer_url_source"`
	InstallerSizeInBytes    types.Int64  `tfsdk:"installer_size_in_bytes"`
	InstallerMD5Checksum    types.String `tfsdk:"installer_md5_checksum"`
	InstallerSHA256Checksum types.String `tfsdk:"installer_sha256_checksum"`
}

// MobileAppIconResourceModel contains the icon url and file path sources
type MobileAppIconResourceModel struct {
	IconFilePathSource types.String `tfsdk:"icon_file_path_source"`
	IconURLSource      types.String `tfsdk:"icon_url_source"`
}
