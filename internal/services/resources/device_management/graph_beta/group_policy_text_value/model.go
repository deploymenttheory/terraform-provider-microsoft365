// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicypresentationvaluetext?view=graph-rest-beta
package graphBetaGroupPolicyTextValue

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyTextValueResourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	GroupPolicyConfigurationID   types.String   `tfsdk:"group_policy_configuration_id"`
	PolicyName                   types.String   `tfsdk:"policy_name"`
	ClassType                    types.String   `tfsdk:"class_type"`
	CategoryPath                 types.String   `tfsdk:"category_path"`
	GroupPolicyDefinitionValueID types.String   `tfsdk:"group_policy_definition_value_id"`
	PresentationID               types.String   `tfsdk:"presentation_id"`
	Enabled                      types.Bool     `tfsdk:"enabled"`
	Value                        types.String   `tfsdk:"value"`
	CreatedDateTime              types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime         types.String   `tfsdk:"last_modified_date_time"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
	// Internal data for passing instance IDs during updates
	AdditionalData map[string]any `tfsdk:"-"`
}

type ResolvedPresentation struct {
	TemplateID string
	InstanceID string
	Index      int
}
