package graphBetaWindowsUpdatesAutopatchDeploymentAudience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchDeploymentAudienceResourceModel, remoteResource graphmodelswindowsupdates.DeploymentAudienceable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = types.StringValue(*remoteResource.GetId())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state")
}
