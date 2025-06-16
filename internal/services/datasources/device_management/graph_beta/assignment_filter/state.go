package graphBetaAssignmentFilter

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps an Assignment Filter to a model
func MapRemoteStateToDataSource(data graphmodels.DeviceAndAppManagementAssignmentFilterable) AssignmentFilterModel {
	model := AssignmentFilterModel{
		ID:          types.StringPointerValue(data.GetId()),
		DisplayName: types.StringPointerValue(data.GetDisplayName()),
		Description: types.StringPointerValue(data.GetDescription()),
	}

	return model
}
