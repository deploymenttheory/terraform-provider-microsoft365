package graphBetaAssignmentFilter

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps an Assignment Filter to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceAndAppManagementAssignmentFilterable) AssignmentFilterModel {
	model := AssignmentFilterModel{
		ID:          convert.GraphToFrameworkString(data.GetId()),
		DisplayName: convert.GraphToFrameworkString(data.GetDisplayName()),
		Description: convert.GraphToFrameworkString(data.GetDescription()),
	}

	return model
}
