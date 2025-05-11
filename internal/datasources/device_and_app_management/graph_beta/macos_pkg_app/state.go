package graphBetaMacOSPKGApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a MacOS PKG App to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.MacOSPkgAppable) MacOSPKGAppModel {
	model := MacOSPKGAppModel{
		ID:              state.StringPointerValue(data.GetId()),
		DisplayName:     state.StringPointerValue(data.GetDisplayName()),
		Description:     state.StringPointerValue(data.GetDescription()),
		CreatedDateTime: state.TimeToString(data.GetCreatedDateTime()),
	}

	return model
}
