package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// MapRemoteStateToTerraform iterates the members page iterator returned from
// GET .../updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
// and maps member IDs back into entra_device_ids.
//
// When the plan already has a known set (i.e. not null/unknown), only members
// that appear in that planned set are retained. This prevents drift caused by
// devices that were added to the group outside of Terraform.
func MapRemoteStateToTerraform(
	ctx context.Context,
	data *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel,
	pageIterator *graphcore.PageIterator[graphmodelswindowsupdates.UpdatableAssetable],
) {
	if pageIterator == nil {
		tflog.Debug(ctx, "Page iterator is nil, skipping member mapping")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping remote member state to Terraform state for %s", ResourceName))

	filterByPlanned := !data.EntraDeviceIds.IsNull() && !data.EntraDeviceIds.IsUnknown()
	plannedDeviceIDs := make(map[string]bool)
	if filterByPlanned {
		for _, elem := range data.EntraDeviceIds.Elements() {
			if strVal, ok := elem.(types.String); ok {
				plannedDeviceIDs[strVal.ValueString()] = true
			}
		}
	}

	memberIDs := make([]attr.Value, 0)
	_ = pageIterator.Iterate(ctx, func(item graphmodelswindowsupdates.UpdatableAssetable) bool {
		if item == nil {
			return true
		}
		id := item.GetId()
		if id == nil {
			return true
		}
		if !filterByPlanned || plannedDeviceIDs[*id] {
			memberIDs = append(memberIDs, types.StringValue(*id))
		}
		return true
	})

	data.EntraDeviceIds = types.SetValueMust(types.StringType, memberIDs)
	data.ID = data.UpdatableAssetGroupId

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
