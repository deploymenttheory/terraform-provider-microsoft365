package graphBetaLinuxPlatformScript

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune linux platform script resource for the Terraform provider.
func constructResource(ctx context.Context, data *LinuxPlatformScriptResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	Name := data.DisplayName.ValueString()
	requestBody.SetName(&Name)

	description := data.Description.ValueString()
	requestBody.SetDescription(&description)

	// Set platforms to linux (always)
	parsedPlatform, err := graphmodels.ParseDeviceManagementConfigurationPlatforms("linux")
	if err != nil {
		return nil, fmt.Errorf("error parsing platforms: %v", err)
	}
	if platform, ok := parsedPlatform.(*graphmodels.DeviceManagementConfigurationPlatforms); ok {
		requestBody.SetPlatforms(platform)
	}

	// Set technologies to linuxMdm (always)
	parsedTechnologies, err := graphmodels.ParseDeviceManagementConfigurationTechnologies("linuxMdm")
	if err != nil {
		return nil, fmt.Errorf("error parsing technologies: %v", err)
	}
	if technologies, ok := parsedTechnologies.(*graphmodels.DeviceManagementConfigurationTechnologies); ok {
		requestBody.SetTechnologies(technologies)
	}

	if len(data.RoleScopeTagIds) > 0 {
		var tagIds []string
		for _, tag := range data.RoleScopeTagIds {
			tagIds = append(tagIds, tag.ValueString())
		}
		requestBody.SetRoleScopeTagIds(tagIds)
	} else {
		requestBody.SetRoleScopeTagIds([]string{"0"})
	}

	//TODO
	//settings := constructSettingsCatalogSettings(ctx, data.Settings)
	//requestBody.SetSettings(settings)

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructSettingsCatalogSettings is a helper function to construct the linux platform script settings catalog settings from the JSON data.
func constructSettingsCatalogSettings(ctx context.Context, settingsJSON types.String) []graphmodels.DeviceManagementConfigurationSettingable {
	tflog.Debug(ctx, "Constructing settings catalog settings")

	var simplifiedSettings []struct {
		SettingDefinitionID string `json:"settingDefinitionId"`
		Value               string `json:"value"`
		TemplateID          string `json:"templateId"`
		InstanceTemplateID  string `json:"instanceTemplateId"`
		ODataType           string `json:"@odata.type"`
	}

	if err := json.Unmarshal([]byte(settingsJSON.ValueString()), &simplifiedSettings); err != nil {
		tflog.Error(ctx, "Failed to unmarshal settings JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	tflog.Debug(ctx, "Unmarshaled simplified settings", map[string]interface{}{
		"data": simplifiedSettings,
	})

	var settingsCollection []graphmodels.DeviceManagementConfigurationSettingable

	for _, detail := range simplifiedSettings {
		baseSetting := graphmodels.NewDeviceManagementConfigurationSetting()

		switch detail.ODataType {
		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
			// Construct choice setting instance
			settingInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
			settingInstance.SetSettingDefinitionId(&detail.SettingDefinitionID)

			choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
			choiceValue.SetValue(&detail.Value)

			// Attach value template reference if available
			if detail.TemplateID != "" {
				valueTemplateRef := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
				valueTemplateRef.SetSettingValueTemplateId(&detail.TemplateID)
				choiceValue.SetSettingValueTemplateReference(valueTemplateRef)
			}

			settingInstance.SetChoiceSettingValue(choiceValue)

			// Attach instance template reference if available
			if detail.InstanceTemplateID != "" {
				instanceTemplateRef := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
				instanceTemplateRef.SetSettingInstanceTemplateId(&detail.InstanceTemplateID)
				settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			baseSetting.SetSettingInstance(settingInstance)

		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
			// Construct simple setting instance
			settingInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
			settingInstance.SetSettingDefinitionId(&detail.SettingDefinitionID)

			simpleValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			simpleValue.SetValue(&detail.Value)

			// Attach value template reference if available
			if detail.TemplateID != "" {
				valueTemplateRef := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
				valueTemplateRef.SetSettingValueTemplateId(&detail.TemplateID)
				simpleValue.SetSettingValueTemplateReference(valueTemplateRef)
			}

			settingInstance.SetSimpleSettingValue(simpleValue)

			// Attach instance template reference if available
			if detail.InstanceTemplateID != "" {
				instanceTemplateRef := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
				instanceTemplateRef.SetSettingInstanceTemplateId(&detail.InstanceTemplateID)
				settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			baseSetting.SetSettingInstance(settingInstance)
		}

		settingsCollection = append(settingsCollection, baseSetting)
	}

	tflog.Debug(ctx, "Constructed simplified settings collection", map[string]interface{}{
		"count": len(settingsCollection),
	})

	return settingsCollection
}
