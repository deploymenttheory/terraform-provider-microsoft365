package graphBetaBrowserSiteList

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Browser Site List to a model
func MapRemoteStateToDataSource(data graphmodels.BrowserSiteListable) BrowserSiteListResourceModel {
	model := BrowserSiteListResourceModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		DisplayName: convert.GraphToFrameworkString(data.GetDisplayName()),
	}

	return model
}
