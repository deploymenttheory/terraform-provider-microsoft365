package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PreserveComputedStringsWhenUnchanged prevents metadata-only updates after
// optional computed values have been restored by attribute plan modifiers.
// Metadata stays unknown when any other planned value changes, allowing Graph
// timestamps to change during a real write without an inconsistent result.
func PreserveComputedStringsWhenUnchanged(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, names ...string) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || !req.Plan.Raw.IsKnown() {
		return
	}
	var original, prior map[string]tftypes.Value
	if err := req.Plan.Raw.As(&original); err != nil {
		resp.Diagnostics.AddError("Read planned values", err.Error())
		return
	}
	if err := req.State.Raw.As(&prior); err != nil {
		resp.Diagnostics.AddError("Read prior values", err.Error())
		return
	}
	planned := make(map[string]tftypes.Value, len(original))
	for name, value := range original {
		planned[name] = value
	}
	for _, name := range names {
		value, ok := prior[name]
		if !ok {
			resp.Diagnostics.AddError("Unknown computed attribute", name)
			return
		}
		planned[name] = value
	}
	if !tftypes.NewValue(req.Plan.Raw.Type(), planned).Equal(req.State.Raw) {
		return
	}
	for _, name := range names {
		var value types.String
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(name), &value)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(name), value)...)
	}
}
