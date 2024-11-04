// planmodifiers/object.go

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
