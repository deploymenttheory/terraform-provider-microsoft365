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

// Main entry point to construct the intune settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *sharedmodels.SettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	constructors.SetStringProperty(data.Name, requestBody.SetName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	// Set platform (always Windows for this resource). may change in the future
	platform := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
	requestBody.SetPlatforms(&platform)

	technologies := graphmodels.DeviceManagementConfigurationTechnologies(
		graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES |
			graphmodels.ENDPOINTPRIVILEGEMANAGEMENT_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES,
	)
	requestBody.SetTechnologies(&technologies)

	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	settings := sharedConstructor.ConstructSettingsCatalogSettings(ctx, data.Settings)
	requestBody.SetSettings(settings)

	// Set the templateReference based on ConfigurationPolicyTemplateType
	templateReference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
	switch data.ConfigurationPolicyTemplateType.ValueString() {
	case "elevation_settings_policy":
		templateId := "e7dcaba4-959b-46ed-88f0-16ba39b14fd8_1" // Template ID for Elevation Settings Policy
		templateReference.SetTemplateId(&templateId)
	case "elevation_rules_policy":
		templateId := "cff02aad-51b1-498d-83ad-81161a393f56_1" // Template ID for Elevation Rules Policy
		templateReference.SetTemplateId(&templateId)
	default:
		return nil, fmt.Errorf("invalid configuration_policy_template_type: %s", data.ConfigurationPolicyTemplateType.ValueString())
	}

	requestBody.SetTemplateReference(templateReference)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
