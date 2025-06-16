package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
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

	data.ID = types.StringValue(*remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = state.StringPointerValue(remoteResource.GetDescription())
	data.Language = state.StringPointerValue(remoteResource.GetLanguage())
	data.Locale = state.StringPointerValue(remoteResource.GetLocale())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.HardwareHashExtractionEnabled = state.BoolPtrToTypeBool(remoteResource.GetHardwareHashExtractionEnabled())
	data.DeviceNameTemplate = state.StringPointerValue(remoteResource.GetDeviceNameTemplate())
	data.PreprovisioningAllowed = state.BoolPtrToTypeBool(remoteResource.GetPreprovisioningAllowed())
	data.ManagementServiceAppId = state.StringPointerValue(remoteResource.GetManagementServiceAppId())

	// Determine device join type based on the actual profile type returned from API
	if _, ok := remoteResource.(graphmodels.ActiveDirectoryWindowsAutopilotDeploymentProfileable); ok {
		data.DeviceJoinType = types.StringValue("microsoft_entra_hybrid_joined")
	} else if _, ok := remoteResource.(graphmodels.AzureADWindowsAutopilotDeploymentProfileable); ok {
		data.DeviceJoinType = types.StringValue("microsoft_entra_joined")
	} else {
		data.DeviceJoinType = types.StringNull()
	}

	// Check if this is an ActiveDirectoryWindowsAutopilotDeploymentProfile and handle hybrid Azure AD join setting
	if adResource, ok := remoteResource.(graphmodels.ActiveDirectoryWindowsAutopilotDeploymentProfileable); ok {
		data.HybridAzureADJoinSkipConnectivityCheck = state.BoolPtrToTypeBool(adResource.GetHybridAzureADJoinSkipConnectivityCheck())
	} else {
		// For Azure AD profiles, this field is not applicable but we should preserve the configured value
		// to avoid Terraform inconsistency errors. The field is not sent to the API for Azure AD profiles.
		if !data.HybridAzureADJoinSkipConnectivityCheck.IsUnknown() {
			// Keep the existing value from configuration
		} else {
			data.HybridAzureADJoinSkipConnectivityCheck = types.BoolNull()
		}
	}

	if deviceType := remoteResource.GetDeviceType(); deviceType != nil {
		data.DeviceType = state.EnumPtrToTypeString(deviceType)
	}

	if roleScopeTagIds := remoteResource.GetRoleScopeTagIds(); roleScopeTagIds != nil {
		data.RoleScopeTagIds = state.StringSliceToSet(ctx, roleScopeTagIds)
	} else {
		data.RoleScopeTagIds = types.SetNull(types.StringType)
	}

	if oobeSetting := remoteResource.GetOutOfBoxExperienceSetting(); oobeSetting != nil {
		data.OutOfBoxExperienceSetting = &OutOfBoxExperienceSettingModel{
			PrivacySettingsHidden:        state.BoolPtrToTypeBool(oobeSetting.GetPrivacySettingsHidden()),
			EulaHidden:                   state.BoolPtrToTypeBool(oobeSetting.GetEulaHidden()),
			KeyboardSelectionPageSkipped: state.BoolPtrToTypeBool(oobeSetting.GetKeyboardSelectionPageSkipped()),
			EscapeLinkHidden:             state.BoolPtrToTypeBool(oobeSetting.GetEscapeLinkHidden()),
		}

		if userType := oobeSetting.GetUserType(); userType != nil {
			data.OutOfBoxExperienceSetting.UserType = state.EnumPtrToTypeString(userType)
		}

		if deviceUsageType := oobeSetting.GetDeviceUsageType(); deviceUsageType != nil {
			data.OutOfBoxExperienceSetting.DeviceUsageType = state.EnumPtrToTypeString(deviceUsageType)
		}
	}

	if essSettings := remoteResource.GetEnrollmentStatusScreenSettings(); essSettings != nil {
		data.EnrollmentStatusScreenSettings = &WindowsEnrollmentStatusScreenSettingsModel{
			HideInstallationProgress:                         state.BoolPtrToTypeBool(essSettings.GetHideInstallationProgress()),
			AllowDeviceUseBeforeProfileAndAppInstallComplete: state.BoolPtrToTypeBool(essSettings.GetAllowDeviceUseBeforeProfileAndAppInstallComplete()),
			BlockDeviceSetupRetryByUser:                      state.BoolPtrToTypeBool(essSettings.GetBlockDeviceSetupRetryByUser()),
			AllowLogCollectionOnInstallFailure:               state.BoolPtrToTypeBool(essSettings.GetAllowLogCollectionOnInstallFailure()),
			CustomErrorMessage:                               state.StringPointerValue(essSettings.GetCustomErrorMessage()),
			InstallProgressTimeoutInMinutes:                  state.Int32PtrToTypeInt32(essSettings.GetInstallProgressTimeoutInMinutes()),
			AllowDeviceUseOnInstallFailure:                   state.BoolPtrToTypeBool(essSettings.GetAllowDeviceUseOnInstallFailure()),
		}
	}

	tflog.Debug(ctx, "Finished mapping Windows Autopilot Deployment Profile from API to Terraform state")
}
