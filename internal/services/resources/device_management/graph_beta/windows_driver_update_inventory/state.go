// MapRemoteResourceStateToTerraform maps the base properties of a WindowsDriverUpdateInventoryResourceModel to a Terraform state
package graphBetaWindowsDriverUpdateInventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())
	data.Manufacturer = convert.GraphToFrameworkString(remoteResource.GetManufacturer())
	data.DriverClass = convert.GraphToFrameworkString(remoteResource.GetDriverClass())
	data.ApprovalStatus = convert.GraphToFrameworkEnum(remoteResource.GetApprovalStatus())
	data.Category = convert.GraphToFrameworkEnum(remoteResource.GetCategory())
	data.ReleaseDateTime = convert.GraphToFrameworkTime(remoteResource.GetReleaseDateTime())
	data.DeployDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeployDateTime())
	data.ApplicableDeviceCount = convert.GraphToFrameworkInt32(remoteResource.GetApplicableDeviceCount())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
