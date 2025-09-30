package attr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapValue creates a types.Map from attr.Value elements.
// Panics if creation fails, which should only happen with programming errors.
func MapValue(elementType attr.Type, elements map[string]attr.Value) types.Map {
	m, err := types.MapValue(elementType, elements)
	if err != nil {
		panic(err)
	}
	return m
}

// MapValueFrom creates a types.Map from a Go map.
// Panics if creation fails, which should only happen with programming errors.
func MapValueFrom(ctx context.Context, elementType attr.Type, elements any) types.Map {
	m, diags := types.MapValueFrom(ctx, elementType, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert to types.Map", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		panic(diags.Errors()[0].Detail())
	}
	return m
}

// MapNullIfEmpty returns MapNull if the map is empty, otherwise creates a map with the values.
func MapNullIfEmpty(elementType attr.Type, values map[string]attr.Value) types.Map {
	if len(values) == 0 {
		return types.MapNull(elementType)
	}

	m, err := types.MapValue(elementType, values)
	if err != nil {
		return types.MapNull(elementType)
	}
	return m
}
