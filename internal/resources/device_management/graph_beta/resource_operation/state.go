package graphBetaResourceOperation

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote resource operation to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data ResourceOperationResourceModel, resourceOperation graphmodels.ResourceOperationable) ResourceOperationResourceModel {
	if resourceOperation == nil {
		tflog.Debug(ctx, "Remote resource operation is nil")
		return data
	}

	data.ID = state.StringPointerValue(resourceOperation.GetId())
	data.Resource = state.StringPointerValue(resourceOperation.GetResource())
	data.ResourceName = state.StringPointerValue(resourceOperation.GetResourceName())
	data.ActionName = state.StringPointerValue(resourceOperation.GetActionName())
	data.Description = state.StringPointerValue(resourceOperation.GetDescription())
	data.EnabledForScopeValidation = state.BoolPointerValue(resourceOperation.GetEnabledForScopeValidation())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}
