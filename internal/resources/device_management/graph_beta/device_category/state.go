package graphBetaDeviceCategory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of an DeviceCategoryResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *DeviceCategoryResourceModel, remoteResource graphmodels.DeviceCategoryable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": types.StringPointerValue(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
