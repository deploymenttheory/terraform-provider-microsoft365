package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	errExecutionConditionsRequired = errors.New("execution_conditions is required")
	errOnDemandWithScopeOrTrigger  = errors.New("scope and trigger must not be set when on_demand_execution_only is true")
	errScopeAndTriggerRequired     = errors.New("scope and trigger are both required unless on_demand_execution_only is true")
	errMembershipChangeNeedsGroups = errors.New("a membership_change trigger requires scope.group_ids")
	errScopeUnset                  = errors.New("execution_conditions.scope must set rule or group_ids")
	errTriggerUnset                = errors.New("execution_conditions.trigger must set time_based_attribute, attribute_change or membership_change")
	errGroupIdsNotStrings          = errors.New("execution_conditions.scope.group_ids must contain strings")
)

// ValidateConfig checks the execution conditions at plan time so an invalid combination fails
// before any Graph call is made.
func (r *LifecycleWorkflowResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var conditions types.Object
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("execution_conditions"), &conditions)...)
	if resp.Diagnostics.HasError() || conditions.IsNull() || conditions.IsUnknown() {
		return
	}

	// Nested values can still be unknown (for example group IDs from another resource); those
	// configurations are validated again when the request is built.
	var data *WorkflowExecutionConditions
	if diags := req.Config.GetAttribute(ctx, path.Root("execution_conditions"), &data); diags.HasError() || data == nil {
		return
	}

	if err := validateExecutionConditions(data); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("execution_conditions"), "Invalid execution conditions", err.Error())
	}
}

// validateRequest validates the model before it is sent to Graph.
func validateRequest(data *LifecycleWorkflowResourceModel) error {
	if data.ExecutionConditions == nil {
		return errExecutionConditionsRequired
	}
	return validateExecutionConditions(data.ExecutionConditions)
}

// validateExecutionConditions enforces the onDemandExecutionOnly / triggerAndScopeBasedConditions split.
// Values that are still unknown are skipped; they are checked again when the request is built.
func validateExecutionConditions(data *WorkflowExecutionConditions) error {
	if data.OnDemandExecutionOnly.IsUnknown() {
		return nil
	}

	if data.OnDemandExecutionOnly.ValueBool() {
		if data.Scope != nil || data.Trigger != nil {
			return errOnDemandWithScopeOrTrigger
		}
		return nil
	}

	if data.Scope == nil || data.Trigger == nil {
		return errScopeAndTriggerRequired
	}

	if data.Trigger.MembershipChange != nil && data.Scope.GroupIds.IsNull() {
		return errMembershipChangeNeedsGroups
	}

	return nil
}
