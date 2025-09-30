package attr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ObjectValue creates a types.Object with known values from a map.
// Panics if creation fails, which should only happen with programming errors.
func ObjectValue(attrTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	object, err := types.ObjectValue(attrTypes, values)
	if err != nil {
		panic(err)
	}
	return object
}

// ObjectNullIfEmpty returns ObjectNull if the map is empty, otherwise creates an object with the values.
func ObjectNullIfEmpty(attrTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	if len(values) == 0 {
		return types.ObjectNull(attrTypes)
	}

	object, err := types.ObjectValue(attrTypes, values)
	if err != nil {
		return types.ObjectNull(attrTypes)
	}
	return object
}

// ObjectValueFrom creates a types.Object from a Go struct.
// Panics if creation fails, which should only happen with programming errors.
func ObjectValueFrom(ctx context.Context, attrTypes map[string]attr.Type, value any) types.Object {
	object, diags := types.ObjectValueFrom(ctx, attrTypes, value)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert to types.Object", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		panic(diags.Errors()[0].Detail())
	}
	return object
}

// GetObjectAttr safely retrieves an attribute from an object by name.
// Returns the default value if the object is null, unknown, or the attribute doesn't exist.
func GetObjectAttr(obj basetypes.ObjectValue, name string, defaultValue attr.Value) attr.Value {
	if obj.IsNull() || obj.IsUnknown() {
		return defaultValue
	}

	attrs := obj.Attributes()
	if val, ok := attrs[name]; ok {
		return val
	}
	return defaultValue
}

// ObjectSetFromSlice builds a set of objects from a slice
// Returns types.SetNull() if the length is zero
func ObjectSetFromSlice(ctx context.Context, attrTypes map[string]attr.Type, valueFunc func(int) map[string]attr.Value, length int) types.Set {
	if length == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	values := make([]attr.Value, length)
	for i := 0; i < length; i++ {
		obj, err := types.ObjectValue(attrTypes, valueFunc(i))
		if err != nil {
			tflog.Error(ctx, "Failed to create object value", map[string]any{
				"error": err,
				"index": i,
			})
			return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
		}
		values[i] = obj
	}

	set, diags := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create set value", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}
	return set
}
