// MapRemoteResourceStateToTerraform maps the base properties of a WindowsDriverUpdateProfileResourceModel to a Terraform state
package graphBetaWindowsDriverUpdateProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the base properties of a WindowsDriverUpdateProfileResourceModel to a Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsDriverUpdateProfileResourceModel, remoteResource graphmodels.WindowsDriverUpdateProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.ApprovalType = state.EnumPtrToTypeString(remoteResource.GetApprovalType())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.DeviceReporting = state.Int32PtrToTypeInt32(remoteResource.GetDeviceReporting())
	data.NewUpdates = state.Int32PtrToTypeInt32(remoteResource.GetNewUpdates())
	data.DeploymentDeferralInDays = state.Int32PtrToTypeInt32(remoteResource.GetDeploymentDeferralInDays())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Handle inventory sync status
	if inventorySyncStatus := remoteResource.GetInventorySyncStatus(); inventorySyncStatus != nil {
		if data.InventorySyncStatus == nil {
			data.InventorySyncStatus = &WindowsDriverUpdateProfileInventorySyncStatus{}
		}

		data.InventorySyncStatus.LastSuccessfulSyncDateTime = state.TimeToString(inventorySyncStatus.GetLastSuccessfulSyncDateTime())
		data.InventorySyncStatus.DriverInventorySyncState = state.EnumPtrToTypeString(inventorySyncStatus.GetDriverInventorySyncState())
	}

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
