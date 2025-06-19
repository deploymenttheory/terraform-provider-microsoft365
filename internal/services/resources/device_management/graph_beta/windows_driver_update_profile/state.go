package graphBetaWindowsDriverUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.ApprovalType = convert.GraphToFrameworkEnum(remoteResource.GetApprovalType())
	data.DeviceReporting = convert.GraphToFrameworkInt32(remoteResource.GetDeviceReporting())
	data.NewUpdates = convert.GraphToFrameworkInt32(remoteResource.GetNewUpdates())
	data.DeploymentDeferralInDays = convert.GraphToFrameworkInt32(remoteResource.GetDeploymentDeferralInDays())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// inventory_sync_status as an Object
	inventorySyncStatusTypes := map[string]attr.Type{
		"last_successful_sync_date_time": types.StringType,
		"driver_inventory_sync_state":    types.StringType,
	}

	if status := remoteResource.GetInventorySyncStatus(); status != nil {
		inventorySyncStatusValues := map[string]attr.Value{
			"last_successful_sync_date_time": convert.GraphToFrameworkTime(status.GetLastSuccessfulSyncDateTime()),
			"driver_inventory_sync_state":    convert.GraphToFrameworkEnum(status.GetDriverInventorySyncState()),
		}
		object, diags := types.ObjectValue(inventorySyncStatusTypes, inventorySyncStatusValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create object value", map[string]interface{}{
				"error": diags.Errors()[0].Detail(),
			})
			data.InventorySyncStatus = types.ObjectNull(inventorySyncStatusTypes)
		} else {
			data.InventorySyncStatus = object
		}
	} else {
		data.InventorySyncStatus = types.ObjectNull(inventorySyncStatusTypes)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
