package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel) (graphmodels.DeviceManagementCompliancePolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementCompliancePolicy()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	// Always set platforms to linux for Linux device compliance policies
	platform := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
	requestBody.SetPlatforms(&platform)

	// Always set technologies to linuxMdm for Linux device compliance policies
	technology := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.LINUXMDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
	requestBody.SetTechnologies(&technology)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	settings, err := constructLinuxComplianceSettings(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct Linux compliance settings: %s", err)
	}
	requestBody.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructLinuxComplianceSettings constructs the complex device management configuration settings
func constructLinuxComplianceSettings(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel) ([]graphmodels.DeviceManagementConfigurationSettingable, error) {
	settings := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	// Distribution Allowed Distros
	if !data.DistributionAllowedDistros.IsNull() && !data.DistributionAllowedDistros.IsUnknown() {
		setting, err := constructDistributionAllowedDistrosSetting(ctx, data.DistributionAllowedDistros)
		if err != nil {
			return nil, fmt.Errorf("failed to construct distribution allowed distros setting: %s", err)
		}
		settings = append(settings, setting)
	}

	// Custom Compliance
	if !data.CustomComplianceRequired.IsNull() && !data.CustomComplianceRequired.IsUnknown() {
		setting, err := constructCustomComplianceSetting(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("failed to construct custom compliance setting: %s", err)
		}
		settings = append(settings, setting)
	}

	// Device Encryption
	if !data.DeviceEncryptionRequired.IsNull() && !data.DeviceEncryptionRequired.IsUnknown() {
		setting, err := constructDeviceEncryptionSetting(ctx, data.DeviceEncryptionRequired)
		if err != nil {
			return nil, fmt.Errorf("failed to construct device encryption setting: %s", err)
		}
		settings = append(settings, setting)
	}

	// Password Policy Settings
	passwordSettings, err := constructPasswordPolicySettings(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct password policy settings: %s", err)
	}
	settings = append(settings, passwordSettings...)

	return settings, nil
}

// constructDistributionAllowedDistrosSetting constructs the distribution allowed distros setting
func constructDistributionAllowedDistrosSetting(ctx context.Context, allowedDistros types.List) (graphmodels.DeviceManagementConfigurationSettingable, error) {
	setting := graphmodels.NewDeviceManagementConfigurationSetting()

	// Set the @odata.type
	odataType := "#microsoft.graph.deviceManagementConfigurationSetting"
	setting.SetOdataType(&odataType)

	// Create the group setting collection instance
	settingInstance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	instanceOdataType := "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
	settingInstance.SetOdataType(&instanceOdataType)

	// Set the setting definition ID
	settingDefinitionId := "linux_distribution_alloweddistros"
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	// Parse the allowed distros from the terraform model
	var allowedDistrosModels []AllowedDistributionModel
	diags := allowedDistros.ElementsAs(ctx, &allowedDistrosModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse allowed distros list: %v", diags.Errors())
	}

	// Create the group setting collection value
	groupSettingCollectionValue := make([]graphmodels.DeviceManagementConfigurationGroupSettingValueable, 0)

	for _, distro := range allowedDistrosModels {
		groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()

		// Create children for this distribution
		children := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)

		// Maximum version setting
		if !distro.MaximumVersion.IsNull() && !distro.MaximumVersion.IsUnknown() {
			maxVersionSetting := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
			maxVersionOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
			maxVersionSetting.SetOdataType(&maxVersionOdataType)
			maxVersionDefinitionId := "linux_distribution_alloweddistros_item_maximumversion"
			maxVersionSetting.SetSettingDefinitionId(&maxVersionDefinitionId)

			// Create string setting value
			maxVersionValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			maxVersionValueOdataType := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
			maxVersionValue.SetOdataType(&maxVersionValueOdataType)
			maxVersionStr := distro.MaximumVersion.ValueString()
			maxVersionValue.SetValue(&maxVersionStr)
			maxVersionSetting.SetSimpleSettingValue(maxVersionValue)

			children = append(children, maxVersionSetting)
		}

		// Minimum version setting
		if !distro.MinimumVersion.IsNull() && !distro.MinimumVersion.IsUnknown() {
			minVersionSetting := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
			minVersionOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
			minVersionSetting.SetOdataType(&minVersionOdataType)
			minVersionDefinitionId := "linux_distribution_alloweddistros_item_minimumversion"
			minVersionSetting.SetSettingDefinitionId(&minVersionDefinitionId)

			// Create string setting value
			minVersionValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			minVersionValueOdataType := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
			minVersionValue.SetOdataType(&minVersionValueOdataType)
			minVersionStr := distro.MinimumVersion.ValueString()
			minVersionValue.SetValue(&minVersionStr)
			minVersionSetting.SetSimpleSettingValue(minVersionValue)

			children = append(children, minVersionSetting)
		}

		// Distribution type setting
		if !distro.Type.IsNull() && !distro.Type.IsUnknown() {
			typeSetting := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
			typeOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
			typeSetting.SetOdataType(&typeOdataType)
			typeDefinitionId := "linux_distribution_alloweddistros_item_$type"
			typeSetting.SetSettingDefinitionId(&typeDefinitionId)

			// Create choice setting value
			typeValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
			typeValueOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
			typeValue.SetOdataType(&typeValueOdataType)
			typeStr := fmt.Sprintf("linux_distribution_alloweddistros_item_$type_%s", distro.Type.ValueString())
			typeValue.SetValue(&typeStr)
			emptyChildren := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)
			typeValue.SetChildren(emptyChildren)
			typeSetting.SetChoiceSettingValue(typeValue)

			children = append(children, typeSetting)
		}

		groupValue.SetChildren(children)
		groupSettingCollectionValue = append(groupSettingCollectionValue, groupValue)
	}

	settingInstance.SetGroupSettingCollectionValue(groupSettingCollectionValue)
	setting.SetSettingInstance(settingInstance)

	return setting, nil
}

