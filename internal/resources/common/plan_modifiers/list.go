// planmodifiers/list.go

package planmodifiers

import (
	"context"

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

func ListDefaultValueEmpty() ListModifier {
	emptyList, _ := types.ListValue(types.ObjectType{AttrTypes: map[string]attr.Type{}}, []attr.Value{})
	return defaultValueList{
		listModifier: listModifier{
			description:         "Default value set to empty list",
			markdownDescription: "Default value set to empty list",
		},
		defaultValue: emptyList,
	}
}
