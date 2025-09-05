package graphBetaAppControlForBusinessPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/normalize"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates a new App Control for Business configuration policy with XML content for the Graph API
func constructResource(ctx context.Context, data *AppControlForBusinessPolicyResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing App Control for Business policy with XML content from model")

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

	settings, err := constructAppControlXMLSettings(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct App Control XML settings: %v", err)
	}
	requestBody.SetSettings(settings)

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

// constructAppControlXMLSettings constructs the App Control for Business XML policy settings
func constructAppControlXMLSettings(ctx context.Context, data *AppControlForBusinessPolicyResourceModel) ([]graphmodels.DeviceManagementConfigurationSettingable, error) {
	tflog.Debug(ctx, "Constructing App Control for Business XML settings")

	settings := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	setting := graphmodels.NewDeviceManagementConfigurationSetting()

	// Main choice setting instance for XML configuration
	settingInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingDefinitionId := "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions"
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	choiceSettingValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	value := "device_vendor_msft_policy_config_applicationcontrol_configure_xml_selected"
	choiceSettingValue.SetValue(&value)

	// Create XML content child setting
	xmlContentChild, err := constructXMLContentChild(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct XML content child: %v", err)
	}

	children := []graphmodels.DeviceManagementConfigurationSettingInstanceable{xmlContentChild}
	choiceSettingValue.SetChildren(children)

	// Set template references from GraphXRaySession.go
	settingValueTemplateReference := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	settingValueTemplateId := "b28c7dc4-c7b2-4ce2-8f51-6ebfd3ea69d3"
	settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
	choiceSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)

	settingInstanceTemplateReference := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	settingInstanceTemplateId := "1de98212-6949-42dc-a89c-e0ff6e5da04b"
	settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
	settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)

	setting.SetSettingInstance(settingInstance)
	settings = append(settings, setting)

	return settings, nil
}

// constructXMLContentChild constructs the XML content child setting
func constructXMLContentChild(ctx context.Context, data *AppControlForBusinessPolicyResourceModel) (graphmodels.DeviceManagementConfigurationSettingInstanceable, error) {
	tflog.Debug(ctx, "Constructing XML content child setting")

	xmlContentChild := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionId := "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_xml"
	xmlContentChild.SetSettingDefinitionId(&settingDefinitionId)

	simpleSettingValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
	originalXML := data.PolicyXML.ValueString()

	// Normalize XML content for API submission
	xmlContent := normalize.NormalizeXML(originalXML)

	tflog.Debug(ctx, "Cleaned XML content for Graph API", map[string]interface{}{
		"originalLength": len(originalXML),
		"cleanedLength":  len(xmlContent),
		"hasBOM":         strings.HasPrefix(originalXML, "\ufeff"),
	})

	simpleSettingValue.SetValue(&xmlContent)

	// Set template references from GraphXRaySession.go
	settingValueTemplateReference := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	settingValueTemplateId := "88f6f096-dedb-4cf1-ac2f-4b41e303adb5"
	settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
	simpleSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)

	xmlContentChild.SetSimpleSettingValue(simpleSettingValue)

	settingInstanceTemplateReference := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	settingInstanceTemplateId := "4d709667-63d7-42f2-8e1b-b780f6c3c9c7"
	settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
	xmlContentChild.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)

	return xmlContentChild, nil
}
