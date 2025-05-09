package graphBetaApplicationCategory

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps an Application Category to a model
func MapRemoteStateToDataSource(data graphmodels.MobileAppCategoryable) ApplicationCategoryModel {
	model := ApplicationCategoryModel{
		ID:                   types.StringPointerValue(data.GetId()),
		DisplayName:          types.StringPointerValue(data.GetDisplayName()),
		LastModifiedDateTime: state.TimeToString(data.GetLastModifiedDateTime()),
	}

	return model
}