// constructCustomComplianceSetting constructs the custom compliance setting
func constructCustomComplianceSetting(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel) (graphmodels.DeviceManagementConfigurationSettingable, error) {
	setting := graphmodels.NewDeviceManagementConfigurationSetting()

	// Set the @odata.type
	odataType := "#microsoft.graph.deviceManagementConfigurationSetting"
	setting.SetOdataType(&odataType)

	// Create the choice setting instance
	settingInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	instanceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&instanceOdataType)

	// Set the setting definition ID
	settingDefinitionId := "linux_customcompliance_required"
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	// Create the choice setting value
	choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceValueOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&choiceValueOdataType)

	// Set the value based on the boolean
	var valueStr string
	if data.CustomComplianceRequired.ValueBool() {
		valueStr = "linux_customcompliance_required_true"

		// Add children for discovery script and rules when enabled
		children := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)

		// Discovery script
		if !data.CustomComplianceDiscoveryScript.IsNull() && !data.CustomComplianceDiscoveryScript.IsUnknown() {
			scriptSetting := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
			scriptOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
			scriptSetting.SetOdataType(&scriptOdataType)
			scriptDefinitionId := "linux_customcompliance_discoveryscript"
			scriptSetting.SetSettingDefinitionId(&scriptDefinitionId)

			// Create reference setting value
			scriptValue := graphmodels.NewDeviceManagementConfigurationReferenceSettingValue()
			scriptValueOdataType := "#microsoft.graph.deviceManagementConfigurationReferenceSettingValue"
			scriptValue.SetOdataType(&scriptValueOdataType)
			scriptStr := data.CustomComplianceDiscoveryScript.ValueString()
			scriptValue.SetValue(&scriptStr)
			scriptSetting.SetSimpleSettingValue(scriptValue)

			children = append(children, scriptSetting)
		}

		// Custom compliance rules
		if !data.CustomComplianceRules.IsNull() && !data.CustomComplianceRules.IsUnknown() {
			rulesSetting := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
			rulesOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
			rulesSetting.SetOdataType(&rulesOdataType)
			rulesDefinitionId := "linux_customcompliance_rules"
			rulesSetting.SetSettingDefinitionId(&rulesDefinitionId)

			// Create string setting value
			rulesValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			rulesValueOdataType := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
			rulesValue.SetOdataType(&rulesValueOdataType)
			rulesStr := data.CustomComplianceRules.ValueString()
			rulesValue.SetValue(&rulesStr)
			rulesSetting.SetSimpleSettingValue(rulesValue)

			children = append(children, rulesSetting)
		}

		choiceValue.SetChildren(children)
	} else {
		valueStr = "linux_customcompliance_required_false"
		emptyChildren := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)
		choiceValue.SetChildren(emptyChildren)
	}

	choiceValue.SetValue(&valueStr)
	settingInstance.SetChoiceSettingValue(choiceValue)
	setting.SetSettingInstance(settingInstance)

	return setting, nil
}

