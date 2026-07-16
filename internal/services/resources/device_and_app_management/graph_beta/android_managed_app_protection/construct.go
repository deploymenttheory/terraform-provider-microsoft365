package graphBetaDeviceAndAppManagementAndroidManagedAppProtection

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
func constructResource(ctx context.Context, data *AndroidManagedAppProtectionResourceModel) (graphmodels.AndroidManagedAppProtectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewAndroidManagedAppProtection()

	// Required
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	// Optional strings
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.MinimumRequiredOsVersion, requestBody.SetMinimumRequiredOsVersion)
	convert.FrameworkToGraphString(data.MinimumWarningOsVersion, requestBody.SetMinimumWarningOsVersion)
	convert.FrameworkToGraphString(data.MinimumRequiredAppVersion, requestBody.SetMinimumRequiredAppVersion)
	convert.FrameworkToGraphString(data.MinimumWarningAppVersion, requestBody.SetMinimumWarningAppVersion)
	convert.FrameworkToGraphString(data.MinimumRequiredPatchVersion, requestBody.SetMinimumRequiredPatchVersion)
	convert.FrameworkToGraphString(data.MinimumWarningPatchVersion, requestBody.SetMinimumWarningPatchVersion)
	convert.FrameworkToGraphString(data.CustomBrowserPackageId, requestBody.SetCustomBrowserPackageId)
	convert.FrameworkToGraphString(data.CustomBrowserDisplayName, requestBody.SetCustomBrowserDisplayName)

	// Optional bools
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
	convert.FrameworkToGraphBool(data.ScreenCaptureBlocked, requestBody.SetScreenCaptureBlocked)
	convert.FrameworkToGraphBool(data.DisableAppEncryptionIfDeviceEncryptionIsEnabled, requestBody.SetDisableAppEncryptionIfDeviceEncryptionIsEnabled)
	convert.FrameworkToGraphBool(data.EncryptAppData, requestBody.SetEncryptAppData)

	// Optional int64
	if !data.MaximumPinRetries.IsNull() && !data.MaximumPinRetries.IsUnknown() {
		v := int32(data.MaximumPinRetries.ValueInt64())
		requestBody.SetMaximumPinRetries(&v)
	}
	if !data.MinimumPinLength.IsNull() && !data.MinimumPinLength.IsUnknown() {
		v := int32(data.MinimumPinLength.ValueInt64())
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

	// Enum fields
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
