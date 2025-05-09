package graphBetaBrowserSiteList

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Browser Site List to a model
func MapRemoteStateToDataSource(data graphmodels.BrowserSiteListable) BrowserSiteListResourceModel {
	model := BrowserSiteListResourceModel{
		ID:          state.StringPointerValue(data.GetId()),
		DisplayName: state.StringPointerValue(data.GetDisplayName()),
	}

	return model
}
