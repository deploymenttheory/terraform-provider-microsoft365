package graphBetaReuseablePolicySettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Reusable Policy Setting to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceManagementReusablePolicySettingable) ReuseablePolicySettingModel {
	model := ReuseablePolicySettingModel{
		ID:          types.StringPointerValue(data.GetId()),
		DisplayName: types.StringPointerValue(data.GetDisplayName()),
		Description: types.StringPointerValue(data.GetDescription()),
	}

	return model
}
