package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the properties of a RoleScopeTag from Graph API to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleScopeTagResourceModel, remoteResource graphmodels.RoleScopeTagable) {
	if remoteResource == nil {
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsBuiltIn = convert.GraphToFrameworkBool(remoteResource.GetIsBuiltIn())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
