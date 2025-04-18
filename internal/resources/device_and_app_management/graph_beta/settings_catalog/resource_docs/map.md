# Device Management Configuration Setting (Settings Catalog) Map

This map represents the current supported odata hierarachy supported within this terraform resource.
It will updated as new requirements are identified.

switch detail.SettingInstance.ODataType:
├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│   └── choiceSettingValue.children[] switch ODataType:
│       ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│       │   └── simpleSettingValue switch ODataType:
│       │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       │       └── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│       ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│       │   └── choiceSettingValue
│       │       └── children[] switch ODataType:
│       │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
│       │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
│       ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
│       │   └── groupSettingCollectionValue[]
│       │       └── children if ODataType == "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
│       │           └── simpleSettingValue switch ODataType:
│       │               ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
│       │               └── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
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
│       └── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
│
├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
│   └── choiceSettingCollectionValue[]
│       └── children[].simpleSettingValue handle types:
│           ├── string -> StringSettingValue
│           └── float64 -> IntegerSettingValue
│
└── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
    └── groupSettingCollectionValue[]
        └── children[] switch ODataType:
            ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
            │   └── groupSettingCollectionValue[]
            │       └── children[] switch ODataType:
            │           ├── case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
            │           │   └── groupSettingCollectionValue[]
            │           │       └── children[] switch ODataType:
            │           │           ├── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │           │           │   └── choiceSettingValue
            │           │           │       └── children[] switch ODataType:
            │           │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │           │               └── simpleSettingValue switch ODataType:
            │           │           │                   ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │           │                   └── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │           │   └── simpleSettingValue switch ODataType:
            │           │           │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │           │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │           │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           │           └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │           │               └── simpleSettingCollectionValue[] // String values only
            │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │           │   └── simpleSettingValue switch ODataType:
            │           │       ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │           │       ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │           │       └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            │           ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │           │   └── simpleSettingCollectionValue[] // String values only
            │           └── case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            │               └── choiceSettingValue
            │                   └── children[] switch ODataType:
            │                       └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
            │                           └── simpleSettingValue switch ODataType:
            │                               ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            │                               ├── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
            │                               └── case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            ├── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
            │   └── simpleSettingCollectionValue[] // String values only
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
                        │       └── children[] // Empty array
                        └── case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                            └── simpleSettingValue switch ODataType:
                                ├── case "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                                └── case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
