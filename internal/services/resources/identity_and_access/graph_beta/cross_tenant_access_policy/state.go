package graphBetaCrossTenantAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote crossTenantAccessPolicy API response to Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *CrossTenantAccessPolicyResourceModel, remoteResource graphmodels.CrossTenantAccessPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// This is a singleton — the ID is always the static identifier used throughout the provider.
	data.ID = types.StringValue(singletonID)
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.AllowedCloudEndpoints = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remoteResource.GetAllowedCloudEndpoints())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
