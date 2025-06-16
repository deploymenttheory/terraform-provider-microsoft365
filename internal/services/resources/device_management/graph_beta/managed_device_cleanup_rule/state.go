package graphBetaManagedDeviceCleanupRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = state.StringPointerValue(remoteResource.GetDescription())
	data.DeviceInactivityBeforeRetirementInDays = state.Int32PtrToTypeInt32(remoteResource.GetDeviceInactivityBeforeRetirementInDays())
	data.DeviceCleanupRulePlatformType = state.EnumPtrToTypeString(remoteResource.GetDeviceCleanupRulePlatformType())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}
