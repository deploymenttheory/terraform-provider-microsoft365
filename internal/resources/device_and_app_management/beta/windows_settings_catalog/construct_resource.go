package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructSettingsCatalogProfileRequestBody constructs a DeviceManagementConfigurationPolicyable object from a WindowsSettingsCatalogProfileResourceModel object
// and a boolean indicating whether the profile is being created or updated.
func constructResource(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel, isCreate bool) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	var catalog WindowsSettingsCatalogSettingsResourceModel

	tflog.Debug(ctx, "Constructing WindowsSettingsCatalogProfile resource")
	construct.DebugPrintStruct(ctx, "Constructed WindowsSettingsCatalogProfile Resource from model", data)

	profile := graphmodels.NewDeviceManagementConfigurationPolicy()

	construct.SetStringProperty(data.DisplayName, profile.SetName)
	construct.SetStringProperty(data.Description, profile.SetDescription)

	// Conditionally set settings if isCreate is true else skip for update operation
	if isCreate {

		platforms := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
		profile.SetPlatforms(&platforms)

		technologies := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
		profile.SetTechnologies(&technologies)

		settings, err := constructSettingsCatalogSettings(ctx, catalog.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to construct settings: %v", err)
		}
		profile.SetSettings(settings)
	}

	// Handle RoleScopeTagIds
	construct.SetArrayProperty(data.RoleScopeTagIds, profile.SetRoleScopeTagIds)

	tflog.Debug(ctx, "Finished constructing WindowsSettingsCatalogProfile resource")
	return profile, nil
}

func constructSettingsCatalogSettings(ctx context.Context, settingConfigs WindowsSettingsCatalogSettingsResourceModel) ([]graphmodels.DeviceManagementConfigurationSettingable, error) {
	var settings []graphmodels.DeviceManagementConfigurationSettingable

	for _, settingConfig := range settingConfigs.Settings {
		setting := graphmodels.NewDeviceManagementConfigurationSetting()
		if settingConfig.SettingInstance != nil {
			settingInstance, err := constructSettingInstance(ctx, *settingConfig.SettingInstance)
			if err != nil {
				return nil, fmt.Errorf("failed to construct setting instance: %v", err)
			}
			setting.SetSettingInstance(settingInstance)
			settings = append(settings, setting)
			tflog.Debug(ctx, fmt.Sprintf("Adding setting: %s", settingConfig.SettingInstance.SettingDefinitionID.ValueString()))
		}
	}

	return settings, nil
}

func constructSettingInstance(ctx context.Context, instanceConfig DeviceManagementConfigurationSettingInstance) (graphmodels.DeviceManagementConfigurationSettingInstanceable, error) {
	if instanceConfig.ODataType.IsNull() || instanceConfig.ODataType.IsUnknown() {
		return nil, fmt.Errorf("ODataType is null or unknown")
	}

	switch instanceConfig.ODataType.ValueString() {
	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
		instance.SetSettingDefinitionId(instanceConfig.SettingDefinitionID.ValueStringPointer())

		var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

		if instanceConfig.ChoiceSettingValue == nil || instanceConfig.ChoiceSettingValue.Value.IsNull() || instanceConfig.ChoiceSettingValue.Value.IsUnknown() {
			return nil, fmt.Errorf("ChoiceSettingValue or its Value is null or unknown")
		}

		valueStr := instanceConfig.ChoiceSettingValue.Value.ValueString()

		// Attempt to parse the value as different types
		if intValue, err := strconv.Atoi(valueStr); err == nil {
			int32Value := int32(intValue)
			intSettingValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
			intSettingValue.SetValue(&int32Value)
			simpleValue = intSettingValue
		} else if int64Value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			int32Value := int32(int64Value)
			intSettingValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
			intSettingValue.SetValue(&int32Value)
			simpleValue = intSettingValue
		} else if boolValue, err := strconv.ParseBool(valueStr); err == nil {
			boolStr := strconv.FormatBool(boolValue)
			stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			stringValue.SetValue(&boolStr)
			simpleValue = stringValue
		} else if floatValue, err := strconv.ParseFloat(valueStr, 64); err == nil {
			floatStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
			stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			stringValue.SetValue(&floatStr)
			simpleValue = stringValue
		} else {
			// For any other types, use as string
			stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			stringValue.SetValue(&valueStr)
			simpleValue = stringValue
		}

		instance.SetSimpleSettingValue(simpleValue)
		return instance, nil

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
		instance.SetSettingDefinitionId(instanceConfig.SettingDefinitionID.ValueStringPointer())
		choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()

		if instanceConfig.ChoiceSettingValue == nil || instanceConfig.ChoiceSettingValue.Value.IsNull() || instanceConfig.ChoiceSettingValue.Value.IsUnknown() {
			return nil, fmt.Errorf("ChoiceSettingValue or its Value is null or unknown")
		}

		valueStr := instanceConfig.ChoiceSettingValue.Value.ValueString()
		choiceValue.SetValue(&valueStr)

		instance.SetChoiceSettingValue(choiceValue)
		return instance, nil

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
		instance.SetSettingDefinitionId(instanceConfig.SettingDefinitionID.ValueStringPointer())

		if instanceConfig.ChoiceSettingValue == nil {
			return nil, fmt.Errorf("ChoiceSettingValue is null")
		}

		var collectionValues []graphmodels.DeviceManagementConfigurationSimpleSettingValueable

		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			if child.ChoiceSettingValue == nil || child.ChoiceSettingValue.Value.IsNull() || child.ChoiceSettingValue.Value.IsUnknown() {
				continue // Skip null or unknown values
			}

			valueStr := child.ChoiceSettingValue.Value.ValueString()
			var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

			// Attempt to parse the value as different types
			if intValue, err := strconv.Atoi(valueStr); err == nil {
				int32Value := int32(intValue)
				intSettingValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
				intSettingValue.SetValue(&int32Value)
				simpleValue = intSettingValue
			} else if int64Value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
				int32Value := int32(int64Value)
				intSettingValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
				intSettingValue.SetValue(&int32Value)
				simpleValue = intSettingValue
			} else if boolValue, err := strconv.ParseBool(valueStr); err == nil {
				boolStr := strconv.FormatBool(boolValue)
				stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
				stringValue.SetValue(&boolStr)
				simpleValue = stringValue
			} else if floatValue, err := strconv.ParseFloat(valueStr, 64); err == nil {
				floatStr := strconv.FormatFloat(floatValue, 'f', -1, 64)
				stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
				stringValue.SetValue(&floatStr)
				simpleValue = stringValue
			} else {
				// For any other types, use as string
				stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
				stringValue.SetValue(&valueStr)
				simpleValue = stringValue
			}

			collectionValues = append(collectionValues, simpleValue)
		}

		instance.SetSimpleSettingCollectionValue(collectionValues)
		return instance, nil

	default:
		tflog.Warn(ctx, fmt.Sprintf("Unsupported setting type '%s' for '%s'. Skipping this setting.", instanceConfig.ODataType.ValueString(), instanceConfig.SettingDefinitionID.ValueString()))
		return nil, nil
	}
}
