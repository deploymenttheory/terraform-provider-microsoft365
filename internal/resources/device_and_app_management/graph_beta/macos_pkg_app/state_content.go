package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state" // Assuming this contains helper funcs like StringPointerValue etc.
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapCommittedContentVersionStateToTerraform maps the committed content version to Terraform state
// This consolidated function directly processes a single version and its files
func MapCommittedContentVersionStateToTerraform(
	ctx context.Context,
	committedVersionId string,
	respFiles interface{}, // This is typically MobileAppContentFileCollectionResponseable
	err error,
) types.List {
	// --- Added Logging ---
	tflog.Debug(ctx, "Entering MapCommittedContentVersionStateToTerraform", map[string]interface{}{
		"committedVersionId": committedVersionId,
		"respFilesType":      fmt.Sprintf("%T", respFiles),
		"errorReceived":      err,
	})
	// --- End Logging ---

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

	// Handle files based on response
	var files []graphmodels.MobileAppContentFileable
	if err == nil && respFiles != nil {
		// Assume respFiles is a MobileAppContentFileCollectionResponseable
		// Need to use type assertion to get the actual value
		fileCollection, ok := respFiles.(graphmodels.MobileAppContentFileCollectionResponseable)
		if ok {
			files = fileCollection.GetValue()
			// --- Added Logging ---
			tflog.Debug(ctx, "Successfully extracted files from response", map[string]interface{}{
				"fileCount": len(files),
			})
			// --- End Logging ---
		} else {
			// Log if the type assertion failed, indicating an unexpected response type
			tflog.Warn(ctx, "Response type assertion to MobileAppContentFileCollectionResponseable failed", map[string]interface{}{
				"actualType": fmt.Sprintf("%T", respFiles),
			})
			files = []graphmodels.MobileAppContentFileable{}
		}
	} else {
		// Log details if an error was passed in or respFiles was nil
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Mapping received an error for content version %s, resulting in empty file list.", committedVersionId), map[string]interface{}{
				"error": err.Error(),
			})
		} else if respFiles == nil {
			tflog.Warn(ctx, fmt.Sprintf("Mapping received nil respFiles for content version %s, resulting in empty file list.", committedVersionId))
		}
		files = []graphmodels.MobileAppContentFileable{} // Ensure files is empty slice on error or nil response
	}

	// Process files for this version
	var fileElements []attr.Value
	// --- Added Logging ---
	tflog.Debug(ctx, "Processing files for content version", map[string]interface{}{
		"committedVersionId": committedVersionId,
		"numberOfFiles":      len(files),
	})
	// --- End Logging ---
	for i, file := range files { // Added index for logging
		if file == nil {
			tflog.Warn(ctx, "Encountered nil file in files slice", map[string]interface{}{"index": i})
			continue
		}

		// --- Added Logging ---
		tflog.Debug(ctx, "Processing file", map[string]interface{}{
			"index":       i,
			"fileName":    state.StringPointerValue(file.GetName()).ValueString(), // Log actual values where safe
			"fileSize":    state.Int64PointerValue(file.GetSize()).ValueInt64(),
			"uploadState": state.EnumPtrToTypeString(file.GetUploadState()).ValueString(),
			"isCommitted": state.BoolPointerValue(file.GetIsCommitted()).ValueBool(),
		})
		// --- End Logging ---

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
			tflog.Warn(ctx, "Error creating file object value", map[string]interface{}{
				"index":    i,
				"fileName": state.StringPointerValue(file.GetName()).ValueString(),
				"errors":   fmt.Sprintf("%v", diags.Errors()), // Convert errors for logging
			})
			continue // Skip this file if object creation fails
		}

		fileElements = append(fileElements, fileObj)
	}

	// Create the files set
	filesSet, diags := types.SetValue(fileObjectType, fileElements)
	if diags.HasError() {
		// Log error and default to an empty set
		tflog.Warn(ctx, "Error creating files set value", map[string]interface{}{
			"committedVersionId": committedVersionId,
			"errors":             fmt.Sprintf("%v", diags.Errors()),
		})
		// Use SetValueMust for guaranteed empty set on error
		filesSet = types.SetValueMust(fileObjectType, []attr.Value{})
	} else {
		// --- Added Logging ---
		tflog.Debug(ctx, "Successfully created files set", map[string]interface{}{
			"committedVersionId": committedVersionId,
			"numberOfElements":   len(fileElements),
		})
		// --- End Logging ---
	}

	// Create the version object
	versionValues := map[string]attr.Value{
		"id":    types.StringValue(committedVersionId),
		"files": filesSet,
	}

	versionObj, diags := types.ObjectValue(contentVersionObjectType.AttrTypes, versionValues)
	if diags.HasError() {
		// Log error and return null list
		tflog.Warn(ctx, "Error creating version object value", map[string]interface{}{
			"committedVersionId": committedVersionId,
			"errors":             fmt.Sprintf("%v", diags.Errors()),
		})
		return types.ListNull(contentVersionObjectType)
	}

	// Create the list with a single element
	contentVersionsList, diags := types.ListValue(contentVersionObjectType, []attr.Value{versionObj})
	if diags.HasError() {
		// Log error and return null list
		tflog.Error(ctx, "Error creating final content versions list value", map[string]interface{}{
			"committedVersionId": committedVersionId,
			"errors":             fmt.Sprintf("%v", diags.Errors()),
		})
		return types.ListNull(contentVersionObjectType)
	}

	// --- Added Logging ---
	tflog.Debug(ctx, "Exiting MapCommittedContentVersionStateToTerraform successfully", map[string]interface{}{
		"committedVersionId": committedVersionId,
		"listIsNull":         contentVersionsList.IsNull(),
		"numberOfVersions":   len(contentVersionsList.Elements()),
	})
	// --- End Logging ---

	return contentVersionsList
}
