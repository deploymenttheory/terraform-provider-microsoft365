package configurationPolicyTemplateBuilders

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// detectBooleanChoice determines if a value represents a boolean choice (_0/_1 suffix)
// Returns: (baseDefinitionId, boolValue, isBooleanChoice)
func detectBooleanChoice(value string) (string, bool, bool) {
	if strings.HasSuffix(value, "_0") {
		return strings.TrimSuffix(value, "_0"), false, true
	}
	if strings.HasSuffix(value, "_1") {
		return strings.TrimSuffix(value, "_1"), true, true
	}
	return value, false, false
}

// constructSettingInstanceTemplateReference creates a setting instance template reference if templateId is provided
// This creates the settingInstanceTemplateReference object that contains only settingInstanceTemplateId
func constructSettingInstanceTemplateReference(instanceTemplateId string) models.DeviceManagementConfigurationSettingInstanceTemplateReferenceable {
	if instanceTemplateId == "" {
		return nil
	}
	templateRef := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()

	// Convert string to types.String and use constructor helper
	templateIdValue := types.StringValue(instanceTemplateId)
	constructors.SetStringProperty(templateIdValue, templateRef.SetSettingInstanceTemplateId)

	return templateRef
}

// constructSettingValueTemplateReference creates a setting value template reference if templateId is provided
// This creates the settingValueTemplateReference object that contains both settingValueTemplateId and useTemplateDefault
// useTemplateDefault defaults to false if not specified
func constructSettingValueTemplateReference(valueTemplateId string, useTemplateDefault ...bool) models.DeviceManagementConfigurationSettingValueTemplateReferenceable {
	if valueTemplateId == "" {
		return nil
	}
	templateRef := models.NewDeviceManagementConfigurationSettingValueTemplateReference()

	// Convert string to types.String and use constructor helper
	templateIdValue := types.StringValue(valueTemplateId)
	constructors.SetStringProperty(templateIdValue, templateRef.SetSettingValueTemplateId)

	// Handle useTemplateDefault - default to false if not provided
	useDefault := false
	if len(useTemplateDefault) > 0 {
		useDefault = useTemplateDefault[0]
	}
	templateRef.SetUseTemplateDefault(&useDefault)

	return templateRef
}

// createOptionalInstanceTemplateReference creates a setting instance template reference if templateId is provided
func createOptionalInstanceTemplateReference(instanceTemplateId string) models.DeviceManagementConfigurationSettingInstanceTemplateReferenceable {
	if instanceTemplateId == "" {
		return nil
	}
	templateRef := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
	templateRef.SetSettingInstanceTemplateId(&instanceTemplateId)
	return templateRef
}

// constructSettingValueTemplateId creates a setting value template reference if templateId is provided
func constructSettingValueTemplateId(valueTemplateId string) models.DeviceManagementConfigurationSettingValueTemplateReferenceable {
	if valueTemplateId == "" {
		return nil
	}
	templateRef := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
	templateRef.SetSettingValueTemplateId(&valueTemplateId)
	return templateRef
}

// ConstructSimpleChoiceSetting creates a simple choice setting
func ConstructSimpleChoiceSetting(
	ctx context.Context,
	settingDefinitionId string,
	value string,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("value cannot be empty")
	}

	tflog.Debug(ctx, "Constructing simple choice setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"value":               value,
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceValue.SetValue(&value)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&odataTypeValue)

	// Set empty children array
	var children []models.DeviceManagementConfigurationSettingInstanceable
	choiceValue.SetChildren(children)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		choiceValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	settingInstance.SetChoiceSettingValue(choiceValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed simple choice setting")
	return setting, nil
}

