package graphBetaLinuxPlatformScript

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Linux Platform Script to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceManagementConfigurationPolicyable) LinuxPlatformScriptModel {
	model := LinuxPlatformScriptModel{
		ID:          types.StringPointerValue(data.GetId()),
		DisplayName: types.StringPointerValue(data.GetName()),
		Description: types.StringPointerValue(data.GetDescription()),
	}

	return model
}
