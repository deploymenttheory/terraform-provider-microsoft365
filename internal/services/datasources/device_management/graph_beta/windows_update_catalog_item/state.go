package graphBetaWindowsUpdateCatalogItem

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Update Catalog Item to a model
func MapRemoteStateToDataSource(data graphmodels.WindowsUpdateCatalogItemable) WindowsUpdateCatalogItemModel {
	model := WindowsUpdateCatalogItemModel{
		ID:               convert.GraphToFrameworkString(data.GetId()),
		DisplayName:      convert.GraphToFrameworkString(data.GetDisplayName()),
		ReleaseDateTime:  convert.GraphToFrameworkTime(data.GetReleaseDateTime()),
		EndOfSupportDate: convert.GraphToFrameworkTime(data.GetEndOfSupportDate()),
	}

	return model
}
