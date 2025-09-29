// planmodifiers/set.go
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// defaultValueSet is a Set plan modifier that sets a default value when the config is null or empty

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

// defaultValueSet is a Set plan modifier that sets a default value when the config is null or empty
type defaultValueSet struct {
	setModifier
	defaultValue types.Set
}

// PlanModifySet sets the plan value to the default set if the config is null or empty.
func (m defaultValueSet) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsNull() && len(req.PlanValue.Elements()) > 0 {
		return
	}
	resp.PlanValue = m.defaultValue
}

// DefaultSetValue returns a SetModifier that sets the default value to the specified set.
func DefaultSetValue(defaultValue []attr.Value) planmodifier.Set {
	return defaultValueSet{
		setModifier: setModifier{
			description:         fmt.Sprintf("Default value set to %v", defaultValue),
			markdownDescription: fmt.Sprintf("Default value set to `%v`", defaultValue),
		},
		defaultValue: types.SetValueMust(types.StringType, defaultValue),
	}
}

// DefaultSetEmptyValue returns a SetModifier that sets the default value to an empty set.
func DefaultSetEmptyValue() planmodifier.Set {
	emptySet, _ := types.SetValue(types.StringType, []attr.Value{})
	return defaultValueSet{
		setModifier: setModifier{
			description:         "Default value set to empty set",
			markdownDescription: "Default value set to empty set",
		},
		defaultValue: emptySet,
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
	if req.PlanValue.IsNull() {
		return
	}

	var dependencyValue types.Bool
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && !dependencyValue.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is enabled (true)", m.dependencyPath),
		)
	}
}

// UseStateForUnknownOrNullSet returns a plan modifier that copies a known prior state
// Set value into the planned value if the planned value is null or unknown.
// This is useful for fields that are populated during creation but may not be
// explicitly set in configuration.
func UseStateForUnknownOrNullSet() planmodifier.Set {
	return useStateForUnknownOrNullSetModifier{}
}

// useStateForUnknownOrNullSetModifier implements the modifier
type useStateForUnknownOrNullSetModifier struct{}

// Description returns a plain text description of the modifier's behavior.
func (m useStateForUnknownOrNullSetModifier) Description(ctx context.Context) string {
	return "If the Set is unknown or null after plan creation, use the value from the state."
}

// MarkdownDescription returns a markdown formatted description of the modifier's behavior.
func (m useStateForUnknownOrNullSetModifier) MarkdownDescription(ctx context.Context) string {
	return "If the Set is unknown or null after plan creation, use the value from the state."
}

// PlanModifySet implements the plan modification logic.
func (m useStateForUnknownOrNullSetModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	tflog.Debug(ctx, "Using state value instead of null/unknown plan value for set", map[string]any{
		"path": req.Path.String(),
	})

	resp.PlanValue = req.StateValue
}
