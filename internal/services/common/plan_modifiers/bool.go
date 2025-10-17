// planmodifiers/bool.go
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// Add new default value modifier
type defaultValueBool struct {
	boolModifier
	defaultValue bool
}

func (m defaultValueBool) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = types.BoolValue(m.defaultValue)
}

// BoolDefaultValue returns a plan modifier that sets a default value if the planned value is null
func BoolDefaultValue(defaultValue bool) BoolModifier {
	return defaultValueBool{
		boolModifier: boolModifier{
			description:         "Set default value if null",
			markdownDescription: "Set default value if null",
		},
		defaultValue: defaultValue,
	}
}

// RequiresOtherAttributeValueBool returns a plan modifier that ensures a Bool attribute
// can only be used when another specified attribute has a specific string value.
func RequiresOtherAttributeValueBool(dependencyPath path.Path, requiredValue string) planmodifier.Bool {
	return &requiresOtherAttributeValueBoolModifier{
		dependencyPath: dependencyPath,
		requiredValue:  requiredValue,
	}
}

type requiresOtherAttributeValueBoolModifier struct {
	dependencyPath path.Path
	requiredValue  string
}

func (m *requiresOtherAttributeValueBoolModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is set to %q", m.dependencyPath, m.requiredValue)
}

func (m *requiresOtherAttributeValueBoolModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is set to `%s`", m.dependencyPath, m.requiredValue)
}

func (m *requiresOtherAttributeValueBoolModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.PlanValue.IsNull() {
		return
	}

	var dependencyValue types.String
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && dependencyValue.ValueString() != m.requiredValue {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is set to %q", m.dependencyPath, m.requiredValue),
		)
	}
}

// RequiresReplaceIfChangedBool is a custom PlanModifier for BoolAttribute
type RequiresReplaceIfChangedBool struct{}

func (m RequiresReplaceIfChangedBool) Description(_ context.Context) string {
	return "Requires resource replacement if the boolean value changes."
}

func (m RequiresReplaceIfChangedBool) MarkdownDescription(_ context.Context) string {
	return "Requires resource replacement if the boolean value changes."
}

func (m RequiresReplaceIfChangedBool) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	stateVal := req.StateValue.ValueBool()
	planVal := req.PlanValue.ValueBool()

	if stateVal != planVal {
		resp.RequiresReplace = true
	}
}

// NewRequiresReplaceIfChangedBool returns a new instance of the RequiresReplaceIfChangedBool plan modifier.
func NewRequiresReplaceIfChangedBool() planmodifier.Bool {
	return RequiresReplaceIfChangedBool{}
}

// RequiresReplaceIfFalseToTrue is a custom PlanModifier that prevents changing from false to true
type RequiresReplaceIfFalseToTrue struct{}

func (m RequiresReplaceIfFalseToTrue) Description(_ context.Context) string {
	return "Requires resource replacement if the boolean value changes from false to true. Changes from true to false are allowed."
}

func (m RequiresReplaceIfFalseToTrue) MarkdownDescription(_ context.Context) string {
	return "Requires resource replacement if the boolean value changes from false to true. Changes from true to false are allowed."
}

func (m RequiresReplaceIfFalseToTrue) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	stateVal := req.StateValue.ValueBool()
	planVal := req.PlanValue.ValueBool()

	if !stateVal && planVal {
		resp.RequiresReplace = true
	}
}

// NewRequiresReplaceIfFalseToTrue returns a new instance of the RequiresReplaceIfFalseToTrue plan modifier.
func NewRequiresReplaceIfFalseToTrue() planmodifier.Bool {
	return RequiresReplaceIfFalseToTrue{}
}
