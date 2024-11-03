// planmodifiers/bool.go
package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// BoolModifier defines the interface for boolean plan modifiers
type BoolModifier interface {
	planmodifier.Bool
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

type boolModifier struct {
	description         string
	markdownDescription string
}

func (m boolModifier) Description(ctx context.Context) string {
	return m.description
}

func (m boolModifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

type useStateForUnknownBool struct {
	boolModifier
}

func (m useStateForUnknownBool) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}

// UseStateForUnknownBool returns a plan modifier that copies a known prior state bool
// value into a planned unknown value.
func UseStateForUnknownBool() BoolModifier {
	return useStateForUnknownBool{
		boolModifier: boolModifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}
