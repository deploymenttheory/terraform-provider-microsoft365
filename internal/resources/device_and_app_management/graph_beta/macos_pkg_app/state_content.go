package graphBetaMacOSPKGApp

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func mapContentFilesStateToTerraform(ctx context.Context, files []graphmodels.MobileAppContentFileable) []sharedmodels.MobileAppContentFileResourceModel {
	var result []sharedmodels.MobileAppContentFileResourceModel

	for _, file := range files {
		if file == nil {
			continue
		}

		result = append(result, sharedmodels.MobileAppContentFileResourceModel{
			Name:                      state.StringPointerValue(file.GetName()),
			Size:                      state.Int64PointerValue(file.GetSize()),
			SizeEncrypted:             state.Int64PointerValue(file.GetSizeEncrypted()),
			UploadState:               state.EnumPtrToTypeString(file.GetUploadState()),
			IsCommitted:               state.BoolPointerValue(file.GetIsCommitted()),
			IsDependency:              state.BoolPointerValue(file.GetIsDependency()),
			AzureStorageUri:           state.StringPointerValue(file.GetAzureStorageUri()),
			AzureStorageUriExpiration: state.TimeToString(file.GetAzureStorageUriExpirationDateTime()),
			CreatedDateTime:           state.TimeToString(file.GetCreatedDateTime()),
		})
	}

	return result
}
