package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NoAllUsersMixValidator validates that within a targets set, the special value "AllUsers"
// (target_type = "user", target = "AllUsers") is not mixed with specific user or group GUIDs.
type NoAllUsersMixValidator struct {
	// ParentPath is the top-level block path used as a guard — validation is skipped when absent.
	ParentPath path.Path
	// TargetsPath is the full path to the targets set attribute.
	TargetsPath path.Path
}

// NoAllUsersMix returns a resource.ConfigValidator that ensures "AllUsers" is not combined
// with specific user or group GUIDs in the same targets set.
//
// Example usage:
//
//	resourcevalidator.NoAllUsersMix(
//	    path.Root("b2b_collaboration_inbound"),
//	    path.Root("b2b_collaboration_inbound").AtName("users_and_groups").AtName("targets"),
//	)
func NoAllUsersMix(parentPath, targetsPath path.Path) resource.ConfigValidator {
	return &NoAllUsersMixValidator{
		ParentPath:  parentPath,
		TargetsPath: targetsPath,
	}
}

// Description describes the validation in plain text formatting.
func (v NoAllUsersMixValidator) Description(_ context.Context) string {
	return fmt.Sprintf("%s: 'AllUsers' cannot be combined with specific user or group GUIDs in the same targets set", v.TargetsPath)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v NoAllUsersMixValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("`%s`: `AllUsers` cannot be combined with specific user or group GUIDs in the same `targets` set", v.TargetsPath)
}

// ValidateResource performs the validation.
func (v NoAllUsersMixValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if req.Config.Raw.IsNull() {
		return
	}

	// Skip when the parent block is absent.
	var parentObj types.Object
	if diags := req.Config.GetAttribute(ctx, v.ParentPath, &parentObj); diags.HasError() || parentObj.IsNull() || parentObj.IsUnknown() {
		return
	}

	var targetsSet types.Set
	if diags := req.Config.GetAttribute(ctx, v.TargetsPath, &targetsSet); diags.HasError() || targetsSet.IsNull() || targetsSet.IsUnknown() {
		return
	}

	hasAllUsers := false
	hasSpecificUserOrGroup := false

	for _, elem := range targetsSet.Elements() {
		targetObj, ok := elem.(types.Object)
		if !ok || targetObj.IsNull() || targetObj.IsUnknown() {
			continue
		}

		attrs := targetObj.Attributes()

		targetVal, ok := attrs["target"].(types.String)
		if !ok || targetVal.IsNull() || targetVal.IsUnknown() {
			continue
		}

		targetTypeVal, ok := attrs["target_type"].(types.String)
		if !ok || targetTypeVal.IsNull() || targetTypeVal.IsUnknown() {
			continue
		}

		targetType := targetTypeVal.ValueString()

		if targetVal.ValueString() == "AllUsers" && targetType == "user" {
			hasAllUsers = true
		} else if targetType == "user" || targetType == "group" {
			hasSpecificUserOrGroup = true
		}
	}

	if hasAllUsers && hasSpecificUserOrGroup {
		resp.Diagnostics.AddAttributeError(
			v.TargetsPath,
			"Invalid targets combination",
			"'AllUsers' cannot be combined with specific user or group GUIDs in the same targets set. "+
				"Use either 'AllUsers' to target all users and groups, or list specific user/group GUIDs, but not both.",
		)
	}
}
