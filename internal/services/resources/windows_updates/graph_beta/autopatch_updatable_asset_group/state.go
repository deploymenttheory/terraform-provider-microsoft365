package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// MapRemoteStateToTerraform maps the remote Graph API state for an updatable asset group
// and its members back into the Terraform model.
//
// remoteResource is the result of GET /admin/windows/updates/updatableAssets/{id}.
// pageIterator is the result of GET .../microsoft.graph.windowsUpdates.updatableAssetGroup/members.
//
// When the model already has a known, non-null set for entra_device_object_ids (i.e. state
// set from plan before a ReadWithRetry call), only the subset of members returned by the API
// that appear in that planned set are retained. This prevents out-of-band additions from
// causing drift.
func MapRemoteStateToTerraform(
	ctx context.Context,
	data *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel,
	remoteResource graphmodelswindowsupdates.UpdatableAssetable,
	pageIterator *graphcore.PageIterator[graphmodelswindowsupdates.UpdatableAssetable],
) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil, skipping state mapping")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	if pageIterator == nil {
		tflog.Debug(ctx, "Page iterator is nil, skipping member mapping")
		return
	}

	filterByPlanned := !data.EntraDeviceObjectIds.IsNull() && !data.EntraDeviceObjectIds.IsUnknown()
	plannedIDs := make(map[string]bool)
	if filterByPlanned {
		for _, elem := range data.EntraDeviceObjectIds.Elements() {
			if strVal, ok := elem.(types.String); ok {
				plannedIDs[strVal.ValueString()] = true
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
		if !filterByPlanned || plannedIDs[*id] {
			memberIDs = append(memberIDs, types.StringValue(*id))
		}
		return true
	})

	// Only overwrite if there are members to record, or the config explicitly set an
	// empty set (i.e. original was non-null). When the attribute was never set (null)
	// and the API returns no members, preserve null so the plan/state comparison
	// sees null==null rather than null!=[] and avoids an inconsistent-result error.
	if len(memberIDs) > 0 || !data.EntraDeviceObjectIds.IsNull() {
		data.EntraDeviceObjectIds = types.SetValueMust(types.StringType, memberIDs)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
