package graphBetaLinuxPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates a new Linux script resource model for the Graph API
func constructResource(ctx context.Context, data *LinuxPlatformScriptResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing Linux script resource from model")

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	// Set platform (always Linux for this resource)
	platform := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
	requestBody.SetPlatforms(&platform)

	// Set technologies (always linuxMdm for this resource)
	technologies := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.LINUXMDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
	requestBody.SetTechnologies(&technologies)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	settings, err := constructSettingsCatalogSettings(data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct settings: %v", err)
	}
	requestBody.SetSettings(settings)

	templateReference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
	templateId := "92439f26-2b30-4503-8429-6d40f7e172dd_1" // This is the template ID for Linux Platform Script
	templateReference.SetTemplateId(&templateId)
	requestBody.SetTemplateReference(templateReference)

	if err := constructors.DebugLogGraphObject(ctx, "Final JSON to be sent to Graph API", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructSettingsCatalogSettings creates the settings configuration for the Linux script
// this constructor requires specific template and template instance IDs to be set correctly.
// The values are taken from Microsoft Graph X-Ray.
func constructSettingsCatalogSettings(data *LinuxPlatformScriptResourceModel) ([]graphmodels.DeviceManagementConfigurationSettingable, error) {
	var settings []graphmodels.DeviceManagementConfigurationSettingable

	encodedScript, err := helpers.StringToBase64(data.ScriptContent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to encode script content: %v", err)
	}

	executionContextSetting := graphmodels.NewDeviceManagementConfigurationSetting()
	executionContextInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	executionContextDefId := "linux_customconfig_executioncontext"
	executionContextInstance.SetSettingDefinitionId(&executionContextDefId)

	executionContextValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	executionContextChoice := fmt.Sprintf("linux_customconfig_executioncontext_%s", data.ExecutionContext.ValueString())
	executionContextValue.SetValue(&executionContextChoice)

	children := []graphmodels.DeviceManagementConfigurationSettingInstanceable{}
	executionContextValue.SetChildren(children)

	executionContextValueTemplate := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	executionContextValueTemplateId := "119f0327-4114-444a-b53d-4b55fd579e43"
	executionContextValueTemplate.SetSettingValueTemplateId(&executionContextValueTemplateId)
	executionContextValue.SetSettingValueTemplateReference(executionContextValueTemplate)

	executionContextInstance.SetChoiceSettingValue(executionContextValue)

	executionContextInstanceTemplate := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	executionContextInstanceTemplateId := "2c59a6c5-e874-445b-ac5a-d53688ef838e"
	executionContextInstanceTemplate.SetSettingInstanceTemplateId(&executionContextInstanceTemplateId)
	executionContextInstance.SetSettingInstanceTemplateReference(executionContextInstanceTemplate)

	executionContextSetting.SetSettingInstance(executionContextInstance)
	settings = append(settings, executionContextSetting)

	frequencySetting := graphmodels.NewDeviceManagementConfigurationSetting()
	frequencyInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	frequencyDefId := "linux_customconfig_executionfrequency"
	frequencyInstance.SetSettingDefinitionId(&frequencyDefId)

	frequencyValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	frequencyChoice := fmt.Sprintf("linux_customconfig_executionfrequency_%s", data.ExecutionFrequency.ValueString())
	frequencyValue.SetValue(&frequencyChoice)

	frequencyValueTemplate := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	frequencyValueTemplateId := "d0fb527e-606e-455f-891d-2a4de6a5db90"
	frequencyValueTemplate.SetSettingValueTemplateId(&frequencyValueTemplateId)
	frequencyValue.SetSettingValueTemplateReference(frequencyValueTemplate)

	frequencyInstance.SetChoiceSettingValue(frequencyValue)

	frequencyInstanceTemplate := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	frequencyInstanceTemplateId := "f42b866f-ff2b-4d19-bef8-63e7c763d49b"
	frequencyInstanceTemplate.SetSettingInstanceTemplateId(&frequencyInstanceTemplateId)
	frequencyInstance.SetSettingInstanceTemplateReference(frequencyInstanceTemplate)

	frequencySetting.SetSettingInstance(frequencyInstance)
	settings = append(settings, frequencySetting)

	retriesSetting := graphmodels.NewDeviceManagementConfigurationSetting()
	retriesInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	retriesDefId := "linux_customconfig_executionretries"
	retriesInstance.SetSettingDefinitionId(&retriesDefId)

	retriesValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	retriesChoice := fmt.Sprintf("linux_customconfig_executionretries_%s", data.ExecutionRetries.ValueString())
	retriesValue.SetValue(&retriesChoice)

	retriesValueTemplate := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	retriesValueTemplateId := "92b31053-6ebb-4d2d-9e4d-081fe15d5d21"
	retriesValueTemplate.SetSettingValueTemplateId(&retriesValueTemplateId)
	retriesValue.SetSettingValueTemplateReference(retriesValueTemplate)

	retriesInstance.SetChoiceSettingValue(retriesValue)

	retriesInstanceTemplate := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	retriesInstanceTemplateId := "a3326517-152b-4b32-bc11-8772b5b4fe6a"
	retriesInstanceTemplate.SetSettingInstanceTemplateId(&retriesInstanceTemplateId)
	retriesInstance.SetSettingInstanceTemplateReference(retriesInstanceTemplate)

	retriesSetting.SetSettingInstance(retriesInstance)
	settings = append(settings, retriesSetting)

	scriptSetting := graphmodels.NewDeviceManagementConfigurationSetting()
	scriptInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	scriptDefId := "linux_customconfig_script"
	scriptInstance.SetSettingDefinitionId(&scriptDefId)

	scriptValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
	scriptValue.SetValue(&encodedScript)

	scriptValueTemplate := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	scriptValueTemplateId := "18dc8a98-2ecd-4753-8baf-3ab7a1d677a9"
	scriptValueTemplate.SetSettingValueTemplateId(&scriptValueTemplateId)
	scriptValue.SetSettingValueTemplateReference(scriptValueTemplate)

	scriptInstance.SetSimpleSettingValue(scriptValue)

	scriptInstanceTemplate := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	scriptInstanceTemplateId := "add4347a-f9aa-4202-a497-34a4c178d013"
	scriptInstanceTemplate.SetSettingInstanceTemplateId(&scriptInstanceTemplateId)
	scriptInstance.SetSettingInstanceTemplateReference(scriptInstanceTemplate)

	scriptSetting.SetSettingInstance(scriptInstance)
	settings = append(settings, scriptSetting)

	return settings, nil
}
