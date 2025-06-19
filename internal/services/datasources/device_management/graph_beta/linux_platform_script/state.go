package graphBetaLinuxPlatformScript

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Linux Platform Script to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceManagementConfigurationPolicyable) LinuxPlatformScriptModel {
	model := LinuxPlatformScriptModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		DisplayName: convert.GraphToFrameworkString(data.GetName()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
	}

	return model
}
