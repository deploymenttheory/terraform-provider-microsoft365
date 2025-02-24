package graphBetaMacOSPKGApp

// import (
// 	"context"
// 	"encoding/base64"

// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
// )

// // MapContentVersionsToState maps the Graph API MobileAppContent response to Terraform state.
// //
// // This function converts the properties of the latest MobileAppContent version into a
// // ContentVersionResourceModel and assigns it directly to the Terraform state without using a pointer.
// //
// // Parameters:
// //   - ctx: The context for logging or cancellation.
// //   - data: The Terraform resource model to populate.
// //   - contentVersions: The Graph API collection response for MobileAppContent.
// func MapContentVersionsToState(ctx context.Context, data *MacOSPKGAppResourceModel, contentVersions graphmodels.MobileAppContentCollectionResponseable) {
// 	if contentVersions == nil || len(contentVersions.GetValue()) == 0 {
// 		tflog.Debug(ctx, "No content versions found")
// 		return
// 	}

// 	versions := contentVersions.GetValue()
// 	tflog.Debug(ctx, "Mapping content versions", map[string]interface{}{
// 		"resourceId":   data.ID.ValueString(),
// 		"versionCount": len(versions),
// 	})

// 	latest := versions[len(versions)-1]

// 	var files []ContentFileResourceModel
// 	for _, file := range latest.GetFiles() {
// 		uploadState := types.StringNull()
// 		if file.GetUploadState() != nil {
// 			uploadState = types.StringValue((*file.GetUploadState()).String())
// 		}

// 		manifest := types.StringNull()
// 		if mb := file.GetManifest(); mb != nil && len(mb) > 0 {
// 			manifest = types.StringValue(base64.StdEncoding.EncodeToString(mb))
// 		}

// 		files = append(files, ContentFileResourceModel{
// 			Id:                        types.StringValue(state.StringPtrToString(file.GetId())),
// 			Name:                      types.StringValue(state.StringPtrToString(file.GetName())),
// 			IsDependency:              state.BoolPtrToTypeBool(file.GetIsDependency()),
// 			IsCommitted:               state.BoolPtrToTypeBool(file.GetIsCommitted()),
// 			Size:                      state.Int64PtrToTypeInt64(file.GetSize()),
// 			SizeEncrypted:             state.Int64PtrToTypeInt64(file.GetSizeEncrypted()),
// 			UploadState:               uploadState,
// 			CreatedDateTime:           state.TimeToString(file.GetCreatedDateTime()),
// 			AzureStorageUri:           types.StringValue(state.StringPtrToString(file.GetAzureStorageUri())),
// 			AzureStorageUriExpiration: state.TimeToString(file.GetAzureStorageUriExpirationDateTime()),
// 			IsFrameworkFile:           state.BoolPtrToTypeBool(file.GetIsFrameworkFile()),
// 			Manifest:                  manifest,
// 			SizeEncryptedInBytes:      state.Int64PtrToTypeInt64(file.GetSizeEncryptedInBytes()),
// 			SizeInBytes:               state.Int64PtrToTypeInt64(file.GetSizeInBytes()),
// 		})
// 	}

// 	// Assign the constructed ContentVersionResourceModel directly without using a pointer.
// 	data.ContentVersion = ContentVersionResourceModel{
// 		ContentVersionId: types.StringValue(state.StringPtrToString(latest.GetId())),
// 		FileCount:        types.Int64Value(int64(len(latest.GetFiles()))),
// 		Files:            files,
// 	}

// 	tflog.Debug(ctx, "Finished mapping content versions", map[string]interface{}{
// 		"contentVersionId": data.ContentVersion.ContentVersionId.ValueString(),
// 		"fileCount":        data.ContentVersion.FileCount.ValueInt64(),
// 	})
// }
