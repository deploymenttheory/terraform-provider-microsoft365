package convert

import "github.com/hashicorp/terraform-plugin-framework/types"

// GraphToFrameworkBoolChoice reads the settings-catalog choice encoding written
// by ConstructBoolChoiceSettingInstance. Unknown choices remain null, not false.
func GraphToFrameworkBoolChoice(value *string, settingDefinitionID string) types.Bool {
	if value == nil || settingDefinitionID == "" {
		return types.BoolNull()
	}
	switch *value {
	case settingDefinitionID + "_0":
		return types.BoolValue(false)
	case settingDefinitionID + "_1":
		return types.BoolValue(true)
	default:
		return types.BoolNull()
	}
}
