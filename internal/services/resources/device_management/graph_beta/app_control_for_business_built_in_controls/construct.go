package graphBetaAppControlForBusinessBuiltInControls

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates a new App Control for Business configuration policy for the Graph API
func constructResource(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing App Control for Business configuration policy from model")

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	// Set platform (always Windows 10 for App Control for Business)
	platform := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
	requestBody.SetPlatforms(&platform)

	// Set technologies (always mdm for this resource)
	technologies := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
	requestBody.SetTechnologies(&technologies)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	settings, err := constructAppControlSettings(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct App Control settings: %v", err)
	}
	requestBody.SetSettings(settings)

	// Create template reference with hardcoded values
	templateReference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateId := "4321b946-b76b-4450-8afd-769c08b16ffc_1"
	templateFamily := graphmodels.ENDPOINTSECURITYAPPLICATIONCONTROL_DEVICEMANAGEMENTCONFIGURATIONTEMPLATEFAMILY
	templateDisplayName := "App Control for Business"
	templateDisplayVersion := "Version 1"

	templateReference.SetTemplateId(&templateId)
	templateReference.SetTemplateFamily(&templateFamily)
	templateReference.SetTemplateDisplayName(&templateDisplayName)
	templateReference.SetTemplateDisplayVersion(&templateDisplayVersion)
	requestBody.SetTemplateReference(templateReference)

	if err := constructors.DebugLogGraphObject(ctx, "Final JSON to be sent to Graph API", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAppControlSettings constructs the App Control for Business specific settings
func constructAppControlSettings(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel) ([]graphmodels.DeviceManagementConfigurationSettingable, error) {
	tflog.Debug(ctx, "Constructing App Control for Business settings")

	settings := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	// Create the main App Control setting
	setting := graphmodels.NewDeviceManagementConfigurationSetting()

	// Create the choice setting instance
	settingInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingDefinitionId := "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions"
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	// Create the choice setting value
	choiceSettingValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	value := "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_selected"
	choiceSettingValue.SetValue(&value)

	// Create children settings for the built-in controls
	children, err := constructBuiltInControlsChildren(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct built-in controls children: %v", err)
	}
	choiceSettingValue.SetChildren(children)

	// Set setting value template reference
	settingValueTemplateReference := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	settingValueTemplateId := "b28c7dc4-c7b2-4ce2-8f51-6ebfd3ea69d3"
	settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
	choiceSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)

	// Set setting instance template reference
	settingInstanceTemplateReference := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	settingInstanceTemplateId := "1de98212-6949-42dc-a89c-e0ff6e5da04b"
	settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
	settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)

	setting.SetSettingInstance(settingInstance)
	settings = append(settings, setting)

	return settings, nil
}

// constructBuiltInControlsChildren constructs the built-in controls children settings
func constructBuiltInControlsChildren(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel) ([]graphmodels.DeviceManagementConfigurationSettingInstanceable, error) {
	tflog.Debug(ctx, "Constructing built-in controls children")

	children := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)

	// Create the group setting collection instance
	groupSettingCollection := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	settingDefinitionId := "device_vendor_msft_policy_config_applicationcontrol_built_in_controls"
	groupSettingCollection.SetSettingDefinitionId(&settingDefinitionId)

	// Create the group setting collection value
	groupSettingCollectionValue := make([]graphmodels.DeviceManagementConfigurationGroupSettingValueable, 0)

	groupSettingValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
	groupChildren, err := constructGroupChildren(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct group children: %v", err)
	}
	groupSettingValue.SetChildren(groupChildren)

	groupSettingCollectionValue = append(groupSettingCollectionValue, groupSettingValue)
	groupSettingCollection.SetGroupSettingCollectionValue(groupSettingCollectionValue)

	children = append(children, groupSettingCollection)

	return children, nil
}

// constructGroupChildren constructs the group children settings for enable_app_control and trust_apps
func constructGroupChildren(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel) ([]graphmodels.DeviceManagementConfigurationSettingInstanceable, error) {
	tflog.Debug(ctx, "Constructing group children settings")

	children := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)

	// Add enable_app_control setting
	enableAppControlSetting := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	enableAppControlDefId := "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control"
	enableAppControlSetting.SetSettingDefinitionId(&enableAppControlDefId)

	enableAppControlValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	var enableValue string
	switch data.EnableAppControl.ValueString() {
	case "audit":
		enableValue = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_0"
	case "enforce":
		enableValue = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_1"
	default:
		return nil, fmt.Errorf("invalid enable_app_control value: %s, must be 'audit' or 'enforce'", data.EnableAppControl.ValueString())
	}
	enableAppControlValue.SetValue(&enableValue)
	enableAppControlValue.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{})

	enableAppControlSetting.SetChoiceSettingValue(enableAppControlValue)
	children = append(children, enableAppControlSetting)

	// Add trust_apps setting
	trustAppsSetting := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
	trustAppsDefId := "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps"
	trustAppsSetting.SetSettingDefinitionId(&trustAppsDefId)

	trustAppsCollectionValue := make([]graphmodels.DeviceManagementConfigurationChoiceSettingValueable, 0)

	// Convert trust apps from Terraform set to collection values.
	// using custom field name for clarity.
	if !data.AdditionalRulesForTrustingApps.IsNull() && !data.AdditionalRulesForTrustingApps.IsUnknown() {
		var trustApps []string
		diags := data.AdditionalRulesForTrustingApps.ElementsAs(ctx, &trustApps, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract trust apps: %v", diags.Errors())
		}

		for _, trustApp := range trustApps {
			trustAppValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
			var appValue string
			switch trustApp {
			case "trust_apps_with_good_reputation":
				appValue = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_0"
			case "trust_apps_from_managed_installers":
				appValue = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_1"
			default:
				tflog.Warn(ctx, "Unknown trust app value", map[string]interface{}{
					"value": trustApp,
				})
				continue
			}
			trustAppValue.SetValue(&appValue)
			trustAppValue.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{})
			trustAppsCollectionValue = append(trustAppsCollectionValue, trustAppValue)
		}
	}

	trustAppsSetting.SetChoiceSettingCollectionValue(trustAppsCollectionValue)
	children = append(children, trustAppsSetting)

	return children, nil
}
