package graphBetaDeviceAndAppManagementIosManagedAppProtection

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
func MapRemoteStateToTerraform(ctx context.Context, data *IosManagedAppProtectionResourceModel, remoteResource graphmodels.IosManagedAppProtectionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	//Computed-only fields
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
		data.DeployedAppCount = types.Int32Value(*v)
	} else {
		data.DeployedAppCount = types.Int32Value(0)
	}

	//Required
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())

	// --- Optional plain strings — shared with Android ---
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.MinimumRequiredOsVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumRequiredOsVersion())
	data.MinimumWarningOsVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWarningOsVersion())
	data.MinimumRequiredAppVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumRequiredAppVersion())
	data.MinimumWarningAppVersion = convert.GraphToFrameworkString(remoteResource.GetMinimumWarningAppVersion())

	//Optional plain strings — iOS-specific
	data.CustomBrowserProtocol = convert.GraphToFrameworkString(remoteResource.GetCustomBrowserProtocol())

	//Optional bools — shared with Android
	if v := remoteResource.GetOrganizationalCredentialsRequired(); v != nil {
		data.OrganizationalCredentialsRequired = types.BoolValue(*v)
	} else {
		data.OrganizationalCredentialsRequired = types.BoolValue(false)
	}
	if v := remoteResource.GetDataBackupBlocked(); v != nil {
		data.DataBackupBlocked = types.BoolValue(*v)
	} else {
		data.DataBackupBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetDeviceComplianceRequired(); v != nil {
		data.DeviceComplianceRequired = types.BoolValue(*v)
	} else {
		data.DeviceComplianceRequired = types.BoolValue(true)
	}
	if v := remoteResource.GetManagedBrowserToOpenLinksRequired(); v != nil {
		data.ManagedBrowserToOpenLinksRequired = types.BoolValue(*v)
	} else {
		data.ManagedBrowserToOpenLinksRequired = types.BoolValue(false)
	}
	if v := remoteResource.GetSaveAsBlocked(); v != nil {
		data.SaveAsBlocked = types.BoolValue(*v)
	} else {
		data.SaveAsBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetPinRequired(); v != nil {
		data.PinRequired = types.BoolValue(*v)
	} else {
		data.PinRequired = types.BoolValue(true)
	}
	if v := remoteResource.GetSimplePinBlocked(); v != nil {
		data.SimplePinBlocked = types.BoolValue(*v)
	} else {
		data.SimplePinBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetContactSyncBlocked(); v != nil {
		data.ContactSyncBlocked = types.BoolValue(*v)
	} else {
		data.ContactSyncBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetPrintBlocked(); v != nil {
		data.PrintBlocked = types.BoolValue(*v)
	} else {
		data.PrintBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetFingerprintBlocked(); v != nil {
		data.FingerprintBlocked = types.BoolValue(*v)
	} else {
		data.FingerprintBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetDisableAppPinIfDevicePinIsSet(); v != nil {
		data.DisableAppPinIfDevicePinIsSet = types.BoolValue(*v)
	} else {
		data.DisableAppPinIfDevicePinIsSet = types.BoolValue(false)
	}

	//Optional bools — iOS-specific
	if v := remoteResource.GetFaceIdBlocked(); v != nil {
		data.FaceIdBlocked = types.BoolValue(*v)
	} else {
		data.FaceIdBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetThirdPartyKeyboardsBlocked(); v != nil {
		data.ThirdPartyKeyboardsBlocked = types.BoolValue(*v)
	} else {
		data.ThirdPartyKeyboardsBlocked = types.BoolValue(false)
	}
	if v := remoteResource.GetFilterOpenInToOnlyManagedApps(); v != nil {
		data.FilterOpenInToOnlyManagedApps = types.BoolValue(*v)
	} else {
		data.FilterOpenInToOnlyManagedApps = types.BoolValue(false)
	}
	if v := remoteResource.GetDisableProtectionOfManagedOutboundOpenInData(); v != nil {
		data.DisableProtectionOfManagedOutboundOpenInData = types.BoolValue(*v)
	} else {
		data.DisableProtectionOfManagedOutboundOpenInData = types.BoolValue(false)
	}

	//Optional int32
	if v := remoteResource.GetMaximumPinRetries(); v != nil {
		data.MaximumPinRetries = types.Int32Value(*v)
	} else {
		data.MaximumPinRetries = types.Int32Value(5)
	}
	if v := remoteResource.GetMinimumPinLength(); v != nil {
		data.MinimumPinLength = types.Int32Value(*v)
	} else {
		data.MinimumPinLength = types.Int32Value(4)
	}

	//Duration fields
	if v := remoteResource.GetPeriodOfflineBeforeAccessCheck(); v != nil {
		data.PeriodOfflineBeforeAccessCheck = types.StringValue(v.String())
	} else {
		data.PeriodOfflineBeforeAccessCheck = types.StringNull()
	}
	if v := remoteResource.GetPeriodOnlineBeforeAccessCheck(); v != nil {
		data.PeriodOnlineBeforeAccessCheck = types.StringValue(v.String())
	} else {
		data.PeriodOnlineBeforeAccessCheck = types.StringNull()
	}
	if v := remoteResource.GetPeriodOfflineBeforeWipeIsEnforced(); v != nil {
		data.PeriodOfflineBeforeWipeIsEnforced = types.StringValue(v.String())
	} else {
		data.PeriodOfflineBeforeWipeIsEnforced = types.StringNull()
	}
	if v := remoteResource.GetPeriodBeforePinReset(); v != nil {
		data.PeriodBeforePinReset = types.StringValue(v.String())
	} else {
		data.PeriodBeforePinReset = types.StringNull()
	}

	//Enum fields — shared with Android
	if v := remoteResource.GetAllowedInboundDataTransferSources(); v != nil {
		data.AllowedInboundDataTransferSources = types.StringValue(v.String())
	} else {
		data.AllowedInboundDataTransferSources = types.StringNull()
	}
	if v := remoteResource.GetAllowedOutboundDataTransferDestinations(); v != nil {
		data.AllowedOutboundDataTransferDestinations = types.StringValue(v.String())
	} else {
		data.AllowedOutboundDataTransferDestinations = types.StringNull()
	}
	if v := remoteResource.GetAllowedOutboundClipboardSharingLevel(); v != nil {
		data.AllowedOutboundClipboardSharingLevel = types.StringValue(v.String())
	} else {
		data.AllowedOutboundClipboardSharingLevel = types.StringNull()
	}
	if v := remoteResource.GetPinCharacterSet(); v != nil {
		data.PinCharacterSet = types.StringValue(v.String())
	} else {
		data.PinCharacterSet = types.StringNull()
	}
	if v := remoteResource.GetManagedBrowser(); v != nil {
		data.ManagedBrowser = types.StringValue(v.String())
	} else {
		data.ManagedBrowser = types.StringNull()
	}

	//Enum fields — iOS-specific
	if v := remoteResource.GetAppDataEncryptionType(); v != nil {
		data.AppDataEncryptionType = types.StringValue(v.String())
	} else {
		data.AppDataEncryptionType = types.StringNull()
	}

	//List fields
	if locations := remoteResource.GetAllowedDataStorageLocations(); locations != nil {
		locationStrings := make([]attr.Value, 0, len(locations))
		for _, loc := range locations {
			locationStrings = append(locationStrings, types.StringValue(loc.String()))
		}
		listVal, diags := types.ListValue(types.StringType, locationStrings)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to map allowed_data_storage_locations from remote state")
			data.AllowedDataStorageLocations, _ = types.ListValue(types.StringType, []attr.Value{})
		} else {
			data.AllowedDataStorageLocations = listVal
		}
	} else {
		data.AllowedDataStorageLocations, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}