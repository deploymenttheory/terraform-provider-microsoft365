package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the properties of a RoleScopeTag from Graph API to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *RoleScopeTagResourceModel, remoteResource graphmodels.RoleScopeTagable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		// Initialize with empty/default values
		data.ID = types.StringNull()
		data.DisplayName = types.StringNull()
		data.Description = types.StringNull()
		data.IsBuiltIn = types.BoolNull()
		data.Assignments = make([]types.String, 0)
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.IsBuiltIn = types.BoolPointerValue(remoteResource.GetIsBuiltIn())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
