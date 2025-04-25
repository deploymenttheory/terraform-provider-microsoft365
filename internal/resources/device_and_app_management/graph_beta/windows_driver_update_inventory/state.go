// MapRemoteResourceStateToTerraform maps the base properties of a WindowsDriverUpdateInventoryResourceModel to a Terraform state
package graphBetaWindowsDriverUpdateInventory

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the base properties of a WindowsDriverUpdateInventoryResourceModel to a Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsDriverUpdateInventoryResourceModel, remoteResource graphmodels.WindowsDriverUpdateInventoryable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.Name = types.StringPointerValue(remoteResource.GetName())
	data.Version = types.StringPointerValue(remoteResource.GetVersion())
	data.Manufacturer = types.StringPointerValue(remoteResource.GetManufacturer())
	data.DriverClass = types.StringPointerValue(remoteResource.GetDriverClass())
	data.ApprovalStatus = state.EnumPtrToTypeString(remoteResource.GetApprovalStatus())
	data.Category = state.EnumPtrToTypeString(remoteResource.GetCategory())
	data.ReleaseDateTime = state.TimeToString(remoteResource.GetReleaseDateTime())
	data.DeployDateTime = state.TimeToString(remoteResource.GetDeployDateTime())
	data.ApplicableDeviceCount = state.Int32PtrToTypeInt32(remoteResource.GetApplicableDeviceCount())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
