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
	Platform       graphmodels.DeviceManagementConfigurationPlatforms
	Technologies   graphmodels.DeviceManagementConfigurationTechnologies
	TemplateID     string
	CreationSource string
}

var policyConfigMap = map[string]PolicyTemplateConfig{
	"linux_anti_virus_microsoft_defender_antivirus": {
		Platform:     graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "4cfd164c-5e8a-4ea9-b15d-9aa71e4ffff4_1",
	},
	"linux_anti_virus_microsoft_defender_antivirus_exclusions": {
		Platform:     graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "8a17a1e5-3df4-4e07-9d20-3878267a79b8_1",
	},
	"linux_endpoint_detection_and_response": {
		Platform:     graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "3514388a-d4d1-4aa8-bd64-c317776008f5_1",
	},
	"macOS_anti_virus_microsoft_defender_antivirus": {
		Platform:     graphmodels.MACOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "2d345ec2-c817-49e5-9156-3ed416dc972a_1",
	},
	"macOS_anti_virus_microsoft_defender_antivirus_exclusions": {
		Platform:     graphmodels.MACOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "43397174-2244-4006-b5ad-421b369e90d4_1",
	},
	"windows_anti_virus_defender_update_controls": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "e3f74c5a-a6de-411d-aef6-eb15628f3a0a_1",
	},
	"windows_anti_virus_microsoft_defender_antivirus": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "804339ad-1553-4478-a742-138fb5807418_1",
	},
	"windows_anti_virus_microsoft_defender_antivirus_exclusions": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "45fea5e9-280d-4da1-9792-fb5736da0ca9_1",
	},
	"windows_anti_virus_security_experience": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "d948ff9b-99cb-4ee0-8012-1fbc09685377_1",
	},
	"windows_(config_mgr)_anti_virus_microsoft_defender_antivirus": {
		Platform:       graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies:   graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		CreationSource: "SccmAV",
	},
	"windows_(config_mgr)_anti_virus_windows_security_experience": {
		Platform:       graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies:   graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		CreationSource: "WindowsSecurity",
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

	// Only set CreationSource / TemplateReference if it exists
	if config.CreationSource != "" {
		requestBody.SetCreationSource(&config.CreationSource)
	}

	if config.TemplateID != "" {
		templateReference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
		templateReference.SetTemplateId(&config.TemplateID)
		requestBody.SetTemplateReference(templateReference)
	}

	return nil
}
