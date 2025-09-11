// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicypresentationvaluetext?view=graph-rest-beta
package graphBetaGroupPolicyTextValue

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyTextValueResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	GroupPolicyConfigurationID types.String `tfsdk:"group_policy_configuration_id"`

	// Input fields for auto-discovery
	PolicyName        types.String `tfsdk:"policy_name"`
	ClassType         types.String `tfsdk:"class_type"`
	CategoryPath      types.String `tfsdk:"category_path"`
	PresentationIndex types.Int64  `tfsdk:"presentation_index"`

	// Computed/resolved fields (backward compatibility)
	GroupPolicyDefinitionValueID types.String `tfsdk:"group_policy_definition_value_id"`
	PresentationID               types.String `tfsdk:"presentation_id"`

	Enabled              types.Bool     `tfsdk:"enabled"`
	Value                types.String   `tfsdk:"value"`
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`

	// Internal data for passing instance IDs during updates
	AdditionalData map[string]interface{} `tfsdk:"-"`
}
