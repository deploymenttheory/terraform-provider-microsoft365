package graphBetaDeviceAndAppManagementWindowsManagedAppProtection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the Graph API response to the Terraform state model.
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsManagedAppProtectionResourceModel, remoteResource graphmodels.WindowsManagedAppProtectionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Computed-only fields — always mapped from API, never written by user
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())

	if v := remoteResource.GetCreatedDateTime(); v != nil {
		data.CreatedDateTime = types.StringValue(v.Format("2006-01-02T15:04:05Z"))
	}
	if v := remoteResource.GetLastModifiedDateTime(); v != nil {
		data.LastModifiedDateTime = types.StringValue(v.Format("2006-01-02T15:04:05Z"))
	}
	if v := remoteResource.GetIsAssigned(); v != nil {
		data.IsAssigned = types.BoolValue(*v)
	}
	if v := remoteResource.GetDeployedAppCount(); v != nil {
		data.DeployedAppCount = types.Int64Value(int64(*v))
	}

	// Writable fields — mapped from API response to keep state in sync
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.MinimumRequiredSdkVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumRequiredSdkVersion())
	data.MinimumWipeSdkVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWipeSdkVersion())
	data.MinimumRequiredOsVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumRequiredOsVersion())
	data.MinimumWarningOsVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWarningOsVersion())
	data.MinimumWipeOsVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWipeOsVersion())
	data.MinimumRequiredAppVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumRequiredAppVersion())
	data.MinimumWarningAppVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWarningAppVersion())
	data.MinimumWipeAppVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWipeAppVersion())
	data.MaximumRequiredOsVersion = convert.GraphToFrameworkString(remoteResource.GetMaximumRequiredOsVersion())
	data.MaximumWarningOsVersion = convert.GraphToFrameworkString(remoteResource.GetMaximumWarningOsVersion())
	data.MaximumWipeOsVersion = convert.GraphToFrameworkString(remoteResource.GetMaximumWipeOsVersion())

	if v := remoteResource.GetPrintBlocked(); v != nil {
		data.PrintBlocked = types.BoolValue(*v)
	}

	// Duration fields — stored as ISO 8601 strings
	if v := remoteResource.GetPeriodOfflineBeforeWipeIsEnforced(); v != nil {
		data.PeriodOfflineBeforeWipeIsEnforced = types.StringValue(v.String())
	}
	if v := remoteResource.GetPeriodOfflineBeforeAccessCheck(); v != nil {
		data.PeriodOfflineBeforeAccessCheck = types.StringValue(v.String())
	}

	// Enum fields — convert Graph enum to string
	if v := remoteResource.GetAllowedInboundDataTransferSources(); v != nil {
		data.AllowedInboundDataTransferSources = types.StringValue(v.String())
	}
	if v := remoteResource.GetAllowedOutboundClipboardSharingLevel(); v != nil {
		data.AllowedOutboundClipboardSharingLevel = types.StringValue(v.String())
	}
	if v := remoteResource.GetAllowedOutboundDataTransferDestinations(); v != nil {
		data.AllowedOutboundDataTransferDestinations = types.StringValue(v.String())
	}
	if v := remoteResource.GetAppActionIfUnableToAuthenticateUser(); v != nil {
		data.AppActionIfUnableToAuthenticateUser = types.StringValue(v.String())
	}
	if v := remoteResource.GetMaximumAllowedDeviceThreatLevel(); v != nil {
		data.MaximumAllowedDeviceThreatLevel = types.StringValue(v.String())
	}
	if v := remoteResource.GetMobileThreatDefenseRemediationAction(); v != nil {
		data.MobileThreatDefenseRemediationAction = types.StringValue(v.String())
	}

	// List fields
	if tags := remoteResource.GetRoleScopeTagIds(); tags != nil {
		listVal, diags := types.ListValueFrom(ctx, types.StringType, tags)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to map role_scope_tag_ids from remote state")
		} else {
			data.RoleScopeTagIds = listVal
		}
	} else {
		data.RoleScopeTagIds, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
