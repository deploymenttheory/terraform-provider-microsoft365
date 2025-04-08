// planmodifiers/string.go

package planmodifiers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
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

// caseInsensitiveString handles case-insensitive string comparisons
type caseInsensitiveString struct {
	stringModifier
}

func (m caseInsensitiveString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// For config values that don't match state, preserve their case
	if req.ConfigValue.IsNull() || req.StateValue.IsNull() {
		resp.PlanValue = req.ConfigValue
		return
	}

	if strings.EqualFold(req.PlanValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
		return
	}

	// Allow either case from config
	if strings.EqualFold(req.PlanValue.ValueString(), req.ConfigValue.ValueString()) {
		resp.PlanValue = req.ConfigValue
		return
	}

	resp.PlanValue = types.StringValue(strings.ToUpper(req.PlanValue.ValueString()))
}

// CaseInsensitiveString returns a plan modifier for case-insensitive string handling
func CaseInsensitiveString() StringModifier {
	return caseInsensitiveString{
		stringModifier: stringModifier{
			description:         "Handles case-insensitive string comparisons",
			markdownDescription: "Handles case-insensitive string comparisons",
		},
	}
}

// RequiresOtherAttributeEnabled returns a plan modifier that ensures an attribute
// can only be used when another specified attribute is enabled (set to true).
func RequiresOtherAttributeEnabled(dependencyPath path.Path) planmodifier.String {
	return &requiresOtherAttributeEnabledModifier{
		dependencyPath: dependencyPath,
	}
}

// requiresOtherAttributeEnabledModifier implements the plan modifier.
type requiresOtherAttributeEnabledModifier struct {
	dependencyPath path.Path
}

// Description returns a human-readable description of the plan modifier.
func (m *requiresOtherAttributeEnabledModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is enabled", m.dependencyPath)
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m *requiresOtherAttributeEnabledModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is enabled", m.dependencyPath)
}

// PlanModifyString implements the plan modification logic.
func (m *requiresOtherAttributeEnabledModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Skip if the attribute is null in the plan
	if req.PlanValue.IsNull() {
		return
	}

	// Get the dependency attribute's value from the plan
	var dependencyValue types.Bool
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If dependency is defined, not null, and false, this attribute should not be used
	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && !dependencyValue.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is enabled (true)", m.dependencyPath),
		)
	}
}
