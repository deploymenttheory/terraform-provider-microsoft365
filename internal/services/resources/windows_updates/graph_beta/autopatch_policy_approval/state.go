package graphBetaWindowsUpdatesAutopatchPolicyApproval

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchPolicyApprovalResourceModel, remoteResource graphmodelswindowsupdates.PolicyApprovalable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	if statusPtr := remoteResource.GetStatus(); statusPtr != nil {
		data.Status = types.StringValue(statusPtr.String())
	}

	if catalogEntryId := remoteResource.GetCatalogEntryId(); catalogEntryId != nil {
		data.CatalogEntryId = types.StringValue(*catalogEntryId)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
