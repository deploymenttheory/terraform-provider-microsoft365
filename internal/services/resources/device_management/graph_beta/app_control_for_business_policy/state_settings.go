package graphBetaAppControlForBusinessPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapAppControlPolicyXMLSettingsToTerraform maps App Control for Business XML policy settings from Graph API response to Terraform state
func MapAppControlPolicyXMLSettingsToTerraform(ctx context.Context, data *AppControlForBusinessPolicyResourceModel, settingsResponse graphmodels.DeviceManagementConfigurationSettingCollectionResponseable) error {
	if settingsResponse == nil {
		tflog.Debug(ctx, "No settings response data to process")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "Settings array is empty")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Processing %d settings from API response", len(settings)))

	for _, setting := range settings {
		if setting == nil {
			continue
		}

		settingInstance := setting.GetSettingInstance()
		if settingInstance == nil {
			continue
		}

		settingDefId := settingInstance.GetSettingDefinitionId()
		if settingDefId == nil {
			continue
		}

		// Look for the main App Control policy setting
		if *settingDefId == "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions" {
			err := extractAppControlPolicyXMLSettings(ctx, data, settingInstance)
			if err != nil {
				return fmt.Errorf("failed to extract app control policy XML settings: %v", err)
			}
		}
	}

	return nil
}

// extractAppControlPolicyXMLSettings extracts XML policy settings from the main setting instance
func extractAppControlPolicyXMLSettings(ctx context.Context, data *AppControlForBusinessPolicyResourceModel, settingInstance graphmodels.DeviceManagementConfigurationSettingInstanceable) error {
	choiceInstance, ok := settingInstance.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return fmt.Errorf("setting instance is not a choice setting")
	}

	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return nil
	}

	// Verify this is an XML configuration type
	value := choiceValue.GetValue()
	if value == nil || *value != "device_vendor_msft_policy_config_applicationcontrol_configure_xml_selected" {
		tflog.Debug(ctx, "Setting is not XML configuration type", map[string]any{
			"value": value,
		})
		return nil
	}

	children := choiceValue.GetChildren()
	if children == nil {
		return nil
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		childDefId := child.GetSettingDefinitionId()
		if childDefId == nil {
			continue
		}

		// Look for the XML content setting
		if *childDefId == "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_xml" {
			err := extractPolicyXMLContent(ctx, data, child)
			if err != nil {
				return fmt.Errorf("failed to extract policy XML content: %v", err)
			}
		}
	}

	return nil
}

// extractPolicyXMLContent extracts the XML policy content from the setting
func extractPolicyXMLContent(ctx context.Context, data *AppControlForBusinessPolicyResourceModel, child graphmodels.DeviceManagementConfigurationSettingInstanceable) error {
	simpleInstance, ok := child.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable)
	if !ok {
		return fmt.Errorf("child setting is not a simple setting")
	}

	simpleValue := simpleInstance.GetSimpleSettingValue()
	if simpleValue == nil {
		return nil
	}

	stringValue, ok := simpleValue.(graphmodels.DeviceManagementConfigurationStringSettingValueable)
	if !ok {
		return fmt.Errorf("simple setting value is not a string value")
	}

	xmlContent := stringValue.GetValue()
	if xmlContent != nil {
		// Normalize XML content to match what we send to the API
		normalizedXML := normalize.ReverseNormalizeXML(*xmlContent)

		data.PolicyXML = types.StringValue(normalizedXML)
		tflog.Debug(ctx, "Successfully extracted and normalized XML policy content", map[string]any{
			"xmlLength": len(normalizedXML),
			"hasBOM":    strings.HasPrefix(*xmlContent, "\ufeff"),
		})
	} else {
		data.PolicyXML = types.StringNull()
		tflog.Debug(ctx, "XML policy content is null")
	}

	return nil
}
