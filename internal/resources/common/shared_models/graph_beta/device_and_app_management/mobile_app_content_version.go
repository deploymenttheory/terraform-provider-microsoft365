// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcontentfile?view=graph-rest-beta
package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

type MobileAppContentVersionResourceModel struct {
	ID    types.String                        `tfsdk:"id"`
	Files []MobileAppContentFileResourceModel `tfsdk:"files"`
}

type MobileAppContentFileResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Size                      types.Int64  `tfsdk:"size"`
	SizeEncrypted             types.Int64  `tfsdk:"size_encrypted"`
	UploadState               types.String `tfsdk:"upload_state"`
	IsCommitted               types.Bool   `tfsdk:"is_committed"`
	IsDependency              types.Bool   `tfsdk:"is_dependency"`
	AzureStorageUri           types.String `tfsdk:"azure_storage_uri"`
	AzureStorageUriExpiration types.String `tfsdk:"azure_storage_uri_expiration"`
	CreatedDateTime           types.String `tfsdk:"created_date_time"`
}
