package graphBetaUsersUserManager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote manager state to the Terraform resource model.
func MapRemoteStateToTerraform(ctx context.Context, state *UserManagerResourceModel, manager graphmodels.DirectoryObjectable) {
	if manager == nil {
		tflog.Debug(ctx, "Manager is nil, cannot map state")
		return
	}

	// The manager ID comes from the returned directoryObject
	if manager.GetId() != nil {
		state.ManagerID = types.StringValue(*manager.GetId())
		tflog.Debug(ctx, "Mapped manager ID", map[string]any{
			"manager_id": *manager.GetId(),
		})
	}
}
