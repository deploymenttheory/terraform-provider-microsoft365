// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsAutopilotDevicePreparationPolicyDataSourceModel represents the Terraform data source model
type WindowsAutopilotDevicePreparationPolicyDataSourceModel struct {
	ID              types.String                                       `tfsdk:"id"`
	PolicyId        types.String                                       `tfsdk:"policy_id"`
	Name            types.String                                       `tfsdk:"name"`
	ListAll         types.Bool                                         `tfsdk:"list_all"`
	ListAssignments types.Bool                                         `tfsdk:"list_assignments"`
	ODataQuery      types.String                                       `tfsdk:"odata_query"`
	Items           []WindowsAutopilotDevicePreparationPolicyItemModel `tfsdk:"items"`
	Assignments     []PolicyAssignmentModel                            `tfsdk:"assignments"`
	Timeouts        timeouts.Value                                     `tfsdk:"timeouts"`
}

// WindowsAutopilotDevicePreparationPolicyItemModel represents an individual device preparation policy
type WindowsAutopilotDevicePreparationPolicyItemModel struct {
	ID                                types.String            `tfsdk:"id"`
	Name                              types.String            `tfsdk:"name"`
	Description                       types.String            `tfsdk:"description"`
	CreatedDateTime                   types.String            `tfsdk:"created_date_time"`
	LastModifiedDateTime              types.String            `tfsdk:"last_modified_date_time"`
	CreationSource                    types.String            `tfsdk:"creation_source"`
	Platforms                         types.String            `tfsdk:"platforms"`
	Technologies                      types.String            `tfsdk:"technologies"`
	SettingCount                      types.Int64             `tfsdk:"setting_count"`
	IsAssigned                        types.Bool              `tfsdk:"is_assigned"`
	DisableEntraGroupPolicyAssignment types.Bool              `tfsdk:"disable_entra_group_policy_assignment"`
	Priority                          types.Int64             `tfsdk:"priority"`
	RoleScopeTagIds                   []types.String          `tfsdk:"role_scope_tag_ids"`
	TemplateReference                 *TemplateReferenceModel `tfsdk:"template_reference"`
}

// TemplateReferenceModel represents the settings catalog template a policy was created from
type TemplateReferenceModel struct {
	TemplateId             types.String `tfsdk:"template_id"`
	TemplateFamily         types.String `tfsdk:"template_family"`
	TemplateDisplayName    types.String `tfsdk:"template_display_name"`
	TemplateDisplayVersion types.String `tfsdk:"template_display_version"`
	DeploymentMode         types.String `tfsdk:"deployment_mode"`
}

// PolicyAssignmentModel represents a single assignment of a device preparation policy
type PolicyAssignmentModel struct {
	ID         types.String `tfsdk:"id"`
	Type       types.String `tfsdk:"type"`
	GroupId    types.String `tfsdk:"group_id"`
	FilterId   types.String `tfsdk:"filter_id"`
	FilterType types.String `tfsdk:"filter_type"`
}
