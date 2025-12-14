package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapEnumCollectionToSet converts an enum slice to a Terraform Framework set of strings
// Preserves empty slices as empty sets to maintain Terraform state consistency
func mapEnumCollectionToSet[T fmt.Stringer](ctx context.Context, enums []T, fieldName string) types.Set {
	values := make([]string, len(enums))
	for i, enum := range enums {
		values[i] = enum.String()
	}

	elemType := types.StringType
	if len(values) == 0 {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	elements := make([]attr.Value, len(values))
	for i, v := range values {
		elements[i] = types.StringValue(v)
	}

	return types.SetValueMust(elemType, elements)
}

// mapStringSliceToSetPreserveEmpty converts a string slice to a Terraform Framework set
// Preserves empty slices as empty sets to maintain Terraform state consistency
func mapStringSliceToSetPreserveEmpty(ctx context.Context, values []string) types.Set {
	elemType := types.StringType
	if len(values) == 0 {
		return types.SetValueMust(elemType, []attr.Value{})
	}

	elements := make([]attr.Value, len(values))
	for i, v := range values {
		elements[i] = types.StringValue(v)
	}

	return types.SetValueMust(elemType, elements)
}

// mapEnumCollectionToSetNullIfEmpty converts an enum slice to a Terraform Framework set of strings
// Returns null for empty slices (used for fields where API removes empty arrays)
func mapEnumCollectionToSetNullIfEmpty[T fmt.Stringer](ctx context.Context, enums []T, fieldName string) types.Set {
	if len(enums) == 0 {
		return types.SetNull(types.StringType)
	}

	values := make([]string, len(enums))
	for i, enum := range enums {
		values[i] = enum.String()
	}

	elements := make([]attr.Value, len(values))
	for i, v := range values {
		elements[i] = types.StringValue(v)
	}

	return types.SetValueMust(types.StringType, elements)
}

// mapStringSliceToSetNullIfEmpty converts a string slice to a Terraform Framework set
// Returns null for empty slices (used for optional fields where API returns [] as default)
func mapStringSliceToSetNullIfEmpty(ctx context.Context, values []string) types.Set {
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}

	elements := make([]attr.Value, len(values))
	for i, v := range values {
		elements[i] = types.StringValue(v)
	}

	return types.SetValueMust(types.StringType, elements)
}
