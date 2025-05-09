package graphBetaWindowsRemediationScript

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Remediation Script to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceHealthScriptable) WindowsRemediationScriptModel {
	model := WindowsRemediationScriptModel{
		ID:          types.StringPointerValue(data.GetId()),
		DisplayName: types.StringPointerValue(data.GetDisplayName()),
		Description: types.StringPointerValue(data.GetDescription()),
	}

	return model
}
