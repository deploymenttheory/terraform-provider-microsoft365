package sharedStater

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// getFileObjectType returns the type definition for file objects
func getFileObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":                         types.StringType,
			"size":                         types.Int32Type,
			"size_encrypted":               types.Int32Type,
			"upload_state":                 types.StringType,
			"is_committed":                 types.BoolType,
			"is_dependency":                types.BoolType,
			"is_framework_file":            types.BoolType,
			"azure_storage_uri":            types.StringType,
			"azure_storage_uri_expiration": types.StringType,
			"created_date_time":            types.StringType,
		},
	}
}

// getContentVersionObjectType returns the type definition for content version objects
func getContentVersionObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"files": types.SetType{ElemType: getFileObjectType()},
		},
	}
}

// MapCommittedContentVersionStateToTerraform maps the committed content version to Terraform state
// If installerFileName is provided, it will filter to only include files matching that name
func MapCommittedContentVersionStateToTerraform(
	ctx context.Context,
	committedVersionId string,
	respFiles any,
	err error,
	installerFileName string, // Optional - if provided, will filter to this file only
) types.List {
	fileObjectType := getFileObjectType()
	contentVersionObjectType := getContentVersionObjectType()

	var files []graphmodels.MobileAppContentFileable
	if err == nil && respFiles != nil {
		fileCollection, ok := respFiles.(graphmodels.MobileAppContentFileCollectionResponseable)
		if ok {
			files = fileCollection.GetValue()
		} else {
			tflog.Warn(ctx, "Response type assertion failed", map[string]any{
				"actualType": fmt.Sprintf("%T", respFiles),
			})
			files = []graphmodels.MobileAppContentFileable{}
		}
	} else {
		files = []graphmodels.MobileAppContentFileable{}
	}

	var fileElements []attr.Value
	for i, file := range files {
		if file == nil {
			continue
		}

		// Only include the file if it matches the installer file name or if no installer file name was provided
		fileName := convert.GraphToFrameworkString(file.GetName()).ValueString()
		if installerFileName != "" && fileName != installerFileName && filepath.Base(fileName) != installerFileName {
			tflog.Debug(ctx, fmt.Sprintf("Skipping file %s as it doesn't match installer name %s", fileName, installerFileName))
			continue
		}

		fileValues := map[string]attr.Value{
			"name":                         convert.GraphToFrameworkString(file.GetName()),
			"size":                         convert.GraphToFrameworkInt64(file.GetSize()),
			"size_encrypted":               convert.GraphToFrameworkInt64(file.GetSizeEncrypted()),
			"upload_state":                 convert.GraphToFrameworkEnum(file.GetUploadState()),
			"is_committed":                 convert.GraphToFrameworkBool(file.GetIsCommitted()),
			"is_dependency":                convert.GraphToFrameworkBool(file.GetIsDependency()),
			"is_framework_file":            convert.GraphToFrameworkBool(file.GetIsFrameworkFile()),
			"azure_storage_uri":            convert.GraphToFrameworkString(file.GetAzureStorageUri()),
			"azure_storage_uri_expiration": convert.GraphToFrameworkTime(file.GetAzureStorageUriExpirationDateTime()),
			"created_date_time":            convert.GraphToFrameworkTime(file.GetCreatedDateTime()),
		}

		fileObj, diags := types.ObjectValue(fileObjectType.AttrTypes, fileValues)
		if diags.HasError() {
			tflog.Warn(ctx, "Error creating file object", map[string]any{
				"index": i,
				"error": diags.Errors(),
			})
			continue
		}

		fileElements = append(fileElements, fileObj)
	}

	// Create files set
	filesSet, diags := types.SetValue(fileObjectType, fileElements)
	if diags.HasError() {
		filesSet = types.SetValueMust(fileObjectType, []attr.Value{})
	}

	// Create version object
	versionValues := map[string]attr.Value{
		"id":    types.StringValue(committedVersionId),
		"files": filesSet,
	}

	versionObj, diags := types.ObjectValue(contentVersionObjectType.AttrTypes, versionValues)
	if diags.HasError() {
		return types.ListNull(contentVersionObjectType)
	}

	// Create list with single element
	contentVersionsList, diags := types.ListValue(contentVersionObjectType, []attr.Value{versionObj})
	if diags.HasError() {
		return types.ListNull(contentVersionObjectType)
	}

	return contentVersionsList
}
