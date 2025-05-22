# Device Management Configuration Setting (Settings Catalog) Map

This map represents the current supported odata hierarchy supported within this terraform resource.
It will be updated as new requirements are identified.

**Note:** All setting instances include:

- `@odata.type` field
- `settingDefinitionId` field  
- Optional `settingInstanceTemplateReference` with `settingInstanceTemplateId`

**Note:** All setting values include:

- `@odata.type` field
- Optional `settingValueTemplateReference` with `settingValueTemplateId` and `useTemplateDefault` boolean
- Optional `valueState` field (for simple setting values)

```bash
switch detail.SettingInstance.ODataType:
├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│   └── choiceSettingValue.children[] switch ODataType:
│       ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│       │   └── simpleSettingValue switch ODataType:
│       │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│       │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
│       ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│       │   └── choiceSettingValue
│       │       └── children[] switch ODataType:
│       │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│       │           │   └── choiceSettingValue.children[] (empty array)
│       │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│       │               └── simpleSettingValue switch ODataType:
│       │                   ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       │                   ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│       │                   └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
│       ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
│       │   └── groupSettingCollectionValue[]
│       │       └── children[] switch ODataType:
│       │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│       │           │   └── simpleSettingValue switch ODataType:
│       │           │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       │           │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│       │           │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
│       │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
│       │           │   └── simpleSettingCollectionValue[]
│       │           │       └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
│       │           └── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│       │               └── choiceSettingValue
│       │                   └── children[] switch ODataType:
│       │                       └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│       │                           └── simpleSettingValue switch ODataType:
│       │                               ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       │                               ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│       │                               └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
│       └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
│           └── simpleSettingCollectionValue[]
│               └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
│
├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
│   └── simpleSettingCollectionValue[]
│       └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
│
├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│   └── simpleSettingValue switch ODataType:
│       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
│
├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
│   └── choiceSettingCollectionValue[]
│       └── children[] switch ODataType:
│           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│           │   └── simpleSettingValue handle types:
│           │       ├── string -> StringSettingValue
│           │       ├── float64 -> IntegerSettingValue
│           │       └── secret -> SecretSettingValue
│           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
│               └── simpleSettingCollectionValue[]
│                   └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
│
└── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
    └── groupSettingCollectionValue[]
        └── children[] switch ODataType:
            ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance" (Level 2)
            │   └── groupSettingCollectionValue[]
            │       └── children[] switch ODataType:
            │           ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance" (Level 3)
            │           │   └── groupSettingCollectionValue[]
            │           │       └── children[] switch ODataType:
            │           │           ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance" (Level 4)
            │           │           │   └── groupSettingCollectionValue[]
            │           │           │       └── children[] switch ODataType:
            │           │           │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │           │           │           │   └── choiceSettingValue
            │           │           │           │       └── children[] switch ODataType:
            │           │           │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance" (Level 5)
            │           │           │           │               └── simpleSettingValue switch ODataType:
            │           │           │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │           │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │           │           │                   └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           │           │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │           │           │   └── simpleSettingValue switch ODataType:
            │           │           │           │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │           │           │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │           │           │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │           │           │               └── simpleSettingCollectionValue[]
            │           │           │                   └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
            │           │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │           │           │   └── choiceSettingValue
            │           │           │       └── children[] switch ODataType:
            │           │           │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │           │           │           │   └── choiceSettingValue.children[] (empty array)
            │           │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │           │               └── simpleSettingValue switch ODataType:
            │           │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │           │                   └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │           │   └── simpleSettingValue switch ODataType:
            │           │           │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │           │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │           │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │           │               └── simpleSettingCollectionValue[]
            │           │                   └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
            │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │           │   └── choiceSettingValue
            │           │       └── children[] switch ODataType:
            │           │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │           │           │   └── choiceSettingValue.children[] (empty array)
            │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │               └── simpleSettingValue switch ODataType:
            │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │                   └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │   └── simpleSettingValue switch ODataType:
            │           │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │               └── simpleSettingCollectionValue[]
            │                   └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
            ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │   └── simpleSettingCollectionValue[]
            │       └── (@ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue")
            ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │   └── simpleSettingValue switch ODataType:
            │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            └── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                └── choiceSettingValue
                    └── children[] switch ODataType:
                        ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                        │   └── choiceSettingValue
                        │       └── children[] (empty array)
                        └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                            └── simpleSettingValue switch ODataType:
                                ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                                ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                                └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
```

## Template Reference Structure

All setting instances and values may include template references:

### Setting Instance Template Reference
```
settingInstanceTemplateReference: {
    settingInstanceTemplateId: string
}
```

### Setting Value Template Reference  
```
settingValueTemplateReference: {
    settingValueTemplateId: string
    useTemplateDefault: boolean
}
```

## Additional Properties

- **valueState**: Optional string field present on simple setting values
- **@odata.type**: Required on all setting instances and values
- **settingDefinitionId**: Required on all setting instances
- **value**: The actual setting value (type varies by setting value type)

## Supported Value Types

- **String**: `#microsoft.graph.deviceManagementConfigurationStringSettingValue`
- **Integer**: `#microsoft.graph.deviceManagementConfigurationIntegerSettingValue`  
- **Secret**: `#microsoft.graph.deviceManagementConfigurationSecretSettingValue`

## Maximum Nesting Levels

The structure supports up to 5 levels of nesting within group setting collections, with the deepest level being simple settings within choice settings within the most nested group collections