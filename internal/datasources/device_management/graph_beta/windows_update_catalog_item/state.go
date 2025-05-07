package graphBetaWindowsUpdateCatalogItem

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Update Catalog Item to a model
func MapRemoteStateToDataSource(item graphmodels.WindowsUpdateCatalogItemable) WindowsUpdateCatalogItemModel {
	model := WindowsUpdateCatalogItemModel{
		ID:          types.StringPointerValue(item.GetId()),
		DisplayName: types.StringPointerValue(item.GetDisplayName()),
	}

	if releaseDateTime := item.GetReleaseDateTime(); releaseDateTime != nil {
		model.ReleaseDateTime = types.StringValue(releaseDateTime.Format(time.RFC3339))
	}

	if endOfSupportDate := item.GetEndOfSupportDate(); endOfSupportDate != nil {
		model.EndOfSupportDate = types.StringValue(endOfSupportDate.Format(time.RFC3339))
	}

	return model
}
