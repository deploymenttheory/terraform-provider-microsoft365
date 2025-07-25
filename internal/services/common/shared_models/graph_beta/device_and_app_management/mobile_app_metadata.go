package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// MobileAppMetaDataResourceModel contains common metadata relevant to all mobile app types
// extracted during app upload
type MobileAppMetaDataResourceModel struct {
	InstallerFilePathSource types.String `tfsdk:"installer_file_path_source"`
	InstallerURLSource      types.String `tfsdk:"installer_url_source"`
}

// MobileAppIconResourceModel contains the icon url and file path sources
type MobileAppIconResourceModel struct {
	IconFilePathSource types.String `tfsdk:"icon_file_path_source"`
	IconURLSource      types.String `tfsdk:"icon_url_source"`
}

// ImageResourceModel contains the image url and file path sources
type ImageResourceModel struct {
	ImageFilePathSource types.String `tfsdk:"image_file_path_source"`
	ImageURLSource      types.String `tfsdk:"image_url_source"`
}
