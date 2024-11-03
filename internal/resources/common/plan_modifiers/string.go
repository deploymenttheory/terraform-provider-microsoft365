// planmodifiers/string.go

package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringModifier defines the interface for string plan modifiers
type StringModifier interface {
	planmodifier.String
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

// stringModifier implements the base functionality for string plan modifiers
type stringModifier struct {
	description         string
	markdownDescription string
}

func (m stringModifier) Description(ctx context.Context) string {
	return m.description
}

func (m stringModifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

// UseStateForUnknown returns a plan modifier that copies a known prior state string value
// into a planned unknown value
type useStateForUnknownString struct {
	stringModifier
}

func (m useStateForUnknownString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is unknown and there's a state value, use the state value
	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}

// UseStateForUnknownString returns a plan modifier that copies a known prior state string
// value into a planned unknown value.
func UseStateForUnknownString() StringModifier {
	return useStateForUnknownString{
		stringModifier: stringModifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}

// RequiresReplaceString returns a plan modifier that requires resource replacement if
// the value changes.
type requiresReplaceString struct {
	stringModifier
}

func (m requiresReplaceString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.RequiresReplace = true
}

// RequiresReplaceIfString returns a plan modifier that requires resource replacement if
// the string value changes.
func RequiresReplaceString() StringModifier {
	return requiresReplaceString{
		stringModifier: stringModifier{
			description:         "Requires resource replacement if value changes",
			markdownDescription: "Requires resource replacement if value changes",
		},
	}
}

// DefaultValueString returns a plan modifier that sets a default value if the planned
// value is null.
type defaultValueString struct {
	stringModifier
	defaultValue types.String
}

func (m defaultValueString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = m.defaultValue
}

// DefaultValue returns a plan modifier that sets a default string value.
func DefaultValueString(defaultValue string) StringModifier {
	return defaultValueString{
		stringModifier: stringModifier{
			description:         fmt.Sprintf("Default value set to %q", defaultValue),
			markdownDescription: fmt.Sprintf("Default value set to `%s`", defaultValue),
		},
		defaultValue: types.StringValue(defaultValue),
	}
}
