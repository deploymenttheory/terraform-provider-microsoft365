package graphBetaMacOSPKGApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the properties of a MacOSPkgApp to Terraform state.
//
// This function handles the conversion of Graph API model properties to Terraform state.
// It follows a direct mapping approach with proper logging.
//
// Parameters:
//   - ctx: Context for logging or cancellation
//   - data: Terraform model to populate
//   - remoteResource: Graph API model containing source data
func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSPKGAppResourceModel, remoteResource graphmodels.MacOSPkgAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = state.StringPointerValue(remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = state.StringPointerValue(remoteResource.GetDescription())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())

	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId":  data.ID.ValueString(),
		"displayName": data.DisplayName.ValueString(),
	})
}