// ConstructSimpleChoiceSettingWithTemplate creates a choice setting
func ConstructSimpleChoiceSettingWithTemplate(
	ctx context.Context,
	settingDefinitionId string,
	choiceValue string,
	instanceTemplateId string,
	valueTemplateId string,
	useTemplateDefault bool,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if choiceValue == "" {
		return nil, fmt.Errorf("choiceValue cannot be empty")
	}

	tflog.Debug(ctx, "Constructing simple choice setting with template", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"choiceValue":         choiceValue,
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
		"useTemplateDefault":  useTemplateDefault,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceSettingValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceSettingValue.SetValue(&choiceValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceSettingValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference with useTemplateDefault
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId, useTemplateDefault); valueTemplateRef != nil {
		choiceSettingValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	// Simple choice settings have no children
	var children []models.DeviceManagementConfigurationSettingInstanceable
	choiceSettingValue.SetChildren(children)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed simple choice setting with template")
	return setting, nil
}

// ConstructSimpleStringSetting creates a simple string setting
// Example: custom_message = "Contact IT for assistance"
func ConstructSimpleStringSetting(
	ctx context.Context,
	settingDefinitionId string,
	value string,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}

	tflog.Debug(ctx, "Constructing simple string setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"value":               value,
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create simple setting instance
	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create string value
	stringValue := models.NewDeviceManagementConfigurationStringSettingValue()
	stringValue.SetValue(&value)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	stringValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		stringValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	settingInstance.SetSimpleSettingValue(stringValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed simple string setting")
	return setting, nil
}

// ConstructSimpleIntegerSetting creates a simple integer setting
// Example: timeout_minutes = "30"
func ConstructSimpleIntegerSetting(
	ctx context.Context,
	settingDefinitionId string,
	value string,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("value cannot be empty for integer setting")
	}

	tflog.Debug(ctx, "Constructing simple integer setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"value":               value,
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
	})

	// Parse string to integer
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		tflog.Error(ctx, "Failed to parse integer value", map[string]interface{}{
			"value": value,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to parse integer value '%s': %w", value, err)
	}

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create simple setting instance
	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create integer value
	integerValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
	int32Value := int32(intVal)
	integerValue.SetValue(&int32Value)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	integerValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		integerValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	settingInstance.SetSimpleSettingValue(integerValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed simple integer setting")
	return setting, nil
}

// ConstructSimpleSecretSetting creates a simple secret setting
// Example: password = "secret123", valueState = "notEncrypted"
func ConstructSimpleSecretSetting(
	ctx context.Context,
	settingDefinitionId string,
	value string,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	// Use default value state of "notEncrypted"
	return ConstructSimpleSecretSettingWithState(ctx, settingDefinitionId, value, "notEncrypted", instanceTemplateId, valueTemplateId)
}

// ConstructSimpleSecretSettingWithState creates a simple secret setting with explicit value state
// Example: password = "secret123", valueState = "notEncrypted"
func ConstructSimpleSecretSettingWithState(
	ctx context.Context,
	settingDefinitionId string,
	value string,
	valueState string,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("value cannot be empty for secret setting")
	}
	if valueState == "" {
		valueState = "notEncrypted" // Default value state
	}

	tflog.Debug(ctx, "Constructing simple secret setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"valueState":          valueState,
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create simple setting instance
	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create secret value
	secretValue := models.NewDeviceManagementConfigurationSecretSettingValue()
	secretValue.SetValue(&value)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
	secretValue.SetOdataType(&odataTypeValue)

	// Parse and set value state
	parsedValueState, err := models.ParseDeviceManagementConfigurationSecretSettingValueState(valueState)
	if err != nil {
		tflog.Error(ctx, "Failed to parse secret value state", map[string]interface{}{
			"valueState": valueState,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to parse secret value state '%s': %w", valueState, err)
	}

	if typedValueState, ok := parsedValueState.(*models.DeviceManagementConfigurationSecretSettingValueState); ok && typedValueState != nil {
		secretValue.SetValueState(typedValueState)
	}

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		secretValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	settingInstance.SetSimpleSettingValue(secretValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed simple secret setting")
	return setting, nil
}

// ConstructSimpleBooleanChoiceSetting creates a boolean choice setting (detects _0/_1 suffix automatically)
// Example: allow_skip = "enrollment_autopilot_dpp_allowskip_1" (true) or "enrollment_autopilot_dpp_allowskip_0" (false)
func ConstructSimpleBooleanChoiceSetting(
	ctx context.Context,
	settingDefinitionId string,
	value string,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("value cannot be empty")
	}

	// Detect if this is a boolean choice pattern
	_, boolValue, isBooleanChoice := detectBooleanChoice(value)
	if !isBooleanChoice {
		return nil, fmt.Errorf("value '%s' does not follow boolean choice pattern (must end with _0 or _1)", value)
	}

	tflog.Debug(ctx, "Constructing simple boolean choice setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"value":               value,
		"detectedBoolValue":   boolValue,
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
	})

	// Use the ConstructSimpleChoiceSetting since boolean settings are represented as choice settings in the API
	setting, err := ConstructSimpleChoiceSetting(ctx, settingDefinitionId, value, instanceTemplateId, valueTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to construct boolean choice setting: %w", err)
	}

	tflog.Debug(ctx, "Successfully constructed simple boolean choice setting", map[string]interface{}{
		"detectedBoolValue": boolValue,
	})
	return setting, nil
}

// Helper function to create a string child setting instance
func createStringChildSettingInstance(
	ctx context.Context,
	childSettingDefinitionId string,
	childValue string,
	childInstanceTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingInstanceable, error) {
	childInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	childInstance.SetSettingDefinitionId(&childSettingDefinitionId)

	childOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	childInstance.SetOdataType(&childOdataType)

	// Set optional child instance template reference
	if childInstanceTemplateRef := createOptionalInstanceTemplateReference(childInstanceTemplateId); childInstanceTemplateRef != nil {
		childInstance.SetSettingInstanceTemplateReference(childInstanceTemplateRef)
	}

	// Create string value for child
	stringValue := models.NewDeviceManagementConfigurationStringSettingValue()
	stringValue.SetValue(&childValue)

	stringOdataType := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	stringValue.SetOdataType(&stringOdataType)

	// Set optional child value template reference
	if childValueTemplateRef := constructSettingValueTemplateId(childValueTemplateId); childValueTemplateRef != nil {
		stringValue.SetSettingValueTemplateReference(childValueTemplateRef)
	}

	childInstance.SetSimpleSettingValue(stringValue)
	return childInstance, nil
}

// Helper function to create an integer child setting instance
func createIntegerChildSettingInstance(
	ctx context.Context,
	childSettingDefinitionId string,
	childValue string,
	childInstanceTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingInstanceable, error) {
	// Parse string to integer
	intVal, err := strconv.ParseInt(childValue, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse child integer value '%s': %w", childValue, err)
	}

	childInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	childInstance.SetSettingDefinitionId(&childSettingDefinitionId)

	childOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	childInstance.SetOdataType(&childOdataType)

	// Set optional child instance template reference
	if childInstanceTemplateRef := createOptionalInstanceTemplateReference(childInstanceTemplateId); childInstanceTemplateRef != nil {
		childInstance.SetSettingInstanceTemplateReference(childInstanceTemplateRef)
	}

	// Create integer value for child
	integerValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
	int32Value := int32(intVal)
	integerValue.SetValue(&int32Value)

	intOdataType := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	integerValue.SetOdataType(&intOdataType)

	// Set optional child value template reference
	if childValueTemplateRef := constructSettingValueTemplateId(childValueTemplateId); childValueTemplateRef != nil {
		integerValue.SetSettingValueTemplateReference(childValueTemplateRef)
	}

	childInstance.SetSimpleSettingValue(integerValue)
	return childInstance, nil
}

// Helper function to create a secret child setting instance
func createSecretChildSettingInstance(
	ctx context.Context,
	childSettingDefinitionId string,
	childValue string,
	childValueState string,
	childInstanceTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingInstanceable, error) {
	if childValueState == "" {
		childValueState = "notEncrypted" // Default value state
	}

	childInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	childInstance.SetSettingDefinitionId(&childSettingDefinitionId)

	childOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	childInstance.SetOdataType(&childOdataType)

	// Set optional child instance template reference
	if childInstanceTemplateRef := createOptionalInstanceTemplateReference(childInstanceTemplateId); childInstanceTemplateRef != nil {
		childInstance.SetSettingInstanceTemplateReference(childInstanceTemplateRef)
	}

	// Create secret value for child
	secretValue := models.NewDeviceManagementConfigurationSecretSettingValue()
	secretValue.SetValue(&childValue)

	secretOdataType := "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
	secretValue.SetOdataType(&secretOdataType)

	// Parse and set value state
	parsedValueState, err := models.ParseDeviceManagementConfigurationSecretSettingValueState(childValueState)
	if err != nil {
		return nil, fmt.Errorf("failed to parse child secret value state '%s': %w", childValueState, err)
	}

	if typedValueState, ok := parsedValueState.(*models.DeviceManagementConfigurationSecretSettingValueState); ok && typedValueState != nil {
		secretValue.SetValueState(typedValueState)
	}

	// Set optional child value template reference
	if childValueTemplateRef := constructSettingValueTemplateId(childValueTemplateId); childValueTemplateRef != nil {
		secretValue.SetSettingValueTemplateReference(childValueTemplateRef)
	}

	childInstance.SetSimpleSettingValue(secretValue)
	return childInstance, nil
}

// ConstructChoiceWithStringSetting creates a choice setting with one string child (FSLogix cache directory pattern)
// Example: choice_value = "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~ccd_ccdcachedirectory_1", child_value = "\\\\server\\share\\cache"
func ConstructChoiceWithStringSetting(
	ctx context.Context,
	settingDefinitionId string,
	choiceValue string,
	childSettingDefinitionId string,
	childValue string,
	instanceTemplateId string,
	valueTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if choiceValue == "" {
		return nil, fmt.Errorf("choiceValue cannot be empty")
	}
	if childSettingDefinitionId == "" {
		return nil, fmt.Errorf("childSettingDefinitionId cannot be empty")
	}

	tflog.Debug(ctx, "Constructing choice with string child setting", map[string]interface{}{
		"settingDefinitionId":      settingDefinitionId,
		"choiceValue":              choiceValue,
		"childSettingDefinitionId": childSettingDefinitionId,
		"childValue":               childValue,
		"instanceTemplateId":       instanceTemplateId,
		"valueTemplateId":          valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceSettingValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceSettingValue.SetValue(&choiceValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceSettingValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		choiceSettingValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	// Create string child setting instance
	childInstance, err := createStringChildSettingInstance(ctx, childSettingDefinitionId, childValue, "", childValueTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to create string child setting: %w", err)
	}

	// Set children array with single child
	children := []models.DeviceManagementConfigurationSettingInstanceable{childInstance}
	choiceSettingValue.SetChildren(children)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed choice with string child setting")
	return setting, nil
}

// ConstructChoiceWithIntegerSetting creates a choice setting with one integer child (FSLogix timeout pattern)
// Example: choice_value = "device_vendor_msft_policy_config_fslogixv1~policy~fslogix~profiles~profiles_ccd_profilesccdunregistertimeout_1", child_value = "30"
func ConstructChoiceWithIntegerSetting(
	ctx context.Context,
	settingDefinitionId string,
	choiceValue string,
	childSettingDefinitionId string,
	childValue string,
	instanceTemplateId string,
	valueTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if choiceValue == "" {
		return nil, fmt.Errorf("choiceValue cannot be empty")
	}
	if childSettingDefinitionId == "" {
		return nil, fmt.Errorf("childSettingDefinitionId cannot be empty")
	}
	if childValue == "" {
		return nil, fmt.Errorf("childValue cannot be empty for integer setting")
	}

	tflog.Debug(ctx, "Constructing choice with integer child setting", map[string]interface{}{
		"settingDefinitionId":      settingDefinitionId,
		"choiceValue":              choiceValue,
		"childSettingDefinitionId": childSettingDefinitionId,
		"childValue":               childValue,
		"instanceTemplateId":       instanceTemplateId,
		"valueTemplateId":          valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceSettingValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceSettingValue.SetValue(&choiceValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceSettingValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		choiceSettingValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	// Create integer child setting instance
	childInstance, err := createIntegerChildSettingInstance(ctx, childSettingDefinitionId, childValue, "", childValueTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to create integer child setting: %w", err)
	}

	// Set children array with single child
	children := []models.DeviceManagementConfigurationSettingInstanceable{childInstance}
	choiceSettingValue.SetChildren(children)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed choice with integer child setting")
	return setting, nil
}

// ConstructChoiceWithBooleanSetting creates a choice setting with one boolean child
// Example: choice_value = "some_setting_1", child_value = "some_child_setting_1" (auto-detects _1 as true)
func ConstructChoiceWithBooleanSetting(
	ctx context.Context,
	settingDefinitionId string,
	choiceValue string,
	childSettingDefinitionId string,
	childValue string,
	instanceTemplateId string,
	valueTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if choiceValue == "" {
		return nil, fmt.Errorf("choiceValue cannot be empty")
	}
	if childSettingDefinitionId == "" {
		return nil, fmt.Errorf("childSettingDefinitionId cannot be empty")
	}
	if childValue == "" {
		return nil, fmt.Errorf("childValue cannot be empty")
	}

	// Detect if child value is a boolean choice pattern
	_, boolValue, isBooleanChoice := detectBooleanChoice(childValue)
	if !isBooleanChoice {
		return nil, fmt.Errorf("childValue '%s' does not follow boolean choice pattern (must end with _0 or _1)", childValue)
	}

	tflog.Debug(ctx, "Constructing choice with boolean child setting", map[string]interface{}{
		"settingDefinitionId":      settingDefinitionId,
		"choiceValue":              choiceValue,
		"childSettingDefinitionId": childSettingDefinitionId,
		"childValue":               childValue,
		"detectedBoolValue":        boolValue,
		"instanceTemplateId":       instanceTemplateId,
		"valueTemplateId":          valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceSettingValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceSettingValue.SetValue(&choiceValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceSettingValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		choiceSettingValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	// Create boolean child setting instance as a choice setting (Microsoft's pattern)
	childInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	childInstance.SetSettingDefinitionId(&childSettingDefinitionId)

	childOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	childInstance.SetOdataType(&childOdataType)

	// Create choice value for the boolean child
	childChoiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	childChoiceValue.SetValue(&childValue)

	childChoiceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	childChoiceValue.SetOdataType(&childChoiceOdataType)

	// Boolean choices have no children of their own
	var childChildren []models.DeviceManagementConfigurationSettingInstanceable
	childChoiceValue.SetChildren(childChildren)

	childInstance.SetChoiceSettingValue(childChoiceValue)

	// Set children array with single boolean child
	children := []models.DeviceManagementConfigurationSettingInstanceable{childInstance}
	choiceSettingValue.SetChildren(children)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed choice with boolean child setting", map[string]interface{}{
		"detectedBoolValue": boolValue,
	})
	return setting, nil
}

// ConstructChoiceWithSecretSetting creates a choice setting with one secret child
// Example: choice_value = "some_setting_1", child_value = "secret123", childValueState = "notEncrypted"
func ConstructChoiceWithSecretSetting(
	ctx context.Context,
	settingDefinitionId string,
	choiceValue string,
	childSettingDefinitionId string,
	childValue string,
	childValueState string,
	instanceTemplateId string,
	valueTemplateId string,
	childValueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if choiceValue == "" {
		return nil, fmt.Errorf("choiceValue cannot be empty")
	}
	if childSettingDefinitionId == "" {
		return nil, fmt.Errorf("childSettingDefinitionId cannot be empty")
	}
	if childValue == "" {
		return nil, fmt.Errorf("childValue cannot be empty for secret setting")
	}

	tflog.Debug(ctx, "Constructing choice with secret child setting", map[string]interface{}{
		"settingDefinitionId":      settingDefinitionId,
		"choiceValue":              choiceValue,
		"childSettingDefinitionId": childSettingDefinitionId,
		"childValueState":          childValueState,
		"instanceTemplateId":       instanceTemplateId,
		"valueTemplateId":          valueTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceSettingValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceSettingValue.SetValue(&choiceValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceSettingValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId); valueTemplateRef != nil {
		choiceSettingValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	// Create secret child setting instance
	childInstance, err := createSecretChildSettingInstance(ctx, childSettingDefinitionId, childValue, childValueState, "", childValueTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret child setting: %w", err)
	}

	// Set children array with single child
	children := []models.DeviceManagementConfigurationSettingInstanceable{childInstance}
	choiceSettingValue.SetChildren(children)

	settingInstance.SetChoiceSettingValue(choiceSettingValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed choice with secret child setting")
	return setting, nil
}

// ========================================================================================
// COLLECTION BUILDERS
// ========================================================================================

// ConstructStringCollectionSetting creates a collection of string values
// Example: allowed_domains = ["domain1.com", "domain2.com"]
func ConstructStringCollectionSetting(
	ctx context.Context,
	settingDefinitionId string,
	values []string,
	instanceTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("values cannot be empty for string collection")
	}

	tflog.Debug(ctx, "Constructing string collection setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"valuesCount":         len(values),
		"instanceTemplateId":  instanceTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create simple setting collection instance
	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create collection values
	var simpleSettingCollectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
	for _, val := range values {
		stringValue := models.NewDeviceManagementConfigurationStringSettingValue()
		stringValue.SetValue(&val)

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
		stringValue.SetOdataType(&odataTypeValue)

		simpleSettingCollectionValues = append(simpleSettingCollectionValues, stringValue)
	}

	settingInstance.SetSimpleSettingCollectionValue(simpleSettingCollectionValues)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed string collection setting")
	return setting, nil
}

// ConstructIntegerCollectionSetting creates a collection of integer values
// Example: retry_intervals = ["30", "60", "120"]
func ConstructIntegerCollectionSetting(
	ctx context.Context,
	settingDefinitionId string,
	values []string,
	instanceTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("values cannot be empty for integer collection")
	}

	tflog.Debug(ctx, "Constructing integer collection setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"valuesCount":         len(values),
		"instanceTemplateId":  instanceTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create simple setting collection instance
	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create collection values with integer parsing
	var simpleSettingCollectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
	for i, val := range values {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			tflog.Error(ctx, "Failed to parse integer value in collection", map[string]interface{}{
				"index": i,
				"value": val,
				"error": err.Error(),
			})
			return nil, fmt.Errorf("failed to parse integer value '%s' at index %d: %w", val, i, err)
		}

		integerValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
		int32Value := int32(intVal)
		integerValue.SetValue(&int32Value)

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
		integerValue.SetOdataType(&odataTypeValue)

		simpleSettingCollectionValues = append(simpleSettingCollectionValues, integerValue)
	}

	settingInstance.SetSimpleSettingCollectionValue(simpleSettingCollectionValues)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed integer collection setting")
	return setting, nil
}

// ConstructFormattedJSONCollectionSetting creates a collection of pre-formatted JSON strings
// Example: apps = ["{\"id\":\"app1\",\"type\":\"#microsoft.graph.win32LobApp\"}", ...]
func ConstructFormattedJSONCollectionSetting(
	ctx context.Context,
	settingDefinitionId string,
	jsonValues []string,
	instanceTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if len(jsonValues) == 0 {
		return nil, fmt.Errorf("jsonValues cannot be empty")
	}

	tflog.Debug(ctx, "Constructing formatted JSON collection setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"valuesCount":         len(jsonValues),
		"instanceTemplateId":  instanceTemplateId,
	})

	// Use the string collection builder since JSON strings are stored as strings
	setting, err := ConstructStringCollectionSetting(ctx, settingDefinitionId, jsonValues, instanceTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to construct JSON collection setting: %w", err)
	}

	tflog.Debug(ctx, "Successfully constructed formatted JSON collection setting")
	return setting, nil
}

// ConstructChoiceCollectionSetting creates a choice setting collection
// Example: printing_options = ["printingsettings_2", "printingsettings_1", "printingsettings_0"]
func ConstructChoiceCollectionSetting(
	ctx context.Context,
	settingDefinitionId string,
	choiceValues []string,
	instanceTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if len(choiceValues) == 0 {
		return nil, fmt.Errorf("choiceValues cannot be empty")
	}

	tflog.Debug(ctx, "Constructing choice collection setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"valuesCount":         len(choiceValues),
		"instanceTemplateId":  instanceTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting collection instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice collection values
	var choiceSettingCollectionValues []models.DeviceManagementConfigurationChoiceSettingValueable
	for _, val := range choiceValues {
		choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
		choiceValue.SetValue(&val)

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
		choiceValue.SetOdataType(&odataTypeValue)

		// Choice collection values typically have no children
		var children []models.DeviceManagementConfigurationSettingInstanceable
		choiceValue.SetChildren(children)

		choiceSettingCollectionValues = append(choiceSettingCollectionValues, choiceValue)
	}

	settingInstance.SetChoiceSettingCollectionValue(choiceSettingCollectionValues)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed choice collection setting")
	return setting, nil
}

// ========================================================================================
// ADVANCED CHOICE BUILDERS (for Windows Attack Surface Reduction patterns)
// ========================================================================================

// ConstructChoiceWithMultipleChoiceChildren creates a choice setting with multiple choice children (Windows ASR pattern)
// This handles the complex Windows Attack Surface Reduction pattern where a choice setting has many choice children
func ConstructChoiceWithMultipleChoiceChildren(
	ctx context.Context,
	settingDefinitionId string,
	choiceValue string,
	children []ChoiceChildConfig,
	instanceTemplateId string,
	valueTemplateId string,
	useTemplateDefault bool,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if choiceValue == "" {
		return nil, fmt.Errorf("choiceValue cannot be empty")
	}
	if len(children) == 0 {
		return nil, fmt.Errorf("children cannot be empty")
	}

	tflog.Debug(ctx, "Constructing choice with multiple choice children", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"choiceValue":         choiceValue,
		"childrenCount":       len(children),
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
		"useTemplateDefault":  useTemplateDefault,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create choice setting instance
	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create choice value
	choiceSettingValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceSettingValue.SetValue(&choiceValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceSettingValue.SetOdataType(&odataTypeValue)

	// Set optional value template reference with useTemplateDefault
	if valueTemplateRef := constructSettingValueTemplateReference(valueTemplateId, useTemplateDefault); valueTemplateRef != nil {
		choiceSettingValue.SetSettingValueTemplateReference(valueTemplateRef)
	}

	// Create children
	childInstances, err := createChoiceChildren(ctx, children)
	if err != nil {
		return nil, fmt.Errorf("failed to create choice children: %w", err)
	}

	choiceSettingValue.SetChildren(childInstances)
	settingInstance.SetChoiceSettingValue(choiceSettingValue)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed choice with multiple choice children")
	return setting, nil
}

// Helper function to create children for choice settings
func createChoiceChildren(ctx context.Context, childConfigs []ChoiceChildConfig) ([]models.DeviceManagementConfigurationSettingInstanceable, error) {
	var children []models.DeviceManagementConfigurationSettingInstanceable

	for _, childConfig := range childConfigs {
		var childInstance models.DeviceManagementConfigurationSettingInstanceable
		var err error

		switch childConfig.SettingType {
		case "choice":
			childInstance = models.NewDeviceManagementConfigurationChoiceSettingInstance()
			childInstance.SetSettingDefinitionId(&childConfig.SettingDefinitionId)

			choiceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
			childInstance.SetOdataType(&choiceOdataType)

			// Set optional instance template reference
			if instanceTemplateRef := createOptionalInstanceTemplateReference(childConfig.InstanceTemplateId); instanceTemplateRef != nil {
				childInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			// Create choice value
			choiceInstance := childInstance.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
			choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
			choiceValue.SetValue(&childConfig.Value)

			choiceValueOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
			choiceValue.SetOdataType(&choiceValueOdataType)

			// Set optional value template reference with useTemplateDefault
			if valueTemplateRef := constructSettingValueTemplateReference(childConfig.ValueTemplateId, childConfig.UseTemplateDefault); valueTemplateRef != nil {
				choiceValue.SetSettingValueTemplateReference(valueTemplateRef)
			}

			// Choice settings in this pattern typically have no children
			var choiceChildren []models.DeviceManagementConfigurationSettingInstanceable
			choiceValue.SetChildren(choiceChildren)

			choiceInstance.SetChoiceSettingValue(choiceValue)

		case "simple_string":
			childInstance, err = createStringChildSettingInstance(
				ctx,
				childConfig.SettingDefinitionId,
				childConfig.Value,
				childConfig.InstanceTemplateId,
				childConfig.ValueTemplateId,
			)

		case "simple_integer":
			childInstance, err = createIntegerChildSettingInstance(
				ctx,
				childConfig.SettingDefinitionId,
				childConfig.Value,
				childConfig.InstanceTemplateId,
				childConfig.ValueTemplateId,
			)

		case "simple_collection":
			childInstance = models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
			childInstance.SetSettingDefinitionId(&childConfig.SettingDefinitionId)

			collectionOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
			childInstance.SetOdataType(&collectionOdataType)

			// Set optional instance template reference
			if instanceTemplateRef := createOptionalInstanceTemplateReference(childConfig.InstanceTemplateId); instanceTemplateRef != nil {
				childInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			// Create collection values
			var collectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
			for _, val := range childConfig.Values {
				stringValue := models.NewDeviceManagementConfigurationStringSettingValue()
				stringValue.SetValue(&val)

				stringOdataType := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
				stringValue.SetOdataType(&stringOdataType)

				collectionValues = append(collectionValues, stringValue)
			}

			collectionInstance := childInstance.(models.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
			collectionInstance.SetSimpleSettingCollectionValue(collectionValues)

		case "choice_collection":
			childInstance = models.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
			childInstance.SetSettingDefinitionId(&childConfig.SettingDefinitionId)

			collectionOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
			childInstance.SetOdataType(&collectionOdataType)

			// Set optional instance template reference
			if instanceTemplateRef := createOptionalInstanceTemplateReference(childConfig.InstanceTemplateId); instanceTemplateRef != nil {
				childInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			// Create choice collection values
			var choiceCollectionValues []models.DeviceManagementConfigurationChoiceSettingValueable
			for _, val := range childConfig.ChoiceValues {
				choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
				choiceValue.SetValue(&val)

				choiceValueOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
				choiceValue.SetOdataType(&choiceValueOdataType)

				// Choice collection values have no children
				var choiceChildren []models.DeviceManagementConfigurationSettingInstanceable
				choiceValue.SetChildren(choiceChildren)

				choiceCollectionValues = append(choiceCollectionValues, choiceValue)
			}

			groupCollectionInstance := childInstance.(models.DeviceManagementConfigurationChoiceSettingCollectionInstanceable)
			groupCollectionInstance.SetChoiceSettingCollectionValue(choiceCollectionValues)

		default:
			return nil, fmt.Errorf("unsupported choice child setting type: %s", childConfig.SettingType)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create choice child setting %s: %w", childConfig.SettingDefinitionId, err)
		}

		children = append(children, childInstance)
	}

	return children, nil
}

// ========================================================================================
// GROUP COLLECTION BUILDERS
// ========================================================================================

// ConstructGroupCollectionSetting creates a group setting collection with mixed children
// This handles the complex macOS DDM patterns with nested groups and mixed setting types
func ConstructGroupCollectionSetting(
	ctx context.Context,
	settingDefinitionId string,
	groupConfigs [][]GroupChildConfig, // Array of group instances, each with their children
	instanceTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if len(groupConfigs) == 0 {
		return nil, fmt.Errorf("groupConfigs cannot be empty")
	}

	tflog.Debug(ctx, "Constructing group collection setting", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"groupInstanceCount":  len(groupConfigs),
		"instanceTemplateId":  instanceTemplateId,
	})

	setting := models.NewDeviceManagementConfigurationSetting()

	// Create group setting collection instance
	settingInstance := models.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	// Set optional instance template reference
	if instanceTemplateRef := constructSettingInstanceTemplateReference(instanceTemplateId); instanceTemplateRef != nil {
		settingInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
	}

	// Create group values
	var groupValues []models.DeviceManagementConfigurationGroupSettingValueable
	for i, groupChildren := range groupConfigs {
		groupValue := models.NewDeviceManagementConfigurationGroupSettingValue()

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationGroupSettingValue"
		groupValue.SetOdataType(&odataTypeValue)

		// Create children for this group instance
		children, err := createGroupChildren(ctx, groupChildren)
		if err != nil {
			return nil, fmt.Errorf("failed to create children for group instance %d: %w", i, err)
		}

		groupValue.SetChildren(children)
		groupValues = append(groupValues, groupValue)
	}

	settingInstance.SetGroupSettingCollectionValue(groupValues)
	setting.SetSettingInstance(settingInstance)

	tflog.Debug(ctx, "Successfully constructed group collection setting")
	return setting, nil
}

// Helper function to create children for group settings
func createGroupChildren(ctx context.Context, childConfigs []GroupChildConfig) ([]models.DeviceManagementConfigurationSettingInstanceable, error) {
	var children []models.DeviceManagementConfigurationSettingInstanceable

	for _, childConfig := range childConfigs {
		var childInstance models.DeviceManagementConfigurationSettingInstanceable
		var err error

		switch childConfig.SettingType {
		case "simple_string":
			childInstance, err = createStringChildSettingInstance(
				ctx,
				childConfig.SettingDefinitionId,
				childConfig.Value,
				childConfig.InstanceTemplateId,
				childConfig.ValueTemplateId,
			)

		case "simple_integer":
			childInstance, err = createIntegerChildSettingInstance(
				ctx,
				childConfig.SettingDefinitionId,
				childConfig.Value,
				childConfig.InstanceTemplateId,
				childConfig.ValueTemplateId,
			)

		case "choice":
			childInstance = models.NewDeviceManagementConfigurationChoiceSettingInstance()
			childInstance.SetSettingDefinitionId(&childConfig.SettingDefinitionId)

			choiceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
			childInstance.SetOdataType(&choiceOdataType)

			// Set optional instance template reference
			if instanceTemplateRef := createOptionalInstanceTemplateReference(childConfig.InstanceTemplateId); instanceTemplateRef != nil {
				childInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			// Create choice value
			choiceInstance := childInstance.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
			choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
			choiceValue.SetValue(&childConfig.Value)

			choiceValueOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
			choiceValue.SetOdataType(&choiceValueOdataType)

			// Set optional value template reference with useTemplateDefault if specified
			if childConfig.ValueTemplateId != "" {
				if valueTemplateRef := constructSettingValueTemplateReference(childConfig.ValueTemplateId, childConfig.UseTemplateDefault); valueTemplateRef != nil {
					choiceValue.SetSettingValueTemplateReference(valueTemplateRef)
				}
			}

			// Choice settings in groups typically have no children
			var choiceChildren []models.DeviceManagementConfigurationSettingInstanceable
			choiceValue.SetChildren(choiceChildren)

			choiceInstance.SetChoiceSettingValue(choiceValue)

		case "simple_collection":
			childInstance = models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
			childInstance.SetSettingDefinitionId(&childConfig.SettingDefinitionId)

			collectionOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
			childInstance.SetOdataType(&collectionOdataType)

			// Set optional instance template reference
			if instanceTemplateRef := createOptionalInstanceTemplateReference(childConfig.InstanceTemplateId); instanceTemplateRef != nil {
				childInstance.SetSettingInstanceTemplateReference(instanceTemplateRef)
			}

			// Create collection values
			var collectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
			for _, val := range childConfig.Values {
				stringValue := models.NewDeviceManagementConfigurationStringSettingValue()
				stringValue.SetValue(&val)

				stringOdataType := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
				stringValue.SetOdataType(&stringOdataType)

				collectionValues = append(collectionValues, stringValue)
			}

			collectionInstance := childInstance.(models.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
			collectionInstance.SetSimpleSettingCollectionValue(collectionValues)

		case "group_collection":
			// Handle nested group collections recursively
			if childConfig.ChildConfig == nil {
				return nil, fmt.Errorf("childConfig required for nested group_collection")
			}

			// This would need to recursively call ConstructGroupCollectionSetting
			// For now, return an error as this requires more complex handling
			return nil, fmt.Errorf("nested group collections not yet supported in this builder")

		default:
			return nil, fmt.Errorf("unsupported setting type: %s", childConfig.SettingType)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create child setting %s: %w", childConfig.SettingDefinitionId, err)
		}

		children = append(children, childInstance)
	}

	return children, nil
}

// ConstructSimpleGroupWithMixedChildren creates a single group setting (not collection) with mixed children
// This is for simpler group patterns that don't need collection semantics
func ConstructSimpleGroupWithMixedChildren(
	ctx context.Context,
	settingDefinitionId string,
	children []GroupChildConfig,
	instanceTemplateId string,
	valueTemplateId string,
) (models.DeviceManagementConfigurationSettingable, error) {
	if settingDefinitionId == "" {
		return nil, fmt.Errorf("settingDefinitionId cannot be empty")
	}
	if len(children) == 0 {
		return nil, fmt.Errorf("children cannot be empty")
	}

	tflog.Debug(ctx, "Constructing simple group with mixed children", map[string]interface{}{
		"settingDefinitionId": settingDefinitionId,
		"childrenCount":       len(children),
		"instanceTemplateId":  instanceTemplateId,
		"valueTemplateId":     valueTemplateId,
	})

	// Use the group collection builder with a single group instance
	groupConfigs := [][]GroupChildConfig{children}

	setting, err := ConstructGroupCollectionSetting(ctx, settingDefinitionId, groupConfigs, instanceTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to construct simple group: %w", err)
	}

	tflog.Debug(ctx, "Successfully constructed simple group with mixed children")
	return setting, nil
}
