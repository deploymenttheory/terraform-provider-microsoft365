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
//       "5180aeab-886e-4589-97d4-40855c646315",
//       "5874c2f6-bcf1-463b-a9eb-bee64e2f2d82",
//   )
//
//   // Construct a string simple setting instance
//   stringSetting := builders.ConstructStringSimpleSettingInstance(
//       "enrollment_autopilot_dpp_customerrormessage",
//       "Contact your IT department for assistance.",
//       "2ddf0619-2b7a-46de-b29b-c6191e9dda6e",
//       "fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97",
//   )
//
//   // Construct a boolean choice setting instance
//   boolSetting := builders.ConstructBoolChoiceSettingInstance(
//       "enrollment_autopilot_dpp_allowskip",
//       true,
//       "2a71dc89-0f17-4ba9-bb27-af2521d34710",
//       "a2323e5e-ac56-4517-8847-b0a6fdb467e7",
//   )
//
//   // Construct a collection setting instance
//   collectionSetting := builders.ConstructSimpleSettingCollectionInstance(
//       "enrollment_autopilot_dpp_allowedappids",
//       []string{"{\"id\":\"app-guid\",\"type\":\"#microsoft.graph.win32LobApp\"}"},
//       "70d22a8a-a03c-4f62-b8df-dded3e327639",
//   )
//

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
