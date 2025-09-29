package attr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SetValue creates a types.Set from attr.Value elements.
// Panics if creation fails, which should only happen with programming errors.
func SetValue(elementType attr.Type, elements []attr.Value) types.Set {
	set, err := types.SetValue(elementType, elements)
	if err != nil {
		panic(err)
	}
	return set
}

// SetValueFrom creates a types.Set from a Go slice.
// Panics if creation fails, which should only happen with programming errors.
func SetValueFrom(ctx context.Context, elementType attr.Type, elements interface{}) types.Set {
	set, diags := types.SetValueFrom(ctx, elementType, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert to types.Set", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		panic(diags.Errors()[0].Detail())
	}
	return set
}

// SetNullIfEmpty returns SetNull if the slice is empty, otherwise creates a set with the values.
func SetNullIfEmpty(ctx context.Context, elementType attr.Type, values []attr.Value) types.Set {
	if len(values) == 0 {
		return types.SetNull(elementType)
	}

	set, diags := types.SetValue(elementType, values)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create set value", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		return types.SetNull(elementType)
	}
	return set
}
