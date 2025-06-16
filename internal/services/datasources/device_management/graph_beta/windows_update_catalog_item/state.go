package graphBetaWindowsUpdateCatalogItem

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Update Catalog Item to a model
func MapRemoteStateToDataSource(data graphmodels.WindowsUpdateCatalogItemable) WindowsUpdateCatalogItemModel {
	model := WindowsUpdateCatalogItemModel{
		ID:               types.StringPointerValue(data.GetId()),
		DisplayName:      types.StringPointerValue(data.GetDisplayName()),
		ReleaseDateTime:  state.TimeToString(data.GetReleaseDateTime()),
		EndOfSupportDate: state.TimeToString(data.GetEndOfSupportDate()),
	}

	return model
}
