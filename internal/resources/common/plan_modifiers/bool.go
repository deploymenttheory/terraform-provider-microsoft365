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
	// Skip if the attribute is null in the plan
	if req.PlanValue.IsNull() {
		return
	}

	// Get the dependency attribute's value from the plan
	var dependencyValue types.String
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If dependency is defined, not null, and not the required value, this attribute should not be used
	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && dependencyValue.ValueString() != m.requiredValue {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is set to %q", m.dependencyPath, m.requiredValue),
		)
	}
}
