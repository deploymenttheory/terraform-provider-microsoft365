package graphBetaReuseablePolicySettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *sharedmodels.SettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementReusablePolicySettingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementReusablePolicySetting()

	// Set display name
	displayName := data.Name.ValueString()
	requestBody.SetDisplayName(&displayName)

	// Set description
	description := data.Description.ValueString()
	requestBody.SetDescription(&description)

	// Set setting definition ID
	settingDefinitionId := "device_vendor_msft_policy_privilegemanagement_reusablesettings_certificatefile"
	requestBody.SetSettingDefinitionId(&settingDefinitionId)

	// Create and set setting instance
	settingInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	// Create and configure the setting value
	settingValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
	value := data.Settings.ValueString()
	settingValue.SetValue(&value)
	settingInstance.SetSimpleSettingValue(settingValue)

	requestBody.SetSettingInstance(settingInstance)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
