// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcontentfile?view=graph-rest-beta
package sharedmodels

import "github.com/hashicorp/terraform-plugin-framework/types"

// MobileAppContentFileResourceModel represents a file within a content version
// Based on the JSON structure from the Microsoft Graph API
type MobileAppContentFileResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Size                      types.Int32  `tfsdk:"size"`
	SizeEncrypted             types.Int32  `tfsdk:"size_encrypted"`
	UploadState               types.String `tfsdk:"upload_state"`
	IsCommitted               types.Bool   `tfsdk:"is_committed"`
	IsDependency              types.Bool   `tfsdk:"is_dependency"`
	IsFrameworkFile           types.Bool   `tfsdk:"is_framework_file"`
	AzureStorageUri           types.String `tfsdk:"azure_storage_uri"`
	AzureStorageUriExpiration types.String `tfsdk:"azure_storage_uri_expiration"`
	CreatedDateTime           types.String `tfsdk:"created_date_time"`
}
