package graphBetaDeviceAndAppManagementWindowsManagedAppProtection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema model to the Graph SDK request body.
func constructResource(ctx context.Context, data *WindowsManagedAppProtectionResourceModel) (graphmodels.WindowsManagedAppProtectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsManagedAppProtection()

	// Required
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	// Optional strings
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.MinimumRequiredSdkVersion, requestBody.SetMinimumRequiredSdkVersion)
	convert.FrameworkToGraphString(data.MinimumWipeSdkVersion, requestBody.SetMinimumWipeSdkVersion)
	convert.FrameworkToGraphString(data.MinimumRequiredOsVersion, requestBody.SetMinimumRequiredOsVersion)
	convert.FrameworkToGraphString(data.MinimumWarningOsVersion, requestBody.SetMinimumWarningOsVersion)
	convert.FrameworkToGraphString(data.MinimumWipeOsVersion, requestBody.SetMinimumWipeOsVersion)
	convert.FrameworkToGraphString(data.MinimumRequiredAppVersion, requestBody.SetMinimumRequiredAppVersion)
	convert.FrameworkToGraphString(data.MinimumWarningAppVersion, requestBody.SetMinimumWarningAppVersion)
	convert.FrameworkToGraphString(data.MinimumWipeAppVersion, requestBody.SetMinimumWipeAppVersion)
	convert.FrameworkToGraphString(data.MaximumRequiredOsVersion, requestBody.SetMaximumRequiredOsVersion)
	convert.FrameworkToGraphString(data.MaximumWarningOsVersion, requestBody.SetMaximumWarningOsVersion)
	convert.FrameworkToGraphString(data.MaximumWipeOsVersion, requestBody.SetMaximumWipeOsVersion)
	convert.FrameworkToGraphString(data.PeriodOfflineBeforeWipeIsEnforced, requestBody.SetPeriodOfflineBeforeWipeIsEnforced)
	convert.FrameworkToGraphString(data.PeriodOfflineBeforeAccessCheck, requestBody.SetPeriodOfflineBeforeAccessCheck)

	// Optional bools
	convert.FrameworkToGraphBool(data.PrintBlocked, requestBody.SetPrintBlocked)

	// Optional enums — AllowedInboundDataTransferSources
	if !data.AllowedInboundDataTransferSources.IsNull() && !data.AllowedInboundDataTransferSources.IsUnknown() {
		val, err := graphmodels.ParseWindowsManagedAppDataTransferLevel(data.AllowedInboundDataTransferSources.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_inbound_data_transfer_sources value: %s", err)
		}
		requestBody.SetAllowedInboundDataTransferSources(val.(*graphmodels.WindowsManagedAppDataTransferLevel))
	}

	// Optional enums — AllowedOutboundClipboardSharingLevel
	if !data.AllowedOutboundClipboardSharingLevel.IsNull() && !data.AllowedOutboundClipboardSharingLevel.IsUnknown() {
		val, err := graphmodels.ParseWindowsManagedAppClipboardSharingLevel(data.AllowedOutboundClipboardSharingLevel.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_outbound_clipboard_sharing_level value: %s", err)
		}
		requestBody.SetAllowedOutboundClipboardSharingLevel(val.(*graphmodels.WindowsManagedAppClipboardSharingLevel))
	}

	// Optional enums — AllowedOutboundDataTransferDestinations
	if !data.AllowedOutboundDataTransferDestinations.IsNull() && !data.AllowedOutboundDataTransferDestinations.IsUnknown() {
		val, err := graphmodels.ParseWindowsManagedAppDataTransferLevel(data.AllowedOutboundDataTransferDestinations.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid allowed_outbound_data_transfer_destinations value: %s", err)
		}
		requestBody.SetAllowedOutboundDataTransferDestinations(val.(*graphmodels.WindowsManagedAppDataTransferLevel))
	}

	// Optional enums — AppActionIfUnableToAuthenticateUser
	if !data.AppActionIfUnableToAuthenticateUser.IsNull() && !data.AppActionIfUnableToAuthenticateUser.IsUnknown() {
		val, err := graphmodels.ParseManagedAppRemediationAction(data.AppActionIfUnableToAuthenticateUser.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid app_action_if_unable_to_authenticate_user value: %s", err)
		}
		requestBody.SetAppActionIfUnableToAuthenticateUser(val.(*graphmodels.ManagedAppRemediationAction))
	}

	// Optional enums — MaximumAllowedDeviceThreatLevel
	if !data.MaximumAllowedDeviceThreatLevel.IsNull() && !data.MaximumAllowedDeviceThreatLevel.IsUnknown() {
		val, err := graphmodels.ParseManagedAppDeviceThreatLevel(data.MaximumAllowedDeviceThreatLevel.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid maximum_allowed_device_threat_level value: %s", err)
		}
		requestBody.SetMaximumAllowedDeviceThreatLevel(val.(*graphmodels.ManagedAppDeviceThreatLevel))
	}

	// Optional enums — MobileThreatDefenseRemediationAction
	if !data.MobileThreatDefenseRemediationAction.IsNull() && !data.MobileThreatDefenseRemediationAction.IsUnknown() {
		val, err := graphmodels.ParseManagedAppRemediationAction(data.MobileThreatDefenseRemediationAction.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid mobile_threat_defense_remediation_action value: %s", err)
		}
		requestBody.SetMobileThreatDefenseRemediationAction(val.(*graphmodels.ManagedAppRemediationAction))
	}

	// Optional list — RoleScopeTagIds
	if !data.RoleScopeTagIds.IsNull() && !data.RoleScopeTagIds.IsUnknown() {
		var tagIds []string
		data.RoleScopeTagIds.ElementsAs(ctx, &tagIds, false)
		requestBody.SetRoleScopeTagIds(tagIds)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
