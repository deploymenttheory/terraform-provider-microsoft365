package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Windows Quality Update Expedite Policy to a model
func MapRemoteStateToDataSource(data graphmodels.WindowsQualityUpdateProfileable) WindowsQualityUpdateExpeditePolicyModel {
	model := WindowsQualityUpdateExpeditePolicyModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		DisplayName: convert.GraphToFrameworkString(data.GetDisplayName()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
	}

	return model
}
