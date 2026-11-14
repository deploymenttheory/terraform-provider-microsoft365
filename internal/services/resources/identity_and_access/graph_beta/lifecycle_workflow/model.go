// REF: https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflow?view=graph-rest-beta
package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LifecycleWorkflowResourceModel represents the schema for the lifecycle workflow resource
type LifecycleWorkflowResourceModel struct {
	ID                   types.String                 `tfsdk:"id"`
	Category             types.String                 `tfsdk:"category"`
	DisplayName          types.String                 `tfsdk:"display_name"`
	Description          types.String                 `tfsdk:"description"`
	IsEnabled            types.Bool                   `tfsdk:"is_enabled"`
	IsSchedulingEnabled  types.Bool                   `tfsdk:"is_scheduling_enabled"`
	ExecutionConditions  *WorkflowExecutionConditions `tfsdk:"execution_conditions"`
	Tasks                []WorkflowTask               `tfsdk:"tasks"`
	Version              types.Int64                  `tfsdk:"version"`
	CreatedDateTime      types.String                 `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String                 `tfsdk:"last_modified_date_time"`
	Timeouts             timeouts.Value               `tfsdk:"timeouts"`
}

// WorkflowExecutionConditions represents the polymorphic executionConditions of a workflow.
// Either on_demand_execution_only is true (onDemandExecutionOnly), or scope and trigger are set
// (triggerAndScopeBasedConditions).
type WorkflowExecutionConditions struct {
	OnDemandExecutionOnly types.Bool       `tfsdk:"on_demand_execution_only"`
	Scope                 *WorkflowScope   `tfsdk:"scope"`
	Trigger               *WorkflowTrigger `tfsdk:"trigger"`
}

// WorkflowScope represents the subject set the workflow runs against.
// rule maps to ruleBasedSubjectSet and group_ids maps to groupBasedSubjectSet.
type WorkflowScope struct {
	Rule     types.String `tfsdk:"rule"`
	GroupIds types.Set    `tfsdk:"group_ids"`
}

// WorkflowTrigger represents the polymorphic workflowExecutionTrigger. Exactly one block is set.
type WorkflowTrigger struct {
	TimeBasedAttribute *TimeBasedAttributeTrigger `tfsdk:"time_based_attribute"`
	AttributeChange    *AttributeChangeTrigger    `tfsdk:"attribute_change"`
	MembershipChange   *MembershipChangeTrigger   `tfsdk:"membership_change"`
}

// TimeBasedAttributeTrigger represents a timeBasedAttributeTrigger
type TimeBasedAttributeTrigger struct {
	Attribute    types.String `tfsdk:"attribute"`
	OffsetInDays types.Int64  `tfsdk:"offset_in_days"`
}

// AttributeChangeTrigger represents an attributeChangeTrigger. Graph accepts a single trigger attribute.
type AttributeChangeTrigger struct {
	TriggerAttribute types.String `tfsdk:"trigger_attribute"`
}

// MembershipChangeTrigger represents a membershipChangeTrigger
type MembershipChangeTrigger struct {
	ChangeType types.String `tfsdk:"change_type"`
}

// WorkflowTask represents a built-in task in the workflow
type WorkflowTask struct {
	TaskDefinitionId  types.String           `tfsdk:"task_definition_id"`
	DisplayName       types.String           `tfsdk:"display_name"`
	Description       types.String           `tfsdk:"description"`
	IsEnabled         types.Bool             `tfsdk:"is_enabled"`
	ContinueOnError   types.Bool             `tfsdk:"continue_on_error"`
	ExecutionSequence types.Int64            `tfsdk:"execution_sequence"`
	Arguments         []WorkflowTaskArgument `tfsdk:"arguments"`
}

// WorkflowTaskArgument represents a task argument key/value pair
type WorkflowTaskArgument struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
