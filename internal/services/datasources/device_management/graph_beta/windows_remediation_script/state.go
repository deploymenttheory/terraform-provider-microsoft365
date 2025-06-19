package graphBetaWindowsRemediationScript

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Remediation Script to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceHealthScriptable) WindowsRemediationScriptModel {
	model := WindowsRemediationScriptModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		DisplayName: convert.GraphToFrameworkString(data.GetDisplayName()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
	}

	return model
}
