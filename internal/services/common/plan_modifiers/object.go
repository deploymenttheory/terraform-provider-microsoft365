// planmodifiers/object.go

package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ObjectModifier provides plan modifier functionality for objects.
type ObjectModifier interface {
	planmodifier.Object
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

// Common structure for object modifiers with descriptions
type objectModifier struct {
	description         string
	markdownDescription string
}

func (m objectModifier) Description(ctx context.Context) string {
	return m.description
}

func (m objectModifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

// UseStateForUnknownObject sets the plan value to the state value if the plan is unknown.
type useStateForUnknownObject struct {
	objectModifier
}

func (m useStateForUnknownObject) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.PlanValue.IsUnknown() && !req.StateValue.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

// UseStateForUnknownObject constructor
func UseStateForUnknownObject() ObjectModifier {
	return useStateForUnknownObject{
		objectModifier: objectModifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}

// DefaultValueObject sets a default value to an object if the plan value is null.
type defaultValueObject struct {
	objectModifier
	defaultValue types.Object
}

func (m defaultValueObject) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.defaultValue
	}
}

// DefaultValueObject constructor, creates an object modifier with a specified default value
func DefaultValueObject(defaultValue map[string]attr.Value) ObjectModifier {
	return defaultValueObject{
		objectModifier: objectModifier{
			description:         "Default value set to specified object",
			markdownDescription: "Default value set to specified object",
		},
		defaultValue: createDefaultObject(defaultValue),
	}
}

// Helper function to create a default empty object if needed
func createDefaultObject(defaultValue map[string]attr.Value) types.Object {
	return types.ObjectValueMust(map[string]attr.Type{}, defaultValue)
}

// RequiresReplaceIfStateNonNullObject requires replacement whenever the prior state value is
// non-null — i.e. once the attribute has been configured, any change (including removal)
// forces a new resource. When the state is null (first-time configuration), in-place updates
// are allowed.
type requiresReplaceIfStateNonNullObject struct {
	objectModifier
}

func (m requiresReplaceIfStateNonNullObject) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	resp.RequiresReplace = true
}

// RequiresReplaceIfStateNonNullObject returns a plan modifier that requires resource replacement
// when the state already has the attribute configured and the plan differs from the state.
// First-time configuration (state is null) allows in-place update.
func RequiresReplaceIfStateNonNullObject() planmodifier.Object {
	return requiresReplaceIfStateNonNullObject{
		objectModifier: objectModifier{
			description:         "Requires replacement if the attribute was previously configured and has changed.",
			markdownDescription: "Requires replacement if the attribute was previously configured and has changed.",
		},
	}
}

// ---- Object Attribute Implementation ----

// RequiresOtherAttributeEnabledObject returns a plan modifier that ensures an object attribute
// can only be used when another specified attribute is enabled (set to true).
func RequiresOtherAttributeEnabledObject(dependencyPath path.Path) planmodifier.Object {
	return &requiresOtherAttributeEnabledObjectModifier{
		dependencyPath: dependencyPath,
	}
}

type requiresOtherAttributeEnabledObjectModifier struct {
	dependencyPath path.Path
}

func (m *requiresOtherAttributeEnabledObjectModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledObjectModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
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
