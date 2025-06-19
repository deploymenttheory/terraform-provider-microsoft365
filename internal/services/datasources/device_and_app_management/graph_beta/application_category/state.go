package graphBetaApplicationCategory

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps an Application Category to a model
func MapRemoteStateToDataSource(data graphmodels.MobileAppCategoryable) ApplicationCategoryModel {
	model := ApplicationCategoryModel{
		ID:                   convert.GraphToFrameworkString(data.GetId()),
		DisplayName:          convert.GraphToFrameworkString(data.GetDisplayName()),
		LastModifiedDateTime: convert.GraphToFrameworkTime(data.GetLastModifiedDateTime()),
	}

	return model
}
