// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicydefinitionvalue?view=graph-rest-beta
package graphBetaGroupPolicyDefinition

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupPolicyDefinitionResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	GroupPolicyConfigurationID types.String   `tfsdk:"group_policy_configuration_id"`
	PolicyName                 types.String   `tfsdk:"policy_name"`
	ClassType                  types.String   `tfsdk:"class_type"`
	CategoryPath               types.String   `tfsdk:"category_path"`
	Enabled                    types.Bool     `tfsdk:"enabled"`
	Values                     types.Set      `tfsdk:"values"`
	CreatedDateTime            types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime       types.String   `tfsdk:"last_modified_date_time"`
	Timeouts                   timeouts.Value `tfsdk:"timeouts"`

	// Internal data for passing resolved information
	AdditionalData map[string]any `tfsdk:"-"`
}

type PresentationValue struct {
	ID    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	Value types.String `tfsdk:"value"`
}

type ResolvedPresentation struct {
	TemplateID string
	InstanceID string
	Label      string
	Type       string // OData type
	Index      int
}
