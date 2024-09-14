package construct

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func DebugPrintStruct(ctx context.Context, prefix string, v interface{}) {
	m := structToMap(reflect.ValueOf(v))

	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling struct to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, prefix, map[string]interface{}{
		"data": string(jsonData),
	})
}

func structToMap(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return structToMap(v.Elem())
	case reflect.Struct:
		return handleStruct(v)
	case reflect.Slice, reflect.Array:
		return handleSlice(v)
	case reflect.Map:
		return handleMap(v)
	default:
		return handlePrimitive(v)
	}
}

func handleStruct(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue // Skip unexported fields
		}

		tag := field.Tag.Get("tfsdk")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = field.Name
		}

		fieldValue := v.Field(i)

		var value interface{}
		if fieldValue.Type().Implements(reflect.TypeOf((*attr.Value)(nil)).Elem()) {
			value = handleTerraformValue(fieldValue.Interface().(attr.Value))
		} else {
			value = structToMap(fieldValue)
		}

		if value != nil {
			result[tag] = value
		}
	}

	return result
}

func handleSlice(v reflect.Value) []interface{} {
	var result []interface{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Type().Implements(reflect.TypeOf((*attr.Value)(nil)).Elem()) {
			// If the element implements attr.Value, use handleTerraformValue
			value := handleTerraformValue(elem.Interface().(attr.Value))
			if value != nil {
				result = append(result, value)
			}
		} else {
			// For other types, use structToMap
			value := structToMap(elem)
			if value != nil {
				result = append(result, value)
			}
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func handleMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		value := structToMap(v.MapIndex(key))
		if value != nil {
			result[fmt.Sprint(key.Interface())] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func handlePrimitive(v reflect.Value) interface{} {
	if v.Type().Implements(reflect.TypeOf((*attr.Value)(nil)).Elem()) {
		return handleTerraformValue(v.Interface().(attr.Value))
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.String:
		return v.String()
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func handleTerraformValue(v attr.Value) interface{} {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	switch value := v.(type) {
	case types.String:
		return value.ValueString()
	case types.Int64:
		return value.ValueInt64()
	case types.Float64:
		return value.ValueFloat64()
	case types.Bool:
		return value.ValueBool()
	case types.List:
		return handleTerraformList(value)
	case types.Set:
		return handleTerraformSet(value)
	case types.Map:
		return handleTerraformMap(value)
	case types.Object:
		return handleTerraformObject(value)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func handleTerraformList(list types.List) interface{} {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	elements := list.Elements()
	result := make([]interface{}, 0, len(elements))
	for _, elem := range elements {
		if value := handleTerraformValue(elem); value != nil {
			result = append(result, value)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func handleTerraformSet(set types.Set) interface{} {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	elements := set.Elements()
	result := make([]interface{}, 0, len(elements))
	for _, elem := range elements {
		if value := handleTerraformValue(elem); value != nil {
			result = append(result, value)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func handleTerraformMap(m types.Map) interface{} {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	elements := m.Elements()
	result := make(map[string]interface{})
	for k, v := range elements {
		if value := handleTerraformValue(v); value != nil {
			result[k] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func handleTerraformObject(obj types.Object) interface{} {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}
	attrs := obj.Attributes()
	result := make(map[string]interface{})
	for k, v := range attrs {
		if value := handleTerraformValue(v); value != nil {
			result[k] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
