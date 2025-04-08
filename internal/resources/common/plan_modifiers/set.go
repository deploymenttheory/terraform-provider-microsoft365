// planmodifiers/set.go
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SetModifier interface {
	planmodifier.Set
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

type setModifier struct {
	description         string
	markdownDescription string
}

func (m setModifier) Description(ctx context.Context) string {
	return m.description
}

func (m setModifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

type useStateForUnknownSet struct {
	setModifier
}

func (m useStateForUnknownSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}

func UseStateForUnknownSet() SetModifier {
	return useStateForUnknownSet{
		setModifier: setModifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}

type defaultValueSet struct {
	setModifier
	defaultValue types.Set
}

func (m defaultValueSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = m.defaultValue
}

func DefaultValueSet(defaultValue []attr.Value) SetModifier {
	return defaultValueSet{
		setModifier: setModifier{
			description:         fmt.Sprintf("Default value set to %v", defaultValue),
			markdownDescription: fmt.Sprintf("Default value set to `%v`", defaultValue),
		},
		defaultValue: types.SetValueMust(types.StringType, defaultValue),
	}
}

// RequiresOtherAttributeEnabledSet returns a plan modifier that ensures a set attribute
// can only be used when another specified attribute is enabled (set to true).
func RequiresOtherAttributeEnabledSet(dependencyPath path.Path) planmodifier.Set {
	return &requiresOtherAttributeEnabledSetModifier{
		dependencyPath: dependencyPath,
	}
}

type requiresOtherAttributeEnabledSetModifier struct {
	dependencyPath path.Path
}

func (m *requiresOtherAttributeEnabledSetModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledSetModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledSetModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
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
