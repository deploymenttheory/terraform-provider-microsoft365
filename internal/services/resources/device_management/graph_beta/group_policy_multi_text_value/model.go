// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicypresentationvaluemultitext?view=graph-rest-beta
package graphBetaGroupPolicyMultiTextValue

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyMultiTextValueResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	GroupPolicyConfigurationID types.String `tfsdk:"group_policy_configuration_id"`
	// input fields for auto-discovery
	PolicyName                   types.String   `tfsdk:"policy_name"`
	ClassType                    types.String   `tfsdk:"class_type"`
	PresentationIndex            types.Int64    `tfsdk:"presentation_index"`
	GroupPolicyDefinitionValueID types.String   `tfsdk:"group_policy_definition_value_id"`
	PresentationID               types.String   `tfsdk:"presentation_id"`
	Values                       types.List     `tfsdk:"values"`
	CreatedDateTime              types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime         types.String   `tfsdk:"last_modified_date_time"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}
