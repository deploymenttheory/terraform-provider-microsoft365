package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote security defaults enforcement policy from Kiota SDK to Terraform state.
// The fixed singleton ID is always used instead of the object id Graph returns.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data SecurityDefaultsEnforcementPolicyResourceModel,
	remoteResource graphmodels.IdentitySecurityDefaultsEnforcementPolicyable,
) SecurityDefaultsEnforcementPolicyResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": singletonID,
	})

	data.ID = types.StringValue(singletonID)
	data.IsEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsEnabled())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}
