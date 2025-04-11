package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// MobileAppMetaDataResourceModel contains common metadata relevant to all mobile app types
// extracted during app upload
type MobileAppMetaDataResourceModel struct {
	InstallerSizeInBytes    types.Int64  `tfsdk:"installer_size_in_bytes"`
	InstallerMD5Checksum    types.String `tfsdk:"installer_md5_checksum"`
	InstallerSHA256Checksum types.String `tfsdk:"installer_sha256_checksum"`
}
