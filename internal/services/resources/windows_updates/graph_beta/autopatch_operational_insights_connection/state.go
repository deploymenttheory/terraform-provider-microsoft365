package graphBetaWindowsUpdatesAutopatchOperationalInsightsConnection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchOperationalInsightsConnectionResourceModel, remoteResource graphmodelswindowsupdates.ResourceConnectionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	if statePtr := remoteResource.GetState(); statePtr != nil {
		data.State = convert.GraphToFrameworkEnum(statePtr)
	}

	if oic, ok := remoteResource.(graphmodelswindowsupdates.OperationalInsightsConnectionable); ok {
		data.AzureResourceGroupName = convert.GraphToFrameworkString(oic.GetAzureResourceGroupName())
		data.AzureSubscriptionId = convert.GraphToFrameworkString(oic.GetAzureSubscriptionId())
		data.WorkspaceName = convert.GraphToFrameworkString(oic.GetWorkspaceName())
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
