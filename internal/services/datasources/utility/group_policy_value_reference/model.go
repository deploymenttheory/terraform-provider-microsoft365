package utilityGroupPolicyValueReference

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupPolicyValueReferenceDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	PolicyName  types.String   `tfsdk:"policy_name"`
	Definitions types.List     `tfsdk:"definitions"` // List of DefinitionModel
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

type DefinitionModel struct {
	Id            types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	ClassType     types.String `tfsdk:"class_type"`
	CategoryPath  types.String `tfsdk:"category_path"`
	ExplainText   types.String `tfsdk:"explain_text"`
	SupportedOn   types.String `tfsdk:"supported_on"`
	PolicyType    types.String `tfsdk:"policy_type"`
	Presentations types.List   `tfsdk:"presentations"` // List of PresentationModel
}

type PresentationModel struct {
	Id       types.String `tfsdk:"id"`
	Label    types.String `tfsdk:"label"`
	Type     types.String `tfsdk:"type"`
	Required types.Bool   `tfsdk:"required"`
}

