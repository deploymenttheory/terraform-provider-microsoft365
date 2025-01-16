package graphBetaEndpointPrivilegeManagement

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	sharedConstructor "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors/graph_beta/device_and_app_management"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// PolicyTemplateConfig is a struct that contains the platform, technologies, and template id for a device management policy.
type PolicyTemplateConfig struct {
	Platform     graphmodels.DeviceManagementConfigurationPlatforms
	Technologies graphmodels.DeviceManagementConfigurationTechnologies
	TemplateID   string
}

// Mapping of supported settings catalog template types to their configuration.
var policyConfigMap = map[string]PolicyTemplateConfig{
	"elevation_settings_policy": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.ENDPOINTPRIVILEGEMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "e7dcaba4-959b-46ed-88f0-16ba39b14fd8_1",
	},
	"elevation_rules_policy": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.ENDPOINTPRIVILEGEMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "cff02aad-51b1-498d-83ad-81161a393f56_1",
	},
}

// Main entry point to construct the intune settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *sharedmodels.SettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	// Set general properties
	constructors.SetStringProperty(data.Name, requestBody.SetName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	// Configure template context
	if err := setTemplateContext(ctx, data, requestBody); err != nil {
		return nil, err
	}

	// Set role scope tag IDs
	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Construct settings
	settings := sharedConstructor.ConstructSettingsCatalogSettings(ctx, data.Settings)
	requestBody.SetSettings(settings)

	// Log the constructed resource
	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// setTemplateContext sets the template-specific settings for the device management template.
// It sets the platform, technologies, and template ID reference based on the `SettingsCatalogTemplateType`.
func setTemplateContext(ctx context.Context, data *sharedmodels.SettingsCatalogProfileResourceModel, requestBody graphmodels.DeviceManagementConfigurationPolicyable) error {
	config, exists := policyConfigMap[data.SettingsCatalogTemplateType.ValueString()]
	if !exists {
		tflog.Error(ctx, "Invalid settings catalog template type", map[string]interface{}{
			"settings_catalog_template_type": data.SettingsCatalogTemplateType.ValueString(),
		})
		return fmt.Errorf("invalid settings_catalog_template_type: %s", data.SettingsCatalogTemplateType.ValueString())
	}

	// Platform and Technologies are always required
	requestBody.SetPlatforms(&config.Platform)
	requestBody.SetTechnologies(&config.Technologies)

	// Set TemplateReference if TemplateID exists
	if config.TemplateID != "" {
		templateReference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
		templateReference.SetTemplateId(&config.TemplateID)
		requestBody.SetTemplateReference(templateReference)
	}

	return nil
}
