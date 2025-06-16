package sharedConstructors

import (
	"fmt"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// This package contains reusable builder functions for creating settings catalog objects (such as choice, string, integer, boolean, and collection settings)
// for Microsoft Graph device management configuration policies.
//
//
// # When and Why to Use
//
// Use these builders when you need to:
//   - Programmatically generate device management configuration policy settings for Intune/Endpoint Manager.
//   - Map HCL (Terraform) or other configuration data into the correct Graph API object model.
//
// # How to Use
//
// These builders require a caller function that passes the values from the hcl to the builders.

// Example usage:
//
//   // Construct a choice setting instance
//   setting := builders.ConstructChoiceSettingInstance(
//       "enrollment_autopilot_dpp_deploymentmode",
//       "enrollment_autopilot_dpp_deploymentmode_0",
//       "5180aeab-886e-4589-97d4-40855c646315", // settingInstanceTemplateId
//       "5874c2f6-bcf1-463b-a9eb-bee64e2f2d82", // settingValueTemplateId
//   )
//
//   // Construct a string simple setting instance
//   stringSetting := builders.ConstructStringSimpleSettingInstance(
//       "enrollment_autopilot_dpp_customerrormessage",
//       "Contact your IT department for assistance.",
//       "2ddf0619-2b7a-46de-b29b-c6191e9dda6e", // settingInstanceTemplateId
//       "fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97", // settingValueTemplateId
//   )
//
//   // Construct a boolean choice setting instance
//   boolSetting := builders.ConstructBoolChoiceSettingInstance(
//       "enrollment_autopilot_dpp_allowskip",
//       true,
//       "2a71dc89-0f17-4ba9-bb27-af2521d34710", // settingInstanceTemplateId
//       "a2323e5e-ac56-4517-8847-b0a6fdb467e7", // settingValueTemplateId
//   )
//
//   // Construct a collection setting instance
//   collectionSetting := builders.ConstructSimpleSettingCollectionInstance(
//       "enrollment_autopilot_dpp_allowedappids",
//       []string{"{\"id\":\"app-guid\",\"type\":\"#microsoft.graph.win32LobApp\"}"},
//       "70d22a8a-a03c-4f62-b8df-dded3e327639",
//   )
//
// Complex example with nested settings:
//
//   // 1. Create child settings for a group setting collection
//   var childSettings []models.DeviceManagementConfigurationSettingInstanceable
//
//   // 1.1. Add a choice setting to the children
//   digitsSetting := models.NewDeviceManagementConfigurationChoiceSettingInstance()
//   digitsDefinitionId := "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_digits"
//   digitsSetting.SetSettingDefinitionId(&digitsDefinitionId)
//   digitsOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
//   digitsSetting.SetOdataType(&digitsOdataType)
//
//   digitsChoiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
//   digitsValue := "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_digits_0"
//   digitsChoiceValue.SetValue(&digitsValue)
//   digitsChoiceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
//   digitsChoiceValue.SetOdataType(&digitsChoiceOdataType)
//   digitsChoiceValue.SetChildren([]models.DeviceManagementConfigurationSettingInstanceable{})
//
//   digitsSetting.SetChoiceSettingValue(digitsChoiceValue)
//   childSettings = append(childSettings, digitsSetting)
//
//   // 1.2. Add an integer simple setting to the children
//   minLengthSetting := models.NewDeviceManagementConfigurationSimpleSettingInstance()
//   minLengthDefinitionId := "user_vendor_msft_passportforwork_{tenantid}_policies_pincomplexity_minimumpinlength"
//   minLengthSetting.SetSettingDefinitionId(&minLengthDefinitionId)
//   minLengthOdataType := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
//   minLengthSetting.SetOdataType(&minLengthOdataType)
//
//   minLengthValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
//   intValue := int32(6)
//   minLengthValue.SetValue(&intValue)
//   minLengthOdataType := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
//   minLengthValue.SetOdataType(&minLengthOdataType)
//
//   minLengthSetting.SetSimpleSettingValue(minLengthValue)
//   childSettings = append(childSettings, minLengthSetting)
//
//   // 2. Create a group setting value with the child settings
//   groupValue := builders.ConstructGroupSettingValue(
//       childSettings,
//       "", // settingValueTemplateId (empty if not needed)
//   )
//
//   // 3. Create a group setting collection with the group value
//   var groupValues []models.DeviceManagementConfigurationGroupSettingValueable
//   groupValues = append(groupValues, groupValue)
//
//   groupSetting := builders.ConstructGroupSettingCollectionInstance(
//       "user_vendor_msft_passportforwork_{tenantid}",
//       "", // settingInstanceTemplateId (empty if not needed)
//       groupValues,
//   )
//
//   // 4. Create a choice setting with children using the builder
//   extensionSetting := builders.ConstructChoiceSettingWithChildren(
//       "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~extensions_extensioninstallblocklist",
//       "device_vendor_msft_policy_config_microsoft_edge~policy~microsoft_edge~extensions_extensioninstallblocklist_1",
//       "fb2f16e0-2804-45a0-9982-fe709d59fef8", // settingInstanceTemplateId
//       "e9f334db-ca88-4b09-9ccf-d7b9b3142210", // settingValueTemplateId
//       []models.DeviceManagementConfigurationSettingInstanceable{
//           // You would create and add child settings here, similar to the ones created above
//       },
//   )
//
//   // 5. Add all settings to the settings array
//   var settings []models.DeviceManagementConfigurationSettingable
//   settings = append(settings, groupSetting)
//   settings = append(settings, extensionSetting)
//
//   // 6. Use the settings in your configuration policy
//   configurationPolicy.SetSettings(settings)

// ConstructChoiceSettingInstance creates a choice setting.
func ConstructChoiceSettingInstance(
	settingDefinitionId string,
	value string,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	valuePtr := value
	choiceValue.SetValue(&valuePtr)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&odataTypeValue)

	var children []models.DeviceManagementConfigurationSettingInstanceable
	choiceValue.SetChildren(children)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		choiceValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetChoiceSettingValue(choiceValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructStringSimpleSettingInstance creates a simple string setting.
func ConstructStringSimpleSettingInstance(
	settingDefinitionId string,
	value string,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	simpleSettingValue := models.NewDeviceManagementConfigurationStringSettingValue()
	valuePtr := value
	simpleSettingValue.SetValue(&valuePtr)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	simpleSettingValue.SetOdataType(&odataTypeValue)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		simpleSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetSimpleSettingValue(simpleSettingValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructIntSimpleSettingInstance creates an integer setting.
func ConstructIntSimpleSettingInstance(
	settingDefinitionId string,
	value int64,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	intSettingValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
	intValue := int32(value)
	intSettingValue.SetValue(&intValue)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	intSettingValue.SetOdataType(&odataTypeValue)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		intSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetSimpleSettingValue(intSettingValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructSimpleSettingCollectionInstance creates a collection setting for string values.
func ConstructSimpleSettingCollectionInstance(
	settingDefinitionId string,
	values []string,
	settingInstanceTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	var simpleSettingCollectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
	for _, val := range values {
		simpleSettingValue := models.NewDeviceManagementConfigurationStringSettingValue()
		valuePtr := val
		simpleSettingValue.SetValue(&valuePtr)

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
		simpleSettingValue.SetOdataType(&odataTypeValue)

		simpleSettingCollectionValues = append(simpleSettingCollectionValues, simpleSettingValue)
	}

	settingInstance.SetSimpleSettingCollectionValue(simpleSettingCollectionValues)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructBoolChoiceSettingInstance creates a boolean setting using the choice setting format.
func ConstructBoolChoiceSettingInstance(settingDefinitionId string, value bool, settingInstanceTemplateId string, settingValueTemplateId string) models.DeviceManagementConfigurationSettingable {
	// Convert bool to appropriate format for the API
	// The Graph API uses strings with _0 or _1 suffixes to represent boolean values
	strValue := fmt.Sprintf("%s_0", settingDefinitionId) // Default to false value
	if value {
		strValue = fmt.Sprintf("%s_1", settingDefinitionId) // True value
	}

	// Use ConstructChoiceSettingInstance since boolean settings are represented as choice settings in the API
	return ConstructChoiceSettingInstance(settingDefinitionId, strValue, settingInstanceTemplateId, settingValueTemplateId)
}

// ConstructChoiceSettingWithChildren creates a choice setting with children settings.
func ConstructChoiceSettingWithChildren(
	settingDefinitionId string,
	value string,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
	children []models.DeviceManagementConfigurationSettingInstanceable,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	valuePtr := value
	choiceValue.SetValue(&valuePtr)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&odataTypeValue)

	choiceValue.SetChildren(children)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		choiceValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetChoiceSettingValue(choiceValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructSecretSimpleSettingInstance creates a secret string setting.
func ConstructSecretSimpleSettingInstance(
	settingDefinitionId string,
	value string,
	valueState string,
	settingInstanceTemplateId string,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	secretSettingValue := models.NewDeviceManagementConfigurationSecretSettingValue()
	valuePtr := value
	secretSettingValue.SetValue(&valuePtr)

	// Set the value state (e.g., "invalid", "notEncrypted", "encryptedValueToken")
	// Parse the string valueState to the correct enum type
	parsedValueState, err := models.ParseDeviceManagementConfigurationSecretSettingValueState(valueState)
	if err == nil {
		if typedValueState, ok := parsedValueState.(*models.DeviceManagementConfigurationSecretSettingValueState); ok && typedValueState != nil {
			secretSettingValue.SetValueState(typedValueState)
		}
	}

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
	secretSettingValue.SetOdataType(&odataTypeValue)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		secretSettingValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	settingInstance.SetSimpleSettingValue(secretSettingValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructGroupSettingCollectionInstance creates a group setting collection with child settings.
func ConstructGroupSettingCollectionInstance(
	settingDefinitionId string,
	settingInstanceTemplateId string,
	groupValues []models.DeviceManagementConfigurationGroupSettingValueable,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	settingInstance.SetGroupSettingCollectionValue(groupValues)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructGroupSettingValue creates a group setting value with child settings.
func ConstructGroupSettingValue(
	children []models.DeviceManagementConfigurationSettingInstanceable,
	settingValueTemplateId string,
) models.DeviceManagementConfigurationGroupSettingValueable {
	groupValue := models.NewDeviceManagementConfigurationGroupSettingValue()

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationGroupSettingValue"
	groupValue.SetOdataType(&odataTypeValue)

	groupValue.SetChildren(children)

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		groupValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	return groupValue
}

// ConstructChoiceSettingCollectionInstance creates a choice setting collection.
func ConstructChoiceSettingCollectionInstance(
	settingDefinitionId string,
	settingInstanceTemplateId string,
	choiceValues []models.DeviceManagementConfigurationChoiceSettingValueable,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	settingInstance.SetChoiceSettingCollectionValue(choiceValues)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructChoiceSettingValue creates a choice setting value.
func ConstructChoiceSettingValue(
	value string,
	settingValueTemplateId string,
	children []models.DeviceManagementConfigurationSettingInstanceable,
) models.DeviceManagementConfigurationChoiceSettingValueable {
	choiceValue := models.NewDeviceManagementConfigurationChoiceSettingValue()
	valuePtr := value
	choiceValue.SetValue(&valuePtr)

	odataTypeValue := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	choiceValue.SetOdataType(&odataTypeValue)

	choiceValue.SetChildren(children)

	if settingValueTemplateId != "" {
		settingValueTemplateReference := models.NewDeviceManagementConfigurationSettingValueTemplateReference()
		settingValueTemplateReference.SetSettingValueTemplateId(&settingValueTemplateId)
		choiceValue.SetSettingValueTemplateReference(settingValueTemplateReference)
	}

	return choiceValue
}

// ConstructIntCollectionSettingInstance creates a collection of integer values.
func ConstructIntCollectionSettingInstance(
	settingDefinitionId string,
	values []int64,
	settingInstanceTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	var simpleSettingCollectionValues []models.DeviceManagementConfigurationSimpleSettingValueable
	for _, val := range values {
		intSettingValue := models.NewDeviceManagementConfigurationIntegerSettingValue()
		intValue := int32(val)
		intSettingValue.SetValue(&intValue)

		odataTypeValue := "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
		intSettingValue.SetOdataType(&odataTypeValue)

		simpleSettingCollectionValues = append(simpleSettingCollectionValues, intSettingValue)
	}

	settingInstance.SetSimpleSettingCollectionValue(simpleSettingCollectionValues)
	setting.SetSettingInstance(settingInstance)

	return setting
}

// ConstructSimpleSettingInstance creates a generic simple setting with the provided simple setting value.
func ConstructSimpleSettingInstance(
	settingDefinitionId string,
	settingValue models.DeviceManagementConfigurationSimpleSettingValueable,
	settingInstanceTemplateId string,
) models.DeviceManagementConfigurationSettingable {
	setting := models.NewDeviceManagementConfigurationSetting()

	settingInstance := models.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefinitionIdValue := settingDefinitionId
	settingInstance.SetSettingDefinitionId(&settingDefinitionIdValue)

	odataTypeInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeInstance)

	if settingInstanceTemplateId != "" {
		settingInstanceTemplateReference := models.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		settingInstanceTemplateReference.SetSettingInstanceTemplateId(&settingInstanceTemplateId)
		settingInstance.SetSettingInstanceTemplateReference(settingInstanceTemplateReference)
	}

	settingInstance.SetSimpleSettingValue(settingValue)
	setting.SetSettingInstance(settingInstance)

	return setting
}
