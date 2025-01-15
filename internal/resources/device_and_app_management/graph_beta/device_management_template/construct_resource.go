package graphBetaDeviceManagementTemplate

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

var policyConfigMap = map[string]PolicyTemplateConfig{
	"windows_anti_virus_defender_update_controls": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "e3f74c5a-a6de-411d-aef6-eb15628f3a0a_1",
	},
	"windows_anti_virus_microsoft_defender_antivirus_exclusions": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "45fea5e9-280d-4da1-9792-fb5736da0ca9_1",
	},
}

// Main entry point to construct the intune settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *sharedmodels.SettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	constructors.SetStringProperty(data.Name, requestBody.SetName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	if err := setTemplateContext(ctx, data, requestBody); err != nil {
		return nil, err
	}

	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	settings := sharedConstructor.ConstructSettingsCatalogSettings(ctx, data.Settings)
	requestBody.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// setTemplateContext sets the template specific settings for the device management template.
// It sets the platform, technologies, and template id reference.
func setTemplateContext(ctx context.Context, data *sharedmodels.SettingsCatalogProfileResourceModel, requestBody graphmodels.DeviceManagementConfigurationPolicyable) error {
	config, exists := policyConfigMap[data.ConfigurationPolicyTemplateType.ValueString()]
	if !exists {
		tflog.Error(ctx, "Invalid configuration policy template type", map[string]interface{}{
			"configuration_policy_template_type": data.ConfigurationPolicyTemplateType.ValueString(),
		})
		return fmt.Errorf("invalid configuration_policy_template_type: %s", data.ConfigurationPolicyTemplateType.ValueString())
	}

	requestBody.SetPlatforms(&config.Platform)
	requestBody.SetTechnologies(&config.Technologies)

	templateReference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateReference.SetTemplateId(&config.TemplateID)
	requestBody.SetTemplateReference(templateReference)

	return nil
}
