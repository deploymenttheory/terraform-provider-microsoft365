package graphBetaRoleScopeTag

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Role Scope Tag to a model
func MapRemoteStateToDataSource(data graphmodels.RoleScopeTagable) RoleScopeTagModel {
	model := RoleScopeTagModel{
		ID:          types.StringPointerValue(data.GetId()),
		DisplayName: types.StringPointerValue(data.GetDisplayName()),
		Description: types.StringPointerValue(data.GetDescription()),
	}

	return model
}
