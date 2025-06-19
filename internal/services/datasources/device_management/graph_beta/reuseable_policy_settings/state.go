package graphBetaReuseablePolicySettings

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Reusable Policy Setting to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceManagementReusablePolicySettingable) ReuseablePolicySettingModel {
	model := ReuseablePolicySettingModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		DisplayName: convert.GraphToFrameworkString(data.GetDisplayName()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
	}

	return model
}
