// planmodifiers/list.go

package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ListModifier interface {
	planmodifier.List
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

type listModifier struct {
	description         string
	markdownDescription string
}

func (m listModifier) Description(ctx context.Context) string {
	return m.description
}

func (m listModifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

type defaultValueList struct {
	listModifier
	defaultValue types.List
}

func (m defaultValueList) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if !req.PlanValue.IsNull() {
		return
	}
	resp.PlanValue = m.defaultValue
}

// DefaultListValue sets a custom default list value
func DefaultListValue(defaultValue []attr.Value) ListModifier {
	return defaultValueList{
		listModifier: listModifier{
			description:         fmt.Sprintf("Default value set to %v", defaultValue),
			markdownDescription: fmt.Sprintf("Default value set to `%v`", defaultValue),
		},
		defaultValue: types.ListValueMust(types.StringType, defaultValue),
	}
}

// DefaultListEmptyValue sets the default value to an empty list
func DefaultListEmptyValue() ListModifier {
	emptyList, _ := types.ListValue(types.StringType, []attr.Value{})
	return defaultValueList{
		listModifier: listModifier{
			description:         "Default value set to empty list",
			markdownDescription: "Default value set to empty list",
		},
		defaultValue: emptyList,
	}
}