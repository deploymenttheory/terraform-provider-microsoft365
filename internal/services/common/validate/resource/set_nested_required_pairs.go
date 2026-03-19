package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetNestedRequiredPairsValidator validates that within each element of a set-nested attribute,
// when FieldA holds a specific value, FieldB must hold the corresponding required value,
// and vice-versa. This enforces bidirectional pairing constraints at plan time.
type SetNestedRequiredPairsValidator struct {
	// ParentPath is a guard path — validation is skipped when the parent block is absent.
	ParentPath path.Path
	// SetPath is the full path to the set-nested attribute.
	SetPath path.Path
	// FieldA is the name of the first field in each set element object.
	FieldA string
	// FieldB is the name of the second field in each set element object.
	FieldB string
	// Pairs maps FieldA value → required FieldB value.
	// The reverse is also enforced: if FieldB equals a value in the map, FieldA must be the
	// corresponding key.
	Pairs map[string]string
}

// SetNestedRequiredPairs returns a resource.ConfigValidator that enforces bidirectional value
// pairing constraints within each element of a set-nested attribute.
//
// Example usage — ineligible signal requires offerFallback action (and vice-versa):
//
//	resourcevalidator.SetNestedRequiredPairs(
//	    path.Root("settings"),
//	    path.Root("settings").AtName("monitoring").AtName("monitoring_rules"),
//	    "signal",
//	    "action",
//	    map[string]string{
//	        "ineligible": "offerFallback",
//	    },
//	)
func SetNestedRequiredPairs(parentPath, setPath path.Path, fieldA, fieldB string, pairs map[string]string) resource.ConfigValidator {
	return &SetNestedRequiredPairsValidator{
		ParentPath: parentPath,
		SetPath:    setPath,
		FieldA:     fieldA,
		FieldB:     fieldB,
		Pairs:      pairs,
	}
}

// Description describes the validation in plain text formatting.
func (v SetNestedRequiredPairsValidator) Description(_ context.Context) string {
	descs := make([]string, 0, len(v.Pairs))
	for k, val := range v.Pairs {
		descs = append(descs, fmt.Sprintf("%s=%q requires %s=%q (and vice-versa)", v.FieldA, k, v.FieldB, val))
	}
	return fmt.Sprintf("%s: %s", v.SetPath, strings.Join(descs, "; "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v SetNestedRequiredPairsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateResource performs the validation.
func (v SetNestedRequiredPairsValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if req.Config.Raw.IsNull() {
		return
	}

	// Skip when the parent block is absent.
	var parentObj types.Object
	if diags := req.Config.GetAttribute(ctx, v.ParentPath, &parentObj); diags.HasError() || parentObj.IsNull() || parentObj.IsUnknown() {
		return
	}

	var rulesSet types.Set
	if diags := req.Config.GetAttribute(ctx, v.SetPath, &rulesSet); diags.HasError() || rulesSet.IsNull() || rulesSet.IsUnknown() {
		return
	}

	// Build a reverse map: FieldB value → required FieldA value.
	reverse := make(map[string]string, len(v.Pairs))
	for k, val := range v.Pairs {
		reverse[val] = k
	}

	for _, elem := range rulesSet.Elements() {
		obj, ok := elem.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}

		attrs := obj.Attributes()

		aVal, ok := attrs[v.FieldA].(types.String)
		if !ok || aVal.IsNull() || aVal.IsUnknown() {
			continue
		}

		bVal, ok := attrs[v.FieldB].(types.String)
		if !ok || bVal.IsNull() || bVal.IsUnknown() {
			continue
		}

		a := aVal.ValueString()
		b := bVal.ValueString()

		// Check forward: if FieldA is a constrained key, FieldB must be the required value.
		if required, constrained := v.Pairs[a]; constrained && b != required {
			resp.Diagnostics.AddAttributeError(
				v.SetPath,
				"Invalid field combination in set element",
				fmt.Sprintf(
					"When %s is %q, %s must be %q (got %q).",
					v.FieldA, a, v.FieldB, required, b,
				),
			)
		}

		// Check reverse: if FieldB is a constrained value, FieldA must be the required key.
		if requiredKey, constrained := reverse[b]; constrained && a != requiredKey {
			resp.Diagnostics.AddAttributeError(
				v.SetPath,
				"Invalid field combination in set element",
				fmt.Sprintf(
					"When %s is %q, %s must be %q (got %q).",
					v.FieldB, b, v.FieldA, requiredKey, a,
				),
			)
		}
	}
}
