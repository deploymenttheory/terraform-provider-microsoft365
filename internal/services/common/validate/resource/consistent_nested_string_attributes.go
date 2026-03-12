package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ConsistentNestedStringAttributesValidator validates that two string attributes share the same
// value whenever their common parent block is configured.
type ConsistentNestedStringAttributesValidator struct {
	ParentPath path.Path
	PathA      path.Path
	PathB      path.Path
}

// ConsistentNestedStringAttributes returns a resource validator that ensures two string attributes
// within a nested block always carry the same value.
//
// Example usage:
//
//	func (r *MyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
//	    return []resource.ConfigValidator{
//	        resource_level.ConsistentNestedStringAttributes(
//	            path.Root("tenant_restrictions"),
//	            path.Root("tenant_restrictions").AtName("users_and_groups").AtName("access_type"),
//	            path.Root("tenant_restrictions").AtName("applications").AtName("access_type"),
//	        ),
//	    }
//	}
func ConsistentNestedStringAttributes(parentPath, pathA, pathB path.Path) resource.ConfigValidator {
	return &ConsistentNestedStringAttributesValidator{
		ParentPath: parentPath,
		PathA:      pathA,
		PathB:      pathB,
	}
}

// Description describes the validation in plain text formatting.
func (v ConsistentNestedStringAttributesValidator) Description(_ context.Context) string {
	return fmt.Sprintf("%s and %s must have the same value", v.PathA, v.PathB)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v ConsistentNestedStringAttributesValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("`%s` and `%s` must have the same value", v.PathA, v.PathB)
}

// ValidateResource performs the validation.
func (v ConsistentNestedStringAttributesValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if req.Config.Raw.IsNull() {
		return
	}

	// Skip when the parent block is absent.
	var parentObj basetypes.ObjectValue
	if diags := req.Config.GetAttribute(ctx, v.ParentPath, &parentObj); diags.HasError() || parentObj.IsNull() || parentObj.IsUnknown() {
		return
	}

	var valA, valB basetypes.StringValue
	if diags := req.Config.GetAttribute(ctx, v.PathA, &valA); diags.HasError() {
		return
	}
	if diags := req.Config.GetAttribute(ctx, v.PathB, &valB); diags.HasError() {
		return
	}

	if valA.IsNull() || valA.IsUnknown() || valB.IsNull() || valB.IsUnknown() {
		return
	}

	if valA.ValueString() != valB.ValueString() {
		msg := fmt.Sprintf(
			"Both attributes must have the same value: got %q and %q. They must both be \"allowed\" or both be \"blocked\".",
			valA.ValueString(), valB.ValueString(),
		)
		resp.Diagnostics.AddAttributeError(v.PathA, "Inconsistent access_type values", msg)
		resp.Diagnostics.AddAttributeError(v.PathB, "Inconsistent access_type values", msg)
	}
}
