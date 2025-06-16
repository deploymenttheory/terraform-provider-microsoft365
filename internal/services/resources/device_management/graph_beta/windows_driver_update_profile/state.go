package graphBetaWindowsDriverUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsDriverUpdateProfileResourceModel, remoteResource graphmodels.WindowsDriverUpdateProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": remoteResource.GetId()})

	// Scalar fields
	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.ApprovalType = state.EnumPtrToTypeString(remoteResource.GetApprovalType())
	data.DeviceReporting = state.Int32PtrToTypeInt32(remoteResource.GetDeviceReporting())
	data.NewUpdates = state.Int32PtrToTypeInt32(remoteResource.GetNewUpdates())
	data.DeploymentDeferralInDays = state.Int32PtrToTypeInt32(remoteResource.GetDeploymentDeferralInDays())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	// inventory_sync_status as an Object
	inventorySyncStatusTypes := map[string]attr.Type{
		"last_successful_sync_date_time": types.StringType,
		"driver_inventory_sync_state":    types.StringType,
	}

	if status := remoteResource.GetInventorySyncStatus(); status != nil {
		inventorySyncStatusValues := map[string]attr.Value{
			"last_successful_sync_date_time": state.TimeToString(status.GetLastSuccessfulSyncDateTime()),
			"driver_inventory_sync_state":    state.EnumPtrToTypeString(status.GetDriverInventorySyncState()),
		}
		data.InventorySyncStatus = state.ObjectValueMust(inventorySyncStatusTypes, inventorySyncStatusValues)
	} else {
		data.InventorySyncStatus = types.ObjectNull(inventorySyncStatusTypes)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
