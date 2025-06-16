package state

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// BuildObjectSetFromSlice creates a Set from a slice using an extractor function.
//
// This function converts a Go slice into a Terraform Framework Set type with proper
// error handling and logging. It's commonly used for mapping Graph API collections
// to Terraform state.
//
// Parameters:
//   - ctx: Context for logging
//   - attrTypes: Map defining the attribute types for objects in the set
//   - extract: Function that extracts attribute values for each slice element
//   - length: Length of the source slice
//
// Returns:
//   - types.Set: The constructed set, or a null set if errors occur
func BuildObjectSetFromSlice(
	ctx context.Context,
	attrTypes map[string]attr.Type,
	extract func(i int) map[string]attr.Value,
	length int,
) types.Set {
	objectType := types.ObjectType{AttrTypes: attrTypes}

	if length == 0 {
		emptySet, diags := types.SetValue(objectType, []attr.Value{})
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create empty set", map[string]interface{}{
				"errors": diags.Errors(),
			})
			return types.SetNull(objectType)
		}
		return emptySet
	}

	var elements []attr.Value
	for i := 0; i < length; i++ {
		values := extract(i)
		obj, diags := types.ObjectValue(attrTypes, values)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to build object for Set", map[string]interface{}{
				"index":  i,
				"errors": diags.Errors(),
			})
			continue
		}
		elements = append(elements, obj)
	}

	set, diags := types.SetValue(objectType, elements)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to build Set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(objectType)
	}

	return set
}
