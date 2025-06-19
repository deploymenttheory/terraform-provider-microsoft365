package attr

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// StringListElements converts a basetypes.ListValue to a slice of strings.
// Returns nil if the list is null or unknown.
func StringListElements(list basetypes.ListValue) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elements := list.Elements()
	result := make([]string, 0, len(elements))

	for _, element := range elements {
		if stringVal, ok := element.(basetypes.StringValue); ok {
			if stringVal.IsNull() || stringVal.IsUnknown() {
				continue
			}
			result = append(result, stringVal.ValueString())
		}
	}

	return result
}

// StringSetElements converts a basetypes.SetValue to a slice of strings.
// Returns nil if the set is null or unknown.
func StringSetElements(set basetypes.SetValue) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	elements := set.Elements()
	result := make([]string, 0, len(elements))

	for _, element := range elements {
		if stringVal, ok := element.(basetypes.StringValue); ok {
			if stringVal.IsNull() || stringVal.IsUnknown() {
				continue
			}
			result = append(result, stringVal.ValueString())
		}
	}

	return result
}

// StringMapElements converts a basetypes.MapValue to a map of strings.
// Returns nil if the map is null or unknown.
func StringMapElements(m basetypes.MapValue) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}

	elements := m.Elements()
	result := make(map[string]string, len(elements))

	for key, element := range elements {
		if stringVal, ok := element.(basetypes.StringValue); ok {
			if stringVal.IsNull() || stringVal.IsUnknown() {
				continue
			}
			result[key] = stringVal.ValueString()
		}
	}

	return result
}

// ObjectSetElements converts a basetypes.SetValue containing objects to a slice of maps.
// Returns nil if the set is null or unknown.
func ObjectSetElements(set basetypes.SetValue) []map[string]attr.Value {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	elements := set.Elements()
	result := make([]map[string]attr.Value, 0, len(elements))

	for _, element := range elements {
		if objVal, ok := element.(basetypes.ObjectValue); ok {
			if objVal.IsNull() || objVal.IsUnknown() {
				continue
			}
			result = append(result, objVal.Attributes())
		}
	}

	return result
}

// ObjectListElements converts a basetypes.ListValue containing objects to a slice of maps.
// Returns nil if the list is null or unknown.
func ObjectListElements(list basetypes.ListValue) []map[string]attr.Value {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	elements := list.Elements()
	result := make([]map[string]attr.Value, 0, len(elements))

	for _, element := range elements {
		if objVal, ok := element.(basetypes.ObjectValue); ok {
			if objVal.IsNull() || objVal.IsUnknown() {
				continue
			}
			result = append(result, objVal.Attributes())
		}
	}

	return result
}

// ObjectMapElements converts a basetypes.MapValue containing objects to a map of maps.
// Returns nil if the map is null or unknown.
func ObjectMapElements(m basetypes.MapValue) map[string]map[string]attr.Value {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}

	elements := m.Elements()
	result := make(map[string]map[string]attr.Value, len(elements))

	for key, element := range elements {
		if objVal, ok := element.(basetypes.ObjectValue); ok {
			if objVal.IsNull() || objVal.IsUnknown() {
				continue
			}
			result[key] = objVal.Attributes()
		}
	}

	return result
}
