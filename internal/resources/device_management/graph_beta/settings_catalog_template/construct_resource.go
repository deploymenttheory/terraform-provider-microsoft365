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
	"macOS_endpoint_detection_and_response": {
		Platform:     graphmodels.MACOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "a6ff37f6-c841-4264-9249-1ecf793d94ef_1",
	},
	"security_baseline_for_windows_10_and_later_version_24H2": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "66df8dce-0166-4b82-92f7-1f74e3ca17a3_4",
	},
	"security_baseline_for_microsoft_defender_for_endpoint_version_24H1": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "49b8320f-e179-472e-8e2c-2fde00289ca2_1",
	},
	"security_baseline_for_microsoft_edge_version_128": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "c66347b7-8325-4954-a235-3bf2233dfbfd_3",
	},
	"security_baseline_for_windows_365": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "b00e1a0f-19dd-41de-8243-e6733ca7b4ae_1",
	},
	"security_baseline_for_microsoft_365_apps": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "90316f12-246d-44c6-a767-f87692e86083_2",
	},
	"windows_account_protection": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "fcef01f2-439d-4c3f-9184-823fd6e97646_1",
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
	"windows_app_control_for_business": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "4321b946-b76b-4450-8afd-769c08b16ffc_1",
	},
	"windows_attack_surface_reduction_app_and_browser_isolation": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "9f667e40-8f3c-4f88-80d8-457f16906315_1",
	},
	"windows_attack_surface_reduction_attack_surface_reduction_rules": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "e8c053d6-9f95-42b1-a7f1-ebfd71c67a4b_1",
	},
	"windows_attack_surface_reduction_app_device_control": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "0f2034c6-3cd6-4ee1-bd37-f3c0693e9548_1",
	},
	"windows_attack_surface_reduction_exploit_protection": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "d02f2162-fcac-48db-9b7b-b0a3f160d2c2_1",
	},
	"windows_disk_encryption_bitlocker": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "46ddfc50-d10f-4867-b852-9434254b3bff_1",
	},
	"windows_disk_encryption_personal_data": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "0b5708d9-9bc2-49a9-b4f7-ec463fcc41e0_1",
	},
	"windows_endpoint_detection_and_response": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "0385b795-0f2f-44ac-8602-9f65bf6adede_1",
	},
	"windows_firewall": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "6078910e-d808-4a9f-a51d-1b8a7bacb7c0_1",
	},
	"windows_firewall_rules": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES | graphmodels.MICROSOFTSENSE_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "19c8aa67-f286-4861-9aa0-f23541d31680_1",
	},
	"windows_hyper-v_firewall_rules": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "a5481c22-7a2a-4f59-a33e-6eee30d02f94_1",
	},
	"windows_local_admin_password_solution_(windows_LAPS)": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "adc46e5a-f4aa-4ff6-aeff-4f27bc525796_1",
	},
	"windows_local_user_group_membership": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "22968f54-45fa-486c-848e-f8224aa69772_1",
	},
	"windows_(config_mgr)_attack_surface_reduction_app_and_browser_isolation": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "e373ebb7-c1c5-4ffb-9ce0-698f1834fd9d_1",
	},
	"windows_(config_mgr)_attack_surface_reduction_attack_surface_reduction_rules": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "5dd36540-eb22-4e7e-b19c-2a07772ba627_1",
	},
	"windows_(config_mgr)_attack_surface_reduction_exploit_protection": {
		Platform:       graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies:   graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		CreationSource: "ASR_ExploitProtection",
	},
	"windows_(config_mgr)_attack_surface_reduction_web_protection": {
		Platform:       graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies:   graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		CreationSource: "ASR_WebProtection",
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
	"windows_(config_mgr)_endpoint_detection_and_response": {
		Platform:       graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies:   graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		CreationSource: "SccmEDR",
	},
	"windows_(config_mgr)_firewall": {
		Platform:       graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies:   graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		CreationSource: "Firewall",
	},
	"windows_(config_mgr)_firewall_profile": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "c2791bb6-ad62-412d-99dc-cb179ef72dee_1",
	},
	"windows_(config_mgr)_firewall_rules": {
		Platform:     graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS,
		Technologies: graphmodels.CONFIGMANAGER_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
		TemplateID:   "48da42ed-5df7-485e-8b9d-4844ed5a92bd_1",
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

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
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