// constructDeviceEncryptionSetting constructs the device encryption setting
func constructDeviceEncryptionSetting(ctx context.Context, encryptionRequired types.Bool) (graphmodels.DeviceManagementConfigurationSettingable, error) {
	setting := graphmodels.NewDeviceManagementConfigurationSetting()

	// Set the @odata.type
	odataType := "#microsoft.graph.deviceManagementConfigurationSetting"
	setting.SetOdataType(&odataType)

	// Create the choice setting instance
	settingInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	instanceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&instanceOdataType)

	// Set the setting definition ID
	settingDefinitionId := "linux_deviceencryption_required"
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	// Create the choice setting value
	choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceValueOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&choiceValueOdataType)

	// Set the value based on the boolean
	var valueStr string
	if encryptionRequired.ValueBool() {
		valueStr = "linux_deviceencryption_required_true"
	} else {
		valueStr = "linux_deviceencryption_required_false"
	}

	choiceValue.SetValue(&valueStr)
	emptyChildren := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)
	choiceValue.SetChildren(emptyChildren)
	settingInstance.SetChoiceSettingValue(choiceValue)
	setting.SetSettingInstance(settingInstance)

	return setting, nil
}

// constructPasswordPolicySettings constructs all password policy settings
func constructPasswordPolicySettings(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel) ([]graphmodels.DeviceManagementConfigurationSettingable, error) {
	settings := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	// Password policy settings mapping
	passwordSettings := map[string]types.Int32{
		"linux_passwordpolicy_minimumdigits":    data.PasswordPolicyMinimumDigits,
		"linux_passwordpolicy_minimumlength":    data.PasswordPolicyMinimumLength,
		"linux_passwordpolicy_minimumlowercase": data.PasswordPolicyMinimumLowercase,
		"linux_passwordpolicy_minimumsymbols":   data.PasswordPolicyMinimumSymbols,
		"linux_passwordpolicy_minimumuppercase": data.PasswordPolicyMinimumUppercase,
	}

	for settingDefinitionId, value := range passwordSettings {
		if !value.IsNull() && !value.IsUnknown() {
			setting, err := constructPasswordPolicySetting(settingDefinitionId, value)
			if err != nil {
				return nil, fmt.Errorf("failed to construct password policy setting %s: %s", settingDefinitionId, err)
			}
			settings = append(settings, setting)
		}
	}

	return settings, nil
}

// constructPasswordPolicySetting constructs a single password policy setting
func constructPasswordPolicySetting(settingDefinitionId string, value types.Int32) (graphmodels.DeviceManagementConfigurationSettingable, error) {
	setting := graphmodels.NewDeviceManagementConfigurationSetting()

	// Set the @odata.type
	odataType := "#microsoft.graph.deviceManagementConfigurationSetting"
	setting.SetOdataType(&odataType)

	// Create the simple setting instance
	settingInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	instanceOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&instanceOdataType)

	// Set the setting definition ID
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	// Create the integer setting value
	intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
	intValueOdataType := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	intValue.SetOdataType(&intValueOdataType)

	intVal := int32(value.ValueInt32())
	intValue.SetValue(&intVal)

	settingInstance.SetSimpleSettingValue(intValue)
	setting.SetSettingInstance(settingInstance)

	return setting, nil
}
