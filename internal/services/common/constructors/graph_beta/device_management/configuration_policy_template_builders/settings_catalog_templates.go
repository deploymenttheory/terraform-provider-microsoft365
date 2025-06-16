package configurationPolicyTemplateBuilders

// ========================================================================================
// CONFIGURATION TYPES
// ========================================================================================

// HCLSettingsInput represents the flat key-value input from HCL
type HCLSettingsInput map[string]string

// SettingDefinition defines how a setting should be processed
type SettingDefinition struct {
	SettingType        string   // "simple_choice", "simple_string", "simple_integer", "choice_with_child", etc.
	ValueType          string   // "string", "integer", "choice", "boolean"
	ParentSetting      string   // For child settings, what's the parent setting ID
	ChildSettings      []string // For parent settings, what are the child setting IDs
	InstanceTemplateID string   // Template ID for the setting instance
	ValueTemplateID    string   // Template ID for the setting value
	UseTemplateDefault bool     // Whether to use template default
}

// SettingsCatalogTemplate contains the schema definitions for all known settings
type SettingsCatalogTemplate struct {
	definitions map[string]SettingDefinition
}

// ========================================================================================
// REGISTRY INITIALIZATION
// ========================================================================================

// NewSettingsCatalogTemplate creates a new registry with known setting definitions
func NewSettingsCatalogTemplate() *SettingsCatalogTemplate {
	registry := &SettingsCatalogTemplate{
		definitions: make(map[string]SettingDefinition),
	}

	// Register LAPS settings (example)
	registry.registerLAPSSettings()

	// Add more setting registrations here as needed
	// registry.registerFSLogixSettings()
	// registry.registerWindowsASRSettings()
	// registry.registermacOSDDMSettings()

	return registry
}

// registerLAPSSettings registers all Windows LAPS setting definitions
func (r *SettingsCatalogTemplate) registerLAPSSettings() {
	// Setting 0: Choice with integer child
	r.definitions["device_vendor_msft_laps_policies_backupdirectory"] = SettingDefinition{
		SettingType:        "choice_with_child",
		ValueType:          "choice",
		ChildSettings:      []string{"device_vendor_msft_laps_policies_passwordagedays_aad"},
		InstanceTemplateID: "a3270f64-e493-499d-8900-90290f61ed8a",
		ValueTemplateID:    "4d90f03d-e14c-43c4-86da-681da96a2f92",
		UseTemplateDefault: false,
	}

	r.definitions["device_vendor_msft_laps_policies_passwordagedays_aad"] = SettingDefinition{
		SettingType:        "child_integer",
		ValueType:          "integer",
		ParentSetting:      "device_vendor_msft_laps_policies_backupdirectory",
		ValueTemplateID:    "4d90f03d-e14c-43c4-86da-681da96a2f92",
		UseTemplateDefault: false,
	}

	// Setting 1: Simple string
	r.definitions["device_vendor_msft_laps_policies_administratoraccountname"] = SettingDefinition{
		SettingType:        "simple_string",
		ValueType:          "string",
		InstanceTemplateID: "d3d7d492-0019-4f56-96f8-1967f7deabeb",
		ValueTemplateID:    "992c7fce-f9e4-46ab-ac11-e167398859ea",
		UseTemplateDefault: false,
	}

	// Setting 2: Simple choice
	r.definitions["device_vendor_msft_laps_policies_passwordcomplexity"] = SettingDefinition{
		SettingType:        "simple_choice",
		ValueType:          "choice",
		InstanceTemplateID: "8a7459e8-1d1c-458a-8906-7b27d216de52",
		ValueTemplateID:    "aa883ab5-625e-4e3b-b830-a37a4bb8ce01",
		UseTemplateDefault: false,
	}

	// Setting 3: Simple integer
	r.definitions["device_vendor_msft_laps_policies_passwordlength"] = SettingDefinition{
		SettingType:        "simple_integer",
		ValueType:          "integer",
		InstanceTemplateID: "da7a1dbd-caf7-4341-ab63-ece6f994ff02",
		ValueTemplateID:    "d08f1266-5345-4f53-8ae1-4c20e6cb5ec9",
		UseTemplateDefault: false,
	}

	// Setting 4: Simple choice
	r.definitions["device_vendor_msft_laps_policies_postauthenticationactions"] = SettingDefinition{
		SettingType:        "simple_choice",
		ValueType:          "choice",
		InstanceTemplateID: "d9282eb1-d187-42ae-b366-7081f32dcfff",
		ValueTemplateID:    "68ff4f78-baa8-4b32-bf3d-5ad5566d8142",
		UseTemplateDefault: false,
	}

	// Setting 5: Simple integer
	r.definitions["device_vendor_msft_laps_policies_postauthenticationresetdelay"] = SettingDefinition{
		SettingType:        "simple_integer",
		ValueType:          "integer",
		InstanceTemplateID: "a9e21166-4055-4042-9372-efaf3ef41868",
		ValueTemplateID:    "0deb6aee-8dac-40c4-a9dd-c3718e5c1d52",
		UseTemplateDefault: false,
	}
}
