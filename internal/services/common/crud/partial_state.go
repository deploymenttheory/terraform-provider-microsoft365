package crud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// SetPartialState records an allocated resource before subsequent API operations
// can fail. Preserve known configuration and identity values, but replace unknown
// planned values with typed nulls: Terraform cannot persist unknowns after apply.
// The caller must populate the resource identity before calling this function.
// A later successful Read remains responsible for filling remote computed values.
func SetPartialState(ctx context.Context, state *tfsdk.State, value any) diag.Diagnostics {
	snapshot := tfsdk.State{Schema: state.Schema}
	diagnostics := snapshot.Set(ctx, value)
	if diagnostics.HasError() {
		return diagnostics
	}
	known, err := tftypes.Transform(snapshot.Raw, func(_ *tftypes.AttributePath, value tftypes.Value) (tftypes.Value, error) {
		if !value.IsKnown() {
			return tftypes.NewValue(value.Type(), nil), nil
		}
		return value, nil
	})
	if err != nil {
		diagnostics.AddError("Unable to save partial resource state", fmt.Sprintf("Could not resolve unknown planned values: %s", err))
		return diagnostics
	}
	snapshot.Raw = known
	*state = snapshot
	return diagnostics
}
