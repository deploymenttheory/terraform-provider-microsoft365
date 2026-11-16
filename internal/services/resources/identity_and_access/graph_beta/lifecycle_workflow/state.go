package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	igmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

const (
	ruleBasedSubjectSetODataType  = "#microsoft.graph.identityGovernance.ruleBasedSubjectSet"
	groupBasedSubjectSetODataType = "#microsoft.graph.identityGovernance.groupBasedSubjectSet"
)

// MapRemoteResourceStateToTerraform maps the remote lifecycle workflow from the Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data LifecycleWorkflowResourceModel, remoteResource igmodels.Workflowable) LifecycleWorkflowResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Category = convert.GraphToFrameworkEnum(remoteResource.GetCategory())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkStringWithDefault(remoteResource.GetDescription(), "")
	data.IsEnabled = convert.GraphToFrameworkBoolWithDefault(remoteResource.GetIsEnabled(), false)
	data.IsSchedulingEnabled = convert.GraphToFrameworkBoolWithDefault(remoteResource.GetIsSchedulingEnabled(), false)
	data.Version = convert.GraphToFrameworkInt32AsInt64(remoteResource.GetVersion())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	data.ExecutionConditions = stateExecutionConditions(ctx, remoteResource.GetExecutionConditions())
	data.Tasks = stateTasks(remoteResource.GetTasks(), data.Tasks)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}

// stateExecutionConditions maps the polymorphic executionConditions.
func stateExecutionConditions(ctx context.Context, remote igmodels.WorkflowExecutionConditionsable) *WorkflowExecutionConditions {
	if remote == nil {
		return nil
	}

	switch conditions := remote.(type) {
	case *igmodels.OnDemandExecutionOnly:
		return &WorkflowExecutionConditions{OnDemandExecutionOnly: types.BoolValue(true)}
	case *igmodels.TriggerAndScopeBasedConditions:
		return &WorkflowExecutionConditions{
			OnDemandExecutionOnly: types.BoolValue(false),
			Scope:                 stateScope(conditions),
			Trigger:               stateTrigger(conditions.GetTrigger()),
		}
	}

	tflog.Warn(ctx, fmt.Sprintf("Unsupported execution conditions type %v", remote.GetOdataType()))
	return nil
}

// stateScope maps the subject set. The SDK deserializes scope with the base models.SubjectSet
// discriminator, which does not know the identityGovernance subtypes, so rule and groups arrive in
// AdditionalData (*string, and []any of map[string]any) and the type is read from @odata.type.
func stateScope(conditions igmodels.TriggerAndScopeBasedConditionsable) *WorkflowScope {
	scope := conditions.GetScope()
	if scope == nil || scope.GetOdataType() == nil {
		return nil
	}
	additionalData := scope.GetAdditionalData()

	switch *scope.GetOdataType() {
	case ruleBasedSubjectSetODataType:
		rule := types.StringNull()
		if typed, ok := scope.(*igmodels.RuleBasedSubjectSet); ok {
			rule = convert.GraphToFrameworkString(typed.GetRule())
		} else if value, ok := additionalData["rule"].(*string); ok {
			rule = convert.GraphToFrameworkString(value)
		}
		return &WorkflowScope{Rule: rule, GroupIds: types.SetNull(types.StringType)}
	case groupBasedSubjectSetODataType:
		ids := make([]attr.Value, 0)
		if typed, ok := scope.(*igmodels.GroupBasedSubjectSet); ok {
			for _, group := range typed.GetGroups() {
				if group != nil && group.GetId() != nil {
					ids = append(ids, types.StringValue(*group.GetId()))
				}
			}
		} else if groups, ok := additionalData["groups"].([]any); ok {
			for _, item := range groups {
				group, _ := item.(map[string]any)
				if id, ok := group["id"].(*string); ok && id != nil {
					ids = append(ids, types.StringValue(*id))
				}
			}
		}
		return &WorkflowScope{Rule: types.StringNull(), GroupIds: types.SetValueMust(types.StringType, ids)}
	}
	return nil
}

func stateTrigger(remote igmodels.WorkflowExecutionTriggerable) *WorkflowTrigger {
	switch trigger := remote.(type) {
	case *igmodels.TimeBasedAttributeTrigger:
		return &WorkflowTrigger{
			TimeBasedAttribute: &TimeBasedAttributeTrigger{
				Attribute:    convert.GraphToFrameworkEnum(trigger.GetTimeBasedAttribute()),
				OffsetInDays: convert.GraphToFrameworkInt32AsInt64(trigger.GetOffsetInDays()),
			},
		}
	case *igmodels.AttributeChangeTrigger:
		attribute := types.StringNull()
		if attributes := trigger.GetTriggerAttributes(); len(attributes) > 0 && attributes[0] != nil {
			attribute = convert.GraphToFrameworkString(attributes[0].GetName())
		}
		return &WorkflowTrigger{AttributeChange: &AttributeChangeTrigger{TriggerAttribute: attribute}}
	case *igmodels.MembershipChangeTrigger:
		return &WorkflowTrigger{MembershipChange: &MembershipChangeTrigger{ChangeType: convert.GraphToFrameworkEnum(trigger.GetChangeType())}}
	}
	return nil
}

// stateTasks maps tasks in execution order. Graph returns arguments as [] when none were sent,
// so an empty argument list stays null unless the prior state at the same position held an
// explicit empty set.
func stateTasks(remote []igmodels.Taskable, prior []WorkflowTask) []WorkflowTask {
	ordered := make([]igmodels.Taskable, 0, len(remote))
	for _, task := range remote {
		if task != nil {
			ordered = append(ordered, task)
		}
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		return sequenceOf(ordered[i]) < sequenceOf(ordered[j])
	})

	tasks := make([]WorkflowTask, 0, len(ordered))
	for i, task := range ordered {
		item := WorkflowTask{
			TaskDefinitionId:  convert.GraphToFrameworkString(task.GetTaskDefinitionId()),
			DisplayName:       convert.GraphToFrameworkString(task.GetDisplayName()),
			Description:       convert.GraphToFrameworkStringWithDefault(task.GetDescription(), ""),
			IsEnabled:         convert.GraphToFrameworkBoolWithDefault(task.GetIsEnabled(), false),
			ContinueOnError:   convert.GraphToFrameworkBoolWithDefault(task.GetContinueOnError(), false),
			ExecutionSequence: convert.GraphToFrameworkInt32AsInt64(task.GetExecutionSequence()),
		}

		remoteArguments := task.GetArguments()
		switch {
		case len(remoteArguments) > 0:
			item.Arguments = make([]WorkflowTaskArgument, 0, len(remoteArguments))
			for _, argument := range remoteArguments {
				if argument == nil {
					continue
				}
				item.Arguments = append(item.Arguments, WorkflowTaskArgument{
					Name:  convert.GraphToFrameworkString(argument.GetName()),
					Value: convert.GraphToFrameworkString(argument.GetValue()),
				})
			}
		case i < len(prior) && prior[i].Arguments != nil:
			item.Arguments = []WorkflowTaskArgument{}
		}

		tasks = append(tasks, item)
	}

	return tasks
}

func sequenceOf(task igmodels.Taskable) int32 {
	if task.GetExecutionSequence() == nil {
		return 0
	}
	return *task.GetExecutionSequence()
}
