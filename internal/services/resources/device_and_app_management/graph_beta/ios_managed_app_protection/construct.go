package graphBetaDeviceAndAppManagementIosManagedAppProtection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema model to the Graph SDK request body.
func constructResource(ctx context.Context, data *IosManagedAppProtectionResourceModel) (graphmodels.IosManagedAppProtectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewIosManagedAppProtection()

	// Required
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	// Optional strings — shared with Android
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.MinimumRequiredOsVersion, requestBody.SetMinimumRequiredOsVersion)
	convert.FrameworkToGraphString(data.MinimumWarningOsVersion, requestBody.SetMinimumWarningOsVersion)
	convert.FrameworkToGraphString(data.MinimumRequiredAppVersion, requestBody.SetMinimumRequiredAppVersion)
	convert.FrameworkToGraphString(data.MinimumWarningAppVersion, requestBody.SetMinimumWarningAppVersion)

	// Optional strings — iOS-specific
	convert.FrameworkToGraphString(data.CustomBrowserProtocol, requestBody.SetCustomBrowserProtocol)

	// Optional bools — shared with Android
	convert.FrameworkToGraphBool(data.OrganizationalCredentialsRequired, requestBody.SetOrganizationalCredentialsRequired)
	convert.FrameworkToGraphBool(data.DataBackupBlocked, requestBody.SetDataBackupBlocked)
	convert.FrameworkToGraphBool(data.DeviceComplianceRequired, requestBody.SetDeviceComplianceRequired)
	convert.FrameworkToGraphBool(data.ManagedBrowserToOpenLinksRequired, requestBody.SetManagedBrowserToOpenLinksRequired)
	convert.FrameworkToGraphBool(data.SaveAsBlocked, requestBody.SetSaveAsBlocked)
	convert.FrameworkToGraphBool(data.PinRequired, requestBody.SetPinRequired)
	convert.FrameworkToGraphBool(data.SimplePinBlocked, requestBody.SetSimplePinBlocked)
	convert.FrameworkToGraphBool(data.ContactSyncBlocked, requestBody.SetContactSyncBlocked)
	convert.FrameworkToGraphBool(data.PrintBlocked, requestBody.SetPrintBlocked)
	convert.FrameworkToGraphBool(data.FingerprintBlocked, requestBody.SetFingerprintBlocked)
	convert.FrameworkToGraphBool(data.DisableAppPinIfDevicePinIsSet, requestBody.SetDisableAppPinIfDevicePinIsSet)

	// Optional bools — iOS-specific
	convert.FrameworkToGraphBool(data.FaceIdBlocked, requestBody.SetFaceIdBlocked)
	convert.FrameworkToGraphBool(data.ThirdPartyKeyboardsBlocked, requestBody.SetThirdPartyKeyboardsBlocked)
	convert.FrameworkToGraphBool(data.FilterOpenInToOnlyManagedApps, requestBody.SetFilterOpenInToOnlyManagedApps)
	convert.FrameworkToGraphBool(data.DisableProtectionOfManagedOutboundOpenInData, requestBody.SetDisableProtectionOfManagedOutboundOpenInData)

	// Optional int32
	if !data.MaximumPinRetries.IsNull() && !data.MaximumPinRetries.IsUnknown() {
		v := data.MaximumPinRetries.ValueInt32()
		requestBody.SetMaximumPinRetries(&v)
	}
	if !data.MinimumPinLength.IsNull() && !data.MinimumPinLength.IsUnknown() {
		v := data.MinimumPinLength.ValueInt32()
		requestBody.SetMinimumPinLength(&v)
	}

	// Duration fields
	if !data.PeriodOfflineBeforeAccessCheck.IsNull() && !data.PeriodOfflineBeforeAccessCheck.IsUnknown() {
		duration, err := serialization.ParseISODuration(data.PeriodOfflineBeforeAccessCheck.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid period_offline_before_access_check value: %s", err)
		}
		requestBody.SetPeriodOfflineBeforeAccessCheck(duration)
	}
	if !data.PeriodOnlineBeforeAccessCheck.IsNull() && !data.PeriodOnlineBeforeAccessCheck.IsUnknown() {
		duration, err := serialization.ParseISODuration(data.PeriodOnlineBeforeAccessCheck.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid period_online_before_access_check value: %s", err)
		}
		requestBody.SetPeriodOnlineBeforeAccessCheck(duration)
	}
	if !data.PeriodOfflineBeforeWipeIsEnforced.IsNull() && !data.PeriodOfflineBeforeWipeIsEnforced.IsUnknown() {
		duration, err := serialization.ParseISODuration(data.PeriodOfflineBeforeWipeIsEnforced.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid period_offline_before_wipe_is_enforced value: %s", err)
		}
		requestBody.SetPeriodOfflineBeforeWipeIsEnforced(duration)
	}
	if !data.PeriodBeforePinReset.IsNull() && !data.PeriodBeforePinReset.IsUnknown() {
		duration, err := serialization.ParseISODuration(data.PeriodBeforePinReset.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid period_before_pin_reset value: %s", err)
		}
		requestBody.SetPeriodBeforePinReset(duration)
	}

	// Enum fields — shared with Android
	if !data.AllowedInboundDataTransferSources.IsNull() && !data.AllowedInboundDataTransferSources.IsUnknown() {
		val, err := graphmodels.ParseManagedAppDataTransferLevel(data.AllowedInboundDataTransferSources.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_inbound_data_transfer_sources value: %s", err)
		}
		requestBody.SetAllowedInboundDataTransferSources(val.(*graphmodels.ManagedAppDataTransferLevel))
	}
	if !data.AllowedOutboundDataTransferDestinations.IsNull() && !data.AllowedOutboundDataTransferDestinations.IsUnknown() {
		val, err := graphmodels.ParseManagedAppDataTransferLevel(data.AllowedOutboundDataTransferDestinations.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_outbound_data_transfer_destinations value: %s", err)
		}
		requestBody.SetAllowedOutboundDataTransferDestinations(val.(*graphmodels.ManagedAppDataTransferLevel))
	}
	if !data.AllowedOutboundClipboardSharingLevel.IsNull() && !data.AllowedOutboundClipboardSharingLevel.IsUnknown() {
		val, err := graphmodels.ParseManagedAppClipboardSharingLevel(data.AllowedOutboundClipboardSharingLevel.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_outbound_clipboard_sharing_level value: %s", err)
		}
		requestBody.SetAllowedOutboundClipboardSharingLevel(val.(*graphmodels.ManagedAppClipboardSharingLevel))
	}
	if !data.PinCharacterSet.IsNull() && !data.PinCharacterSet.IsUnknown() {
		val, err := graphmodels.ParseManagedAppPinCharacterSet(data.PinCharacterSet.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid pin_character_set value: %s", err)
		}
		requestBody.SetPinCharacterSet(val.(*graphmodels.ManagedAppPinCharacterSet))
	}
	if !data.ManagedBrowser.IsNull() && !data.ManagedBrowser.IsUnknown() {
		val, err := graphmodels.ParseManagedBrowserType(data.ManagedBrowser.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid managed_browser value: %s", err)
		}
		requestBody.SetManagedBrowser(val.(*graphmodels.ManagedBrowserType))
	}

	// Enum fields — iOS-specific
	if !data.AppDataEncryptionType.IsNull() && !data.AppDataEncryptionType.IsUnknown() {
		val, err := graphmodels.ParseManagedAppDataEncryptionType(data.AppDataEncryptionType.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid app_data_encryption_type value: %s", err)
		}
		requestBody.SetAppDataEncryptionType(val.(*graphmodels.ManagedAppDataEncryptionType))
	}

	// List fields
	if !data.AllowedDataStorageLocations.IsNull() && !data.AllowedDataStorageLocations.IsUnknown() {
		var locations []string
		data.AllowedDataStorageLocations.ElementsAs(ctx, &locations, false)
		parsedLocations := make([]graphmodels.ManagedAppDataStorageLocation, 0, len(locations))
		for _, loc := range locations {
			val, err := graphmodels.ParseManagedAppDataStorageLocation(loc)
			if err != nil {
				return nil, fmt.Errorf("invalid allowed_data_storage_locations value '%s': %s", loc, err)
			}
			parsedLocations = append(parsedLocations, *val.(*graphmodels.ManagedAppDataStorageLocation))
		}
		requestBody.SetAllowedDataStorageLocations(parsedLocations)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}