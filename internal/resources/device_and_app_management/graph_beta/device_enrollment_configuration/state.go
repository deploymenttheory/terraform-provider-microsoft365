package graphBetaDeviceEnrollmentConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": remoteResource.GetId()})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Priority = state.Int32PtrToTypeInt64(remoteResource.GetPriority())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.Version = state.Int32PtrToTypeInt64(remoteResource.GetVersion())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Map device enrollment configuration type if available
	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		data.DeviceEnrollmentConfigurationType = types.StringValue(string(*configType))
	}

	// Map platform restrictions if available
	if platformRestriction := remoteResource.GetPlatformRestriction(); platformRestriction != nil {
		data.PlatformRestriction = &PlatformRestrictionModel{
			PlatformBlocked:                 state.BoolPtrToTypeBool(platformRestriction.GetPlatformBlocked()),
			PersonalDeviceEnrollmentBlocked: state.BoolPtrToTypeBool(platformRestriction.GetPersonalDeviceEnrollmentBlocked()),
			OSMinimumVersion:                types.StringPointerValue(platformRestriction.GetOsMinimumVersion()),
			OSMaximumVersion:                types.StringPointerValue(platformRestriction.GetOsMaximumVersion()),
			BlockedManufacturers:            state.StringSliceToSet(ctx, platformRestriction.GetBlockedManufacturers()),
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform", map[string]interface{}{"resourceId": data.ID.ValueString()})
}
