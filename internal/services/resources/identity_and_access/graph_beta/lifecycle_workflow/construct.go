package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	igmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK workflow.
// The result is used both for POST /workflows and as the workflow in a createNewVersion request,
// so it always carries every writable property: Graph requires category on createNewVersion.
func constructResource(ctx context.Context, data *LifecycleWorkflowResourceModel) (igmodels.Workflowable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := igmodels.NewWorkflow()

	if err := convert.FrameworkToGraphEnum(data.Category, igmodels.ParseLifecycleWorkflowCategory, requestBody.SetCategory); err != nil {
		return nil, fmt.Errorf("failed to set category: %w", err)
	}

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsEnabled, requestBody.SetIsEnabled)
	convert.FrameworkToGraphBool(data.IsSchedulingEnabled, requestBody.SetIsSchedulingEnabled)

	executionConditions, err := constructExecutionConditions(data.ExecutionConditions)
	if err != nil {
		return nil, fmt.Errorf("failed to construct execution conditions: %w", err)
	}
	requestBody.SetExecutionConditions(executionConditions)

	requestBody.SetTasks(constructTasks(data.Tasks))

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructBasicUpdate builds the PATCH body. Graph only accepts displayName, description,
// isEnabled and isSchedulingEnabled on PATCH and rejects any other property with
// "Only displayname description isenabled and isschedulingenabled can be patched".
func constructBasicUpdate(ctx context.Context, data *LifecycleWorkflowResourceModel) (igmodels.Workflowable, error) {
	requestBody := igmodels.NewWorkflow()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsEnabled, requestBody.SetIsEnabled)
	convert.FrameworkToGraphBool(data.IsSchedulingEnabled, requestBody.SetIsSchedulingEnabled)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final PATCH JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructExecutionConditions builds either onDemandExecutionOnly or triggerAndScopeBasedConditions.
func constructExecutionConditions(data *WorkflowExecutionConditions) (igmodels.WorkflowExecutionConditionsable, error) {
	if data == nil {
		return nil, errExecutionConditionsRequired
	}

	if data.OnDemandExecutionOnly.ValueBool() {
		return igmodels.NewOnDemandExecutionOnly(), nil
	}

	conditions := igmodels.NewTriggerAndScopeBasedConditions()

	scope, err := constructScope(data.Scope)
	if err != nil {
		return nil, err
	}
	conditions.SetScope(scope)

	trigger, err := constructTrigger(data.Trigger)
	if err != nil {
		return nil, err
	}
	conditions.SetTrigger(trigger)

	return conditions, nil
}

// constructScope builds a ruleBasedSubjectSet or a groupBasedSubjectSet.
func constructScope(data *WorkflowScope) (graphmodels.SubjectSetable, error) {
	if data == nil {
		return nil, errScopeAndTriggerRequired
	}

	if !data.Rule.IsNull() && !data.Rule.IsUnknown() {
		scope := igmodels.NewRuleBasedSubjectSet()
		convert.FrameworkToGraphString(data.Rule, scope.SetRule)
		return scope, nil
	}

	if !data.GroupIds.IsNull() && !data.GroupIds.IsUnknown() {
		groups := make([]graphmodels.Groupable, 0, len(data.GroupIds.Elements()))
		for _, element := range data.GroupIds.Elements() {
			value, ok := element.(interface{ ValueString() string })
			if !ok {
				return nil, errGroupIdsNotStrings
			}
			id := value.ValueString()
			group := graphmodels.NewGroup()
			group.SetId(&id)
			groups = append(groups, group)
		}
		scope := igmodels.NewGroupBasedSubjectSet()
		scope.SetGroups(groups)
		return scope, nil
	}

	return nil, errScopeUnset
}

// constructTrigger builds the polymorphic workflowExecutionTrigger.
func constructTrigger(data *WorkflowTrigger) (igmodels.WorkflowExecutionTriggerable, error) {
	if data == nil {
		return nil, errScopeAndTriggerRequired
	}

	switch {
	case data.TimeBasedAttribute != nil:
		trigger := igmodels.NewTimeBasedAttributeTrigger()
		if err := convert.FrameworkToGraphEnum(data.TimeBasedAttribute.Attribute, igmodels.ParseWorkflowTriggerTimeBasedAttribute, trigger.SetTimeBasedAttribute); err != nil {
			return nil, fmt.Errorf("failed to set time based attribute: %w", err)
		}
		offset := int32(data.TimeBasedAttribute.OffsetInDays.ValueInt64()) //nolint:gosec // offset_in_days is validated to -180..180 by the schema
		trigger.SetOffsetInDays(&offset)
		return trigger, nil
	case data.AttributeChange != nil:
		trigger := igmodels.NewAttributeChangeTrigger()
		attribute := igmodels.NewTriggerAttribute()
		convert.FrameworkToGraphString(data.AttributeChange.TriggerAttribute, attribute.SetName)
		trigger.SetTriggerAttributes([]igmodels.TriggerAttributeable{attribute})
		return trigger, nil
	case data.MembershipChange != nil:
		trigger := igmodels.NewMembershipChangeTrigger()
		if err := convert.FrameworkToGraphEnum(data.MembershipChange.ChangeType, igmodels.ParseMembershipChangeType, trigger.SetChangeType); err != nil {
			return nil, fmt.Errorf("failed to set membership change type: %w", err)
		}
		return trigger, nil
	}

	return nil, errTriggerUnset
}

// constructTasks builds the task list. executionSequence follows list order.
func constructTasks(data []WorkflowTask) []igmodels.Taskable {
	tasks := make([]igmodels.Taskable, 0, len(data))

	for i, item := range data {
		task := igmodels.NewTask()

		convert.FrameworkToGraphString(item.TaskDefinitionId, task.SetTaskDefinitionId)
		convert.FrameworkToGraphString(item.DisplayName, task.SetDisplayName)
		convert.FrameworkToGraphString(item.Description, task.SetDescription)
		convert.FrameworkToGraphBool(item.IsEnabled, task.SetIsEnabled)
		convert.FrameworkToGraphBool(item.ContinueOnError, task.SetContinueOnError)

		sequence := int32(i + 1)
		task.SetExecutionSequence(&sequence)

		// Graph rejects a task without an arguments property ("Task arguments missing or invalid"),
		// so unset arguments are sent as an empty array, never omitted and never null.
		arguments := make([]graphmodels.KeyValuePairable, 0, len(item.Arguments))
		for _, argument := range item.Arguments {
			pair := graphmodels.NewKeyValuePair()
			convert.FrameworkToGraphString(argument.Name, pair.SetName)
			convert.FrameworkToGraphString(argument.Value, pair.SetValue)
			arguments = append(arguments, pair)
		}
		task.SetArguments(arguments)

		tasks = append(tasks, task)
	}

	return tasks
}
