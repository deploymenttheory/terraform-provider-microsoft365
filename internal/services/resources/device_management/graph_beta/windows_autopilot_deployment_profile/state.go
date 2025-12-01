package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote Windows Autopilot Deployment Profile resource state to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsAutopilotDeploymentProfileResourceModel, remoteResource graphmodels.WindowsAutopilotDeploymentProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map Windows Autopilot Deployment Profile from API to Terraform state")

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())

	// Handle locale mapping - convert empty string from API back to "user_select"
	if locale := remoteResource.GetLocale(); locale != nil && *locale == "" {
		data.Locale = types.StringValue("user_select")
	} else {
		data.Locale = convert.GraphToFrameworkString(remoteResource.GetLocale())
	}
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.HardwareHashExtractionEnabled = convert.GraphToFrameworkBool(remoteResource.GetHardwareHashExtractionEnabled())
	data.DeviceNameTemplate = convert.GraphToFrameworkString(remoteResource.GetDeviceNameTemplate())
	data.PreprovisioningAllowed = convert.GraphToFrameworkBool(remoteResource.GetPreprovisioningAllowed())
	data.ManagementServiceAppId = convert.GraphToFrameworkString(remoteResource.GetManagementServiceAppId())

	// Determine device join type based on the actual profile type returned from API
	if _, ok := remoteResource.(graphmodels.ActiveDirectoryWindowsAutopilotDeploymentProfileable); ok {
		data.DeviceJoinType = types.StringValue("microsoft_entra_hybrid_joined")
	} else if _, ok := remoteResource.(graphmodels.AzureADWindowsAutopilotDeploymentProfileable); ok {
		data.DeviceJoinType = types.StringValue("microsoft_entra_joined")
	} else {
		data.DeviceJoinType = types.StringNull()
	}

	// Handle hybrid_azure_ad_join_skip_connectivity_check
	// This property only exists on activeDirectoryWindowsAutopilotDeploymentProfile (Hybrid AD joined)
	// For azureADWindowsAutopilotDeploymentProfile (Azure AD joined), the API does not return this field
	if adResource, ok := remoteResource.(graphmodels.ActiveDirectoryWindowsAutopilotDeploymentProfileable); ok {
		// Hybrid AD joined profile - read from the typed property
		data.HybridAzureADJoinSkipConnectivityCheck = convert.GraphToFrameworkBool(adResource.GetHybridAzureADJoinSkipConnectivityCheck())
	} else {
		// Azure AD joined profile - this field doesn't exist in the API response
		// Keep the configured value from state to avoid inconsistency errors
		// The validator ensures this is always false for Azure AD joined profiles
		if data.HybridAzureADJoinSkipConnectivityCheck.IsNull() || data.HybridAzureADJoinSkipConnectivityCheck.IsUnknown() {
			data.HybridAzureADJoinSkipConnectivityCheck = types.BoolValue(false)
		}
		// Otherwise preserve the existing value from config/state
	}

	if deviceType := remoteResource.GetDeviceType(); deviceType != nil {
		data.DeviceType = convert.GraphToFrameworkEnum(deviceType)
	}

	if roleScopeTagIds := remoteResource.GetRoleScopeTagIds(); roleScopeTagIds != nil {
		data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, roleScopeTagIds)
	} else {
		data.RoleScopeTagIds = types.SetNull(types.StringType)
	}

	if oobeSetting := remoteResource.GetOutOfBoxExperienceSetting(); oobeSetting != nil {
		data.OutOfBoxExperienceSetting = &OutOfBoxExperienceSettingModel{
			PrivacySettingsHidden:        convert.GraphToFrameworkBool(oobeSetting.GetPrivacySettingsHidden()),
			EulaHidden:                   convert.GraphToFrameworkBool(oobeSetting.GetEulaHidden()),
			KeyboardSelectionPageSkipped: convert.GraphToFrameworkBool(oobeSetting.GetKeyboardSelectionPageSkipped()),
			EscapeLinkHidden:             convert.GraphToFrameworkBool(oobeSetting.GetEscapeLinkHidden()),
		}

		if userType := oobeSetting.GetUserType(); userType != nil {
			data.OutOfBoxExperienceSetting.UserType = convert.GraphToFrameworkEnum(userType)
		}

		if deviceUsageType := oobeSetting.GetDeviceUsageType(); deviceUsageType != nil {
			data.OutOfBoxExperienceSetting.DeviceUsageType = convert.GraphToFrameworkEnum(deviceUsageType)
		}
	}

	tflog.Debug(ctx, "Finished mapping Windows Autopilot Deployment Profile from API to Terraform state")
}
