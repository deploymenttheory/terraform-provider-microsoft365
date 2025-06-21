package graphBetaRBACResourceOperation

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote resource operation to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data RBACResourceOperationResourceModel, resourceOperation graphmodels.ResourceOperationable) RBACResourceOperationResourceModel {
	if resourceOperation == nil {
		tflog.Debug(ctx, "Remote resource operation is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(resourceOperation.GetId())
	data.Resource = convert.GraphToFrameworkString(resourceOperation.GetResource())
	data.ResourceName = convert.GraphToFrameworkString(resourceOperation.GetResourceName())
	data.ActionName = convert.GraphToFrameworkString(resourceOperation.GetActionName())
	data.Description = convert.GraphToFrameworkString(resourceOperation.GetDescription())
	data.EnabledForScopeValidation = convert.GraphToFrameworkBool(resourceOperation.GetEnabledForScopeValidation())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}
