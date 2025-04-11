package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// In state_content.go
// MapContentVersionsStateToTerraform maps API content versions to Terraform state
func MapContentVersionsStateToTerraform(
	ctx context.Context,
	versions []graphmodels.MobileAppContentable,
	versionFiles map[string][]graphmodels.MobileAppContentFileable,
) types.List {
	// Define the object types for our complex structure
	fileObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":                         types.StringType,
			"size":                         types.Int64Type,
			"size_encrypted":               types.Int64Type,
			"upload_state":                 types.StringType,
			"is_committed":                 types.BoolType,
			"is_dependency":                types.BoolType,
			"azure_storage_uri":            types.StringType,
			"azure_storage_uri_expiration": types.StringType,
			"created_date_time":            types.StringType,
		},
	}

	contentVersionObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"files": types.SetType{ElemType: fileObjectType},
		},
	}

	// If there are no versions, return null list
	if len(versions) == 0 {
		return types.ListNull(contentVersionObjectType)
	}

	// Process each version
	var contentVersionElements []attr.Value
	for _, version := range versions {
		if version == nil || version.GetId() == nil {
			continue
		}

		versionId := *version.GetId()

		// Get files for this version from our map
		files, exists := versionFiles[versionId]
		if !exists {
			tflog.Debug(ctx, fmt.Sprintf("No files found for content version %s", versionId))
			files = []graphmodels.MobileAppContentFileable{}
		}

		// Create the files set for this version
		filesSet := mapContentFilesStateToTerraform(ctx, files, fileObjectType)

		// Create version object
		versionValues := map[string]attr.Value{
			"id":    state.StringPointerValue(version.GetId()),
			"files": filesSet,
		}

		versionObj, diags := types.ObjectValue(contentVersionObjectType.AttrTypes, versionValues)
		if diags.HasError() {
			tflog.Warn(ctx, fmt.Sprintf("Error creating version object: %v", diags.Errors()))
			continue
		}

		contentVersionElements = append(contentVersionElements, versionObj)
	}

	// Create the final list value
	contentVersionsList, diags := types.ListValue(contentVersionObjectType, contentVersionElements)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error creating content versions list: %v", diags.Errors()))
		return types.ListNull(contentVersionObjectType)
	}

	return contentVersionsList
}

// Update the existing files mapper to use Terraform types
func mapContentFilesStateToTerraform(
	ctx context.Context,
	files []graphmodels.MobileAppContentFileable,
	fileObjectType types.ObjectType,
) types.Set {
	// If there are no files, return an empty set
	if len(files) == 0 {
		return types.SetValueMust(fileObjectType, []attr.Value{})
	}

	// Process each file
	var fileElements []attr.Value
	for _, file := range files {
		if file == nil {
			continue
		}

		fileValues := map[string]attr.Value{
			"name":                         state.StringPointerValue(file.GetName()),
			"size":                         state.Int64PointerValue(file.GetSize()),
			"size_encrypted":               state.Int64PointerValue(file.GetSizeEncrypted()),
			"upload_state":                 state.EnumPtrToTypeString(file.GetUploadState()),
			"is_committed":                 state.BoolPointerValue(file.GetIsCommitted()),
			"is_dependency":                state.BoolPointerValue(file.GetIsDependency()),
			"azure_storage_uri":            state.StringPointerValue(file.GetAzureStorageUri()),
			"azure_storage_uri_expiration": state.TimeToString(file.GetAzureStorageUriExpirationDateTime()),
			"created_date_time":            state.TimeToString(file.GetCreatedDateTime()),
		}

		fileObj, diags := types.ObjectValue(fileObjectType.AttrTypes, fileValues)
		if diags.HasError() {
			tflog.Warn(ctx, fmt.Sprintf("Error creating file object: %v", diags.Errors()))
			continue
		}

		fileElements = append(fileElements, fileObj)
	}

	// Create the files set
	filesSet, diags := types.SetValue(fileObjectType, fileElements)
	if diags.HasError() {
		tflog.Warn(ctx, fmt.Sprintf("Error creating files set: %v", diags.Errors()))
		return types.SetValueMust(fileObjectType, []attr.Value{})
	}

	return filesSet
}
