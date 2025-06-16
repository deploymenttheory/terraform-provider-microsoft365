package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Quality Update Expedite Policy to a model
func MapRemoteStateToDataSource(data graphmodels.WindowsQualityUpdateProfileable) WindowsQualityUpdateExpeditePolicyModel {
	model := WindowsQualityUpdateExpeditePolicyModel{
		ID:          types.StringPointerValue(data.GetId()),
		DisplayName: types.StringPointerValue(data.GetDisplayName()),
		Description: types.StringPointerValue(data.GetDescription()),
	}

	return model
}
