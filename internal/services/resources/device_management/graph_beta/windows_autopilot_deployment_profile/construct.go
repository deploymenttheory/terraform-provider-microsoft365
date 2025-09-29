package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// stringPtr returns a pointer to the given string value
func stringPtr(s string) *string {
	return &s
}

// constructResource constructs a Windows Autopilot Deployment Profile object for API requests
func constructResource(ctx context.Context, data *WindowsAutopilotDeploymentProfileResourceModel, isCreate bool) (graphmodels.WindowsAutopilotDeploymentProfileable, error) {
	// Use base WindowsAutopilotDeploymentProfile type as shown in SDK example
	resource := graphmodels.NewWindowsAutopilotDeploymentProfile()

	// Set the @odata.type based on device join type
	switch data.DeviceJoinType.ValueString() {
	case "microsoft_entra_hybrid_joined":
		resource.SetOdataType(stringPtr("#microsoft.graph.activeDirectoryWindowsAutopilotDeploymentProfile"))
	case "microsoft_entra_joined":
		resource.SetOdataType(stringPtr("#microsoft.graph.azureADWindowsAutopilotDeploymentProfile"))
	default:
		return nil, fmt.Errorf("invalid device join type: %s", data.DeviceJoinType.ValueString())
	}

	convert.FrameworkToGraphString(data.DisplayName, resource.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, resource.SetDescription)

	// user_select requires empty string for api call - graph xray.
	// performings api call with "" results in a state error.
	// locale: was cty.StringVal("user_select"), but now cty.StringVal("os-default").
	// api appears to be broken. old api doubt it will ever be fixed.
	if data.Locale.ValueString() == "user_select" {
		resource.SetLocale(stringPtr(""))
	} else {
		convert.FrameworkToGraphString(data.Locale, resource.SetLocale)
	}
	convert.FrameworkToGraphString(data.DeviceNameTemplate, resource.SetDeviceNameTemplate)
	convert.FrameworkToGraphString(data.ManagementServiceAppId, resource.SetManagementServiceAppId)
	convert.FrameworkToGraphBool(data.HardwareHashExtractionEnabled, resource.SetHardwareHashExtractionEnabled)
	convert.FrameworkToGraphBool(data.PreprovisioningAllowed, resource.SetPreprovisioningAllowed)

	// hybridAzureADJoinSkipConnectivityCheck is only set during create as per graph xray
	if isCreate && !data.HybridAzureADJoinSkipConnectivityCheck.IsNull() {
		additionalData := map[string]any{
			"hybridAzureADJoinSkipConnectivityCheck": data.HybridAzureADJoinSkipConnectivityCheck.ValueBool(),
		}
		resource.SetAdditionalData(additionalData)
	}

	// deviceType is only set during create as per graph xray
	if isCreate {
		if err := convert.FrameworkToGraphEnum(data.DeviceType, graphmodels.ParseWindowsAutopilotDeviceType, resource.SetDeviceType); err != nil {
			return nil, fmt.Errorf("error setting device type: %v", err)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, resource.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("error setting role scope tag IDs: %v", err)
	}

	if data.OutOfBoxExperienceSetting != nil {
		oobe := graphmodels.NewOutOfBoxExperienceSetting()

		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.PrivacySettingsHidden, oobe.SetPrivacySettingsHidden)
		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.EulaHidden, oobe.SetEulaHidden)
		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.KeyboardSelectionPageSkipped, oobe.SetKeyboardSelectionPageSkipped)
		// Field is always required to be set to TRUE but doesnt configure anything in the gui.
		oobe.SetEscapeLinkHidden(&[]bool{true}[0])

		if err := convert.FrameworkToGraphEnum(data.OutOfBoxExperienceSetting.UserType, graphmodels.ParseWindowsUserType, oobe.SetUserType); err != nil {
			return nil, fmt.Errorf("error setting OOBE setting user type: %v", err)
		}

		if err := convert.FrameworkToGraphEnum(data.OutOfBoxExperienceSetting.DeviceUsageType, graphmodels.ParseWindowsDeviceUsageType, oobe.SetDeviceUsageType); err != nil {
			return nil, fmt.Errorf("error setting OOBE setting device usage type: %v", err)
		}

		resource.SetOutOfBoxExperienceSetting(oobe)
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed Windows Autopilot Deployment Profile Resource", resource); err != nil {
		tflog.Error(ctx, "Failed to log Windows Autopilot Deployment Profile", map[string]any{
			"error": err.Error(),
		})
	}

	return resource, nil
}
