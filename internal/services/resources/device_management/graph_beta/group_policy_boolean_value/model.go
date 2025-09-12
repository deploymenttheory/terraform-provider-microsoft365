// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicypresentationvalueboolean?view=graph-rest-beta
package graphBetaGroupPolicyBooleanValue

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyBooleanValueResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	GroupPolicyConfigurationID types.String `tfsdk:"group_policy_configuration_id"`

	// Input fields for auto-discovery
	PolicyName   types.String `tfsdk:"policy_name"`
	ClassType    types.String `tfsdk:"class_type"`
	CategoryPath types.String `tfsdk:"category_path"`

	// Computed/resolved fields (backward compatibility)
	GroupPolicyDefinitionValueID types.String `tfsdk:"group_policy_definition_value_id"`

	Enabled              types.Bool     `tfsdk:"enabled"`
	Values               types.List     `tfsdk:"values"` // List of BooleanPresentationValue
	CreatedDateTime      types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String   `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`

	// Internal data for passing instance IDs during updates and resolved presentations
	AdditionalData map[string]any `tfsdk:"-"`
}

type BooleanPresentationValue struct {
	PresentationID types.String `tfsdk:"presentation_id"`
	Value          types.Bool   `tfsdk:"value"`
}

type ResolvedPresentation struct {
	TemplateID string
	InstanceID string
	Index      int
}
