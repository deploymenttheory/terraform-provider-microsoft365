// planmodifiers/set.go
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
