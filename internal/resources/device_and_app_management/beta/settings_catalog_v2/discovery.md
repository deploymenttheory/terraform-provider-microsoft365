# Device Management Configuration Instance Types

## 1. Choice Setting Instance
```go
settingInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
```
Supported Value Type:
```go
choiceSettingValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
// Properties:
// - value: String (e.g., "softwareupdate_notifications_true")
// - children: Can contain other setting instances
```
Example Usage:
- Boolean settings (true/false)
- Enumerated options (e.g., filter grades, proxy types)
- Mode selections (e.g., calculator modes)

## 2. Simple Setting Instance
```go
settingInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
```
Supported Value Types:
```go
// Integer Values
simpleSettingValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
// Example: periodindays, port numbers

// String Values
simpleSettingValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
// Example: serveraddress, username

// Secret Values
simpleSettingValue := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
valueState := graphmodels.NOTENCRYPTED_DEVICEMANAGEMENTCONFIGURATIONSECRETSETTINGVALUESTATE
// Example: passwords, security tokens
```

## 3. Simple Setting Collection Instance
```go
settingInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
```
Supported Value Types:
```go
simpleSettingCollectionValue := []graphmodels.DeviceManagementConfigurationSimpleSettingValueable{
    // Collections of:
    graphmodels.NewDeviceManagementConfigurationStringSettingValue()  // For string arrays
}
```
Example Usage:
- Lists of domains (allowed/denied)
- Exception lists
- Multiple string configurations

## 4. Group Setting Collection Instance
```go
settingInstance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
```
Supported Value Types:
```go
groupSettingCollectionValue := []graphmodels.DeviceManagementConfigurationGroupSettingValueable{
    graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
}
```
Properties:
- Can contain children of any other instance type
- Used for logical grouping and hierarchy
- Supports nested configurations

## Common Implementation Notes

1. All instances require a settingDefinitionId:
```go
settingInstance.SetSettingDefinitionId(&settingDefinitionId)
```

2. Value Assignment Pattern:
```go
// For simple values
settingInstance.SetSimpleSettingValue(simpleSettingValue)

// For choice values
settingInstance.SetChoiceSettingValue(choiceSettingValue)

// For collections
settingInstance.SetSimpleSettingCollectionValue(simpleSettingCollectionValue)

// For groups
settingInstance.SetGroupSettingCollectionValue(groupSettingCollectionValue)
```

3. Children Pattern:
```go
children := []graphmodels.DeviceManagementConfigurationSettingInstanceable{
    // Can contain any instance type
}
```

This structure represents all observed instance types and their supported value types from the provided examples, focused specifically on the technical implementation details and type relationships.