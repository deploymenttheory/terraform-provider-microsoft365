package graphBetaWindowsUpdatesAutopatchDeploymentState

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchDeploymentStateResourceModel, remoteResource graphmodelswindowsupdates.Deploymentable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	if deploymentState := remoteResource.GetState(); deploymentState != nil {
		data.RequestedValue = convert.GraphToFrameworkEnum(deploymentState.GetRequestedValue())
		data.EffectiveValue = convert.GraphToFrameworkEnum(deploymentState.GetEffectiveValue())
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
