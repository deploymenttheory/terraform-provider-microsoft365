package graphBetaManagedDeviceCleanupRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data ManagedDeviceCleanupRuleResourceModel, remoteResource graphmodels.ManagedDeviceCleanupRuleable) ManagedDeviceCleanupRuleResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state")

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.DeviceInactivityBeforeRetirementInDays = convert.GraphToFrameworkInt32(remoteResource.GetDeviceInactivityBeforeRetirementInDays())
	data.DeviceCleanupRulePlatformType = convert.GraphToFrameworkEnum(remoteResource.GetDeviceCleanupRulePlatformType())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}
