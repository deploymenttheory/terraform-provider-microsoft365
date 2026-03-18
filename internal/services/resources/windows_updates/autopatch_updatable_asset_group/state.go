package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel, remoteResource graphmodelswindowsupdates.UpdatableAssetable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
