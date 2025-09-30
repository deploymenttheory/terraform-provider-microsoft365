package attr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ListValue creates a types.List from attr.Value elements.
// Panics if creation fails, which should only happen with programming errors.
func ListValue(elementType attr.Type, elements []attr.Value) types.List {
	list, err := types.ListValue(elementType, elements)
	if err != nil {
		panic(err)
	}
	return list
}

// ListValueFrom creates a types.List from a Go slice.
// Panics if creation fails, which should only happen with programming errors.
func ListValueFrom(ctx context.Context, elementType attr.Type, elements any) types.List {
	list, diags := types.ListValueFrom(ctx, elementType, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert to types.List", map[string]any{
			"error": diags.Errors()[0].Detail(),
		})
		panic(diags.Errors()[0].Detail())
	}
	return list
}

// ListNullIfEmpty returns ListNull if the slice is empty, otherwise creates a list with the values.
func ListNullIfEmpty(elementType attr.Type, values []attr.Value) types.List {
	if len(values) == 0 {
		return types.ListNull(elementType)
	}

	list, err := types.ListValue(elementType, values)
	if err != nil {
		return types.ListNull(elementType)
	}
	return list
}
