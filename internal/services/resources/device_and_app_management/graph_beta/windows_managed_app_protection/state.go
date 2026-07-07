package graphBetaDeviceAndAppManagementWindowsManagedAppProtection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the Graph API response to the Terraform state model.
// Every optional field that the API may return as null is explicitly set to types.StringNull()
// or types.BoolNull() in the else branch — this prevents Terraform from seeing "unknown"
// values after apply, which would cause a provider error.
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsManagedAppProtectionResourceModel, remoteResource graphmodels.WindowsManagedAppProtectionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// --- Computed-only fields ---
	// Always set from API, never written by the user. Safe to use convert directly
	// as id and version are always present on a valid response.
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())

	if v := remoteResource.GetCreatedDateTime(); v != nil {
		data.CreatedDateTime = types.StringValue(v.Format("2006-01-02T15:04:05Z"))
	} else {
		data.CreatedDateTime = types.StringNull()
	}

	if v := remoteResource.GetLastModifiedDateTime(); v != nil {
		data.LastModifiedDateTime = types.StringValue(v.Format("2006-01-02T15:04:05Z"))
	} else {
		data.LastModifiedDateTime = types.StringNull()
	}

	if v := remoteResource.GetIsAssigned(); v != nil {
		data.IsAssigned = types.BoolValue(*v)
	} else {
		data.IsAssigned = types.BoolValue(false)
	}

	if v := remoteResource.GetDeployedAppCount(); v != nil {
		data.DeployedAppCount = types.Int64Value(int64(*v))
	} else {
		data.DeployedAppCount = types.Int64Value(0)
	}

	// --- Required string fields ---
	// display_name is required so will always be present, but still use convert for safety.
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())

	// --- Optional plain string fields ---
	// All version/OS strings are optional — API returns null when not configured.
	// convert.GraphToFrameworkString already handles nil pointers by returning types.StringNull(),
	// so these are safe to use directly without an explicit else branch.
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

	// --- Optional bool fields ---
	if v := remoteResource.GetPrintBlocked(); v != nil {
		data.PrintBlocked = types.BoolValue(*v)
	} else {
		data.PrintBlocked = types.BoolValue(false)
	}

	// --- Duration fields ---
	// API returns these as ISODuration objects. Convert to string for Terraform state.
	// Explicitly null when not set by the API.
	if v := remoteResource.GetPeriodOfflineBeforeWipeIsEnforced(); v != nil {
		data.PeriodOfflineBeforeWipeIsEnforced = types.StringValue(v.String())
	} else {
		data.PeriodOfflineBeforeWipeIsEnforced = types.StringNull()
	}

	if v := remoteResource.GetPeriodOfflineBeforeAccessCheck(); v != nil {
		data.PeriodOfflineBeforeAccessCheck = types.StringValue(v.String())
	} else {
		data.PeriodOfflineBeforeAccessCheck = types.StringNull()
	}

	// --- Enum fields ---
	// Every enum field is explicitly nulled in the else branch so Terraform never
	// sees an unknown value after apply, which would cause a provider error.
	if v := remoteResource.GetAllowedInboundDataTransferSources(); v != nil {
		data.AllowedInboundDataTransferSources = types.StringValue(v.String())
	} else {
		data.AllowedInboundDataTransferSources = types.StringNull()
	}

	if v := remoteResource.GetAllowedOutboundClipboardSharingLevel(); v != nil {
		data.AllowedOutboundClipboardSharingLevel = types.StringValue(v.String())
	} else {
		data.AllowedOutboundClipboardSharingLevel = types.StringNull()
	}

	if v := remoteResource.GetAllowedOutboundDataTransferDestinations(); v != nil {
		data.AllowedOutboundDataTransferDestinations = types.StringValue(v.String())
	} else {
		data.AllowedOutboundDataTransferDestinations = types.StringNull()
	}

	// This was the field that caused the original provider error — API returns null
	// when no action is configured, which must be mapped to types.StringNull() explicitly.
	if v := remoteResource.GetAppActionIfUnableToAuthenticateUser(); v != nil {
		data.AppActionIfUnableToAuthenticateUser = types.StringValue(v.String())
	} else {
		data.AppActionIfUnableToAuthenticateUser = types.StringNull()
	}

	if v := remoteResource.GetMaximumAllowedDeviceThreatLevel(); v != nil {
		data.MaximumAllowedDeviceThreatLevel = types.StringValue(v.String())
	} else {
		data.MaximumAllowedDeviceThreatLevel = types.StringNull()
	}

	if v := remoteResource.GetMobileThreatDefenseRemediationAction(); v != nil {
		data.MobileThreatDefenseRemediationAction = types.StringValue(v.String())
	} else {
		data.MobileThreatDefenseRemediationAction = types.StringNull()
	}

	// --- List fields ---
	// RoleScopeTagIds returns an empty slice rather than nil when no tags are set,
	// but we handle both cases to be safe.
	if tags := remoteResource.GetRoleScopeTagIds(); tags != nil {
		listVal, diags := types.ListValueFrom(ctx, types.StringType, tags)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to map role_scope_tag_ids from remote state")
			data.RoleScopeTagIds, _ = types.ListValue(types.StringType, []attr.Value{})
		} else {
			data.RoleScopeTagIds = listVal
		}
	} else {
		data.RoleScopeTagIds, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
